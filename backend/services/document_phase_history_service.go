// ============================================================================
// document_phase_history_service.go - IstorijaFazaDokumenta CRUD Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type DocumentPhaseHistoryService struct{ db *sql.DB }

func NewDocumentPhaseHistoryService(db *sql.DB) *DocumentPhaseHistoryService {
	return &DocumentPhaseHistoryService{db: db}
}

func (s *DocumentPhaseHistoryService) ListByDocument(documentID int) ([]models.IstorijaFazaDokumenta, error) {
	rows, err := s.db.Query(`SELECT istorija_id, dokument_id, prethodna_faza_id, nova_faza_id, korisnik_id, datum_promene FROM istorijafazadokumenta WHERE dokument_id = $1 ORDER BY datum_promene DESC`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.IstorijaFazaDokumenta
	for rows.Next() {
		var h models.IstorijaFazaDokumenta
		if err := rows.Scan(&h.IstorijaID, &h.DokumentID, &h.PrethodnaFazaID, &h.NovaFazaID, &h.KorisnikID, &h.DatumPromene); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, rows.Err()
}

func (s *DocumentPhaseHistoryService) Create(h *models.IstorijaFazaDokumenta) error {
	return s.db.QueryRow(`INSERT INTO istorijafazadokumenta (dokument_id, prethodna_faza_id, nova_faza_id, korisnik_id) VALUES ($1, $2, $3, $4) RETURNING istorija_id`, h.DokumentID, h.PrethodnaFazaID, h.NovaFazaID, h.KorisnikID).Scan(&h.IstorijaID)
}

func (s *DocumentPhaseHistoryService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM istorijafazadokumenta WHERE istorija_id = $1`, id)
	return err
}
