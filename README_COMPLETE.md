# Research Institute Information System
## Istraživačko razvojni centar - Sistem za upravljanje informacijama

A comprehensive information management system built with Wails, Go, and PostgreSQL for research and development institutes.

### 🎯 Project Overview

This desktop application provides a complete solution for research institutes, supporting three main subsystems:

1. **Document Management** (Upravljanje dokumentima)
2. **Project Documentation Preparation** (Priprema projektne dokumentacije) 
3. **Project Realization** (Realizacija projekata)

### ✨ Key Features

#### 🔐 User Management & Roles
- **Administrator**: Full system access and user management
- **Project Manager** (Rukovodilac): Project oversight and team management
- **Researcher** (Istraživač): Document and task management

#### 📁 Project Management
- Complete project lifecycle tracking
- Team member assignment and collaboration
- Progress monitoring and reporting
- Customizable workflow phases

#### 📋 Task Management
- Kanban-style task boards
- Priority and deadline management
- Phase change requests and approvals
- Task comments and collaboration

#### 📄 Document Management
- Hierarchical folder organization
- Document versioning and history
- AI-powered document summaries (LLM integration)
- Flexible metadata and tagging system
- Permission-based access control

#### 📊 Analytics & Reporting
- Comprehensive activity logging
- Dashboard with key metrics
- Project and user statistics
- System audit trails

### 🛠 Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML5, CSS3, JavaScript
- **Database**: PostgreSQL with full schema
- **Framework**: Wails v2 (Desktop Application)
- **Architecture**: Clean Architecture with separated concerns

### 📁 Project Structure

```
├── app/                    # Wails application context
├── backend/               
│   ├── models/            # Database models and DTOs
│   ├── repositories/      # Data access layer
│   └── services/          # Business logic layer
├── database/
│   └── schema.sql         # Complete PostgreSQL schema
├── frontend/              # Web frontend assets
├── wireframes/            # UI/UX wireframes
│   ├── index.html         # Wireframe showcase
│   ├── 01-login.html      # Login interface
│   ├── 02-dashboard.html  # Main dashboard
│   ├── 03-projects.html   # Project management
│   ├── 04-documents.html  # Document management  
│   ├── 05-tasks-kanban.html # Task kanban board
│   └── 06-users-admin.html  # User administration
└── build/                 # Build artifacts

```

### 🗄️ Database Schema

The system uses a comprehensive PostgreSQL schema with 4 main modules:

#### Module 1: User & Role Management
- `Uloge` - User roles definition
- `Korisnici` - System users with authentication

#### Module 2: Project, Task & Documentation Management  
- `RadniTokovi` - Workflow definitions
- `Faze` - Workflow phases
- `Projekti` - Research projects
- `ClanoviProjekta` - Project team members
- `Zadaci` - Tasks within projects
- `KomentariZadataka` - Task comments
- `ZahteviPromeneFaze` - Phase change requests

#### Module 3: Document & Metadata Management
- `Folderi` - Hierarchical folder structure
- `Dokumenti` - Document registry
- `VerzijeDokumenata` - Document versions
- `LLMSazeci` - AI-generated summaries
- `MetaPodaci` - Flexible metadata
- `Tagovi` - Document tagging system
- `DokumentTagovi` - Document-tag relationships
- `DozvoleDokumenata` - Access permissions
- `IstorijaFazaDokumenta` - Document phase history

#### Module 4: Analytics & Logging
- `LogAktivnosti` - Comprehensive activity logging

### 🎨 Wireframes

The project includes detailed wireframes showcasing all major interfaces:

- **Login System**: Secure authentication interface
- **Dashboard**: Overview of projects, tasks, and activities  
- **Project Management**: Grid view with filtering and team management
- **Document Library**: Hierarchical organization with search capabilities
- **Kanban Board**: Agile task management with drag-drop functionality
- **User Administration**: Complete user and role management

View wireframes: Open `wireframes/index.html` in your browser

### 🚀 Getting Started

#### Prerequisites
- Go 1.19 or higher
- Node.js 16 or higher  
- PostgreSQL 12 or higher
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

#### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/cane122/Research-Institute-Information-System.git
   cd Research-Institute-Information-System
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup PostgreSQL database**
   ```sql
   CREATE DATABASE institute_db;
   \c institute_db;
   \i database/schema.sql
   ```

4. **Configure environment**
   ```bash
   # Create .env file with database connection
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_user
   DB_PASSWORD=your_password
   DB_NAME=institute_db
   ```

5. **Build and run**
   ```bash
   # Development mode
   wails dev
   
   # Production build
   wails build
   ```

### 🧪 Development

#### Project Structure
- **Clean Architecture**: Separated concerns with models, repositories, services
- **Serbian Language**: All database fields and API responses in Serbian
- **Role-based Access**: Granular permissions based on user roles
- **Audit Trail**: Comprehensive logging of all system activities

#### Key Models
```go
// User with Serbian field names
type Korisnici struct {
    KorisnikID       int        `json:"korisnik_id" db:"korisnik_id"`
    KorisnickoIme    string     `json:"korisnicko_ime" db:"korisnicko_ime"`
    Email            string     `json:"email" db:"email"`
    // ... other fields
}

// Project management
type Projekti struct {
    ProjekatID      int        `json:"projekat_id" db:"projekat_id"`
    NazivProjekta   string     `json:"naziv_projekta" db:"naziv_projekta"`
    // ... other fields
}
```

### 📊 System Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend       │    │   Database      │
│                 │    │                  │    │                 │
│ • HTML/CSS/JS   │◄──►│ • Go Services    │◄──►│ • PostgreSQL    │
│ • Wails Runtime │    │ • Repositories   │    │ • Full Schema   │
│ • Reactive UI   │    │ • Authentication │    │ • Indexes       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 📝 API Overview

The system provides RESTful APIs for all major operations:

- **Authentication**: Login, logout, session management
- **Users**: CRUD operations, role management  
- **Projects**: Project lifecycle, team management
- **Tasks**: Kanban operations, phase changes
- **Documents**: Upload, versioning, permissions
- **Analytics**: Statistics, activity logs

### 🔐 Security Features

- **Authentication**: Secure password hashing with bcrypt
- **Authorization**: Role-based access control (RBAC)
- **Session Management**: Secure session handling
- **Audit Logging**: All actions are logged with user attribution
- **Data Validation**: Input validation and sanitization
- **SQL Injection Protection**: Parameterized queries

### 📈 Performance Optimizations

- **Database Indexing**: Strategic indexes for fast queries
- **Connection Pooling**: Efficient database connection management
- **Caching**: In-memory caching for frequent operations
- **Pagination**: Efficient handling of large datasets
- **Lazy Loading**: Load data only when needed

### 🧪 Testing Strategy

- **Unit Tests**: Individual component testing
- **Integration Tests**: Database and API testing  
- **End-to-End Tests**: Complete user workflow testing
- **Performance Tests**: Load and stress testing
- **Security Tests**: Vulnerability and penetration testing

### 📝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### 🐛 Bug Reports

When reporting bugs, please include:
- Operating system and version
- Go version
- Steps to reproduce the issue
- Expected vs actual behavior
- Screenshots if applicable

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### 👥 Authors

- **Cane122** - *Initial work* - [cane122](https://github.com/cane122)

### 🙏 Acknowledgments

- Wails framework for excellent Go-to-desktop development
- PostgreSQL for robust data management
- Research community for requirements and feedback

---

**Note**: This system is specifically designed for Serbian research institutes with Serbian language support throughout the interface and database schema.

### 🔗 Related Links

- [Wails Documentation](https://wails.io/)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Project Wireframes](./wireframes/index.html)
