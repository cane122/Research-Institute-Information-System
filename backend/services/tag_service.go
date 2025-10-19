// ============================================================================
// tag_service.go - Tagovi CRUD Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type TagService struct{ db *sql.DB }

func NewTagService(db *sql.DB) *TagService { return &TagService{db: db} }

func (s *TagService) List() ([]models.Tagovi, error) {
	rows, err := s.db.Query(`SELECT tag_id, naziv_taga FROM tagovi ORDER BY naziv_taga`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Tagovi
	for rows.Next() {
		var t models.Tagovi
		if err := rows.Scan(&t.TagID, &t.NazivTaga); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (s *TagService) GetByID(id int) (*models.Tagovi, error) {
	var t models.Tagovi
	err := s.db.QueryRow(`SELECT tag_id, naziv_taga FROM tagovi WHERE tag_id = $1`, id).Scan(&t.TagID, &t.NazivTaga)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TagService) Create(t *models.Tagovi) error {
	return s.db.QueryRow(`INSERT INTO tagovi (naziv_taga) VALUES ($1) RETURNING tag_id`, t.NazivTaga).Scan(&t.TagID)
}

func (s *TagService) Update(t models.Tagovi) error {
	_, err := s.db.Exec(`UPDATE tagovi SET naziv_taga = $1 WHERE tag_id = $2`, t.NazivTaga, t.TagID)
	return err
}

func (s *TagService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM tagovi WHERE tag_id = $1`, id)
	return err
}
