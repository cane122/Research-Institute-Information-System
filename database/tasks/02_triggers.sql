-- ============================================================================
-- PL/SQL Triggers for Document Management System
-- Automatsko logovanje kada se kreira dokument
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
