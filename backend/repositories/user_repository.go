package repositories

import (
	"database/sql"
	"time"

	"github.com/cane/research-institute-system/backend/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT k.korisnik_id, k.korisnicko_ime, k.email, k.hash_sifre, k.ime, k.prezime, 
		       k.uloga_id, k.status, k.poslednja_prijava, k.kreiran_datuma,
		       u.naziv_uloge
		FROM Korisnici k
		JOIN Uloge u ON k.uloga_id = u.uloga_id
		WHERE k.korisnik_id = $1
	`

	var user models.User
	var role models.Role
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&user.KorisnikID, &user.KorisnickoIme, &user.Email, &user.HashSifre,
		&user.Ime, &user.Prezime, &user.UlogaID, &user.Status,
		&lastLogin, &user.KreiranDatuma, &role.NazivUloge,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.PoslednajaPrijava = &lastLogin.Time
	}

	role.UlogaID = user.UlogaID
	user.NazivUloge = role.NazivUloge

	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT k.korisnik_id, k.korisnicko_ime, k.email, k.hash_sifre, k.ime, k.prezime, 
		       k.uloga_id, k.status, k.poslednja_prijava, k.kreiran_datuma,
		       u.naziv_uloge
		FROM Korisnici k
		JOIN Uloge u ON k.uloga_id = u.uloga_id
		WHERE k.korisnicko_ime = $1
	`

	var user models.User
	var role models.Role
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, username).Scan(
		&user.KorisnikID, &user.KorisnickoIme, &user.Email, &user.HashSifre,
		&user.Ime, &user.Prezime, &user.UlogaID, &user.Status,
		&lastLogin, &user.KreiranDatuma, &role.NazivUloge,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.PoslednajaPrijava = &lastLogin.Time
	}

	role.UlogaID = user.UlogaID
	user.NazivUloge = role.NazivUloge

	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING korisnik_id, kreiran_datuma
	`

	err := r.db.QueryRow(query, user.KorisnickoIme, user.Email, user.HashSifre,
		user.Ime, user.Prezime, user.UlogaID, user.Status).Scan(&user.KorisnikID, &user.KreiranDatuma)

	return err
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE Korisnici 
		SET korisnicko_ime = $1, email = $2, ime = $3, prezime = $4, uloga_id = $5, status = $6
		WHERE korisnik_id = $7
	`

	_, err := r.db.Exec(query, user.KorisnickoIme, user.Email, user.Ime,
		user.Prezime, user.UlogaID, user.Status, user.KorisnikID)

	return err
}

func (r *UserRepository) UpdatePassword(userID int, passwordHash string) error {
	query := `UPDATE Korisnici SET hash_sifre = $1 WHERE korisnik_id = $2`
	_, err := r.db.Exec(query, passwordHash, userID)
	return err
}

func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := `UPDATE Korisnici SET poslednja_prijava = $1 WHERE korisnik_id = $2`
	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
		SELECT k.korisnik_id, k.korisnicko_ime, k.email, k.ime, k.prezime, 
		       k.uloga_id, k.status, k.poslednja_prijava, k.kreiran_datuma,
		       u.naziv_uloge
		FROM Korisnici k
		JOIN Uloge u ON k.uloga_id = u.uloga_id
		ORDER BY k.kreiran_datuma DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var role models.Role
		var lastLogin sql.NullTime

		err := rows.Scan(
			&user.KorisnikID, &user.KorisnickoIme, &user.Email, &user.Ime,
			&user.Prezime, &user.UlogaID, &user.Status, &lastLogin,
			&user.KreiranDatuma, &role.NazivUloge,
		)

		if err != nil {
			return nil, err
		}

		if lastLogin.Valid {
			user.PoslednajaPrijava = &lastLogin.Time
		}

		role.UlogaID = user.UlogaID
		user.NazivUloge = role.NazivUloge

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetRoles() ([]models.Role, error) {
	query := `SELECT uloga_id, naziv_uloge FROM Uloge ORDER BY uloga_id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.UlogaID, &role.NazivUloge)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
