-- Dummy podaci za Research Institute Information System
-- Encoding: UTF-8

SET client_encoding = 'UTF8';

-- Početak transakcije
BEGIN;

-- 1. KREIRANJE TEST KORISNIKA
INSERT INTO Korisnici (korisnik_id, korisnicko_ime, email, hash_sifre, ime, prezime, uloga_id, status) VALUES
(1, 'admin', 'admin@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$fg2s2gmJk7467X4Hd4AggCjOQ3NwUhaWTQhwz7vfxhc$CDebI9A9GMptWEpovQO0YB+P6C3gSbSeM7GoIRqXWSU', 'Marko', 'Petrovic', 1, 'aktivan'),
(2, 'researcher1', 'researcher1@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Ana', 'Petrovic', 3, 'aktivan'),
(3, 'manager1', 'manager1@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Petar', 'Jovanovic', 2, 'aktivan'),
(4, 'manager2', 'manager2@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Milica', 'Nikolic', 2, 'aktivan'),
(5, 'researcher2', 'researcher2@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Stefan', 'Milic', 3, 'aktivan'),
(6, 'researcher3', 'researcher3@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Jovana', 'Stojanovic', 3, 'aktivan'),
(7, 'researcher4', 'researcher4@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Marko', 'Lazic', 3, 'aktivan'),
(8, 'researcher5', 'researcher5@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Tamara', 'Radic', 3, 'aktivan'),
(9, 'researcher6', 'researcher6@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Nikola', 'Peric', 3, 'aktivan'),
(10, 'researcher7', 'researcher7@institut.rs', '$argon2id$v=19$m=65536,t=1,p=4$6qe2HIn5Ekg/JbfLixqCvtSJPc4wBC68CyI6YwVf61M$u3Ms2NZr0QSPv4IQCRAVNrpNFY29/YeMr6CHyWx8+JE', 'Milena', 'Savic', 3, 'aktivan');
-- 2. KREIRANJE TEST PROJEKATA
INSERT INTO Projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id, radni_tok_id) VALUES 
('AI u Zdravstvu', 'Implementacija vestacke inteligencije u dijagnostici medicinskih slika', '2025-01-15', '2025-12-31', 'Aktivan', 3, 2),
('Pametni Gradovi IoT', 'Razvoj IoT sistema za upravljanje javnim osvetljenjem i prometom', '2025-02-01', '2026-01-31', 'Aktivan', 4, 1),
('Kvantno Racunarstvo', 'Istrazivanje primene kvantnih algoritama u kriptografiji', '2024-09-01', '2025-08-31', 'Aktivan', 3, 2),
('Blockchain Identiteti', 'Decentralizovani sistem za upravljanje digitalnim identitetima', '2025-03-01', '2025-11-30', 'Aktivan', 4, 1),
('Obnovljiva Energija', 'Optimizacija solarnih panela pomocu machine learning algoritma', '2024-11-01', '2025-10-31', 'Aktivan', 3, 2);

-- 3. DODAVANJE CLANOVA PROJEKATA
INSERT INTO ClanoviProjekta (projekat_id, korisnik_id) VALUES 
-- AI u Zdravstvu tim
(1, 3), (1, 5), (1, 6), (1, 9),
-- Pametni Gradovi tim  
(2, 4), (2, 7), (2, 8), (2, 10),
-- Kvantno Racunarstvo tim
(3, 3), (3, 5), (3, 7),
-- Blockchain tim
(4, 4), (4, 6), (4, 8), (4, 9),
-- Obnovljiva Energija tim
(5, 3), (5, 5), (5, 6), (5, 7), (5, 10);

-- 4. KREIRANJE ZADATAKA
INSERT INTO Zadaci (projekat_id, faza_id, naziv_zadatka, opis, dodeljen_korisniku_id, rok, prioritet, progres) VALUES 
-- AI u Zdravstvu zadaci
(1, 6, 'Definisanje zahteva za AI model', 'Analiza medicinskih standarda i zahteva za dijagnostiku', 5, '2025-02-15', 'Visok', 100),
(1, 7, 'Prikupljanje medicinskih slika', 'Kreiranje dataseta za treniranje AI modela', 6, '2025-03-30', 'Visok', 80),
(1, 8, 'Implementacija CNN algoritma', 'Razvoj konvolucijskog neuronskog modela', 5, '2025-05-15', 'Visok', 60),

-- Pametni Gradovi zadaci
(2, 1, 'Analiza postojece infrastrukture', 'Mapiranje trenutnih sistema javnog osvetljenja', 7, '2025-03-15', 'Srednji', 90),
(2, 2, 'Dizajn IoT senzora', 'Specifikacija senzora za monitoring prometa', 8, '2025-04-30', 'Visok', 70),
(2, 3, 'Prototip mobilne aplikacije', 'Razvoj aplikacije za gradjanе', 7, '2025-06-30', 'Srednji', 40),

-- Kvantno Racunarstvo zadaci
(3, 7, 'Implementacija Shor algoritma', 'Kvantni algoritam za faktorizaciju velikih brojeva', 5, '2025-04-15', 'Visok', 85),
(3, 8, 'Testiranje na kvantnom simulatoru', 'Validacija algoritma na IBM Quantum simulatoru', 7, '2025-05-30', 'Visok', 45),

-- Blockchain zadaci
(4, 2, 'Smart contract za identitete', 'Ethereum smart contract za decentralizovane ID', 6, '2025-04-20', 'Visok', 75),
(4, 3, 'Web3 frontend aplikacija', 'React aplikacija za upravljanje identitetima', 8, '2025-06-15', 'Srednji', 50);

-- 5. KREIRANJE FOLDERA
INSERT INTO Folderi (naziv_foldera, roditelj_folder_id, vlasnik_id) VALUES 
('Projekti 2025', NULL, 1),
('AI Dokumenti', 1, 3),
('IoT Dokumenti', 1, 4),
('Kvantni Algoritmi', 1, 3),
('Blockchain Docs', 1, 4),
('Izvestaji', NULL, 2),
('Finansijski Izvestaji', 6, 2),
('Tehnicki Izvestaji', 6, 3);

-- 6. KREIRANJE DOKUMENATA
INSERT INTO Dokumenti (projekat_id, naziv_dokumenta, folder_id, opis, tip_dokumenta, jezik_dokumenta, radni_tok_id, trenutna_faza_id, kreirao_korisnik_id) VALUES 
(1, 'Specifikacija AI Modela v1.2', 2, 'Detaljne specifikacije za CNN model dijagnostike', 'Specifikacija', 'srpski', 3, 11, 5),
(1, 'Dataset Medicinskih Slika', 2, 'Kolekcija od 10000 anotiranih medicinskih slika', 'Dataset', 'engleski', 3, 12, 6),
(2, 'IoT Sensor Protokol', 3, 'Komunikacijski protokol za IoT senzore', 'Protokol', 'srpski', 3, 11, 7),
(2, 'Mobilna App Wireframes', 3, 'UI/UX dizajn za gradjansku aplikaciju', 'Dizajn', 'srpski', 3, 12, 8),
(3, 'Kvantni Algoritmi - Implementacija', 4, 'Python kod za Shor i Grover algoritme', 'Kod', 'engleski', 3, 13, 5),
(4, 'Blockchain Arhitektura', 5, 'Sistemska arhitektura decentralizovanog ID sistema', 'Arhitektura', 'srpski', 3, 13, 6),
(NULL, 'Godisnji Izvestaj 2024', 7, 'Finansijski izvestaj instituta za 2024. godinu', 'Izvestaj', 'srpski', 3, 14, 2);

-- 7. VERZIJE DOKUMENATA
INSERT INTO VerzijeDokumenata (dokument_id, verzija_oznaka, putanja_do_fajla, velicina_fajla_MB, postavio_korisnik_id) VALUES 
(1, 'v1.0', '/docs/ai_specifikacija_v1.0.pdf', 2.5, 5),
(1, 'v1.1', '/docs/ai_specifikacija_v1.1.pdf', 2.7, 5),
(1, 'v1.2', '/docs/ai_specifikacija_v1.2.pdf', 3.1, 5),
(2, 'v1.0', '/datasets/medical_images.zip', 1024.5, 6),
(3, 'v1.0', '/docs/iot_protocol.pdf', 1.8, 7),
(3, 'v1.1', '/docs/iot_protocol_v1.1.pdf', 2.2, 7),
(4, 'v1.0', '/designs/mobile_wireframes.sketch', 15.3, 8),
(5, 'v1.0', '/code/quantum_algorithms.py', 0.5, 5),
(6, 'v1.0', '/docs/blockchain_architecture.pdf', 4.2, 6),
(7, 'v1.0', '/reports/godisnji_izvestaj_2024.pdf', 8.7, 2);

-- 8. TAGOVI
INSERT INTO Tagovi (naziv_taga) VALUES 
('AI'),
('Machine Learning'),
('IoT'),
('Blockchain'),
('Kvantno Racunarstvo'),
('Pametni Gradovi'),
('Zdravstvo'),
('Finansije'),
('Izvestavanje'),
('Prototip');

-- 9. LINKOVANJE DOKUMENATA I TAGOVA
INSERT INTO DokumentTagovi (dokument_id, tag_id) VALUES 
(1, 1), (1, 2), (1, 7),  -- AI spec: AI, ML, Zdravstvo
(2, 1), (2, 2), (2, 7),  -- Dataset: AI, ML, Zdravstvo  
(3, 3), (3, 6),          -- IoT protokol: IoT, Pametni Gradovi
(4, 3), (4, 6), (4, 10), -- Wireframes: IoT, Pametni Gradovi, Prototip
(5, 5),                  -- Kvantni kod: Kvantno Racunarstvo
(6, 4),                  -- Blockchain: Blockchain
(7, 8), (7, 9);         -- Izvestaj: Finansije, Izvestavanje

-- 10. LLM SAZECI
INSERT INTO LLMSazeci (dokument_id, verzija_oznaka, sazetak) VALUES 
(1, 'v1.2', 'Specifikacija definise CNN arhitekturu sa 5 konvolucijskih slojeva za klasifikaciju medicinskih slika. Model koristi ResNet backbone sa accuracy od 94.2% na test datasetu. Ukljucuje data augmentation tehnike i transfer learning pristup.'),
(2, 'v1.0', 'Dataset sadrzi 10,000 anotiranih rendgenskih slika grudnog kosa kategorizovanih u 14 klasa patologija. Slike su u DICOM formatu, rezolucije 1024x1024 piksela. Dataset je podeljen 70/15/15 za train/validation/test skupove.'),
(3, 'v1.1', 'Protokol definise MQTT komunikaciju između IoT senzora i centralne platforme. Koristi JSON format za razmenu podataka sa enkriptovanjem AES-256. Implementira heartbeat mehanizam svakih 30 sekundi.'),
(6, 'v1.0', 'Arhitektura koristi Ethereum blockchain sa custom ERC-721 tokenima za digitalne identitete. Implementira zero-knowledge proof protokol za privatnost. Frontend koristi Web3.js biblioteku za interakciju sa smart contract-ima.');

-- 11. KOMENTARI NA ZADACIMA
INSERT INTO KomentariZadataka (zadatak_id, korisnik_id, tekst_komentara) VALUES 
(1, 3, 'Zahtevi su uspesno definisani u saradnji sa klinickim partnerima.'),
(1, 5, 'Dodao sam dodatne metrike za evaluaciju modela u finalnu verziju.'),
(2, 6, 'Prikupljeno je 8,500 slika do sada. Potrebno jos 1,500 za kompletan dataset.'),
(3, 5, 'Implementacija je gotova. Pocinje testiranje performansi na test datasetu.'),
(4, 7, 'Analiza je gotova. Identifikovano je 1,200 lampi za upgrade na pametne senzore.'),
(5, 8, 'Prototip senzora je spreman. Testiranje u realnim uslovima slede nedelju.'),
(7, 5, 'Shor algoritam uspesno implementiran. Testiranje na IBM Quantum Cloud sledi.');

-- 12. LOG AKTIVNOSTI
INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, opis, ciljani_entitet, ciljani_id) VALUES 
(5, 'KREIRAN_DOKUMENT', 'Kreiran dokument Specifikacija AI Modela v1.0', 'Dokument', 1),
(5, 'AZURIRAN_DOKUMENT', 'Azurirana specifikacija na verziju v1.2', 'Dokument', 1),
(6, 'UPLOAD_DOKUMENTA', 'Postavljen dataset medicinskih slika', 'Dokument', 2),
(3, 'KREIRAN_PROJEKAT', 'Kreiran novi projekat AI u Zdravstvu', 'Projekat', 1),
(4, 'KREIRAN_PROJEKAT', 'Kreiran novi projekat Pametni Gradovi IoT', 'Projekat', 2),
(5, 'PROMENA_FAZE', 'Zadatak premesten u fazu Implementacija', 'Zadatak', 3),
(7, 'ZAVRSIO_ZADATAK', 'Kompletirana analiza postojece infrastrukture', 'Zadatak', 4),
(8, 'KOMENTAR_ZADATAK', 'Dodat komentar na zadatak dizajn senzora', 'Zadatak', 5),
(2, 'KREIRAN_DOKUMENT', 'Kreiran godisnji finansijski izvestaj', 'Dokument', 7),
(1, 'LOGIN', 'Administrator se ulogovao u sistem', 'Korisnik', 1);

-- Commit transakcije
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

-- Prikaz osnovnih statistika
SELECT 
    'PROJEKTI PO STATUSU' as info,
    status,
    COUNT(*) as broj
FROM Projekti 
GROUP BY status
UNION ALL
SELECT 
    'ZADACI PO PRIORITETU',
    prioritet,
    COUNT(*)
FROM Zadaci 
GROUP BY prioritet
ORDER BY info, broj DESC;
