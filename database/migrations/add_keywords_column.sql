-- Migration: Add kljucne_reci column to Dokumenti table
-- Date: 2025-10-09

-- Add keywords column to documents table
ALTER TABLE Dokumenti ADD COLUMN IF NOT EXISTS kljucne_reci TEXT;

-- Add comment for documentation
COMMENT ON COLUMN Dokumenti.kljucne_reci IS 'Free-form keywords separated by commas';
