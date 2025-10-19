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
		INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id)
		VALUES ($1, $2, $3, $4, COALESCE($5,'Aktivan'), $6, $7)
		RETURNING projekat_id
	`
	return r.db.QueryRow(query, project.NazivProjekta, project.Opis, project.DatumPocetka, project.DatumZavrsetka, project.Status, project.RukovodilaID, project.RadniTokID).Scan(&project.ProjekatID)
}

func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	query := `
		SELECT p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka,
			   p.datum_zavrsetka, p.status, p.rukovodilac_id, p.radni_tok_id,
			   COALESCE(k.korisnicko_ime, '') as rukovodilac_ime
		FROM projekti p
		LEFT JOIN korisnici k ON p.rukovodilac_id = k.korisnik_id
		WHERE p.projekat_id = $1
	`
	var pr models.Project
	err := r.db.QueryRow(query, id).Scan(
		&pr.ProjekatID, &pr.NazivProjekta, &pr.Opis, &pr.DatumPocetka,
		&pr.DatumZavrsetka, &pr.Status, &pr.RukovodilaID, &pr.RadniTokID, &pr.RukovodilaIme,
	)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *ProjectRepository) GetByUserID(userID int) ([]models.Project, error) {
	query := `
		SELECT DISTINCT p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka,
			   p.datum_zavrsetka, p.status, p.rukovodilac_id, p.radni_tok_id,
			   COALESCE(k.korisnicko_ime, '') as rukovodilac_ime
		FROM projekti p
		LEFT JOIN korisnici k ON p.rukovodilac_id = k.korisnik_id
		LEFT JOIN clanoviprojekta cp ON p.projekat_id = cp.projekat_id
		WHERE p.rukovodilac_id = $1 OR cp.korisnik_id = $1
		ORDER BY p.projekat_id DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Project
	for rows.Next() {
		var pr models.Project
		if err := rows.Scan(&pr.ProjekatID, &pr.NazivProjekta, &pr.Opis, &pr.DatumPocetka, &pr.DatumZavrsetka, &pr.Status, &pr.RukovodilaID, &pr.RadniTokID, &pr.RukovodilaIme); err != nil {
			return nil, err
		}
		list = append(list, pr)
	}
	return list, rows.Err()
}

func (r *ProjectRepository) Update(project *models.Project) error {
	_, err := r.db.Exec(`
		UPDATE projekti SET naziv_projekta = $1, opis = $2, datum_pocetka = $3,
			datum_zavrsetka = $4, status = $5, rukovodilac_id = $6, radni_tok_id = $7
		WHERE projekat_id = $8
	`, project.NazivProjekta, project.Opis, project.DatumPocetka, project.DatumZavrsetka, project.Status, project.RukovodilaID, project.RadniTokID, project.ProjekatID)
	return err
}

func (r *ProjectRepository) GetMembers(projectID int) ([]models.User, error) {
	query := `
		SELECT k.korisnik_id, k.korisnicko_ime, k.email, k.hash_sifre, k.ime, k.prezime,
			   k.uloga_id, k.status, k.poslednja_prijava, k.kreiran_datuma,
			   u.naziv_uloge
		FROM korisnici k
		JOIN uloge u ON k.uloga_id = u.uloga_id
		JOIN clanoviprojekta cp ON k.korisnik_id = cp.korisnik_id
		WHERE cp.projekat_id = $1
		ORDER BY k.korisnicko_ime
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		var lastLogin sql.NullTime
		if err := rows.Scan(&u.KorisnikID, &u.KorisnickoIme, &u.Email, &u.HashSifre, &u.Ime, &u.Prezime, &u.UlogaID, &u.Status, &lastLogin, &u.KreiranDatuma, &u.NazivUloge); err != nil {
			return nil, err
		}
		if lastLogin.Valid {
			u.PoslednajaPrijava = &lastLogin.Time
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
