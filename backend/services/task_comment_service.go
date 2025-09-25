// ============================================================================
// task_comment_service.go - Task Comments Service
// ============================================================================

package services

import (
	"database/sql"
	"fmt"

	"github.com/cane/research-institute-system/backend/models"
)

type TaskCommentService struct {
	db *sql.DB
}

func NewTaskCommentService(db *sql.DB) *TaskCommentService {
	return &TaskCommentService{db: db}
}

// GetCommentsByTask retrieves all comments for a specific task
func (s *TaskCommentService) GetCommentsByTask(taskID int) ([]models.KomentariZadataka, error) {
	query := `
		SELECT kz.komentar_id, kz.zadatak_id, kz.korisnik_id, kz.tekst_komentara, kz.datuma_kreiranja,
		       COALESCE(k.ime || ' ' || k.prezime, k.korisnicko_ime) as ime_korisnika
		FROM KomentariZadataka kz
		JOIN Korisnici k ON kz.korisnik_id = k.korisnik_id
		WHERE kz.zadatak_id = $1
		ORDER BY kz.datuma_kreiranja ASC
	`

	rows, err := s.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("error querying task comments: %v", err)
	}
	defer rows.Close()

	var comments []models.KomentariZadataka
	for rows.Next() {
		var comment models.KomentariZadataka
		err := rows.Scan(
			&comment.KomentarID,
			&comment.ZadatakID,
			&comment.KorisnikID,
			&comment.TekstKomentara,
			&comment.DatumaKreiranja,
			&comment.ImeKorisnika,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning task comment: %v", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// GetCommentByID retrieves a specific comment by ID
func (s *TaskCommentService) GetCommentByID(commentID int) (*models.KomentariZadataka, error) {
	query := `
		SELECT kz.komentar_id, kz.zadatak_id, kz.korisnik_id, kz.tekst_komentara, kz.datuma_kreiranja,
		       COALESCE(k.ime || ' ' || k.prezime, k.korisnicko_ime) as ime_korisnika
		FROM KomentariZadataka kz
		JOIN Korisnici k ON kz.korisnik_id = k.korisnik_id
		WHERE kz.komentar_id = $1
	`

	var comment models.KomentariZadataka
	err := s.db.QueryRow(query, commentID).Scan(
		&comment.KomentarID,
		&comment.ZadatakID,
		&comment.KorisnikID,
		&comment.TekstKomentara,
		&comment.DatumaKreiranja,
		&comment.ImeKorisnika,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("error querying comment: %v", err)
	}

	return &comment, nil
}

// CreateComment adds a new comment to a task
func (s *TaskCommentService) CreateComment(taskID int, userID int, text string) (*models.KomentariZadataka, error) {
	// First verify the task exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Zadaci WHERE zadatak_id = $1)", taskID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking task existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("task with ID %d does not exist", taskID)
	}

	query := `
		INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara, datuma_kreiranja)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING komentar_id, datuma_kreiranja
	`

	var comment models.KomentariZadataka
	comment.ZadatakID = taskID
	comment.KorisnikID = userID
	comment.TekstKomentara = text

	err = s.db.QueryRow(query, taskID, userID, text).Scan(
		&comment.KomentarID,
		&comment.DatumaKreiranja,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating comment: %v", err)
	}

	// Get user name for the response
	userQuery := `
		SELECT COALESCE(ime || ' ' || prezime, korisnicko_ime) as ime_korisnika
		FROM Korisnici WHERE korisnik_id = $1
	`
	err = s.db.QueryRow(userQuery, userID).Scan(&comment.ImeKorisnika)
	if err != nil {
		comment.ImeKorisnika = "Unknown User"
	}

	return &comment, nil
}

// UpdateComment updates an existing comment (only if user is the owner or admin)
func (s *TaskCommentService) UpdateComment(commentID int, userID int, text string, isAdmin bool) error {
	// Check if user is owner of the comment or admin
	if !isAdmin {
		var ownerID int
		err := s.db.QueryRow("SELECT korisnik_id FROM KomentariZadataka WHERE komentar_id = $1", commentID).Scan(&ownerID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("comment not found")
			}
			return fmt.Errorf("error checking comment ownership: %v", err)
		}
		if ownerID != userID {
			return fmt.Errorf("unauthorized: can only edit own comments")
		}
	}

	query := `
		UPDATE KomentariZadataka 
		SET tekst_komentara = $1
		WHERE komentar_id = $2
	`

	result, err := s.db.Exec(query, text, commentID)
	if err != nil {
		return fmt.Errorf("error updating comment: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// DeleteComment removes a comment (only if user is the owner or admin)
func (s *TaskCommentService) DeleteComment(commentID int, userID int, isAdmin bool) error {
	// Check if user is owner of the comment or admin
	if !isAdmin {
		var ownerID int
		err := s.db.QueryRow("SELECT korisnik_id FROM KomentariZadataka WHERE komentar_id = $1", commentID).Scan(&ownerID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("comment not found")
			}
			return fmt.Errorf("error checking comment ownership: %v", err)
		}
		if ownerID != userID {
			return fmt.Errorf("unauthorized: can only delete own comments")
		}
	}

	query := `DELETE FROM KomentariZadataka WHERE komentar_id = $1`

	result, err := s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("error deleting comment: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// GetCommentsByUser retrieves all comments made by a specific user
func (s *TaskCommentService) GetCommentsByUser(userID int) ([]models.KomentariZadataka, error) {
	query := `
		SELECT kz.komentar_id, kz.zadatak_id, kz.korisnik_id, kz.tekst_komentara, kz.datuma_kreiranja,
		       COALESCE(k.ime || ' ' || k.prezime, k.korisnicko_ime) as ime_korisnika
		FROM KomentariZadataka kz
		JOIN Korisnici k ON kz.korisnik_id = k.korisnik_id
		WHERE kz.korisnik_id = $1
		ORDER BY kz.datuma_kreiranja DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user comments: %v", err)
	}
	defer rows.Close()

	var comments []models.KomentariZadataka
	for rows.Next() {
		var comment models.KomentariZadataka
		err := rows.Scan(
			&comment.KomentarID,
			&comment.ZadatakID,
			&comment.KorisnikID,
			&comment.TekstKomentara,
			&comment.DatumaKreiranja,
			&comment.ImeKorisnika,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user comment: %v", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}