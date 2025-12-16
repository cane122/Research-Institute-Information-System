-- Dummy podaci za Research Institute Information System
-- Created for Oracle Database
-- Encoding: UTF-8

-- 1. KREIRANJE TEST KORISNIKA (bez eksplicitnog ID-a zbog IDENTITY kolona)
INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('admin', 'admin@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$fg2s2gmJk7467X4Hd4AggCjOQ3NwUhaWTQhwz7vfxhc$CDebI9A9GMptWEpovQO0YB+P6C3gSbSeM7GoIRqXWSU', 'Marko', 'Petrovic', 1, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('researcher1', 'researcher1@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Ana', 'Petrovic', 3, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('manager1', 'manager1@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Petar', 'Jovanovic', 2, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('manager2', 'manager2@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Milica', 'Nikolic', 2, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('researcher2', 'researcher2@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Stefan', 'Milic', 3, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('researcher3', 'researcher3@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Jovana', 'Stojanovic', 3, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('researcher4', 'researcher4@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Marko', 'Lazic', 3, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('researcher5', 'researcher5@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Tamara', 'Radic', 3, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('researcher6', 'researcher6@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Nikola', 'Peric', 3, 'aktivan');

INSERT INTO Korisnici (korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
('researcher7', 'researcher7@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Milena', 'Savic', 3, 'aktivan');

-- 2. KREIRANJE TEST PROJEKATA
INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id) VALUES 
('AI u Zdravstvu', 'Implementacija vestacke inteligencije u dijagnostici medicinskih slika', TO_DATE('2025-01-15', 'YYYY-MM-DD'), TO_DATE('2025-12-31', 'YYYY-MM-DD'), 'Aktivan', 3, 2);

INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id) VALUES 
('Pametni Gradovi IoT', 'Razvoj IoT sistema za upravljanje javnim osvetljenjem i prometom', TO_DATE('2025-02-01', 'YYYY-MM-DD'), TO_DATE('2026-01-31', 'YYYY-MM-DD'), 'Aktivan', 4, 1);

INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id) VALUES 
('Kvantno Racunarstvo', 'Istrazivanje primene kvantnih algoritama u kriptografiji', TO_DATE('2024-09-01', 'YYYY-MM-DD'), TO_DATE('2025-08-31', 'YYYY-MM-DD'), 'Aktivan', 3, 2);

INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id) VALUES 
('Blockchain Identiteti', 'Decentralizovani sistem za upravljanje digitalnim identitetima', TO_DATE('2025-03-01', 'YYYY-MM-DD'), TO_DATE('2025-11-30', 'YYYY-MM-DD'), 'Aktivan', 4, 1);

INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id) VALUES 
('Obnovljiva Energija', 'Optimizacija solarnih panela pomocu machine learning algoritma', TO_DATE('2024-11-01', 'YYYY-MM-DD'), TO_DATE('2025-10-31', 'YYYY-MM-DD'), 'Aktivan', 3, 2);

-- 3. DODAVANJE CLANOVA PROJEKATA
-- AI u Zdravstvu tim
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (1, 3);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (1, 5);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (1, 6);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (1, 9);

-- Pametni Gradovi tim  
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (2, 4);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (2, 7);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (2, 8);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (2, 10);

-- Kvantno Racunarstvo tim
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (3, 3);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (3, 5);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (3, 7);

-- Blockchain tim
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (4, 4);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (4, 6);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (4, 8);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (4, 9);

-- Obnovljiva Energija tim
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (5, 3);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (5, 5);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (5, 6);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (5, 7);
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES (5, 10);

-- 4. KREIRANJE ZADATAKA
-- AI u Zdravstvu zadaci
INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(1, 6, 'Definisanje zahteva za AI model', 'Analiza medicinskih standarda i zahteva za dijagnostiku', 5, TO_DATE('2025-02-15', 'YYYY-MM-DD'), 'Visok', 100);

INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(1, 7, 'Prikupljanje medicinskih slika', 'Kreiranje dataseta za treniranje AI modela', 6, TO_DATE('2025-03-30', 'YYYY-MM-DD'), 'Visok', 80);

INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(1, 8, 'Implementacija CNN algoritma', 'Razvoj konvolucijskog neuronskog modela', 5, TO_DATE('2025-05-15', 'YYYY-MM-DD'), 'Visok', 60);

-- Pametni Gradovi zadaci
INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(2, 1, 'Analiza postojece infrastrukture', 'Mapiranje trenutnih sistema javnog osvetljenja', 7, TO_DATE('2025-03-15', 'YYYY-MM-DD'), 'Srednji', 90);

INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(2, 2, 'Dizajn IoT senzora', 'Specifikacija senzora za monitoring prometa', 8, TO_DATE('2025-04-30', 'YYYY-MM-DD'), 'Visok', 70);

INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(2, 3, 'Prototip mobilne aplikacije', 'Razvoj aplikacije za gradjane', 7, TO_DATE('2025-06-30', 'YYYY-MM-DD'), 'Srednji', 40);

-- Kvantno Racunarstvo zadaci
INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(3, 7, 'Implementacija Shor algoritma', 'Kvantni algoritam za faktorizaciju velikih brojeva', 5, TO_DATE('2025-04-15', 'YYYY-MM-DD'), 'Visok', 85);

INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(3, 8, 'Testiranje na kvantnom simulatoru', 'Validacija algoritma na IBM Quantum simulatoru', 7, TO_DATE('2025-05-30', 'YYYY-MM-DD'), 'Visok', 45);

-- Blockchain zadaci
INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(4, 2, 'Smart contract za identitete', 'Ethereum smart contract za decentralizovane ID', 6, TO_DATE('2025-04-20', 'YYYY-MM-DD'), 'Visok', 75);

INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
(4, 3, 'Web3 frontend aplikacija', 'React aplikacija za upravljanje identitetima', 8, TO_DATE('2025-06-15', 'YYYY-MM-DD'), 'Srednji', 50);

-- 5. KREIRANJE FOLDERA
INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('Projekti 2025', NULL, 1);

INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('AI Dokumenti', 1, 3);

INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('IoT Dokumenti', 1, 4);

INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('Kvantni Algoritmi', 1, 3);

INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('Blockchain Docs', 1, 4);

INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('Izvestaji', NULL, 2);

INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('Finansijski Izvestaji', 6, 2);

INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('Tehnicki Izvestaji', 6, 3);

-- 6. KREIRANJE DOKUMENATA
INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id) VALUES 
(1, 'Specifikacija AI Modela v1.2', 2, 'Detaljne specifikacije za CNN model dijagnostike', 'Specifikacija', 'srpski', 3, 11, 5);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id) VALUES 
(1, 'Dataset Medicinskih Slika', 2, 'Kolekcija od 10000 anotiranih medicinskih slika', 'Dataset', 'engleski', 3, 12, 6);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id) VALUES 
(2, 'IoT Sensor Protokol', 3, 'Komunikacijski protokol za IoT senzore', 'Protokol', 'srpski', 3, 11, 7);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id) VALUES 
(2, 'Mobilna App Wireframes', 3, 'UI/UX dizajn za gradjansku aplikaciju', 'Dizajn', 'srpski', 3, 12, 8);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id) VALUES 
(3, 'Kvantni Algoritmi - Implementacija', 4, 'Python kod za Shor i Grover algoritme', 'Kod', 'engleski', 3, 13, 5);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id) VALUES 
(4, 'Blockchain Arhitektura', 5, 'Sistemska arhitektura decentralizovanog ID sistema', 'Arhitektura', 'srpski', 3, 13, 6);

INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id) VALUES 
(NULL, 'Godisnji Izvestaj 2024', 7, 'Finansijski izvestaj instituta za 2024. godinu', 'Izvestaj', 'srpski', 3, 14, 2);

-- 7. VERZIJE DOKUMENATA
INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(1, 'v1.0', '/docs/ai_specifikacija_v1.0.pdf', 2.5, 5);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(1, 'v1.1', '/docs/ai_specifikacija_v1.1.pdf', 2.7, 5);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(1, 'v1.2', '/docs/ai_specifikacija_v1.2.pdf', 3.1, 5);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(2, 'v1.0', '/datasets/medical_images.zip', 1024.5, 6);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(3, 'v1.0', '/docs/iot_protocol.pdf', 1.8, 7);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(3, 'v1.1', '/docs/iot_protocol_v1.1.pdf', 2.2, 7);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(4, 'v1.0', '/designs/mobile_wireframes.sketch', 15.3, 8);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(5, 'v1.0', '/code/quantum_algorithms.py', 0.5, 5);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(6, 'v1.0', '/docs/blockchain_architecture.pdf', 4.2, 6);

INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(7, 'v1.0', '/reports/godisnji_izvestaj_2024.pdf', 8.7, 2);

-- 8. TAGOVI
INSERT INTO Tagovi (naziv_taga) VALUES ('AI');
INSERT INTO Tagovi (naziv_taga) VALUES ('Machine Learning');
INSERT INTO Tagovi (naziv_taga) VALUES ('IoT');
INSERT INTO Tagovi (naziv_taga) VALUES ('Blockchain');
INSERT INTO Tagovi (naziv_taga) VALUES ('Kvantno Racunarstvo');
INSERT INTO Tagovi (naziv_taga) VALUES ('Pametni Gradovi');
INSERT INTO Tagovi (naziv_taga) VALUES ('Zdravstvo');
INSERT INTO Tagovi (naziv_taga) VALUES ('Finansije');
INSERT INTO Tagovi (naziv_taga) VALUES ('Izvestavanje');
INSERT INTO Tagovi (naziv_taga) VALUES ('Prototip');

-- 9. LINKOVANJE DOKUMENATA I TAGOVA
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (1, 1);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (1, 2);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (1, 7);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (2, 1);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (2, 2);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (2, 7);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (3, 3);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (3, 6);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (4, 3);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (4, 6);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (4, 10);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (5, 5);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (6, 4);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (7, 8);
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES (7, 9);

-- 10. LLM SAZECI
INSERT INTO LLMSazeci (dokument_id, verzija_oznaka, sazetak) VALUES 
(1, 'v1.2', 'Specifikacija definise CNN arhitekturu sa 5 konvolucijskih slojeva za klasifikaciju medicinskih slika. Model koristi ResNet backbone sa accuracy od 94.2% na test datasetu. Ukljucuje data augmentation tehnike i transfer learning pristup.');

INSERT INTO LLMSazeci (dokument_id, verzija_oznaka, sazetak) VALUES 
(2, 'v1.0', 'Dataset sadrzi 10,000 anotiranih rendgenskih slika grudnog kosa kategorizovanih u 14 klasa patologija. Slike su u DICOM formatu, rezolucije 1024x1024 piksela. Dataset je podeljen 70/15/15 za train/validation/test skupove.');

INSERT INTO LLMSazeci (dokument_id, verzija_oznaka, sazetak) VALUES 
(3, 'v1.1', 'Protokol definise MQTT komunikaciju između IoT senzora i centralne platforme. Koristi JSON format za razmenu podataka sa enkriptovanjem AES-256. Implementira heartbeat mehanizam svakih 30 sekundi.');

INSERT INTO LLMSazeci (dokument_id, verzija_oznaka, sazetak) VALUES 
(6, 'v1.0', 'Arhitektura koristi Ethereum blockchain sa custom ERC-721 tokenima za digitalne identitete. Implementira zero-knowledge proof protokol za privatnost. Frontend koristi Web3.js biblioteku za interakciju sa smart contract-ima.');

-- 11. KOMENTARI NA ZADACIMA
INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES 
(1, 3, 'Zahtevi su uspesno definisani u saradnji sa klinickim partnerima.');

INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES 
(1, 5, 'Dodao sam dodatne metrike za evaluaciju modela u finalnu verziju.');

INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES 
(2, 6, 'Prikupljeno je 8,500 slika do sada. Potrebno jos 1,500 za kompletan dataset.');

INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES 
(3, 5, 'Implementacija je gotova. Pocinje testiranje performansi na test datasetu.');

INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES 
(4, 7, 'Analiza je gotova. Identifikovano je 1,200 lampi za upgrade na pametne senzore.');

INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES 
(5, 8, 'Prototip senzora je spreman. Testiranje u realnim uslovima slede nedelju.');

INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES 
(7, 5, 'Shor algoritam uspesno implementiran. Testiranje na IBM Quantum Cloud sledi.');

-- 12. LOG AKTIVNOSTI
INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(5, 'KREIRAN_DOKUMENT', 'Kreiran dokument Specifikacija AI Modela v1.0', 'Dokument', 1);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(5, 'AZURIRAN_DOKUMENT', 'Azurirana specifikacija na verziju v1.2', 'Dokument', 1);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(6, 'UPLOAD_DOKUMENTA', 'Postavljen dataset medicinskih slika', 'Dokument', 2);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(3, 'KREIRAN_PROJEKAT', 'Kreiran novi projekat AI u Zdravstvu', 'Projekat', 1);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(4, 'KREIRAN_PROJEKAT', 'Kreiran novi projekat Pametni Gradovi IoT', 'Projekat', 2);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(5, 'PROMENA_FAZE', 'Zadatak premesten u fazu Implementacija', 'Zadatak', 3);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(7, 'ZAVRSIO_ZADATAK', 'Kompletirana analiza postojece infrastrukture', 'Zadatak', 4);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(8, 'KOMENTAR_ZADATAK', 'Dodat komentar na zadatak dizajn senzora', 'Zadatak', 5);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(2, 'KREIRAN_DOKUMENT', 'Kreiran godisnji finansijski izvestaj', 'Dokument', 7);

INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(1, 'LOGIN', 'Administrator se ulogovao u sistem', 'Korisnik', 1);

COMMIT;

-- Verifikacija podataka
SELECT 'KORISNICI' as tabela, COUNT(*) as broj_redova FROM Korisnici
UNION ALL
SELECT 'PROJEKTI', COUNT(*) FROM Projekti
UNION ALL
SELECT 'ZADACI', COUNT(*) FROM Zadaci
UNION ALL
SELECT 'DOKUMENTI', COUNT(*) FROM Dokumenti
UNION ALL
SELECT 'VERZIJE', COUNT(*) FROM VerzijeDokumenata
UNION ALL
SELECT 'TAGOVI', COUNT(*) FROM Tagovi
UNION ALL
SELECT 'LOG AKTIVNOSTI', COUNT(*) FROM LogAktivnosti;

-- End of dummy data
