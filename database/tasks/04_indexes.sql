-- ============================================================================
-- Index Creation Script for Oracle Database - Document Management System
-- This script is idempotent - safe to run multiple times
-- ============================================================================

-- Drop existing indexes if they exist (ignore errors if they don't exist)
BEGIN EXECUTE IMMEDIATE 'DROP INDEX idx_logaktivnosti_korisnik'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- ============================================================================
-- INDEKS ZA LOGOVE (Activity Logs by User)
-- ============================================================================

-- Index on activity logs by user for faster user activity queries
CREATE INDEX idx_logaktivnosti_korisnik ON LogAktivnosti(korisnik_id);


