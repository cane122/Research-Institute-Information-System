package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// Test osnovne konekcije na bazu
func TestDatabaseConnection(t *testing.T) {
	// Test različitih konfiguracija
	testCases := []struct {
		name     string
		host     string
		user     string
		password string
		dbname   string
	}{
		{
			name:     "Default Config",
			host:     "localhost:5432",
			user:     "postgres",
			password: "password",
			dbname:   "research_institute",
		},
		{
			name:     "Alternative Password",
			host:     "localhost:5432",
			user:     "postgres",
			password: "postgres",
			dbname:   "research_institute",
		},
		{
			name:     "Different Port",
			host:     "localhost:5433",
			user:     "postgres",
			password: "password",
			dbname:   "research_institute",
		},
		{
			name:     "Environment Variables",
			host:     getEnvOrDefault("DB_HOST", "localhost:5432"),
			user:     getEnvOrDefault("DB_USER", "postgres"),
			password: getEnvOrDefault("DB_PASSWORD", "password"),
			dbname:   getEnvOrDefault("DB_NAME", "research_institute"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Separated host and port for proper DSN format
			hostParts := strings.Split(tc.host, ":")
			host := hostParts[0]
			port := "5432"
			if len(hostParts) > 1 {
				port = hostParts[1]
			}

			dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				host, port, tc.user, tc.password, tc.dbname)

			t.Logf("Pokušavam konekciju sa: host=%s port=%s user=%s dbname=%s", host, port, tc.user, tc.dbname)

			db, err := sql.Open("postgres", dsn)
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
			err = db.QueryRow("SELECT version()").Scan(&version)
			if err != nil {
				t.Logf("GREŠKA pri izvršavanju upita: %v", err)
				return
			}

			t.Logf("PostgreSQL verzija: %s", version)
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
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)`

		err := db.QueryRow(query, tableName).Scan(&exists)
		if err != nil {
			t.Errorf("Greška pri proveri tabele %s: %v", tableName, err)
			continue
		}

		if !exists {
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

// Test PostgreSQL servisa
func TestPostgreSQLService(t *testing.T) {
	t.Log("=== DIJAGNOSTIKA PostgreSQL SERVISA ===")

	// Proveri da li je servis pokrenut (Windows specifično)
	t.Log("Proveravam PostgreSQL servis...")

	// Pokušaj konekcije na različitim portovima
	ports := []string{"5432", "5433", "5434"}

	for _, port := range ports {
		t.Logf("Testiram port %s...", port)

		dsn := fmt.Sprintf("host=localhost:%s user=postgres password=password dbname=postgres sslmode=disable", port)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			t.Logf("❌ Port %s: Greška pri otvaranju - %v", port, err)
			continue
		}

		err = db.Ping()
		db.Close()

		if err != nil {
			t.Logf("❌ Port %s: Ping neuspešan - %v", port, err)
		} else {
			t.Logf("✅ Port %s: PostgreSQL je aktivan!", port)
		}
	}
}

// Helper funkcija za konekciju
func connectToDatabase(t *testing.T) *sql.DB {
	configs := []string{
		"host=localhost:5432 user=postgres password=password dbname=research_institute sslmode=disable",
		"host=localhost:5432 user=postgres password=postgres dbname=research_institute sslmode=disable",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			getEnvOrDefault("DB_HOST", "localhost:5432"),
			getEnvOrDefault("DB_USER", "postgres"),
			getEnvOrDefault("DB_PASSWORD", "password"),
			getEnvOrDefault("DB_NAME", "research_institute")),
	}

	for i, dsn := range configs {
		t.Logf("Pokušavam konfiguraciju %d: %s", i+1, dsn)

		db, err := sql.Open("postgres", dsn)
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
