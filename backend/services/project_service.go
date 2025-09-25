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
		       p.datum_zavrsetka, p.status, p.rukovodilac_id, p.radni_tok_id,
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
			&project.RukovodilaID, &project.RadniTokID, &project.RukovodilaIme,
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
		       p.datum_zavrsetka, p.status, p.rukovodilac_id, p.radni_tok_id,
		       COALESCE(k.korisnicko_ime, '') as rukovodilac_ime
		FROM projekti p
		LEFT JOIN korisnici k ON p.rukovodilac_id = k.korisnik_id
		WHERE p.projekat_id = $1
	`

	err := s.db.QueryRow(query, projectID).Scan(
		&project.ProjekatID, &project.NazivProjekta, &project.Opis,
		&project.DatumPocetka, &project.DatumZavrsetka, &project.Status,
		&project.RukovodilaID, &project.RadniTokID, &project.RukovodilaIme,
	)

	return project, err
}

func (s *ProjectService) CreateProject(req models.CreateProjectRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert project
	var projectID int
	query := `
		INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, radni_tok_id, status)
		VALUES ($1, $2, $3, $4, $5, 'aktivan')
		RETURNING projekat_id
	`

	err = tx.QueryRow(query, req.NazivProjekta, req.Opis, req.DatumPocetka,
		req.DatumZavrsetka, req.RadniTokID).Scan(&projectID)
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
		    datum_zavrsetka = $4, status = $5, rukovodilac_id = $6, radni_tok_id = $7
		WHERE projekat_id = $8
	`

	_, err := s.db.Exec(query, project.NazivProjekta, project.Opis,
		project.DatumPocetka, project.DatumZavrsetka, project.Status,
		project.RukovodilaID, project.RadniTokID, projectID)

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
