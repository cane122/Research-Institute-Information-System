-- Migration: Add resursi column to Projekti table
-- Date: 2025-10-19

-- Add resursi column to Projekti table
ALTER TABLE Projekti ADD COLUMN IF NOT EXISTS resursi TEXT;

-- Add comment to describe the column
COMMENT ON COLUMN Projekti.resursi IS 'Resources needed for the project';
