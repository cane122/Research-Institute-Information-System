-- Research Institute Information System Database Schema
-- Created for PostgreSQL

-- Module 1: User and Role Management

-- Table for defining user roles (Administrator, Manager, Researcher)
CREATE TABLE Uloge (
    uloga_id SERIAL PRIMARY KEY,
    naziv_uloge VARCHAR(50) UNIQUE NOT NULL -- e.g., 'Administrator', 'Rukovodilac projekta', 'Istraživač'
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
('Istraživač'),
('Organizator projekta');
