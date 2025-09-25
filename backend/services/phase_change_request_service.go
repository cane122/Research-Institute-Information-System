// ============================================================================
// phase_change_request_service.go - Phase Change Requests Service
// ============================================================================

package services

import (
	"database/sql"
	"fmt"

	"github.com/cane/research-institute-system/backend/models"
)

type PhaseChangeRequestService struct {
	db *sql.DB
}

func NewPhaseChangeRequestService(db *sql.DB) *PhaseChangeRequestService {
	return &PhaseChangeRequestService{db: db}
}

// GetRequestsByTask retrieves all phase change requests for a specific task
func (s *PhaseChangeRequestService) GetRequestsByTask(taskID int) ([]models.ZahteviPromeneFaze, error) {
	query := `
		SELECT zpr.zahtev_id, zpr.zadatak_id, zpr.podnosilac_zahteva_id, zpr.zahtevana_faza_id,
		       zpr.status, zpr.komentar, zpr.datum_kreiranja,
		       COALESCE(k.ime || ' ' || k.prezime, k.korisnicko_ime) as podnosilac_ime,
		       f.naziv_faze as zahtevana_faza_naziv
		FROM ZahteviPromeneFaze zpr
		JOIN Korisnici k ON zpr.podnosilac_zahteva_id = k.korisnik_id
		JOIN Faze f ON zpr.zahtevana_faza_id = f.faza_id
		WHERE zpr.zadatak_id = $1
		ORDER BY zpr.datum_kreiranja DESC
	`

	rows, err := s.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("error querying phase change requests: %v", err)
	}
	defer rows.Close()

	var requests []models.ZahteviPromeneFaze
	for rows.Next() {
		var request models.ZahteviPromeneFaze
		var podnosilacIme, zahtevanaFazaNaziv string

		err := rows.Scan(
			&request.ZahtevID,
			&request.ZadatakID,
			&request.PodnosilacZahtevaID,
			&request.ZahtevanaFazaID,
			&request.Status,
			&request.Komentar,
			&request.DatumKreiranja,
			&podnosilacIme,
			&zahtevanaFazaNaziv,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning phase change request: %v", err)
		}

		// Note: You might want to extend the model to include these fields
		// For now, they're available but not stored in the struct
		requests = append(requests, request)
	}

	return requests, nil
}

// GetRequestByID retrieves a specific phase change request by ID
func (s *PhaseChangeRequestService) GetRequestByID(requestID int) (*models.ZahteviPromeneFaze, error) {
	query := `
		SELECT zpr.zahtev_id, zpr.zadatak_id, zpr.podnosilac_zahteva_id, zpr.zahtevana_faza_id,
		       zpr.status, zpr.komentar, zpr.datum_kreiranja
		FROM ZahteviPromeneFaze zpr
		WHERE zpr.zahtev_id = $1
	`

	var request models.ZahteviPromeneFaze
	err := s.db.QueryRow(query, requestID).Scan(
		&request.ZahtevID,
		&request.ZadatakID,
		&request.PodnosilacZahtevaID,
		&request.ZahtevanaFazaID,
		&request.Status,
		&request.Komentar,
		&request.DatumKreiranja,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("phase change request not found")
		}
		return nil, fmt.Errorf("error querying phase change request: %v", err)
	}

	return &request, nil
}

// CreateRequest creates a new phase change request
func (s *PhaseChangeRequestService) CreateRequest(taskID int, requesterID int, requestedPhaseID int, comment string) (*models.ZahteviPromeneFaze, error) {
	// First verify the task exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Zadaci WHERE zadatak_id = $1)", taskID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking task existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("task with ID %d does not exist", taskID)
	}

	// Verify the requested phase exists
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Faze WHERE faza_id = $1)", requestedPhaseID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking phase existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("phase with ID %d does not exist", requestedPhaseID)
	}

	// Check if there's already a pending request for this task
	var pendingExists bool
	err = s.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM ZahteviPromeneFaze 
		WHERE zadatak_id = $1 AND status = 'Na cekanju')
	`, taskID).Scan(&pendingExists)
	if err != nil {
		return nil, fmt.Errorf("error checking pending requests: %v", err)
	}
	if pendingExists {
		return nil, fmt.Errorf("there is already a pending phase change request for this task")
	}

	query := `
		INSERT INTO ZahteviPromeneFaze (zadatak_id, podnosilac_zahteva_id, zahtevana_faza_id, status, komentar, datum_kreiranja)
		VALUES ($1, $2, $3, 'Na cekanju', $4, CURRENT_TIMESTAMP)
		RETURNING zahtev_id, datum_kreiranja
	`

	var request models.ZahteviPromeneFaze
	request.ZadatakID = taskID
	request.PodnosilacZahtevaID = requesterID
	request.ZahtevanaFazaID = requestedPhaseID
	request.Status = "Na cekanju"
	request.Komentar = &comment

	err = s.db.QueryRow(query, taskID, requesterID, requestedPhaseID, comment).Scan(
		&request.ZahtevID,
		&request.DatumKreiranja,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating phase change request: %v", err)
	}

	return &request, nil
}

// UpdateRequestStatus updates the status of a phase change request (approve/reject)
func (s *PhaseChangeRequestService) UpdateRequestStatus(requestID int, status string, reviewerID int) error {
	// Validate status
	if status != "Odobren" && status != "Odbijen" {
		return fmt.Errorf("invalid status: must be 'Odobren' or 'Odbijen'")
	}

	// Get the request details first
	request, err := s.GetRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != "Na cekanju" {
		return fmt.Errorf("request has already been processed")
	}

	// Start transaction for atomic update
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Update the request status
	query := `UPDATE ZahteviPromeneFaze SET status = $1 WHERE zahtev_id = $2`
	_, err = tx.Exec(query, status, requestID)
	if err != nil {
		return fmt.Errorf("error updating request status: %v", err)
	}

	// If approved, update the task's phase
	if status == "Odobren" {
		taskUpdateQuery := `UPDATE Zadaci SET faza_id = $1 WHERE zadatak_id = $2`
		_, err = tx.Exec(taskUpdateQuery, request.ZahtevanaFazaID, request.ZadatakID)
		if err != nil {
			return fmt.Errorf("error updating task phase: %v", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// GetPendingRequests retrieves all pending phase change requests
func (s *PhaseChangeRequestService) GetPendingRequests() ([]models.ZahteviPromeneFaze, error) {
	query := `
		SELECT zpr.zahtev_id, zpr.zadatak_id, zpr.podnosilac_zahteva_id, zpr.zahtevana_faza_id,
		       zpr.status, zpr.komentar, zpr.datum_kreiranja
		FROM ZahteviPromeneFaze zpr
		WHERE zpr.status = 'Na cekanju'
		ORDER BY zpr.datum_kreiranja ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying pending requests: %v", err)
	}
	defer rows.Close()

	var requests []models.ZahteviPromeneFaze
	for rows.Next() {
		var request models.ZahteviPromeneFaze
		err := rows.Scan(
			&request.ZahtevID,
			&request.ZadatakID,
			&request.PodnosilacZahtevaID,
			&request.ZahtevanaFazaID,
			&request.Status,
			&request.Komentar,
			&request.DatumKreiranja,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning pending request: %v", err)
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// GetRequestsByUser retrieves all phase change requests made by a specific user
func (s *PhaseChangeRequestService) GetRequestsByUser(userID int) ([]models.ZahteviPromeneFaze, error) {
	query := `
		SELECT zpr.zahtev_id, zpr.zadatak_id, zpr.podnosilac_zahteva_id, zpr.zahtevana_faza_id,
		       zpr.status, zpr.komentar, zpr.datum_kreiranja
		FROM ZahteviPromeneFaze zpr
		WHERE zpr.podnosilac_zahteva_id = $1
		ORDER BY zpr.datum_kreiranja DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user requests: %v", err)
	}
	defer rows.Close()

	var requests []models.ZahteviPromeneFaze
	for rows.Next() {
		var request models.ZahteviPromeneFaze
		err := rows.Scan(
			&request.ZahtevID,
			&request.ZadatakID,
			&request.PodnosilacZahtevaID,
			&request.ZahtevanaFazaID,
			&request.Status,
			&request.Komentar,
			&request.DatumKreiranja,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user request: %v", err)
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// DeleteRequest removes a phase change request (only if pending and user is the requester or admin)
func (s *PhaseChangeRequestService) DeleteRequest(requestID int, userID int, isAdmin bool) error {
	request, err := s.GetRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != "Na cekanju" {
		return fmt.Errorf("can only delete pending requests")
	}

	// Check if user is the requester or admin
	if !isAdmin && request.PodnosilacZahtevaID != userID {
		return fmt.Errorf("unauthorized: can only delete own requests")
	}

	query := `DELETE FROM ZahteviPromeneFaze WHERE zahtev_id = $1`

	result, err := s.db.Exec(query, requestID)
	if err != nil {
		return fmt.Errorf("error deleting phase change request: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("phase change request not found")
	}

	return nil
}