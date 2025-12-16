-- ============================================================================
-- MASOVNI UNOS PODATAKA - Za demonstraciju performance poboljšanja
-- ============================================================================

-- Dodavanje dodatnih projekata (50 novih)
BEGIN
    FOR i IN 1..50 LOOP
        INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, status, rukovodilac_id)
        VALUES (
            'Projekat Test ' || i,
            'Test projekat broj ' || i || ' za demonstraciju performance optimizacija',
            SYSDATE - DBMS_RANDOM.VALUE(1, 365),
            CASE WHEN MOD(i, 3) = 0 THEN 'Zavrsen' ELSE 'Aktivan' END,
            2 + MOD(i, 10)
        );
    END LOOP;
    COMMIT;
END;
/

-- Dodavanje zadataka za svaki projekat (200 novih zadataka)
DECLARE
    v_proj_id NUMBER;
BEGIN
    FOR proj IN (SELECT projekat_id FROM projekti WHERE projekat_id > 50) LOOP
        FOR i IN 1..4 LOOP
            INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, progres, faza_id)
            VALUES (
                proj.projekat_id,
                'Zadatak ' || i || ' za projekat ' || proj.projekat_id,
                'Detaljan opis zadatka broj ' || i,
                CASE MOD(i, 3) WHEN 0 THEN 'Visok' WHEN 1 THEN 'Srednji' ELSE 'Nizak' END,
                DBMS_RANDOM.VALUE(0, 100),
                1
            );
        END LOOP;
    END LOOP;
    COMMIT;
END;
/

-- Dodavanje dokumenata za svaki projekat (300 novih dokumenata)
DECLARE
    v_types DBMS_SQL.VARCHAR2_TABLE;
    v_langs DBMS_SQL.VARCHAR2_TABLE;
BEGIN
    v_types(1) := 'Technical'; v_types(2) := 'Report'; v_types(3) := 'Plan';
    v_types(4) := 'Requirements'; v_types(5) := 'Specification';
    v_langs(1) := 'EN'; v_langs(2) := 'SR';
    
    FOR proj IN (SELECT projekat_id FROM projekti WHERE projekat_id > 20) LOOP
        FOR i IN 1..6 LOOP
            INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
            VALUES (
                proj.projekat_id,
                'Dokument ' || i || ' - Projekat ' || proj.projekat_id,
                'Opis dokumenta broj ' || i || ' za projekat ' || proj.projekat_id,
                v_types(MOD(i, 5) + 1),
                v_langs(MOD(i, 2) + 1),
                2 + MOD(i, 10)
            );
        END LOOP;
    END LOOP;
    COMMIT;
END;
/

-- Dodavanje tagova za dokumente (1000+ tag veza)
DECLARE
    v_tag_id NUMBER;
    v_tags DBMS_SQL.VARCHAR2_TABLE;
BEGIN
    v_tags(1) := 'urgent'; v_tags(2) := 'review'; v_tags(3) := 'draft';
    v_tags(4) := 'final'; v_tags(5) := 'technical'; v_tags(6) := 'business';
    v_tags(7) := 'security'; v_tags(8) := 'performance'; v_tags(9) := 'testing';
    v_tags(10) := 'documentation';
    
    -- Kreiraj tagove ako ne postoje
    FOR i IN 1..10 LOOP
        BEGIN
            SELECT tag_id INTO v_tag_id FROM tagovi WHERE naziv_taga = v_tags(i);
        EXCEPTION
            WHEN NO_DATA_FOUND THEN
                INSERT INTO tagovi (naziv_taga) VALUES (v_tags(i)) RETURNING tag_id INTO v_tag_id;
        END;
    END LOOP;
    COMMIT;
    
    -- Dodaj tagove dokumentima
    FOR doc IN (SELECT dokument_id FROM dokumenti WHERE dokument_id > 70) LOOP
        FOR i IN 1..3 LOOP
            BEGIN
                SELECT tag_id INTO v_tag_id FROM tagovi WHERE naziv_taga = v_tags(MOD(doc.dokument_id * i, 10) + 1);
                INSERT INTO dokumenttagovi (dokument_id, tag_id) VALUES (doc.dokument_id, v_tag_id);
            EXCEPTION
                WHEN DUP_VAL_ON_INDEX THEN NULL;
            END;
        END LOOP;
    END LOOP;
    COMMIT;
END;
/

-- Dodavanje članova projekata (500+ veza)
BEGIN
    FOR proj IN (SELECT projekat_id FROM projekti WHERE projekat_id > 20) LOOP
        FOR user_id IN 2..5 LOOP
            BEGIN
                INSERT INTO ClanoviProjekta (projekat_id, korisnik_id)
                VALUES (proj.projekat_id, user_id);
            EXCEPTION
                WHEN DUP_VAL_ON_INDEX THEN NULL;
            END;
        END LOOP;
    END LOOP;
    COMMIT;
END;
/

-- Prikaz statistike nakon unosa
SET SERVEROUTPUT ON
DECLARE
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count FROM projekti;
    DBMS_OUTPUT.PUT_LINE('Ukupno projekata: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM zadaci;
    DBMS_OUTPUT.PUT_LINE('Ukupno zadataka: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM dokumenti;
    DBMS_OUTPUT.PUT_LINE('Ukupno dokumenata: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM dokumenttagovi;
    DBMS_OUTPUT.PUT_LINE('Ukupno tag veza: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM ClanoviProjekta;
    DBMS_OUTPUT.PUT_LINE('Ukupno članova projekata: ' || v_count);
    
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('✅ Masovni unos podataka završen!');
    DBMS_OUTPUT.PUT_LINE('Sada možete testirati performance sa većim skupom podataka.');
END;
/
