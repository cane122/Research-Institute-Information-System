package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cane/research-institute-system/backend/models"
	"github.com/cane/research-institute-system/backend/repositories"
	"github.com/cane/research-institute-system/backend/services"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/joho/godotenv"
	_ "github.com/sijms/go-ora/v2"
)

//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx              context.Context
	db               *sql.DB
	authService      *services.AuthService
	documentService  *services.DocumentService
	llmService       *services.LLMService
	analyticsService *services.AnalyticsService
	userRepo         *repositories.UserRepository
	projectRepo      *repositories.ProjectRepository
	currentUser      *models.User
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// maskPassword masks password for logging
func maskPassword(password string) string {
	if len(password) <= 2 {
		return "***"
	}
	return password[:2] + "***"
}

// testDatabaseConnections tries different connection configurations
func (a *App) testDatabaseConnections() {
	log.Printf("=== DIJAGNOSTIKA KONEKCIJE ===")

	// Test različitih konfiguracija using SID
	configs := []struct {
		name     string
		host     string
		port     string
		user     string
		password string
		sid      string
	}{
		{"Default XE SID", "localhost", "1521", "SYSTEM", "", "xe"},
		{"XE with password", "localhost", "1521", "SYSTEM", "oracle", "xe"},
		{"Alternative Port", "localhost", "1522", "SYSTEM", "", "xe"},
	}

	for _, config := range configs {
		log.Printf("Testiram %s: %s@%s:%s (SID=%s)", config.name, config.user, config.host, config.port, config.sid)

		// go-ora connection URL format: oracle://user:password@host:port/service_name
		// For SID, we use the SID as service name
		connStr := fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
			config.user, config.password, config.host, config.port, config.sid)

		db, err := sql.Open("oracle", connStr)
		if err != nil {
			log.Printf("  ❌ Open greška: %v", err)
			continue
		}

		err = db.Ping()
		db.Close()

		if err != nil {
			log.Printf("  ❌ Ping greška: %v", err)
		} else {
			log.Printf("  ✅ USPEŠNO!")
		}
	}

	log.Printf("=== KRAJ DIJAGNOSTIKE ===")
}

// OnStartup is called when the app starts up
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx

	log.Printf("=== POKRETANJE APLIKACIJE ===")

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Upozorenje: .env fajl nije pronađen (%v), koristim default/environment varijable", err)
	} else {
		log.Printf("✅ .env fajl uspešno učitan")
	}

	// Initialize database
	a.initializeDatabase()
}

// OnDomReady is called after the front-end dom has been loaded
func (a *App) OnDomReady(ctx context.Context) {
	// Optional: Initialize frontend-specific stuff
}

// OnShutdown is called when the app is terminating
func (a *App) OnShutdown(ctx context.Context) {
	// Cleanup database connections
	if a.db != nil {
		a.db.Close()
	}
}

// initializeDatabase initializes database connection
func (a *App) initializeDatabase() {
	// Initialize database connection with better error handling
	// Prioritet: .env file > environment variables > defaults
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "1521")
	dbUser := getEnvOrDefault("DB_USER", "SYSTEM")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "")
	dbSID := getEnvOrDefault("DB_SID", "xe")

	log.Printf("Pokušavam konekciju na bazu:")
	log.Printf("  Host: %s", dbHost)
	log.Printf("  Port: %s", dbPort)
	log.Printf("  User: %s", dbUser)
	log.Printf("  SID: %s", dbSID)
	log.Printf("  Password: %s", maskPassword(dbPassword))

	// go-ora connection URL format: oracle://user:password@host:port/service_name
	// For SID, we use the SID as service name
	connStr := fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		dbUser, dbPassword, dbHost, dbPort, dbSID)

	db, err := sql.Open("oracle", connStr)
	if err != nil {
		log.Printf("❌ GREŠKA: Failed to open database connection: %v", err)
		log.Printf("Application will continue without database. To configure database:")
		log.Printf("1. Install Oracle Database (XE or Standard)")
		log.Printf("2. Create user with appropriate privileges")
		log.Printf("3. Set environment variables: DB_HOST, DB_USER, DB_PASSWORD, DB_SID")
		log.Printf("4. Run the SQL schema from database/schema_oracle.sql")
		a.testDatabaseConnections()
		return
	}

	// Test the connection
	log.Printf("Testiram konekciju...")
	if err := db.Ping(); err != nil {
		log.Printf("❌ GREŠKA: Failed to ping database: %v", err)
		log.Printf("Database connection string (masked): %s:***@%s:%s:%s", dbUser, dbHost, dbPort, dbSID)
		log.Printf("Application will continue without database.")
		a.testDatabaseConnections()
		return
	}

	a.db = db
	log.Printf("✅ Successfully connected to Oracle database (SID: %s)", dbSID) // Initialize repositories
	a.userRepo = repositories.NewUserRepository(db)
	a.projectRepo = repositories.NewProjectRepository(db)

	// Initialize services
	a.authService = services.NewAuthService(a.userRepo)
	a.documentService = services.NewDocumentService(db)
	a.llmService = services.NewLLMService()
	a.analyticsService = services.NewAnalyticsService(db)

	// Check if OpenAI API key is configured
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		log.Printf("✅ OpenAI API key configured")
		a.llmService.SetAPIKey(apiKey)
	} else {
		log.Printf("⚠️  OpenAI API key not found. Set OPENAI_API_KEY environment variable to use LLM features.")
	}
}

// Login authenticates a user
func (a *App) Login(username, password string) (*services.LoginResponse, error) {
	if a.authService == nil {
		return &services.LoginResponse{
			Success: false,
			Message: "Sistem nije povezan sa bazom podataka",
		}, nil
	}

	response, err := a.authService.Login(services.LoginRequest{
		Username: username,
		Password: password,
	})

	if err == nil && response.Success {
		a.currentUser = response.User
	}

	return response, err
}

// Logout logs out the current user
func (a *App) Logout() {
	a.currentUser = nil
}

// GetCurrentUser returns the currently logged in user
func (a *App) GetCurrentUser() *models.User {
	return a.currentUser
}

// TestConnection tests if the backend is working
func (a *App) TestConnection() map[string]interface{} {
	result := make(map[string]interface{})

	result["backend_status"] = "ok"
	result["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	if a.db != nil {
		err := a.db.Ping()
		if err != nil {
			result["database_status"] = "error"
			result["database_error"] = err.Error()
		} else {
			result["database_status"] = "connected"

			// Test basic query
			var count int
			err = a.db.QueryRow("SELECT COUNT(*) FROM Korisnici").Scan(&count)
			if err != nil {
				result["query_test"] = "error: " + err.Error()
			} else {
				result["query_test"] = fmt.Sprintf("ok - %d users in database", count)
			}
		}
	} else {
		result["database_status"] = "not_connected"
	}

	if a.authService != nil {
		result["auth_service"] = "initialized"
	} else {
		result["auth_service"] = "not_initialized"
	}

	return result
}

// CreateUser creates a new user (Admin only)
func (a *App) CreateUser(user *models.User, tempPassword string) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za kreiranje korisnika")
	}

	if a.authService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.authService.CreateUser(user, tempPassword)
}

// CompleteFirstTimeSetup completes first-time login setup
func (a *App) CompleteFirstTimeSetup(username, newPassword string) map[string]interface{} {
	result := make(map[string]interface{})

	log.Printf("🔧 CompleteFirstTimeSetup pozvana za korisnika: %s", username)

	if a.authService == nil {
		log.Printf("❌ AuthService nije inicijalizovan")
		result["success"] = false
		result["message"] = "Sistem nije povezan sa bazom podataka"
		return result
	}

	log.Printf("🔍 Pozivam CompleteFirstTimeSetupByUsername...")
	err := a.authService.CompleteFirstTimeSetupByUsername(username, newPassword)
	if err != nil {
		log.Printf("❌ Greška u CompleteFirstTimeSetupByUsername: %v", err)
		result["success"] = false
		result["message"] = err.Error()
		return result
	}

	log.Printf("✅ CompleteFirstTimeSetup uspešno završena")
	result["success"] = true
	result["message"] = "Lozinka je uspešno postavljena"
	return result
}

// GetAllUsers returns all users (Admin only)
func (a *App) GetAllUsers() ([]models.User, error) {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return nil, errors.New("nemate dozvolu za pregled korisnika")
	}

	if a.userRepo == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.userRepo.GetAll()
}

// GetUserProjects returns projects for the current user
func (a *App) GetUserProjects() ([]models.Project, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.projectRepo == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.projectRepo.GetByUserID(a.currentUser.KorisnikID)
}

// CreateProject creates a new project
func (a *App) CreateProject(project *models.Project) error {
	if a.currentUser == nil || (a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator") {
		return errors.New("nemate dozvolu za kreiranje projekata")
	}

	if a.projectRepo == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	project.RukovodilaID = &a.currentUser.KorisnikID
	project.Status = "Aktivan"

	return a.projectRepo.Create(project)
}

// Document Management Methods

// GetAllDocuments returns all documents
func (a *App) GetAllDocuments() ([]models.Dokumenti, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.documentService.GetAllDocuments()
}

// GetDocumentsForUser returns only documents that the current user can access
func (a *App) GetDocumentsForUser() ([]models.Dokumenti, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.documentService.GetDocumentsForUser(a.currentUser.KorisnikID)
}

// GetAllTags returns all available tags from the database
func (a *App) GetAllTags() ([]models.Tag, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.documentService.GetAllTags()
}

// GetDocumentByID returns a specific document by ID
func (a *App) GetDocumentByID(documentID int) (models.Dokumenti, error) {
	if a.currentUser == nil {
		return models.Dokumenti{}, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return models.Dokumenti{}, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user has read permission
	hasPermission, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "read")
	if err != nil {
		return models.Dokumenti{}, fmt.Errorf("greška pri proveri dozvola: %w", err)
	}

	if !hasPermission {
		return models.Dokumenti{}, errors.New("nemate dozvolu za pregled ovog dokumenta")
	}

	return a.documentService.GetDocumentByID(documentID)
}

// GetDocumentVersions returns all versions of a document
func (a *App) GetDocumentVersions(documentID int) ([]models.VerzijeDokumenata, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user has read permission
	hasPermission, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "read")
	if err != nil {
		return nil, fmt.Errorf("greška pri proveri dozvola: %w", err)
	}

	if !hasPermission {
		return nil, errors.New("nemate dozvolu za pregled verzija ovog dokumenta")
	}

	return a.documentService.GetDocumentVersions(documentID)
}

// GetDocumentTags returns all tags for a document
func (a *App) GetDocumentTags(documentID int) ([]models.Tagovi, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user has read permission
	hasPermission, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "read")
	if err != nil {
		return nil, fmt.Errorf("greška pri proveri dozvola: %w", err)
	}

	if !hasPermission {
		return nil, errors.New("nemate dozvolu za pregled tagova ovog dokumenta")
	}

	return a.documentService.GetDocumentTags(documentID)
}

// UploadDocument uploads a new document
func (a *App) UploadDocument(req models.UploadDocumentRequest, fileData []byte, fileName string) (int, error) {
	if a.currentUser == nil {
		return 0, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.documentService.UploadDocument(req, fileData, fileName, a.currentUser.KorisnikID)
}

// UpdateDocument updates an existing document
func (a *App) UpdateDocument(documentID int, req models.UploadDocumentRequest) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user has write permission
	hasPermission, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "write")
	if err != nil {
		return fmt.Errorf("greška pri proveri dozvola: %w", err)
	}

	if !hasPermission {
		return errors.New("nemate dozvolu za izmenu ovog dokumenta")
	}

	return a.documentService.UpdateDocument(documentID, req)
}

// DeleteDocument deletes a document
func (a *App) DeleteDocument(documentID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user has delete permission
	hasPermission, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "delete")
	if err != nil {
		return fmt.Errorf("greška pri proveri dozvola: %w", err)
	}

	if !hasPermission {
		return errors.New("nemate dozvolu za brisanje ovog dokumenta")
	}

	return a.documentService.DeleteDocument(documentID)
}

// ============================================================================
// Document Permissions Management
// ============================================================================

// GetDocumentPermissions returns all permissions for a document
func (a *App) GetDocumentPermissions(documentID int) ([]models.DocumentPermissionResponse, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.documentService.GetDocumentPermissions(documentID)
}

// SetDocumentPermission sets or updates permission for a user on a document
func (a *App) SetDocumentPermission(req models.DocumentPermissionRequest) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.documentService.SetDocumentPermission(req)
}

// RemoveDocumentPermission removes a user's permission from a document
func (a *App) RemoveDocumentPermission(documentID int, userID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.documentService.RemoveDocumentPermission(documentID, userID)
}

// CheckUserPermission checks if current user has specific permission on a document
func (a *App) CheckUserPermission(documentID int, permissionType string) (bool, error) {
	if a.currentUser == nil {
		return false, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return false, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, permissionType)
}

// ============================================================================
// Analytics and Activity Log Operations
// ============================================================================

// LogActivity logs a user activity
func (a *App) LogActivity(req models.ActivityLogRequest) error {
	if a.analyticsService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	var korisnikID *int
	if a.currentUser != nil {
		korisnikID = &a.currentUser.KorisnikID
	}

	return a.analyticsService.LogActivity(korisnikID, req)
}

// GetRecentActivity retrieves recent activities
func (a *App) GetRecentActivity(limit int) ([]models.SkornjeAktivnosti, error) {
	if a.analyticsService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.analyticsService.GetRecentActivity(limit)
}

// GetDocumentStatistics retrieves document statistics
func (a *App) GetDocumentStatistics() (*models.StatistikaDokumenata, error) {
	if a.analyticsService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.analyticsService.GetDocumentStatistics()
}

// GetActivityStatistics retrieves activity statistics
func (a *App) GetActivityStatistics() ([]models.StatistikaAktivnosti, error) {
	if a.analyticsService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.analyticsService.GetActivityStatistics()
}

// GetDocumentsByType returns document count grouped by type
func (a *App) GetDocumentsByType() (map[string]int, error) {
	if a.analyticsService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.analyticsService.GetDocumentsByType()
}

// GetDocumentTrends returns document creation trends over time
func (a *App) GetDocumentTrends(days int) ([]map[string]interface{}, error) {
	if a.analyticsService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.analyticsService.GetDocumentTrends(days)
}

// GetTopContributors returns users with most document uploads
func (a *App) GetTopContributors(limit int) ([]map[string]interface{}, error) {
	if a.analyticsService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.analyticsService.GetTopContributors(limit)
}

// ============================================================================
// LLM Operations
// ============================================================================

// SetOpenAIKey sets the OpenAI API key for LLM operations
func (a *App) SetOpenAIKey(apiKey string) error {
	if a.llmService == nil {
		return errors.New("LLM servis nije inicijalizovan")
	}

	a.llmService.SetAPIKey(apiKey)

	// Optionally save to environment or config file
	os.Setenv("OPENAI_API_KEY", apiKey)

	return nil
}

// GenerateDocumentSummary generates an AI summary of a document
func (a *App) GenerateDocumentSummary(documentID int, maxLength int) (string, error) {
	if a.currentUser == nil {
		return "", errors.New("niste prijavljeni")
	}

	if a.llmService == nil {
		return "", errors.New("LLM servis nije inicijalizovan")
	}

	if a.documentService == nil {
		return "", errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get document content
	document, err := a.documentService.GetDocumentByID(documentID)
	if err != nil {
		return "", fmt.Errorf("greška pri učitavanju dokumenta: %w", err)
	}

	// For now, use document description and name as content
	// In a real implementation, you'd extract text from the actual file
	documentText := fmt.Sprintf("Document: %s\nDescription: %s",
		document.NazivDokumenta,
		func() string {
			if document.Opis != nil {
				return *document.Opis
			}
			return "No description"
		}())

	// Generate summary
	summary, err := a.llmService.GenerateSummary(documentText, maxLength)
	if err != nil {
		return "", fmt.Errorf("greška pri generisanju sažetka: %w", err)
	}

	// Save summary to database
	err = a.documentService.SaveLLMSummary(documentID, summary)
	if err != nil {
		log.Printf("Upozorenje: Nije moguće sačuvati sažetak u bazu: %v", err)
		// Continue anyway and return the summary
	}

	return summary, nil
}

// GenerateDocumentTags generates AI tags for a document
func (a *App) GenerateDocumentTags(documentID int, maxTags int) ([]string, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.llmService == nil {
		return nil, errors.New("LLM servis nije inicijalizovan")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get document content
	document, err := a.documentService.GetDocumentByID(documentID)
	if err != nil {
		return nil, fmt.Errorf("greška pri učitavanju dokumenta: %w", err)
	}

	// For now, use document description and name as content
	documentText := fmt.Sprintf("Document: %s\nDescription: %s",
		document.NazivDokumenta,
		func() string {
			if document.Opis != nil {
				return *document.Opis
			}
			return "No description"
		}())

	// Generate tags
	tags, err := a.llmService.GenerateTags(documentText, maxTags)
	if err != nil {
		return nil, fmt.Errorf("greška pri generisanju tagova: %w", err)
	}

	return tags, nil
}

// GenerateTagsFromText generates AI tags from document text (before upload)
func (a *App) GenerateTagsFromText(documentName string, description string, documentType string, maxTags int) ([]string, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.llmService == nil {
		return nil, errors.New("LLM servis nije inicijalizovan")
	}

	// Build document text from provided information
	documentText := fmt.Sprintf("Document Name: %s\nDocument Type: %s\nDescription: %s",
		documentName, documentType, description)

	// Generate tags
	tags, err := a.llmService.GenerateTags(documentText, maxTags)
	if err != nil {
		return nil, fmt.Errorf("greška pri generisanju tagova: %w", err)
	}

	return tags, nil
}

// GenerateDescriptionFromText generates AI description from document name and type (before upload)
func (a *App) GenerateDescriptionFromText(documentName string, documentType string, fileName string) (string, error) {
	if a.currentUser == nil {
		return "", errors.New("niste prijavljeni")
	}

	if a.llmService == nil {
		return "", errors.New("LLM servis nije inicijalizovan")
	}

	// Generate description using the specialized function
	description, err := a.llmService.GenerateDescription(documentName, documentType, fileName)
	if err != nil {
		return "", fmt.Errorf("greška pri generisanju opisa: %w", err)
	}

	return description, nil
}

// AskDocumentQuestion asks a question about a document using AI
func (a *App) AskDocumentQuestion(documentID int, question string) (string, error) {
	if a.currentUser == nil {
		return "", errors.New("niste prijavljeni")
	}

	if a.llmService == nil {
		return "", errors.New("LLM servis nije inicijalizovan")
	}

	if a.documentService == nil {
		return "", errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get document content
	document, err := a.documentService.GetDocumentByID(documentID)
	if err != nil {
		return "", fmt.Errorf("greška pri učitavanju dokumenta: %w", err)
	}

	// For now, use document description and name as content
	documentText := fmt.Sprintf("Document: %s\nDescription: %s",
		document.NazivDokumenta,
		func() string {
			if document.Opis != nil {
				return *document.Opis
			}
			return "No description"
		}())

	// Ask question
	answer, err := a.llmService.AnswerQuestion(documentText, question)
	if err != nil {
		return "", fmt.Errorf("greška pri postavljanju pitanja: %w", err)
	}

	return answer, nil
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Research Institute System",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.OnStartup,
		OnDomReady:       app.OnDomReady,
		OnShutdown:       app.OnShutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
