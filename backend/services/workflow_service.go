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

func (s *WorkflowService) CreateWorkflow(workflow models.RadniTokovi) error {
	query := `
		INSERT INTO radnitokovi (naziv, tip_toka, opis, da_li_je_sablon)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(query, workflow.Naziv, workflow.TipToka,
		workflow.Opis, workflow.DaLiJeSablon)

	return err
}

func (s *WorkflowService) CreatePhase(phase models.Faze) error {
	query := `
		INSERT INTO faze (radni_tok_id, naziv_faze, redosled)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(query, phase.RadniTokID, phase.NazivFaze, phase.Redosled)
	return err
}

func (s *WorkflowService) UpdateWorkflow(wf models.RadniTokovi) error {
	_, err := s.db.Exec(`UPDATE radnitokovi SET naziv = $1, tip_toka = $2, opis = $3, da_li_je_sablon = $4 WHERE radni_tok_id = $5`,
		wf.Naziv, wf.TipToka, wf.Opis, wf.DaLiJeSablon, wf.RadniTokID)
	return err
}

func (s *WorkflowService) DeleteWorkflow(id int) error {
	_, err := s.db.Exec(`DELETE FROM radnitokovi WHERE radni_tok_id = $1`, id)
	return err
}

func (s *WorkflowService) UpdatePhase(phase models.Faze) error {
	_, err := s.db.Exec(`UPDATE faze SET naziv_faze = $1, redosled = $2 WHERE faza_id = $3`, phase.NazivFaze, phase.Redosled, phase.FazaID)
	return err
}

func (s *WorkflowService) DeletePhase(phaseID int) error {
	_, err := s.db.Exec(`DELETE FROM faze WHERE faza_id = $1`, phaseID)
	return err
}
