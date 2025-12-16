-- ============================================================================
-- Oracle PL/SQL Procedures for Complex Reports
-- Demonstrates: Custom Types, Cursors, Complex SQL with WITH clause,
--               JOIN (3+ tables), GROUP BY, HAVING, COUNT, SUM
-- ============================================================================

SET SERVEROUTPUT ON;

-- ============================================================================
-- KOMPLEKSAN IZVEŠTAJ SA SVIM ZAHTEVIMA
-- ============================================================================

CREATE OR REPLACE PROCEDURE kompleksan_izvestaj_projekata IS
    -- ========================================
    -- SLOŽENI PL/SQL TIPOVI (Custom Types)
    -- ========================================
    
    -- Tip 1: RECORD za statistiku projekta
    TYPE t_projekat_statistika IS RECORD (
        projekat_id NUMBER,
        naziv_projekta VARCHAR2(255),
        broj_clanova NUMBER,
        ukupno_zadataka NUMBER,
        zavrsenih_zadataka NUMBER,
        suma_progresa NUMBER,
        broj_dokumenata NUMBER,
        procenat_zavrsenosti NUMBER(5,2)
    );
    
    -- Tip 2: TABLE OF za kolekciju projekata
    TYPE t_projekti_tabela IS TABLE OF t_projekat_statistika INDEX BY PLS_INTEGER;
    
    -- Promenljive
    v_projekti t_projekti_tabela;
    v_index PLS_INTEGER := 0;
    v_ukupno_projekata NUMBER := 0;
    v_ukupno_zadataka NUMBER := 0;
    
    -- ========================================
    -- KURSOR SA SLOŽENIM SQL UPITOM
    -- Koristi: WITH klauzulu, JOIN 4 tabele,
    --          GROUP BY, HAVING, COUNT, SUM
    -- ========================================
    CURSOR c_projekat_stats IS
        WITH projekat_agregati AS (
            -- CTE (Common Table Expression) za agregaciju podataka
            SELECT 
                p.projekat_id,
                p.naziv_projekta,
                p.status,
                COUNT(DISTINCT cp.korisnik_id) as broj_clanova,
                COUNT(DISTINCT z.zadatak_id) as ukupno_zadataka,
                SUM(CASE WHEN z.progres = 100 THEN 1 ELSE 0 END) as zavrsenih_zadataka,
                SUM(NVL(z.progres, 0)) as suma_progresa,
                COUNT(DISTINCT d.dokument_id) as broj_dokumenata
            FROM Projekti p
            LEFT JOIN ClanoviProjekta cp ON p.projekat_id = cp.projekat_id
            LEFT JOIN Zadaci z ON p.projekat_id = z.projekat_id
            LEFT JOIN Dokumenti d ON p.projekat_id = d.projekat_id
            WHERE p.status = 'Aktivan'
            GROUP BY p.projekat_id, p.naziv_projekta, p.status
            -- HAVING: Pokazujemo samo projekte sa zadacima (ali možete i bez toga)
            -- Zakomentarisano za test: HAVING COUNT(DISTINCT z.zadatak_id) > 0
        )
        SELECT 
            projekat_id,
            naziv_projekta,
            broj_clanova,
            ukupno_zadataka,
            zavrsenih_zadataka,
            suma_progresa,
            broj_dokumenata,
            CASE 
                WHEN ukupno_zadataka = 0 THEN 0 
                ELSE ROUND((zavrsenih_zadataka / ukupno_zadataka) * 100, 2) 
            END as procenat_zavrsenosti
        FROM projekat_agregati
        ORDER BY procenat_zavrsenosti DESC, naziv_projekta;
    
    r_projekat t_projekat_statistika;
    
BEGIN
    DBMS_OUTPUT.PUT_LINE('==========================================================');
    DBMS_OUTPUT.PUT_LINE('    KOMPLEKSAN IZVEŠTAJ O PROJEKTIMA I STATISTIKAMA');
    DBMS_OUTPUT.PUT_LINE('==========================================================');
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('Napomena: Ako nema podataka, učitajte seed data:');
    DBMS_OUTPUT.PUT_LINE('  @06_seed_data.sql');
    DBMS_OUTPUT.PUT_LINE('');
    
    -- Otvaranje kursora i prikupljanje podataka u kolekciju
    OPEN c_projekat_stats;
    LOOP
        FETCH c_projekat_stats INTO r_projekat;
        EXIT WHEN c_projekat_stats%NOTFOUND;
        
        v_index := v_index + 1;
        v_projekti(v_index) := r_projekat;
        
        -- Akumulacija za ukupnu statistiku
        v_ukupno_projekata := v_ukupno_projekata + 1;
        v_ukupno_zadataka := v_ukupno_zadataka + r_projekat.ukupno_zadataka;
    END LOOP;
    CLOSE c_projekat_stats;
    
    -- Prikaz podataka iz kolekcije
    IF v_projekti.COUNT > 0 THEN
        DBMS_OUTPUT.PUT_LINE('Pronađeno projekata: ' || v_projekti.COUNT);
        DBMS_OUTPUT.PUT_LINE('----------------------------------------------------------');
        
        FOR i IN 1..v_projekti.COUNT LOOP
            DBMS_OUTPUT.PUT_LINE('');
            DBMS_OUTPUT.PUT_LINE('PROJEKAT: ' || v_projekti(i).naziv_projekta);
            DBMS_OUTPUT.PUT_LINE('  - ID: ' || v_projekti(i).projekat_id);
            DBMS_OUTPUT.PUT_LINE('  - Članovi tima: ' || v_projekti(i).broj_clanova);
            DBMS_OUTPUT.PUT_LINE('  - Ukupno zadataka: ' || v_projekti(i).ukupno_zadataka);
            DBMS_OUTPUT.PUT_LINE('  - Završeno zadataka: ' || v_projekti(i).zavrsenih_zadataka);
            DBMS_OUTPUT.PUT_LINE('  - Suma progresa: ' || v_projekti(i).suma_progresa || '%');
            DBMS_OUTPUT.PUT_LINE('  - Broj dokumenata: ' || v_projekti(i).broj_dokumenata);
            DBMS_OUTPUT.PUT_LINE('  - Procenat završenosti: ' || v_projekti(i).procenat_zavrsenosti || '%');
            DBMS_OUTPUT.PUT_LINE('----------------------------------------------------------');
        END LOOP;
        
        -- Agregatna statistika
        DBMS_OUTPUT.PUT_LINE('');
        DBMS_OUTPUT.PUT_LINE('==========================================================');
        DBMS_OUTPUT.PUT_LINE('UKUPNA STATISTIKA:');
        DBMS_OUTPUT.PUT_LINE('  - Ukupno aktivnih projekata: ' || v_ukupno_projekata);
        DBMS_OUTPUT.PUT_LINE('  - Ukupno zadataka: ' || v_ukupno_zadataka);
        DBMS_OUTPUT.PUT_LINE('  - Prosečno zadataka po projektu: ' || 
            ROUND(v_ukupno_zadataka / v_ukupno_projekata, 2));
        DBMS_OUTPUT.PUT_LINE('==========================================================');
    ELSE
        DBMS_OUTPUT.PUT_LINE('Nema podataka za prikaz.');
    END IF;
    
EXCEPTION
    WHEN OTHERS THEN
        IF c_projekat_stats%ISOPEN THEN
            CLOSE c_projekat_stats;
        END IF;
        DBMS_OUTPUT.PUT_LINE('GREŠKA: ' || SQLERRM);
        RAISE;
END kompleksan_izvestaj_projekata;
/