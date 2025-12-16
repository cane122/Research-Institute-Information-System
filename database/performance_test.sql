-- ============================================================================
-- Performance Testing Script - Comparing Queries With/Without Indexes
-- ============================================================================

SET TIMING ON
SET SERVEROUTPUT ON

-- Disable output for cleaner timing results
SET SERVEROUTPUT OFF

PROMPT ========================================
PROMPT PERFORMANCE TEST: Searching Tasks by Priority
PROMPT ========================================
PROMPT

-- Test 1: Query using idx_zadaci_prioritet index
PROMPT Test 1: Querying tasks with HIGH priority (using index)
SELECT COUNT(*), AVG(progres), SUM(progres)
FROM zadaci
WHERE prioritet = 'Visok';

PROMPT
PROMPT Test 2: Querying tasks and joining with projects
SELECT z.zadatak_id, z.naziv_zadatka, p.naziv_projekta, z.prioritet, z.progres
FROM zadaci z
JOIN projekti p ON z.projekat_id = p.projekat_id
WHERE z.prioritet = 'Visok'
ORDER BY z.zadatak_id DESC;

PROMPT
PROMPT ========================================
PROMPT PERFORMANCE TEST: Searching Project Members
PROMPT ========================================
PROMPT

-- Test 3: Query using project members
PROMPT Test 3: Finding all projects and counting members
SELECT p.projekat_id, p.naziv_projekta, COUNT(*) as broj_dokumenata
FROM projekti p
LEFT JOIN dokumenti d ON p.projekat_id = d.projekat_id
GROUP BY p.projekat_id, p.naziv_projekta
HAVING COUNT(*) > 0
ORDER BY broj_dokumenata DESC;

PROMPT
PROMPT ========================================
PROMPT PERFORMANCE TEST: Complex Query with Functions
PROMPT ========================================
PROMPT

-- Test 4: Using PL/SQL functions in query
PROMPT Test 4: Projects with calculated statistics (using PL/SQL functions)
SELECT 
    p.projekat_id,
    p.naziv_projekta,
    broj_aktivnih_clanova(p.projekat_id) as aktivni_clanovi,
    procenat_zavrsenih_zadataka(p.projekat_id) as procenat_zavrsenosti
FROM projekti p
WHERE p.status = 'Aktivan'
ORDER BY procenat_zavrsenosti DESC;

PROMPT
PROMPT ========================================
PROMPT PERFORMANCE TEST: Document Search with Joins
PROMPT ========================================
PROMPT

-- Test 5: Document query joining multiple tables
PROMPT Test 5: Documents with project and user info (multiple joins)
SELECT 
    d.dokument_id,
    d.naziv_dokumenta,
    p.naziv_projekta,
    k.korisnicko_ime,
    COUNT(dt.tag_id) as broj_tagova
FROM dokumenti d
LEFT JOIN projekti p ON d.projekat_id = p.projekat_id
JOIN korisnici k ON d.kreirao_korisnik_id = k.korisnik_id
LEFT JOIN dokumenttagovi dt ON d.dokument_id = dt.dokument_id
GROUP BY d.dokument_id, d.naziv_dokumenta, p.naziv_projekta, k.korisnicko_ime
HAVING COUNT(dt.tag_id) > 0
ORDER BY broj_tagova DESC;

PROMPT
PROMPT ========================================
PROMPT Performance tests completed
PROMPT Check timing results above
PROMPT ========================================

SET TIMING OFF
SET SERVEROUTPUT ON
