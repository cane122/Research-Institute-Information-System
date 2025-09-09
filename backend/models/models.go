package models

import "time"

// =============================================================================
// Modul 1: Upravljanje Korisnicima i Ulogama
// =============================================================================

// Uloge represents user roles in the system
type Uloge struct {
	UlogaID    int    `json:"uloga_id" db:"uloga_id"`
	NazivUloge string `json:"naziv_uloge" db:"naziv_uloge"`
}

// Korisnici represents system users
type Korisnici struct {
	KorisnikID        int        `json:"korisnik_id" db:"korisnik_id"`
	KorisnickoIme     string     `json:"korisnicko_ime" db:"korisnicko_ime"`
	Email             string     `json:"email" db:"email"`
	HashSifre         string     `json:"-" db:"hash_sifre"` // Ne vraćamo hash u JSON
	Ime               *string    `json:"ime" db:"ime"`
	Prezime           *string    `json:"prezime" db:"prezime"`
	UlogaID           int        `json:"uloga_id" db:"uloga_id"`
	Status            string     `json:"status" db:"status"`
	PoslednajaPrijava *time.Time `json:"poslednja_prijava" db:"poslednja_prijava" ts_type:"string"`
	KreiranDatuma     time.Time  `json:"kreiran_datuma" db:"kreiran_datuma" ts_type:"string"`

	// Joined fields
	NazivUloge string `json:"naziv_uloge,omitempty" db:"naziv_uloge"`
}

// =============================================================================
// Modul 2: Upravljanje Projektima, Zadacima i Dokumentacijom
// =============================================================================

// RadniTokovi represents workflow definitions
type RadniTokovi struct {
	RadniTokID   int     `json:"radni_tok_id" db:"radni_tok_id"`
	Naziv        string  `json:"naziv" db:"naziv"`
	TipToka      string  `json:"tip_toka" db:"tip_toka"`
	Opis         *string `json:"opis" db:"opis"`
	DaLiJeSablon bool    `json:"da_li_je_sablon" db:"da_li_je_sablon"`
}

// Faze represents phases within workflows
type Faze struct {
	FazaID     int    `json:"faza_id" db:"faza_id"`
	RadniTokID int    `json:"radni_tok_id" db:"radni_tok_id"`
	NazivFaze  string `json:"naziv_faze" db:"naziv_faze"`
	Redosled   int    `json:"redosled" db:"redosled"`
}

// Projekti represents research projects
type Projekti struct {
	ProjekatID     int        `json:"projekat_id" db:"projekat_id"`
	NazivProjekta  string     `json:"naziv_projekta" db:"naziv_projekta"`
	Opis           *string    `json:"opis" db:"opis"`
	DatumPocetka   *time.Time `json:"datum_pocetka" db:"datum_pocetka" ts_type:"string"`
	DatumZavrsetka *time.Time `json:"datum_zavrsetka" db:"datum_zavrsetka" ts_type:"string"`
	Status         string     `json:"status" db:"status"`
	RukovodilaID   *int       `json:"rukovodilac_id" db:"rukovodilac_id"`
	RadniTokID     *int       `json:"radni_tok_id" db:"radni_tok_id"`

	// Joined fields
	RukovodilaIme string `json:"rukovodilac_ime,omitempty" db:"rukovodilac_ime"`
	BrojZadataka  int    `json:"broj_zadataka,omitempty" db:"broj_zadataka"`
	BrojClanova   int    `json:"broj_clanova,omitempty" db:"broj_clanova"`
}

// ClanoviProjekta represents project team members
type ClanoviProjekta struct {
	ProjekatID int `json:"projekat_id" db:"projekat_id"`
	KorisnikID int `json:"korisnik_id" db:"korisnik_id"`
}

// Zadaci represents tasks within projects
type Zadaci struct {
	ZadatakID            int        `json:"zadatak_id" db:"zadatak_id"`
	ProjekatID           int        `json:"projekat_id" db:"projekat_id"`
	FazaID               int        `json:"faza_id" db:"faza_id"`
	NazivZadatka         string     `json:"naziv_zadatka" db:"naziv_zadatka"`
	Opis                 *string    `json:"opis" db:"opis"`
	DodjeljenKorisnikuID *int       `json:"dodeljen_korisniku_id" db:"dodeljen_korisniku_id"`
	Rok                  *time.Time `json:"rok" db:"rok"`
	Prioritet            *string    `json:"prioritet" db:"prioritet"`
	Progres              int        `json:"progres" db:"progres"`
	KreiranDatuma        time.Time  `json:"kreiran_datuma" db:"kreiran_datuma"`

	// Joined fields
	NazivProjekta      string `json:"naziv_projekta,omitempty" db:"naziv_projekta"`
	NazivFaze          string `json:"naziv_faze,omitempty" db:"naziv_faze"`
	DodjeljenKorisniku string `json:"dodeljen_korisniku,omitempty" db:"dodeljen_korisniku"`
}

// KomentariZadataka represents task comments
type KomentariZadataka struct {
	KomentarID      int       `json:"komentar_id" db:"komentar_id"`
	ZadatakID       int       `json:"zadatak_id" db:"zadatak_id"`
	KorisnikID      int       `json:"korisnik_id" db:"korisnik_id"`
	TekstKomentara  string    `json:"tekst_komentara" db:"tekst_komentara"`
	DatumaKreiranja time.Time `json:"datuma_kreiranja" db:"datuma_kreiranja"`

	// Joined fields
	ImeKorisnika string `json:"ime_korisnika,omitempty" db:"ime_korisnika"`
}

// ZahteviPromeneFaze represents phase change requests
type ZahteviPromeneFaze struct {
	ZahtevID            int       `json:"zahtev_id" db:"zahtev_id"`
	ZadatakID           int       `json:"zadatak_id" db:"zadatak_id"`
	PodnosilacZahtevaID int       `json:"podnosilac_zahteva_id" db:"podnosilac_zahteva_id"`
	ZahtevanaFazaID     int       `json:"zahtevana_faza_id" db:"zahtevana_faza_id"`
	Status              string    `json:"status" db:"status"`
	Komentar            *string   `json:"komentar" db:"komentar"`
	DatumKreiranja      time.Time `json:"datum_kreiranja" db:"datum_kreiranja"`
}

// =============================================================================
// Modul 3: Upravljanje Dokumentima i Meta-podacima
// =============================================================================

// Folderi represents document folders
type Folderi struct {
	FolderID         int    `json:"folder_id" db:"folder_id"`
	NazivFoldera     string `json:"naziv_foldera" db:"naziv_foldera"`
	RoditeljFolderID *int   `json:"roditelj_folder_id" db:"roditelj_folder_id"`
	VlasnikID        int    `json:"vlasnik_id" db:"vlasnik_id"`
}

// Dokumenti represents documents in the system
type Dokumenti struct {
	DokumentID        int        `json:"dokument_id" db:"dokument_id"`
	ProjekatID        *int       `json:"projekat_id" db:"projekat_id"`
	NazivDokumenta    string     `json:"naziv_dokumenta" db:"naziv_dokumenta"`
	FolderID          *int       `json:"folder_id" db:"folder_id"`
	Opis              *string    `json:"opis" db:"opis"`
	TipDokumenta      *string    `json:"tip_dokumenta" db:"tip_dokumenta"`
	JezikDokumenta    *string    `json:"jezik_dokumenta" db:"jezik_dokumenta"`
	RadniTokID        *int       `json:"radni_tok_id" db:"radni_tok_id"`
	TrenutnaFazaID    *int       `json:"trenutna_faza_id" db:"trenutna_faza_id"`
	KreiraoKorisnikID int        `json:"kreirao_korisnik_id" db:"kreirao_korisnik_id"`
	DatumaPostavke    time.Time  `json:"datuma_postavke" db:"datuma_postavke"`
	PoslednjaIzmena   *time.Time `json:"poslednja_izmena" db:"poslednja_izmena"`

	// Joined fields
	NazivProjekta string `json:"naziv_projekta,omitempty" db:"naziv_projekta"`
	ImeKreirao    string `json:"ime_kreirao,omitempty" db:"ime_kreirao"`
	NazivFaze     string `json:"naziv_faze,omitempty" db:"naziv_faze"`
	BrojVerzija   int    `json:"broj_verzija,omitempty" db:"broj_verzija"`
}

// VerzijeDokumenata represents document versions
type VerzijeDokumenata struct {
	VerzijaID          int       `json:"verzija_id" db:"verzija_id"`
	DokumentID         int       `json:"dokument_id" db:"dokument_id"`
	VerzijaOznaka      *string   `json:"verzija_oznaka" db:"verzija_oznaka"`
	PutanjaDoFajla     string    `json:"putanja_do_fajla" db:"putanja_do_fajla"`
	VelicinafajlaMB    *float64  `json:"velicina_fajla_mb" db:"velicina_fajla_mb"`
	PostavioKorisnikID int       `json:"postavio_korisnik_id" db:"postavio_korisnik_id"`
	DatumaPostavke     time.Time `json:"datuma_postavke" db:"datuma_postavke"`
}

// LLMSazeci represents AI-generated document summaries
type LLMSazeci struct {
	SazetakID      int       `json:"sazetak_id" db:"sazetak_id"`
	DokumentID     int       `json:"dokument_id" db:"dokument_id"`
	VerzijaOznaka  *string   `json:"verzija_oznaka" db:"verzija_oznaka"`
	Sazetak        string    `json:"sazetak" db:"sazetak"`
	DatumKreiranja time.Time `json:"datum_kreiranja" db:"datum_kreiranja"`
}

// MetaPodaci represents flexible document metadata
type MetaPodaci struct {
	MetaID     int     `json:"meta_id" db:"meta_id"`
	DokumentID int     `json:"dokument_id" db:"dokument_id"`
	Kljuc      string  `json:"kljuc" db:"kljuc"`
	Vrednost   *string `json:"vrednost" db:"vrednost"`
}

// Tagovi represents document tags
type Tagovi struct {
	TagID     int    `json:"tag_id" db:"tag_id"`
	NazivTaga string `json:"naziv_taga" db:"naziv_taga"`
}

// DokumentTagovi represents many-to-many relationship between documents and tags
type DokumentTagovi struct {
	DokumentID int `json:"dokument_id" db:"dokument_id"`
	TagID      int `json:"tag_id" db:"tag_id"`
}

// DozvoleDokumenata represents document access permissions
type DozvoleDokumenata struct {
	DozvoljID   int  `json:"dozvola_id" db:"dozvola_id"`
	DokumentID  int  `json:"dokument_id" db:"dokument_id"`
	KorisnikID  int  `json:"korisnik_id" db:"korisnik_id"`
	MozeCitati  bool `json:"moze_citati" db:"moze_citati"`
	MozeMenjati bool `json:"moze_menjati" db:"moze_menjati"`
	MozeBrisati bool `json:"moze_brisati" db:"moze_brisati"`
}

// IstorijaFazaDokumenta represents document phase history
type IstorijaFazaDokumenta struct {
	IstorijaID      int       `json:"istorija_id" db:"istorija_id"`
	DokumentID      int       `json:"dokument_id" db:"dokument_id"`
	PrethodnaFazaID *int      `json:"prethodna_faza_id" db:"prethodna_faza_id"`
	NovaFazaID      int       `json:"nova_faza_id" db:"nova_faza_id"`
	KorisnikID      int       `json:"korisnik_id" db:"korisnik_id"`
	DatumPromene    time.Time `json:"datum_promene" db:"datum_promene"`
}

// =============================================================================
// Modul 4: Analitika i Logovanje
// =============================================================================

// LogAktivnosti represents system activity logs
type LogAktivnosti struct {
	LogID          int64     `json:"log_id" db:"log_id"`
	KorisnikID     *int      `json:"korisnik_id" db:"korisnik_id"`
	TipAktivnosti  string    `json:"tip_aktivnosti" db:"tip_aktivnosti"`
	Opis           *string   `json:"opis" db:"opis"`
	CiljaniEntitet *string   `json:"ciljani_entitet" db:"ciljani_entitet"`
	CiljaniID      *int      `json:"ciljani_id" db:"ciljani_id"`
	Datuma         time.Time `json:"datuma" db:"datuma"`

	// Joined fields
	ImeKorisnika string `json:"ime_korisnika,omitempty" db:"ime_korisnika"`
}

// =============================================================================
// DTO strukture za API
// =============================================================================

// LoginRequest represents login request data
type LoginRequest struct {
	KorisnickoIme string `json:"korisnicko_ime" validate:"required"`
	Sifra         string `json:"sifra" validate:"required"`
	ZapamtiMe     bool   `json:"zapamti_me"`
}

// LoginResponse represents login response data
type LoginResponse struct {
	Token    string    `json:"token"`
	Korisnik Korisnici `json:"korisnik"`
	Expires  time.Time `json:"expires"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	AktivniProjekti  int `json:"aktivni_projekti"`
	UkupnoDokumenata int `json:"ukupno_dokumenata"`
	ZadaciUToku      int `json:"zadaci_u_toku"`
	AktivniKorisnici int `json:"aktivni_korisnici"`
}

// CreateProjectRequest represents new project creation data
type CreateProjectRequest struct {
	NazivProjekta  string     `json:"naziv_projekta" validate:"required"`
	Opis           string     `json:"opis"`
	DatumPocetka   *time.Time `json:"datum_pocetka"`
	DatumZavrsetka *time.Time `json:"datum_zavrsetka"`
	RadniTokID     *int       `json:"radni_tok_id"`
	ClanoviTima    []int      `json:"clanovi_tima"`
}

// CreateTaskRequest represents new task creation data
type CreateTaskRequest struct {
	ProjekatID           int        `json:"projekat_id" validate:"required"`
	NazivZadatka         string     `json:"naziv_zadatka" validate:"required"`
	Opis                 string     `json:"opis"`
	DodjeljenKorisnikuID *int       `json:"dodeljen_korisniku_id"`
	Rok                  *time.Time `json:"rok"`
	Prioritet            string     `json:"prioritet"`
}

// UpdateTaskRequest represents task update data
type UpdateTaskRequest struct {
	NazivZadatka         *string    `json:"naziv_zadatka"`
	Opis                 *string    `json:"opis"`
	DodjeljenKorisnikuID *int       `json:"dodeljen_korisniku_id"`
	Rok                  *time.Time `json:"rok"`
	Prioritet            *string    `json:"prioritet"`
	Progres              *int       `json:"progres"`
	FazaID               *int       `json:"faza_id"`
}

// UploadDocumentRequest represents document upload data
type UploadDocumentRequest struct {
	NazivDokumenta string   `json:"naziv_dokumenta" validate:"required"`
	ProjekatID     *int     `json:"projekat_id"`
	FolderID       *int     `json:"folder_id"`
	Opis           string   `json:"opis"`
	TipDokumenta   string   `json:"tip_dokumenta"`
	JezikDokumenta string   `json:"jezik_dokumenta"`
	Tagovi         []string `json:"tagovi"`
}

// =============================================================================
// English aliases for compatibility with existing code
// =============================================================================

// Type aliases for English names
type User = Korisnici
type Role = Uloge
type Project = Projekti
type Task = Zadaci
type Document = Dokumenti
type Folder = Folderi
type Workflow = RadniTokovi
type Phase = Faze
type ProjectMember = ClanoviProjekta
type TaskComment = KomentariZadataka
type PhaseChangeRequest = ZahteviPromeneFaze
type DocumentVersion = VerzijeDokumenata
type LLMSummary = LLMSazeci
type Metadata = MetaPodaci
type Tag = Tagovi
type DocumentTag = DokumentTagovi
type DocumentPermission = DozvoleDokumenata
type DocumentPhaseHistory = IstorijaFazaDokumenta
type ActivityLog = LogAktivnosti
