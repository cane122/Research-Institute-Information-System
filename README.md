# Research Institute Information System

Informacioni sistem za istraživačko razvojni institut kreiran sa Wails framework-om (Go backend) i Vue.js frontend-om.

## � Funkcionalnosti

### ✅ Implementirane funkcionalnosti:

1. **Autentifikacija i autorizacija**
   - Login sistem sa demo kredencijalima
   - Uloge korisnika (Admin/Korisnik)
   - Session management

2. **Dashboard**
   - Pregled statistika (projekti, dokumenti, zadaci, korisnici)
   - Poslednje aktivnosti
   - Brze akcije
   - Status projekata sa progress bar-om

3. **Upravljanje projektima**
   - CRUD operacije za projekte
   - Filteriranje i pretraga
   - Grid prikaz sa detaljnim informacijama
   - Modal forme za kreiranje/editovanje

4. **Kanban board za zadatke**
   - Drag & drop funkcionalnost
   - 4 kolone statusa (Za rad, U toku, Na proveri, Završeno)
   - Lista i Kanban prikaz
   - Filteriranje po projektu, prioritetu, korisniku
   - CRUD operacije za zadatke
   - Uslovi faze (phase conditions)
   - Zahtevi za promenu faze (phase change requests)
   - **Lista zahteva** - stranica za rukovodioce projekata

5. **Upravljanje dokumenata**
   - Upload funkcionalnost
   - Grid i lista prikaz
   - Kategorije i tipovi dokumenata
   - Pretraga i filteriranje
   - Akcije (preuzmi, podeli, obriši)

6. **Administracija korisnika** (samo za admin)
   - CRUD operacije za korisnike
   - Bulk akcije (aktivacija, deaktivacija, brisanje)
   - Statistike korisnika
   - Reset lozinke

### 🎨 UI/UX Karakteristike:

- **Responsivan dizajn** - prilagođava se svim uređajima
- **Moderni UI** - inspirisan wireframe-ovima
- **Intuitivna navigacija** - sidebar sa ikonama
- **Smooth animacije** - hover efekti i tranzicije
- **Kanban drag & drop** - intuitivno upravljanje zadacima
- **Modal dialog-zi** - za kreiranje/editovanje
- **Loading states** - indikatori učitavanja

## 🚀 BRZO POKRETANJE - STARTUP SKRIPTE

### Windows (jedan klik):
```batch
# Automatski pokreće bazu, frontend i backend
./start-project.bat

# PowerShell verzija (naprednije opcije)  
powershell -ExecutionPolicy Bypass -File start-project.ps1

# Zaustavljanje svih servisa
./stop-project.bat
```

### Linux/macOS:
```bash
# Dodeli dozvole i pokreni
chmod +x start-project.sh
./start-project.sh
```

### NPM komande:
```bash
npm run start      # Windows batch
npm run start:ps   # PowerShell  
npm run stop       # Zaustavi sve
npm run build      # Build ceo projekat
```

### 📋 Šta rade startup skripte:

1. **Proveravaju i pokreću PostgreSQL** servis
2. **Kreiraju bazu** `research_institute` (ako ne postoji)  
3. **Učitavaju dummy podatke** automatski
4. **Pokreću frontend** server (Vue.js na portu 5173)
5. **Pokreću backend** (Wails development mode)
6. **Otvaraju aplikaciju** u browseru
7. **Omogućavaju čisto zaustavljanje** svih servisa

---

## 🔧 Tehnologije

### Frontend:
- **Vue.js 3** - Composition API
- **Vue Router 4** - Rutiranje
- **Pinia** - State management
- **Vite** - Build tool
- **CSS3** - Custom styling bez framework-a

### Backend:
- **Go** - Server-side logika
- **Wails v2** - Desktop aplikacija framework
- **PostgreSQL** - Baza podataka (spremna za integraciju)

## 📦 Instaliranje

### Preduslovi:

#### 1. Node.js instalacija
- Preuzmite i instalirajte Node.js 18+ sa [zvaničnog sajta](https://nodejs.org/)
- Proverite instalaciju: `node --version` i `npm --version`

#### 2. Go instalacija
- **Windows:**
  1. Idite na [Go downloads](https://golang.org/dl/)
  2. Preuzmite Windows installer (go1.21.x.windows-amd64.msi)
  3. Pokrenite installer i pratite instrukcije
  4. Dodajte `C:\Go\bin` u PATH environment variable
  5. Restartujte Command Prompt/PowerShell
  6. Proverite instalaciju: `go version`

- **macOS:**
  ```bash
  # Sa Homebrew
  brew install go
  
  # Ili preuzmite installer sa golang.org
  ```

- **Linux:**
  ```bash
  # Ubuntu/Debian
  sudo apt update
  sudo apt install golang-go
  
  # CentOS/RHEL/Fedora
  sudo dnf install golang
  ```

#### 3. Wails CLI instalacija
```bash
# Instalirajte Wails v2
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Proverite instalaciju
wails doctor

# Inicijalizujte Wails dependencies
wails init
```

#### 4. PostgreSQL instalacija
- **Windows:** Preuzmite sa [postgresql.org](https://www.postgresql.org/download/windows/)
- **macOS:** `brew install postgresql`
- **Linux:** `sudo apt install postgresql postgresql-contrib`

### Frontend setup:

```bash
cd frontend
npm install
```

### Development:

```bash
# Pokretanje frontend-a u development modu
cd frontend
npm run dev
```

Aplikacija će biti dostupna na: http://localhost:5173/

### Production build:

```bash
cd frontend
npm run build
```

### Wails Desktop Aplikacija:

#### Development mode (preporučeno):
```bash
# Pozicionirajte se u root direktorijum
cd "c:\Users\cane\Downloads\iis\Research Institute Information System"

# Pokrenite Wails development server
wails dev
```

#### Production build:
```bash
# Windows executable
wails build -platform windows/amd64

# Cross-platform builds
wails build -platform windows/amd64,darwin/amd64,linux/amd64

# Rezultat će biti u build/bin/ direktorijumu
```

#### Pokretanje gotove aplikacije:
```bash
# Direktno pokretanje
./research-institute-system.exe

# Ili iz build direktorijuma
./build/bin/research-institute-system.exe
```

#### Napomene za Wails:
- Frontend mora biti build-ovan pre Wails kompajliranja
- Wails će automatski embed-ovati frontend fajlove
- Desktop aplikacija neće zavisiti od eksternog web server-a

## 🔐 Demo kredencijali

**Administrator:**
- Korisničko ime: `admin`
- Lozinka: `admin`

**Korisnik:**
- Korisničko ime: `user`
- Lozinka: `user`

## 📱 Responsive Design

Aplikacija je potpuno responzivna i prilagođava se:
- **Desktop** (1200px+) - Pun sidebar i grid layout
- **Tablet** (768px - 1199px) - Kompaktni sidebar
- **Mobile** (< 768px) - Collapsed sidebar, stack layout

## 📁 Struktura frontend projekta

```
frontend/
├── src/
│   ├── components/
│   │   └── Layout.vue          # Glavni layout sa sidebar-om
│   ├── views/
│   │   ├── Login.vue           # Login stranica
│   │   ├── Dashboard.vue       # Dashboard sa statistikama
│   │   ├── Projects.vue        # Upravljanje projektima
│   │   ├── Tasks.vue           # Kanban board za zadatke
│   │   ├── Documents.vue       # Upravljanje dokumentima
│   │   └── Users.vue           # Admin korisnici
│   ├── stores/
│   │   └── auth.js             # Pinia auth store
│   ├── router.js               # Vue Router konfiguracija
│   ├── style.css              # Globalni stilovi
│   ├── App.vue                # Root komponenta
│   └── main.js                # Entry point
├── index.html                 # HTML template
├── package.json
└── vite.config.js
```

## 🎯 Demo funkcionalnosti

Kada pokrenete aplikaciju, možete testirati:

1. **Login** - koristite demo kredencijale
2. **Dashboard** - pregledajte statistike i aktivnosti
3. **Projekti** - kreirajte, editujte projekte
4. **Zadaci** - testirajte Kanban board sa drag & drop
5. **Dokumenti** - upload dokumenata (mock funkcionalnost)
6. **Korisnici** - admin funkcionalnosti (samo sa admin nalogom)

## 🚧 Budući razvoj

### Planirana proširenja:

1. **Backend integracija**
   - Povezivanje sa Go backend-om
   - Stvarne API pozive umesto mock podataka
   - Autentifikacija sa JWT tokenima

2. **Dodatne funkcionalnosti**
   - Real-time notifikacije
   - Advanced search i filtering
   - Export/Import podataka
   - File viewer za dokumente
   - Kalendar događaja

3. **Optimizacije**
   - Code splitting
   - Lazy loading komponenti
   - Performance monitoring
   - PWA karakteristike

---

## 🗄️ Podešavanje Baze Podataka (Original Backend)

### 1. Instalacija PostgreSQL
- Preuzmite i instalirajte PostgreSQL sa [zvaničnog sajta](https://www.postgresql.org/download/)
- Zapamtite lozinku koju postavite za `postgres` korisnika

### 2. Kreiranje Baze Podataka
```sql
CREATE DATABASE research_institute;
```

### 3. Kreiranje Šeme i Učitavanje Dummy Podataka

#### Automatsko učitavanje (Windows):
```batch
# Pokrenite batch fajl za učitavanje dummy podataka
./load-dummy-data.bat
```

#### Manuelno učitavanje:

**Korak 1: Kreiranje šeme**
```bash
# Windows (PowerShell/CMD)
psql -U postgres -d research_institute -f database/schema.sql

# macOS/Linux
psql -U postgres -d research_institute -f database/schema.sql
```

**Korak 2: Učitavanje dummy podataka**
```bash
# Windows
psql -U postgres -d research_institute -f database/dummy_data.sql

# macOS/Linux  
psql -U postgres -d research_institute -f database/dummy_data.sql
```

**Korak 3: Verifikacija podataka**
```bash
# Pokretanje skripte za pregled podataka
psql -U postgres -d research_institute -f database/view_data.sql

# Ili manuelno povezivanje i pregled
psql -U postgres -d research_institute
\dt  -- Lista tabela
SELECT COUNT(*) FROM users;    -- Broj korisnika
SELECT COUNT(*) FROM projects; -- Broj projekata
SELECT COUNT(*) FROM tasks;    -- Broj zadataka
\q   -- Izlaz iz psql
```

#### Dummy podaci uključuju:

1. **Korisnici (10 korisnika)**
   - 2 administratora (admin, marko.petrovic)
   - 8 običnih korisnika
   - Različite uloge: Senior Developer, Project Manager, QA Engineer, itd.

2. **Projekti (8 projekata)**
   - Različiti statusi: Active, Planning, Completed, On Hold
   - Različiti prioriteti: High, Medium, Low
   - Dodeljeni project manageri

3. **Zadaci (25 zadataka)**
   - Raspoređeni po projektima
   - 4 statusa: Todo, In Progress, Review, Done  
   - Različiti prioriteti i dodeljeni korisnici

4. **Dokumenti (15 dokumenata)**
   - Različiti tipovi: Specifikacija, Dizajn, Test Plan
   - Povezani sa projektima
   - Različite verzije

#### Konfiguracija baze podataka:

Ako koristite različite kredencijale za PostgreSQL, uredite fajlove:
- `database/schema.sql` - za kreiranje strukture
- `database/dummy_data.sql` - za test podatke
- `load-dummy-data.bat` - za automatsko učitavanje (Windows)
- `update_db.ps1` - PowerShell skripta za update

**Primer konfiguracije:**
```bash
# Ako koristite drugačiji port ili host
psql -h localhost -p 5432 -U your_username -d research_institute -f database/schema.sql
```

### 4. Wails Development
```bash
# Instalacija Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Development mode
wails dev

# Production build
wails build
```

## 📁 Kompletna Struktura Projekta

```
Research Institute Information System/
├── backend/                      # Go backend kod
│   ├── models/                   # Data modeli
│   ├── repositories/             # Database repository sloj
│   └── services/                 # Business logika
├── build/                        # Build output
├── database/                     # Database šema i migracije
├── frontend/                     # Vue.js frontend aplikacija
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/
│   │   └── router.js
│   ├── dist/                     # Production build
│   ├── package.json
│   └── vite.config.js
├── wireframes/                   # UI wireframes
├── go.mod                       # Go dependencies
├── main.go                      # Glavna aplikacija
├── wails.json                   # Wails konfiguracija
├── DEPLOYMENT.md                # Deployment instrukcije
└── README.md                    # Ovaj fajl
```

## 🎨 Wireframes

Projekat sadrži kompletne wireframes za sve module:

- **Login** - `wireframes/01-login.html`
- **Dashboard** - `wireframes/02-dashboard.html`
- **Projekti** - `wireframes/03-projects.html`
- **Dokumenti** - `wireframes/04-documents.html`
- **Kanban Taskovi** - `wireframes/05-tasks-kanban.html`
- **Administracija Korisnika** - `wireframes/06-users-admin.html`

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

## ⚡ Status Implementacije

### ✅ Kompletno Implementirano
- **Database šema** - PostgreSQL tabele i relacije
- **Go backend** - modeli, repositories, servisi
- **Vue.js frontend** - kompletne funkcionalnosti
- **Wails desktop** - aplikacija framework
- **UI wireframes** - svi ekrani dizajnirani
- **Authentication** - login sistem sa ulogama
- **CRUD operacije** - za sve module
- **Responsive design** - mobile-first pristup

### 🔄 Backend Integracija (U toku)
- REST API endpoints
- Povezivanje frontend-a sa backend-om
- File upload funkcionalnost
- Real-time notifikacije

### � Budući Razvoj
- Email notifikacije
- Advanced reporting
- Export/Import podataka
- PWA karakteristike

## 🛠️ Troubleshooting

### Frontend Problemi

**Greška:** `npm run dev` ne radi

**Rešenje:**
```bash
cd frontend
npm install
npm run dev
```

**Greška:** Vue komponente se ne učitavaju

**Rešenje:**
1. Proverite da li je Layout.vue komponenta kreirana
2. Proverite importovanje u router.js
3. Restartujte dev server

### Backend Problemi

**Greška:** `Failed to connect to database`

**Rešenje:**
1. Proverite da li je PostgreSQL servis pokrenut
2. Proverite da li baza `research_institute` postoji
3. Verifikujte username i password u .env fajlu
4. Proverite da li port 5432 nije blokiran

**Greška:** `wails: command not found`

**Rešenje:**
```bash
# Instaliraj Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Dodaj Go bin direktorijum u PATH (Windows)
# Dodajte u System Environment Variables:
# %USERPROFILE%\go\bin

# Ili postavite GOPATH manuelno
go env GOPATH
# Dodajte %GOPATH%\bin u PATH
```

**Greška:** `Go installation issues`

**Detaljno rešenje za Windows:**
1. Preuzmite Go sa https://golang.org/dl/
2. Pokrenite installer kao Administrator  
3. Dodajte u PATH environment variables:
   - `C:\Go\bin` (za standardnu instalaciju)
   - `%USERPROFILE%\go\bin` (za user packages)
4. Otvorite novi Command Prompt/PowerShell
5. Test: `go version`

**Verifikacija Go instalacije:**
```bash
# Proverite verziju
go version

# Proverite Go environment
go env

# Proverite GOPATH i GOROOT
go env GOPATH
go env GOROOT

# Test jednostavnog programa
echo 'package main; import "fmt"; func main() { fmt.Println("Hello Go!") }' > test.go
go run test.go
del test.go
```

**Wails diagnostika:**
```bash
# Proverite Wails instalaciju
wails version

# Dijagnoza sistema
wails doctor

# Kreiranje test Wails app
wails init -n testapp -t vanilla
cd testapp
wails dev
```

### Build Problemi

**Greška:** Go compilation errors

**Rešenje:**
```bash
go mod tidy
go mod download
```

**Greška:** Frontend build fails

**Rešenje:**
```bash
cd frontend
npm install
npm run build
```

## 📊 Demo i Testing

**Live Demo:** http://localhost:5173/ (nakon `npm run dev`)

**Demo kredencijali:**
- **Administrator:** admin / admin
- **Korisnik:** user / user

**Test scenario:**
1. Login sa admin nalogom
2. Testirajte Dashboard statistike
3. Kreirajte novi projekat
4. Dodajte zadatke u Kanban board
5. Upload dokument (mock)
6. Upravljajte korisnicima (admin only)

## 📄 Dokumentacija

- **README.md** - Osnovne instrukcije (ovaj fajl)
- **DEPLOYMENT.md** - Deployment instrukcije
- **database/schema.sql** - Database šema
- **wireframes/** - UI dizajn wireframes
- **.env.example** - Environment varijable template

## 📝 Licenca

MIT License - videti LICENSE fajl za detalje.

## 👥 Kontakt i Podrška

Za pitanja o projektu:
- **GitHub Issues:** https://github.com/cane122/Research-Institute-Information-System/issues
- **Email:** [kontakt email]

## 🔄 Change Log

### v2.0.0 (Septembar 2025) - CURRENT
- ✅ **Vue.js frontend** - kompletna implementacija
- ✅ **Responsive design** - desktop/tablet/mobile
- ✅ **Kanban board** - drag & drop funkcionalnost
- ✅ **Authentication** - login sa ulogama
- ✅ **CRUD operacije** - svi moduli
- ✅ **Modern UI/UX** - animacije i efekti

### v1.0.0 (Početak 2025)
- Početna verzija sa database šemom
- Go backend struktura
- Wails setup
- UI wireframes

---

**Napomena:** Aplikacija sada ima potpuno funkcionalan Vue.js frontend koji radi sa mock podacima. Za potpunu funkcionalnost sa pravom bazom podataka, potrebno je integrisati backend API-je.

**🚀 Za brzo pokretanje:**
```bash
cd frontend && npm install && npm run dev
```

**GitHub Repository:** https://github.com/cane122/Research-Institute-Information-System
