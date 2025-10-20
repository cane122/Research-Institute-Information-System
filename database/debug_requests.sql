-- Check current user and their role
SELECT korisnik_id, korisnicko_ime, ime, prezime, naziv_uloge 
FROM korisnici 
WHERE korisnicko_ime = 'YOUR_USERNAME';

-- Check projects and their leaders
SELECT p.projekat_id, p.naziv_projekta, p.rukovodilac_id, 
       k.korisnicko_ime as rukovodilac_korisnicko_ime,
       k.ime || ' ' || k.prezime as rukovodilac_ime_prezime
FROM projekti p
LEFT JOIN korisnici k ON p.rukovodilac_id = k.korisnik_id
ORDER BY p.projekat_id;

-- Check phase change requests
SELECT z.zahtev_id, z.zadatak_id, z.status, z.komentar,
       zd.naziv_zadatka,
       p.projekat_id, p.naziv_projekta, p.rukovodilac_id,
       k.korisnicko_ime as podnosilac
FROM zahtevipromenefaze z
JOIN zadaci zd ON z.zadatak_id = zd.zadatak_id
JOIN projekti p ON zd.projekat_id = p.projekat_id
LEFT JOIN korisnici k ON z.podnosilac_zahteva_id = k.korisnik_id
ORDER BY z.datum_kreiranja DESC;
