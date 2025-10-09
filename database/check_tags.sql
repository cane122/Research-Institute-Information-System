-- Check tags data to verify proper saving
-- This script will help verify if tags are being saved correctly

-- 1. Check all tags in the system
SELECT 
    tag_id,
    naziv_taga,
    datum_kreiranja
FROM tagovi
ORDER BY datum_kreiranja DESC;

-- 2. Check document-tag relationships
SELECT 
    dt.dokument_id,
    d.naziv_dokumenta,
    dt.tag_id,
    t.naziv_taga,
    dt.datum_dodavanja
FROM dokumenttagovi dt
JOIN dokumenti d ON dt.dokument_id = d.dokument_id
JOIN tagovi t ON dt.tag_id = t.tag_id
ORDER BY dt.datum_dodavanja DESC;

-- 3. Count tags per document
SELECT 
    d.dokument_id,
    d.naziv_dokumenta,
    COUNT(dt.tag_id) as broj_tagova
FROM dokumenti d
LEFT JOIN dokumenttagovi dt ON d.dokument_id = dt.dokument_id
GROUP BY d.dokument_id, d.naziv_dokumenta
ORDER BY d.datum_kreiranja DESC;

-- 4. Check if there are any orphaned tags (tags not associated with any document)
SELECT 
    t.tag_id,
    t.naziv_taga,
    COUNT(dt.dokument_id) as broj_dokumenata
FROM tagovi t
LEFT JOIN dokumenttagovi dt ON t.tag_id = dt.tag_id
GROUP BY t.tag_id, t.naziv_taga
HAVING COUNT(dt.dokument_id) = 0;

-- 5. Show recent documents with their tags
SELECT 
    d.dokument_id,
    d.naziv_dokumenta,
    d.datum_kreiranja,
    STRING_AGG(t.naziv_taga, ', ') as tagovi
FROM dokumenti d
LEFT JOIN dokumenttagovi dt ON d.dokument_id = dt.dokument_id
LEFT JOIN tagovi t ON dt.tag_id = t.tag_id
GROUP BY d.dokument_id, d.naziv_dokumenta, d.datum_kreiranja
ORDER BY d.datum_kreiranja DESC
LIMIT 10;
