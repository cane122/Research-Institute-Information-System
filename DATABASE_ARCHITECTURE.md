# Database Architecture - Research Institute Information System

## Overview
This Go application implements a **three-tier architecture** for database operations with **Oracle Database XE 21c**, using the **Repository Pattern** and **Service Layer Pattern** for clean separation of concerns.

---

## 1. Database Connection Layer

### Technology Stack
- **Database**: Oracle Database XE 21c
- **Driver**: `github.com/sijms/go-ora/v2` (native Go Oracle driver)
- **Connection Pool**: Go's built-in `database/sql` package

### Connection Configuration (main.go)
```go
connStr := "oracle://SYSTEM:123@localhost:1521/xe"
db, err := sql.Open("oracle", connStr)
```

**Key Features**:
- Connection pooling automatically managed by `database/sql`
- Environment variables loaded via `godotenv` for configuration
- Connection test using `db.Ping()` before application starts
- Structured error handling and logging

---

## 2. Data Models Layer (`backend/models/models.go`)

### Purpose
Define Go structs that map to Oracle database tables using **struct tags** for JSON serialization and database mapping.

### Example Models

**Korisnici (Users)**:
```go
type Korisnici struct {
    KorisnikID        int        `json:"korisnik_id" db:"korisnik_id"`
    KorisnickoIme     string     `json:"korisnicko_ime"`
    Email             string     `json:"email"`
    HashSifre         string     `json:"-"`  // Hidden from JSON
    UlogaID           int        `json:"uloga_id"`
    Status            string     `json:"status"`
    PoslednajaPrijava *time.Time `json:"poslednja_prijava"`
}
```

**Dokumenti (Documents)**:
```go
type Dokumenti struct {
    DokumentID         int        `json:"dokument_id"`
    ProjekatID         *int       `json:"projekat_id"`
    NazivDokumenta     string     `json:"naziv_dokumenta"`
    TipDokumenta       string     `json:"tip_dokumenta"`
    KreiraoKorisnikID  int        `json:"kreirao_korisnik_id"`
    DatumaPostavke     time.Time  `json:"datuma_postavke"`
}
```

**Key Features**:
- Pointers (`*int`, `*time.Time`) for nullable fields
- JSON tags control frontend serialization
- Type aliases for English/Serbian compatibility

---

## 3. Repository Layer (`backend/repositories/`)

### Purpose
Encapsulates all direct SQL queries and database operations. Provides clean abstraction over data access.

### Structure
```
repositories/
├── user_repository.go       # User CRUD operations
└── project_repository.go    # Project CRUD operations
```

### Example: UserRepository (`user_repository.go`)

**Initialization**:
```go
type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}
```

**Query Example with Oracle-specific syntax**:
```go
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
    query := `
        SELECT k.korisnik_id, k.korisnicko_ime, k.email, k.hash_sifre, 
               k.uloga_id, k.status, u.naziv_uloge
        FROM Korisnici k
        JOIN Uloge u ON k.uloga_id = u.uloga_id
        WHERE k.korisnicko_ime = :1   -- Oracle positional binding
    `
    
    var user models.User
    err := r.db.QueryRow(query, username).Scan(
        &user.KorisnikID, &user.KorisnickoIme, &user.Email,
        &user.HashSifre, &user.UlogaID, &user.Status, &user.NazivUloge,
    )
    
    return &user, err
}
```

**Key Features**:
- **Oracle bind variables**: `:1, :2, :3` (positional parameters)
- **sql.NullTime** and **sql.NullString** for nullable columns
- **LEFT JOIN** for optional relationships
- Error propagation to service layer

---

## 4. Service Layer (`backend/services/`)

### Purpose
Contains business logic, orchestrates multiple repository calls, handles transactions, and enforces authorization rules.

### Structure
```
services/
├── auth_service.go         # Authentication, password hashing
├── document_service.go     # Document management, permissions
├── analitics_service.go    # Activity logging, statistics
├── llm_service.go          # AI integration
├── project_service.go      # Project management
├── task_service.go         # Task management
├── user_service.go         # User management
└── workflow_service.go     # Workflow management
```

### Example: DocumentService (`document_service.go`)

**Initialization**:
```go
type DocumentService struct {
    db         *sql.DB
    uploadPath string
}

func NewDocumentService(db *sql.DB) *DocumentService {
    cfg := config.LoadConfig()
    return &DocumentService{
        db:         db,
        uploadPath: cfg.UploadPath,
    }
}
```

**Complex Query with Multiple JOINs**:
```go
func (s *DocumentService) GetAllDocuments() ([]models.Dokumenti, error) {
    query := `
        SELECT d.dokument_id, d.projekat_id, d.naziv_dokumenta, 
               d.tip_dokumenta, d.kreirao_korisnik_id,
               COALESCE(p.naziv_projekta, '') as naziv_projekta,
               k.korisnicko_ime as ime_kreirao,
               COALESCE(v.version_count, 0) as broj_verzija
        FROM dokumenti d
        LEFT JOIN projekti p ON d.projekat_id = p.projekat_id
        JOIN korisnici k ON d.kreirao_korisnik_id = k.korisnik_id
        LEFT JOIN (
            SELECT dokument_id, COUNT(*) as version_count 
            FROM verzijedokumenata 
            GROUP BY dokument_id
        ) v ON d.dokument_id = v.dokument_id
        ORDER BY d.datuma_postavke DESC
    `
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var documents []models.Dokumenti
    for rows.Next() {
        var doc models.Dokumenti
        err := rows.Scan(/* ... */)
        documents = append(documents, doc)
    }
    
    return documents, nil
}
```

**Permission Check with Business Logic**:
```go
func (s *DocumentService) CheckUserPermission(documentID, userID int, permissionType string) (bool, error) {
    // Check if user is creator or admin
    checkQuery := `
        SELECT d.kreirao_korisnik_id, u.naziv_uloge
        FROM dokumenti d
        CROSS JOIN korisnici k
        LEFT JOIN uloge u ON k.uloga_id = u.uloga_id
        WHERE d.dokument_id = :1 AND k.korisnik_id = :2
    `
    
    var creatorID int
    var userRole sql.NullString
    err := s.db.QueryRow(checkQuery, documentID, userID).Scan(&creatorID, &userRole)
    
    // Admin or creator has full access
    if creatorID == userID || strings.ToLower(userRole.String) == "administrator" {
        return true, nil
    }
    
    // Check specific permissions in dozvoledokumenata table
    // ...
}
```

**Key Features**:
- **COALESCE()** for handling NULLs
- **Subqueries** for aggregations (version counts)
- **Transaction support** with `db.Begin()`
- **Business logic** (permission checks, role-based access)

---

## 5. Application Layer (main.go)

### App Structure
```go
type App struct {
    ctx              context.Context
    db               *sql.DB
    authService      *services.AuthService
    documentService  *services.DocumentService
    analyticsService *services.AnalyticsService
    userRepo         *repositories.UserRepository
    projectRepo      *repositories.ProjectRepository
    currentUser      *models.User
}
```

### Initialization Flow
```go
func main() {
    app := NewApp()
    
    // 1. Load environment variables
    godotenv.Load()
    
    // 2. Connect to Oracle database
    db, err := sql.Open("oracle", connStr)
    app.db = db
    
    // 3. Initialize repositories
    app.userRepo = repositories.NewUserRepository(db)
    app.projectRepo = repositories.NewProjectRepository(db)
    
    // 4. Initialize services (dependency injection)
    app.authService = services.NewAuthService(app.userRepo)
    app.documentService = services.NewDocumentService(db)
    app.analyticsService = services.NewAnalyticsService(db)
    
    // 5. Start Wails application
    wails.Run(&options.App{
        OnStartup: app.startup,
        Bind: []interface{}{app},
    })
}
```

### Frontend-Backend Bridge
Methods in `App` are **automatically exposed** to frontend via Wails:

```go
// Frontend calls: GetDocumentsForUser()
func (a *App) GetDocumentsForUser() ([]models.Dokumenti, error) {
    userID := a.currentUser.KorisnikID
    return a.documentService.GetDocumentsForUser(userID)
}

// Frontend calls: CheckUserPermission(docID, "read")
func (a *App) CheckUserPermission(documentID int, permissionType string) (bool, error) {
    return a.documentService.CheckUserPermission(
        documentID, 
        a.currentUser.KorisnikID, 
        permissionType,
    )
}
```

---

## 6. Oracle-Specific Implementation Details

### Parameter Binding
- **Oracle style**: `:1, :2, :3` (positional)
- **Standard SQL**: `$1, $2, $3` (PostgreSQL style not supported)

### Null Handling
```go
var lastLogin sql.NullTime
var description sql.NullString

// After scan:
if lastLogin.Valid {
    user.PoslednajaPrijava = &lastLogin.Time
}
```

### INSERT with Auto-increment
Oracle uses `GENERATED ALWAYS AS IDENTITY`:
```sql
INSERT INTO dokumenti (naziv_dokumenta, kreirao_korisnik_id) 
VALUES (:1, :2)
```

### Activity Logging
```go
// Insert into LogAktivnosti table
query := `
    INSERT INTO logaktivnosti (
        korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis
    ) VALUES (:1, :2, :3, :4, :5)
`
_, err := s.db.Exec(query, korisnikID, "LOGIN", "KORISNIK", userID, "User logged in")
```

---

## 7. Key Database Tables

### Core Tables
- **Korisnici** - Users (14 users: admin, managers, researchers)
- **Uloge** - Roles (Administrator, Manager, Researcher, Organizer)
- **Projekti** - Projects (research projects)
- **Dokumenti** - Documents (PDFs, DOCX files)
- **Zadaci** - Tasks
- **LogAktivnosti** - Activity logs (LOGIN, VIEW, EDIT, DELETE, etc.)
- **DozvoleDokumenata** - Document permissions (read, write, delete)

### Relationships
- `Korisnici` → `Uloge` (many-to-one)
- `Dokumenti` → `Projekti` (many-to-one, nullable)
- `Dokumenti` → `Korisnici` (creator, many-to-one)
- `DozvoleDokumenata` → `Dokumenti` + `Korisnici` (many-to-many with permissions)

---

## 8. Transaction Handling

```go
func (s *DocumentService) DeleteDocument(documentID int) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Auto-rollback if not committed
    
    // Delete related records
    _, err = tx.Exec("DELETE FROM verzijedokumenata WHERE dokument_id = :1", documentID)
    if err != nil {
        return err
    }
    
    _, err = tx.Exec("DELETE FROM dokumenti WHERE dokument_id = :1", documentID)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}
```

---

## Summary

**Architecture Pattern**: Three-tier (Presentation → Service → Repository → Database)

**Database**: Oracle XE 21c with `go-ora/v2` native driver

**Key Components**:
1. **Models** - Struct definitions with tags
2. **Repositories** - SQL queries and data access
3. **Services** - Business logic and orchestration
4. **App** - Dependency injection and frontend binding

**Special Features**:
- Oracle-specific positional binding (`:1, :2`)
- Null-safe handling with `sql.NullTime`, `sql.NullString`
- Role-based access control (RBAC)
- Activity logging for audit trail
- Transaction support for data consistency

**Performance**:
- Connection pooling
- Prepared statement reuse
- Efficient JOINs and subqueries
- Index usage on foreign keys

This architecture provides **clean separation of concerns**, **testability**, and **maintainability** while leveraging Oracle Database's features effectively.
