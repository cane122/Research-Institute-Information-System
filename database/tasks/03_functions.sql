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
