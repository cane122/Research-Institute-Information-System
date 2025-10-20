package services

import (
	"database/sql"
)

type FazeService struct {
	db *sql.DB
}

func NewFazeService(db *sql.DB) *FazeService {
	return &FazeService{db: db}
}

// GetFazaNaziv returns the name of a phase by its ID
func (s *FazeService) GetFazaNaziv(fazaID int) (string, error) {
	var naziv string
	err := s.db.QueryRow("SELECT naziv_faze FROM faze WHERE faza_id = $1", fazaID).Scan(&naziv)
	if err != nil {
		return "", err
	}
	return naziv, nil
}
