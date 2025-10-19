// ============================================================================
// role_service.go - Role (Uloge) CRUD Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type RoleService struct {
	db *sql.DB
}

func NewRoleService(db *sql.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) GetAll() ([]models.Uloge, error) {
	rows, err := s.db.Query(`SELECT uloga_id, naziv_uloge FROM uloge ORDER BY naziv_uloge`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Uloge
	for rows.Next() {
		var r models.Uloge
		if err := rows.Scan(&r.UlogaID, &r.NazivUloge); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, rows.Err()
}

func (s *RoleService) GetByID(id int) (*models.Uloge, error) {
	var r models.Uloge
	err := s.db.QueryRow(`SELECT uloga_id, naziv_uloge FROM uloge WHERE uloga_id = $1`, id).Scan(&r.UlogaID, &r.NazivUloge)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *RoleService) Create(role *models.Uloge) error {
	return s.db.QueryRow(`INSERT INTO uloge (naziv_uloge) VALUES ($1) RETURNING uloga_id`, role.NazivUloge).Scan(&role.UlogaID)
}

func (s *RoleService) Update(role models.Uloge) error {
	_, err := s.db.Exec(`UPDATE uloge SET naziv_uloge = $1 WHERE uloga_id = $2`, role.NazivUloge, role.UlogaID)
	return err
}

func (s *RoleService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM uloge WHERE uloga_id = $1`, id)
	return err
}
