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
		       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.resursi, z.kreiran_datuma,
		       p.naziv_projekta, f.naziv_faze,
		       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
		FROM zadaci z
		JOIN projekti p ON z.projekat_id = p.projekat_id
		JOIN faze f ON z.faza_id = f.faza_id
		LEFT JOIN korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
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
			&task.Progres, &task.Resursi, &task.KreiranDatuma, &task.NazivProjekta,
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
		       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.resursi, z.kreiran_datuma,
		       p.naziv_projekta, f.naziv_faze,
		       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
		FROM zadaci z
		JOIN projekti p ON z.projekat_id = p.projekat_id
		JOIN faze f ON z.faza_id = f.faza_id
		LEFT JOIN korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
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
			&task.Progres, &task.Resursi, &task.KreiranDatuma, &task.NazivProjekta,
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
		       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.resursi, z.kreiran_datuma,
		       p.naziv_projekta, f.naziv_faze,
		       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
		FROM zadaci z
		JOIN projekti p ON z.projekat_id = p.projekat_id
		JOIN faze f ON z.faza_id = f.faza_id
		LEFT JOIN korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
		WHERE z.zadatak_id = $1
	`

	err := s.db.QueryRow(query, taskID).Scan(
		&task.ZadatakID, &task.ProjekatID, &task.FazaID, &task.NazivZadatka,
		&task.Opis, &task.DodjeljenKorisnikuID, &task.Rok, &task.Prioritet,
		&task.Progres, &task.Resursi, &task.KreiranDatuma, &task.NazivProjekta,
		&task.NazivFaze, &task.DodjeljenKorisniku,
	)

	return task, err
}

func (s *TaskService) CreateTask(req models.CreateTaskRequest) error {
	// Determine which phase to use
	var faseID int
	
	// If FazaID is provided, use it
	if req.FazaID != nil {
		faseID = *req.FazaID
	} else {
		// Otherwise, get first phase of project workflow
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
	}

	query := `
		INSERT INTO zadaci (projekat_id, faza_id, naziv_zadatka, opis, 
		                   dodeljen_korisniku_id, rok, prioritet, resursi)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.Exec(query, req.ProjekatID, faseID, req.NazivZadatka,
		req.Opis, req.DodjeljenKorisnikuID, req.Rok, req.Prioritet, req.Resursi)

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

	if req.Resursi != nil {
		setParts = append(setParts, fmt.Sprintf("resursi = $%d", argCount))
		args = append(args, *req.Resursi)
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

// MoveTaskToPhase moves a task to a different phase (workflow state)
func (s *TaskService) MoveTaskToPhase(taskID, newPhaseID int) error {
	query := `UPDATE zadaci SET faza_id = $1 WHERE zadatak_id = $2`
	_, err := s.db.Exec(query, newPhaseID, taskID)
	return err
}

// GetTasksByPhase returns all tasks in a specific phase
func (s *TaskService) GetTasksByPhase(phaseID int) ([]models.Zadaci, error) {
	query := `
		SELECT z.zadatak_id, z.projekat_id, z.faza_id, z.naziv_zadatka, z.opis,
		       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.resursi, z.kreiran_datuma,
		       p.naziv_projekta, f.naziv_faze,
		       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
		FROM zadaci z
		JOIN projekti p ON z.projekat_id = p.projekat_id
		JOIN faze f ON z.faza_id = f.faza_id
		LEFT JOIN korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
		WHERE z.faza_id = $1
		ORDER BY z.rok ASC NULLS LAST, z.prioritet DESC
	`

	rows, err := s.db.Query(query, phaseID)
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
			&task.Progres, &task.Resursi, &task.KreiranDatuma, &task.NazivProjekta,
			&task.NazivFaze, &task.DodjeljenKorisniku,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// RequestPhaseChange creates a phase change request for a task
func (s *TaskService) RequestPhaseChange(taskID, userID, requestedPhaseID int, comment string) error {
	query := `
		INSERT INTO zahtevipromenefaze (zadatak_id, podnosilac_zahteva_id, zahtevana_faza_id, status, komentar)
		VALUES ($1, $2, $3, 'na čekanju', $4)
	`

	_, err := s.db.Exec(query, taskID, userID, requestedPhaseID, comment)
	return err
}

// GetPhaseChangeRequests returns all phase change requests for a project or task
func (s *TaskService) GetPhaseChangeRequests(projectID *int, taskID *int) ([]models.ZahteviPromeneFaze, error) {
	var query string
	var args []interface{}

	if taskID != nil {
		query = `
			SELECT z.zahtev_id, z.zadatak_id, z.podnosilac_zahteva_id, z.zahtevana_faza_id,
			       z.status, z.komentar, z.datum_kreiranja
			FROM zahtevipromenefaze z
			WHERE z.zadatak_id = $1
			ORDER BY z.datum_kreiranja DESC
		`
		args = append(args, *taskID)
	} else if projectID != nil {
		query = `
			SELECT z.zahtev_id, z.zadatak_id, z.podnosilac_zahteva_id, z.zahtevana_faza_id,
			       z.status, z.komentar, z.datum_kreiranja
			FROM zahtevipromenefaze z
			JOIN zadaci zd ON z.zadatak_id = zd.zadatak_id
			WHERE zd.projekat_id = $1
			ORDER BY z.datum_kreiranja DESC
		`
		args = append(args, *projectID)
	} else {
		return nil, fmt.Errorf("either projectID or taskID must be provided")
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ZahteviPromeneFaze
	for rows.Next() {
		var request models.ZahteviPromeneFaze
		err := rows.Scan(
			&request.ZahtevID, &request.ZadatakID, &request.PodnosilacZahtevaID,
			&request.ZahtevanaFazaID, &request.Status, &request.Komentar,
			&request.DatumKreiranja,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// ApprovePhaseChangeRequest approves and executes a phase change request
func (s *TaskService) ApprovePhaseChangeRequest(requestID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get request details
	var taskID, newPhaseID int
	err = tx.QueryRow(`
		SELECT zadatak_id, zahtevana_faza_id 
		FROM zahtevipromenefaze 
		WHERE zahtev_id = $1
	`, requestID).Scan(&taskID, &newPhaseID)
	if err != nil {
		return err
	}

	// Update task phase
	_, err = tx.Exec(`UPDATE zadaci SET faza_id = $1 WHERE zadatak_id = $2`, newPhaseID, taskID)
	if err != nil {
		return err
	}

	// Update request status
	_, err = tx.Exec(`UPDATE zahtevipromenefaze SET status = 'odobren' WHERE zahtev_id = $1`, requestID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RejectPhaseChangeRequest rejects a phase change request
func (s *TaskService) RejectPhaseChangeRequest(requestID int) error {
	query := `UPDATE zahtevipromenefaze SET status = 'odbijen' WHERE zahtev_id = $1`
	_, err := s.db.Exec(query, requestID)
	return err
}

// GetOverdueTasks returns tasks that are past their deadline
func (s *TaskService) GetOverdueTasks(projectID *int) ([]models.Zadaci, error) {
	var query string
	var args []interface{}

	if projectID != nil {
		query = `
			SELECT z.zadatak_id, z.projekat_id, z.faza_id, z.naziv_zadatka, z.opis,
			       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.resursi, z.kreiran_datuma,
			       p.naziv_projekta, f.naziv_faze,
			       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
			FROM zadaci z
			JOIN projekti p ON z.projekat_id = p.projekat_id
			JOIN faze f ON z.faza_id = f.faza_id
			LEFT JOIN korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
			WHERE z.rok < NOW() AND z.progres < 100 AND z.projekat_id = $1
			ORDER BY z.rok ASC
		`
		args = append(args, *projectID)
	} else {
		query = `
			SELECT z.zadatak_id, z.projekat_id, z.faza_id, z.naziv_zadatka, z.opis,
			       z.dodeljen_korisniku_id, z.rok, z.prioritet, z.progres, z.resursi, z.kreiran_datuma,
			       p.naziv_projekta, f.naziv_faze,
			       COALESCE(k.korisnicko_ime, '') as dodeljen_korisniku
			FROM zadaci z
			JOIN projekti p ON z.projekat_id = p.projekat_id
			JOIN faze f ON z.faza_id = f.faza_id
			LEFT JOIN korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
			WHERE z.rok < NOW() AND z.progres < 100
			ORDER BY z.rok ASC
		`
	}

	rows, err := s.db.Query(query, args...)
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
			&task.Progres, &task.Resursi, &task.KreiranDatuma, &task.NazivProjekta,
			&task.NazivFaze, &task.DodjeljenKorisniku,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTaskProgress updates only the progress of a task
func (s *TaskService) UpdateTaskProgress(taskID, progress int) error {
	query := `UPDATE zadaci SET progres = $1 WHERE zadatak_id = $2`
	_, err := s.db.Exec(query, progress, taskID)
	return err
}
