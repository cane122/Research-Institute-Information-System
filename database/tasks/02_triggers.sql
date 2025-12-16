/* Oracle PL/SQL triggers for SQL Developer / SQL*Plus */

-- ============================================================================
-- BEFORE INSERT Triggers - Auto-populate Primary Keys from Sequences
-- ============================================================================

CREATE OR REPLACE TRIGGER trg_projekti_bi
BEFORE INSERT ON Projekti
FOR EACH ROW
BEGIN
    IF :NEW.projekat_id IS NULL THEN
        :NEW.projekat_id := projekti_seq.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_dokumenti_bi
BEFORE INSERT ON Dokumenti
FOR EACH ROW
BEGIN
    IF :NEW.dokument_id IS NULL THEN
        :NEW.dokument_id := dokumenti_seq.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_zadaci_bi
BEFORE INSERT ON Zadaci
FOR EACH ROW
BEGIN
    IF :NEW.zadatak_id IS NULL THEN
        :NEW.zadatak_id := zadaci_seq.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_korisnici_bi
BEFORE INSERT ON Korisnici
FOR EACH ROW
BEGIN
    IF :NEW.korisnik_id IS NULL THEN
        :NEW.korisnik_id := korisnici_seq.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_log_aktivnosti_bi
BEFORE INSERT ON LogAktivnosti
FOR EACH ROW
BEGIN
    IF :NEW.log_id IS NULL THEN
        :NEW.log_id := log_aktivnosti_seq.NEXTVAL;
    END IF;
END;
/

-- ============================================================================
-- AFTER INSERT/UPDATE/DELETE Triggers - Activity Logging
-- ============================================================================

CREATE OR REPLACE TRIGGER trg_log_dokumenta
AFTER INSERT OR UPDATE OR DELETE ON Dokumenti
FOR EACH ROW
DECLARE
    v_action VARCHAR2(20);
    v_opis VARCHAR2(255);
BEGIN
    IF INSERTING THEN
        v_action := 'IMPORT';
        v_opis := 'Uvoz početne verzije';
    ELSIF UPDATING THEN
        v_action := 'UPDATE';
        v_opis := 'Ažuriran dokument';
    ELSIF DELETING THEN
        v_action := 'DELETE';
        v_opis := 'Obrisan dokument';
    END IF;

    INSERT INTO LogAktivnosti(korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, naziv_entiteta, opis)
    VALUES (
        NVL(:NEW.kreirao_korisnik_id, :OLD.kreirao_korisnik_id),
        v_action,
        'DOKUMENT',
        NVL(:NEW.dokument_id, :OLD.dokument_id),
        NVL(:NEW.naziv_dokumenta, :OLD.naziv_dokumenta),
        v_opis
    );
END;
/

CREATE OR REPLACE TRIGGER trg_log_projekta
AFTER INSERT OR UPDATE OR DELETE ON Projekti
FOR EACH ROW
DECLARE
    v_action VARCHAR2(20);
    v_opis VARCHAR2(255);
BEGIN
    IF INSERTING THEN
        v_action := 'KREIRANJE';
        v_opis := 'Kreiran projekat';
    ELSIF UPDATING THEN
        v_action := 'UPDATE';
        v_opis := 'Ažuriran projekat';
    ELSIF DELETING THEN
        v_action := 'DELETE';
        v_opis := 'Obrisan projekat';
    END IF;

    INSERT INTO LogAktivnosti(korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, naziv_entiteta, opis)
    VALUES (
        NVL(:NEW.rukovodilac_id, :OLD.rukovodilac_id),
        v_action,
        'PROJEKAT',
        NVL(:NEW.projekat_id, :OLD.projekat_id),
        NVL(:NEW.naziv_projekta, :OLD.naziv_projekta),
        v_opis
    );
END;
/
