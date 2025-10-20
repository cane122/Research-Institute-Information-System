// ============================================================================
// condition_service.go - Phase Transition Conditions Service
// ============================================================================

package services

import (
	"database/sql"
	"fmt"

	"github.com/cane/research-institute-system/backend/models"
)

type ConditionService struct {
	db *sql.DB
}

func NewConditionService(db *sql.DB) *ConditionService {
	return &ConditionService{db: db}
}

// =============================================================================
// Uslov (Condition) Management
// =============================================================================

// GetConditionsByPhase retrieves all conditions for a specific phase
func (s *ConditionService) GetConditionsByPhase(phaseID int) ([]models.Uslovi, error) {
	query := `
		SELECT u.uslov_id, u.faza_id, u.opis, u.kriterijum, u.kreiran_datuma,
		       f.naziv_faze
		FROM uslovi u
		JOIN faze f ON u.faza_id = f.faza_id
		WHERE u.faza_id = $1
		ORDER BY u.uslov_id
	`

	rows, err := s.db.Query(query, phaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize with empty slice instead of nil to ensure JavaScript gets []
	conditions := make([]models.Uslovi, 0)
	for rows.Next() {
		var condition models.Uslovi
		err := rows.Scan(
			&condition.UslovID, &condition.FazaID, &condition.Opis,
			&condition.Kriterijum, &condition.KreiranDatuma, &condition.NazivFaze,
		)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, condition)
	}

	return conditions, nil
}

// GetConditionByID retrieves a single condition by ID
func (s *ConditionService) GetConditionByID(conditionID int) (*models.Uslovi, error) {
	query := `
		SELECT u.uslov_id, u.faza_id, u.opis, u.kriterijum, u.kreiran_datuma,
		       f.naziv_faze
		FROM uslovi u
		JOIN faze f ON u.faza_id = f.faza_id
		WHERE u.uslov_id = $1
	`

	var condition models.Uslovi
	err := s.db.QueryRow(query, conditionID).Scan(
		&condition.UslovID, &condition.FazaID, &condition.Opis,
		&condition.Kriterijum, &condition.KreiranDatuma, &condition.NazivFaze,
	)
	if err != nil {
		return nil, err
	}

	return &condition, nil
}

// CreateCondition creates a new condition for a phase
func (s *ConditionService) CreateCondition(condition *models.Uslovi) error {
	query := `
		INSERT INTO uslovi (faza_id, opis, kriterijum)
		VALUES ($1, $2, $3)
		RETURNING uslov_id, kreiran_datuma
	`

	err := s.db.QueryRow(
		query,
		condition.FazaID,
		condition.Opis,
		condition.Kriterijum,
	).Scan(&condition.UslovID, &condition.KreiranDatuma)

	return err
}

// UpdateCondition updates an existing condition
func (s *ConditionService) UpdateCondition(condition *models.Uslovi) error {
	query := `
		UPDATE uslovi
		SET opis = $1, kriterijum = $2
		WHERE uslov_id = $3
	`

	result, err := s.db.Exec(
		query,
		condition.Opis,
		condition.Kriterijum,
		condition.UslovID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("condition with ID %d not found", condition.UslovID)
	}

	return nil
}

// DeleteCondition deletes a condition
func (s *ConditionService) DeleteCondition(conditionID int) error {
	query := `DELETE FROM uslovi WHERE uslov_id = $1`

	result, err := s.db.Exec(query, conditionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("condition with ID %d not found", conditionID)
	}

	return nil
}

// =============================================================================
// ProcenaUslova (Condition Assessment) Management
// =============================================================================

// GetConditionAssessmentsByTask retrieves all condition assessments for a specific task
func (s *ConditionService) GetConditionAssessmentsByTask(taskID int) ([]models.ProcenaUslova, error) {
	query := `
		SELECT p.procena_id, p.zadatak_id, p.uslov_id, p.ispunjen, p.napomena,
		       p.promenio_korisnik_id, p.datum_procene,
		       u.opis as opis_uslova, u.kriterijum as kriterijum_uslova,
		       z.naziv_zadatka,
		       COALESCE(k.korisnicko_ime, '') as ime_korisnika
		FROM procenauslova p
		JOIN uslovi u ON p.uslov_id = u.uslov_id
		JOIN zadaci z ON p.zadatak_id = z.zadatak_id
		LEFT JOIN korisnici k ON p.promenio_korisnik_id = k.korisnik_id
		WHERE p.zadatak_id = $1
		ORDER BY u.uslov_id
	`

	rows, err := s.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []models.ProcenaUslova
	for rows.Next() {
		var assessment models.ProcenaUslova
		err := rows.Scan(
			&assessment.ProcenaID, &assessment.ZadatakID, &assessment.UslovID,
			&assessment.Ispunjen, &assessment.Napomena, &assessment.PromenioKorisnikID,
			&assessment.DatumProcene, &assessment.OpisUslova, &assessment.KriterijumUslova,
			&assessment.NazivZadatka, &assessment.ImeKorisnika,
		)
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// GetConditionAssessmentsByPhase retrieves all assessments for conditions in a specific phase
func (s *ConditionService) GetConditionAssessmentsByPhase(taskID, phaseID int) ([]models.ProcenaUslova, error) {
	query := `
		SELECT p.procena_id, p.zadatak_id, p.uslov_id, p.ispunjen, p.napomena,
		       p.promenio_korisnik_id, p.datum_procene,
		       u.opis as opis_uslova, u.kriterijum as kriterijum_uslova,
		       z.naziv_zadatka,
		       COALESCE(k.korisnicko_ime, '') as ime_korisnika
		FROM procenauslova p
		JOIN uslovi u ON p.uslov_id = u.uslov_id
		JOIN zadaci z ON p.zadatak_id = z.zadatak_id
		LEFT JOIN korisnici k ON p.promenio_korisnik_id = k.korisnik_id
		WHERE p.zadatak_id = $1 AND u.faza_id = $2
		ORDER BY u.uslov_id
	`

	rows, err := s.db.Query(query, taskID, phaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []models.ProcenaUslova
	for rows.Next() {
		var assessment models.ProcenaUslova
		err := rows.Scan(
			&assessment.ProcenaID, &assessment.ZadatakID, &assessment.UslovID,
			&assessment.Ispunjen, &assessment.Napomena, &assessment.PromenioKorisnikID,
			&assessment.DatumProcene, &assessment.OpisUslova, &assessment.KriterijumUslova,
			&assessment.NazivZadatka, &assessment.ImeKorisnika,
		)
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// CreateOrUpdateConditionAssessment creates or updates a condition assessment for a task
func (s *ConditionService) CreateOrUpdateConditionAssessment(assessment *models.ProcenaUslova) error {
	query := `
		INSERT INTO procenauslova (zadatak_id, uslov_id, ispunjen, napomena, promenio_korisnik_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (zadatak_id, uslov_id)
		DO UPDATE SET
			ispunjen = EXCLUDED.ispunjen,
			napomena = EXCLUDED.napomena,
			promenio_korisnik_id = EXCLUDED.promenio_korisnik_id,
			datum_procene = CURRENT_TIMESTAMP
		RETURNING procena_id, datum_procene
	`

	err := s.db.QueryRow(
		query,
		assessment.ZadatakID,
		assessment.UslovID,
		assessment.Ispunjen,
		assessment.Napomena,
		assessment.PromenioKorisnikID,
	).Scan(&assessment.ProcenaID, &assessment.DatumProcene)

	return err
}

// DeleteConditionAssessment deletes a condition assessment
func (s *ConditionService) DeleteConditionAssessment(assessmentID int) error {
	query := `DELETE FROM procenauslova WHERE procena_id = $1`

	result, err := s.db.Exec(query, assessmentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("condition assessment with ID %d not found", assessmentID)
	}

	return nil
}

// CheckAllConditionsFulfilled checks if all conditions for a phase are fulfilled for a task
func (s *ConditionService) CheckAllConditionsFulfilled(taskID, phaseID int) (bool, error) {
	query := `
		SELECT COUNT(*) as total_conditions,
		       COUNT(*) FILTER (WHERE p.ispunjen = TRUE) as fulfilled_conditions
		FROM uslovi u
		LEFT JOIN procenauslova p ON u.uslov_id = p.uslov_id AND p.zadatak_id = $1
		WHERE u.faza_id = $2
	`

	var totalConditions, fulfilledConditions int
	err := s.db.QueryRow(query, taskID, phaseID).Scan(&totalConditions, &fulfilledConditions)
	if err != nil {
		return false, err
	}

	// If there are no conditions, consider it as fulfilled
	if totalConditions == 0 {
		return true, nil
	}

	return totalConditions == fulfilledConditions, nil
}

// GetConditionFulfillmentStatus returns the fulfillment status for a task's current phase
func (s *ConditionService) GetConditionFulfillmentStatus(taskID int) (map[string]interface{}, error) {
	// First get the task's current phase
	var phaseID int
	err := s.db.QueryRow(`SELECT faza_id FROM zadaci WHERE zadatak_id = $1`, taskID).Scan(&phaseID)
	if err != nil {
		return nil, err
	}

	// Get all conditions and their assessments
	query := `
		SELECT u.uslov_id, u.opis, u.kriterijum,
		       COALESCE(p.ispunjen, FALSE) as ispunjen,
		       p.napomena
		FROM uslovi u
		LEFT JOIN procenauslova p ON u.uslov_id = p.uslov_id AND p.zadatak_id = $1
		WHERE u.faza_id = $2
		ORDER BY u.uslov_id
	`

	rows, err := s.db.Query(query, taskID, phaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type conditionStatus struct {
		UslovID    int     `json:"uslov_id"`
		Opis       string  `json:"opis"`
		Kriterijum string  `json:"kriterijum"`
		Ispunjen   bool    `json:"ispunjen"`
		Napomena   *string `json:"napomena"`
	}

	var conditions []conditionStatus
	totalConditions := 0
	fulfilledCount := 0

	for rows.Next() {
		var cs conditionStatus
		err := rows.Scan(&cs.UslovID, &cs.Opis, &cs.Kriterijum, &cs.Ispunjen, &cs.Napomena)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, cs)
		totalConditions++
		if cs.Ispunjen {
			fulfilledCount++
		}
	}

	allFulfilled := totalConditions > 0 && totalConditions == fulfilledCount

	return map[string]interface{}{
		"task_id":          taskID,
		"phase_id":         phaseID,
		"conditions":       conditions,
		"total_conditions": totalConditions,
		"fulfilled_count":  fulfilledCount,
		"all_fulfilled":    allFulfilled,
		"can_proceed":      allFulfilled,
	}, nil
}
