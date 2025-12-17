-- ============================================================================
-- DODAVANJE DODATNIH 20000 TEST PODATAKA
-- ============================================================================

SET SERVEROUTPUT ON;

PROMPT ============================================================================
PROMPT Dodavanje 20000 dodatnih test logova...
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
    
    FOR i IN 1..20000 LOOP
        v_korisnik_id := TRUNC(DBMS_RANDOM.VALUE(1, 6));
        v_tip_aktivnosti := v_aktivnosti(TRUNC(DBMS_RANDOM.VALUE(1, 9)));
        v_entitet_tip := v_entiteti(TRUNC(DBMS_RANDOM.VALUE(1, 5)));
        v_entitet_id := TRUNC(DBMS_RANDOM.VALUE(1, 1000));
        
        INSERT INTO SYSTEM.LogAktivnosti (
            korisnik_id,
            tip_aktivnosti,
            entitet_tip,
            entitet_id,
            naziv_entiteta,
            opis,
            ip_adresa,
            user_agent,
            rezultat,
            dodatne_informacije,
            kreiran_datuma
        ) VALUES (
            v_korisnik_id,
            v_tip_aktivnosti,
            v_entitet_tip,
            v_entitet_id,
            'Test Entitet ' || v_entitet_id,
            'Detaljni opis test loga broj ' || i || ' za testiranje performansi indeksa sa većim brojem podataka',
            '192.168.' || TRUNC(DBMS_RANDOM.VALUE(1, 255)) || '.' || TRUNC(DBMS_RANDOM.VALUE(1, 255)),
            'Mozilla/5.0 Test User Agent String For Performance Testing',
            'SUCCESS',
            'Extra informacije za testiranje: Dodatni podaci, meta informacije, tracking info',
            SYSTIMESTAMP - DBMS_RANDOM.VALUE(0, 730)
        );
        
        v_count := v_count + 1;
        
        IF MOD(v_count, 5000) = 0 THEN
            COMMIT;
            DBMS_OUTPUT.PUT_LINE('Ubačeno ' || v_count || ' logova...');
        END IF;
    END LOOP;
    
    COMMIT;
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('✅ Ukupno ubačeno: ' || v_count || ' dodatnih logova');
    
    DBMS_OUTPUT.PUT_LINE('Ažuriram statistike...');
    DBMS_STATS.GATHER_TABLE_STATS(
        ownname => 'SYSTEM',
        tabname => 'LOGAKTIVNOSTI',
        estimate_percent => DBMS_STATS.AUTO_SAMPLE_SIZE,
        cascade => TRUE
    );
    
    SELECT COUNT(*) INTO v_count FROM SYSTEM.LogAktivnosti;
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('Ukupan broj logova u tabeli: ' || v_count);
END;
/
