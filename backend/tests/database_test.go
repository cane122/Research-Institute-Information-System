package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

// Test osnovne konekcije na bazu
func TestDatabaseConnection(t *testing.T) {
	// Test različitih konfiguracija
	testCases := []struct {
		name     string
		host     string
		port     string
		user     string
		password string
		sid      string
	}{
		{
			name:     "Default XE",
			host:     "localhost",
			port:     "1521",
			user:     "SYSTEM",
			password: "",
			sid:      "xe",
		},
		{
			name:     "XE with password",
			host:     "localhost",
			port:     "1521",
			user:     "SYSTEM",
			password: "oracle",
			sid:      "xe",
		},
		{
			name:     "Environment Variables",
			host:     getEnvOrDefault("DB_HOST", "localhost"),
			port:     getEnvOrDefault("DB_PORT", "1521"),
			user:     getEnvOrDefault("DB_USER", "SYSTEM"),
			password: getEnvOrDefault("DB_PASSWORD", ""),
			sid:      getEnvOrDefault("DB_SID", "xe"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// go-ora connection URL format: oracle://user:password@host:port/service_name
			connStr := fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
				tc.user, tc.password, tc.host, tc.port, tc.sid)

			t.Logf("Pokušavam konekciju sa: host=%s port=%s user=%s sid=%s",
				tc.host, tc.port, tc.user, tc.sid)

			db, err := sql.Open("oracle", connStr)
			if err != nil {
				t.Logf("GREŠKA pri otvaranju konekcije: %v", err)
				return
			}
			defer db.Close()

			// Test ping sa timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = db.PingContext(ctx)
			if err != nil {
				t.Logf("GREŠKA pri ping-u: %v", err)
				return
			}

			t.Logf("✅ USPEŠNA konekcija sa konfiguracijom: %s", tc.name)

			// Test jednostavnog upita
			var version string
			err = db.QueryRow("SELECT BANNER FROM v$version WHERE ROWNUM = 1").Scan(&version)
			if err != nil {
				t.Logf("GREŠKA pri izvršavanju upita: %v", err)
				return
			}

			t.Logf("Oracle Database verzija: %s", version)
		})
	}
}

// Test postojanja tabela
func TestDatabaseTables(t *testing.T) {
	db := connectToDatabase(t)
	if db == nil {
		t.Skip("Preskačem test - nema konekcije na bazu")
		return
	}
	defer db.Close()

	expectedTables := []string{
		"uloge", "korisnici", "radnitokovi", "faze", "projekti",
		"clanoviprojekta", "zadaci", "komentarizadataka", "zahtevipromenefaze",
		"folderi", "dokumenti", "verzijedokumenata", "llmsazeci", "metapodaci",
		"tagovi", "dokumenttagovi", "dozvoledokumenata", "istorijafazadokumenta",
		"logaktivnosti",
	}

	for _, tableName := range expectedTables {
		var tableCount int
		query := `SELECT COUNT(*) FROM user_tables WHERE UPPER(table_name) = UPPER(:1)`

		err := db.QueryRow(query, tableName).Scan(&tableCount)
		if err != nil {
			t.Errorf("Greška pri proveri tabele %s: %v", tableName, err)
			continue
		}

		if tableCount == 0 {
			t.Errorf("❌ Tabela '%s' ne postoji", tableName)
		} else {
			t.Logf("✅ Tabela '%s' postoji", tableName)
		}
	}
}

// Test osnovnih podataka
func TestDefaultData(t *testing.T) {
	db := connectToDatabase(t)
	if db == nil {
		t.Skip("Preskačem test - nema konekcije na bazu")
		return
	}
	defer db.Close()

	// Proveri uloge
	var roleCount int
	err := db.QueryRow("SELECT COUNT(*) FROM uloge").Scan(&roleCount)
	if err != nil {
		t.Errorf("Greška pri brojanju uloga: %v", err)
	} else if roleCount < 4 {
		t.Errorf("❌ Premalo uloga u bazi: %d (očekivano 4)", roleCount)
	} else {
		t.Logf("✅ Uloge u bazi: %d", roleCount)
	}

	// Proveri radne tokove
	var workflowCount int
	err = db.QueryRow("SELECT COUNT(*) FROM radnitokovi").Scan(&workflowCount)
	if err != nil {
		t.Errorf("Greška pri brojanju radnih tokova: %v", err)
	} else if workflowCount < 3 {
		t.Errorf("❌ Premalo radnih tokova u bazi: %d (očekivano 3)", workflowCount)
	} else {
		t.Logf("✅ Radni tokovi u bazi: %d", workflowCount)
	}

	// Proveri faze
	var phaseCount int
	err = db.QueryRow("SELECT COUNT(*) FROM faze").Scan(&phaseCount)
	if err != nil {
		t.Errorf("Greška pri brojanju faza: %v", err)
	} else if phaseCount < 15 {
		t.Errorf("❌ Premalo faza u bazi: %d (očekivano 15)", phaseCount)
	} else {
		t.Logf("✅ Faze u bazi: %d", phaseCount)
	}
}

// Test Oracle Database servisa
func TestOracleDatabaseService(t *testing.T) {
	t.Log("=== DIJAGNOSTIKA Oracle SERVISA ===")

	// Pokušaj konekcije na različitim portovima
	ports := []string{"1521", "1522"}

	for _, port := range ports {
		t.Logf("Testiram port %s...", port)

		connStr := fmt.Sprintf("oracle://SYSTEM:@localhost:%s/xe", port)
		db, err := sql.Open("oracle", connStr)
		if err != nil {
			t.Logf("❌ Port %s: Greška pri otvaranju - %v", port, err)
			continue
		}

		err = db.Ping()
		db.Close()

		if err != nil {
			t.Logf("❌ Port %s: Ping neuspešan - %v", port, err)
		} else {
			t.Logf("✅ Port %s: Oracle je aktivan!", port)
		}
	}
}

// Helper funkcija za konekciju
func connectToDatabase(t *testing.T) *sql.DB {
	configs := []string{
		"oracle://SYSTEM:@localhost:1521/xe",
		"oracle://SYSTEM:oracle@localhost:1521/xe",
		fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
			getEnvOrDefault("DB_USER", "SYSTEM"),
			getEnvOrDefault("DB_PASSWORD", ""),
			getEnvOrDefault("DB_HOST", "localhost"),
			getEnvOrDefault("DB_PORT", "1521"),
			getEnvOrDefault("DB_SID", "xe")),
	}

	for i, connStr := range configs {
		t.Logf("Pokušavam konfiguraciju %d", i+1)

		db, err := sql.Open("oracle", connStr)
		if err != nil {
			t.Logf("Greška pri otvaranju %d: %v", i+1, err)
			continue
		}

		err = db.Ping()
		if err != nil {
			t.Logf("Ping greška %d: %v", i+1, err)
			db.Close()
			continue
		}

		t.Logf("✅ Uspešna konekcija sa konfiguracijom %d", i+1)
		return db
	}

	t.Log("❌ Sve konfiguracije neuspešne")
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
