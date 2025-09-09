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
	ctx         context.Context
	db          *sql.DB
	authService *services.AuthService
	userRepo    *repositories.UserRepository
	projectRepo *repositories.ProjectRepository
	currentUser *models.User
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
