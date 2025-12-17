-- Oracle PL/SQL implementations

CREATE OR REPLACE FUNCTION procenat_zavrsenih_zadataka(p_projekat_id IN NUMBER)
RETURN NUMBER
IS
    v_ukupno NUMBER := 0;
    v_zavrseno NUMBER := 0;
BEGIN
    SELECT COUNT(*) INTO v_ukupno FROM Zadaci WHERE projekat_id = p_projekat_id;
    IF v_ukupno = 0 THEN
        RETURN 0;
    END IF;

    SELECT COUNT(*) INTO v_zavrseno FROM Zadaci WHERE projekat_id = p_projekat_id AND progres = 100;

    RETURN (v_zavrseno / v_ukupno) * 100;
END procenat_zavrsenih_zadataka;
/

CREATE OR REPLACE FUNCTION broj_aktivnih_clanova(p_projekat_id IN NUMBER)
RETURN NUMBER
IS
    v_broj NUMBER := 0;
BEGIN
    SELECT COUNT(*) INTO v_broj
    FROM ClanoviProjekta cp
    JOIN Korisnici k ON cp.korisnik_id = k.korisnik_id
    WHERE cp.projekat_id = p_projekat_id AND k.status = 'aktivan';

    RETURN v_broj;
END broj_aktivnih_clanova;
/

-- ============================================================================
-- FUNKCIJE ZA DOKUMENTE (Document Management System)
-- ============================================================================

-- Function to count documents for a project
CREATE OR REPLACE FUNCTION broj_dokumenata_projekta(p_projekat_id IN NUMBER)
RETURN NUMBER
IS
    v_broj NUMBER := 0;
BEGIN
    SELECT COUNT(*) INTO v_broj
    FROM Dokumenti
    WHERE projekat_id = p_projekat_id;
    
    RETURN v_broj;
END broj_dokumenata_projekta;
/

-- Function to calculate average document size for a project
CREATE OR REPLACE FUNCTION prosecna_velicina_dokumenata(p_projekat_id IN NUMBER)
RETURN NUMBER
IS
    v_prosek NUMBER := 0;
BEGIN
    SELECT COALESCE(AVG(vd.velicina_fajla_mb), 0) INTO v_prosek
    FROM Dokumenti d
    JOIN VerzijeDokumenata vd ON d.dokument_id = vd.dokument_id
    WHERE d.projekat_id = p_projekat_id;
    
    RETURN ROUND(v_prosek, 2);
END prosecna_velicina_dokumenata;
/

-- Function to count total document versions for a user
CREATE OR REPLACE FUNCTION broj_verzija_korisnika(p_korisnik_id IN NUMBER)
RETURN NUMBER
IS
    v_broj NUMBER := 0;
BEGIN
    SELECT COUNT(*) INTO v_broj
    FROM VerzijeDokumenata vd
    WHERE vd.postavio_korisnik_id = p_korisnik_id;
    
    RETURN v_broj;
END broj_verzija_korisnika;
/

-- Function to calculate total storage used by user's documents (in MB)
CREATE OR REPLACE FUNCTION ukupna_velicina_korisnika(p_korisnik_id IN NUMBER)
RETURN NUMBER
IS
    v_velicina NUMBER := 0;
BEGIN
    SELECT COALESCE(SUM(vd.velicina_fajla_mb), 0) INTO v_velicina
    FROM VerzijeDokumenata vd
    WHERE vd.postavio_korisnik_id = p_korisnik_id;
    
    RETURN ROUND(v_velicina, 2);
END ukupna_velicina_korisnika;
/
