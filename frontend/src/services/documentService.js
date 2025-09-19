// ============================================================================
// documentService.js - Frontend Document Service
// ============================================================================

import { GetAllDocuments, GetDocumentByID, GetDocumentVersions, GetDocumentTags, UploadDocument, UpdateDocument, DeleteDocument } from '../../wailsjs/go/main/App.js'

/**
 * Document service for handling API calls to the backend
 */
export class DocumentService {
  
  /**
   * Get all documents from database
   * @returns {Promise<Array>} Array of documents
   */
  static async getAllDocuments() {
    try {
      const documents = await GetAllDocuments()
      
      // Transform the data to match frontend expectations
      return documents.map(doc => ({
        id: doc.dokument_id,
        name: doc.naziv_dokumenta,
        author: doc.ime_kreirao,
        type: doc.tip_dokumenta || 'Document',
        modified: doc.poslednja_izmena || doc.datuma_postavke,
        project: doc.naziv_projekta || 'Unknown',
        size: this.formatFileSize(0), // Size info from versions
        description: doc.opis || '',
        language: doc.jezik_dokumenta || 'Serbian',
        created: doc.datuma_postavke,
        tags: [], // Will be loaded separately if needed
        versions: doc.broj_verzija || 0,
        currentPhase: doc.naziv_faze || 'Draft'
      }))
    } catch (error) {
      console.error('Error fetching documents:', error)
      throw new Error('Greška pri dohvatanju dokumenata: ' + error.message)
    }
  }

  /**
   * Get a specific document by ID
   * @param {number} documentId - Document ID
   * @returns {Promise<Object>} Document object
   */
  static async getDocumentById(documentId) {
    try {
      const doc = await GetDocumentByID(documentId)
      
      return {
        id: doc.dokument_id,
        name: doc.naziv_dokumenta,
        author: doc.ime_kreirao,
        type: doc.tip_dokumenta || 'Document',
        modified: doc.poslednja_izmena || doc.datuma_postavke,
        project: doc.naziv_projekta || 'Unknown',
        size: this.formatFileSize(0),
        description: doc.opis || '',
        language: doc.jezik_dokumenta || 'Serbian',
        created: doc.datuma_postavke,
        tags: [],
        versions: doc.broj_verzija || 0,
        currentPhase: doc.naziv_faze || 'Draft',
        projectId: doc.projekat_id,
        folderId: doc.folder_id,
        workflowId: doc.radni_tok_id,
        currentPhaseId: doc.trenutna_faza_id,
        createdBy: doc.kreirao_korisnik_id
      }
    } catch (error) {
      console.error('Error fetching document:', error)
      throw new Error('Greška pri dohvatanju dokumenta: ' + error.message)
    }
  }

  /**
   * Get document versions/history
   * @param {number} documentId - Document ID
   * @returns {Promise<Array>} Array of document versions
   */
  static async getDocumentVersions(documentId) {
    try {
      const versions = await GetDocumentVersions(documentId)
      
      return versions.map(version => ({
        id: version.verzija_id,
        version: version.verzija_oznaka || 'v1.0',
        filePath: version.putanja_do_fajla,
        size: this.formatFileSize(version.velicina_fajla_mb * 1024 * 1024), // Convert MB to bytes
        uploadedBy: version.postavio_korisnik_id,
        uploadDate: version.datuma_postavke,
        isCurrent: false // This should be determined by business logic
      }))
    } catch (error) {
      console.error('Error fetching document versions:', error)
      throw new Error('Greška pri dohvatanju verzija dokumenta: ' + error.message)
    }
  }

  /**
   * Get document tags
   * @param {number} documentId - Document ID
   * @returns {Promise<Array>} Array of tags
   */
  static async getDocumentTags(documentId) {
    try {
      const tags = await GetDocumentTags(documentId)
      
      return tags.map(tag => ({
        id: tag.tag_id,
        name: tag.naziv_taga
      }))
    } catch (error) {
      console.error('Error fetching document tags:', error)
      return [] // Return empty array if tags can't be loaded
    }
  }

  /**
   * Upload a new document
   * @param {Object} documentData - Document information
   * @param {File} file - File to upload
   * @returns {Promise<void>}
   */
  static async uploadDocument(documentData, file) {
    try {
      // Convert file to byte array
      const fileData = await this.fileToByteArray(file)
      
      // Prepare request object
      const request = {
        naziv_dokumenta: documentData.name,
        projekat_id: documentData.projectId || null,
        folder_id: documentData.folderId || null,
        opis: documentData.description || '',
        tip_dokumenta: documentData.type || 'Document',
        jezik_dokumenta: documentData.language || 'Serbian',
        tagovi: documentData.tags || []
      }

      await UploadDocument(request, fileData, file.name)
    } catch (error) {
      console.error('Error uploading document:', error)
      throw new Error('Greška pri učitavanju dokumenta: ' + error.message)
    }
  }

  /**
   * Update an existing document
   * @param {number} documentId - Document ID
   * @param {Object} documentData - Updated document information
   * @returns {Promise<void>}
   */
  static async updateDocument(documentId, documentData) {
    try {
      const request = {
        naziv_dokumenta: documentData.name,
        projekat_id: documentData.projectId || null,
        folder_id: documentData.folderId || null,
        opis: documentData.description || '',
        tip_dokumenta: documentData.type || 'Document',
        jezik_dokumenta: documentData.language || 'Serbian',
        tagovi: documentData.tags || []
      }

      await UpdateDocument(documentId, request)
    } catch (error) {
      console.error('Error updating document:', error)
      throw new Error('Greška pri ažuriranju dokumenta: ' + error.message)
    }
  }

  /**
   * Delete a document
   * @param {number} documentId - Document ID
   * @returns {Promise<void>}
   */
  static async deleteDocument(documentId) {
    try {
      await DeleteDocument(documentId)
    } catch (error) {
      console.error('Error deleting document:', error)
      throw new Error('Greška pri brisanju dokumenta: ' + error.message)
    }
  }

  /**
   * Helper method to convert file to byte array
   * @param {File} file - File object
   * @returns {Promise<Array>} Byte array
   */
  static fileToByteArray(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      
      reader.onload = function(event) {
        const arrayBuffer = event.target.result
        const byteArray = new Uint8Array(arrayBuffer)
        resolve(Array.from(byteArray))
      }
      
      reader.onerror = function(error) {
        reject(error)
      }
      
      reader.readAsArrayBuffer(file)
    })
  }

  /**
   * Helper method to format file size
   * @param {number} bytes - File size in bytes
   * @returns {string} Formatted size string
   */
  static formatFileSize(bytes) {
    if (!bytes || bytes === 0) return '0 B'
    
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
  }

  /**
   * Helper method to format date for display
   * @param {string} dateString - ISO date string
   * @returns {string} Formatted date string
   */
  static formatDate(dateString) {
    if (!dateString) return 'Unknown'
    
    const date = new Date(dateString)
    return date.toLocaleDateString('sr-RS') + ' ' + date.toLocaleTimeString('sr-RS', { 
      hour: '2-digit', 
      minute: '2-digit' 
    })
  }
}

export default DocumentService