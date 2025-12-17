-- ============================================================================
-- PL/SQL Functions for Document Management System
-- ============================================================================

-- Function to calculate total storage used by user's documents (in MB)
CREATE OR REPLACE FUNCTION ukupna_velicina_korisnika(p_korisnik_id IN NUMBER)
RETURN NUMBER
IS
    v_velicina NUMBER := 0;
BEGIN
    SELECT COALESCE(SUM(vd.velicina_fajla_mb), 0) INTO v_velicina
    FROM VerzijeDokumenata vd
    JOIN Dokumenti d ON vd.dokument_id = d.dokument_id
    WHERE d.kreirao_korisnik_id = p_korisnik_id;
    
    RETURN ROUND(v_velicina, 2);
END ukupna_velicina_korisnika;
/

