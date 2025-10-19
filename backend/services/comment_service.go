// ============================================================================
// comment_service.go - KomentariZadataka CRUD Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type CommentService struct{ db *sql.DB }

func NewCommentService(db *sql.DB) *CommentService { return &CommentService{db: db} }

func (s *CommentService) ListByTask(taskID int) ([]models.KomentariZadataka, error) {
	rows, err := s.db.Query(`SELECT komentar_id, zadatak_id, korisnik_id, tekst_komentara, datuma_kreiranja FROM komentarizadataka WHERE zadatak_id = $1 ORDER BY datuma_kreiranja DESC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.KomentariZadataka
	for rows.Next() {
		var c models.KomentariZadataka
		if err := rows.Scan(&c.KomentarID, &c.ZadatakID, &c.KorisnikID, &c.TekstKomentara, &c.DatumaKreiranja); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (s *CommentService) Create(c *models.KomentariZadataka) error {
	return s.db.QueryRow(`INSERT INTO komentarizadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES ($1, $2, $3) RETURNING komentar_id`, c.ZadatakID, c.KorisnikID, c.TekstKomentara).Scan(&c.KomentarID)
}

func (s *CommentService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM komentarizadataka WHERE komentar_id = $1`, id)
	return err
}
