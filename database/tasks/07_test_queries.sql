-- ============================================================================
-- Test Queries for SBP Project
-- Database: sbp_test
-- ============================================================================

-- SQL*Plus / SQL Developer test script
SET SERVEROUTPUT ON SIZE 1000000

PROMPT === Testing Functions ===

PROMPT Test 1: Procenat zavrsenih zadataka za Projekt A
SELECT procenat_zavrsenih_zadataka((SELECT projekat_id FROM Projekti WHERE naziv_projekta='Projekt A')) procenat FROM DUAL;

PROMPT Test 2: Broj aktivnih clanova u Projektu A
SELECT broj_aktivnih_clanova((SELECT projekat_id FROM Projekti WHERE naziv_projekta='Projekt A')) broj_clanova FROM DUAL;

PROMPT Test 3: Procenat zavrsenih zadataka za Projekt B
SELECT procenat_zavrsenih_zadataka((SELECT projekat_id FROM Projekti WHERE naziv_projekta='Projekt B')) procenat FROM DUAL;

PROMPT === Testing Triggers ===

PROMPT Test 4: Broj log unosa PRE testiranja trigera
SELECT COUNT(*) as logova_pre_testa FROM LogAktivnosti;

PROMPT Test 5: Detalji postojecih log unosa
SELECT tip_aktivnosti, entitet_tip, naziv_entiteta, opis, kreiran_datuma FROM (
  SELECT tip_aktivnosti, entitet_tip, naziv_entiteta, opis, kreiran_datuma FROM LogAktivnosti ORDER BY kreiran_datuma DESC
) WHERE ROWNUM <= 10;

PROMPT Test 6: TRIGGER INSERT - Kreiranje novog projekta (treba da kreira log sa tip_aktivnosti=KREIRANJE)
INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, status, rukovodilac_id, radni_tok_id)
VALUES (
    'Test Projekat',
    'Projekat za testiranje INSERT trigera',
    SYSDATE,
    'Aktivan',
    (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='lead'),
    (SELECT radni_tok_id FROM RadniTokovi WHERE naziv = 'Standardni projektni tok')
);

PROMPT Provera: Da li je kreiran log unos za INSERT?
SELECT COUNT(*) as novih_logova FROM LogAktivnosti WHERE entitet_tip = 'PROJEKAT' AND naziv_entiteta = 'Test Projekat';

PROMPT Test 7: TRIGGER UPDATE - Azuriranje projekta (treba da kreira log sa tip_aktivnosti=UPDATE)
UPDATE Projekti 
SET opis = 'AZURIRANI OPIS - Trigger test' 
WHERE naziv_projekta = 'Test Projekat';

PROMPT Provera: Da li je kreiran log unos za UPDATE?
SELECT tip_aktivnosti, naziv_entiteta, opis, kreiran_datuma 
FROM LogAktivnosti 
WHERE entitet_tip = 'PROJEKAT' AND naziv_entiteta = 'Test Projekat'
ORDER BY kreiran_datuma DESC;

PROMPT Test 8: TRIGGER na Dokumentima - Kreiranje dokumenta (treba da kreira log sa tip_aktivnosti=IMPORT)
INSERT INTO Dokumenti (naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (
    'Test_Dokument.pdf',
    'Dokument za testiranje triggera',
    'PDF',
    'Srpski',
    (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra1')
);

PROMPT Provera: Da li je kreiran log unos za dokument?
SELECT tip_aktivnosti, entitet_tip, naziv_entiteta, opis 
FROM LogAktivnosti 
WHERE naziv_entiteta = 'Test_Dokument.pdf';

PROMPT Test 9: UKUPAN BROJ LOGOVA nakon svih trigger testova
SELECT COUNT(*) as ukupno_logova_posle FROM LogAktivnosti;

PROMPT Rollback test podataka (opciono - zakomentarisano)
-- ROLLBACK;

PROMPT === Testing Reports ===

PROMPT Pokretanje kompleksnog izvestaja sa svim zahtevima (PL/SQL tipovi, kursor, WITH, JOIN, GROUP BY, HAVING, COUNT, SUM)
EXEC kompleksan_izvestaj_projekata;

PROMPT Pokretanje izvestaja o projektima (ispis se pojavljuje u DBMS_OUTPUT)
EXEC izvestaj_projekti;

PROMPT Pokretanje izvestaja o istrazivacima
EXEC izvestaj_istrazivaci;

PROMPT === Testing Indexes ===

PROMPT Test 10: Upit koristeci indeks na prioritet
EXPLAIN PLAN FOR SELECT * FROM Zadaci WHERE prioritet = 'Visok';
SELECT * FROM TABLE(DBMS_XPLAN.DISPLAY()) WHERE ROWNUM <= 10;

PROMPT Test 11: Upit koristeci indeks na korisnik_id
EXPLAIN PLAN FOR SELECT * FROM ClanoviProjekta WHERE korisnik_id = (SELECT korisnik_id FROM Korisnici WHERE korisnicko_ime='istra1');
SELECT * FROM TABLE(DBMS_XPLAN.DISPLAY()) WHERE ROWNUM <= 10;

PROMPT === Complex Report Query (WITH clause, GROUP BY, HAVING, COUNT, SUM) ===

PROMPT Test 12: Kompleksan analiticki upit
WITH projekat_stats AS (
    SELECT 
        p.projekat_id,
        p.naziv_projekta,
        COUNT(DISTINCT cp.korisnik_id) as broj_clanova,
        COUNT(DISTINCT z.zadatak_id) as ukupno_zadataka,
        SUM(CASE WHEN z.progres = 100 THEN 1 ELSE 0 END) as zavrsenih_zadataka,
        COUNT(DISTINCT d.dokument_id) as broj_dokumenata
    FROM Projekti p
    LEFT JOIN ClanoviProjekta cp ON p.projekat_id = cp.projekat_id
    LEFT JOIN Zadaci z ON p.projekat_id = z.projekat_id
    LEFT JOIN Dokumenti d ON p.projekat_id = d.projekat_id
    WHERE p.status = 'Aktivan'
    GROUP BY p.projekat_id, p.naziv_projekta
    HAVING COUNT(DISTINCT z.zadatak_id) > 0
)
SELECT 
    naziv_projekta,
    broj_clanova,
    ukupno_zadataka,
    zavrsenih_zadataka,
    broj_dokumenata,
    CASE WHEN ukupno_zadataka = 0 THEN 0 ELSE ROUND((zavrsenih_zadataka / ukupno_zadataka) * 100, 2) END as procenat_zavrsenosti
FROM projekat_stats
ORDER BY procenat_zavrsenosti DESC;

PROMPT === Testing Views ===

PROMPT Pregled aktivnih projekata
SELECT * FROM v_aktivni_projekti WHERE ROWNUM <= 5;

PROMPT Pregled dokumenata sa verzijama
SELECT * FROM v_dokumenti_sa_verzijama WHERE ROWNUM <= 5;

PROMPT Pregled zadataka sa detaljima
SELECT * FROM v_zadaci_sa_detaljima WHERE ROWNUM <= 5;

PROMPT === Summary Statistics ===

PROMPT Ukupne statistike sistema
SELECT 
    (SELECT COUNT(*) FROM Korisnici) as ukupno_korisnika,
    (SELECT COUNT(*) FROM Projekti) as ukupno_projekata,
    (SELECT COUNT(*) FROM Zadaci) as ukupno_zadataka,
    (SELECT COUNT(*) FROM Dokumenti) as ukupno_dokumenata,
    (SELECT COUNT(*) FROM LogAktivnosti) as ukupno_log_unosa
FROM DUAL;

PROMPT === All tests completed! ===
