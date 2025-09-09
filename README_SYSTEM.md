# Research Institute Information System

Informacioni sistem za istraživačko-razvojni institut sa tri glavna podsistema:

1. **Upravljanje dokumentima** - Upravljanje i organizacija dokumenata sa meta-podacima
2. **Priprema projektne dokumentacije** - Podrška za životni ciklus projektne dokumentacije  
3. **Realizacija projekata** - Kreiranje, upravljanje i praćenje projekata i zadataka

## Tehnologije

- **Backend**: Go (Wails framework)
- **Frontend**: HTML, CSS, JavaScript
- **Baza podataka**: PostgreSQL
- **Desktop aplikacija**: Wails v2

## Struktura projekta

```
├── backend/               # Go backend kod
│   ├── models/           # Data modeli
│   ├── repositories/     # Baza podataka pristup
│   └── services/         # Biznis logika
├── frontend/             # Web frontend
│   └── dist/            # Statični fajlovi (HTML, CSS, JS)
├── database/            # Baza podataka skema
├── build/               # Kompajlirani fajlovi
└── main.go             # Glavna aplikacija
```

## Instalacija i pokretanje

### Preduslovi

1. **Go** (verzija 1.21 ili novija)
   - Preuzmite sa https://golang.org/dl/
   - Dodajte `C:\Program Files\Go\bin` u PATH

2. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
   - Dodajte `%USERPROFILE%\go\bin` u PATH

3. **PostgreSQL** (opciono za punu funkcionalnost)
   - Instalirajte PostgreSQL server
   - Kreirajte bazu podataka `research_institute`
   - Pokrenite skriptu `database/schema.sql`

### Pokretanje aplikacije

#### Development mode (sa live reload):
```bash
wails dev
```

#### Production build:
```bash
wails build
```
Izvršna datoteka će biti u `build/bin/research-institute-system.exe`

#### Direktno pokretanje:
```bash
./build/bin/research-institute-system.exe
```

## Konfiguracija baze podataka

Za povezivanje sa PostgreSQL bazom, promenite connection string u `main.go`:

```go
db, err := sql.Open("postgres", "postgres://username:password@localhost/research_institute?sslmode=disable")
```

## Test korisnici (development mode)

Kada aplikacija radi bez baze podataka, možete se prijaviti sa:
- **Username**: `admin`
- **Password**: `admin`

## Funkcionalnosti po ulogama

### Administrator
- Kreiranje i upravljanje korisnicima
- Resetovanje lozinki
- Pristup svim modulima
- Upravljanje ulogama

### Rukovodilac projekta  
- Kreiranje i upravljanje projektima
- Dodela zadataka
- Definisanje radnih tokova
- Upravljanje projektnom dokumentacijom
- Analitika projekata

### Istraživač
- Rad na dodeljenim zadacima
- Upload i izmena dokumenata
- Komentarisanje zadataka
- Zahtev za promenu faze zadataka

### Organizator projekta
- Pregled projekata i napretka
- Pristup analitici
- Čitanje dokumenata

## Baza podataka

Sistem koristi PostgreSQL bazu sa sledećim glavnim tabelama:

- `Korisnici` - Korisnici sistema
- `Uloge` - Korisničke uloge
- `Projekti` - Osnovni podaci o projektima
- `Zadaci` - Zadaci unutar projekata
- `Dokumenti` - Upravljanje dokumentima
- `RadniTokovi` - Definicija radnih tokova
- `Faze` - Faze radnih tokova

Kompletnu šemu možete naći u `database/schema.sql`.

## Razvoj

### Dodavanje novih funkcionalnosti

1. **Backend** - Dodajte nova polja u modele (`backend/models/`)
2. **Repository** - Implementirajte pristup bazi (`backend/repositories/`)
3. **Service** - Dodajte biznis logiku (`backend/services/`)
4. **Main App** - Eksponajte funkcije za frontend (`main.go`)
5. **Frontend** - Implementirajte UI (`frontend/dist/`)

### Stilizovanje

Stil aplikacije je definisan u `frontend/dist/styles.css`. Koristi se responsive dizajn sa CSS Grid layoutom.

## Trenutne mogućnosti

✅ Prijava/odjava korisnika  
✅ Kreiranje korisnika (Admin)  
✅ Kreiranje projekata (Rukovodilac)  
✅ Pregled projekata  
✅ Role-based pristup  
✅ Responsive dizajn  

### U razvoju
🔄 Upravljanje zadacima  
🔄 Upload dokumenata  
🔄 Radni tokovi  
🔄 Analitika  
🔄 Pretraga dokumenata  

## Licenca

Ovaj projekat je kreiran za potrebe istraživačko-razvojnog instituta.

## Kontakt

Za pitanja i podršku kontaktirajte razvojni tim.
