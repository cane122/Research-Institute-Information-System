package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type ZadacicService struct {
	db *sql.DB
}

func NewZadacicService(db *sql.DB) *ZadacicService {
	return &ZadacicService{db: db}
}

func (s *ZadacicService) ListByDocument(dokumentID int) ([]models.Zadacic, error) {
	rows, err := s.db.Query(`SELECT zadacic_id, dokument_id, opis, izvrsen FROM zadacici WHERE dokument_id = $1 ORDER BY zadacic_id`, dokumentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Zadacic
	for rows.Next() {
		var z models.Zadacic
		if err := rows.Scan(&z.ZadacicID, &z.DokumentID, &z.Opis, &z.Izvrsen); err != nil {
			return nil, err
		}
		result = append(result, z)
	}
	return result, nil
}

func (s *ZadacicService) Add(z *models.Zadacic) error {
	return s.db.QueryRow(`INSERT INTO zadacici (dokument_id, opis, izvrsen) VALUES ($1, $2, $3) RETURNING zadacic_id`, z.DokumentID, z.Opis, z.Izvrsen).Scan(&z.ZadacicID)
}

func (s *ZadacicService) UpdateStatus(zadacicID int, izvrsen bool) error {
	_, err := s.db.Exec(`UPDATE zadacici SET izvrsen = $1 WHERE zadacic_id = $2`, izvrsen, zadacicID)
	return err
}

func (s *ZadacicService) Delete(zadacicID int) error {
	_, err := s.db.Exec(`DELETE FROM zadacici WHERE zadacic_id = $1`, zadacicID)
	return err
}
