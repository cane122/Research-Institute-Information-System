-- Update script to add conditions for all phases
-- This script adds comprehensive conditions for all workflow phases

-- Remove old conditions (optional - uncomment if you want to replace existing)
-- DELETE FROM ProcenaUslova;
-- DELETE FROM Uslovi;

-- Add conditions only if they don't already exist for each phase

-- Uslovi za STANDARDNI PROJEKTNI TOK (radni_tok_id = 1)

-- Faza 1: Planiranje
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(1, 'Definisati ciljeve projekta', 'Svi ciljevi projekta moraju biti jasno definisani i dokumentovani'),
(1, 'Identifikovati resurse', 'Potrebni resursi (ljudski, finansijski, tehnički) moraju biti identifikovani'),
(1, 'Kreirati plan projekta', 'Plan projekta sa vremenskom linijom mora biti kreiran i odobren')
ON CONFLICT DO NOTHING;

-- Faza 2: Analiza
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(2, 'Prikupiti zahteve', 'Svi funkcionalni i nefunkcionalni zahtevi moraju biti prikupljeni'),
(2, 'Analizirati rizike', 'Analiza rizika mora biti završena sa planovima za mitigaciju'),
(2, 'Odobriti specifikaciju', 'Specifikacija zahteva mora biti odobrena od strane stejkholdera')
ON CONFLICT DO NOTHING;

-- Faza 3: Razvoj
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(3, 'Implementirati osnovne funkcionalnosti', 'Sve osnovne funkcionalnosti moraju biti kodirane'),
(3, 'Napisati unit testove', 'Pokrivenost testovima mora biti minimum 70%'),
(3, 'Code review', 'Kod mora proći code review proces')
ON CONFLICT DO NOTHING;

-- Faza 4: Testiranje
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(4, 'Sve testove prolaze', 'Svi automatski testovi moraju biti zeleni'),
(4, 'QA provera', 'QA tim mora odobriti funkcionalnost'),
(4, 'Performance testiranje', 'Aplikacija mora zadovoljiti performance kriterijume')
ON CONFLICT DO NOTHING;

-- Faza 5: Završeno
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(5, 'Deployment uspešan', 'Projekat mora biti uspešno deployovan u produkciju'),
(5, 'Dokumentacija kompletna', 'Sva tehnička i korisnička dokumentacija mora biti završena'),
(5, 'Sign-off od klijenta', 'Klijent mora potpisati prihvatanje projekta')
ON CONFLICT DO NOTHING;

-- Uslovi za ISTRAŽIVAČKI TOK (radni_tok_id = 2)

-- Faza 6: Definisanje istraživanja
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(6, 'Formulisati istraživačko pitanje', 'Jasno istraživačko pitanje mora biti formulisano'),
(6, 'Pregled literature', 'Pregled postojeće literature mora biti završen'),
(6, 'Metodologija definisana', 'Metodologija istraživanja mora biti definisana i odobrena')
ON CONFLICT DO NOTHING;

-- Faza 7: Prikupljanje podataka
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(7, 'Završiti dokumentaciju', 'Sva potrebna tehnička dokumentacija mora biti kompletna'),
(7, 'Prikupiti sve izvore', 'Svi relevantni izvori i reference moraju biti sakupljeni'),
(7, 'Završiti analizu', 'Inicijalna analiza podataka mora biti završena')
ON CONFLICT DO NOTHING;

-- Faza 8: Analiza podataka
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(8, 'Statistička analiza završena', 'Sve planirane statističke analize moraju biti izvršene'),
(8, 'Validacija rezultata', 'Rezultati moraju biti validovani i provereni'),
(8, 'Vizualizacija podataka', 'Grafici i tabele moraju biti pripremljeni')
ON CONFLICT DO NOTHING;

-- Faza 9: Pisanje izveštaja
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(9, 'Napisati draft izveštaja', 'Prvi draft istraživačkog izveštaja mora biti napisan'),
(9, 'Peer review', 'Izveštaj mora proći peer review od strane kolega'),
(9, 'Revizije završene', 'Sve sugerisane revizije moraju biti implementirane')
ON CONFLICT DO NOTHING;

-- Faza 10: Publikovanje
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(10, 'Izbor časopisa', 'Ciljni časopis ili konferencija mora biti izabran'),
(10, 'Formatiranje prema zahtevima', 'Rad mora biti formatiran prema zahtevima časopisa'),
(10, 'Submission', 'Rad mora biti poslat (submitted) za publikaciju')
ON CONFLICT DO NOTHING;

-- Uslovi za DOKUMENTACIONI TOK (radni_tok_id = 3)

-- Faza 11: Kreiranje
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(11, 'Definisati strukturu dokumenta', 'Struktura i sadržaj dokumenta moraju biti definisani'),
(11, 'Napisati prvi draft', 'Inicijalni sadržaj dokumenta mora biti napisan'),
(11, 'Dodati metapodatke', 'Svi potrebni metapodaci moraju biti dodati')
ON CONFLICT DO NOTHING;

-- Faza 12: Revizija
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(12, 'Tehnička provera', 'Dokument mora proći tehničku proveru na tačnost'),
(12, 'Jezička korekcija', 'Jezička i gramatička korekcija mora biti završena'),
(12, 'Provera formatiranja', 'Formatiranje mora biti usklađeno sa standardima')
ON CONFLICT DO NOTHING;

-- Faza 13: Odobravanje
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(13, 'Odobrenje rukovodioca', 'Dokument mora biti odobren od strane rukovodioca'),
(13, 'Usklađenost sa propisima', 'Dokument mora biti u skladu sa svim relevantnim propisima'),
(13, 'Sign-off od stejkholdera', 'Svi ključni stejkholderi moraju odobriti dokument')
ON CONFLICT DO NOTHING;

-- Faza 14: Finalizovanje
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(14, 'Finalna verzija kreirana', 'Finalna verzija dokumenta mora biti kreirana'),
(14, 'PDF export', 'Dokument mora biti exportovan u finalni format (PDF)'),
(14, 'Digitalni potpis', 'Dokument mora biti digitalno potpisan ako je potrebno')
ON CONFLICT DO NOTHING;

-- Faza 15: Arhiviranje
INSERT INTO Uslovi (faza_id, opis, kriterijum) VALUES 
(15, 'Uploadovati u arhivu', 'Dokument mora biti uploadovan u arhivski sistem'),
(15, 'Metadata kompletna', 'Svi metadata za pretragu moraju biti popunjeni'),
(15, 'Backup kreiran', 'Backup kopija mora biti kreirana na sigurnoj lokaciji')
ON CONFLICT DO NOTHING;

-- Verify the data
SELECT 
    f.faza_id,
    f.naziv_faze,
    rt.naziv as radni_tok,
    COUNT(u.uslov_id) as broj_uslova
FROM Faze f
JOIN RadniTokovi rt ON f.radni_tok_id = rt.radni_tok_id
LEFT JOIN Uslovi u ON f.faza_id = u.faza_id
GROUP BY f.faza_id, f.naziv_faze, rt.naziv, f.redosled, rt.radni_tok_id
ORDER BY rt.radni_tok_id, f.redosled;
