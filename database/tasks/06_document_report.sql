-- ============================================================================
-- DOCUMENT MANAGEMENT COMPLEX REPORT PROCEDURE
-- Demonstrates: Custom Types, Cursors, Complex SQL with WITH clause,
--               JOIN (4 tables), GROUP BY, HAVING, COUNT, SUM, AVG
-- ============================================================================

SET SERVEROUTPUT ON;

CREATE OR REPLACE PROCEDURE izvestaj_dokumenata_korisnika IS
    -- ========================================
    -- SLOŽENI PL/SQL TIPOVI (Custom Types)
    -- ========================================
    
    -- Tip: RECORD za statistiku korisnika i dokumenata
    TYPE t_korisnik_doc_stat IS RECORD (
        korisnik_id NUMBER,
        korisnicko_ime VARCHAR2(100),
        ime_prezime VARCHAR2(200),
        broj_kreiranih_dokumenata NUMBER,
        broj_postavljenih_verzija NUMBER,
        ukupna_velicina_mb NUMBER,
        prosecna_velicina_mb NUMBER,
        broj_projekata NUMBER
    );
    
    -- Tip: TABLE OF za kolekciju korisnika
    TYPE t_korisnici_tabela IS TABLE OF t_korisnik_doc_stat INDEX BY PLS_INTEGER;
    
    -- Promenljive
    v_korisnici t_korisnici_tabela;
    v_index PLS_INTEGER := 0;
    v_ukupno_dokumenata NUMBER := 0;
    v_ukupno_mb NUMBER := 0;
    
    -- ========================================
    -- KURSOR SA SLOŽENIM SQL UPITOM
    -- Koristi: WITH klauzulu, JOIN 4 tabele,
    --          GROUP BY, COUNT, SUM, AVG
    -- ========================================
    CURSOR c_korisnik_stats IS
        WITH korisnik_agregati AS (
            -- CTE za agregaciju podataka o korisnicima i dokumentima
            SELECT 
                k.korisnik_id,
                k.korisnicko_ime,
                NVL(k.ime || ' ' || k.prezime, k.korisnicko_ime) as ime_prezime,
                COUNT(DISTINCT d.dokument_id) as broj_kreiranih_dokumenata,
                COUNT(DISTINCT vd.verzija_id) as broj_postavljenih_verzija,
                SUM(NVL(vd.velicina_fajla_mb, 0)) as ukupna_velicina_mb,
                COUNT(DISTINCT d.projekat_id) as broj_projekata
            FROM Korisnici k
            LEFT JOIN Dokumenti d ON k.korisnik_id = d.kreirao_korisnik_id
            LEFT JOIN VerzijeDokumenata vd ON d.dokument_id = vd.dokument_id
            LEFT JOIN Projekti p ON d.projekat_id = p.projekat_id
            WHERE k.status = 'aktivan'
            GROUP BY k.korisnik_id, k.korisnicko_ime, k.ime, k.prezime
            HAVING COUNT(DISTINCT d.dokument_id) > 0  -- Samo korisnici koji su kreirali dokumente
        )
        SELECT 
            korisnik_id,
            korisnicko_ime,
            ime_prezime,
            broj_kreiranih_dokumenata,
            broj_postavljenih_verzija,
            ukupna_velicina_mb,
            CASE 
                WHEN broj_kreiranih_dokumenata = 0 THEN 0 
                ELSE ROUND(ukupna_velicina_mb / broj_kreiranih_dokumenata, 2) 
            END as prosecna_velicina_mb,
            broj_projekata
        FROM korisnik_agregati
        ORDER BY broj_kreiranih_dokumenata DESC, ukupna_velicina_mb DESC;
    
    r_korisnik t_korisnik_doc_stat;

BEGIN
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('KOMPLEKSAN IZVEŠTAJ O KORISNICIMA I DOKUMENTIMA');
    DBMS_OUTPUT.PUT_LINE('========================================================');
    
    -- Otvaranje kursora i prolazak kroz rezultate
    OPEN c_korisnik_stats;
    LOOP
        FETCH c_korisnik_stats INTO r_korisnik;
        EXIT WHEN c_korisnik_stats%NOTFOUND;
        
        v_index := v_index + 1;
        v_korisnici(v_index) := r_korisnik;
        
        v_ukupno_dokumenata := v_ukupno_dokumenata + r_korisnik.broj_kreiranih_dokumenata;
        v_ukupno_mb := v_ukupno_mb + r_korisnik.ukupna_velicina_mb;
    END LOOP;
    CLOSE c_korisnik_stats;
    
    DBMS_OUTPUT.PUT_LINE('Pronađeno korisnika sa dokumentima: ' || v_index);
    DBMS_OUTPUT.PUT_LINE('');
    
    -- Ispis rezultata iz kolekcije
    FOR i IN 1..v_index LOOP
        DBMS_OUTPUT.PUT_LINE('KORISNIK: ' || v_korisnici(i).ime_prezime);
        DBMS_OUTPUT.PUT_LINE('  - Username: ' || v_korisnici(i).korisnicko_ime);
        DBMS_OUTPUT.PUT_LINE('  - Kreiranih dokumenata: ' || v_korisnici(i).broj_kreiranih_dokumenata);
        DBMS_OUTPUT.PUT_LINE('  - Postavljenih verzija: ' || v_korisnici(i).broj_postavljenih_verzija);
        DBMS_OUTPUT.PUT_LINE('  - Ukupna veličina: ' || v_korisnici(i).ukupna_velicina_mb || ' MB');
        DBMS_OUTPUT.PUT_LINE('  - Prosečna veličina: ' || v_korisnici(i).prosecna_velicina_mb || ' MB');
        DBMS_OUTPUT.PUT_LINE('  - Broj projekata: ' || v_korisnici(i).broj_projekata);
        DBMS_OUTPUT.PUT_LINE('----------------------------------------------------------');
    END LOOP;
    
    DBMS_OUTPUT.PUT_LINE('========================================================');
    DBMS_OUTPUT.PUT_LINE('UKUPNA STATISTIKA:');
    DBMS_OUTPUT.PUT_LINE('  - Ukupno aktivnih korisnika sa dokumentima: ' || v_index);
    DBMS_OUTPUT.PUT_LINE('  - Ukupno dokumenata: ' || v_ukupno_dokumenata);
    DBMS_OUTPUT.PUT_LINE('  - Ukupna zauzetost: ' || ROUND(v_ukupno_mb, 2) || ' MB');
    IF v_index > 0 THEN
        DBMS_OUTPUT.PUT_LINE('  - Prosečno dokumenata po korisniku: ' || 
            ROUND(v_ukupno_dokumenata / v_index, 2));
        DBMS_OUTPUT.PUT_LINE('  - Prosečna zauzetost po korisniku: ' || 
            ROUND(v_ukupno_mb / v_index, 2) || ' MB');
    END IF;
    DBMS_OUTPUT.PUT_LINE('========================================================');
    
EXCEPTION
    WHEN OTHERS THEN
        IF c_korisnik_stats%ISOPEN THEN
            CLOSE c_korisnik_stats;
        END IF;
        DBMS_OUTPUT.PUT_LINE('GREŠKA: ' || SQLERRM);
END izvestaj_dokumenata_korisnika;
/
