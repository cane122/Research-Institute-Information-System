package models

import "time"

// User represents a system user
type User struct {
	ID           int       `json:"id" db:"korisnik_id"`
	Username     string    `json:"username" db:"korisnicko_ime"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"hash_sifre"`
	FirstName    string    `json:"firstName" db:"ime"`
	LastName     string    `json:"lastName" db:"prezime"`
	RoleID       int       `json:"roleId" db:"uloga_id"`
	Status       string    `json:"status" db:"status"`
	LastLogin    *time.Time `json:"lastLogin" db:"poslednja_prijava"`
	CreatedAt    time.Time `json:"createdAt" db:"kreiran_datuma"`
	Role         *Role     `json:"role,omitempty"`
}

// Role represents a user role
type Role struct {
	ID   int    `json:"id" db:"uloga_id"`
	Name string `json:"name" db:"naziv_uloge"`
}

// Project represents a research project
type Project struct {
	ID          int       `json:"id" db:"projekat_id"`
	Name        string    `json:"name" db:"naziv_projekta"`
	Description string    `json:"description" db:"opis"`
	StartDate   *time.Time `json:"startDate" db:"datum_pocetka"`
	EndDate     *time.Time `json:"endDate" db:"datum_zavrsetka"`
	Status      string    `json:"status" db:"status"`
	ManagerID   int       `json:"managerId" db:"rukovodilac_id"`
	WorkflowID  *int      `json:"workflowId" db:"radni_tok_id"`
	Manager     *User     `json:"manager,omitempty"`
	Members     []User    `json:"members,omitempty"`
	Tasks       []Task    `json:"tasks,omitempty"`
}

// Task represents a project task
type Task struct {
	ID          int       `json:"id" db:"zadatak_id"`
	ProjectID   int       `json:"projectId" db:"projekat_id"`
	PhaseID     int       `json:"phaseId" db:"faza_id"`
	Name        string    `json:"name" db:"naziv_zadatka"`
	Description string    `json:"description" db:"opis"`
	AssigneeID  *int      `json:"assigneeId" db:"dodeljen_korisniku_id"`
	Deadline    *time.Time `json:"deadline" db:"rok"`
	Priority    string    `json:"priority" db:"prioritet"`
	Progress    int       `json:"progress" db:"progres"`
	CreatedAt   time.Time `json:"createdAt" db:"kreiran_datuma"`
	Assignee    *User     `json:"assignee,omitempty"`
	Phase       *Phase    `json:"phase,omitempty"`
	Comments    []TaskComment `json:"comments,omitempty"`
}

// TaskComment represents a comment on a task
type TaskComment struct {
	ID        int       `json:"id" db:"komentar_id"`
	TaskID    int       `json:"taskId" db:"zadatak_id"`
	UserID    int       `json:"userId" db:"korisnik_id"`
	Content   string    `json:"content" db:"tekst_komentara"`
	CreatedAt time.Time `json:"createdAt" db:"datuma_kreiranja"`
	User      *User     `json:"user,omitempty"`
}

// Document represents a document in the system
type Document struct {
	ID              int       `json:"id" db:"dokument_id"`
	ProjectID       *int      `json:"projectId" db:"projekat_id"`
	Name            string    `json:"name" db:"naziv_dokumenta"`
	FolderID        *int      `json:"folderId" db:"folder_id"`
	Description     string    `json:"description" db:"opis"`
	Type            string    `json:"type" db:"tip_dokumenta"`
	Language        string    `json:"language" db:"jezik_dokumenta"`
	WorkflowID      *int      `json:"workflowId" db:"radni_tok_id"`
	CurrentPhaseID  *int      `json:"currentPhaseId" db:"trenutna_faza_id"`
	CreatedByID     int       `json:"createdById" db:"kreirao_korisnik_id"`
	CreatedAt       time.Time `json:"createdAt" db:"datuma_postavke"`
	LastModified    *time.Time `json:"lastModified" db:"poslednja_izmena"`
	CreatedBy       *User     `json:"createdBy,omitempty"`
	CurrentPhase    *Phase    `json:"currentPhase,omitempty"`
	Versions        []DocumentVersion `json:"versions,omitempty"`
	Tags            []Tag     `json:"tags,omitempty"`
	Metadata        []Metadata `json:"metadata,omitempty"`
}

// DocumentVersion represents a version of a document
type DocumentVersion struct {
	ID         int       `json:"id" db:"verzija_id"`
	DocumentID int       `json:"documentId" db:"dokument_id"`
	Version    string    `json:"version" db:"verzija_oznaka"`
	FilePath   string    `json:"filePath" db:"putanja_do_fajla"`
	FileSizeMB float64   `json:"fileSizeMB" db:"velicina_fajla_MB"`
	UploadedBy int       `json:"uploadedBy" db:"postavio_korisnik_id"`
	UploadedAt time.Time `json:"uploadedAt" db:"datuma_postavke"`
	User       *User     `json:"user,omitempty"`
}

// Workflow represents a workflow template
type Workflow struct {
	ID          int     `json:"id" db:"radni_tok_id"`
	Name        string  `json:"name" db:"naziv"`
	Type        string  `json:"type" db:"tip_toka"`
	Description string  `json:"description" db:"opis"`
	IsTemplate  bool    `json:"isTemplate" db:"da_li_je_sablon"`
	Phases      []Phase `json:"phases,omitempty"`
}

// Phase represents a phase in a workflow
type Phase struct {
	ID         int    `json:"id" db:"faza_id"`
	WorkflowID int    `json:"workflowId" db:"radni_tok_id"`
	Name       string `json:"name" db:"naziv_faze"`
	Order      int    `json:"order" db:"redosled"`
}

// Tag represents a document tag
type Tag struct {
	ID   int    `json:"id" db:"tag_id"`
	Name string `json:"name" db:"naziv_taga"`
}

// Metadata represents document metadata
type Metadata struct {
	ID         int    `json:"id" db:"meta_id"`
	DocumentID int    `json:"documentId" db:"dokument_id"`
	Key        string `json:"key" db:"kljuc"`
	Value      string `json:"value" db:"vrednost"`
}

// Folder represents a document folder
type Folder struct {
	ID       int    `json:"id" db:"folder_id"`
	Name     string `json:"name" db:"naziv_foldera"`
	ParentID *int   `json:"parentId" db:"roditelj_folder_id"`
	OwnerID  int    `json:"ownerId" db:"vlasnik_id"`
	Owner    *User  `json:"owner,omitempty"`
}
