package main

import (
	"github.com/cane/research-institute-system/backend/handlers"
	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	projectHandler *handlers.ProjectHandler,
	taskHandler *handlers.TaskHandler,
	documentHandler *handlers.DocumentHandler,
	workflowHandler *handlers.WorkflowHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	taskCommentHandler *handlers.TaskCommentHandler,
	phaseChangeRequestHandler *handlers.PhaseChangeRequestHandler) {

	// API versioning
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes (no auth required)
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/first-time-setup", authHandler.FirstTimeSetup).Methods("POST")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware)

	// Auth routes
	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/auth/change-password", authHandler.ChangePassword).Methods("POST")
	protected.HandleFunc("/auth/reset-password/{userId}", authHandler.ResetPassword).Methods("POST")

	// User routes
	protected.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	protected.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	protected.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	protected.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	protected.HandleFunc("/roles", userHandler.GetAllRoles).Methods("GET")

	// Project routes
	protected.HandleFunc("/projects", projectHandler.GetAllProjects).Methods("GET")
	protected.HandleFunc("/projects", projectHandler.CreateProject).Methods("POST")
	protected.HandleFunc("/projects/{id}", projectHandler.GetProject).Methods("GET")
	protected.HandleFunc("/projects/{id}", projectHandler.UpdateProject).Methods("PUT")
	protected.HandleFunc("/projects/{id}", projectHandler.DeleteProject).Methods("DELETE")
	protected.HandleFunc("/projects/{id}/members", projectHandler.GetProjectMembers).Methods("GET")
	protected.HandleFunc("/projects/{id}/members", projectHandler.AddProjectMember).Methods("POST")
	protected.HandleFunc("/projects/{id}/members/{userId}", projectHandler.RemoveProjectMember).Methods("DELETE")

	// Task routes
	protected.HandleFunc("/projects/{id}/tasks", taskHandler.GetTasksByProject).Methods("GET")
	protected.HandleFunc("/users/{id}/tasks", taskHandler.GetTasksByUser).Methods("GET")
	protected.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	protected.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	protected.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	// Task Comment routes
	protected.HandleFunc("/tasks/{id}/comments", taskCommentHandler.GetTaskComments).Methods("GET")
	protected.HandleFunc("/tasks/{id}/comments", taskCommentHandler.CreateTaskComment).Methods("POST")
	protected.HandleFunc("/comments/{id}", taskCommentHandler.GetCommentByID).Methods("GET")
	protected.HandleFunc("/comments/{id}", taskCommentHandler.UpdateTaskComment).Methods("PUT")
	protected.HandleFunc("/comments/{id}", taskCommentHandler.DeleteTaskComment).Methods("DELETE")
	protected.HandleFunc("/users/{id}/comments", taskCommentHandler.GetUserComments).Methods("GET")

	// Phase Change Request routes
	protected.HandleFunc("/tasks/{id}/phase-requests", phaseChangeRequestHandler.GetTaskPhaseChangeRequests).Methods("GET")
	protected.HandleFunc("/tasks/{id}/phase-requests", phaseChangeRequestHandler.CreatePhaseChangeRequest).Methods("POST")
	protected.HandleFunc("/phase-requests/pending", phaseChangeRequestHandler.GetPendingPhaseChangeRequests).Methods("GET")
	protected.HandleFunc("/phase-requests/{id}", phaseChangeRequestHandler.GetPhaseChangeRequestByID).Methods("GET")
	protected.HandleFunc("/phase-requests/{id}/status", phaseChangeRequestHandler.UpdatePhaseChangeRequestStatus).Methods("PUT")
	protected.HandleFunc("/phase-requests/{id}", phaseChangeRequestHandler.DeletePhaseChangeRequest).Methods("DELETE")
	protected.HandleFunc("/users/{id}/phase-requests", phaseChangeRequestHandler.GetUserPhaseChangeRequests).Methods("GET")

	// Document routes
	protected.HandleFunc("/documents", documentHandler.GetAllDocuments).Methods("GET")
	protected.HandleFunc("/documents", documentHandler.UploadDocument).Methods("POST")
	protected.HandleFunc("/projects/{id}/documents", documentHandler.GetDocumentsByProject).Methods("GET")
	protected.HandleFunc("/documents/{id}", documentHandler.GetDocument).Methods("GET")
	protected.HandleFunc("/documents/{id}", documentHandler.UpdateDocument).Methods("PUT")
	protected.HandleFunc("/documents/{id}", documentHandler.DeleteDocument).Methods("DELETE")
	protected.HandleFunc("/documents/{id}/versions", documentHandler.GetDocumentVersions).Methods("GET")
	protected.HandleFunc("/documents/{id}/tags", documentHandler.GetDocumentTags).Methods("GET")
	protected.HandleFunc("/documents/{id}/tags", documentHandler.AddDocumentTag).Methods("POST")
	protected.HandleFunc("/documents/{id}/tags/{tagId}", documentHandler.RemoveDocumentTag).Methods("DELETE")
	protected.HandleFunc("/documents/{id}/metadata", documentHandler.GetDocumentMetadata).Methods("GET")
	protected.HandleFunc("/documents/{id}/metadata", documentHandler.UpdateDocumentMetadata).Methods("PUT")
	protected.HandleFunc("/folders", documentHandler.GetAllFolders).Methods("GET")
	protected.HandleFunc("/folders", documentHandler.CreateFolder).Methods("POST")

	// Workflow routes
	protected.HandleFunc("/workflows", workflowHandler.GetAllWorkflows).Methods("GET")
	protected.HandleFunc("/workflows", workflowHandler.CreateWorkflow).Methods("POST")
	protected.HandleFunc("/workflows/{id}/phases", workflowHandler.GetWorkflowPhases).Methods("GET")
	protected.HandleFunc("/workflows/{id}/phases", workflowHandler.CreatePhase).Methods("POST")

	// Analytics routes
	protected.HandleFunc("/analytics/dashboard", analyticsHandler.GetDashboardStats).Methods("GET")
	protected.HandleFunc("/analytics/activity-logs", analyticsHandler.GetActivityLogs).Methods("GET")
}
