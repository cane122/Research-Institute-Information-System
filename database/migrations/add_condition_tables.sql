-- Migration: Add Uslovi and ProcenaUslova tables
-- Description: Adds tables for managing phase transition conditions
-- Date: 2025-10-20

-- Table for defining conditions that must be met for phase transitions
CREATE TABLE IF NOT EXISTS Uslovi (
    uslov_id SERIAL PRIMARY KEY,
    faza_id INT NOT NULL,
    opis TEXT NOT NULL, -- Description of the condition (e.g., "Završiti dokumentaciju")
    kriterijum TEXT NOT NULL, -- Criteria that must be met
    kreiran_datuma TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (faza_id) REFERENCES Faze(faza_id) ON DELETE CASCADE
);

-- Table for tracking condition fulfillment for specific tasks
CREATE TABLE IF NOT EXISTS ProcenaUslova (
    procena_id SERIAL PRIMARY KEY,
    zadatak_id INT NOT NULL,
    uslov_id INT NOT NULL,
    ispunjen BOOLEAN DEFAULT FALSE, -- Whether the condition is fulfilled
    napomena TEXT, -- Additional notes about fulfillment (optional)
    promenio_korisnik_id INT, -- User who evaluated the condition
    datum_procene TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (zadatak_id, uslov_id), -- Each condition can be evaluated once per task
    FOREIGN KEY (zadatak_id) REFERENCES Zadaci(zadatak_id) ON DELETE CASCADE,
    FOREIGN KEY (uslov_id) REFERENCES Uslovi(uslov_id) ON DELETE CASCADE,
    FOREIGN KEY (promenio_korisnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_uslovi_faza ON Uslovi(faza_id);
CREATE INDEX IF NOT EXISTS idx_procena_zadatak ON ProcenaUslova(zadatak_id);
CREATE INDEX IF NOT EXISTS idx_procena_uslov ON ProcenaUslova(uslov_id);
CREATE INDEX IF NOT EXISTS idx_procena_ispunjen ON ProcenaUslova(ispunjen);

-- Insert sample conditions for existing phases
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
-- Conditions for phase "Prikupljanje podataka" (assuming faza_id 7 exists)
(7, 'Završiti dokumentaciju', 'Sva potrebna tehnička dokumentacija mora biti kompletna'),
(7, 'Prikupiti sve izvore', 'Svi relevantni izvori i reference moraju biti sakupljeni'),
(7, 'Završiti analizu', 'Inicijalna analiza podataka mora biti završena')
ON CONFLICT DO NOTHING;

-- Conditions for phase "Razvoj" (assuming faza_id 3 exists)
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(3, 'Implementirati osnovne funkcionalnosti', 'Sve osnovne funkcionalnosti moraju biti kodirane'),
(3, 'Napisati unit testove', 'Pokrivenost testovima mora biti minimum 70%'),
(3, 'Code review', 'Kod mora proci code review proces')
ON CONFLICT DO NOTHING;

-- Conditions for phase "Testiranje" (assuming faza_id 4 exists)
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(4, 'Sve testove prolaze', 'Svi automatski testovi moraju biti zeleni'),
(4, 'QA provera', 'QA tim mora odobriti funkcionalnost'),
(4, 'Performance testiranje', 'Aplikacija mora zadovoljiti performance kriterijume')
ON CONFLICT DO NOTHING;
