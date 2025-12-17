-- ============================================================================
-- INDEX PERFORMANCE TEST - Kratak i Fokusiran
-- ============================================================================
SET ECHO OFF
SET SERVEROUTPUT ON SIZE UNLIMITED
SET TIMING ON
SET LINESIZE 150
SET PAGESIZE 100

PROMPT ============================================================================
PROMPT        INDEX PERFORMANCE TEST - Research Institute System
PROMPT ============================================================================
PROMPT 

-- Kreiraj privremenu tabelu za rezultate
CREATE TABLE temp_perf_results (
    test_id NUMBER,
    test_name VARCHAR2(100),
    scan_type VARCHAR2(20),
    execution_time_ms NUMBER(10,2),
    num_rows NUMBER
);

-- ============================================================================
-- TEST 1: LogAktivnosti - Pretraga po korisniku
-- ============================================================================
DECLARE
    v_start TIMESTAMP;
    v_end TIMESTAMP;
    v_count NUMBER;
    v_time_ms NUMBER;
BEGIN
    DBMS_OUTPUT.PUT_LINE('TEST 1: LogAktivnosti (korisnik_id = 2)');
    
    -- BEZ indeksa
    v_start := SYSTIMESTAMP;
    SELECT /*+ FULL(LogAktivnosti) */ COUNT(*) INTO v_count FROM LogAktivnosti WHERE korisnik_id = 2;
    v_end := SYSTIMESTAMP;
    v_time_ms := EXTRACT(SECOND FROM (v_end - v_start)) * 1000;
    INSERT INTO temp_perf_results VALUES (1, 'LogAktivnosti - korisnik_id', 'FULL SCAN', v_time_ms, v_count);
    DBMS_OUTPUT.PUT_LINE('  FULL SCAN: ' || ROUND(v_time_ms, 2) || ' ms, ' || v_count || ' redova');
    
    -- SA indeksom
    v_start := SYSTIMESTAMP;
    SELECT /*+ INDEX(LogAktivnosti idx_log_korisnik) */ COUNT(*) INTO v_count FROM LogAktivnosti WHERE korisnik_id = 2;
    v_end := SYSTIMESTAMP;
    v_time_ms := EXTRACT(SECOND FROM (v_end - v_start)) * 1000;
    INSERT INTO temp_perf_results VALUES (1, 'LogAktivnosti - korisnik_id', 'INDEX SCAN', v_time_ms, v_count);
    DBMS_OUTPUT.PUT_LINE('  INDEX SCAN: ' || ROUND(v_time_ms, 2) || ' ms, ' || v_count || ' redova');
/

-- ============================================================================
-- TEST 2: Zadaci - Pretraga po prioritetu
-- ============================================================================
DECLARE
    v_start TIMESTAMP;
    v_end TIMESTAMP;
    v_count NUMBER;
    v_time_ms NUMBER;
BEGIN
    DBMS_OUTPUT.PUT_LINE('TEST 2: Zadaci (prioritet = Visok)');
    
    -- BEZ indeksa
    v_start := SYSTIMESTAMP;
    SELECT /*+ FULL(Zadaci) */ COUNT(*) INTO v_count FROM Zadaci WHERE prioritet = 'Visok';
    v_end := SYSTIMESTAMP;
    v_time_ms := EXTRACT(SECOND FROM (v_end - v_start)) * 1000;
    INSERT INTO temp_perf_results VALUES (2, 'Zadaci - prioritet', 'FULL SCAN', v_time_ms, v_count);
    DBMS_OUTPUT.PUT_LINE('  FULL SCAN: ' || ROUND(v_time_ms, 2) || ' ms, ' || v_count || ' redova');
    
    -- SA indeksom
    v_start := SYSTIMESTAMP;
    SELECT /*+ INDEX(Zadaci idx_zadaci_prioritet) */ COUNT(*) INTO v_count FROM Zadaci WHERE prioritet = 'Visok';
    v_end := SYSTIMESTAMP;
    v_time_ms := EXTRACT(SECOND FROM (v_end - v_start)) * 1000;
    INSERT INTO temp_perf_results VALUES (2, 'Zadaci - prioritet', 'INDEX SCAN', v_time_ms, v_count);
    DBMS_OUTPUT.PUT_LINE('  INDEX SCAN: ' || ROUND(v_time_ms, 2) || ' ms, ' || v_count || ' redova');
    COMMIT;
END;
/

-- ============================================================================
-- TEST 3: Dokumenti - Pretraga po tipu
-- ============================================================================
DECLARE
    v_start TIMESTAMP;
    v_end TIMESTAMP;
    v_count NUMBER;
    v_time_ms NUMBER;
BEGIN
    DBMS_OUTPUT.PUT_LINE('TEST 3: Dokumenti (tip = PDF)');
    
    -- BEZ indeksa
    v_start := SYSTIMESTAMP;
    SELECT /*+ FULL(Dokumenti) */ COUNT(*) INTO v_count FROM Dokumenti WHERE tip_dokumenta = 'PDF';
    v_end := SYSTIMESTAMP;
    v_time_ms := EXTRACT(SECOND FROM (v_end - v_start)) * 1000;
    INSERT INTO temp_perf_results VALUES (3, 'Dokumenti - tip', 'FULL SCAN', v_time_ms, v_count);
    DBMS_OUTPUT.PUT_LINE('  FULL SCAN: ' || ROUND(v_time_ms, 2) || ' ms, ' || v_count || ' redova');
    
    -- SA indeksom
    v_start := SYSTIMESTAMP;
    SELECT /*+ INDEX(Dokumenti idx_dokumenti_tip) */ COUNT(*) INTO v_count FROM Dokumenti WHERE tip_dokumenta = 'PDF';
    v_end := SYSTIMESTAMP;
    v_time_ms := EXTRACT(SECOND FROM (v_end - v_start)) * 1000;
    INSERT INTO temp_perf_results VALUES (3, 'Dokumenti - tip', 'INDEX SCAN', v_time_ms, v_count);
    DBMS_OUTPUT.PUT_LINE('  INDEX SCAN: ' || ROUND(v_time_ms, 2) || ' ms, ' || v_count || ' redova');
    COMMIT;
END;
/

-- ============================================================================
-- TEST 4: ClanoviProjekta - JOIN
-- ============================================================================
DECLARE
    v_start TIMESTAMP;
    v_end TIMESTAMP;
    v_count NUMBER;
    v_time_ms NUMBER;
BEGIN
    DBMS_OUTPUT.PUT_LINE('TEST 4: ClanoviProjekta JOIN');
    
    -- BEZ indeksa
    v_start := SYSTIMESTAMP;
    SELECT /*+ FULL(cp) */ COUNT(*) INTO v_count FROM ClanoviProjekta cp JOIN Korisnici k ON cp.korisnik_id = k.korisnik_id;
    v_end := SYSTIMESTAMP;
    v_time_ms := EXTRACT(SECOND FROM (v_end - v_start)) * 1000;
    INSERT INTO temp_perf_results VALUES (4, 'ClanoviProjekta - JOIN', 'FULL SCAN', v_time_ms, v_count);
    DBMS_OUTPUT.PUT_LINE('  FULL SCAN: ' || ROUND(v_time_ms, 2) || ' ms, ' || v_count || ' redova');
    
    -- SA indeksom
    v_start := SYSTIMESTAMP;
    SELECT /*+ INDEX(cp idx_clanovi_projekta_korisnik) */ COUNT(*) INTO v_count FROM ClanoviProjekta cp JOIN Korisnici k ON cp.korisnik_id = k.korisnik_id;
    v_end := SYSTIMESTAMP;
    v_time_ms := EXTRACT(SECOND FROM (v_end - v_start)) * 1000;
    INSERT INTO temp_perf_results VALUES (4, 'ClanoviProjekta - JOIN', 'INDEX SCAN', v_time_ms, v_count);
    DBMS_OUTPUT.PUT_LINE('  INDEX SCAN: ' || ROUND(v_time_ms, 2) || ' ms, ' || v_count || ' redova');
    COMMIT;
END;
/

PROMPT 
PROMPT ============================================================================
PROMPT                        REZULTATI - TABELA
PROMPT ============================================================================

-- Prikaži uporednu tabelu
SELECT 
    test_name AS "Test",
    MAX(CASE WHEN scan_type = 'FULL SCAN' THEN execution_time_ms END) AS "Full Scan (ms)",
    MAX(CASE WHEN scan_type = 'INDEX SCAN' THEN execution_time_ms END) AS "Index Scan (ms)",
    ROUND((MAX(CASE WHEN scan_type = 'FULL SCAN' THEN execution_time_ms END) - 
           MAX(CASE WHEN scan_type = 'INDEX SCAN' THEN execution_time_ms END)) / 
           MAX(CASE WHEN scan_type = 'FULL SCAN' THEN execution_time_ms END) * 100, 2) AS "Ubrzanje (%)",
    MAX(num_rows) AS "Broj Redova"
FROM temp_perf_results
GROUP BY test_id, test_name
ORDER BY test_id;

PROMPT 
PROMPT ============================================================================
PROMPT                    STATISTIKA INDEKSA
PROMPT ============================================================================

SELECT 
    SUBSTR(table_name, 1, 20) AS "Tabela",
    SUBSTR(index_name, 1, 30) AS "Index",
    num_rows AS "Redova",
    CASE status WHEN 'VALID' THEN 'OK' ELSE status END AS "Status"
FROM user_indexes
WHERE table_name IN ('LOGAKTIVNOSTI', 'ZADACI', 'DOKUMENTI', 'CLANOVIPROJEKTA')
AND index_name LIKE 'IDX%'
ORDER BY table_name;

-- Očisti privremenu tabelu
DROP TABLE temp_perf_results;

PROMPT 
PROMPT ============================================================================
PROMPT                           TEST ZAVRŠEN
PROMPT ============================================================================
