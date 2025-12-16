// ============================================================================
// analytics_service.go - Analytics and Logging Service
// ============================================================================

package services

import (
	"database/sql"
	"fmt"

	"github.com/cane/research-institute-system/backend/models"
)

type AnalyticsService struct {
	db *sql.DB
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// ============================================================================
// Activity Logging
// ============================================================================

// ExecuteComplexReport calls the PL/SQL procedure for complex project statistics report
func (s *AnalyticsService) ExecuteComplexReport() (string, error) {
	// Call the stored procedure
	query := `BEGIN kompleksan_izvestaj_projekata; END;`

	_, err := s.db.Exec(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute complex report: %w", err)
	}

	return "Complex report executed successfully. Check database output.", nil
}

// LogActivity logs a user activity
func (s *AnalyticsService) LogActivity(korisnikID *int, req models.ActivityLogRequest) error {
	// Use actual column names from database: korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis
	query := `
		INSERT INTO logaktivnosti (
			korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis
		) VALUES (:1, :2, :3, :4, :5)
	`

	var entitetTip *string
	var entitetID *int
	var opis *string

	if req.EntitetTip != "" {
		entitetTip = &req.EntitetTip
	}
	if req.EntitetID > 0 {
		entitetID = &req.EntitetID
	}
	if req.Opis != "" {
		opis = &req.Opis
	}

	_, err := s.db.Exec(
		query,
		korisnikID,
		req.TipAktivnosti,
		entitetTip,
		entitetID,
		opis,
	)

	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	return nil
}

// GetRecentActivity retrieves recent activities
func (s *AnalyticsService) GetRecentActivity(limit int) ([]models.SkornjeAktivnosti, error) {
	if limit <= 0 {
		limit = 100
	}

	// Query actual database structure
	query := `
		SELECT 
			la.log_id,
			la.korisnik_id,
			COALESCE(k.korisnicko_ime, '') as korisnik_ime,
			la.tip_aktivnosti,
			COALESCE(la.ciljani_entitet, '') as entitet_tip,
			COALESCE(la.ciljani_id, 0) as entitet_id,
			COALESCE(la.opis, '') as naziv_entiteta,
			COALESCE(la.opis, '') as opis,
			'SUCCESS' as rezultat,
			la.datuma as kreiran_datuma
		FROM logaktivnosti la
		LEFT JOIN korisnici k ON la.korisnik_id = k.korisnik_id
		ORDER BY la.datuma DESC
		LIMIT :1
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent activity: %w", err)
	}
	defer rows.Close()

	var activities []models.SkornjeAktivnosti
	for rows.Next() {
		var activity models.SkornjeAktivnosti
		err := rows.Scan(
			&activity.LogID,
			&activity.KorisnikID,
			&activity.KorisnikIme,
			&activity.TipAktivnosti,
			&activity.EntitetTip,
			&activity.EntitetID,
			&activity.NazivEntiteta,
			&activity.Opis,
			&activity.Rezultat,
			&activity.KreiranDatuma,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// ============================================================================
// Statistics
// ============================================================================

// GetDocumentStatistics retrieves document statistics
func (s *AnalyticsService) GetDocumentStatistics() (*models.StatistikaDokumenata, error) {
	// Calculate statistics directly from tables
	query := `
		SELECT 
			(SELECT COUNT(*) FROM dokumenti) as ukupno_dokumenata,
			(SELECT COUNT(*) FROM dokumenti WHERE datuma_postavke >= CURRENT_DATE - INTERVAL '1 month') as novih_dokumenata_mesecno,
			(SELECT COUNT(DISTINCT kreirao_korisnik_id) FROM dokumenti) as broj_autora,
			(SELECT COUNT(DISTINCT projekat_id) FROM dokumenti WHERE projekat_id IS NOT NULL) as broj_projekata_sa_dokumentima,
			(SELECT COALESCE(AVG(verzija_count), 0) FROM (SELECT dokument_id, COUNT(*) as verzija_count FROM verzijedokumenata GROUP BY dokument_id) v) as prosecno_verzija_po_dokumentu
	`

	var stats models.StatistikaDokumenata
	err := s.db.QueryRow(query).Scan(
		&stats.UkupnoDokumenata,
		&stats.NovihDokumenataMesecno,
		&stats.BrojAutora,
		&stats.BrojProjekataSaDokumentima,
		&stats.ProsecnoVerzijaPoDokumentu,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get document statistics: %w", err)
	}

	return &stats, nil
}

// GetActivityStatistics retrieves activity statistics
func (s *AnalyticsService) GetActivityStatistics() ([]models.StatistikaAktivnosti, error) {
	// Calculate statistics directly from LogAktivnosti table
	query := `
		SELECT 
			tip_aktivnosti,
			COUNT(*) as broj_aktivnosti,
			COUNT(DISTINCT korisnik_id) as broj_korisnika,
			COUNT(CASE WHEN datuma >= CURRENT_DATE THEN 1 END) as danas,
			COUNT(CASE WHEN datuma >= CURRENT_DATE - INTERVAL '7 days' THEN 1 END) as ove_nedelje,
			COUNT(CASE WHEN datuma >= CURRENT_DATE - INTERVAL '30 days' THEN 1 END) as ovog_meseca,
			MAX(datuma) as poslednja_aktivnost
		FROM logaktivnosti
		GROUP BY tip_aktivnosti
		ORDER BY broj_aktivnosti DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity statistics: %w", err)
	}
	defer rows.Close()

	var statistics []models.StatistikaAktivnosti
	for rows.Next() {
		var stat models.StatistikaAktivnosti
		err := rows.Scan(
			&stat.TipAktivnosti,
			&stat.BrojAktivnosti,
			&stat.BrojKorisnika,
			&stat.Danas,
			&stat.OveNedelje,
			&stat.OvogMeseca,
			&stat.PoslednjaAktivnost,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity statistic: %w", err)
		}
		statistics = append(statistics, stat)
	}

	return statistics, nil
}

// GetDocumentsByType returns document count grouped by type
func (s *AnalyticsService) GetDocumentsByType() (map[string]int, error) {
	query := `
		SELECT 
			COALESCE(tip_dokumenta, 'Unknown') as tip,
			COUNT(*) as count
		FROM dokumenti
		GROUP BY tip_dokumenta
		ORDER BY count DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents by type: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var tip string
		var count int
		if err := rows.Scan(&tip, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result[tip] = count
	}

	return result, nil
}

// GetDocumentTrends returns document creation trends over time
func (s *AnalyticsService) GetDocumentTrends(days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}

	query := fmt.Sprintf(`
		SELECT 
			DATE(datuma_postavke) as datum,
			COUNT(*) as broj_dokumenata
		FROM dokumenti
		WHERE datuma_postavke >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY DATE(datuma_postavke)
		ORDER BY datum DESC
	`, days)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query document trends: %w", err)
	}
	defer rows.Close()

	var trends []map[string]interface{}
	for rows.Next() {
		var datum sql.NullTime
		var count int
		if err := rows.Scan(&datum, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if datum.Valid {
			trends = append(trends, map[string]interface{}{
				"datum": datum.Time.Format("2006-01-02"),
				"count": count,
			})
		}
	}

	return trends, nil
}

// GetTopContributors returns users with most document uploads
func (s *AnalyticsService) GetTopContributors(limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT 
			k.korisnik_id,
			COALESCE(k.ime || ' ' || k.prezime, k.korisnicko_ime) as ime,
			COUNT(d.dokument_id) as broj_dokumenata
		FROM korisnici k
		JOIN dokumenti d ON k.korisnik_id = d.kreirao_korisnik_id
		GROUP BY k.korisnik_id, k.ime, k.prezime, k.korisnicko_ime
		ORDER BY broj_dokumenata DESC
		LIMIT :1
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top contributors: %w", err)
	}
	defer rows.Close()

	var contributors []map[string]interface{}
	for rows.Next() {
		var korisnikID int
		var ime string
		var count int
		if err := rows.Scan(&korisnikID, &ime, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		contributors = append(contributors, map[string]interface{}{
			"korisnik_id":     korisnikID,
			"ime":             ime,
			"broj_dokumenata": count,
		})
	}

	return contributors, nil
}
