-- Dodavanje javnih dokumenata sa pristupom za sve korisnike
-- Pokrenite ovu skriptu u SQL Developer-u (F5)

-- Prvo proveri da li postoji korisnik sa ID 1
DECLARE
    v_user_count NUMBER;
    v_user_id NUMBER := 1;
BEGIN
    SELECT COUNT(*) INTO v_user_count FROM Korisnici WHERE korisnik_id = 1;
    IF v_user_count = 0 THEN
        -- Koristi prvog dostupnog korisnika
        SELECT MIN(korisnik_id) INTO v_user_id FROM Korisnici;
        IF v_user_id IS NULL THEN
            DBMS_OUTPUT.PUT_LINE('GREŠKA: Nema korisnika u bazi. Prvo pokrenite dummy_data_oracle.sql');
            RETURN;
        END IF;
    END IF;
    DBMS_OUTPUT.PUT_LINE('Koristim korisnik_id: ' || v_user_id);
END;
/

-- Ubaci testne dokumente
INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kljucne_reci, kreirao_korisnik_id, datuma_postavke, poslednja_izmena)
VALUES (NULL, 'Uputstvo za korišćenje sistema', 'Kompletno korisničko uputstvo za rad sa informacionim sistemom istraživačkog instituta.', 'PDF', 'Srpski', 'uputstvo, pomoć, dokumentacija, vodič', 1, SYSTIMESTAMP, SYSTIMESTAMP);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kljucne_reci, kreirao_korisnik_id, datuma_postavke, poslednja_izmena)
VALUES (NULL, 'Pravilnik o zaštiti podataka', 'Pravilnik o zaštiti ličnih podataka i informacionoj bezbednosti instituta.', 'PDF', 'Srpski', 'pravilnik, GDPR, zaštita podataka, bezbednost', 1, SYSTIMESTAMP, SYSTIMESTAMP);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kljucne_reci, kreirao_korisnik_id, datuma_postavke, poslednja_izmena)
VALUES (NULL, 'Šablon istraživačkog izveštaja', 'Standardni šablon za pisanje istraživačkih izveštaja.', 'DOCX', 'Srpski', 'šablon, izveštaj, istraživanje, template', 1, SYSTIMESTAMP, SYSTIMESTAMP);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kljucne_reci, kreirao_korisnik_id, datuma_postavke, poslednja_izmena)
VALUES (NULL, 'Metodologija naučnog istraživanja', 'Priručnik o metodama i tehnikama naučnog istraživanja.', 'PDF', 'Srpski', 'metodologija, istraživanje, nauka, metode', 1, SYSTIMESTAMP, SYSTIMESTAMP);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kljucne_reci, kreirao_korisnik_id, datuma_postavke, poslednja_izmena)
VALUES (NULL, 'Etički kodeks istraživača', 'Etički principi i standardi za istraživače instituta.', 'PDF', 'Srpski', 'etika, kodeks, standardi, principi', 1, SYSTIMESTAMP, SYSTIMESTAMP);

-- Dodaj dozvole za sve korisnike (javni pristup)
-- Uzmi ID-jeve upravo ubačenih dokumenata
DECLARE
    CURSOR c_docs IS 
        SELECT dokument_id FROM Dokumenti 
        WHERE naziv_dokumenta IN (
            'Uputstvo za korišćenje sistema',
            'Pravilnik o zaštiti podataka',
            'Šablon istraživačkog izveštaja',
            'Metodologija naučnog istraživanja',
            'Etički kodeks istraživača'
        );
    CURSOR c_users IS SELECT korisnik_id FROM Korisnici;
BEGIN
    FOR doc IN c_docs LOOP
        FOR usr IN c_users LOOP
            -- Dodaj dozvolu za čitanje svim korisnicima
            INSERT INTO DozvoleDokumenata (dokument_id, korisnik_id, moze_citati, moze_menjati, moze_brisati)
            VALUES (doc.dokument_id, usr.korisnik_id, 1, 0, 0);
        END LOOP;
    END LOOP;
    DBMS_OUTPUT.PUT_LINE('Dozvole dodeljene svim korisnicima.');
END;
/

COMMIT;

-- Prikaži rezultat
SELECT d.dokument_id, d.naziv_dokumenta, d.tip_dokumenta, 
       (SELECT COUNT(*) FROM DozvoleDokumenata dd WHERE dd.dokument_id = d.dokument_id) as broj_dozvola
FROM Dokumenti d
WHERE d.naziv_dokumenta IN (
    'Uputstvo za korišćenje sistema',
    'Pravilnik o zaštiti podataka',
    'Šablon istraživačkog izveštaja',
    'Metodologija naučnog istraživanja',
    'Etički kodeks istraživača'
)
ORDER BY d.dokument_id;

PROMPT 'Uspešno dodato 5 javnih dokumenata sa pristupom za sve korisnike!';
