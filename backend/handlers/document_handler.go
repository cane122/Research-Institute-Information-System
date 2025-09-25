package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cane/research-institute-system/backend/models"
	"github.com/cane/research-institute-system/backend/services"
	"github.com/gorilla/mux"
)

type DocumentHandler struct {
	documentService  *services.DocumentService
	analyticsService *services.AnalyticsService
}

func NewDocumentHandler(documentService *services.DocumentService, analyticsService *services.AnalyticsService) *DocumentHandler {
	return &DocumentHandler{
		documentService:  documentService,
		analyticsService: analyticsService,
	}
}

func (h *DocumentHandler) GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := h.documentService.GetAllDocuments()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, documents)
}

func (h *DocumentHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (32MB max)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	// Read file data
	fileData := make([]byte, fileHeader.Size)
	_, err = file.Read(fileData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Parse document metadata from form
	var req models.UploadDocumentRequest
	req.NazivDokumenta = r.FormValue("nazivDokumenta")
	req.Opis = r.FormValue("opis")
	req.TipDokumenta = r.FormValue("tipDokumenta")
	req.JezikDokumenta = r.FormValue("jezikDokumenta")

	if projectID := r.FormValue("projekatId"); projectID != "" {
		if pid, err := strconv.Atoi(projectID); err == nil {
			req.ProjekatID = &pid
		}
	}

	if folderID := r.FormValue("folderId"); folderID != "" {
		if fid, err := strconv.Atoi(folderID); err == nil {
			req.FolderID = &fid
		}
	}

	// Parse tags (comma-separated)
	if tags := r.FormValue("tagovi"); tags != "" {
		req.Tagovi = strings.Split(tags, ",")
		for i, tag := range req.Tagovi {
			req.Tagovi[i] = strings.TrimSpace(tag)
		}
	}

	userID := getUserIDFromRequest(r)
	err = h.documentService.UploadDocument(req, fileData, fileHeader.Filename, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.analyticsService.LogActivity(&userID, "DOCUMENT_UPLOAD", "Uploaded document: "+req.NazivDokumenta, "document", nil)

	respondWithJSON(w, http.StatusCreated, map[string]bool{"success": true})
}

func (h *DocumentHandler) GetDocumentsByProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	documents, err := h.documentService.GetDocumentsByProject(projectID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, documents)
}

func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	document, err := h.documentService.GetDocumentByID(documentID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Document not found")
		return
	}

	respondWithJSON(w, http.StatusOK, document)
}

func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req models.UploadDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.documentService.UpdateDocument(documentID, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "DOCUMENT_UPDATE", "Updated document", "document", &documentID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	err = h.documentService.DeleteDocument(documentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "DOCUMENT_DELETE", "Deleted document", "document", &documentID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *DocumentHandler) GetDocumentVersions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	versions, err := h.documentService.GetDocumentVersions(documentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, versions)
}

func (h *DocumentHandler) GetDocumentTags(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	tags, err := h.documentService.GetDocumentTags(documentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tags)
}

func (h *DocumentHandler) AddDocumentTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req struct {
		TagName string `json:"tagName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.documentService.AddDocumentTag(documentID, req.TagName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "DOCUMENT_TAG_ADD", "Added tag to document", "document", &documentID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *DocumentHandler) RemoveDocumentTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	tagID, err := strconv.Atoi(vars["tagId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	err = h.documentService.RemoveDocumentTag(documentID, tagID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "DOCUMENT_TAG_REMOVE", "Removed tag from document", "document", &documentID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *DocumentHandler) GetDocumentMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	metadata, err := h.documentService.GetDocumentMetadata(documentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, metadata)
}

func (h *DocumentHandler) UpdateDocumentMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var metadata []models.MetaPodaci
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.documentService.UpdateDocumentMetadata(documentID, metadata)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "DOCUMENT_METADATA_UPDATE", "Updated document metadata", "document", &documentID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *DocumentHandler) GetAllFolders(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromRequest(r)
	folders, err := h.documentService.GetAllFolders(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, folders)
}

func (h *DocumentHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	var folder models.Folderi
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	folder.VlasnikID = getUserIDFromRequest(r)
	err := h.documentService.CreateFolder(folder)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "FOLDER_CREATE", "Created folder: "+folder.NazivFoldera, "folder", nil)

	respondWithJSON(w, http.StatusCreated, map[string]bool{"success": true})
}
