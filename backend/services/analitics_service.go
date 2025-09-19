// ============================================================================
// analytics_service.go - Analytics and Logging Service
// ============================================================================

package services

import (
	"database/sql"

	"github.com/cane/research-institute-system/backend/models"
)

type AnalyticsService struct {
	db *sql.DB
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

func (s *AnalyticsService) GetDashboardStats() (models.DashboardStats, error) {
	var stats models.DashboardStats

	// Count active projects
	err := s.db.QueryRow("SELECT COUNT(*) FROM projekti WHERE status = 'aktivan'").Scan(&stats.AktivniProjekti)
	if err != nil {
		return stats, err
	}

	// Count total documents
	err = s.db.QueryRow("SELECT COUNT(*) FROM dokumenti").Scan(&stats.UkupnoDokumenata)
	if err != nil {
		return stats, err
	}

	// Count tasks in progress
	err = s.db.QueryRow("SELECT COUNT(*) FROM zadaci WHERE progres < 100").Scan(&stats.ZadaciUToku)
	if err != nil {
		return stats, err
	}

	// Count active users (logged in last 30 days)
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM korisnici 
		WHERE poslednja_prijava > CURRENT_TIMESTAMP - INTERVAL '30 days'
	`).Scan(&stats.AktivniKorisnici)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func (s *AnalyticsService) GetActivityLogs(limit int) ([]models.LogAktivnosti, error) {
	query := `
		SELECT l.log_id, l.korisnik_id, l.tip_aktivnosti, l.opis,
		       l.ciljani_entitet, l.ciljani_id, l.datuma,
		       COALESCE(k.korisnicko_ime, 'System') as ime_korisnika
		FROM log_aktivnosti l
		LEFT JOIN korisnici k ON l.korisnik_id = k.korisnik_id
		ORDER BY l.datuma DESC
		LIMIT $1
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.LogAktivnosti
	for rows.Next() {
		var log models.LogAktivnosti
		err := rows.Scan(
			&log.LogID, &log.KorisnikID, &log.TipAktivnosti, &log.Opis,
			&log.CiljaniEntitet, &log.CiljaniID, &log.Datuma, &log.ImeKorisnika,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *AnalyticsService) LogActivity(userID *int, activityType, description, targetEntity string, targetID *int) error {
	query := `
		INSERT INTO log_aktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Exec(query, userID, activityType, description, targetEntity, targetID)
	return err
}
