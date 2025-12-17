-- ============================================================================
-- DODAVANJE TEST PODATAKA ZA TESTIRANJE PERFORMANSI INDEKSA
-- Dodaje 5000 logova aktivnosti
-- ============================================================================

SET SERVEROUTPUT ON;

PROMPT ============================================================================
PROMPT Dodavanje 5000 test logova...
PROMPT ============================================================================

DECLARE
    v_count NUMBER := 0;
    v_korisnik_id NUMBER;
    v_tip_aktivnosti VARCHAR2(50);
    v_entitet_tip VARCHAR2(50);
    v_entitet_id NUMBER;
    v_aktivnosti DBMS_SQL.VARCHAR2_TABLE;
    v_entiteti DBMS_SQL.VARCHAR2_TABLE;
BEGIN
    -- Inicijalizuj tipove aktivnosti i entiteta
    v_aktivnosti(1) := 'LOGIN';
    v_aktivnosti(2) := 'LOGOUT';
    v_aktivnosti(3) := 'KREIRANJE';
    v_aktivnosti(4) := 'UPDATE';
    v_aktivnosti(5) := 'DELETE';
    v_aktivnosti(6) := 'PREGLED';
    v_aktivnosti(7) := 'EXPORT';
    v_aktivnosti(8) := 'IMPORT';
    
    v_entiteti(1) := 'DOKUMENT';
    v_entiteti(2) := 'PROJEKAT';
    v_entiteti(3) := 'ZADATAK';
    v_entiteti(4) := 'KORISNIK';
    
    -- Dodaj 5000 random logova
    FOR i IN 1..5000 LOOP
        v_korisnik_id := TRUNC(DBMS_RANDOM.VALUE(1, 6)); -- Korisnici od 1 do 5
        v_tip_aktivnosti := v_aktivnosti(TRUNC(DBMS_RANDOM.VALUE(1, 9)));
        v_entitet_tip := v_entiteti(TRUNC(DBMS_RANDOM.VALUE(1, 5)));
        v_entitet_id := TRUNC(DBMS_RANDOM.VALUE(1, 100));
        
        INSERT INTO SYSTEM.LogAktivnosti (
            korisnik_id,
            tip_aktivnosti,
            entitet_tip,
            entitet_id,
            naziv_entiteta,
            opis,
            rezultat,
            kreiran_datuma
        ) VALUES (
            v_korisnik_id,
            v_tip_aktivnosti,
            v_entitet_tip,
            v_entitet_id,
            'Test Entitet ' || v_entitet_id,
            'Test log za performanse indeksa',
            'SUCCESS',
            SYSTIMESTAMP - DBMS_RANDOM.VALUE(0, 365)
        );
        
        v_count := v_count + 1;
        
        -- Commit svakih 1000 redova
        IF MOD(v_count, 1000) = 0 THEN
            COMMIT;
            DBMS_OUTPUT.PUT_LINE('Ubačeno ' || v_count || ' logova...');
        END IF;
    END LOOP;
    
    COMMIT;
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('✅ Ukupno ubačeno: ' || v_count || ' logova');
    DBMS_OUTPUT.PUT_LINE('');
    
    -- Ažuriraj statistike za tabelu
    DBMS_OUTPUT.PUT_LINE('Ažuriram statistike tabele...');
    DBMS_STATS.GATHER_TABLE_STATS(
        ownname => 'SYSTEM',
        tabname => 'LOGAKTIVNOSTI',
        estimate_percent => DBMS_STATS.AUTO_SAMPLE_SIZE,
        cascade => TRUE
    );
    
    DBMS_OUTPUT.PUT_LINE('✅ Statistike ažurirane');
    
    -- Prikaži ukupan broj logova
    SELECT COUNT(*) INTO v_count FROM SYSTEM.LogAktivnosti;
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('Ukupan broj logova u tabeli: ' || v_count);
END;
/

PROMPT ============================================================================
PROMPT Test podaci uspešno dodati!
PROMPT Sada možete pokrenuti: @test_index_performance.sql
PROMPT ============================================================================
