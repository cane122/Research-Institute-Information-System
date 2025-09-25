// ============================================================================
// task_service.go - Task Management Service
// ============================================================================

package services

import (
	"database/sql"
	"fmt"

	"github.com/cane/research-institute-system/backend/models"
)

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) GetTasksByProject(projectID int) ([]models.Zadaci, error) {
	query := `
		SELECT z.zadatak_id, z.projekat_id, z.faza_id, z.naziv_zadatka, z.opis,
		       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.kreiran_datuma,
		       p.naziv_projekta, f.naziv_faze,
		       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
		FROM Zadaci z
		JOIN Projekti p ON z.projekat_id = p.projekat_id
		JOIN Faze f ON z.faza_id = f.faza_id
		LEFT JOIN Korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
		WHERE z.projekat_id = $1
		ORDER BY z.kreiran_datuma DESC
	`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Zadaci
	for rows.Next() {
		var task models.Zadaci
		err := rows.Scan(
			&task.ZadatakID, &task.ProjekatID, &task.FazaID, &task.NazivZadatka,
			&task.Opis, &task.DodjeljenKorisnikuID, &task.Rok, &task.Prioritet,
			&task.Progres, &task.KreiranDatuma, &task.NazivProjekta,
			&task.NazivFaze, &task.DodjeljenKorisniku,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *TaskService) GetTasksByUser(userID int) ([]models.Zadaci, error) {
	query := `
		SELECT z.zadatak_id, z.projekat_id, z.faza_id, z.naziv_zadatka, z.opis,
		       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.kreiran_datuma,
		       p.naziv_projekta, f.naziv_faze,
		       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
		FROM Zadaci z
		JOIN Projekti p ON z.projekat_id = p.projekat_id
		JOIN Faze f ON z.faza_id = f.faza_id
		LEFT JOIN Korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
		WHERE z.dodeljen_korisniku_id = $1
		ORDER BY z.rok ASC NULLS LAST, z.prioritet DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Zadaci
	for rows.Next() {
		var task models.Zadaci
		err := rows.Scan(
			&task.ZadatakID, &task.ProjekatID, &task.FazaID, &task.NazivZadatka,
			&task.Opis, &task.DodjeljenKorisnikuID, &task.Rok, &task.Prioritet,
			&task.Progres, &task.KreiranDatuma, &task.NazivProjekta,
			&task.NazivFaze, &task.DodjeljenKorisniku,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *TaskService) GetTaskByID(taskID int) (models.Zadaci, error) {
	var task models.Zadaci
	query := `
		SELECT z.zadatak_id, z.projekat_id, z.faza_id, z.naziv_zadatka, z.opis,
		       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.kreiran_datuma,
		       p.naziv_projekta, f.naziv_faze,
		       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
		FROM Zadaci z
		JOIN Projekti p ON z.projekat_id = p.projekat_id
		JOIN Faze f ON z.faza_id = f.faza_id
		LEFT JOIN Korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
		WHERE z.zadatak_id = $1
	`

	err := s.db.QueryRow(query, taskID).Scan(
		&task.ZadatakID, &task.ProjekatID, &task.FazaID, &task.NazivZadatka,
		&task.Opis, &task.DodjeljenKorisnikuID, &task.Rok, &task.Prioritet,
		&task.Progres, &task.KreiranDatuma, &task.NazivProjekta,
		&task.NazivFaze, &task.DodjeljenKorisniku,
	)

	return task, err
}

func (s *TaskService) CreateTask(req models.CreateTaskRequest) error {
	// Get first phase of project workflow
	var faseID int
	phaseQuery := `
		SELECT f.faza_id 
		FROM faze f
		JOIN projekti p ON f.radni_tok_id = p.radni_tok_id
		WHERE p.projekat_id = $1
		ORDER BY f.redosled ASC
		LIMIT 1
	`
	err := s.db.QueryRow(phaseQuery, req.ProjekatID).Scan(&faseID)
	if err != nil {
		// If no workflow, use default phase 1
		faseID = 1
	}

	query := `
		INSERT INTO zadaci (projekat_id, faza_id, naziv_zadatka, opis, 
		                   dodeljen_korisniku_id, rok, prioritet)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = s.db.Exec(query, req.ProjekatID, faseID, req.NazivZadatka,
		req.Opis, req.DodjeljenKorisnikuID, req.Rok, req.Prioritet)

	return err
}

func (s *TaskService) UpdateTask(taskID int, req models.UpdateTaskRequest) error {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	if req.NazivZadatka != nil {
		setParts = append(setParts, fmt.Sprintf("naziv_zadatka = $%d", argCount))
		args = append(args, *req.NazivZadatka)
		argCount++
	}

	if req.Opis != nil {
		setParts = append(setParts, fmt.Sprintf("opis = $%d", argCount))
		args = append(args, *req.Opis)
		argCount++
	}

	if req.DodjeljenKorisnikuID != nil {
		setParts = append(setParts, fmt.Sprintf("dodeljen_korisniku_id = $%d", argCount))
		args = append(args, *req.DodjeljenKorisnikuID)
		argCount++
	}

	if req.Rok != nil {
		setParts = append(setParts, fmt.Sprintf("rok = $%d", argCount))
		args = append(args, *req.Rok)
		argCount++
	}

	if req.Prioritet != nil {
		setParts = append(setParts, fmt.Sprintf("prioritet = $%d", argCount))
		args = append(args, *req.Prioritet)
		argCount++
	}

	if req.Progres != nil {
		setParts = append(setParts, fmt.Sprintf("progres = $%d", argCount))
		args = append(args, *req.Progres)
		argCount++
	}

	if req.FazaID != nil {
		setParts = append(setParts, fmt.Sprintf("faza_id = $%d", argCount))
		args = append(args, *req.FazaID)
		argCount++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE zadaci SET %s WHERE zadatak_id = $%d",
		fmt.Sprintf("%s", setParts[0]), argCount)
	for i := 1; i < len(setParts); i++ {
		query = query[:len(query)-len(fmt.Sprintf(" WHERE zadatak_id = $%d", argCount))] +
			fmt.Sprintf(", %s WHERE zadatak_id = $%d", setParts[i], argCount)
	}
	args = append(args, taskID)

	_, err := s.db.Exec(query, args...)
	return err
}

func (s *TaskService) DeleteTask(taskID int) error {
	query := `DELETE FROM zadaci WHERE zadatak_id = $1`
	result, err := s.db.Exec(query, taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", taskID)
	}

	return nil
}

func (s *TaskService) GetTaskComments(taskID int) ([]models.KomentariZadataka, error) {
	query := `
		SELECT kz.komentar_id, kz.zadatak_id, kz.korisnik_id, kz.tekst_komentara,
		       kz.datuma_kreiranja, k.korisnicko_ime as ime_korisnika
		FROM komentarizadataka kz
		JOIN korisnici k ON kz.korisnik_id = k.korisnik_id
		WHERE kz.zadatak_id = $1
		ORDER BY kz.datuma_kreiranja DESC
	`

	rows, err := s.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.KomentariZadataka
	for rows.Next() {
		var comment models.KomentariZadataka
		err := rows.Scan(
			&comment.KomentarID, &comment.ZadatakID, &comment.KorisnikID,
			&comment.TekstKomentara, &comment.DatumaKreiranja, &comment.ImeKorisnika,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *TaskService) AddTaskComment(taskID, userID int, comment string) error {
	query := `
		INSERT INTO komentarizadataka (zadatak_id, korisnik_id, tekst_komentara)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(query, taskID, userID, comment)
	return err
}
