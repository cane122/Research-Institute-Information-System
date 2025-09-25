package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/cane/research-institute-system/backend/handlers"
	"github.com/cane/research-institute-system/backend/repositories"
	"github.com/cane/research-institute-system/backend/services"
)

func main() {
	// Database connection
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Set database for middleware
	setDB(db)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(db)
	projectService := services.NewProjectService(db)
	taskService := services.NewTaskService(db)
	documentService := services.NewDocumentService(db)
	workflowService := services.NewWorkflowService(db)
	analyticsService := services.NewAnalyticsService(db)
	taskCommentService := services.NewTaskCommentService(db)
	phaseChangeRequestService := services.NewPhaseChangeRequestService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, analyticsService)
	userHandler := handlers.NewUserHandler(userService, analyticsService)
	projectHandler := handlers.NewProjectHandler(projectService, analyticsService)
	taskHandler := handlers.NewTaskHandler(taskService, analyticsService)
	documentHandler := handlers.NewDocumentHandler(documentService, analyticsService)
	workflowHandler := handlers.NewWorkflowHandler(workflowService, analyticsService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	taskCommentHandler := handlers.NewTaskCommentHandler(taskCommentService, analyticsService)
	phaseChangeRequestHandler := handlers.NewPhaseChangeRequestHandler(phaseChangeRequestService, analyticsService)

	// Setup routes
	router := mux.NewRouter()
	setupRoutes(router, authHandler, userHandler, projectHandler, taskHandler, documentHandler, workflowHandler, analyticsHandler, taskCommentHandler, phaseChangeRequestHandler)

	// CORS setup
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func connectDB() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "123")
	dbname := getEnv("DB_NAME", "research_institute")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
