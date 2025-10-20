// ============================================================================
// workflow_service.go - Workflow Management Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type WorkflowService struct {
	db *sql.DB
}

func NewWorkflowService(db *sql.DB) *WorkflowService {
	return &WorkflowService{db: db}
}

func (s *WorkflowService) GetAllWorkflows() ([]models.RadniTokovi, error) {
	query := `
		SELECT radni_tok_id, naziv, tip_toka, opis, da_li_je_sablon
		FROM radnitokovi
		ORDER BY naziv
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.RadniTokovi
	for rows.Next() {
		var workflow models.RadniTokovi
		err := rows.Scan(
			&workflow.RadniTokID, &workflow.Naziv, &workflow.TipToka,
			&workflow.Opis, &workflow.DaLiJeSablon,
		)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

func (s *WorkflowService) GetWorkflowByID(id int) (*models.RadniTokovi, error) {
	var wf models.RadniTokovi
	err := s.db.QueryRow(`SELECT radni_tok_id, naziv, tip_toka, opis, da_li_je_sablon FROM radnitokovi WHERE radni_tok_id = $1`, id).
		Scan(&wf.RadniTokID, &wf.Naziv, &wf.TipToka, &wf.Opis, &wf.DaLiJeSablon)
	if err != nil {
		return nil, err
	}
	return &wf, nil
}

func (s *WorkflowService) GetWorkflowPhases(workflowID int) ([]models.Faze, error) {
	query := `
		SELECT faza_id, radni_tok_id, naziv_faze, redosled
		FROM faze
		WHERE radni_tok_id = $1
		ORDER BY redosled
	`

	rows, err := s.db.Query(query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phases []models.Faze
	for rows.Next() {
		var phase models.Faze
		err := rows.Scan(
			&phase.FazaID, &phase.RadniTokID, &phase.NazivFaze, &phase.Redosled,
		)
		if err != nil {
			return nil, err
		}
		phases = append(phases, phase)
	}

	return phases, nil
}

// GetPhasesByWorkflow is an alias for GetWorkflowPhases
func (s *WorkflowService) GetPhasesByWorkflow(workflowID int) ([]models.Faze, error) {
	return s.GetWorkflowPhases(workflowID)
}

func (s *WorkflowService) CreateWorkflow(workflow models.RadniTokovi) error {
	query := `
		INSERT INTO radnitokovi (naziv, tip_toka, opis, da_li_je_sablon)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(query, workflow.Naziv, workflow.TipToka,
		workflow.Opis, workflow.DaLiJeSablon)

	return err
}

// UpdateWorkflow updates an existing workflow
func (s *WorkflowService) UpdateWorkflow(workflowID int, workflow models.RadniTokovi) error {
	query := `
		UPDATE radnitokovi 
		SET naziv = $1, tip_toka = $2, opis = $3, da_li_je_sablon = $4
		WHERE radni_tok_id = $5
	`

	_, err := s.db.Exec(query, workflow.Naziv, workflow.TipToka,
		workflow.Opis, workflow.DaLiJeSablon, workflowID)

	return err
}

// DeleteWorkflow deletes a workflow
func (s *WorkflowService) DeleteWorkflow(workflowID int) error {
	query := `DELETE FROM radnitokovi WHERE radni_tok_id = $1`
	_, err := s.db.Exec(query, workflowID)
	return err
}

// CreatePhase creates a new phase in a workflow
func (s *WorkflowService) CreatePhase(phase models.Faze) error {
	query := `
		INSERT INTO faze (radni_tok_id, naziv_faze, redosled)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(query, phase.RadniTokID, phase.NazivFaze, phase.Redosled)
	return err
}

// UpdatePhase updates an existing phase
func (s *WorkflowService) UpdatePhase(phaseID int, phase models.Faze) error {
	query := `
		UPDATE faze 
		SET naziv_faze = $1, redosled = $2
		WHERE faza_id = $3
	`

	_, err := s.db.Exec(query, phase.NazivFaze, phase.Redosled, phaseID)
	return err
}

// DeletePhase deletes a phase from a workflow
func (s *WorkflowService) DeletePhase(phaseID int) error {
	query := `DELETE FROM faze WHERE faza_id = $1`
	_, err := s.db.Exec(query, phaseID)
	return err
}

// ReorderPhases updates the order of phases in a workflow
func (s *WorkflowService) ReorderPhases(phaseOrders []struct {
	PhaseID  int
	NewOrder int
}) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, po := range phaseOrders {
		_, err = tx.Exec(`UPDATE faze SET redosled = $1 WHERE faza_id = $2`, po.NewOrder, po.PhaseID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetWorkflowTemplates returns all workflow templates
func (s *WorkflowService) GetWorkflowTemplates() ([]models.RadniTokovi, error) {
	query := `
		SELECT radni_tok_id, naziv, tip_toka, opis, da_li_je_sablon
		FROM radnitokovi
		WHERE da_li_je_sablon = true
		ORDER BY naziv
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.RadniTokovi
	for rows.Next() {
		var workflow models.RadniTokovi
		err := rows.Scan(
			&workflow.RadniTokID, &workflow.Naziv, &workflow.TipToka,
			&workflow.Opis, &workflow.DaLiJeSablon,
		)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// CloneWorkflow creates a copy of an existing workflow
func (s *WorkflowService) CloneWorkflow(sourceWorkflowID int, newName string) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Get source workflow
	var workflow models.RadniTokovi
	err = tx.QueryRow(`
		SELECT naziv, tip_toka, opis, da_li_je_sablon 
		FROM radnitokovi 
		WHERE radni_tok_id = $1
	`, sourceWorkflowID).Scan(&workflow.Naziv, &workflow.TipToka, &workflow.Opis, &workflow.DaLiJeSablon)
	if err != nil {
		return 0, err
	}

	// Create new workflow
	var newWorkflowID int
	err = tx.QueryRow(`
		INSERT INTO radnitokovi (naziv, tip_toka, opis, da_li_je_sablon)
		VALUES ($1, $2, $3, false)
		RETURNING radni_tok_id
	`, newName, workflow.TipToka, workflow.Opis).Scan(&newWorkflowID)
	if err != nil {
		return 0, err
	}

	// Clone phases
	_, err = tx.Exec(`
		INSERT INTO faze (radni_tok_id, naziv_faze, redosled)
		SELECT $1, naziv_faze, redosled
		FROM faze
		WHERE radni_tok_id = $2
	`, newWorkflowID, sourceWorkflowID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	return newWorkflowID, err
}
