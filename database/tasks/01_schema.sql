-- Research Institute Information System Database Schema
-- Converted for Oracle Database
-- Encoding: UTF-8

-- ============================================================================
-- CLEANUP: Drop existing objects (if re-running script)
-- ============================================================================

-- Drop views first (they depend on tables)
BEGIN EXECUTE IMMEDIATE 'DROP VIEW SkorijeAktivnosti'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP VIEW StatistikaAktivnosti'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP VIEW StatistikaDokumenata'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP VIEW v_zadaci_sa_detaljima'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP VIEW v_dokumenti_sa_verzijama'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP VIEW v_aktivni_projekti'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- Drop tables in reverse order of dependencies
BEGIN EXECUTE IMMEDIATE 'DROP TABLE LogAktivnosti CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE IstorijaFazaDokumenta CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE DozvoleDokumenata CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE DokumentTagovi CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Tagovi CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE MetaPodaci CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE LLMSazeci CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE VerzijeDokumenata CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Dokumenti CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Folderi CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE ZahteviPromeneFaze CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE KomentariZadataka CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Zadaci CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE ClanoviProjekta CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Projekti CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Faze CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE RadniTokovi CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Korisnici CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Uloge CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- Drop sequences
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE log_aktivnosti_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE istorija_faza_dok_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE dozvole_dokumenata_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE tagovi_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE meta_podaci_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE llm_sazeci_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE verzije_dokumenata_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE dokumenti_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE folderi_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE zahtevi_promene_faze_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE komentari_zadataka_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE zadaci_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE projekti_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE faze_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE radni_tokovi_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE korisnici_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE uloge_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- ============================================================================
-- SCHEMA CREATION
-- ============================================================================

-- Module 1: User and Role Management

-- Table for defining user roles (Administrator, Manager, Researcher)
CREATE TABLE Uloge (
    uloga_id NUMBER PRIMARY KEY,
    naziv_uloge VARCHAR2(50) UNIQUE NOT NULL -- e.g., 'Administrator', 'Rukovodilac projekta', 'Istrazivac'
);

CREATE SEQUENCE uloge_seq START WITH 1 INCREMENT BY 1;

-- Main table for system users
CREATE TABLE Korisnici (
    korisnik_id NUMBER PRIMARY KEY,
    korisnicko_ime VARCHAR2(100) UNIQUE NOT NULL,
    email VARCHAR2(100) UNIQUE NOT NULL,
    hash_sifre VARCHAR2(255) NOT NULL,
    ime VARCHAR2(100),
    prezime VARCHAR2(100),
    uloga_id NUMBER NOT NULL,
    status VARCHAR2(20) DEFAULT 'aktivan', -- e.g., 'aktivan', 'neaktivan'
    poslednja_prijava TIMESTAMP,
    kreiran_datuma TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_korisnici_uloga FOREIGN KEY (uloga_id) REFERENCES Uloge(uloga_id)
);

CREATE SEQUENCE korisnici_seq START WITH 1 INCREMENT BY 1;

-- Module 2: Project, Task and Documentation Management

-- Table for defining workflows (e.g., for projects or documents)
CREATE TABLE RadniTokovi (
    radni_tok_id NUMBER PRIMARY KEY,
    naziv VARCHAR2(255) NOT NULL,
    tip_toka VARCHAR2(50) NOT NULL CHECK (tip_toka IN ('PROJEKAT', 'DOKUMENTACIJA')),
    opis VARCHAR2(255),
    da_li_je_sablon NUMBER(1) DEFAULT 0, -- 1 = TRUE, 0 = FALSE; Whether this workflow can be used as a template
    CONSTRAINT uq_radni_tokovi UNIQUE(naziv, tip_toka)
);

CREATE SEQUENCE radni_tokovi_seq START WITH 1 INCREMENT BY 1;

-- Phases within a workflow (e.g., Planning, Development, Completed)
CREATE TABLE Faze (
    faza_id NUMBER PRIMARY KEY,
    radni_tok_id NUMBER NOT NULL,
    naziv_faze VARCHAR2(100) NOT NULL,
    redosled NUMBER NOT NULL, -- For sorting phases within workflow
    CONSTRAINT fk_faze_radni_tok FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id) ON DELETE CASCADE
);

CREATE SEQUENCE faze_seq START WITH 1 INCREMENT BY 1;

-- Table containing basic project information
CREATE TABLE Projekti (
    projekat_id NUMBER PRIMARY KEY,
    naziv_projekta VARCHAR2(255) NOT NULL,
    opis VARCHAR2(255),
    datum_pocetka DATE,
    datum_zavrsetka DATE,
    status VARCHAR2(50) DEFAULT 'Aktivan', -- e.g., 'Aktivan', 'Završen', 'Otkazan'
    rukovodilac_id NUMBER, -- User who created and manages the project
    radni_tok_id NUMBER, -- Workflow applied to tasks in this project
    CONSTRAINT fk_projekti_rukovodilac FOREIGN KEY (rukovodilac_id) REFERENCES Korisnici(korisnik_id),
    CONSTRAINT fk_projekti_radni_tok FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id)
);

CREATE SEQUENCE projekti_seq START WITH 1 INCREMENT BY 1;

-- Table linking users to projects (team members)
CREATE TABLE ClanoviProjekta (
    projekat_id NUMBER NOT NULL,
    korisnik_id NUMBER NOT NULL,
    CONSTRAINT pk_clanovi_projekta PRIMARY KEY (projekat_id, korisnik_id),
    CONSTRAINT fk_clanovi_projekat FOREIGN KEY (projekat_id) REFERENCES Projekti(projekat_id) ON DELETE CASCADE,
    CONSTRAINT fk_clanovi_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id) ON DELETE CASCADE
);

-- Table for tasks within projects
CREATE TABLE Zadaci (
    zadatak_id NUMBER PRIMARY KEY,
    projekat_id NUMBER NOT NULL,
    faza_id NUMBER NOT NULL, -- Current phase of the task
    naziv_zadatka VARCHAR2(255) NOT NULL,
    opis VARCHAR2(255),
    dodeljen_korisniku_id NUMBER,
    rok DATE,
    prioritet VARCHAR2(50), -- e.g., 'Nizak', 'Srednji', 'Visok'
    progres NUMBER DEFAULT 0,
    kreiran_datuma TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_zadaci_projekat FOREIGN KEY (projekat_id) REFERENCES Projekti(projekat_id) ON DELETE CASCADE,
    CONSTRAINT fk_zadaci_faza FOREIGN KEY (faza_id) REFERENCES Faze(faza_id),
    CONSTRAINT fk_zadaci_korisnik FOREIGN KEY (dodeljen_korisniku_id) REFERENCES Korisnici(korisnik_id)
);

CREATE SEQUENCE zadaci_seq START WITH 1 INCREMENT BY 1;

-- Table for task comments
CREATE TABLE KomentariZadataka (
    komentar_id NUMBER PRIMARY KEY,
    zadatak_id NUMBER NOT NULL,
    korisnik_id NUMBER NOT NULL,
    tekst_komentara VARCHAR2(255) NOT NULL,
    datuma_kreiranja TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_komentari_zadatak FOREIGN KEY (zadatak_id) REFERENCES Zadaci(zadatak_id) ON DELETE CASCADE,
    CONSTRAINT fk_komentari_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id)
);

CREATE SEQUENCE komentari_zadataka_seq START WITH 1 INCREMENT BY 1;

-- Table for recording phase change requests from Researchers
CREATE TABLE ZahteviPromeneFaze (
    zahtev_id NUMBER PRIMARY KEY,
    zadatak_id NUMBER NOT NULL,
    podnosilac_zahteva_id NUMBER NOT NULL,
    zahtevana_faza_id NUMBER NOT NULL,
    status VARCHAR2(50) DEFAULT 'Na cekanju', -- e.g., 'Na cekanju', 'Odobren', 'Odbijen'
    komentar VARCHAR2(255),
    datum_kreiranja TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_zahtevi_zadatak FOREIGN KEY (zadatak_id) REFERENCES Zadaci(zadatak_id) ON DELETE CASCADE,
    CONSTRAINT fk_zahtevi_podnosilac FOREIGN KEY (podnosilac_zahteva_id) REFERENCES Korisnici(korisnik_id),
    CONSTRAINT fk_zahtevi_faza FOREIGN KEY (zahtevana_faza_id) REFERENCES Faze(faza_id)
);

CREATE SEQUENCE zahtevi_promene_faze_seq START WITH 1 INCREMENT BY 1;

-- Module 3: Document and Metadata Management

CREATE TABLE Folderi (
    folder_id NUMBER PRIMARY KEY,
    naziv_foldera VARCHAR2(255) NOT NULL,
    roditelj_folder_id NUMBER,
    vlasnik_id NUMBER NOT NULL,
    CONSTRAINT fk_folderi_roditelj FOREIGN KEY (roditelj_folder_id) REFERENCES Folderi(folder_id) ON DELETE CASCADE,
    CONSTRAINT fk_folderi_vlasnik FOREIGN KEY (vlasnik_id) REFERENCES Korisnici(korisnik_id)
);

CREATE SEQUENCE folderi_seq START WITH 1 INCREMENT BY 1;

-- Table for basic document information
CREATE TABLE Dokumenti (
    dokument_id NUMBER PRIMARY KEY,
    projekat_id NUMBER, -- Optional, document can be part of a project
    naziv_dokumenta VARCHAR2(255) NOT NULL,
    folder_id NUMBER,
    opis VARCHAR2(255),
    tip_dokumenta VARCHAR2(50), -- e.g., 'Istraživački rad', 'PDF', 'CSV'
    jezik_dokumenta VARCHAR2(50),
    kljucne_reci VARCHAR2(255), -- Free-form keywords separated by commas
    radni_tok_id NUMBER,
    trenutna_faza_id NUMBER, -- If document follows workflow
    kreirao_korisnik_id NUMBER NOT NULL,
    datuma_postavke TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    poslednja_izmena TIMESTAMP,
    CONSTRAINT fk_dokumenti_projekat FOREIGN KEY (projekat_id) REFERENCES Projekti(projekat_id),
    CONSTRAINT fk_dokumenti_radni_tok FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id),
    CONSTRAINT fk_dokumenti_faza FOREIGN KEY (trenutna_faza_id) REFERENCES Faze(faza_id),
    CONSTRAINT fk_dokumenti_kreirao FOREIGN KEY (kreirao_korisnik_id) REFERENCES Korisnici(korisnik_id),
    CONSTRAINT fk_dokumenti_folder FOREIGN KEY (folder_id) REFERENCES Folderi(folder_id) ON DELETE SET NULL
);

CREATE SEQUENCE dokumenti_seq START WITH 1 INCREMENT BY 1;

-- Table for tracking document versions
CREATE TABLE VerzijeDokumenata (
    verzija_id NUMBER PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    verzija_oznaka VARCHAR2(50), -- e.g., 'v1.0', 'v1.1 final'
    putanja_do_fajla VARCHAR2(1024) NOT NULL,
    velicina_fajla_MB NUMBER(10, 2),
    postavio_korisnik_id NUMBER NOT NULL,
    datuma_postavke TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_verzije_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    CONSTRAINT fk_verzije_postavio FOREIGN KEY (postavio_korisnik_id) REFERENCES Korisnici(korisnik_id)
);

CREATE SEQUENCE verzije_dokumenata_seq START WITH 1 INCREMENT BY 1;

CREATE TABLE LLMSazeci (
    sazetak_id NUMBER PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    verzija_oznaka VARCHAR2(50),
    sazetak VARCHAR2(255) NOT NULL,
    datum_kreiranja TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sazeci_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE
);

CREATE SEQUENCE llm_sazeci_seq START WITH 1 INCREMENT BY 1;

-- Table for additional metadata, allows flexibility
CREATE TABLE MetaPodaci (
    meta_id NUMBER PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    kljuc VARCHAR2(100) NOT NULL, -- e.g., 'ISO Broj', 'Izvorni URL', 'LLM sažetak'
    vrednost VARCHAR2(255),
    CONSTRAINT fk_meta_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE
);

CREATE SEQUENCE meta_podaci_seq START WITH 1 INCREMENT BY 1;

-- Table for tags for easier search
CREATE TABLE Tagovi (
    tag_id NUMBER PRIMARY KEY,
    naziv_taga VARCHAR2(100) UNIQUE NOT NULL
);

CREATE SEQUENCE tagovi_seq START WITH 1 INCREMENT BY 1;

-- Links documents to tags (many-to-many relationship)
CREATE TABLE DokumentTagovi (
    dokument_id NUMBER NOT NULL,
    tag_id NUMBER NOT NULL,
    CONSTRAINT pk_dokument_tagovi PRIMARY KEY (dokument_id, tag_id),
    CONSTRAINT fk_dok_tagovi_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    CONSTRAINT fk_dok_tagovi_tag FOREIGN KEY (tag_id) REFERENCES Tagovi(tag_id) ON DELETE CASCADE
);

-- Table defining access rights (read, write, delete) for users over documents
CREATE TABLE DozvoleDokumenata (
    dozvola_id NUMBER PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    korisnik_id NUMBER NOT NULL,
    moze_citati NUMBER(1) DEFAULT 1, -- 1 = TRUE, 0 = FALSE
    moze_menjati NUMBER(1) DEFAULT 0,
    moze_brisati NUMBER(1) DEFAULT 0,
    CONSTRAINT uq_dozvole_dok_korisnik UNIQUE (dokument_id, korisnik_id), -- Each user has one permission row per document
    CONSTRAINT fk_dozvole_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    CONSTRAINT fk_dozvole_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id) ON DELETE CASCADE
);

CREATE SEQUENCE dozvole_dokumenata_seq START WITH 1 INCREMENT BY 1;

-- Table for tracking phase history - applies ONLY to project documentation
CREATE TABLE IstorijaFazaDokumenta (
    istorija_id NUMBER PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    prethodna_faza_id NUMBER,
    nova_faza_id NUMBER NOT NULL,
    korisnik_id NUMBER NOT NULL,
    datum_promene TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_istorija_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    CONSTRAINT fk_istorija_prethodna FOREIGN KEY (prethodna_faza_id) REFERENCES Faze(faza_id),
    CONSTRAINT fk_istorija_nova FOREIGN KEY (nova_faza_id) REFERENCES Faze(faza_id),
    CONSTRAINT fk_istorija_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id)
);

CREATE SEQUENCE istorija_faza_dok_seq START WITH 1 INCREMENT BY 1;

-- Module 4: Analytics and Logging
-- (LogAktivnosti table defined later with complete schema)


-- Insert default roles (Oracle: single-row INSERTs)
INSERT INTO Uloge (uloga_id, naziv_uloge) VALUES (uloge_seq.NEXTVAL, 'Administrator');
INSERT INTO Uloge (uloga_id, naziv_uloge) VALUES (uloge_seq.NEXTVAL, 'Rukovodilac projekta');
INSERT INTO Uloge (uloga_id, naziv_uloge) VALUES (uloge_seq.NEXTVAL, 'Istrazivac');
INSERT INTO Uloge (uloga_id, naziv_uloge) VALUES (uloge_seq.NEXTVAL, 'Organizator projekta');

-- Create some useful views
CREATE OR REPLACE VIEW v_aktivni_projekti AS
SELECT 
    p.projekat_id,
    p.naziv_projekta,
    p.opis,
    p.datum_pocetka,
    p.datum_zavrsetka,
    k.ime || ' ' || k.prezime as rukovodilac_ime,
    k.email as rukovodilac_email,
    rt.naziv as radni_tok_naziv,
    COUNT(z.zadatak_id) as ukupno_zadataka,
    SUM(CASE WHEN z.progres = 100 THEN 1 ELSE 0 END) as zavrsenih_zadataka
FROM Projekti p
LEFT JOIN Korisnici k ON p.rukovodilac_id = k.korisnik_id
LEFT JOIN RadniTokovi rt ON p.radni_tok_id = rt.radni_tok_id
LEFT JOIN Zadaci z ON p.projekat_id = z.projekat_id
WHERE p.status = 'Aktivan'
GROUP BY p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka, p.datum_zavrsetka, k.ime, k.prezime, k.email, rt.naziv, rt.radni_tok_id;

CREATE OR REPLACE VIEW v_dokumenti_sa_verzijama AS
SELECT 
    d.dokument_id,
    d.naziv_dokumenta,
    d.opis,
    d.tip_dokumenta,
    d.jezik_dokumenta,
    k.ime || ' ' || k.prezime as kreirao_ime,
    d.datuma_postavke,
    d.poslednja_izmena,
    COUNT(vd.verzija_id) as broj_verzija,
    MAX(vd.verzija_oznaka) as poslednja_verzija
FROM Dokumenti d
LEFT JOIN Korisnici k ON d.kreirao_korisnik_id = k.korisnik_id
LEFT JOIN VerzijeDokumenata vd ON d.dokument_id = vd.dokument_id
GROUP BY d.dokument_id, d.naziv_dokumenta, d.opis, d.tip_dokumenta, d.jezik_dokumenta, k.ime, k.prezime, d.datuma_postavke, d.poslednja_izmena, k.korisnik_id;

CREATE OR REPLACE VIEW v_zadaci_sa_detaljima AS
SELECT 
    z.zadatak_id,
    z.naziv_zadatka,
    z.opis,
    z.rok,
    z.prioritet,
    z.progres,
    p.naziv_projekta,
    k.ime || ' ' || k.prezime as dodeljen_korisniku,
    f.naziv_faze,
    f.redosled as faza_redosled
FROM Zadaci z
JOIN Projekti p ON z.projekat_id = p.projekat_id
LEFT JOIN Korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
JOIN Faze f ON z.faza_id = f.faza_id;

-- Insert default workflows (Oracle: use 1 for TRUE, 0 for FALSE in BOOLEAN/NUMBER columns)
INSERT INTO RadniTokovi (radni_tok_id, naziv, tip_toka, opis, da_li_je_sablon) VALUES (radni_tokovi_seq.NEXTVAL, 'Standardni projektni tok', 'PROJEKAT', 'Osnovni radni tok za projekte', 1);
INSERT INTO RadniTokovi (radni_tok_id, naziv, tip_toka, opis, da_li_je_sablon) VALUES (radni_tokovi_seq.NEXTVAL, 'Istrazivacki tok', 'PROJEKAT', 'Tok za istrazivacke projekte', 1);
INSERT INTO RadniTokovi (radni_tok_id, naziv, tip_toka, opis, da_li_je_sablon) VALUES (radni_tokovi_seq.NEXTVAL, 'Dokumentacioni tok', 'DOKUMENTACIJA', 'Tok za upravljanje dokumentima', 1);

-- Insert default phases for project workflow
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 1, 'Planiranje', 1);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 1, 'Analiza', 2);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 1, 'Razvoj', 3);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 1, 'Testiranje', 4);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 1, 'Završeno', 5);

-- Insert default phases for research workflow
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 2, 'Definisanje istrazivanja', 1);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 2, 'Prikupljanje podataka', 2);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 2, 'Analiza podataka', 3);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 2, 'Pisanje izvestaja', 4);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 2, 'Publikovanje', 5);

-- Insert default phases for documentation workflow
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 3, 'Kreiranje', 1);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 3, 'Revizija', 2);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 3, 'Odobravanje', 3);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 3, 'Finalizovanje', 4);
INSERT INTO Faze (faza_id, radni_tok_id, naziv_faze, redosled) VALUES (faze_seq.NEXTVAL, 3, 'Arhiviranje', 5);

-- ============================================================================
-- Activity Log System
-- ============================================================================

-- Table for logging all system activities
CREATE TABLE LogAktivnosti (
    log_id NUMBER PRIMARY KEY,
    korisnik_id NUMBER,
    tip_aktivnosti VARCHAR2(50) NOT NULL, -- 'UPLOAD', 'EDIT', 'DELETE', 'VIEW', 'LOGIN', 'LOGOUT', 'SHARE', 'CREATE_PROJECT', etc.
    entitet_tip VARCHAR2(50), -- 'DOKUMENT', 'PROJEKAT', 'ZADATAK', 'KORISNIK', etc.
    entitet_id NUMBER, -- ID of the affected entity
    naziv_entiteta VARCHAR2(255), -- Name/title of the entity for easier querying
    opis VARCHAR2(255), -- Detailed description of the activity
    ip_adresa VARCHAR2(45), -- IPv4 or IPv6
    user_agent VARCHAR2(255), -- Browser/client information
    rezultat VARCHAR2(20) DEFAULT 'SUCCESS', -- 'SUCCESS', 'FAILED', 'PENDING'
    dodatne_informacije VARCHAR2(255), -- Additional metadata in JSON format
    kreiran_datuma TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_log_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id) ON DELETE SET NULL
);

CREATE SEQUENCE log_aktivnosti_seq START WITH 1 INCREMENT BY 1;

-- ============================================================================
-- Analytics Views
-- ============================================================================

-- View for document statistics
CREATE OR REPLACE VIEW StatistikaDokumenata AS
SELECT 
    COUNT(*) as ukupno_dokumenata,
    SUM(CASE WHEN datuma_postavke >= TRUNC(SYSDATE) - 30 THEN 1 ELSE 0 END) as novih_dokumenata_mesecno,
    COUNT(DISTINCT kreirao_korisnik_id) as broj_autora,
    SUM(CASE WHEN projekat_id IS NOT NULL THEN 1 ELSE 0 END) as broj_projekata_sa_dokumentima,
    AVG(broj_verzija) as prosecno_verzija_po_dokumentu
FROM (
    SELECT 
        d.dokument_id,
        d.kreirao_korisnik_id,
        d.projekat_id,
        d.datuma_postavke,
        COUNT(v.verzija_id) as broj_verzija
    FROM Dokumenti d
    LEFT JOIN VerzijeDokumenata v ON d.dokument_id = v.dokument_id
    GROUP BY d.dokument_id, d.kreirao_korisnik_id, d.projekat_id, d.datuma_postavke
) doc_stats;

-- View for activity statistics
CREATE OR REPLACE VIEW StatistikaAktivnosti AS
SELECT 
    tip_aktivnosti,
    COUNT(*) as broj_aktivnosti,
    COUNT(DISTINCT korisnik_id) as broj_korisnika,
    SUM(CASE WHEN kreiran_datuma >= TRUNC(SYSDATE) THEN 1 ELSE 0 END) as danas,
    SUM(CASE WHEN kreiran_datuma >= TRUNC(SYSDATE) - 7 THEN 1 ELSE 0 END) as ove_nedelje,
    SUM(CASE WHEN kreiran_datuma >= TRUNC(SYSDATE) - 30 THEN 1 ELSE 0 END) as ovog_meseca,
    MAX(kreiran_datuma) as poslednja_aktivnost
FROM LogAktivnosti
GROUP BY tip_aktivnosti;

-- View for recent activity feed
CREATE OR REPLACE VIEW SkorijeAktivnosti AS
SELECT * FROM (
    SELECT 
        l.log_id,
        l.korisnik_id,
        COALESCE(k.ime || ' ' || k.prezime, k.korisnicko_ime, 'System') as korisnik_ime,
        l.tip_aktivnosti,
        l.entitet_tip,
        l.entitet_id,
        l.naziv_entiteta,
        l.opis,
        l.rezultat,
        l.kreiran_datuma
    FROM LogAktivnosti l
    LEFT JOIN Korisnici k ON l.korisnik_id = k.korisnik_id
    ORDER BY l.kreiran_datuma DESC
)
WHERE ROWNUM <= 100;
