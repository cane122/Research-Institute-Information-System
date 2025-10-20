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
	_ "github.com/lib/pq"
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
	projectService   *services.ProjectService
	taskService      *services.TaskService
	workflowService  *services.WorkflowService
	conditionService *services.ConditionService
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

	// Test različitih konfiguracija
	configs := []struct {
		name     string
		host     string
		user     string
		password string
		dbname   string
	}{
		{"Default", "localhost:5432", "postgres", "password", "research_institute"},
		{"Alternative Password", "localhost:5432", "postgres", "postgres", "research_institute"},
		{"Different Port", "localhost:5433", "postgres", "password", "research_institute"},
		{"System DB", "localhost:5432", "postgres", "password", "postgres"},
	}

	for _, config := range configs {
		log.Printf("Testiram %s: %s@%s/%s", config.name, config.user, config.host, config.dbname)

		connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			config.user, config.password, config.host, config.dbname)

		db, err := sql.Open("postgres", connStr)
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
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "123")
	dbName := getEnvOrDefault("DB_NAME", "research_institute")

	log.Printf("Pokušavam konekciju na bazu:")
	log.Printf("  Host: %s", dbHost)
	log.Printf("  Port: %s", dbPort)
	log.Printf("  User: %s", dbUser)
	log.Printf("  Database: %s", dbName)
	log.Printf("  Password: %s", maskPassword(dbPassword))

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("❌ GREŠKA: Failed to open database connection: %v", err)
		log.Printf("Application will continue without database. To configure database:")
		log.Printf("1. Install PostgreSQL")
		log.Printf("2. Create database 'research_institute'")
		log.Printf("3. Set environment variables: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME")
		log.Printf("4. Run the SQL schema from database/schema.sql")
		a.testDatabaseConnections()
		return
	}

	// Test the connection
	log.Printf("Testiram konekciju...")
	if err := db.Ping(); err != nil {
		log.Printf("❌ GREŠKA: Failed to ping database: %v", err)
		log.Printf("Database connection string (masked): postgres://%s:***@%s/%s?sslmode=disable", dbUser, dbHost, dbName)
		log.Printf("Application will continue without database.")
		a.testDatabaseConnections()
		return
	}

	a.db = db
	log.Printf("Successfully connected to PostgreSQL database: %s", dbName) // Initialize repositories
	a.userRepo = repositories.NewUserRepository(db)
	a.projectRepo = repositories.NewProjectRepository(db)

	// Initialize services
	a.authService = services.NewAuthService(a.userRepo)
	a.documentService = services.NewDocumentService(db)
	a.llmService = services.NewLLMService()
	a.analyticsService = services.NewAnalyticsService(db)
	a.projectService = services.NewProjectService(db)
	a.taskService = services.NewTaskService(db)
	a.workflowService = services.NewWorkflowService(db)
	a.conditionService = services.NewConditionService(db)

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

// GetAvailableTeamMembers returns all researchers that can be team members
func (a *App) GetAvailableTeamMembers() ([]models.User, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	// Only Project Managers and Administrators can view potential team members
	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return nil, errors.New("nemate dozvolu za pregled članova tima")
	}

	if a.userRepo == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get all users and filter for researchers only
	allUsers, err := a.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Filter for researchers only
	var researchers []models.User
	for _, user := range allUsers {
		if user.NazivUloge == "Istrazivac" {
			researchers = append(researchers, user)
		}
	}

	return researchers, nil
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

// ============================================================================
// Project Realization Subsystem - Project Management
// ============================================================================

// GetAllProjects returns all projects in the system
func (a *App) GetAllProjects() ([]models.Projekti, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.projectService.GetAllProjects()
}

// GetProjectByID returns a specific project by ID
func (a *App) GetProjectByID(projectID int) (models.Projekti, error) {
	if a.currentUser == nil {
		return models.Projekti{}, errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return models.Projekti{}, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.projectService.GetProjectByID(projectID)
}

// CreateNewProject creates a new project with workflow and team members
func (a *App) CreateNewProject(req models.CreateProjectRequest) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za kreiranje projekata")
	}

	if a.projectService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Create project with current user as project manager
	return a.projectService.CreateProjectWithManager(req, a.currentUser.KorisnikID)
}

// UpdateProject updates an existing project
func (a *App) UpdateProject(projectID int, project models.Projekti) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user is project manager
	existingProject, err := a.projectService.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	if existingProject.RukovodilaID != nil && *existingProject.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("samo rukovodilac projekta može ažurirati projekat")
	}

	return a.projectService.UpdateProject(projectID, project)
}

// DeleteProject deletes a project
func (a *App) DeleteProject(projectID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user is project manager
	existingProject, err := a.projectService.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	if existingProject.RukovodilaID != nil && *existingProject.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("samo rukovodilac projekta može obrisati projekat")
	}

	return a.projectService.DeleteProject(projectID)
}

// GetProjectMembers returns all team members of a project
func (a *App) GetProjectMembers(projectID int) ([]models.Korisnici, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.projectService.GetProjectMembers(projectID)
}

// AddProjectMember adds a user to a project team
func (a *App) AddProjectMember(projectID, userID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user is project manager
	existingProject, err := a.projectService.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	if existingProject.RukovodilaID != nil && *existingProject.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("samo rukovodilac projekta može dodavati članove")
	}

	return a.projectService.AddProjectMember(projectID, userID)
}

// RemoveProjectMember removes a user from a project team
func (a *App) RemoveProjectMember(projectID, userID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user is project manager
	existingProject, err := a.projectService.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	if existingProject.RukovodilaID != nil && *existingProject.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("samo rukovodilac projekta može ukloniti članove")
	}

	return a.projectService.RemoveProjectMember(projectID, userID)
}

// GetProjectsByStatus returns projects filtered by status
func (a *App) GetProjectsByStatus(status string) ([]models.Projekti, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.projectService.GetProjectsByStatus(status)
}

// CompleteProject marks a project as completed
func (a *App) CompleteProject(projectID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user is project manager
	existingProject, err := a.projectService.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	if existingProject.RukovodilaID != nil && *existingProject.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("samo rukovodilac projekta može završiti projekat")
	}

	return a.projectService.CompleteProject(projectID)
}

// GetProjectResources returns resources allocated to a project
func (a *App) GetProjectResources(projectID int) ([]map[string]interface{}, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.projectService.GetProjectResources(projectID)
}

// GetProjectAnalytics returns comprehensive analytics for a project
func (a *App) GetProjectAnalytics(projectID int) (map[string]interface{}, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.projectService.GetProjectAnalytics(projectID)
}

// GetProjectsByCurrentUser returns all projects where current user is member or manager
func (a *App) GetProjectsByCurrentUser() ([]models.Projekti, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.projectService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.projectService.GetProjectsByUser(a.currentUser.KorisnikID)
}

// ============================================================================
// Project Realization Subsystem - Task Management
// ============================================================================

// GetTasksByProject returns all tasks for a specific project
func (a *App) GetTasksByProject(projectID int) ([]models.Zadaci, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetTasksByProject(projectID)
}

// GetTasksByUser returns all tasks assigned to a specific user
func (a *App) GetTasksByUser(userID int) ([]models.Zadaci, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetTasksByUser(userID)
}

// GetTasksByCurrentUser returns all tasks assigned to the current user
func (a *App) GetTasksByCurrentUser() ([]models.Zadaci, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetTasksByUser(a.currentUser.KorisnikID)
}

// GetTaskByID returns a specific task by ID
func (a *App) GetTaskByID(taskID int) (models.Zadaci, error) {
	if a.currentUser == nil {
		return models.Zadaci{}, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return models.Zadaci{}, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetTaskByID(taskID)
}

// CreateTask creates a new task
func (a *App) CreateTask(req models.CreateTaskRequest) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user is project manager
	if a.projectService != nil {
		project, err := a.projectService.GetProjectByID(req.ProjekatID)
		if err != nil {
			return err
		}

		if project.RukovodilaID != nil && *project.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
			return errors.New("samo rukovodilac projekta može kreirati zadatke")
		}
	}

	return a.taskService.CreateTask(req)
}

// UpdateTask updates an existing task
func (a *App) UpdateTask(taskID int, req models.UpdateTaskRequest) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get task to check project
	task, err := a.taskService.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// Check if user is project manager
	if a.projectService != nil {
		project, err := a.projectService.GetProjectByID(task.ProjekatID)
		if err != nil {
			return err
		}

		if project.RukovodilaID != nil && *project.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
			// Allow task owner to update their own task
			if task.DodjeljenKorisnikuID == nil || *task.DodjeljenKorisnikuID != a.currentUser.KorisnikID {
				return errors.New("nemate dozvolu za izmenu ovog zadatka")
			}
		}
	}

	return a.taskService.UpdateTask(taskID, req)
}

// DeleteTask deletes a task
func (a *App) DeleteTask(taskID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get task to check project
	task, err := a.taskService.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// Check if user is project manager
	if a.projectService != nil {
		project, err := a.projectService.GetProjectByID(task.ProjekatID)
		if err != nil {
			return err
		}

		if project.RukovodilaID != nil && *project.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
			return errors.New("samo rukovodilac projekta može obrisati zadatke")
		}
	}

	return a.taskService.DeleteTask(taskID)
}

// GetTaskComments returns all comments for a task
func (a *App) GetTaskComments(taskID int) ([]models.KomentariZadataka, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetTaskComments(taskID)
}

// AddTaskComment adds a comment to a task
func (a *App) AddTaskComment(taskID int, comment string) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.AddTaskComment(taskID, a.currentUser.KorisnikID, comment)
}

// MoveTaskToPhase moves a task to a different phase
func (a *App) MoveTaskToPhase(taskID, newPhaseID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get task to check project
	task, err := a.taskService.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// Check if user is project manager
	if a.projectService != nil {
		project, err := a.projectService.GetProjectByID(task.ProjekatID)
		if err != nil {
			return err
		}

		if project.RukovodilaID != nil && *project.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
			return errors.New("samo rukovodilac projekta može menjati fazu zadatka")
		}
	}

	return a.taskService.MoveTaskToPhase(taskID, newPhaseID)
}

// GetTasksByPhase returns all tasks in a specific phase
func (a *App) GetTasksByPhase(phaseID int) ([]models.Zadaci, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetTasksByPhase(phaseID)
}

// RequestPhaseChange creates a phase change request
func (a *App) RequestPhaseChange(taskID, requestedPhaseID int, comment string) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if there's already a pending request for this task
	hasPending, err := a.HasPendingPhaseChangeRequest(taskID)
	if err != nil {
		return err
	}
	if hasPending {
		return errors.New("već postoji aktivni zahtev za promenu faze za ovaj zadatak")
	}

	return a.taskService.RequestPhaseChange(taskID, a.currentUser.KorisnikID, requestedPhaseID, comment)
}

// HasPendingPhaseChangeRequest checks if there's a pending phase change request for a task
func (a *App) HasPendingPhaseChangeRequest(taskID int) (bool, error) {
	if a.currentUser == nil {
		return false, errors.New("niste prijavljeni")
	}

	if a.db == nil {
		return false, errors.New("sistem nije povezan sa bazom podataka")
	}

	var count int
	err := a.db.QueryRow(`
		SELECT COUNT(*) 
		FROM zahtevipromenefaze 
		WHERE zadatak_id = $1 AND status = 'na čekanju'
	`, taskID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetNextPhaseForTask returns the next phase in sequence for a task
func (a *App) GetNextPhaseForTask(taskID int) (*models.Faze, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil || a.workflowService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get task to find current phase
	task, err := a.taskService.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	// Get workflow for the task's project
	var workflowID int
	err = a.db.QueryRow(`SELECT radni_tok_id FROM projekti WHERE projekat_id = $1`, task.ProjekatID).Scan(&workflowID)
	if err != nil {
		return nil, err
	}

	// Get all phases for workflow ordered by sequence
	phases, err := a.workflowService.GetPhasesByWorkflow(workflowID)
	if err != nil {
		return nil, err
	}

	// Find current phase and return next one
	for i, phase := range phases {
		if phase.FazaID == task.FazaID {
			if i+1 < len(phases) {
				return &phases[i+1], nil
			}
			return nil, errors.New("zadatak je već u poslednjoj fazi")
		}
	}

	return nil, errors.New("trenutna faza nije pronađena u radnom toku")
}

// CheckTaskConditionsFulfilled checks if all conditions for current phase are fulfilled
func (a *App) CheckTaskConditionsFulfilled(taskID int) (bool, error) {
	if a.currentUser == nil {
		return false, errors.New("niste prijavljeni")
	}

	if a.taskService == nil || a.conditionService == nil {
		return false, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get task to find current phase
	task, err := a.taskService.GetTaskByID(taskID)
	if err != nil {
		return false, err
	}

	// Check if all conditions are fulfilled
	return a.conditionService.CheckAllConditionsFulfilled(taskID, task.FazaID)
}

// GetPhaseChangeRequestsForProject returns all phase change requests for a project
func (a *App) GetPhaseChangeRequestsForProject(projectID int) ([]models.ZahteviPromeneFaze, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetPhaseChangeRequests(&projectID, nil)
}

// GetPhaseChangeRequestsForTask returns all phase change requests for a task
func (a *App) GetPhaseChangeRequestsForTask(taskID int) ([]models.ZahteviPromeneFaze, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetPhaseChangeRequests(nil, &taskID)
}

// ApprovePhaseChangeRequest approves a phase change request
func (a *App) ApprovePhaseChangeRequest(requestID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Only managers can approve
	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("samo rukovodilac projekta može odobriti zahtev")
	}

	return a.taskService.ApprovePhaseChangeRequest(requestID)
}

// RejectPhaseChangeRequest rejects a phase change request
func (a *App) RejectPhaseChangeRequest(requestID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Only managers can reject
	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("samo rukovodilac projekta može odbiti zahtev")
	}

	return a.taskService.RejectPhaseChangeRequest(requestID)
}

// GetManagerPhaseChangeRequests returns all phase change requests for projects managed by the current user
func (a *App) GetManagerPhaseChangeRequests() ([]map[string]interface{}, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil || a.projectRepo == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Only managers can view
	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return nil, errors.New("samo rukovodilac projekta može videti zahteve")
	}

	fmt.Printf("🔍 Loading requests for user: %d (%s)\n", a.currentUser.KorisnikID, a.currentUser.KorisnickoIme)

	// Get projects where user is the leader
	projects, err := a.projectRepo.GetByUserID(a.currentUser.KorisnikID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("📁 Found %d projects for user\n", len(projects))

	var allRequests []map[string]interface{}

	// Get phase change requests for each project
	for _, project := range projects {
		fmt.Printf("  📂 Checking project %d (%s), Leader ID: %v\n", project.ProjekatID, project.NazivProjekta, project.RukovodilaID)
		
		// Only include projects where user is the leader
		if project.RukovodilaID == nil || *project.RukovodilaID != a.currentUser.KorisnikID {
			fmt.Printf("    ❌ Skipping - user is not the leader\n")
			continue
		}

		fmt.Printf("    ✅ User is the leader - checking requests\n")

		requests, err := a.taskService.GetPhaseChangeRequests(&project.ProjekatID, nil)
		if err != nil {
			fmt.Printf("    ⚠️  Error getting requests: %v\n", err)
			continue
		}

		fmt.Printf("    📋 Found %d requests for this project\n", len(requests))

		// Enrich requests with additional info
		for _, request := range requests {
			fmt.Printf("      🔸 Request %d - Status: %s\n", request.ZahtevID, request.Status)
			
			// Only include pending requests (check both variants)
			if request.Status != "Na cekanju" && request.Status != "na čekanju" {
				fmt.Printf("        ⏭️  Skipping - not pending (status: %s)\n", request.Status)
				continue
			}

			// Get task details
			task, err := a.taskService.GetTaskByID(request.ZadatakID)
			if err != nil {
				fmt.Printf("        ⚠️  Error getting task: %v\n", err)
				continue
			}

			// Get phase details
			var phaseName string
			err = a.db.QueryRow("SELECT naziv_faze FROM faze WHERE faza_id = $1", request.ZahtevanaFazaID).Scan(&phaseName)
			if err != nil {
				phaseName = "Nepoznata faza"
			}

			// Get submitter details
			var submitterName string
			err = a.db.QueryRow("SELECT COALESCE(ime || ' ' || prezime, korisnicko_ime) FROM korisnici WHERE korisnik_id = $1", request.PodnosilacZahtevaID).Scan(&submitterName)
			if err != nil {
				submitterName = "Nepoznat korisnik"
			}

			enrichedRequest := map[string]interface{}{
				"zahtev_id":             request.ZahtevID,
				"zadatak_id":            request.ZadatakID,
				"naziv_zadatka":         task.NazivZadatka,
				"projekat_id":           project.ProjekatID,
				"naziv_projekta":        project.NazivProjekta,
				"podnosilac_zahteva_id": request.PodnosilacZahtevaID,
				"podnosilac_ime":        submitterName,
				"zahtevana_faza_id":     request.ZahtevanaFazaID,
				"naziv_faze":            phaseName,
				"status":                request.Status,
				"komentar":              request.Komentar,
				"datum_kreiranja":       request.DatumKreiranja,
			}

			fmt.Printf("        ✅ Added request: %s\n", task.NazivZadatka)
			allRequests = append(allRequests, enrichedRequest)
		}
	}

	fmt.Printf("📊 Total requests to return: %d\n", len(allRequests))
	return allRequests, nil
}

// GetOverdueTasksForProject returns overdue tasks for a project
func (a *App) GetOverdueTasksForProject(projectID int) ([]models.Zadaci, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetOverdueTasks(&projectID)
}

// GetAllOverdueTasks returns all overdue tasks
func (a *App) GetAllOverdueTasks() ([]models.Zadaci, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.taskService.GetOverdueTasks(nil)
}

// UpdateTaskProgress updates task progress
func (a *App) UpdateTaskProgress(taskID, progress int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.taskService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get task to check if user is assigned
	task, err := a.taskService.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// Check if user is assigned to task or is manager
	if task.DodjeljenKorisnikuID != nil && *task.DodjeljenKorisnikuID != a.currentUser.KorisnikID {
		if a.projectService != nil {
			project, err := a.projectService.GetProjectByID(task.ProjekatID)
			if err != nil {
				return err
			}

			if project.RukovodilaID != nil && *project.RukovodilaID != a.currentUser.KorisnikID && a.currentUser.NazivUloge != "Administrator" {
				return errors.New("nemate dozvolu za ažuriranje progresa ovog zadatka")
			}
		}
	}

	return a.taskService.UpdateTaskProgress(taskID, progress)
}

// ============================================================================
// Project Realization Subsystem - Workflow Management
// ============================================================================

// GetAllWorkflows returns all workflows
func (a *App) GetAllWorkflows() ([]models.RadniTokovi, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.workflowService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.GetAllWorkflows()
}

// GetWorkflowByID returns a specific workflow by ID
func (a *App) GetWorkflowByID(workflowID int) (*models.RadniTokovi, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.workflowService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.GetWorkflowByID(workflowID)
}

// GetWorkflowPhases returns all phases of a workflow
func (a *App) GetWorkflowPhases(workflowID int) ([]models.Faze, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.workflowService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.GetWorkflowPhases(workflowID)
}

// CreateWorkflow creates a new workflow
func (a *App) CreateWorkflow(workflow models.RadniTokovi) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za kreiranje radnih tokova")
	}

	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.CreateWorkflow(workflow)
}

// UpdateWorkflow updates an existing workflow
func (a *App) UpdateWorkflow(workflowID int, workflow models.RadniTokovi) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za izmenu radnih tokova")
	}

	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.UpdateWorkflow(workflowID, workflow)
}

// DeleteWorkflow deletes a workflow
func (a *App) DeleteWorkflow(workflowID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Administrator" {
		return errors.New("samo administrator može brisati radne tokove")
	}

	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.DeleteWorkflow(workflowID)
}

// CreatePhase creates a new phase in a workflow
func (a *App) CreatePhase(phase models.Faze) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za kreiranje faza")
	}

	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.CreatePhase(phase)
}

// UpdatePhase updates an existing phase
func (a *App) UpdatePhase(phaseID int, phase models.Faze) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za izmenu faza")
	}

	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.UpdatePhase(phaseID, phase)
}

// DeletePhase deletes a phase
func (a *App) DeletePhase(phaseID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za brisanje faza")
	}

	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.DeletePhase(phaseID)
}

// GetWorkflowTemplates returns all workflow templates
func (a *App) GetWorkflowTemplates() ([]models.RadniTokovi, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.workflowService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.GetWorkflowTemplates()
}

// CloneWorkflow creates a copy of an existing workflow
func (a *App) CloneWorkflow(sourceWorkflowID int, newName string) (int, error) {
	if a.currentUser == nil {
		return 0, errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return 0, errors.New("nemate dozvolu za kloniranje radnih tokova")
	}

	if a.workflowService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.workflowService.CloneWorkflow(sourceWorkflowID, newName)
}

// =============================================================================
// Condition Management Methods
// =============================================================================

// GetConditionsByPhase retrieves all conditions for a specific phase
func (a *App) GetConditionsByPhase(phaseID int) ([]models.Uslovi, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.conditionService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.conditionService.GetConditionsByPhase(phaseID)
}

// CreateCondition creates a new condition for a phase
func (a *App) CreateCondition(condition models.Uslovi) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za kreiranje uslova")
	}

	if a.conditionService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.conditionService.CreateCondition(&condition)
}

// UpdateCondition updates an existing condition
func (a *App) UpdateCondition(condition models.Uslovi) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za izmenu uslova")
	}

	if a.conditionService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.conditionService.UpdateCondition(&condition)
}

// DeleteCondition deletes a condition
func (a *App) DeleteCondition(conditionID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.currentUser.NazivUloge != "Rukovodilac projekta" && a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu za brisanje uslova")
	}

	if a.conditionService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.conditionService.DeleteCondition(conditionID)
}

// GetConditionAssessmentsByTask retrieves all condition assessments for a task
func (a *App) GetConditionAssessmentsByTask(taskID int) ([]models.ProcenaUslova, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.conditionService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.conditionService.GetConditionAssessmentsByTask(taskID)
}

// CreateOrUpdateConditionAssessment creates or updates a condition assessment
func (a *App) CreateOrUpdateConditionAssessment(assessment models.ProcenaUslova) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.conditionService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Set the user who is making the assessment
	if a.currentUser != nil {
		userID := a.currentUser.KorisnikID
		assessment.PromenioKorisnikID = &userID
	}

	return a.conditionService.CreateOrUpdateConditionAssessment(&assessment)
}

// GetConditionFulfillmentStatus returns fulfillment status for a task
func (a *App) GetConditionFulfillmentStatus(taskID int) (map[string]interface{}, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.conditionService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	return a.conditionService.GetConditionFulfillmentStatus(taskID)
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