-- Research Institute Information System Database Schema
-- Created for Oracle Database
-- Encoding: UTF-8

-- Module 1: User and Role Management

-- Table for defining user roles (Administrator, Manager, Researcher)
CREATE TABLE Uloge (
    uloga_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    naziv_uloge VARCHAR2(50) UNIQUE NOT NULL -- e.g., 'Administrator', 'Rukovodilac projekta', 'Istrazivac'
);

-- Main table for system users
CREATE TABLE Korisnici (
    korisnik_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    korisnicko_ime VARCHAR2(100) UNIQUE NOT NULL,
    email VARCHAR2(100) UNIQUE NOT NULL,
    hash_sifre VARCHAR2(255) NOT NULL,
    ime VARCHAR2(100),
    prezime VARCHAR2(100),
    uloga_id NUMBER NOT NULL,
    status VARCHAR2(20) DEFAULT 'aktivan', -- e.g., 'aktivan', 'neaktivan'
    poslednja_prijava TIMESTAMP,
    kreiran_datuma TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_korisnici_uloge FOREIGN KEY (uloga_id) REFERENCES Uloge(uloga_id)
);

-- Module 2: Project, Task and Documentation Management

-- Table for defining workflows (e.g., for projects or documents)
CREATE TABLE RadniTokovi (
    radni_tok_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    naziv VARCHAR2(255) NOT NULL,
    tip_toka VARCHAR2(50) NOT NULL CHECK (tip_toka IN ('PROJEKAT', 'DOKUMENTACIJA')),
    opis CLOB,
    da_li_je_sablon NUMBER(1) DEFAULT 0, -- 0=FALSE, 1=TRUE
    CONSTRAINT uk_radni_tokovi UNIQUE(naziv, tip_toka)
);

-- Phases within a workflow (e.g., Planning, Development, Completed)
CREATE TABLE Faze (
    faza_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    radni_tok_id NUMBER NOT NULL,
    naziv_faze VARCHAR2(100) NOT NULL,
    redosled NUMBER NOT NULL, -- For sorting phases within workflow
    CONSTRAINT fk_faze_radni_tok FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id) ON DELETE CASCADE
);

-- Table containing basic project information
CREATE TABLE Projekti (
    projekat_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    naziv_projekta VARCHAR2(255) NOT NULL,
    opis CLOB,
    datum_pocetka DATE,
    datum_zavrsetka DATE,
    status VARCHAR2(50) DEFAULT 'Aktivan', -- e.g., 'Aktivan', 'Završen', 'Otkazan'
    rukovodilac_id NUMBER, -- User who created and manages the project
    radni_tok_id NUMBER, -- Workflow applied to tasks in this project
    CONSTRAINT fk_projekti_rukovodilac FOREIGN KEY (rukovodilac_id) REFERENCES Korisnici(korisnik_id),
    CONSTRAINT fk_projekti_radni_tok FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id)
);

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
    zadatak_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    projekat_id NUMBER NOT NULL,
    faza_id NUMBER NOT NULL, -- Current phase of the task
    naziv_zadatka VARCHAR2(255) NOT NULL,
    opis CLOB,
    dodeljen_korisniku_id NUMBER,
    rok DATE,
    prioritet VARCHAR2(50), -- e.g., 'Nizak', 'Srednji', 'Visok'
    progres NUMBER DEFAULT 0,
    kreiran_datuma TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_zadaci_projekat FOREIGN KEY (projekat_id) REFERENCES Projekti(projekat_id) ON DELETE CASCADE,
    CONSTRAINT fk_zadaci_faza FOREIGN KEY (faza_id) REFERENCES Faze(faza_id),
    CONSTRAINT fk_zadaci_korisnik FOREIGN KEY (dodeljen_korisniku_id) REFERENCES Korisnici(korisnik_id)
);

-- Table for task comments
CREATE TABLE KomentariZadataka (
    komentar_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    zadatak_id NUMBER NOT NULL,
    korisnik_id NUMBER NOT NULL,
    tekst_komentara CLOB NOT NULL,
    datuma_kreiranja TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_komentari_zadatak FOREIGN KEY (zadatak_id) REFERENCES Zadaci(zadatak_id) ON DELETE CASCADE,
    CONSTRAINT fk_komentari_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Table for recording phase change requests from Researchers
CREATE TABLE ZahteviPromeneFaze (
    zahtev_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    zadatak_id NUMBER NOT NULL,
    podnosilac_zahteva_id NUMBER NOT NULL,
    zahtevana_faza_id NUMBER NOT NULL,
    status VARCHAR2(50) DEFAULT 'Na cekanju', -- e.g., 'Na cekanju', 'Odobren', 'Odbijen'
    komentar CLOB,
    datum_kreiranja TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_zahtevi_zadatak FOREIGN KEY (zadatak_id) REFERENCES Zadaci(zadatak_id) ON DELETE CASCADE,
    CONSTRAINT fk_zahtevi_podnosilac FOREIGN KEY (podnosilac_zahteva_id) REFERENCES Korisnici(korisnik_id),
    CONSTRAINT fk_zahtevi_faza FOREIGN KEY (zahtevana_faza_id) REFERENCES Faze(faza_id)
);

-- Module 3: Document and Metadata Management

CREATE TABLE Folderi (
    folder_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    naziv_foldera VARCHAR2(255) NOT NULL,
    roditelj_folder_id NUMBER,
    vlasnik_id NUMBER NOT NULL,
    CONSTRAINT fk_folderi_roditelj FOREIGN KEY (roditelj_folder_id) REFERENCES Folderi(folder_id) ON DELETE CASCADE,
    CONSTRAINT fk_folderi_vlasnik FOREIGN KEY (vlasnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Table for basic document information
CREATE TABLE Dokumenti (
    dokument_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    projekat_id NUMBER, -- Optional, document can be part of a project
    naziv_dokumenta VARCHAR2(255) NOT NULL,
    folder_id NUMBER,
    opis CLOB,
    tip_dokumenta VARCHAR2(50), -- e.g., 'Istraživački rad', 'PDF', 'CSV'
    jezik_dokumenta VARCHAR2(50),
    kljucne_reci CLOB, -- Free-form keywords separated by commas
    radni_tok_id NUMBER,
    trenutna_faza_id NUMBER, -- If document follows workflow
    kreirao_korisnik_id NUMBER NOT NULL,
    datuma_postavke TIMESTAMP DEFAULT SYSTIMESTAMP,
    poslednja_izmena TIMESTAMP,
    CONSTRAINT fk_dokumenti_projekat FOREIGN KEY (projekat_id) REFERENCES Projekti(projekat_id),
    CONSTRAINT fk_dokumenti_radni_tok FOREIGN KEY (radni_tok_id) REFERENCES RadniTokovi(radni_tok_id),
    CONSTRAINT fk_dokumenti_faza FOREIGN KEY (trenutna_faza_id) REFERENCES Faze(faza_id),
    CONSTRAINT fk_dokumenti_korisnik FOREIGN KEY (kreirao_korisnik_id) REFERENCES Korisnici(korisnik_id),
    CONSTRAINT fk_dokumenti_folder FOREIGN KEY (folder_id) REFERENCES Folderi(folder_id) ON DELETE SET NULL
);

-- Table for tracking document versions
CREATE TABLE VerzijeDokumenata (
    verzija_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    verzija_oznaka VARCHAR2(50), -- e.g., 'v1.0', 'v1.1 final'
    putanja_do_fajla VARCHAR2(1024) NOT NULL,
    velicina_fajla_MB NUMBER(10, 2),
    postavio_korisnik_id NUMBER NOT NULL,
    datuma_postavke TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_verzije_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    CONSTRAINT fk_verzije_korisnik FOREIGN KEY (postavio_korisnik_id) REFERENCES Korisnici(korisnik_id)
);

CREATE TABLE LLMSazeci (
    sazetak_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    verzija_oznaka VARCHAR2(50),
    sazetak CLOB NOT NULL,
    datum_kreiranja TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_sazeci_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE
);

-- Table for additional metadata, allows flexibility
CREATE TABLE MetaPodaci (
    meta_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    kljuc VARCHAR2(100) NOT NULL, -- e.g., 'ISO Broj', 'Izvorni URL', 'LLM sažetak'
    vrednost CLOB,
    CONSTRAINT fk_meta_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE
);

-- Table for tags for easier search
CREATE TABLE Tagovi (
    tag_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    naziv_taga VARCHAR2(100) UNIQUE NOT NULL
);

-- Links documents to tags (many-to-many relationship)
CREATE TABLE DokumentTagovi (
    dokument_id NUMBER NOT NULL,
    tag_id NUMBER NOT NULL,
    CONSTRAINT pk_dokument_tagovi PRIMARY KEY (dokument_id, tag_id),
    CONSTRAINT fk_dok_tag_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    CONSTRAINT fk_dok_tag_tag FOREIGN KEY (tag_id) REFERENCES Tagovi(tag_id) ON DELETE CASCADE
);

-- Table defining access rights (read, write, delete) for users over documents
CREATE TABLE DozvoleDokumenata (
    dozvola_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    korisnik_id NUMBER NOT NULL,
    moze_citati NUMBER(1) DEFAULT 1,
    moze_menjati NUMBER(1) DEFAULT 0,
    moze_brisati NUMBER(1) DEFAULT 0,
    CONSTRAINT uk_dozvole_dok UNIQUE (dokument_id, korisnik_id), -- Each user has one permission row per document
    CONSTRAINT fk_dozvole_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    CONSTRAINT fk_dozvole_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id) ON DELETE CASCADE
);

-- Table for tracking phase history - applies ONLY to project documentation
CREATE TABLE IstorijaFazaDokumenta (
    istorija_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dokument_id NUMBER NOT NULL,
    prethodna_faza_id NUMBER,
    nova_faza_id NUMBER NOT NULL,
    korisnik_id NUMBER NOT NULL,
    datum_promene TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_istorija_dokument FOREIGN KEY (dokument_id) REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    CONSTRAINT fk_istorija_prethodna FOREIGN KEY (prethodna_faza_id) REFERENCES Faze(faza_id),
    CONSTRAINT fk_istorija_nova FOREIGN KEY (nova_faza_id) REFERENCES Faze(faza_id),
    CONSTRAINT fk_istorija_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Module 4: Analytics and Logging

-- Table for logging all important system activities
CREATE TABLE LogAktivnosti (
    log_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    korisnik_id NUMBER,
    tip_aktivnosti VARCHAR2(100) NOT NULL, -- e.g., 'KREIRAN_PROJEKAT', 'UPLOAD_DOKUMENTA', 'PROMENA_FAZE'
    opis CLOB,
    ciljani_entitet VARCHAR2(50), -- e.g., 'Projekat', 'Zadatak', 'Dokument'
    ciljani_id NUMBER,
    datuma TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_log_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id)
);

-- Table for enhanced activity logging
CREATE TABLE LogAktivnostiV2 (
    log_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    korisnik_id NUMBER,
    tip_aktivnosti VARCHAR2(50) NOT NULL, -- 'UPLOAD', 'EDIT', 'DELETE', 'VIEW', 'LOGIN', 'LOGOUT', 'SHARE', 'CREATE_PROJECT', etc.
    entitet_tip VARCHAR2(50), -- 'DOKUMENT', 'PROJEKAT', 'ZADATAK', 'KORISNIK', etc.
    entitet_id NUMBER, -- ID of the affected entity
    naziv_entiteta VARCHAR2(255), -- Name/title of the entity for easier querying
    opis CLOB, -- Detailed description of the activity
    ip_adresa VARCHAR2(45), -- IPv4 or IPv6
    user_agent CLOB, -- Browser/client information
    rezultat VARCHAR2(20) DEFAULT 'SUCCESS', -- 'SUCCESS', 'FAILED', 'PENDING'
    dodatne_informacije CLOB, -- Additional metadata in JSON format
    kreiran_datuma TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_logv2_korisnik FOREIGN KEY (korisnik_id) REFERENCES Korisnici(korisnik_id) ON DELETE SET NULL
);

-- Insert default roles
INSERT INTO Uloge (naziv_uloge) VALUES ('Administrator');
INSERT INTO Uloge (naziv_uloge) VALUES ('Rukovodilac projekta');
INSERT INTO Uloge (naziv_uloge) VALUES ('Istrazivac');
INSERT INTO Uloge (naziv_uloge) VALUES ('Organizator projekta');

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

CREATE INDEX idx_logv2_korisnik ON LogAktivnostiV2(korisnik_id);
CREATE INDEX idx_logv2_tip ON LogAktivnostiV2(tip_aktivnosti);
CREATE INDEX idx_logv2_entitet ON LogAktivnostiV2(entitet_tip, entitet_id);
CREATE INDEX idx_logv2_datum ON LogAktivnostiV2(kreiran_datuma);
CREATE INDEX idx_logv2_rezultat ON LogAktivnostiV2(rezultat);

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
GROUP BY p.projekat_id, p.naziv_projekta, p.opis, p.datum_pocetka, p.datum_zavrsetka, k.ime, k.prezime, k.email, rt.naziv;

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
    k.ime || ' ' || k.prezime as dodeljen_korisniku,
    f.naziv_faze,
    f.redosled as faza_redosled
FROM Zadaci z
JOIN Projekti p ON z.projekat_id = p.projekat_id
LEFT JOIN Korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
JOIN Faze f ON z.faza_id = f.faza_id;

-- Insert default workflows
INSERT INTO RadniTokovi (naziv, tip_toka, opis, da_li_je_sablon) VALUES 
('Standardni projektni tok', 'PROJEKAT', 'Osnovni radni tok za projekte', 1);
INSERT INTO RadniTokovi (naziv, tip_toka, opis, da_li_je_sablon) VALUES 
('Istrazivacki tok', 'PROJEKAT', 'Tok za istrazivacke projekte', 1);
INSERT INTO RadniTokovi (naziv, tip_toka, opis, da_li_je_sablon) VALUES 
('Dokumentacioni tok', 'DOKUMENTACIJA', 'Tok za upravljanje dokumentima', 1);

-- Get the IDs for the workflows
DECLARE
    v_workflow_id1 NUMBER;
    v_workflow_id2 NUMBER;
    v_workflow_id3 NUMBER;
BEGIN
    SELECT radni_tok_id INTO v_workflow_id1 FROM RadniTokovi WHERE naziv = 'Standardni projektni tok';
    SELECT radni_tok_id INTO v_workflow_id2 FROM RadniTokovi WHERE naziv = 'Istrazivacki tok';
    SELECT radni_tok_id INTO v_workflow_id3 FROM RadniTokovi WHERE naziv = 'Dokumentacioni tok';

    -- Insert default phases for project workflow
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id1, 'Planiranje', 1);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id1, 'Analiza', 2);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id1, 'Razvoj', 3);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id1, 'Testiranje', 4);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id1, 'Završeno', 5);

    -- Insert default phases for research workflow
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id2, 'Definisanje istrazivanja', 1);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id2, 'Prikupljanje podataka', 2);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id2, 'Analiza podataka', 3);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id2, 'Pisanje izvestaja', 4);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id2, 'Publikovanje', 5);

    -- Insert default phases for documentation workflow
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id3, 'Kreiranje', 1);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id3, 'Revizija', 2);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id3, 'Odobravanje', 3);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id3, 'Finalizovanje', 4);
    INSERT INTO Faze (radni_tok_id, naziv_faze, redosled) VALUES (v_workflow_id3, 'Arhiviranje', 5);
    
    COMMIT;
END;
/

-- Analytics Views

-- View for document statistics
CREATE OR REPLACE VIEW StatistikaDokumenata AS
SELECT 
    COUNT(*) as ukupno_dokumenata,
    SUM(CASE WHEN datuma_postavke >= TRUNC(SYSDATE) - 30 THEN 1 ELSE 0 END) as novih_dokumenata_mesecno,
    COUNT(DISTINCT kreirao_korisnik_id) as broj_autora,
    COUNT(DISTINCT projekat_id) as broj_projekata_sa_dokumentima,
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
);

-- View for activity statistics
CREATE OR REPLACE VIEW StatistikaAktivnosti AS
SELECT 
    tip_aktivnosti,
    COUNT(*) as broj_aktivnosti,
    COUNT(DISTINCT korisnik_id) as broj_korisnika,
    SUM(CASE WHEN TRUNC(kreiran_datuma) = TRUNC(SYSDATE) THEN 1 ELSE 0 END) as danas,
    SUM(CASE WHEN kreiran_datuma >= TRUNC(SYSDATE) - 7 THEN 1 ELSE 0 END) as ove_nedelje,
    SUM(CASE WHEN kreiran_datuma >= TRUNC(SYSDATE) - 30 THEN 1 ELSE 0 END) as ovog_meseca,
    MAX(kreiran_datuma) as poslednja_aktivnost
FROM LogAktivnostiV2
GROUP BY tip_aktivnosti;

-- View for recent activity feed
CREATE OR REPLACE VIEW SkornjeAktivnosti AS
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
    FROM LogAktivnostiV2 l
    LEFT JOIN Korisnici k ON l.korisnik_id = k.korisnik_id
    ORDER BY l.kreiran_datuma DESC
)
WHERE ROWNUM <= 100;

COMMIT;

-- End of schema
