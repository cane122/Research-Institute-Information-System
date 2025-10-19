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

func (s *PhaseChangeRequestService) ListByTask(taskID int) ([]models.ZahteviPromeneFaze, error) {
	rows, err := s.db.Query(`SELECT zahtev_id, zadatak_id, podnosilac_zahteva_id, zahtevana_faza_id, status, komentar, datum_kreiranja FROM zahtevipromenefaze WHERE zadatak_id = $1 ORDER BY datum_kreiranja DESC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ZahteviPromeneFaze
	for rows.Next() {
		var z models.ZahteviPromeneFaze
		if err := rows.Scan(&z.ZahtevID, &z.ZadatakID, &z.PodnosilacZahtevaID, &z.ZahtevanaFazaID, &z.Status, &z.Komentar, &z.DatumKreiranja); err != nil {
			return nil, err
		}
		list = append(list, z)
	}
	return list, rows.Err()
}

func (s *PhaseChangeRequestService) Create(z *models.ZahteviPromeneFaze) error {
	return s.db.QueryRow(`INSERT INTO zahtevipromenefaze (zadatak_id, podnosilac_zahteva_id, zahtevana_faza_id, status, komentar) VALUES ($1, $2, $3, COALESCE($4,'Na cekanju'), $5) RETURNING zahtev_id`, z.ZadatakID, z.PodnosilacZahtevaID, z.ZahtevanaFazaID, z.Status, z.Komentar).Scan(&z.ZahtevID)
}

func (s *PhaseChangeRequestService) UpdateStatus(id int, status string, komentar *string) error {
	_, err := s.db.Exec(`UPDATE zahtevipromenefaze SET status = $1, komentar = COALESCE($2, komentar) WHERE zahtev_id = $3`, status, komentar, id)
	return err
}

func (s *PhaseChangeRequestService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM zahtevipromenefaze WHERE zahtev_id = $1`, id)
	return err
}
