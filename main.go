package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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

// PhaseDurationInfo represents phase name and duration in seconds
type PhaseDurationInfo struct {
	NazivFaze       string `json:"naziv_faze"`
	TrajanjeSekundi int64  `json:"trajanje_sekundi"`
}

// GetDocumentPhaseDurations returns a list of phases and how long the document spent in each
func (a *App) GetDocumentPhaseDurations(documentID int) ([]PhaseDurationInfo, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}
	if a.docPhaseHistSvc == nil || a.fazeService == nil || a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get document to see current phase and creation date
	doc, err := a.documentService.GetDocumentByID(documentID)
	if err != nil {
		return nil, err
	}

	// Get phase history, ordered DESC by datum_promene (newest first)
	history, err := a.docPhaseHistSvc.ListByDocument(documentID)
	if err != nil {
		return nil, err
	}

	// Build phase duration map faza_id -> trajanje
	durations := make(map[int]int64)
	phaseNames := make(map[int]string)

	// If there's no history, document is still in its initial phase
	if len(history) == 0 {
		if doc.TrenutnaFazaID != nil {
			naziv, _ := a.fazeService.GetFazaNaziv(*doc.TrenutnaFazaID)
			trajanje := int64(time.Since(doc.DatumaPostavke).Seconds())
			if trajanje < 0 {
				trajanje = 0
			}
			durations[*doc.TrenutnaFazaID] = trajanje
			phaseNames[*doc.TrenutnaFazaID] = naziv
		}
	} else {
		// Reverse to ASC order (oldest first)
		for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
			history[i], history[j] = history[j], history[i]
		}
		// Initial phase before first change
		if doc.TrenutnaFazaID != nil {
			firstEntry := history[0]
			if firstEntry.PrethodnaFazaID != nil && *firstEntry.PrethodnaFazaID > 0 {
				naziv, _ := a.fazeService.GetFazaNaziv(*firstEntry.PrethodnaFazaID)
				trajanje := int64(firstEntry.DatumPromene.Sub(doc.DatumaPostavke).Seconds())
				if trajanje < 0 {
					trajanje = 0
				}
				durations[*firstEntry.PrethodnaFazaID] = trajanje
				phaseNames[*firstEntry.PrethodnaFazaID] = naziv
			}
		}
		// For each entry, compute duration until next change (or now for last)
		for i, h := range history {
			start := h.DatumPromene
			var end time.Time
			if i+1 < len(history) {
				end = history[i+1].DatumPromene
			} else {
				end = time.Now()
			}
			trajanje := int64(end.Sub(start).Seconds())
			if trajanje < 0 {
				trajanje = 0 // Safety check
			}
			durations[h.NovaFazaID] += trajanje
			naziv, _ := a.fazeService.GetFazaNaziv(h.NovaFazaID)
			phaseNames[h.NovaFazaID] = naziv
		}
	}

	// Get all phases for this workflow
	var allPhases []models.Faze
	if doc.RadniTokID != nil && a.workflowService != nil {
		phases, err := a.workflowService.GetWorkflowPhases(*doc.RadniTokID)
		if err == nil {
			allPhases = phases
		}
	}

	// Compose result: for each phase in workflow, show duration or 0s
	var result []PhaseDurationInfo
	for _, faza := range allPhases {
		dur := durations[faza.FazaID]
		result = append(result, PhaseDurationInfo{
			NazivFaze:       faza.NazivFaze,
			TrajanjeSekundi: dur,
		})
	}
	return result, nil
}

//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx              context.Context
	db               *sql.DB
	authService      *services.AuthService
	documentService  *services.DocumentService
	workflowService  *services.WorkflowService
	roleService      *services.RoleService
	tagService       *services.TagService
	phaseReqService  *services.PhaseChangeRequestService
	docPhaseHistSvc  *services.DocumentPhaseHistoryService
	docVersionSvc    *services.DocumentVersionService
	llmService       *services.LLMService
	analyticsService *services.AnalyticsService
	userRepo         *repositories.UserRepository
	projectRepo      *repositories.ProjectRepository
	currentUser      *models.User
	zadacicService   *services.ZadacicService
	fazeService      *services.FazeService
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
		{"Default", "localhost:5432", "postgres", "postgres", "research_institute_db"},
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
	dbPassword := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "research_institute_db")

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
	a.workflowService = services.NewWorkflowService(db)
	a.roleService = services.NewRoleService(db)
	a.tagService = services.NewTagService(db)
	a.phaseReqService = services.NewPhaseChangeRequestService(db)
	a.docPhaseHistSvc = services.NewDocumentPhaseHistoryService(db)
	a.docVersionSvc = services.NewDocumentVersionService(db)
	a.llmService = services.NewLLMService()
	a.analyticsService = services.NewAnalyticsService(db)
	a.zadacicService = services.NewZadacicService(db)
	a.fazeService = services.NewFazeService(db)
}

// ===================== Zadacici (Checklist) =====================
func (a *App) ListZadaciciByDocument(documentID int) ([]models.Zadacic, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}
	if a.zadacicService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.zadacicService.ListByDocument(documentID)
}

func (a *App) AddZadacic(z models.Zadacic) (int, error) {
	if a.currentUser == nil {
		return 0, errors.New("niste prijavljeni")
	}
	if a.zadacicService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	if z.Opis == "" || z.DokumentID == 0 {
		return 0, errors.New("Opis i dokument su obavezni")
	}
	z.Izvrsen = false
	if err := a.zadacicService.Add(&z); err != nil {
		return 0, err
	}
	return z.ZadacicID, nil
}

func (a *App) UpdateZadacicStatus(zadacicID int, izvrsen bool) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}
	if a.zadacicService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.zadacicService.UpdateStatus(zadacicID, izvrsen)
}

func (a *App) DeleteZadacic(zadacicID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}
	if a.zadacicService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.zadacicService.Delete(zadacicID)
}

// ============================================================================
// Workflow CRUD
// ============================================================================

func (a *App) GetAllWorkflows() ([]models.RadniTokovi, error) {
	if a.workflowService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.GetAllWorkflows()
}

func (a *App) GetWorkflowByID(id int) (*models.RadniTokovi, error) {
	if a.workflowService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.GetWorkflowByID(id)
}

func (a *App) GetWorkflowPhases(workflowID int) ([]models.Faze, error) {
	if a.workflowService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.GetWorkflowPhases(workflowID)
}

func (a *App) CreateWorkflow(wf models.RadniTokovi) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.CreateWorkflow(wf)
}

func (a *App) UpdateWorkflow(wf models.RadniTokovi) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.UpdateWorkflow(wf)
}

func (a *App) DeleteWorkflow(id int) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.DeleteWorkflow(id)
}

func (a *App) CreatePhase(phase models.Faze) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.CreatePhase(phase)
}

func (a *App) UpdatePhase(phase models.Faze) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.UpdatePhase(phase)
}

func (a *App) DeletePhase(phaseID int) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.workflowService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.DeletePhase(phaseID)
}

// CreateWorkflowWithPhasesRequest represents the request to create a workflow with phases
type CreateWorkflowWithPhasesRequest struct {
	Naziv   string   `json:"naziv"`
	TipToka string   `json:"tip_toka"`
	Faze    []string `json:"faze"`
}

// CreateWorkflowWithPhases creates a new workflow with phases in one transaction
func (a *App) CreateWorkflowWithPhases(req CreateWorkflowWithPhasesRequest) (int, error) {
	if a.currentUser == nil {
		return 0, errors.New("niste prijavljeni")
	}
	if a.workflowService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.workflowService.CreateWorkflowWithPhases(req.Naziv, req.TipToka, req.Faze)
}

// ============================================================================
// Roles CRUD (Admin only)
// ============================================================================

func (a *App) ListRoles() ([]models.Uloge, error) {
	if a.roleService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.roleService.GetAll()
}

func (a *App) CreateRole(name string) (int, error) {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return 0, errors.New("nemate dozvolu")
	}
	if a.roleService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	role := models.Uloge{NazivUloge: name}
	if err := a.roleService.Create(&role); err != nil {
		return 0, err
	}
	return role.UlogaID, nil
}

func (a *App) UpdateRole(id int, name string) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.roleService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.roleService.Update(models.Uloge{UlogaID: id, NazivUloge: name})
}

func (a *App) DeleteRole(id int) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.roleService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.roleService.Delete(id)
}

// ============================================================================
// Tags CRUD (Admin for mutations)
// ============================================================================

func (a *App) CreateTag(name string) (int, error) {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return 0, errors.New("nemate dozvolu")
	}
	if a.tagService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	tag := models.Tagovi{NazivTaga: name}
	if err := a.tagService.Create(&tag); err != nil {
		return 0, err
	}
	return tag.TagID, nil
}

func (a *App) UpdateTag(id int, name string) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.tagService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.tagService.Update(models.Tagovi{TagID: id, NazivTaga: name})
}

func (a *App) DeleteTag(id int) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.tagService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.tagService.Delete(id)
}

// ============================================================================
// Phase Change Requests
// ============================================================================

// List phase change requests by task
func (a *App) ListPhaseChangeRequests(taskID int) ([]models.ZahteviPromeneFaze, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}
	if a.phaseReqService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.phaseReqService.ListByTask(taskID)
}

// List phase change requests by document
func (a *App) ListPhaseChangeRequestsByDocument(documentID int) ([]models.ZahteviPromeneFaze, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}
	if a.phaseReqService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.phaseReqService.ListByDocument(documentID)
}

// Create phase change request (for task or document)
func (a *App) CreatePhaseChangeRequest(req models.ZahteviPromeneFaze) (int, error) {
	if a.currentUser == nil {
		return 0, errors.New("niste prijavljeni")
	}
	if a.phaseReqService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	req.PodnosilacZahtevaID = a.currentUser.KorisnikID
	// If both IDs are nil, reject
	if req.ZadatakID == nil && req.DokumentID == nil {
		return 0, errors.New("Mora biti prosleđen zadatak_id ili dokument_id")
	}
	if err := a.phaseReqService.Create(&req); err != nil {
		return 0, err
	}
	return req.ZahtevID, nil
}

func (a *App) UpdatePhaseChangeRequestStatus(id int, status string, komentar *string) error {
	if a.currentUser == nil {
		return errors.New("nemate dozvolu")
	}
	// Allow administrators and project leaders
	role := strings.ToLower(a.currentUser.NazivUloge)
	if role != "administrator" && role != "rukovodilac projekta" {
		return errors.New("nemate dozvolu")
	}
	if a.phaseReqService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get the request details first
	request, err := a.phaseReqService.GetByID(id)
	if err != nil {
		return fmt.Errorf("greška pri preuzimanju zahteva: %v", err)
	}

	// Update the request status
	if err := a.phaseReqService.UpdateStatus(id, status, komentar); err != nil {
		return err
	}

	// If approved and it's a document request, update the document's phase
	if strings.ToLower(status) == "odobren" && request.DokumentID != nil {
		// Get current document to find old phase
		doc, err := a.documentService.GetDocumentByID(*request.DokumentID)
		if err != nil {
			return fmt.Errorf("greška pri preuzimanju dokumenta: %v", err)
		}

		// Update document's current phase
		_, err = a.db.Exec(`UPDATE dokumenti SET trenutna_faza_id = $1, poslednja_izmena = CURRENT_TIMESTAMP WHERE dokument_id = $2`,
			request.ZahtevanaFazaID, *request.DokumentID)
		if err != nil {
			return fmt.Errorf("greška pri ažuriranju faze dokumenta: %v", err)
		}

		// Add entry to phase history
		_, err = a.db.Exec(`INSERT INTO istorijafazadokumenta (dokument_id, prethodna_faza_id, nova_faza_id, korisnik_id) VALUES ($1, $2, $3, $4)`,
			*request.DokumentID, doc.TrenutnaFazaID, request.ZahtevanaFazaID, a.currentUser.KorisnikID)
		if err != nil {
			return fmt.Errorf("greška pri čuvanju istorije faza: %v", err)
		}

		// Log activity for phase change
		if a.analyticsService != nil && a.fazeService != nil {
			var nazivStareFaze, nazivNoveFaze string
			if doc.TrenutnaFazaID != nil {
				nazivStareFaze, _ = a.fazeService.GetFazaNaziv(*doc.TrenutnaFazaID)
			}
			nazivNoveFaze, _ = a.fazeService.GetFazaNaziv(request.ZahtevanaFazaID)
			_ = a.analyticsService.LogActivity(&a.currentUser.KorisnikID, models.ActivityLogRequest{
				TipAktivnosti: "PHASE_CHANGE",
				EntitetTip:    "DOKUMENT",
				EntitetID:     *request.DokumentID,
				NazivEntiteta: doc.NazivDokumenta,
				Opis:          fmt.Sprintf("Promena faze dokumenta sa '%s' na '%s'", nazivStareFaze, nazivNoveFaze),
				Rezultat:      "SUCCESS",
			})
		}
	}

	return nil
}

func (a *App) DeletePhaseChangeRequest(id int) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.phaseReqService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.phaseReqService.Delete(id)
}

// ChangeDocumentPhase directly changes a document's phase (leader/admin only, no request approval needed)
func (a *App) ChangeDocumentPhase(documentID int, newPhaseID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	// Allow administrators and project leaders
	role := strings.ToLower(a.currentUser.NazivUloge)
	if role != "administrator" && role != "rukovodilac projekta" {
		return errors.New("nemate dozvolu za direktnu promenu faze")
	}

	if a.documentService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get current document to find old phase
	doc, err := a.documentService.GetDocumentByID(documentID)
	if err != nil {
		return fmt.Errorf("greška pri preuzimanju dokumenta: %v", err)
	}

	// Update document's current phase
	_, err = a.db.Exec(`UPDATE dokumenti SET trenutna_faza_id = $1, poslednja_izmena = CURRENT_TIMESTAMP WHERE dokument_id = $2`,
		newPhaseID, documentID)
	if err != nil {
		return fmt.Errorf("greška pri ažuriranju faze dokumenta: %v", err)
	}

	// Add entry to phase history
	_, err = a.db.Exec(`INSERT INTO istorijafazadokumenta (dokument_id, prethodna_faza_id, nova_faza_id, korisnik_id) VALUES ($1, $2, $3, $4)`,
		documentID, doc.TrenutnaFazaID, newPhaseID, a.currentUser.KorisnikID)
	if err != nil {
		return fmt.Errorf("greška pri čuvanju istorije faza: %v", err)
	}

	// Log activity for direct phase change
	if a.analyticsService != nil && a.fazeService != nil {
		var nazivStareFaze, nazivNoveFaze string
		if doc.TrenutnaFazaID != nil {
			nazivStareFaze, _ = a.fazeService.GetFazaNaziv(*doc.TrenutnaFazaID)
		}
		nazivNoveFaze, _ = a.fazeService.GetFazaNaziv(newPhaseID)
		_ = a.analyticsService.LogActivity(&a.currentUser.KorisnikID, models.ActivityLogRequest{
			TipAktivnosti: "PHASE_CHANGE",
			EntitetTip:    "DOKUMENT",
			EntitetID:     documentID,
			NazivEntiteta: doc.NazivDokumenta,
			Opis:          fmt.Sprintf("Direktna promena faze dokumenta sa '%s' na '%s'", nazivStareFaze, nazivNoveFaze),
			Rezultat:      "SUCCESS",
		})
	}

	return nil
}

// ============================================================================
// Document Phase History
// ============================================================================

func (a *App) GetDocumentPhaseHistory(documentID int) ([]models.IstorijaFazaDokumenta, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}
	if a.docPhaseHistSvc == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.docPhaseHistSvc.ListByDocument(documentID)
}

func (a *App) AddDocumentPhaseHistory(entry models.IstorijaFazaDokumenta) (int, error) {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return 0, errors.New("nemate dozvolu")
	}
	if a.docPhaseHistSvc == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	if entry.KorisnikID == 0 {
		entry.KorisnikID = a.currentUser.KorisnikID
	}
	if err := a.docPhaseHistSvc.Create(&entry); err != nil {
		return 0, err
	}
	return entry.IstorijaID, nil
}

func (a *App) DeleteDocumentPhaseHistory(id int) error {
	if a.currentUser == nil || a.currentUser.NazivUloge != "Administrator" {
		return errors.New("nemate dozvolu")
	}
	if a.docPhaseHistSvc == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.docPhaseHistSvc.Delete(id)
}

// ============================================================================
// Document Versions (permission-aware)
// ============================================================================

func (a *App) AddDocumentVersion(v models.VerzijeDokumenata) (int, error) {
	if a.currentUser == nil {
		return 0, errors.New("niste prijavljeni")
	}
	if a.docVersionSvc == nil || a.documentService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	// Check permission to write to this document
	ok, err := a.documentService.CheckUserPermission(v.DokumentID, a.currentUser.KorisnikID, "write")
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("nemate dozvolu za izmene dokumenta")
	}
	v.PostavioKorisnikID = a.currentUser.KorisnikID
	if err := a.docVersionSvc.Create(&v); err != nil {
		return 0, err
	}
	// Log version creation as a document edit
	if a.analyticsService != nil {
		if doc, derr := a.documentService.GetDocumentByID(v.DokumentID); derr == nil {
			_ = a.analyticsService.LogActivity(&a.currentUser.KorisnikID, models.ActivityLogRequest{
				TipAktivnosti: "DOCUMENT_EDIT",
				EntitetTip:    "DOKUMENT",
				EntitetID:     v.DokumentID,
				NazivEntiteta: doc.NazivDokumenta,
				Opis:          "Dodata nova verzija dokumenta",
				Rezultat:      "SUCCESS",
			})
		}
	}
	return v.VerzijaID, nil
}

func (a *App) DeleteDocumentVersion(versionID int, documentID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}
	if a.docVersionSvc == nil || a.documentService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	ok, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "delete")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("nemate dozvolu za brisanje")
	}
	if err := a.docVersionSvc.Delete(versionID); err != nil {
		return err
	}
	// Log version deletion as a document edit
	if a.analyticsService != nil {
		if doc, derr := a.documentService.GetDocumentByID(documentID); derr == nil {
			_ = a.analyticsService.LogActivity(&a.currentUser.KorisnikID, models.ActivityLogRequest{
				TipAktivnosti: "DOCUMENT_EDIT",
				EntitetTip:    "DOKUMENT",
				EntitetID:     documentID,
				NazivEntiteta: doc.NazivDokumenta,
				Opis:          "Obrisana verzija dokumenta",
				Rezultat:      "SUCCESS",
			})
		}
	}
	return nil
}

// UploadDocumentVersion creates a new version for a document
func (a *App) UploadDocumentVersion(documentID int, versionLabel *string, fileData []byte, originalFileName string) (int, error) {
	if a.currentUser == nil {
		return 0, errors.New("niste prijavljeni")
	}
	if a.documentService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	// Permission: require write
	ok, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "write")
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("nemate dozvolu za izmene dokumenta")
	}
	vid, err := a.documentService.SaveNewVersion(documentID, a.currentUser.KorisnikID, versionLabel, fileData, originalFileName)
	if err != nil {
		return 0, err
	}
	// Log version upload as document edit
	if a.analyticsService != nil {
		if doc, derr := a.documentService.GetDocumentByID(documentID); derr == nil {
			_ = a.analyticsService.LogActivity(&a.currentUser.KorisnikID, models.ActivityLogRequest{
				TipAktivnosti: "DOCUMENT_EDIT",
				EntitetTip:    "DOKUMENT",
				EntitetID:     documentID,
				NazivEntiteta: doc.NazivDokumenta,
				Opis:          "Postavljena nova verzija dokumenta",
				Rezultat:      "SUCCESS",
			})
		}
	}
	return vid, nil
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

// GetAllUsersForDocuments returns all users for document permission assignment (any logged-in user)
func (a *App) GetAllUsersForDocuments() ([]models.User, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
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

// GetProjectMembers returns users assigned to a given project
func (a *App) GetProjectMembers(projectID int) ([]models.User, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}
	if a.projectRepo == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	return a.projectRepo.GetMembers(projectID)
}

// UpdateProject updates project basic information (leader or admin only)
func (a *App) UpdateProject(project *models.Project) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}
	if a.projectRepo == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	// Only admin or current leader can update
	if a.currentUser.NazivUloge != "Administrator" {
		// fetch existing to verify leader
		existing, err := a.projectRepo.GetByID(project.ProjekatID)
		if err != nil {
			return err
		}
		if existing.RukovodilaID == nil || *existing.RukovodilaID != a.currentUser.KorisnikID {
			return errors.New("nemate dozvolu za izmenu projekta")
		}
	}
	return a.projectRepo.Update(project)
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
	return a.documentService.GetAllDocuments(a.currentUser.KorisnikID)
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
	// Require read permission
	ok, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "read")
	if err != nil {
		return models.Dokumenti{}, err
	}
	if !ok {
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
	// Require read permission
	ok, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "read")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("nemate dozvolu za pregled verzija")
	}
	return a.documentService.GetDocumentVersions(documentID)
}

// DownloadDocumentVersion reads file content for a specific version
func (a *App) DownloadDocumentVersion(versionID int) ([]byte, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Get version info to get file path
	version, err := a.docVersionSvc.GetByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("verzija nije pronađena: %w", err)
	}

	// Require read permission for the owning document
	ok, err := a.documentService.CheckUserPermission(version.DokumentID, a.currentUser.KorisnikID, "read")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("nemate dozvolu za preuzimanje ovog dokumenta")
	}

	// Read file from disk
	fileData, err := os.ReadFile(version.PutanjaDoFajla)
	if err != nil {
		return nil, fmt.Errorf("greška pri čitanju fajla: %w", err)
	}

	return fileData, nil
}

// GetDocumentTags returns all tags for a document
func (a *App) GetDocumentTags(documentID int) ([]models.Tagovi, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	// Require read permission
	ok, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "read")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("nemate dozvolu za pregled tagova")
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

	// If document is being added to a project, only project leader or admin can add
	if req.ProjekatID != nil {
		if a.projectRepo == nil {
			return 0, errors.New("sistem nije povezan sa bazom podataka")
		}
		proj, err := a.projectRepo.GetByID(*req.ProjekatID)
		if err != nil {
			return 0, err
		}
		isLeader := proj.RukovodilaID != nil && *proj.RukovodilaID == a.currentUser.KorisnikID
		isAdmin := a.currentUser.NazivUloge == "Administrator"
		if !isLeader && !isAdmin {
			return 0, errors.New("samo rukovodilac projekta ili administrator mogu da dodaju dokument")
		}
	}

	return a.documentService.UploadDocument(req, fileData, fileName, a.currentUser.KorisnikID)
}

// CanAddProjectDocument returns whether the current user can add a document to the given project
func (a *App) CanAddProjectDocument(projectID int) (bool, error) {
	if a.currentUser == nil {
		return false, errors.New("niste prijavljeni")
	}
	if a.projectRepo == nil {
		return false, errors.New("sistem nije povezan sa bazom podataka")
	}
	proj, err := a.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}
	if proj.RukovodilaID != nil && *proj.RukovodilaID == a.currentUser.KorisnikID {
		return true, nil
	}
	if a.currentUser.NazivUloge == "Administrator" {
		return true, nil
	}
	return false, nil
}

// CreateDocumentWithPermissionsRequest represents the request to create a document with user permissions
type CreateDocumentWithPermissionsRequest struct {
	NazivDokumenta   string  `json:"naziv_dokumenta"`
	ProjekatID       int     `json:"projekat_id"`
	RadniTokID       int     `json:"radni_tok_id"`
	Opis             *string `json:"opis"`
	Rok              *string `json:"rok"`
	KorisniciDozvole []int   `json:"korisnici_dozvole"`
}

// CreateDocumentWithPermissions creates a new document with file and sets permissions for specified users
func (a *App) CreateDocumentWithPermissions(req CreateDocumentWithPermissionsRequest, fileData []byte, fileName string) (int, error) {
	if a.currentUser == nil {
		return 0, errors.New("niste prijavljeni")
	}
	if a.documentService == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}

	// Check if user can add document to this project
	if a.projectRepo == nil {
		return 0, errors.New("sistem nije povezan sa bazom podataka")
	}
	proj, err := a.projectRepo.GetByID(req.ProjekatID)
	if err != nil {
		return 0, err
	}
	isLeader := proj.RukovodilaID != nil && *proj.RukovodilaID == a.currentUser.KorisnikID
	isAdmin := a.currentUser.NazivUloge == "Administrator"
	if !isLeader && !isAdmin {
		return 0, errors.New("samo rukovodilac projekta ili administrator mogu da dodaju dokument")
	}

	// Get first phase of the workflow
	var firstPhaseID *int
	if req.RadniTokID > 0 {
		phases, err := a.workflowService.GetWorkflowPhases(req.RadniTokID)
		if err == nil && len(phases) > 0 {
			firstPhaseID = &phases[0].FazaID
		}
	}

	// Create upload request
	opis := ""
	if req.Opis != nil {
		opis = *req.Opis
	}
	uploadReq := models.UploadDocumentRequest{
		NazivDokumenta: req.NazivDokumenta,
		ProjekatID:     &req.ProjekatID,
		Opis:           opis,
		TipDokumenta:   "Dokument", // Default type
		JezikDokumenta: "Srpski",   // Default language
	}

	// Upload document
	documentID, err := a.documentService.UploadDocument(uploadReq, fileData, fileName, a.currentUser.KorisnikID)
	if err != nil {
		return 0, err
	}

	// Set workflow and phase if provided
	if req.RadniTokID > 0 && firstPhaseID != nil {
		updateQuery := `UPDATE dokumenti SET radni_tok_id = $1, trenutna_faza_id = $2 WHERE dokument_id = $3`
		_, err = a.db.Exec(updateQuery, req.RadniTokID, firstPhaseID, documentID)
		if err != nil {
			log.Printf("Warning: failed to set workflow/phase: %v", err)
		} else {
			// Log initial phase assignment as a phase change
			if a.analyticsService != nil && a.fazeService != nil {
				if doc, derr := a.documentService.GetDocumentByID(documentID); derr == nil {
					nazivFaze, _ := a.fazeService.GetFazaNaziv(*firstPhaseID)
					_ = a.analyticsService.LogActivity(&a.currentUser.KorisnikID, models.ActivityLogRequest{
						TipAktivnosti: "PHASE_CHANGE",
						EntitetTip:    "DOKUMENT",
						EntitetID:     documentID,
						NazivEntiteta: doc.NazivDokumenta,
						Opis:          fmt.Sprintf("Postavljanje početne faze dokumenta na '%s' (radni tok %d)", nazivFaze, req.RadniTokID),
						Rezultat:      "SUCCESS",
					})
				}
			}
		}
	}

	// Set permissions for specified users (give them full access: read, write, delete)
	log.Printf("Setting permissions for document %d, users: %v", documentID, req.KorisniciDozvole)
	for _, userID := range req.KorisniciDozvole {
		permReq := models.DocumentPermissionRequest{
			DokumentID:  documentID,
			KorisnikID:  userID,
			MozeCitati:  true,
			MozeMenjati: true,
			MozeBrisati: true,
		}
		log.Printf("Setting permission for user %d on document %d", userID, documentID)
		if err := a.documentService.SetDocumentPermission(permReq); err != nil {
			log.Printf("ERROR: failed to set permission for user %d: %v", userID, err)
		} else {
			log.Printf("SUCCESS: permission set for user %d on document %d", userID, documentID)
		}
	}

	log.Printf("Document %d created successfully with %d user permissions", documentID, len(req.KorisniciDozvole))
	return documentID, nil
}

// UpdateDocument updates an existing document
func (a *App) UpdateDocument(documentID int, req models.UploadDocumentRequest) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	// Require write permission
	ok, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "write")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("nemate dozvolu za izmenu dokumenta")
	}
	if err := a.documentService.UpdateDocument(documentID, req); err != nil {
		return err
	}
	// Log activity for document edit
	if a.analyticsService != nil {
		_ = a.analyticsService.LogActivity(&a.currentUser.KorisnikID, models.ActivityLogRequest{
			TipAktivnosti: "DOCUMENT_EDIT",
			EntitetTip:    "DOKUMENT",
			EntitetID:     documentID,
			NazivEntiteta: req.NazivDokumenta,
			Opis:          "Izmena meta podataka dokumenta",
			Rezultat:      "SUCCESS",
		})
	}
	return nil
}

// DeleteDocument deletes a document
func (a *App) DeleteDocument(documentID int) error {
	if a.currentUser == nil {
		return errors.New("niste prijavljeni")
	}

	if a.documentService == nil {
		return errors.New("sistem nije povezan sa bazom podataka")
	}
	// Require delete permission
	ok, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "delete")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("nemate dozvolu za brisanje dokumenta")
	}
	return a.documentService.DeleteDocument(documentID)
}

// GetProjectDocuments returns all documents attached to a specific project
func (a *App) GetProjectDocuments(projectID int) ([]models.Dokumenti, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}
	if a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	// List only documents the user can read
	docs, err := a.documentService.GetDocumentsByProject(projectID, a.currentUser.KorisnikID)
	if err != nil {
		return nil, err
	}
	// For each document, calculate progress (percentage of completed zadacici)
	for i := range docs {
		zadacici, err := a.zadacicService.ListByDocument(docs[i].DokumentID)
		if err != nil {
			// If error, just skip progress for this document
			docs[i].Progres = nil
			continue
		}
		total := len(zadacici)
		if total == 0 {
			zero := 0
			docs[i].Progres = &zero
			continue
		}
		completed := 0
		for _, z := range zadacici {
			if z.Izvrsen {
				completed++
			}
		}
		percent := int((float64(completed) / float64(total)) * 100)
		docs[i].Progres = &percent
	}
	return docs, nil
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

	// Check if user is owner/leader/admin OR has read permission
	allowed, err := a.documentService.IsOwnerAdminOrLeader(documentID, a.currentUser.KorisnikID)
	if err != nil {
		return nil, err
	}

	if !allowed {
		// If not owner/leader/admin, check if user has read access to the document
		hasAccess, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "read")
		if err != nil {
			return nil, err
		}
		if !hasAccess {
			return nil, errors.New("nemate dozvolu za pregled ovog dokumenta")
		}
	}

	// User has access to the document, so they can see all permissions
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
	// Only owner/leader/admin may modify permissions
	allowed, err := a.documentService.IsOwnerAdminOrLeader(req.DokumentID, a.currentUser.KorisnikID)
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New("nemate dozvolu za izmenu dozvola")
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
	// Only owner/leader/admin may modify permissions
	allowed, err := a.documentService.IsOwnerAdminOrLeader(documentID, a.currentUser.KorisnikID)
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New("nemate dozvolu za izmenu dozvola")
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

// GetDocumentActivity returns activity logs for a specific document (owner/leader/admin or readers)
func (a *App) GetDocumentActivity(documentID int) ([]models.SkornjeAktivnosti, error) {
	if a.currentUser == nil {
		return nil, errors.New("niste prijavljeni")
	}
	if a.analyticsService == nil || a.documentService == nil {
		return nil, errors.New("sistem nije povezan sa bazom podataka")
	}
	// Allow if user is owner/admin/leader or has read permission on the document
	allowed, err := a.documentService.IsOwnerAdminOrLeader(documentID, a.currentUser.KorisnikID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		hasAccess, err := a.documentService.CheckUserPermission(documentID, a.currentUser.KorisnikID, "read")
		if err != nil {
			return nil, err
		}
		if !hasAccess {
			return nil, errors.New("nemate dozvolu za pregled analitike ovog dokumenta")
		}
	}
	return a.analyticsService.GetActivityForDocument(documentID)
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
