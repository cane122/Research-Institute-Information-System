# 📊 Analytics & PDF Export - Quick Start Guide

## 🚀 Brzo Pokretanje

### 1. Učitavanje Test Podataka

#### Opcija A: PowerShell (Preporučeno)
```powershell
.\load-analytics-data.ps1
```

#### Opcija B: Batch File
```cmd
load-analytics-data.bat
```

#### Opcija C: Direktno PostgreSQL
```bash
psql -U postgres -d research_institute_db -f database/analytics_dummy_data.sql
```

### 2. Build Frontend
```bash
cd frontend
npm install        # Prvi put
npm run build
```

### 3. Pokretanje Aplikacije
```bash
wails dev
```

## 📈 Testiranje Analytics

### Pristup Analytics Dashboard-u
1. Prijavite se u aplikaciju
2. Idite na **Document Management**
3. Kliknite **📊 Analytics** dugme
4. Dashboard će prikazati:
   - ✅ Broj dokumenata
   - ✅ Broj pregleda  
   - ✅ Broj izbrisanih
   - ✅ Broj novih
   - ✅ Skornje aktivnosti
   - ✅ Top autore
   - ✅ Statistiku po tipu

### PDF Export
1. Na Analytics stranici
2. Kliknite **📄 Export PDF** dugme
3. PDF će se preuzeti automatski

## 🎨 PDF Sadržaj

Generisani PDF sadrži:

### Strana 1
- **Header** - Gradient zaglavlje sa naslovom
- **Izvršni Rezime** - Sažetak sistema
- **Ključne Metrike** - 4 obojene kartice
- **Statistika Aktivnosti** - Tabela sa tipovima
- **Top Autori** - Rang lista korisnika
- **Dokumenti po Tipu** - Distribucija

### Strana 2+
- **Skornje Aktivnosti** - Detaljna lista (50 aktivnosti)
- **Footer** - Broj strane na svakoj strani

## 📊 Test Podaci

Skripte generišu:
- **~2500** ukupnih aktivnosti
- **~2000** VIEW aktivnosti
- **~300** UPLOAD aktivnosti
- **~100** EDIT aktivnosti
- **~50** DELETE aktivnosti
- **~50** LOGIN/LOGOUT aktivnosti

Podaci pokrivaju period od **30 dana** unazad.

## 🔧 Konfiguracija

### PostgreSQL Connection

Ažurirajte u skriptama:

**load-analytics-data.ps1:**
```powershell
$env:PGHOST = "localhost"
$env:PGPORT = "5432"
$env:PGDATABASE = "research_institute_db"
$env:PGUSER = "postgres"
$env:PGPASSWORD = "your_password"  # ← Promenite ovo
```

**load-analytics-data.bat:**
```batch
set PGPASSWORD=your_password  # ← Promenite ovo
```

## ✅ Verifikacija

### Provera u Bazi
```sql
-- Ukupan broj aktivnosti
SELECT COUNT(*) FROM LogAktivnosti;

-- Aktivnosti po tipu
SELECT tip_aktivnosti, COUNT(*) 
FROM LogAktivnosti 
GROUP BY tip_aktivnosti 
ORDER BY COUNT(*) DESC;

-- Test statistike
SELECT * FROM StatistikaDokumenata;
SELECT * FROM StatistikaAktivnosti;
SELECT * FROM SkornjeAktivnosti LIMIT 10;
```

### Provera u Aplikaciji
1. Analytics dashboard prikazuje brojeve > 0
2. Recent Activity lista ima stavke
3. Top Contributors lista ima korisnike
4. PDF export radi bez grešaka

## 🐛 Troubleshooting

### Problem: "No data in analytics"
**Rešenje:**
```bash
# Proverite da li su podaci učitani
psql -U postgres -d research_institute_db -c "SELECT COUNT(*) FROM LogAktivnosti;"

# Ponovo učitajte podatke
.\load-analytics-data.ps1
```

### Problem: "PDF generation failed"
**Rešenje:**
```bash
# Reinstalirajte pakete
cd frontend
npm install jspdf jspdf-autotable
npm run build
```

### Problem: "Database connection failed"
**Rešenje:**
- Proverite da li PostgreSQL radi
- Proverite username/password u skripti
- Proverite da li baza postoji

## 📝 Napomene

### Dummy Podaci
- Koriste nasumične vrijednosti
- User ID: 1-5
- Document ID: 1-20
- Vremenske oznake: Poslednja 30 dana

### Performance
- 2500+ aktivnosti se učitava brzo (~2-3 sekunde)
- PDF generisanje: ~1-2 sekunde
- Dashboard loading: < 1 sekunda

### Čišćenje Podataka
```sql
-- PAŽNJA: Ovo briše SVE aktivnosti
DELETE FROM LogAktivnosti;

-- Brisanje samo test podataka (samo ako entitet_id > 20)
DELETE FROM LogAktivnosti WHERE entitet_id > 20 AND entitet_tip = 'DOKUMENT';
```

## 🎯 Sledeći Koraci

Nakon testiranja analytics:

1. **Integriši real logging**
   - ✅ Upload/Edit već integrisano
   - ✅ View/Delete već integrisano
   - ✅ Login/Logout već integrisano

2. **Testiranje sa pravim podacima**
   - Obrišite dummy podatke
   - Koristite aplikaciju normalno
   - Provjerite da li se loguje

3. **Customizacija PDF-a**
   - Dodajte logo organizacije
   - Prilagodite boje
   - Dodajte grafike

## 📚 Dokumentacija

- **ANALYTICS_IMPLEMENTATION.md** - Kompletna implementacija
- **PDF_EXPORT_DOCS.md** - Detaljna PDF dokumentacija
- **database/analytics_dummy_data.sql** - SQL skripta za test podatke

## ✨ Features

- ✅ Real-time analytics dashboard
- ✅ Activity logging za sve operacije
- ✅ Fancy PDF export sa dizajnom
- ✅ Database views za brze upite
- ✅ Responsive design
- ✅ Multi-page PDF podrška
- ✅ Automatski prelom strana
- ✅ Color-coded aktivnosti
- ✅ Top contributors ranking
- ✅ Document type distribution

Uživajte! 🎉
