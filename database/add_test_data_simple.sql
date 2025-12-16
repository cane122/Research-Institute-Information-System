-- ============================================================================
-- SIMPLIFIED MASSIVE DATA INSERT - Works around trigger constraints
-- ============================================================================

-- Temporarily disable triggers for bulk insert
ALTER TRIGGER trg_log_projekta DISABLE;
ALTER TRIGGER trg_log_dokumenta DISABLE;

-- Dodavanje 100 projekata
BEGIN
    FOR i IN 1..100 LOOP
        INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, status, rukovodilac_id)
        VALUES (
            'Performance Test Project ' || i,
            'Projekat za testiranje performansi broj ' || i,
            SYSDATE - DBMS_RANDOM.VALUE(1, 365),
            CASE WHEN MOD(i, 4) = 0 THEN 'Zavrsen' ELSE 'Aktivan' END,
            2  -- Valid user ID
        );
    END LOOP;
    COMMIT;
    DBMS_OUTPUT.PUT_LINE('✅ Dodato 100 projekata');
END;
/

-- Dodavanje 500 zadataka  
DECLARE
    v_count NUMBER := 0;
BEGIN
    FOR proj IN (SELECT projekat_id FROM projekti) LOOP
        FOR i IN 1..5 LOOP
            INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, progres, faza_id)
            VALUES (
                proj.projekat_id,
                'Task ' || i || ' for project ' || proj.projekat_id,
                'Description for task ' || i,
                CASE MOD(i, 3) WHEN 0 THEN 'Visok' WHEN 1 THEN 'Srednji' ELSE 'Nizak' END,
                DBMS_RANDOM.VALUE(0, 100),
                1
            );
            v_count := v_count + 1;
        END LOOP;
    END LOOP;
    COMMIT;
    DBMS_OUTPUT.PUT_LINE('✅ Dodato ' || v_count || ' zadataka');
END;
/

-- Dodavanje 600 dokumenata
DECLARE
    v_count NUMBER := 0;
BEGIN
    FOR proj IN (SELECT projekat_id FROM projekti) LOOP
        FOR i IN 1..5 LOOP
            INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
            VALUES (
                proj.projekat_id,
                'Document ' || i || ' - Project ' || proj.projekat_id,
                'Document description ' || i,
                CASE MOD(i, 4) WHEN 0 THEN 'Technical' WHEN 1 THEN 'Report' WHEN 2 THEN 'Plan' ELSE 'Requirements' END,
                CASE MOD(i, 2) WHEN 0 THEN 'EN' ELSE 'SR' END,
                2
            );
            v_count := v_count + 1;
        END LOOP;
    END LOOP;
    COMMIT;
    DBMS_OUTPUT.PUT_LINE('✅ Dodato ' || v_count || ' dokumenata');
END;
/

-- Re-enable triggers
ALTER TRIGGER trg_log_projekta ENABLE;
ALTER TRIGGER trg_log_dokumenta ENABLE;

-- Final statistics
SET SERVEROUTPUT ON
DECLARE
    v_count NUMBER;
BEGIN
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('========================================');
    DBMS_OUTPUT.PUT_LINE('FINALNA STATISTIKA PODATAKA:');
    DBMS_OUTPUT.PUT_LINE('========================================');
    
    SELECT COUNT(*) INTO v_count FROM projekti;
    DBMS_OUTPUT.PUT_LINE('📁 Ukupno projekata: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM zadaci;
    DBMS_OUTPUT.PUT_LINE('📋 Ukupno zadataka: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM dokumenti;
    DBMS_OUTPUT.PUT_LINE('📄 Ukupno dokumenata: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM ClanoviProjekta;
    DBMS_OUTPUT.PUT_LINE('👥 Ukupno članova: ' || v_count);
    
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('✅ Baza spremna za performance testiranje!');
    DBMS_OUTPUT.PUT_LINE('========================================');
END;
/
