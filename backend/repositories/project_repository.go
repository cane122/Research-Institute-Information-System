package repositories

import (
	"database/sql"
	"github.com/cane/research-institute-system/backend/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *models.Project) error {
	query := `
		INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING projekat_id
	`
	
	err := r.db.QueryRow(query, project.Name, project.Description, project.StartDate,
		project.EndDate, project.Status, project.ManagerID, project.WorkflowID).Scan(&project.ID)
	
	return err
}

func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	query := `
		SELECT p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka, p.datum_zavrsetka, 
		       p.status, p.rukovodilac_id,
		       k.korisnicko_ime, k.ime, k.prezime
		FROM Projekti p
		JOIN Korisnici k ON p.rukovodilac_id = k.korisnik_id
		WHERE p.projekat_id = $1
	`
	
	var project models.Project
	var manager models.User
	
	err := r.db.QueryRow(query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.StartDate,
		&project.EndDate, &project.Status, &project.ManagerID,
		&manager.Username, &manager.FirstName, &manager.LastName,
	)
	
	if err != nil {
		return nil, err
	}
	
	manager.ID = project.ManagerID
	project.Manager = &manager
	
	return &project, nil
}

func (r *ProjectRepository) GetByUserID(userID int) ([]models.Project, error) {
	query := `
		SELECT DISTINCT p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka, 
		       p.datum_zavrsetka, p.status, p.rukovodilac_id,
		       k.korisnicko_ime, k.ime, k.prezime
		FROM Projekti p
		JOIN Korisnici k ON p.rukovodilac_id = k.korisnik_id
		LEFT JOIN ClanoviProjekta cp ON p.projekat_id = cp.projekat_id
		WHERE p.rukovodilac_id = $1 OR cp.korisnik_id = $1
		ORDER BY p.projekat_id DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var projects []models.Project
	
	for rows.Next() {
		var project models.Project
		var manager models.User
		
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.StartDate,
			&project.EndDate, &project.Status, &project.ManagerID,
			&manager.Username, &manager.FirstName, &manager.LastName,
		)
		
		if err != nil {
			return nil, err
		}
		
		manager.ID = project.ManagerID
		project.Manager = &manager
		
		projects = append(projects, project)
	}
	
	return projects, nil
}

func (r *ProjectRepository) Update(project *models.Project) error {
	query := `
		UPDATE Projekti 
		SET naziv_projekta = $1, opis = $2, datum_pocetka = $3, datum_zavrsetka = $4, 
		    status = $5, radni_tok_id = $6
		WHERE projekat_id = $7
	`
	
	_, err := r.db.Exec(query, project.Name, project.Description, project.StartDate,
		project.EndDate, project.Status, project.WorkflowID, project.ID)
	
	return err
}

func (r *ProjectRepository) AddMember(projectID, userID int) error {
	query := `INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, projectID, userID)
	return err
}

func (r *ProjectRepository) RemoveMember(projectID, userID int) error {
	query := `DELETE FROM ClanoviProjekta WHERE projekat_id = $1 AND korisnik_id = $2`
	_, err := r.db.Exec(query, projectID, userID)
	return err
}

func (r *ProjectRepository) GetMembers(projectID int) ([]models.User, error) {
	query := `
		SELECT k.korisnik_id, k.korisnicko_ime, k.email, k.ime, k.prezime, k.uloga_id,
		       u.naziv_uloge
		FROM ClanoviProjekta cp
		JOIN Korisnici k ON cp.korisnik_id = k.korisnik_id
		JOIN Uloge u ON k.uloga_id = u.uloga_id
		WHERE cp.projekat_id = $1
	`
	
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var members []models.User
	
	for rows.Next() {
		var user models.User
		var role models.Role
		
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName,
			&user.LastName, &user.RoleID, &role.Name)
		
		if err != nil {
			return nil, err
		}
		
		role.ID = user.RoleID
		user.Role = &role
		
		members = append(members, user)
	}
	
	return members, nil
}
