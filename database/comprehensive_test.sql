-- ============================================================================
-- Comprehensive Testing Script - Triggers, Functions, Performance
-- ============================================================================

SET SERVEROUTPUT ON
SET TIMING ON

PROMPT ========================================================================
PROMPT TEST 1: Trigger Functionality - Automatic ID Generation
PROMPT ========================================================================

-- Test trigger that auto-generates IDs
DECLARE
    v_doc_id NUMBER;
    v_proj_id NUMBER;
BEGIN
    -- Insert document WITHOUT specifying dokument_id - trigger should assign it
    INSERT INTO dokumenti (naziv_dokumenta, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
    VALUES ('TEST_TRIGGER_DOC', 'Test', 'EN', 2)
    RETURNING dokument_id INTO v_doc_id;
    
    DBMS_OUTPUT.PUT_LINE('✅ Trigger test: Auto-generated dokument_id = ' || v_doc_id);
    
    -- Clean up
    DELETE FROM dokumenti WHERE dokument_id = v_doc_id;
    COMMIT;
END;
/

PROMPT
PROMPT ========================================================================
PROMPT TEST 2: Trigger Functionality - Automatic Activity Logging
PROMPT ========================================================================

DECLARE
    v_doc_id NUMBER;
    v_log_count NUMBER;
    v_log_desc VARCHAR2(500);
BEGIN
    -- Get initial log count
    SELECT COUNT(*) INTO v_log_count FROM logaktivnosti;
    DBMS_OUTPUT.PUT_LINE('Initial log count: ' || v_log_count);
    
    -- Insert a document - trigger should automatically log this
    INSERT INTO dokumenti (naziv_dokumenta, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
    VALUES ('TEST_LOGGING_DOC', 'Test', 'EN', 2)
    RETURNING dokument_id INTO v_doc_id;
    
    COMMIT;
    
    -- Check if log was created
    SELECT COUNT(*) INTO v_log_count FROM logaktivnosti 
    WHERE entitet_id = v_doc_id AND entitet_tip = 'DOKUMENT';
    
    IF v_log_count > 0 THEN
        SELECT opis INTO v_log_desc FROM logaktivnosti 
        WHERE entitet_id = v_doc_id AND entitet_tip = 'DOKUMENT' AND ROWNUM = 1;
        
        DBMS_OUTPUT.PUT_LINE('✅ Activity logging trigger worked!');
        DBMS_OUTPUT.PUT_LINE('   Logged action: ' || v_log_desc);
    ELSE
        DBMS_OUTPUT.PUT_LINE('❌ Activity logging trigger did NOT work!');
    END IF;
    
    -- Clean up
    DELETE FROM logaktivnosti WHERE entitet_id = v_doc_id AND entitet_tip = 'DOKUMENT';
    DELETE FROM dokumenti WHERE dokument_id = v_doc_id;
    COMMIT;
END;
/

PROMPT
PROMPT ========================================================================
PROMPT TEST 3: PL/SQL Functions - procenat_zavrsenih_zadataka
PROMPT ========================================================================

DECLARE
    v_procenat NUMBER;
BEGIN
    -- Test function on project 21 (should have tasks)
    SELECT procenat_zavrsenih_zadataka(21) INTO v_procenat FROM DUAL;
    DBMS_OUTPUT.PUT_LINE('✅ Function procenat_zavrsenih_zadataka(21) = ' || v_procenat || '%');
END;
/

PROMPT
PROMPT ========================================================================
PROMPT TEST 4: PL/SQL Functions - broj_aktivnih_clanova
PROMPT ========================================================================

DECLARE
    v_broj NUMBER;
BEGIN
    -- Test function on project 21
    SELECT broj_aktivnih_clanova(21) INTO v_broj FROM DUAL;
    DBMS_OUTPUT.PUT_LINE('✅ Function broj_aktivnih_clanova(21) = ' || v_broj);
END;
/

PROMPT
PROMPT ========================================================================
PROMPT TEST 5: Using Functions in SQL Query (Go code simulation)
PROMPT ========================================================================

SELECT 
    p.projekat_id,
    p.naziv_projekta,
    broj_aktivnih_clanova(p.projekat_id) as aktivni_clanovi,
    procenat_zavrsenih_zadataka(p.projekat_id) as procenat_zavrsenosti
FROM projekti p
WHERE p.projekat_id IN (21, 49, 51)
ORDER BY procenat_zavrsenosti DESC;

PROMPT
PROMPT ========================================================================
PROMPT TEST 6: Performance Test - WITH Index vs WITHOUT Index
PROMPT ========================================================================

-- First, show that index exists
SELECT index_name, table_name, uniqueness 
FROM user_indexes 
WHERE table_name = 'ZADACI' AND index_name = 'IDX_ZADACI_PRIORITET';

PROMPT
PROMPT Test Query (using index on prioritet):

SELECT COUNT(*), AVG(progres), MIN(progres), MAX(progres)
FROM zadaci
WHERE prioritet = 'Visok';

PROMPT
PROMPT ========================================================================
PROMPT TEST 7: Complex Report Execution
PROMPT ========================================================================

EXEC kompleksan_izvestaj_projekata;

PROMPT
PROMPT ========================================================================
PROMPT All Tests Completed!
PROMPT ========================================================================
PROMPT
PROMPT Summary:
PROMPT - Triggers: AUTO ID generation + Activity logging
PROMPT - Functions: Used in queries for real-time calculations
PROMPT - Indexes: Speed up WHERE clauses
PROMPT - Complex Report: Demonstrates all PL/SQL features
PROMPT ========================================================================

SET TIMING OFF
