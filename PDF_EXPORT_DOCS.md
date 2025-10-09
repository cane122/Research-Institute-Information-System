# PDF Export Dokumentacija

## Pregled

PDF export funkcionalnost generiše profesionalne, detaljne izveštaje o analitici sistema sa lepim dizajnom i kompletnom podelom.

## Sadržaj PDF Izveštaja

### 1. **Header (Zaglavlje)**
- Gradient pozadina u boji sistema
- Naslov: "📊 Analitički Izveštaj"
- Datum generisanja izveštaja

### 2. **Executive Summary (Izvršni Rezime)**
- Opis ukupnog stanja sistema
- Sažetak ključnih metrika
- Okvir sa bord

### 3. **Key Metrics (Ključne Metrike)**
Četiri obojene kartice sa:
- **Ukupno Dokumenata** (ljubičasta)
- **Broj Pregleda** (pink)
- **Izbrisanih Dokumenata** (narandžasta)
- **Novih Dokumenata** (teal)

### 4. **Activity Statistics (Statistika Aktivnosti)**
Tabela sa podacima:
- Tip aktivnosti
- Ukupan broj
- Broj danas
- Broj ove nedelje
- Broj ovog meseca

Tipovi aktivnosti:
- Upload
- Izmena
- Brisanje
- Pregled
- Prijava
- Odjava

### 5. **Top Contributors (Top Autori)**
Tabela sa rangovima:
- Rang (1-5)
- Ime korisnika
- Broj uploadovanih dokumenata

### 6. **Documents by Type (Dokumenti po Tipu)**
Tabela sa:
- Tip dokumenta
- Broj dokumenata
- Procenat od ukupnog broja

### 7. **Recent Activity (Skornje Aktivnosti)**
Detaljna tabela sa poslednjim 50 aktivnosti:
- Tip aktivnosti
- Korisnik koji je izvršio
- Entitet (dokument/korisnik)
- Datum i vreme

### 8. **Footer (Podnožje)**
Na svakoj strani:
- Linija separacije
- Naziv sistema: "Research Institute Information System"
- Broj strane (npr. "Strana 1 od 3")

## Dizajn Karakteristike

### Boje
- **Primarni gradient**: #667eea → #764ba2 (ljubičasta)
- **Sekundarni**: #f093fb → #f5576c (pink)
- **Akcent 1**: #fa709a → #fee140 (narandžasta)
- **Akcent 2**: #30cfd0 → #330867 (teal)
- **Tekst**: #1e293b (tamno siva)
- **Sekundarni tekst**: #94a3b8 (svetlo siva)

### Tipografija
- **Naslov**: 24pt, bold, bela boja
- **Sekcije**: 16pt, bold, tamna
- **Tabele (header)**: 10pt, bold, bela
- **Tabele (sadržaj)**: 9pt, normalno
- **Footer**: 8pt, siva

### Layout
- **Margine**: 20mm sa svih strana
- **Razmak između sekcija**: 15mm
- **Tabele**: Striped theme (naizmjenične boje redova)
- **Okviri**: Zaobljeni uglovi (3mm radius)

## Kako Koristiti

### U Aplikaciji

1. Otvorite Document Management stranicu
2. Kliknite na "📊 Analytics" dugme
3. Na analitičkoj stranici, kliknite "📄 Export PDF"
4. PDF će se automatski preuzeti sa nazivom formata:
   ```
   Analytics_Report_YYYY-MM-DD.pdf
   ```

### Programski

```javascript
// Import funkcije
import { exportPDF } from './DocumentAnalytics.vue'

// Poziv
await exportPDF()
```

## Tehnologija

### Biblioteke
- **jsPDF**: Osnovna PDF generacija
- **jspdf-autotable**: Tabele sa automatskim prelomima strana

### Instalacija
```bash
npm install jspdf jspdf-autotable
```

### Import u komponenti
```javascript
import jsPDF from 'jspdf'
import 'jspdf-autotable'
```

## Primer PDF Strukture

```
┌─────────────────────────────────────────────┐
│  HEADER (gradient pozadina)                 │
│  📊 Analitički Izveštaj                     │
│  Generisano: 9. oktobar 2025                │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  IZVRŠNI REZIME                             │
│  ┌─────────────────────────────────────┐   │
│  │ Ovaj izveštaj pruža...              │   │
│  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  KLJUČNE METRIKE                            │
│  ┌───┐ ┌───┐ ┌───┐ ┌───┐                  │
│  │245│ │2642│ │32│ │56│                   │
│  └───┘ └───┘ └───┘ └───┘                  │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  STATISTIKA AKTIVNOSTI                      │
│  ┌───────────────────────────────────────┐ │
│  │ Tip    │ Ukupno │ Danas │ Nedelja  │ │ │
│  ├───────────────────────────────────────┤ │
│  │ Upload │   56   │   3   │    12    │ │ │
│  │ View   │  2642  │  145  │   892    │ │ │
│  └───────────────────────────────────────┘ │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  TOP AUTORI                                 │
│  ┌───────────────────────────────────────┐ │
│  │ Rang │ Korisnik │ Dokumenti        │ │ │
│  ├───────────────────────────────────────┤ │
│  │  1   │ John Doe │      42          │ │ │
│  └───────────────────────────────────────┘ │
└─────────────────────────────────────────────┘

[Nova strana za Recent Activity...]

┌─────────────────────────────────────────────┐
│  Research Institute Information System      │
│                           Strana 1 od 3     │
└─────────────────────────────────────────────┘
```

## Specijalne Funkcionalnosti

### Auto Page Break
- Automatski prelom strana kada sadržaj prelazi granicu
- Svaka sekcija proverava dostupan prostor
- Nova strana se dodaje po potrebi

### Responsive Tables
- Automatsko prilagođavanje širine kolona
- Automatski prelom na novu stranu
- Alternativne boje redova za čitljivost

### Date Formatting
- Srpski format datuma: "DD. Month YYYY"
- 24-časovni format vremena: "HH:MM"

### Color Coding
- Svaka sekcija ima svoju boju
- Vizuelna hijerarhija kroz tipografiju
- Konzistentna paleta boja

## Testiranje

### Test Scenario 1: Prazan Dataset
```javascript
// Rezultat: PDF sa "0" vrednostima i "Nema podataka" porukama
```

### Test Scenario 2: Mali Dataset (< 50 aktivnosti)
```javascript
// Rezultat: PDF sa 1-2 strane
```

### Test Scenario 3: Veliki Dataset (> 1000 aktivnosti)
```javascript
// Rezultat: PDF sa 3+ strana, prikaz prvih 50 aktivnosti
```

## Troubleshooting

### PDF se ne generiše
- Proverite da li su instaliran jsPDF i jspdf-autotable
- Proverite console za greške
- Osvežite cache: `npm run build`

### Nepotpuni podaci u PDF-u
- Proverite da li backend vraća podatke
- Proverite network tab u DevTools
- Proverite da li su views kreirani u bazi

### Greške sa fontovima
- jsPDF koristi standardne fontove
- Ne zahteva dodatne font fajlove

## Budući Poboljšanja

1. **Grafici i Charts**
   - Line chart za trend dokumenata
   - Pie chart za distribuciju tipova
   - Bar chart za aktivnosti po korisnicima

2. **Customizacija**
   - Izbor perioda (poslednji mesec, 3 meseca, godina)
   - Filter po tipu aktivnosti
   - Filter po korisniku

3. **Export Opcije**
   - Excel export (.xlsx)
   - CSV export
   - JSON export

4. **Email Opcije**
   - Slanje PDF-a na email
   - Zakazano generisanje izveštaja
   - Auto-email sa statistikom

5. **Dodatne Sekcije**
   - Projekat statistika
   - Korisnik performanse
   - Storage analytics
   - Error logs

## Zaključak

PDF export funkcionalnost pruža profesionalan i detaljan pregled aktivnosti sistema sa:
- ✅ Lepim dizajnom sa bojama
- ✅ Jasnom podelom sekcija
- ✅ Kompletnom statistikom
- ✅ Detaljnim aktivnostima
- ✅ Automatskim formatiranjem
- ✅ Multi-page podrškom
