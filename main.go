package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/cane/research-institute-system/backend/models"
	"github.com/cane/research-institute-system/backend/repositories"
	"github.com/cane/research-institute-system/backend/services"

	_ "github.com/lib/pq"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
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

// OnStartup is called when the app starts up
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database connection with better error handling
	// To configure: Set environment variables DB_HOST, DB_USER, DB_PASSWORD, DB_NAME
	// Default: postgres://postgres:password@localhost:5432/research_institute?sslmode=disable
	dbHost := getEnvOrDefault("DB_HOST", "localhost:5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "password")
	dbName := getEnvOrDefault("DB_NAME", "research_institute")

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		log.Printf("Application will continue without database. To configure database:")
		log.Printf("1. Install PostgreSQL")
		log.Printf("2. Create database 'research_institute'")
		log.Printf("3. Set environment variables: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME")
		log.Printf("4. Run the SQL schema from database/schema.sql")
		return
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		log.Printf("Database connection string: %s", connStr)
		log.Printf("Application will continue without database.")
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
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
