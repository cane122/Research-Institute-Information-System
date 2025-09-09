-- Pregled dummy podataka u bazi
\echo '=== PREGLED KORISNIKA ==='
SELECT 
    k.korisnik_id,
    k.korisnicko_ime,
    k.ime || ' ' || k.prezime as puno_ime,
    u.naziv_uloge as uloga,
    k.status
FROM Korisnici k
JOIN Uloge u ON k.uloga_id = u.uloga_id
ORDER BY k.uloga_id, k.ime;

\echo ''
\echo '=== PREGLED PROJEKATA ==='
SELECT 
    p.projekat_id,
    p.naziv_projekta,
    k.ime || ' ' || k.prezime as rukovodilac,
    p.status,
    (SELECT COUNT(*) FROM Zadaci z WHERE z.projekat_id = p.projekat_id) as ukupno_zadataka,
    (SELECT COUNT(*) FROM Zadaci z WHERE z.projekat_id = p.projekat_id AND z.progres = 100) as zavrsenih_zadataka
FROM Projekti p
LEFT JOIN Korisnici k ON p.rukovodilac_id = k.korisnik_id
ORDER BY p.projekat_id;

\echo ''
\echo '=== AKTUELNI ZADACI ==='
SELECT 
    z.zadatak_id,
    z.naziv_zadatka,
    p.naziv_projekta as projekat,
    k.ime || ' ' || k.prezime as dodeljen,
    z.prioritet,
    z.progres || '%' as napredak,
    z.rok
FROM Zadaci z
JOIN Projekti p ON z.projekat_id = p.projekat_id
LEFT JOIN Korisnici k ON z.dodeljen_korisniku_id = k.korisnik_id
WHERE z.progres < 100
ORDER BY 
    CASE z.prioritet WHEN 'Visok' THEN 1 WHEN 'Srednji' THEN 2 ELSE 3 END,
    z.rok;

\echo ''
\echo '=== DOKUMENTI PO PROJEKTIMA ==='
SELECT 
    p.naziv_projekta as projekat,
    d.naziv_dokumenta,
    d.tip_dokumenta,
    k.ime || ' ' || k.prezime as kreirao,
    (SELECT COUNT(*) FROM VerzijeDokumenata vd WHERE vd.dokument_id = d.dokument_id) as broj_verzija
FROM Dokumenti d
LEFT JOIN Projekti p ON d.projekat_id = p.projekat_id
JOIN Korisnici k ON d.kreirao_korisnik_id = k.korisnik_id
ORDER BY p.naziv_projekta NULLS LAST, d.naziv_dokumenta;

\echo ''
\echo '=== POSLEDNJE AKTIVNOSTI ==='
SELECT 
    la.datuma::date as datum,
    k.ime || ' ' || k.prezime as korisnik,
    la.tip_aktivnosti,
    la.opis
FROM LogAktivnosti la
LEFT JOIN Korisnici k ON la.korisnik_id = k.korisnik_id
ORDER BY la.datuma DESC
LIMIT 10;
