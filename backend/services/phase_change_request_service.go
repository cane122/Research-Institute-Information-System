// ============================================================================
// phase_change_request_service.go - ZahteviPromeneFaze CRUD Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type PhaseChangeRequestService struct{ db *sql.DB }

func NewPhaseChangeRequestService(db *sql.DB) *PhaseChangeRequestService {
	return &PhaseChangeRequestService{db: db}
}

// List all phase change requests for a given task
func (s *PhaseChangeRequestService) ListByTask(taskID int) ([]models.ZahteviPromeneFaze, error) {
	rows, err := s.db.Query(`SELECT zahtev_id, zadatak_id, dokument_id, podnosilac_zahteva_id, zahtevana_faza_id, status, komentar, datum_kreiranja FROM zahtevipromenefaze WHERE zadatak_id = $1 ORDER BY datum_kreiranja DESC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ZahteviPromeneFaze
	for rows.Next() {
		var z models.ZahteviPromeneFaze
		if err := rows.Scan(&z.ZahtevID, &z.ZadatakID, &z.DokumentID, &z.PodnosilacZahtevaID, &z.ZahtevanaFazaID, &z.Status, &z.Komentar, &z.DatumKreiranja); err != nil {
			return nil, err
		}
		list = append(list, z)
	}
	return list, rows.Err()
}

// List all phase change requests for a given document
func (s *PhaseChangeRequestService) ListByDocument(documentID int) ([]models.ZahteviPromeneFaze, error) {
	rows, err := s.db.Query(`SELECT zahtev_id, zadatak_id, dokument_id, podnosilac_zahteva_id, zahtevana_faza_id, status, komentar, datum_kreiranja FROM zahtevipromenefaze WHERE dokument_id = $1 ORDER BY datum_kreiranja DESC`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ZahteviPromeneFaze
	for rows.Next() {
		var z models.ZahteviPromeneFaze
		if err := rows.Scan(&z.ZahtevID, &z.ZadatakID, &z.DokumentID, &z.PodnosilacZahtevaID, &z.ZahtevanaFazaID, &z.Status, &z.Komentar, &z.DatumKreiranja); err != nil {
			return nil, err
		}
		list = append(list, z)
	}
	return list, rows.Err()
}

// Create a new phase change request (for task or document)
func (s *PhaseChangeRequestService) Create(z *models.ZahteviPromeneFaze) error {
	return s.db.QueryRow(`INSERT INTO zahtevipromenefaze (zadatak_id, dokument_id, podnosilac_zahteva_id, zahtevana_faza_id, status, komentar) VALUES ($1, $2, $3, $4, COALESCE($5,'Na cekanju'), $6) RETURNING zahtev_id`,
		z.ZadatakID, z.DokumentID, z.PodnosilacZahtevaID, z.ZahtevanaFazaID, z.Status, z.Komentar).Scan(&z.ZahtevID)
}

// Get a phase change request by ID
func (s *PhaseChangeRequestService) GetByID(id int) (*models.ZahteviPromeneFaze, error) {
	var z models.ZahteviPromeneFaze
	err := s.db.QueryRow(`SELECT zahtev_id, zadatak_id, dokument_id, podnosilac_zahteva_id, zahtevana_faza_id, status, komentar, datum_kreiranja FROM zahtevipromenefaze WHERE zahtev_id = $1`, id).
		Scan(&z.ZahtevID, &z.ZadatakID, &z.DokumentID, &z.PodnosilacZahtevaID, &z.ZahtevanaFazaID, &z.Status, &z.Komentar, &z.DatumKreiranja)
	if err != nil {
		return nil, err
	}
	return &z, nil
}

func (s *PhaseChangeRequestService) UpdateStatus(id int, status string, komentar *string) error {
	_, err := s.db.Exec(`UPDATE zahtevipromenefaze SET status = $1, komentar = COALESCE($2, komentar) WHERE zahtev_id = $3`, status, komentar, id)
	return err
}

func (s *PhaseChangeRequestService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM zahtevipromenefaze WHERE zahtev_id = $1`, id)
	return err
}
