# Research Institute Information System

Informacioni sistem za istraživačko-razvojni institut razvijen pomoću Wails framework-a (Go + HTML/CSS/JavaScript).

## 📋 Pregled Projekta

Ovaj sistem omogućava:
- **Upravljanje korisnicima i ulogama** - administracija korisnika, dodela uloga
- **Upravljanje projektima i zadacima** - kreiranje, praćenje i upravljanje projektima
- **Upravljanje dokumentima** - upload, verzioniranje i organizacija dokumenata
- **Izveštavanje i logovanje** - praćenje aktivnosti i generiranje izveštaja

## 🚀 Brza Instalacija

### Preduslovi

Potrebno je da imate instaliran:
- [Go](https://golang.org/dl/) (verzija 1.19 ili novija)
- [Node.js](https://nodejs.org/) (verzija 16 ili novija)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- [PostgreSQL](https://www.postgresql.org/download/)

### Instalacija Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Kloniranje Projekta

```bash
git clone https://github.com/cane122/Research-Institute-Information-System.git
cd "Research Institute Information System"
```

## 🗄️ Podešavanje Baze Podataka

### 1. Instalacija PostgreSQL

- Preuzmite i instalirajte PostgreSQL sa [zvaničnog sajta](https://www.postgresql.org/download/windows/)
- Zapamtite lozinku koju postavite za `postgres` korisnika

### 2. Kreiranje Baze Podataka

Otvorite PostgreSQL command line (psql) ili pgAdmin i izvršite:

```sql
CREATE DATABASE research_institute;
```

### 3. Kreiranje Šeme

Navigirajte do direktorijuma projekta i izvršite:

```bash
psql -U postgres -d research_institute -f database/schema.sql
```

Ili kopirajte sadržaj `database/schema.sql` datoteke u pgAdmin i izvršite SQL komande.

### 4. Konfiguracija Konekcije (Opciono)

Možete podesiti sledeće environment varijable:

```bash
# Windows Command Prompt
set DB_HOST=localhost:5432
set DB_USER=postgres
set DB_PASSWORD=vasa_lozinka
set DB_NAME=research_institute

# Windows PowerShell
$env:DB_HOST="localhost:5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="vasa_lozinka"
$env:DB_NAME="research_institute"
```

**Default konfiguracija:**
- Host: `localhost:5432`
- User: `postgres`
- Password: `password`
- Database: `research_institute`

## 🔧 Build i Pokretanje

### Development Mode

```bash
wails dev
```

### Production Build

```bash
wails build
```

Izvršna datoteka će biti kreirana u `build/bin/research-institute-system.exe`

### Pokretanje Aplikacije

```bash
# Direktno pokretanje
./build/bin/research-institute-system.exe

# Ili dvojklik na .exe datoteku
```

## 📁 Struktura Projekta

```
Research Institute Information System/
├── app/                          # Wails app konfiguracija
├── backend/                      # Go backend kod
│   ├── models/                   # Data modeli
│   ├── repositories/             # Database repository sloj
│   └── services/                 # Business logika
├── build/                        # Build output
├── database/                     # Database šema i migracije
├── frontend/                     # Frontend assets
├── wireframes/                   # UI wireframes
├── go.mod                       # Go dependencies
├── main.go                      # Glavna aplikacija
└── wails.json                   # Wails konfiguracija
```

## 🎨 Wireframes

Projekat sadrži kompletne wireframes za sve module:

- **Login** - `wireframes/login.html`
- **Dashboard** - `wireframes/dashboard.html`
- **Projekti** - `wireframes/projects.html`
- **Dokumenti** - `wireframes/documents.html`
- **Kanban Taskovi** - `wireframes/tasks-kanban.html`
- **Administracija Korisnika** - `wireframes/user-admin.html`

Otvorite bilo koji HTML fajl u browser-u za pregled dizajna.

## 🗃️ Database Šema

Sistem koristi PostgreSQL bazu sa sledećim modulima:

### Modul 1: Upravljanje Korisnicima
- `Uloge` - definisanje korisničkih uloga
- `Korisnici` - informacije o korisnicima sistema

### Modul 2: Upravljanje Projektima
- `RadniTokovi` - definisanje workflow-a
- `Faze` - faze unutar workflow-a
- `Projekti` - osnovne informacije o projektima
- `ClanoviProjekta` - članovi projektnih timova
- `Zadaci` - zadaci unutar projekata

### Modul 3: Upravljanje Dokumentima
- `Folderi` - organizacija dokumenata
- `Dokumenti` - metadata dokumenata
- `VerzijeDokumenata` - verzioniranje
- `MetaPodaci` - dodatni metadata
- `Tagovi` - tagovanje sistema

### Modul 4: Logovanje i Izveštaji
- `LogAktivnosti` - praćenje korisničkih aktivnosti

Kompletna šema se nalazi u `database/schema.sql`.

## ⚡ Funkcionalnosti

### Trenutno Implementirano
- ✅ Kompletna database šema
- ✅ Go backend modeli i strukture
- ✅ Osnovni repository pattern
- ✅ Authentication servis
- ✅ Wails desktop aplikacija
- ✅ UI wireframes za sve module

### Planirano za Razvoj
- 🔄 Frontend implementacija (React/Vue/Vanilla JS)
- 🔄 REST API endpoints
- 🔄 File upload funkcionalnost
- 🔄 Reporting sistem
- 🔄 Email notifikacije

## 🛠️ Troubleshooting

### Problem sa Bazom Podataka

**Greška:** `Failed to connect to database`

**Rešenje:**
1. Proverite da li je PostgreSQL servis pokrenut
2. Proverite da li baza `research_institute` postoji
3. Verifikujte username i password
4. Proverite da li port 5432 nije blokiran
5. Proverite Windows Firewall postavke

### Build Problemi

**Greška:** `wails: command not found`

**Rešenje:**
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

**Greška:** Go compilation errors

**Rešenje:**
```bash
go mod tidy
go mod download
```

## 📊 Screenshots

![Dashboard Wireframe](wireframes/dashboard.html)
![Projects Management](wireframes/projects.html)
![Document Management](wireframes/documents.html)

## 📝 Licenca

Ovaj projekat je razvijen za potrebe istraživačko-razvojnog instituta.

## 👥 Doprinosi

Za pitanja i predloge kontaktirajte autora projekta preko GitHub-a.

## 🔄 Change Log

### v1.0.0 (Septembar 2025)
- Početna verzija sa kompletnom database šemom
- Osnovni Go backend struktura
- Wails desktop aplikacija setup
- UI wireframes za sve module
- Build sistem i deployment instrukcije

---

**Napomena:** Aplikacija će raditi i bez database konekcije, ali većina funkcionalnosti neće biti dostupna. Pratite instrukcije za podešavanje baze podataka za potpunu funkcionalnost.

**GitHub Repository:** https://github.com/cane122/Research-Institute-Information-System
