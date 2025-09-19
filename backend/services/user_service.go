// ============================================================================
// user_service.go - User Management Service
// ============================================================================

package services

import (
	"database/sql"
	"fmt"

	"github.com/cane/research-institute-system/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetAllUsers() ([]models.Korisnici, error) {
	query := `
		SELECT k.korisnik_id, k.korisnicko_ime, k.email, k.ime, k.prezime,
		       k.uloga_id, k.status, k.poslednja_prijava, k.kreiran_datuma,
		       u.naziv_uloge
		FROM korisnici k
		JOIN uloge u ON k.uloga_id = u.uloga_id
		ORDER BY k.kreiran_datuma DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.Korisnici
	for rows.Next() {
		var user models.Korisnici
		err := rows.Scan(
			&user.KorisnikID, &user.KorisnickoIme, &user.Email,
			&user.Ime, &user.Prezime, &user.UlogaID, &user.Status,
			&user.PoslednajaPrijava, &user.KreiranDatuma, &user.NazivUloge,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) CreateUser(user models.Korisnici, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = s.db.Exec(query, user.KorisnickoIme, user.Email, string(hashedPassword),
		user.Ime, user.Prezime, user.UlogaID, user.Status)

	return err
}

func (s *UserService) UpdateUser(userID int, user models.Korisnici) error {
	query := `
		UPDATE korisnici 
		SET korisnicko_ime = $1, email = $2, ime = $3, prezime = $4, 
		    uloga_id = $5, status = $6
		WHERE korisnik_id = $7
	`

	_, err := s.db.Exec(query, user.KorisnickoIme, user.Email, user.Ime,
		user.Prezime, user.UlogaID, user.Status, userID)

	return err
}

func (s *UserService) DeleteUser(userID int) error {
	query := `DELETE FROM korisnici WHERE korisnik_id = $1`
	result, err := s.db.Exec(query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

func (s *UserService) GetAllRoles() ([]models.Uloge, error) {
	query := `SELECT uloga_id, naziv_uloge FROM uloge ORDER BY naziv_uloge`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Uloge
	for rows.Next() {
		var role models.Uloge
		err := rows.Scan(&role.UlogaID, &role.NazivUloge)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
