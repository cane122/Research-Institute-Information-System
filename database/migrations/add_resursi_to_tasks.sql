-- Add resursi column to Zadaci table
-- Migration created: 2025-10-20

ALTER TABLE Zadaci ADD COLUMN IF NOT EXISTS resursi TEXT;

COMMENT ON COLUMN Zadaci.resursi IS 'Resources needed for the task';
