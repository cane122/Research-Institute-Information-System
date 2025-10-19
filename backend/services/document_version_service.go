// ============================================================================
// document_version_service.go - VerzijeDokumenata CRUD Service (partial)
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type DocumentVersionService struct{ db *sql.DB }

func NewDocumentVersionService(db *sql.DB) *DocumentVersionService {
	return &DocumentVersionService{db: db}
}

func (s *DocumentVersionService) ListByDocument(documentID int) ([]models.VerzijeDokumenata, error) {
	rows, err := s.db.Query(`SELECT verzija_id, dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_mb, postavio_korisnik_id, datuma_postavke FROM verzijedokumenata WHERE dokument_id = $1 ORDER BY datuma_postavke DESC`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.VerzijeDokumenata
	for rows.Next() {
		var v models.VerzijeDokumenata
		if err := rows.Scan(&v.VerzijaID, &v.DokumentID, &v.VerzijaOznaka, &v.PutanjaDoFajla, &v.VelicinafajlaMB, &v.PostavioKorisnikID, &v.DatumaPostavke); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (s *DocumentVersionService) GetByID(id int) (*models.VerzijeDokumenata, error) {
	var v models.VerzijeDokumenata
	err := s.db.QueryRow(`SELECT verzija_id, dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_mb, postavio_korisnik_id, datuma_postavke FROM verzijedokumenata WHERE verzija_id = $1`, id).
		Scan(&v.VerzijaID, &v.DokumentID, &v.VerzijaOznaka, &v.PutanjaDoFajla, &v.VelicinafajlaMB, &v.PostavioKorisnikID, &v.DatumaPostavke)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *DocumentVersionService) Create(v *models.VerzijeDokumenata) error {
	return s.db.QueryRow(`INSERT INTO verzijedokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_mb, postavio_korisnik_id) VALUES ($1, $2, $3, $4, $5) RETURNING verzija_id`, v.DokumentID, v.VerzijaOznaka, v.PutanjaDoFajla, v.VelicinafajlaMB, v.PostavioKorisnikID).Scan(&v.VerzijaID)
}

func (s *DocumentVersionService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM verzijedokumenata WHERE verzija_id = $1`, id)
	return err
}
