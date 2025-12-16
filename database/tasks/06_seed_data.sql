-- ============================================================================
-- Seed data for testing triggers, functions, indexes and reports
-- This script is IDEMPOTENT - safe to run multiple times
-- Run this file after the schema (01_schema.sql) has been applied.
-- ============================================================================

-- Delete existing test data (in reverse order of dependencies)
DELETE FROM LogAktivnosti WHERE korisnik_id IN (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime IN ('admin','lead','istra1','istra2','org'));
DELETE FROM DokumentTagovi WHERE dokument_id IN (SELECT dokument_id FROM Dokumenti WHERE naziv_dokumenta IN ('Studija_A.pdf','Plan_B.docx'));
DELETE FROM Tagovi WHERE naziv_taga IN ('AI','ML','Plan','Eksperiment');
DELETE FROM Dokumenti WHERE naziv_dokumenta IN ('Studija_A.pdf','Plan_B.docx');
DELETE FROM Zadaci WHERE projekat_id IN (SELECT projekat_id FROM Projekti WHERE naziv_projekta IN ('Projekt A','Projekt B'));
DELETE FROM ClanoviProjekta WHERE projekat_id IN (SELECT projekat_id FROM Projekti WHERE naziv_projekta IN ('Projekt A','Projekt B'));
DELETE FROM Projekti WHERE naziv_projekta IN ('Projekt A','Projekt B');
DELETE FROM Folderi WHERE naziv_foldera IN ('Root','ProjektA_docs');
DELETE FROM Korisnici WHERE korisnicko_ime IN ('admin','lead','istra1','istra2','org');

COMMIT;

-- Insert test users (passwords are placeholders)
INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status)
VALUES ('admin','admin@example.com','$hash$admin','Nikola','Adminovic',1,'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status)
VALUES ('lead','lead@example.com','$hash$lead','Mitar','Markovic',2,'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status)
VALUES ('istra1','istra1@example.com','$hash$1','Jovan','Jovanovic',3,'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status)
VALUES ('istra2','istra2@example.com','$hash$2','Petar','Petrovic',3,'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status)
VALUES ('org','org@example.com','$hash$org','Ana','Anic',4,'neaktivan');

-- Create some folders
INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id)
VALUES ('Root', NULL, (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='admin'));

INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id)
VALUES ('ProjektA_docs', (SELECT MAX(folder_id) FROM Folderi WHERE naziv_foldera='Root'), (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'));

-- Create projects
INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id)
VALUES ('Projekt A', 'Istraživački projekat A', TO_DATE('2025-02-01','YYYY-MM-DD'), NULL, 'Aktivan', (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'), (SELECT radni_tok_id FROM RadniTokovi WHERE naziv = 'Standardni projektni tok'));

INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id)
VALUES ('Projekt B', 'Razvojni projekat B', TO_DATE('2025-03-15','YYYY-MM-DD'), NULL, 'Aktivan', (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'), (SELECT radni_tok_id FROM RadniTokovi WHERE naziv = 'Istrazivacki tok'));

-- Add members to projects
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt A'), (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'));

INSERT INTO ClanoviProjekta (projekat_id, korisnik_id)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt A'), (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra1'));

INSERT INTO ClanoviProjekta (projekat_id, korisnik_id)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt A'), (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra2'));

INSERT INTO ClanoviProjekta (projekat_id, korisnik_id)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt B'), (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'));

-- Create tasks for Projekt A (use phase names to find faza_id)
INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt A'),
 (SELECT faza_id FROM Faze WHERE radni_tok_id = (SELECT radni_tok_id FROM RadniTokovi WHERE naziv = 'Standardni projektni tok') AND naziv_faze = 'Planiranje'),
 'Definisanje ciljeva', 'Pripremiti ciljeve i metodologiju', (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra1'), TO_DATE('2025-04-01','YYYY-MM-DD'), 'Visok', 20);

INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt A'),
 (SELECT faza_id FROM Faze WHERE radni_tok_id = (SELECT radni_tok_id FROM RadniTokovi WHERE naziv = 'Standardni projektni tok') AND naziv_faze = 'Razvoj'),
 'Eksperimentalni kod', 'Implementacija prototipa', (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra2'), TO_DATE('2025-05-10','YYYY-MM-DD'), 'Srednji', 50);

INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt A'),
 (SELECT faza_id FROM Faze WHERE radni_tok_id = (SELECT radni_tok_id FROM RadniTokovi WHERE naziv = 'Standardni projektni tok') AND naziv_faze = 'Završeno'),
 'Zavrsni izvestaj', 'Finalizovati izvestaj', (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'), TO_DATE('2025-06-01','YYYY-MM-DD'), 'Nizak', 100);

-- Create tasks for Projekt B
INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt B'),
 (SELECT faza_id FROM Faze WHERE radni_tok_id = (SELECT radni_tok_id FROM RadniTokovi WHERE naziv = 'Istrazivacki tok') AND naziv_faze = 'Prikupljanje podataka'),
 'Sakupljanje uzoraka', 'Prikupljati i skladištiti podatke', (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra1'), TO_DATE('2025-05-20','YYYY-MM-DD'), 'Visok', 40);

-- Create some documents
INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, kljucne_reci, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt A'), 'Studija_A.pdf', (SELECT MAX(folder_id) FROM Folderi WHERE naziv_foldera='ProjektA_docs'), 'Studija primera', 'PDF', 'Srpski', 'AI,ML,eksperiment', (SELECT radni_tok_id FROM RadniTokovi WHERE naziv='Dokumentacioni tok'), (SELECT faza_id FROM Faze WHERE radni_tok_id = (SELECT radni_tok_id FROM RadniTokovi WHERE naziv='Dokumentacioni tok') AND naziv_faze='Kreiranje'), (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra1'));

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, kljucne_reci, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id)
VALUES ((SELECT MAX(projekat_id) FROM Projekti WHERE naziv_projekta='Projekt B'), 'Plan_B.docx', NULL, 'Plan rada projekta B', 'DOCX', 'Srpski', 'plan,projekat,B', NULL, NULL, (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'));

-- Tags and document-tag relations
INSERT INTO Tagovi (naziv_taga) VALUES ('AI');
INSERT INTO Tagovi (naziv_taga) VALUES ('ML');
INSERT INTO Tagovi (naziv_taga) VALUES ('Plan');
INSERT INTO Tagovi (naziv_taga) VALUES ('Eksperiment');

INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES ((SELECT MAX(dokument_id) FROM Dokumenti WHERE naziv_dokumenta='Studija_A.pdf'), (SELECT tag_id FROM Tagovi WHERE naziv_taga='AI'));
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES ((SELECT MAX(dokument_id) FROM Dokumenti WHERE naziv_dokumenta='Studija_A.pdf'), (SELECT tag_id FROM Tagovi WHERE naziv_taga='ML'));
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES ((SELECT MAX(dokument_id) FROM Dokumenti WHERE naziv_dokumenta='Plan_B.docx'), (SELECT tag_id FROM Tagovi WHERE naziv_taga='Plan'));

COMMIT;

PROMPT ============================================================================
PROMPT Seed data loaded successfully!
PROMPT ============================================================================

-- A few manual log entries (to have data for activity views and to test triggers)
INSERT INTO LogAktivnosti (log_id, korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, naziv_entiteta, opis)
VALUES (log_aktivnosti_seq.NEXTVAL, (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra1'), 'IMPORT', 'DOKUMENT', (SELECT dokument_id FROM Dokumenti WHERE naziv_dokumenta='Studija_A.pdf'), 'Studija_A.pdf', 'Uvoz početne verzije');

INSERT INTO LogAktivnosti (log_id, korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, naziv_entiteta, opis)
VALUES (log_aktivnosti_seq.NEXTVAL, (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'), 'KREIRANJE', 'PROJEKAT', (SELECT projekat_id FROM Projekti WHERE naziv_projekta='Projekt A'), 'Projekt A', 'Kreiran projekat');

-- Test queries and function/procedure calls (examples to run manually)
-- 1) Percentage of finished tasks for Projekt A
-- SELECT procenat_zavrsenih_zadataka((SELECT projekat_id FROM Projekti WHERE naziv_projekta='Projekt A'));

-- 2) Number of active members in Projekt A
-- SELECT broj_aktivnih_clanova((SELECT projekat_id FROM Projekti WHERE naziv_projekta='Projekt A'));

-- 3) Run reports (will print to DBMS_OUTPUT if enabled)
-- CALL izvestaj_projekti();
-- CALL izvestaj_istrazivaci();

-- Notes:
-- - Adjust hashes or other fields if your schema enforces additional constraints.
-- - If some INSERTs fail due to missing optional columns in your local schema (omitted lines in schema), remove/adjust those rows accordingly.
