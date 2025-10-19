// ============================================================================
// folder_service.go - Folder (Folderi) CRUD Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type FolderService struct{ db *sql.DB }

func NewFolderService(db *sql.DB) *FolderService { return &FolderService{db: db} }

func (s *FolderService) GetByID(id int) (*models.Folderi, error) {
	var f models.Folderi
	err := s.db.QueryRow(`SELECT folder_id, naziv_foldera, roditelj_folder_id, vlasnik_id FROM folderi WHERE folder_id = $1`, id).
		Scan(&f.FolderID, &f.NazivFoldera, &f.RoditeljFolderID, &f.VlasnikID)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *FolderService) ListByOwner(ownerID int) ([]models.Folderi, error) {
	rows, err := s.db.Query(`SELECT folder_id, naziv_foldera, roditelj_folder_id, vlasnik_id FROM folderi WHERE vlasnik_id = $1 ORDER BY naziv_foldera`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Folderi
	for rows.Next() {
		var f models.Folderi
		if err := rows.Scan(&f.FolderID, &f.NazivFoldera, &f.RoditeljFolderID, &f.VlasnikID); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, rows.Err()
}

func (s *FolderService) Create(f *models.Folderi) error {
	return s.db.QueryRow(`INSERT INTO folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES ($1, $2, $3) RETURNING folder_id`,
		f.NazivFoldera, f.RoditeljFolderID, f.VlasnikID).Scan(&f.FolderID)
}

func (s *FolderService) Update(f models.Folderi) error {
	_, err := s.db.Exec(`UPDATE folderi SET naziv_foldera = $1, roditelj_folder_id = $2 WHERE folder_id = $3`, f.NazivFoldera, f.RoditeljFolderID, f.FolderID)
	return err
}

func (s *FolderService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM folderi WHERE folder_id = $1`, id)
	return err
}
