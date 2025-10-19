// ============================================================================
// project_service.go - Project Management Service
// ============================================================================

package services

import (
	"database/sql"
	"fmt"

	"github.com/cane/research-institute-system/backend/models"
)

type ProjectService struct {
	db *sql.DB
}

func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) GetAllProjects() ([]models.Projekti, error) {
	query := `
		SELECT p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka,
		       p.datum_zavrsetka, p.status, p.rukovodilac_id, p.radni_tok_id, p.resursi,
		       COALESCE(k.korisnicko_ime, '') as rukovodilac_ime,
		       COALESCE(task_count.cnt, 0) as broj_zadataka,
		       COALESCE(member_count.cnt, 0) as broj_clanova
		FROM projekti p
		LEFT JOIN korisnici k ON p.rukovodilac_id = k.korisnik_id
		LEFT JOIN (
			SELECT projekat_id, COUNT(*) as cnt 
			FROM zadaci 
			GROUP BY projekat_id
		) task_count ON p.projekat_id = task_count.projekat_id
		LEFT JOIN (
			SELECT projekat_id, COUNT(*) as cnt 
			FROM clanoviprojekta 
			GROUP BY projekat_id
		) member_count ON p.projekat_id = member_count.projekat_id
		ORDER BY p.projekat_id DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Projekti
	for rows.Next() {
		var project models.Projekti
		err := rows.Scan(
			&project.ProjekatID, &project.NazivProjekta, &project.Opis,
			&project.DatumPocetka, &project.DatumZavrsetka, &project.Status,
			&project.RukovodilaID, &project.RadniTokID, &project.Resursi, &project.RukovodilaIme,
			&project.BrojZadataka, &project.BrojClanova,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (s *ProjectService) GetProjectByID(projectID int) (models.Projekti, error) {
	var project models.Projekti
	query := `
		SELECT p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka,
		       p.datum_zavrsetka, p.status, p.rukovodilac_id, p.radni_tok_id, p.resursi,
		       COALESCE(k.korisnicko_ime, '') as rukovodilac_ime
		FROM projekti p
		LEFT JOIN korisnici k ON p.rukovodilac_id = k.korisnik_id
		WHERE p.projekat_id = $1
	`

	err := s.db.QueryRow(query, projectID).Scan(
		&project.ProjekatID, &project.NazivProjekta, &project.Opis,
		&project.DatumPocetka, &project.DatumZavrsetka, &project.Status,
		&project.RukovodilaID, &project.RadniTokID, &project.Resursi, &project.RukovodilaIme,
	)

	return project, err
}

func (s *ProjectService) CreateProject(req models.CreateProjectRequest) error {
	return s.CreateProjectWithManager(req, 0)
}

func (s *ProjectService) CreateProjectWithManager(req models.CreateProjectRequest, managerID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert project with manager
	var projectID int
	query := `
		INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, radni_tok_id, rukovodilac_id, resursi, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'aktivan')
		RETURNING projekat_id
	`

	// Use manager ID if provided, otherwise NULL
	var managerIDPtr *int
	if managerID > 0 {
		managerIDPtr = &managerID
	}

	err = tx.QueryRow(query, req.NazivProjekta, req.Opis, req.DatumPocetka,
		req.DatumZavrsetka, req.RadniTokID, managerIDPtr, req.Resursi).Scan(&projectID)
	if err != nil {
		return err
	}

	// Add team members
	for _, memberID := range req.ClanoviTima {
		memberQuery := `INSERT INTO clanoviprojekta (projekat_id, korisnik_id) VALUES ($1, $2)`
		_, err = tx.Exec(memberQuery, projectID, memberID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *ProjectService) UpdateProject(projectID int, project models.Projekti) error {
	query := `
		UPDATE projekti 
		SET naziv_projekta = $1, opis = $2, datum_pocetka = $3, 
		    datum_zavrsetka = $4, status = $5, rukovodilac_id = $6, radni_tok_id = $7, resursi = $8
		WHERE projekat_id = $9
	`

	_, err := s.db.Exec(query, project.NazivProjekta, project.Opis,
		project.DatumPocetka, project.DatumZavrsetka, project.Status,
		project.RukovodilaID, project.RadniTokID, project.Resursi, projectID)

	return err
}

func (s *ProjectService) DeleteProject(projectID int) error {
	query := `DELETE FROM projekti WHERE projekat_id = $1`
	result, err := s.db.Exec(query, projectID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project with ID %d not found", projectID)
	}

	return nil
}

func (s *ProjectService) GetProjectMembers(projectID int) ([]models.Korisnici, error) {
	query := `
		SELECT k.korisnik_id, k.korisnicko_ime, k.email, k.ime, k.prezime,
			   k.uloga_id, u.naziv_uloge
		FROM korisnici k
		JOIN uloge u ON k.uloga_id = u.uloga_id
		JOIN clanoviprojekta cp ON k.korisnik_id = cp.korisnik_id
		WHERE cp.projekat_id = $1
		ORDER BY k.korisnicko_ime
	`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Korisnici
	for rows.Next() {
		var member models.Korisnici
		err := rows.Scan(
			&member.KorisnikID, &member.KorisnickoIme, &member.Email,
			&member.Ime, &member.Prezime, &member.UlogaID, &member.NazivUloge,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (s *ProjectService) AddProjectMember(projectID, userID int) error {
	query := `INSERT INTO clanoviprojekta (projekat_id, korisnik_id) VALUES ($1, $2)`
	_, err := s.db.Exec(query, projectID, userID)
	return err
}

func (s *ProjectService) RemoveProjectMember(projectID, userID int) error {
	query := `DELETE FROM clanoviprojekta WHERE projekat_id = $1 AND korisnik_id = $2`
	_, err := s.db.Exec(query, projectID, userID)
	return err
}

// GetProjectsByStatus returns projects filtered by status
func (s *ProjectService) GetProjectsByStatus(status string) ([]models.Projekti, error) {
	query := `
		SELECT p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka,
		       p.datum_zavrsetka, p.status, p.rukovodilac_id, p.radni_tok_id, p.resursi,
		       COALESCE(k.korisnicko_ime, '') as rukovodilac_ime,
		       COALESCE(task_count.cnt, 0) as broj_zadataka,
		       COALESCE(member_count.cnt, 0) as broj_clanova
		FROM projekti p
		LEFT JOIN korisnici k ON p.rukovodilac_id = k.korisnik_id
		LEFT JOIN (
			SELECT projekat_id, COUNT(*) as cnt 
			FROM zadaci 
			GROUP BY projekat_id
		) task_count ON p.projekat_id = task_count.projekat_id
		LEFT JOIN (
			SELECT projekat_id, COUNT(*) as cnt 
			FROM clanoviprojekta 
			GROUP BY projekat_id
		) member_count ON p.projekat_id = member_count.projekat_id
		WHERE p.status = $1
		ORDER BY p.projekat_id DESC
	`

	rows, err := s.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Projekti
	for rows.Next() {
		var project models.Projekti
		err := rows.Scan(
			&project.ProjekatID, &project.NazivProjekta, &project.Opis,
			&project.DatumPocetka, &project.DatumZavrsetka, &project.Status,
			&project.RukovodilaID, &project.RadniTokID, &project.Resursi, &project.RukovodilaIme,
			&project.BrojZadataka, &project.BrojClanova,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// CompleteProject marks a project as completed
func (s *ProjectService) CompleteProject(projectID int) error {
	query := `UPDATE projekti SET status = 'završen' WHERE projekat_id = $1`
	_, err := s.db.Exec(query, projectID)
	return err
}

// GetProjectResources returns resources allocated to a project
func (s *ProjectService) GetProjectResources(projectID int) ([]map[string]interface{}, error) {
	// This could query a resources table if it exists
	// For now, return project members as resources
	query := `
		SELECT k.korisnik_id, k.korisnicko_ime, k.ime, k.prezime, u.naziv_uloge
		FROM korisnici k
		JOIN uloge u ON k.uloga_id = u.uloga_id
		JOIN clanoviprojekta cp ON k.korisnik_id = cp.korisnik_id
		WHERE cp.projekat_id = $1
		ORDER BY k.korisnicko_ime
	`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []map[string]interface{}
	for rows.Next() {
		var id int
		var username, ime, prezime, uloga string
		err := rows.Scan(&id, &username, &ime, &prezime, &uloga)
		if err != nil {
			return nil, err
		}
		resources = append(resources, map[string]interface{}{
			"korisnik_id":     id,
			"korisnicko_ime":  username,
			"ime":             ime,
			"prezime":         prezime,
			"uloga":           uloga,
		})
	}

	return resources, nil
}

// GetProjectAnalytics returns comprehensive analytics for a project
func (s *ProjectService) GetProjectAnalytics(projectID int) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Get project details
	project, err := s.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	analytics["project"] = project

	// Get task statistics
	var totalTasks, completedTasks, inProgressTasks, pendingTasks int
	statsQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN progres = 100 THEN 1 END) as completed,
			COUNT(CASE WHEN progres > 0 AND progres < 100 THEN 1 END) as in_progress,
			COUNT(CASE WHEN progres = 0 THEN 1 END) as pending
		FROM zadaci
		WHERE projekat_id = $1
	`
	err = s.db.QueryRow(statsQuery, projectID).Scan(&totalTasks, &completedTasks, &inProgressTasks, &pendingTasks)
	if err != nil {
		return nil, err
	}

	analytics["total_tasks"] = totalTasks
	analytics["completed_tasks"] = completedTasks
	analytics["in_progress_tasks"] = inProgressTasks
	analytics["pending_tasks"] = pendingTasks
	if totalTasks > 0 {
		analytics["completion_percentage"] = float64(completedTasks) / float64(totalTasks) * 100
	} else {
		analytics["completion_percentage"] = 0
	}

	// Get tasks by phase
	phaseQuery := `
		SELECT f.naziv_faze, COUNT(*) as task_count
		FROM zadaci z
		JOIN faze f ON z.faza_id = f.faza_id
		WHERE z.projekat_id = $1
		GROUP BY f.naziv_faze, f.redosled
		ORDER BY f.redosled
	`
	rows, err := s.db.Query(phaseQuery, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasksByPhase := make([]map[string]interface{}, 0)
	for rows.Next() {
		var phaseName string
		var taskCount int
		err := rows.Scan(&phaseName, &taskCount)
		if err != nil {
			return nil, err
		}
		tasksByPhase = append(tasksByPhase, map[string]interface{}{
			"phase_name": phaseName,
			"task_count": taskCount,
		})
	}
	analytics["tasks_by_phase"] = tasksByPhase

	// Get team members count
	var memberCount int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM clanoviprojekta WHERE projekat_id = $1`, projectID).Scan(&memberCount)
	if err != nil {
		return nil, err
	}
	analytics["team_member_count"] = memberCount

	// Get overdue tasks
	var overdueTasks int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM zadaci 
		WHERE projekat_id = $1 AND rok < NOW() AND progres < 100
	`, projectID).Scan(&overdueTasks)
	if err == nil {
		analytics["overdue_tasks"] = overdueTasks
	}

	return analytics, nil
}

// GetProjectsByUser returns all projects where user is member or manager
func (s *ProjectService) GetProjectsByUser(userID int) ([]models.Projekti, error) {
	query := `
		SELECT DISTINCT p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka,
		       p.datum_zavrsetka, p.status, p.rukovodilac_id, p.radni_tok_id, p.resursi,
		       COALESCE(k.korisnicko_ime, '') as rukovodilac_ime,
		       COALESCE(task_count.cnt, 0) as broj_zadataka,
		       COALESCE(member_count.cnt, 0) as broj_clanova
		FROM projekti p
		LEFT JOIN korisnici k ON p.rukovodilac_id = k.korisnik_id
		LEFT JOIN (
			SELECT projekat_id, COUNT(*) as cnt 
			FROM zadaci 
			GROUP BY projekat_id
		) task_count ON p.projekat_id = task_count.projekat_id
		LEFT JOIN (
			SELECT projekat_id, COUNT(*) as cnt 
			FROM clanoviprojekta 
			GROUP BY projekat_id
		) member_count ON p.projekat_id = member_count.projekat_id
		WHERE p.rukovodilac_id = $1 
		   OR EXISTS (SELECT 1 FROM clanoviprojekta WHERE projekat_id = p.projekat_id AND korisnik_id = $1)
		ORDER BY p.projekat_id DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Projekti
	for rows.Next() {
		var project models.Projekti
		err := rows.Scan(
			&project.ProjekatID, &project.NazivProjekta, &project.Opis,
			&project.DatumPocetka, &project.DatumZavrsetka, &project.Status,
			&project.RukovodilaID, &project.RadniTokID, &project.Resursi, &project.RukovodilaIme,
			&project.BrojZadataka, &project.BrojClanova,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}
