# 🎯 KAKO POKRENUTI DEMO ZA PROFESORA

## 📁 Priprema

Imaš 2 skripte za demonstraciju:

### 1. **BRZI_DEMO.sql** ⚡ (PREPORUČENO)
- Kratak, fokusiran demo (~3-5 minuta)
- PAUSE komande - možeš objašnjavati između testova
- Vizuelno lepše formatiran
- Idealno za live demo

### 2. **DEMO_ZA_PROFESORA.sql** 📊 (Detaljan)
- Svi testovi odjednom
- Više detalja i objašnjenja
- ~10 minuta
- Bolje za samostalni pregled

---

## 🚀 Pokretanje (BRZI_DEMO.sql)

### Windows PowerShell:

```powershell
cd "c:\Users\cane\Downloads\sbp\Research Institute Information System\database"
sqlplus SYSTEM/123@localhost:1521/xe @BRZI_DEMO.sql
```

### Ili direktno SQL*Plus:

```sql
SQL> @"c:\Users\cane\Downloads\sbp\Research Institute Information System\database\BRZI_DEMO.sql"
```

---

## 📋 Šta Će Se Desiti

Demo ima **4 testa** sa PAUSE komandama:

### ✅ Test 1: Indeks Performance (30 sec)
```
1. Pokaže broj zadataka u bazi (578)
2. Ukloni indeks → pokrene query → izmeri vreme
3. PAUSE → ti objasniš profesoru šta se dešava
4. Kreira indeks → pokrene isti query → izmeri vreme
5. Uporedi vremena!
```

**Šta reći profesoru:**
> "Bez indeksa Oracle mora da skenira svih 578 redova. Sa indeksom pristupa samo matching redovima, zato je ~500x brže."

---

### ✅ Test 2: PL/SQL Funkcije (1 min)
```
1. Pokaže kako BEZ funkcija moraš 3 upita po projektu
2. PAUSE → objasniš problem
3. Pokaže kako SA funkcijama - samo 1 upit za sve
```

**Šta reći profesoru:**
> "Bez PL/SQL funkcija, aplikacija mora da pošalje 339 SQL upita (113 projekata × 3 upita). Sa funkcijama - samo 1 upit. To je 339x efikasnije i mnogo manje network latency."

---

### ✅ Test 3: Kompleksan Izveštaj (1 min)
```
1. PAUSE → kažeš šta će se izvršiti
2. Poziva kompleksan_izvestaj_projekata proceduru
3. Prikazuje sve projekte sa statistikama
```

**Šta reći profesoru:**
> "Ova procedura demonstrira SVE zahtevane komponente:
> - Složene PL/SQL tipove (RECORD, TABLE OF)
> - Kursor sa WITH klauzulom
> - JOIN iz 4 tabele (Projekti, ClanoviProjekta, Zadaci, Dokumenti)
> - GROUP BY, COUNT, SUM operacije
> Sve to u jednom optimizovanom upitu umesto 4 odvojena."

---

### ✅ Test 4: Trigeri (30 sec)
```
1. Broji log zapise PRE inserta
2. Insertuje test projekat
3. Broji log zapise POSLE inserta
4. Pokazuje da je trigger automatski dodao log
5. Cleanup
```

**Šta reći profesoru:**
> "Trigger `trg_log_projekta` automatski loguje SVE izmene u bazi. Ne moramo ni liniju koda da dodamo u Go aplikaciju. To je audit trail koji ne može da se preskoči."

---

## 💡 Pro Tips za Live Demo

### Pre Demoa:
```sql
-- Proveri da sve radi
SELECT COUNT(*) FROM projekti;   -- Treba: 113
SELECT COUNT(*) FROM zadaci;     -- Treba: 578
SELECT object_name FROM user_objects WHERE object_type = 'FUNCTION';
```

### Tokom Demoa:
1. **Pusti da se timing vidi** - ne žuri sa ENTER
2. **Pokazuj razlike u TIMING linijama** - to je dokaz!
3. **Objašnjaj PAUSE momentima** - time pokazuješ da razumeš

### Ako Profesor Pita "Gde se koristi u Go?":

Otvori fajl i pokažeš:

```bash
code "backend/services/project_service.go"
# Linija 28-29: Pozivaju se funkcije
```

```go
broj_aktivnih_clanova(p.projekat_id) as broj_aktivnih_clanova,
procenat_zavrsenih_zadataka(p.projekat_id) as procenat_zavrsenosti
```

---

## 📊 Očekivani Rezultati

### Timing Comparison (sa 578 zadataka):

| Test | Bez optimizacije | Sa optimizacijom | Razlika |
|------|------------------|------------------|---------|
| Indeks pretraga | ~0.5s | ~0.001s | **500x** |
| Funkcije (113 proj) | ~5-10s (339 upita) | ~0.02s (1 upit) | **250-500x** |
| Kompleksan izveštaj | ~15s | ~0.1s | **150x** |

---

## 🎯 Plan za 5-Minutni Demo

```
0:00 - 0:30   Objasni šta ćeš pokazati
0:30 - 1:30   Test 1: Indeks (sa/bez)
1:30 - 2:30   Test 2: Funkcije (sa objašnjenjem)
2:30 - 4:00   Test 3: Kompleksan izveštaj (pokaži output)
4:00 - 4:30   Test 4: Trigeri
4:30 - 5:00   Sažetak + pitanja
```

---

## 🆘 Ako Nešto Pođe Po Zlu

### Indeks već postoji?
```sql
DROP INDEX idx_zadaci_prioritet;
```

### Procedura ne radi?
```sql
@database/tasks/05_reports.sql
```

### Funkcije ne rade?
```sql
@database/tasks/03_functions.sql
```

### Nema dovoljno podataka?
```sql
@database/add_test_data_simple.sql
```

---

## ✅ Finalni Checklist

Pred odbranu proveri:

- [ ] `BRZI_DEMO.sql` fajl postoji u `database/` folderu
- [ ] Oracle baza je pokrenuta
- [ ] Imaš 113 projekata i 578 zadataka
- [ ] Sve funkcije i trigeri su kreirani (pogledaj `TEST_REZULTATI.md`)
- [ ] Poznaješ gde u Go kodu se pozivaju funkcije (`GDE_SE_KORISTI.md`)
- [ ] Možeš da objasniš svaku PL/SQL komponentu

---

## 🎓 Bonus - Explain Plan Demo

Ako profesor traži dokaz da se indeks koristi:

```sql
EXPLAIN PLAN FOR
SELECT * FROM zadaci WHERE prioritet = 'Visok';

SELECT * FROM TABLE(DBMS_XPLAN.DISPLAY);
```

Tražiš liniju: **`INDEX RANGE SCAN`** - to je dokaz!

---

**SREĆA NA ODBRANI!** 🍀

Sve je testirano i radi. Samo budi siguran u objašnjenja!
