-- Research Institute Information System Database Schema
-- Created for PostgreSQL
-- Encoding: UTF-8

-- Set client encoding to UTF8
SET client_encoding = 'UTF8';

-- Module 1: User and Role Management

-- Table for defining user roles (Administrator, Manager, Researcher)
CREATE TABLE Uloge (
    uloga_id SERIAL PRIMARY KEY,
    naziv_uloge VARCHAR(50) UNIQUE NOT NULL -- e.g., 'Administrator', 'Rukovodilac projekta', 'Istrazivac'
);

-- Main table for system users
CREATE TABLE Korisnici (
    korisnik_id SERIAL PRIMARY KEY,
    korisnicko_ime VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    hash_sifre VARCHAR(255) NOT NULL,
    ime VARCHAR(100),
    prezime VARCHAR(100),
    uloga_id INT NOT NULL,
    status VARCHAR(20) DEFAULT 'aktivan', -- e.g., 'aktivan', 'neaktivan'
    poslednja_prijava TIMESTAMP,
    kreiran_datuma TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (uloga_id) REFERENCES Uloge(uloga_id)
);

-- Module 2: Project, Task and Documentation Management

-- Table for defining workflows (e.g., for projects or documents)
CREATE TABLE RadniTokovi (
    radni_tok_id SERIAL PRIMARY KEY,
    naziv VARCHAR(255) NOT NULL,
    tip_toka VARCHAR(50) NOT NULL CHECK (tip_toka IN ('PROJEKAT', 'DOKUMENTACIJA')),
    opis TEXT,
    da_li_je_sablon BOOLEAN DEFAULT FALSE, -- Whether this workflow can be used as a template
    UNIQUE(naziv, tip_toka)
);

-- Phases within a workflow (e.g., Planning, Development, Completed)
CREATE TABLE Faze (
    faza_id SERIAL PRIMARY KEY,
    radni_tok_id INT NOT NULL,
    naziv_faze VARCHAR(100) NOT NULL,
    redosled INT NOT NULL, -- For sorting phases within workflow
    FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id) ON DELETE CASCADE
);

-- Table containing basic project information
CREATE TABLE Projekti (
    projekat_id SERIAL PRIMARY KEY,
    naziv_projekta VARCHAR(255) NOT NULL,
    opis TEXT,
    datum_pocetka DATE,
    datum_zavrsetka DATE,
    status VARCHAR(50) DEFAULT 'Aktivan', -- e.g., 'Aktivan', 'Završen', 'Otkazan'
    rukovodilac_id INT, -- User who created and manages the project
    radni_tok_id INT, -- Workflow applied to tasks in this project
    FOREIGN KEY (rukovodilac_id) REFERENCES Korisnici(korisnik_id),
    FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id)
);

-- Table linking users to projects (team members)
CREATE TABLE ClanoviProjekta (
    projekat_id INT NOT NULL,
    korisnik_id INT NOT NULL,
    PRIMARY KEY (projekat_id, korisnik_id),
    FOREIGN KEY (projekat_id) REFERENCES Projekti(projekat_id) ON DELETE CASCADE,
    FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id) ON DELETE CASCADE
);

-- Table for tasks within projects
CREATE TABLE Zadaci (
    zadatak_id SERIAL PRIMARY KEY,
    projekat_id INT NOT NULL,
    faza_id INT NOT NULL, -- Current phase of the task
    naziv_zadatka VARCHAR(255) NOT NULL,
    opis TEXT,
    dodeljen_korisniku_id INT,
    rok DATE,
    prioritet VARCHAR(50), -- e.g., 'Nizak', 'Srednji', 'Visok'
    progres INT DEFAULT 0,
    kreiran_datuma TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (projekat_id) REFERENCES Projekti(projekat_id) ON DELETE CASCADE,
    FOREIGN KEY (faza_id) REFERENCES Faze(faza_id),
    FOREIGN KEY (dodeljen_korisniku_id) REFERENCES Korisnici(korisnik_id)
);

-- Table for task comments
CREATE TABLE KomentariZadataka (
    komentar_id SERIAL PRIMARY KEY,
    zadatak_id INT NOT NULL,
    korisnik_id INT NOT NULL,
    tekst_komentara TEXT NOT NULL,
    datuma_kreiranja TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (zadatak_id) REFERENCES Zadaci(zadatak_id) ON DELETE CASCADE,
    FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Table for recording phase change requests from Researchers
CREATE TABLE ZahteviPromeneFaze (
    zahtev_id SERIAL PRIMARY KEY,
    zadatak_id INT NOT NULL,
    podnosilac_zahteva_id INT NOT NULL,
    zahtevana_faza_id INT NOT NULL,
    status VARCHAR(50) DEFAULT 'Na cekanju', -- e.g., 'Na cekanju', 'Odobren', 'Odbijen'
    komentar TEXT,
    datum_kreiranja TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (zadatak_id) REFERENCES Zadaci(zadatak_id) ON DELETE CASCADE,
    FOREIGN KEY (podnosilac_zahteva_id) REFERENCES Korisnici(korisnik_id),
    FOREIGN KEY (zahtevana_faza_id) REFERENCES Faze(faza_id)
);

-- Module 3: Document and Metadata Management

CREATE TABLE Folderi (
    folder_id SERIAL PRIMARY KEY,
    naziv_foldera VARCHAR(255) NOT NULL,
    roditelj_folder_id INT,
    vlasnik_id INT NOT NULL,
    FOREIGN KEY (roditelj_folder_id) REFERENCES Folderi(folder_id) ON DELETE CASCADE,
    FOREIGN KEY (vlasnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Table for basic document information
CREATE TABLE Dokumenti (
    dokument_id SERIAL PRIMARY KEY,
    projekat_id INT, -- Optional, document can be part of a project
    naziv_dokumenta VARCHAR(255) NOT NULL,
    folder_id INT,
    opis TEXT,
    tip_dokumenta VARCHAR(50), -- e.g., 'Istraživački rad', 'PDF', 'CSV'
    jezik_dokumenta VARCHAR(50),
    radni_tok_id INT,
    trenutna_faza_id INT, -- If document follows workflow
    kreirao_korisnik_id INT NOT NULL,
    datuma_postavke TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    poslednja_izmena TIMESTAMP,
    FOREIGN KEY (projekat_id) REFERENCES Projekti(projekat_id),
    FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id),
    FOREIGN KEY (trenutna_faza_id) REFERENCES Faze(faza_id),
    FOREIGN KEY (kreirao_korisnik_id) REFERENCES Korisnici(korisnik_id),
    FOREIGN KEY (folder_id) REFERENCES Folderi(folder_id) ON DELETE SET NULL
);

-- Table for tracking document versions
CREATE TABLE VerzijeDokumenata (
    verzija_id SERIAL PRIMARY KEY,
    dokument_id INT NOT NULL,
    verzija_oznaka VARCHAR(50), -- e.g., 'v1.0', 'v1.1 final'
    putanja_do_fajla VARCHAR(1024) NOT NULL,
    velicina_fajla_MB DECIMAL(10, 2),
    postavio_korisnik_id INT NOT NULL,
    datuma_postavke TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    FOREIGN KEY (postavio_korisnik_id) REFERENCES Korisnici(korisnik_id)
);

CREATE TABLE LLMSazeci (
    sazetak_id SERIAL PRIMARY KEY,
    dokument_id INT NOT NULL,
    verzija_oznaka VARCHAR(50),
    sazetak TEXT NOT NULL,
    datum_kreiranja TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE
);

-- Table for additional metadata, allows flexibility
CREATE TABLE MetaPodaci (
    meta_id SERIAL PRIMARY KEY,
    dokument_id INT NOT NULL,
    kljuc VARCHAR(100) NOT NULL, -- e.g., 'ISO Broj', 'Izvorni URL', 'LLM sažetak'
    vrednost TEXT,
    FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE
);

-- Table for tags for easier search
CREATE TABLE Tagovi (
    tag_id SERIAL PRIMARY KEY,
    naziv_taga VARCHAR(100) UNIQUE NOT NULL
);

-- Links documents to tags (many-to-many relationship)
CREATE TABLE DokumentTagovi (
    dokument_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (dokument_id, tag_id),
    FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES Tagovi(tag_id) ON DELETE CASCADE
);

-- Table defining access rights (read, write, delete) for users over documents
CREATE TABLE DozvoleDokumenata (
    dozvola_id SERIAL PRIMARY KEY,
    dokument_id INT NOT NULL,
    korisnik_id INT NOT NULL,
    moze_citati BOOLEAN DEFAULT TRUE,
    moze_menjati BOOLEAN DEFAULT FALSE,
    moze_brisati BOOLEAN DEFAULT FALSE,
    UNIQUE (dokument_id, korisnik_id), -- Each user has one permission row per document
    FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id) ON DELETE CASCADE
);

-- Table for tracking phase history - applies ONLY to project documentation
CREATE TABLE IstorijaFazaDokumenta (
    istorija_id SERIAL PRIMARY KEY,
    dokument_id INT NOT NULL,
    prethodna_faza_id INT,
    nova_faza_id INT NOT NULL,
    korisnik_id INT NOT NULL,
    datum_promene TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    FOREIGN KEY (prethodna_faza_id) REFERENCES Faze(faza_id),
    FOREIGN KEY (nova_faza_id) REFERENCES Faze(faza_id),
    FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Module 4: Analytics and Logging

-- Table for logging all important system activities
CREATE TABLE LogAktivnosti (
    log_id BIGSERIAL PRIMARY KEY,
    korisnik_id INT,
    tip_aktivnosti VARCHAR(100) NOT NULL, -- e.g., 'KREIRAN_PROJEKAT', 'UPLOAD_DOKUMENTA', 'PROMENA_FAZE'
    opis TEXT,
    ciljani_entitet VARCHAR(50), -- e.g., 'Projekat', 'Zadatak', 'Dokument'
    ciljani_id INT,
    datuma TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Insert default roles
INSERT INTO Uloge (naziv_uloge) VALUES 
('Administrator'),
('Rukovodilac projekta'),
('Istrazivac'),
('Organizator projekta');

-- Create indexes for better performance
CREATE INDEX idx_korisnici_email ON Korisnici(email);
CREATE INDEX idx_korisnici_korisnicko_ime ON Korisnici(korisnicko_ime);
CREATE INDEX idx_korisnici_uloga ON Korisnici(uloga_id);
CREATE INDEX idx_korisnici_status ON Korisnici(status);

CREATE INDEX idx_projekti_rukovodilac ON Projekti(rukovodilac_id);
CREATE INDEX idx_projekti_status ON Projekti(status);
CREATE INDEX idx_projekti_datum_pocetka ON Projekti(datum_pocetka);

CREATE INDEX idx_zadaci_projekat ON Zadaci(projekat_id);
CREATE INDEX idx_zadaci_dodeljen_korisnik ON Zadaci(dodeljen_korisniku_id);
CREATE INDEX idx_zadaci_faza ON Zadaci(faza_id);
CREATE INDEX idx_zadaci_rok ON Zadaci(rok);

CREATE INDEX idx_dokumenti_projekat ON Dokumenti(projekat_id);
CREATE INDEX idx_dokumenti_kreirao ON Dokumenti(kreirao_korisnik_id);
CREATE INDEX idx_dokumenti_folder ON Dokumenti(folder_id);
CREATE INDEX idx_dokumenti_tip ON Dokumenti(tip_dokumenta);

CREATE INDEX idx_verzije_dokument ON VerzijeDokumenata(dokument_id);
CREATE INDEX idx_verzije_postavio ON VerzijeDokumenata(postavio_korisnik_id);

CREATE INDEX idx_log_korisnik ON LogAktivnosti(korisnik_id);
CREATE INDEX idx_log_datum ON LogAktivnosti(datuma);
CREATE INDEX idx_log_tip ON LogAktivnosti(tip_aktivnosti);

-- Create some useful views
CREATE OR REPLACE VIEW v_aktivni_projekti AS
SELECT 
    p.projekat_id,
    p.naziv_projekta,
    p.opis,
    p.datum_pocetka,
    p.datum_zavrsetka,
    CONCAT(k.ime, ' ', k.prezime) as rukovodilac_ime,
    k.email as rukovodilac_email,
    rt.naziv as radni_tok_naziv,
    COUNT(z.zadatak_id) as ukupno_zadataka,
    COUNT(CASE WHEN z.progres = 100 THEN 1 END) as zavrsenih_zadataka
FROM Projekti p
LEFT JOIN Korisnici k ON p.rukovodilac_id = k.korisnik_id
LEFT JOIN RadniTokovi rt ON p.radni_tok_id = rt.radni_tok_id
LEFT JOIN Zadaci z ON p.projekat_id = z.projekat_id
WHERE p.status = 'Aktivan'
GROUP BY p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka, p.datum_zavrsetka, k.ime, k.prezime, k.email, rt.naziv;

CREATE OR REPLACE VIEW v_dokumenti_sa_verzijama AS
SELECT 
    d.dokument_id,
    d.naziv_dokumenta,
    d.opis,
    d.tip_dokumenta,
    d.jezik_dokumenta,
    CONCAT(k.ime, ' ', k.prezime) as kreirao_ime,
    d.datuma_postavke,
    d.poslednja_izmena,
    COUNT(vd.verzija_id) as broj_verzija,
    MAX(vd.verzija_oznaka) as poslednja_verzija
FROM Dokumenti d
LEFT JOIN Korisnici k ON d.kreirao_korisnik_id = k.korisnik_id
LEFT JOIN VerzijeDokumenata vd ON d.dokument_id = vd.dokument_id
GROUP BY d.dokument_id, d.naziv_dokumenta, d.opis, d.tip_dokumenta, d.jezik_dokumenta, k.ime, k.prezime, d.datuma_postavke, d.poslednja_izmena;

CREATE OR REPLACE VIEW v_zadaci_sa_detaljima AS
SELECT 
    z.zadatak_id,
    z.naziv_zadatka,
    z.opis,
    z.rok,
    z.prioritet,
    z.progres,
    p.naziv_projekta,
    CONCAT(k.ime, ' ', k.prezime) as dodeljen_korisniku,
    f.naziv_faze,
    f.redosled as faza_redosled
FROM Zadaci z
JOIN Projekti p ON z.projekat_id = p.projekat_id
LEFT JOIN Korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
JOIN Faze f ON z.faza_id = f.faza_id;

-- Insert default workflows
INSERT INTO RadniTokovi (naziv, tip_toka, opis, da_li_je_sablon) VALUES 
('Standardni projektni tok', 'PROJEKAT', 'Osnovni radni tok za projekte', TRUE),
('Istrazivacki tok', 'PROJEKAT', 'Tok za istrazivacke projekte', TRUE),
('Dokumentacioni tok', 'DOKUMENTACIJA', 'Tok za upravljanje dokumentima', TRUE);

-- Insert default phases for project workflow
INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES 
(1, 'Planiranje', 1),
(1, 'Analiza', 2),
(1, 'Razvoj', 3),
(1, 'Testiranje', 4),
(1, 'Završeno', 5);

-- Insert default phases for research workflow
INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES 
(2, 'Definisanje istrazivanja', 1),
(2, 'Prikupljanje podataka', 2),
(2, 'Analiza podataka', 3),
(2, 'Pisanje izvestaja', 4),
(2, 'Publikovanje', 5);

-- Insert default phases for documentation workflow
INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES 
(3, 'Kreiranje', 1),
(3, 'Revizija', 2),
(3, 'Odobravanje', 3),
(3, 'Finalizovanje', 4),
(3, 'Arhiviranje', 5);
