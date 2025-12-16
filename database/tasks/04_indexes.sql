-- ============================================================================
-- Index Creation Script for Oracle Database
-- This script is idempotent - safe to run multiple times
-- ============================================================================

-- Drop existing indexes if they exist (ignore errors if they don't exist)
BEGIN EXECUTE IMMEDIATE 'DROP INDEX idx_dokumenti_kljucne_reci'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP INDEX idx_zadaci_prioritet'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP INDEX idx_clanovi_projekta_korisnik'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- Oracle Text index for full-text search (requires Oracle Text/CONTEXT)
-- If Oracle Text is not available, consider creating a normal index on the column.
BEGIN
	EXECUTE IMMEDIATE 'CREATE INDEX idx_dokumenti_kljucne_reci ON Dokumenti(kljucne_reci) INDEXTYPE IS CTXSYS.CONTEXT';
EXCEPTION WHEN OTHERS THEN
	NULL; -- ignore if text index cannot be created (e.g. missing privilege or option)
END;
/
-- Index on task priority for faster filtering by priority
CREATE INDEX idx_zadaci_prioritet ON Zadaci(prioritet);

-- Index on project members for faster user lookups
CREATE INDEX idx_clanovi_projekta_korisnik ON ClanoviProjekta(korisnik_id);

