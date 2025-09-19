// ============================================================================
// document_service.go - Document Management Service
// ============================================================================

package services

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cane/research-institute-system/backend/config"
	"github.com/cane/research-institute-system/backend/models"
)

type DocumentService struct {
	db         *sql.DB
	uploadPath string
}

func NewDocumentService(db *sql.DB) *DocumentService {
	cfg := config.LoadConfig()
	return &DocumentService{
		db:         db,
		uploadPath: cfg.UploadPath,
	}
}

func (s *DocumentService) GetAllDocuments() ([]models.Dokumenti, error) {
	query := `
		SELECT d.dokument_id, d.projekat_id, d.naziv_dokumenta, d.folder_id,
		       d.opis, d.tip_dokumenta, d.jezik_dokumenta, d.radni_tok_id,
		       d.trenutna_faza_id, d.kreirao_korisnik_id, d.datuma_postavke,
		       d.poslednja_izmena,
		       COALESCE(p.naziv_projekta, '') as naziv_projekta,
		       k.korisnicko_ime as ime_kreirao,
		       COALESCE(f.naziv_faze, '') as naziv_faze,
		       COALESCE(v.version_count, 0) as broj_verzija
		FROM dokumenti d
		LEFT JOIN projekti p ON d.projekat_id = p.projekat_id
		JOIN korisnici k ON d.kreirao_korisnik_id = k.korisnik_id
		LEFT JOIN faze f ON d.trenutna_faza_id = f.faza_id
		LEFT JOIN (
			SELECT dokument_id, COUNT(*) as version_count 
			FROM verzijedokumenata 
			GROUP BY dokument_id
		) v ON d.dokument_id = v.dokument_id
		ORDER BY d.datuma_postavke DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []models.Dokumenti
	for rows.Next() {
		var doc models.Dokumenti
		err := rows.Scan(
			&doc.DokumentID, &doc.ProjekatID, &doc.NazivDokumenta, &doc.FolderID,
			&doc.Opis, &doc.TipDokumenta, &doc.JezikDokumenta, &doc.RadniTokID,
			&doc.TrenutnaFazaID, &doc.KreiraoKorisnikID, &doc.DatumaPostavke,
			&doc.PoslednjaIzmena, &doc.NazivProjekta, &doc.ImeKreirao,
			&doc.NazivFaze, &doc.BrojVerzija,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func (s *DocumentService) GetDocumentsByProject(projectID int) ([]models.Dokumenti, error) {
	query := `
		SELECT d.dokument_id, d.projekat_id, d.naziv_dokumenta, d.folder_id,
		       d.opis, d.tip_dokumenta, d.jezik_dokumenta, d.radni_tok_id,
		       d.trenutna_faza_id, d.kreirao_korisnik_id, d.datuma_postavke,
		       d.poslednja_izmena,
		       p.naziv_projekta, k.korisnicko_ime as ime_kreirao,
		       COALESCE(f.naziv_faze, '') as naziv_faze,
		       COALESCE(v.version_count, 0) as broj_verzija
		FROM dokumenti d
		JOIN projekti p ON d.projekat_id = p.projekat_id
		JOIN korisnici k ON d.kreirao_korisnik_id = k.korisnik_id
		LEFT JOIN faze f ON d.trenutna_faza_id = f.faza_id
		LEFT JOIN (
			SELECT dokument_id, COUNT(*) as version_count 
			FROM verzijedokumenata 
			GROUP BY dokument_id
		) v ON d.dokument_id = v.dokument_id
		WHERE d.projekat_id = $1
		ORDER BY d.datuma_postavke DESC
	`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []models.Dokumenti
	for rows.Next() {
		var doc models.Dokumenti
		err := rows.Scan(
			&doc.DokumentID, &doc.ProjekatID, &doc.NazivDokumenta, &doc.FolderID,
			&doc.Opis, &doc.TipDokumenta, &doc.JezikDokumenta, &doc.RadniTokID,
			&doc.TrenutnaFazaID, &doc.KreiraoKorisnikID, &doc.DatumaPostavke,
			&doc.PoslednjaIzmena, &doc.NazivProjekta, &doc.ImeKreirao,
			&doc.NazivFaze, &doc.BrojVerzija,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func (s *DocumentService) GetDocumentByID(documentID int) (models.Dokumenti, error) {
	var doc models.Dokumenti
	query := `
		SELECT d.dokument_id, d.projekat_id, d.naziv_dokumenta, d.folder_id,
		       d.opis, d.tip_dokumenta, d.jezik_dokumenta, d.radni_tok_id,
		       d.trenutna_faza_id, d.kreirao_korisnik_id, d.datuma_postavke,
		       d.poslednja_izmena,
		       COALESCE(p.naziv_projekta, '') as naziv_projekta,
		       k.korisnicko_ime as ime_kreirao,
		       COALESCE(f.naziv_faze, '') as naziv_faze
		FROM dokumenti d
		LEFT JOIN projekti p ON d.projekat_id = p.projekat_id
		JOIN korisnici k ON d.kreirao_korisnik_id = k.korisnik_id
		LEFT JOIN faze f ON d.trenutna_faza_id = f.faza_id
		WHERE d.dokument_id = $1
	`

	err := s.db.QueryRow(query, documentID).Scan(
		&doc.DokumentID, &doc.ProjekatID, &doc.NazivDokumenta, &doc.FolderID,
		&doc.Opis, &doc.TipDokumenta, &doc.JezikDokumenta, &doc.RadniTokID,
		&doc.TrenutnaFazaID, &doc.KreiraoKorisnikID, &doc.DatumaPostavke,
		&doc.PoslednjaIzmena, &doc.NazivProjekta, &doc.ImeKreirao, &doc.NazivFaze,
	)

	return doc, err
}

func (s *DocumentService) UploadDocument(req models.UploadDocumentRequest, fileData []byte, fileName string, userID int) error {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert document record
	var documentID int
	docQuery := `
		INSERT INTO dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, 
		                      tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING dokument_id
	`

	err = tx.QueryRow(docQuery, req.ProjekatID, req.NazivDokumenta, req.FolderID,
		req.Opis, req.TipDokumenta, req.JezikDokumenta, userID).Scan(&documentID)
	if err != nil {
		return err
	}

	// Generate unique file path
	timestamp := time.Now().Format("20060102_150405")
	fileExt := filepath.Ext(fileName)
	uniqueFileName := fmt.Sprintf("%d_%s_%s%s", documentID, timestamp,
		strings.ReplaceAll(req.NazivDokumenta, " ", "_"), fileExt)
	filePath := filepath.Join(s.uploadPath, uniqueFileName)

	// Save file to disk
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(string(fileData)))
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Calculate file size in MB
	fileSizeMB := float64(len(fileData)) / (1024 * 1024)

	// Insert document version
	versionQuery := `
		INSERT INTO verzijedokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, 
		                               velicina_fajla_mb, postavio_korisnik_id)
		VALUES ($1, '1.0', $2, $3, $4)
	`

	_, err = tx.Exec(versionQuery, documentID, filePath, fileSizeMB, userID)
	if err != nil {
		return err
	}

	// Add tags if provided
	for _, tagName := range req.Tagovi {
		if err := s.addDocumentTagInTx(tx, documentID, tagName); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DocumentService) addDocumentTagInTx(tx *sql.Tx, documentID int, tagName string) error {
	// Check if tag exists, if not create it
	var tagID int
	err := tx.QueryRow("SELECT tag_id FROM tagovi WHERE naziv_taga = $1", tagName).Scan(&tagID)
	if err == sql.ErrNoRows {
		// Create new tag
		err = tx.QueryRow("INSERT INTO tagovi (naziv_taga) VALUES ($1) RETURNING tag_id", tagName).Scan(&tagID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Link tag to document
	_, err = tx.Exec("INSERT INTO dokumenttagovi (dokument_id, tag_id) VALUES ($1, $2)", documentID, tagID)
	return err
}

func (s *DocumentService) UpdateDocument(documentID int, req models.UploadDocumentRequest) error {
	query := `
		UPDATE dokumenti 
		SET naziv_dokumenta = $1, projekat_id = $2, folder_id = $3, opis = $4,
		    tip_dokumenta = $5, jezik_dokumenta = $6, poslednja_izmena = CURRENT_TIMESTAMP
		WHERE dokument_id = $7
	`

	_, err := s.db.Exec(query, req.NazivDokumenta, req.ProjekatID, req.FolderID,
		req.Opis, req.TipDokumenta, req.JezikDokumenta, documentID)

	return err
}

func (s *DocumentService) DeleteDocument(documentID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get all file paths for deletion
	var filePaths []string
	versionQuery := `SELECT putanja_do_fajla FROM verzijedokumenata WHERE dokument_id = $1`
	rows, err := tx.Query(versionQuery, documentID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err != nil {
			return err
		}
		filePaths = append(filePaths, filePath)
	}

	// Delete document (cascade will handle related records)
	deleteQuery := `DELETE FROM dokumenti WHERE dokument_id = $1`
	result, err := tx.Exec(deleteQuery, documentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document with ID %d not found", documentID)
	}

	// Commit transaction before deleting files
	if err := tx.Commit(); err != nil {
		return err
	}

	// Delete physical files
	for _, filePath := range filePaths {
		if err := os.Remove(filePath); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to delete file %s: %v\n", filePath, err)
		}
	}

	return nil
}

func (s *DocumentService) GetDocumentVersions(documentID int) ([]models.VerzijeDokumenata, error) {
	query := `
		SELECT v.verzija_id, v.dokument_id, v.verzija_oznaka, v.putanja_do_fajla,
		       v.velicina_fajla_mb, v.postavio_korisnik_id, v.datuma_postavke
		FROM verzijedokumenata v
		WHERE v.dokument_id = $1
		ORDER BY v.datuma_postavke DESC
	`

	rows, err := s.db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.VerzijeDokumenata
	for rows.Next() {
		var version models.VerzijeDokumenata
		err := rows.Scan(
			&version.VerzijaID, &version.DokumentID, &version.VerzijaOznaka,
			&version.PutanjaDoFajla, &version.VelicinafajlaMB, &version.PostavioKorisnikID,
			&version.DatumaPostavke,
		)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

func (s *DocumentService) GetDocumentTags(documentID int) ([]models.Tagovi, error) {
	query := `
		SELECT t.tag_id, t.naziv_taga
		FROM tagovi t
		JOIN dokumenttagovi dt ON t.tag_id = dt.tag_id
		WHERE dt.dokument_id = $1
		ORDER BY t.naziv_taga
	`

	rows, err := s.db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tagovi
	for rows.Next() {
		var tag models.Tagovi
		err := rows.Scan(&tag.TagID, &tag.NazivTaga)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *DocumentService) AddDocumentTag(documentID int, tagName string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.addDocumentTagInTx(tx, documentID, tagName)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DocumentService) RemoveDocumentTag(documentID, tagID int) error {
	query := `DELETE FROM dokumenttagovi WHERE dokument_id = $1 AND tag_id = $2`
	_, err := s.db.Exec(query, documentID, tagID)
	return err
}

func (s *DocumentService) GetDocumentMetadata(documentID int) ([]models.MetaPodaci, error) {
	query := `
		SELECT meta_id, dokument_id, kljuc, vrednost
		FROM metapodaci
		WHERE dokument_id = $1
		ORDER BY kljuc
	`

	rows, err := s.db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []models.MetaPodaci
	for rows.Next() {
		var meta models.MetaPodaci
		err := rows.Scan(&meta.MetaID, &meta.DokumentID, &meta.Kljuc, &meta.Vrednost)
		if err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, nil
}

func (s *DocumentService) UpdateDocumentMetadata(documentID int, metadata []models.MetaPodaci) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing metadata
	_, err = tx.Exec("DELETE FROM metapodaci WHERE dokument_id = $1", documentID)
	if err != nil {
		return err
	}

	// Insert new metadata
	for _, meta := range metadata {
		_, err = tx.Exec("INSERT INTO metapodaci (dokument_id, kljuc, vrednost) VALUES ($1, $2, $3)",
			documentID, meta.Kljuc, meta.Vrednost)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DocumentService) GetAllFolders(userID int) ([]models.Folderi, error) {
	query := `
		SELECT folder_id, naziv_foldera, roditelj_folder_id, vlasnik_id
		FROM folderi
		WHERE vlasnik_id = $1
		ORDER BY naziv_foldera
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []models.Folderi
	for rows.Next() {
		var folder models.Folderi
		err := rows.Scan(&folder.FolderID, &folder.NazivFoldera,
			&folder.RoditeljFolderID, &folder.VlasnikID)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (s *DocumentService) CreateFolder(folder models.Folderi) error {
	query := `
		INSERT INTO folderi (naziv_foldera, roditelj_folder_id, vlasnik_id)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(query, folder.NazivFoldera, folder.RoditeljFolderID, folder.VlasnikID)
	return err
}
