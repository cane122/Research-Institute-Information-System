// ============================================================================
// metadata_service.go - MetaPodaci CRUD Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type MetadataService struct{ db *sql.DB }

func NewMetadataService(db *sql.DB) *MetadataService { return &MetadataService{db: db} }

func (s *MetadataService) ListByDocument(documentID int) ([]models.MetaPodaci, error) {
	rows, err := s.db.Query(`SELECT meta_id, dokument_id, kljuc, vrednost FROM metapodaci WHERE dokument_id = $1 ORDER BY kljuc`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.MetaPodaci
	for rows.Next() {
		var m models.MetaPodaci
		if err := rows.Scan(&m.MetaID, &m.DokumentID, &m.Kljuc, &m.Vrednost); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (s *MetadataService) GetByID(id int) (*models.MetaPodaci, error) {
	var m models.MetaPodaci
	err := s.db.QueryRow(`SELECT meta_id, dokument_id, kljuc, vrednost FROM metapodaci WHERE meta_id = $1`, id).Scan(&m.MetaID, &m.DokumentID, &m.Kljuc, &m.Vrednost)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *MetadataService) Create(m *models.MetaPodaci) error {
	return s.db.QueryRow(`INSERT INTO metapodaci (dokument_id, kljuc, vrednost) VALUES ($1, $2, $3) RETURNING meta_id`, m.DokumentID, m.Kljuc, m.Vrednost).Scan(&m.MetaID)
}

func (s *MetadataService) Update(m models.MetaPodaci) error {
	_, err := s.db.Exec(`UPDATE metapodaci SET kljuc = $1, vrednost = $2 WHERE meta_id = $3`, m.Kljuc, m.Vrednost, m.MetaID)
	return err
}

func (s *MetadataService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM metapodaci WHERE meta_id = $1`, id)
	return err
}
