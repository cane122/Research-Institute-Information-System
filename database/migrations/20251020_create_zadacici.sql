-- Migration: Create zadacici table for document checklists
CREATE TABLE IF NOT EXISTS zadacici (
    zadacic_id SERIAL PRIMARY KEY,
    dokument_id INT NOT NULL REFERENCES dokumenti(dokument_id) ON DELETE CASCADE,
    opis TEXT NOT NULL,
    izvrsen BOOLEAN NOT NULL DEFAULT FALSE
);
