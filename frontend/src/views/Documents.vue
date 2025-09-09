<template>
  <Layout>
    <div class="documents">
      <!-- Header -->
      <div class="page-header">
        <div>
          <h2>Dokumenti</h2>
          <div class="breadcrumb">Početna > Dokumenti</div>
        </div>
        <button class="btn btn-primary" @click="showUploadModal = true">
          <span class="btn-icon">⬆️</span>
          Upload dokument
        </button>
      </div>
      
      <!-- Filters and Search -->
      <div class="filters card">
        <div class="filter-group">
          <label>Tip dokumenta:</label>
          <select v-model="filters.type">
            <option value="">Svi tipovi</option>
            <option value="pdf">PDF</option>
            <option value="doc">Word dokument</option>
            <option value="xls">Excel tabela</option>
            <option value="ppt">PowerPoint</option>
            <option value="txt">Tekst fajl</option>
            <option value="other">Ostalo</option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>Projekat:</label>
          <select v-model="filters.project">
            <option value="">Svi projekti</option>
            <option 
              v-for="project in projects" 
              :key="project.id" 
              :value="project.id"
            >
              {{ project.name }}
            </option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>Kategorija:</label>
          <select v-model="filters.category">
            <option value="">Sve kategorije</option>
            <option value="research">Istraživanje</option>
            <option value="report">Izveštaj</option>
            <option value="specification">Specifikacija</option>
            <option value="manual">Priručnik</option>
            <option value="presentation">Prezentacija</option>
            <option value="other">Ostalo</option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>Pretraga:</label>
          <input 
            type="text" 
            v-model="filters.search" 
            placeholder="Pretraži dokumente..."
            class="search-input"
          >
        </div>
      </div>
      
      <!-- View Mode Toggle -->
      <div class="view-controls">
        <div class="view-toggle">
          <button 
            class="view-btn"
            :class="{ active: viewMode === 'grid' }"
            @click="viewMode = 'grid'"
          >
            🔲 Grid
          </button>
          <button 
            class="view-btn"
            :class="{ active: viewMode === 'list' }"
            @click="viewMode = 'list'"
          >
            📋 Lista
          </button>
        </div>
        
        <div class="sort-options">
          <label>Sortiranje:</label>
          <select v-model="sortBy">
            <option value="name">Naziv</option>
            <option value="size">Veličina</option>
            <option value="created">Datum kreiranja</option>
            <option value="updated">Poslednja izmena</option>
          </select>
        </div>
      </div>
      
      <!-- Documents Grid View -->
      <div v-if="viewMode === 'grid'" class="documents-grid">
        <div 
          v-for="document in filteredDocuments" 
          :key="document.id"
          class="document-card"
          @click="openDocument(document)"
        >
          <div class="document-icon">
            <span class="file-icon">{{ getFileIcon(document.type) }}</span>
            <div class="document-actions">
              <button 
                class="action-btn" 
                @click.stop="downloadDocument(document)"
                title="Preuzmi"
              >
                ⬇️
              </button>
              <button 
                class="action-btn" 
                @click.stop="shareDocument(document)"
                title="Podeli"
              >
                🔗
              </button>
              <button 
                class="action-btn danger" 
                @click.stop="deleteDocument(document)"
                title="Obriši"
              >
                🗑️
              </button>
            </div>
          </div>
          
          <div class="document-info">
            <h4 class="document-name" :title="document.name">
              {{ document.name }}
            </h4>
            <div class="document-meta">
              <span class="document-size">{{ formatFileSize(document.size) }}</span>
              <span class="document-type">{{ document.type.toUpperCase() }}</span>
            </div>
            <div class="document-details">
              <div class="document-project">{{ getProjectName(document.projectId) }}</div>
              <div class="document-category">
                <span 
                  class="category-badge" 
                  :class="`category-${document.category}`"
                >
                  {{ getCategoryText(document.category) }}
                </span>
              </div>
            </div>
            <div class="document-footer">
              <div class="document-author">{{ document.author }}</div>
              <div class="document-date">{{ formatDate(document.created) }}</div>
            </div>
          </div>
        </div>
        
        <!-- Empty State -->
        <div v-if="filteredDocuments.length === 0" class="empty-state">
          <div class="empty-icon">📄</div>
          <h3>Nema dokumenata</h3>
          <p>Upload-ujte prvi dokument klikom na dugme "Upload dokument"</p>
        </div>
      </div>
      
      <!-- Documents List View -->
      <div v-else class="documents-list">
        <div class="documents-table">
          <table class="table">
            <thead>
              <tr>
                <th>Naziv</th>
                <th>Tip</th>
                <th>Veličina</th>
                <th>Projekat</th>
                <th>Kategorija</th>
                <th>Autor</th>
                <th>Datum</th>
                <th>Akcije</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="document in filteredDocuments" 
                :key="document.id"
                @click="openDocument(document)"
                class="document-row"
              >
                <td>
                  <div class="document-name-cell">
                    <span class="file-icon-small">{{ getFileIcon(document.type) }}</span>
                    <span class="document-name">{{ document.name }}</span>
                  </div>
                </td>
                <td>{{ document.type.toUpperCase() }}</td>
                <td>{{ formatFileSize(document.size) }}</td>
                <td>{{ getProjectName(document.projectId) }}</td>
                <td>
                  <span 
                    class="category-badge" 
                    :class="`category-${document.category}`"
                  >
                    {{ getCategoryText(document.category) }}
                  </span>
                </td>
                <td>{{ document.author }}</td>
                <td>{{ formatDate(document.created) }}</td>
                <td>
                  <div class="action-buttons">
                    <button 
                      class="btn-icon-small" 
                      @click.stop="downloadDocument(document)"
                      title="Preuzmi"
                    >
                      ⬇️
                    </button>
                    <button 
                      class="btn-icon-small" 
                      @click.stop="shareDocument(document)"
                      title="Podeli"
                    >
                      🔗
                    </button>
                    <button 
                      class="btn-icon-small danger" 
                      @click.stop="deleteDocument(document)"
                      title="Obriši"
                    >
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <!-- Upload Modal -->
      <div v-if="showUploadModal" class="modal-overlay" @click="closeModal">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h3 class="modal-title">Upload dokumenta</h3>
            <button class="modal-close" @click="closeModal">×</button>
          </div>
          
          <form @submit.prevent="uploadDocument">
            <div class="upload-area">
              <div class="upload-zone" @click="triggerFileInput">
                <div class="upload-icon">📁</div>
                <h4>Kliknite ili prevucite fajlove ovde</h4>
                <p>Podržani formati: PDF, DOC, XLS, PPT, TXT</p>
                <input 
                  type="file" 
                  ref="fileInput"
                  @change="handleFileSelect"
                  multiple
                  accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt"
                  style="display: none"
                >
              </div>
              
              <div v-if="selectedFiles.length > 0" class="selected-files">
                <h4>Izabrani fajlovi:</h4>
                <div 
                  v-for="(file, index) in selectedFiles" 
                  :key="index"
                  class="file-item"
                >
                  <span class="file-icon">{{ getFileIconFromName(file.name) }}</span>
                  <span class="file-name">{{ file.name }}</span>
                  <span class="file-size">{{ formatFileSize(file.size) }}</span>
                  <button 
                    type="button"
                    class="remove-file"
                    @click="removeFile(index)"
                  >
                    ×
                  </button>
                </div>
              </div>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label>Projekat</label>
                <select v-model="uploadForm.projectId" required>
                  <option value="">Izaberite projekat</option>
                  <option 
                    v-for="project in projects" 
                    :key="project.id" 
                    :value="project.id"
                  >
                    {{ project.name }}
                  </option>
                </select>
              </div>
              
              <div class="form-group">
                <label>Kategorija</label>
                <select v-model="uploadForm.category" required>
                  <option value="research">Istraživanje</option>
                  <option value="report">Izveštaj</option>
                  <option value="specification">Specifikacija</option>
                  <option value="manual">Priručnik</option>
                  <option value="presentation">Prezentacija</option>
                  <option value="other">Ostalo</option>
                </select>
              </div>
            </div>
            
            <div class="form-group">
              <label>Opis (opciono)</label>
              <textarea 
                v-model="uploadForm.description" 
                placeholder="Unesite opis dokumenta..."
                rows="3"
              ></textarea>
            </div>
            
            <div class="form-actions">
              <button type="button" class="btn btn-secondary" @click="closeModal">
                Otkaži
              </button>
              <button 
                type="submit" 
                class="btn btn-primary"
                :disabled="selectedFiles.length === 0"
              >
                Upload {{ selectedFiles.length }} fajl(ova)
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Layout from '../components/Layout.vue'

// Reactive data
const viewMode = ref('grid')
const showUploadModal = ref(false)
const selectedFiles = ref([])
const fileInput = ref(null)
const sortBy = ref('name')

const filters = ref({
  type: '',
  project: '',
  category: '',
  search: ''
})

const uploadForm = ref({
  projectId: '',
  category: 'other',
  description: ''
})

// Mock data
const documents = ref([
  {
    id: 1,
    name: 'Projektna specifikacija v2.1.pdf',
    type: 'pdf',
    size: 2540000,
    projectId: 1,
    category: 'specification',
    author: 'Marko Petrović',
    created: '2024-09-01',
    updated: '2024-09-05'
  },
  {
    id: 2,
    name: 'Analiza zahteva.docx',
    type: 'doc',
    size: 850000,
    projectId: 1,
    category: 'research',
    author: 'Ana Jovanović',
    created: '2024-08-28',
    updated: '2024-08-30'
  },
  {
    id: 3,
    name: 'Mesečni izveštaj avgust.xlsx',
    type: 'xls',
    size: 456000,
    projectId: 2,
    category: 'report',
    author: 'Stefan Nikolić',
    created: '2024-09-02',
    updated: '2024-09-02'
  },
  {
    id: 4,
    name: 'Prezentacija rezultata.pptx',
    type: 'ppt',
    size: 12400000,
    projectId: 2,
    category: 'presentation',
    author: 'Milica Stojković',
    created: '2024-08-25',
    updated: '2024-08-26'
  },
  {
    id: 5,
    name: 'API dokumentacija.pdf',
    type: 'pdf',
    size: 1200000,
    projectId: 1,
    category: 'manual',
    author: 'Marko Petrović',
    created: '2024-08-20',
    updated: '2024-08-22'
  }
])

const projects = ref([
  { id: 1, name: 'Web Portal Refactoring' },
  { id: 2, name: 'AI Analitika Modula' },
  { id: 3, name: 'Mobile Aplikacija' }
])

// Computed
const filteredDocuments = computed(() => {
  let filtered = documents.value

  // Apply filters
  if (filters.value.type) {
    filtered = filtered.filter(d => d.type === filters.value.type)
  }
  
  if (filters.value.project) {
    filtered = filtered.filter(d => d.projectId == filters.value.project)
  }
  
  if (filters.value.category) {
    filtered = filtered.filter(d => d.category === filters.value.category)
  }
  
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(d => 
      d.name.toLowerCase().includes(search) ||
      d.author.toLowerCase().includes(search)
    )
  }

  // Sort
  filtered.sort((a, b) => {
    switch (sortBy.value) {
      case 'name':
        return a.name.localeCompare(b.name)
      case 'size':
        return b.size - a.size
      case 'created':
        return new Date(b.created) - new Date(a.created)
      case 'updated':
        return new Date(b.updated) - new Date(a.updated)
      default:
        return 0
    }
  })

  return filtered
})

// Methods
function getFileIcon(type) {
  const icons = {
    pdf: '📋',
    doc: '📝',
    xls: '📊',
    ppt: '📊',
    txt: '📄',
    other: '📄'
  }
  return icons[type] || icons.other
}

function getFileIconFromName(fileName) {
  const extension = fileName.split('.').pop().toLowerCase()
  if (['pdf'].includes(extension)) return '📋'
  if (['doc', 'docx'].includes(extension)) return '📝'
  if (['xls', 'xlsx'].includes(extension)) return '📊'
  if (['ppt', 'pptx'].includes(extension)) return '📊'
  if (['txt'].includes(extension)) return '📄'
  return '📄'
}

function getProjectName(projectId) {
  return projects.value.find(p => p.id === projectId)?.name || 'Nepoznat projekat'
}

function getCategoryText(category) {
  const categoryMap = {
    research: 'Istraživanje',
    report: 'Izveštaj',
    specification: 'Specifikacija',
    manual: 'Priručnik',
    presentation: 'Prezentacija',
    other: 'Ostalo'
  }
  return categoryMap[category] || category
}

function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString('sr-RS')
}

function openDocument(document) {
  console.log('Opening document:', document)
  // Here you would typically open the document in a viewer or download it
}

function downloadDocument(document) {
  console.log('Downloading document:', document)
  // Here you would trigger the download
}

function shareDocument(document) {
  console.log('Sharing document:', document)
  // Here you would show sharing options or copy link
}

function deleteDocument(document) {
  if (confirm(`Da li ste sigurni da želite da obrišete dokument "${document.name}"?`)) {
    const index = documents.value.findIndex(d => d.id === document.id)
    if (index > -1) {
      documents.value.splice(index, 1)
    }
  }
}

function triggerFileInput() {
  fileInput.value.click()
}

function handleFileSelect(event) {
  const files = Array.from(event.target.files)
  selectedFiles.value = [...selectedFiles.value, ...files]
}

function removeFile(index) {
  selectedFiles.value.splice(index, 1)
}

function uploadDocument() {
  if (selectedFiles.value.length === 0) return
  
  // In a real app, this would upload files to the server
  selectedFiles.value.forEach(file => {
    const extension = file.name.split('.').pop().toLowerCase()
    let type = 'other'
    
    if (['pdf'].includes(extension)) type = 'pdf'
    else if (['doc', 'docx'].includes(extension)) type = 'doc'
    else if (['xls', 'xlsx'].includes(extension)) type = 'xls'
    else if (['ppt', 'pptx'].includes(extension)) type = 'ppt'
    else if (['txt'].includes(extension)) type = 'txt'
    
    const newDocument = {
      id: Date.now() + Math.random(),
      name: file.name,
      type: type,
      size: file.size,
      projectId: uploadForm.value.projectId,
      category: uploadForm.value.category,
      author: 'Trenutni korisnik', // Replace with actual user
      created: new Date().toISOString().split('T')[0],
      updated: new Date().toISOString().split('T')[0]
    }
    
    documents.value.unshift(newDocument)
  })
  
  closeModal()
}

function closeModal() {
  showUploadModal.value = false
  selectedFiles.value = []
  uploadForm.value = {
    projectId: '',
    category: 'other',
    description: ''
  }
}

// Lifecycle
onMounted(() => {
  console.log('Documents loaded')
})
</script>

<style scoped>
.documents {
  padding: 30px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 30px;
}

.page-header h2 {
  color: #2c3e50;
  margin-bottom: 8px;
  font-size: 28px;
}

.breadcrumb {
  color: #7f8c8d;
  font-size: 14px;
}

.btn-icon {
  margin-right: 8px;
}

/* Filters */
.filters {
  display: flex;
  gap: 20px;
  align-items: flex-end;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  min-width: 150px;
}

.filter-group label {
  margin-bottom: 5px;
  font-size: 14px;
  color: #2c3e50;
}

.filter-group select,
.search-input {
  padding: 8px 12px;
  border: 1px solid #bdc3c7;
  border-radius: 4px;
  font-size: 14px;
}

.search-input {
  min-width: 250px;
}

/* View Controls */
.view-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 15px 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.view-toggle {
  display: flex;
  border: 1px solid #bdc3c7;
  border-radius: 4px;
  overflow: hidden;
}

.view-btn {
  padding: 8px 16px;
  border: none;
  background: white;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.view-btn.active {
  background: #3498db;
  color: white;
}

.view-btn:hover:not(.active) {
  background: #f8f9fa;
}

.sort-options {
  display: flex;
  align-items: center;
  gap: 10px;
}

.sort-options label {
  font-size: 14px;
  color: #2c3e50;
}

/* Documents Grid */
.documents-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.document-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.document-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.document-icon {
  text-align: center;
  margin-bottom: 15px;
  position: relative;
}

.file-icon {
  font-size: 48px;
}

.document-actions {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  gap: 5px;
  opacity: 0;
  transition: opacity 0.3s;
}

.document-card:hover .document-actions {
  opacity: 1;
}

.action-btn {
  background: white;
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.3s;
}

.action-btn:hover {
  transform: scale(1.1);
}

.action-btn.danger:hover {
  background: #e74c3c;
  color: white;
}

.document-name {
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}

.document-meta {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 12px;
  color: #7f8c8d;
}

.document-type {
  background: #ecf0f1;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
}

.document-details {
  margin-bottom: 15px;
}

.document-project {
  font-size: 13px;
  color: #7f8c8d;
  margin-bottom: 8px;
}

.category-badge {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.category-research {
  background: #d1ecf1;
  color: #0c5460;
}

.category-report {
  background: #d4edda;
  color: #155724;
}

.category-specification {
  background: #fff3cd;
  color: #856404;
}

.category-manual {
  background: #f8d7da;
  color: #721c24;
}

.category-presentation {
  background: #e2e3f1;
  color: #383d41;
}

.category-other {
  background: #e9ecef;
  color: #495057;
}

.document-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #95a5a6;
  border-top: 1px solid #ecf0f1;
  padding-top: 12px;
}

/* Documents List */
.documents-table {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.document-row {
  cursor: pointer;
  transition: background 0.3s;
}

.document-row:hover {
  background: #f8f9fa;
}

.document-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.file-icon-small {
  font-size: 16px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-icon-small {
  background: none;
  border: none;
  font-size: 14px;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: background 0.3s;
}

.btn-icon-small:hover {
  background: #f8f9fa;
}

.btn-icon-small.danger:hover {
  background: #f8d7da;
}

/* Empty State */
.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 20px;
  color: #7f8c8d;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 24px;
  margin-bottom: 10px;
  color: #2c3e50;
}

/* Upload Modal */
.modal {
  max-width: 700px;
  width: 90%;
}

.upload-area {
  margin-bottom: 30px;
}

.upload-zone {
  border: 2px dashed #bdc3c7;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 20px;
}

.upload-zone:hover {
  border-color: #3498db;
  background: #f8f9fa;
}

.upload-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.upload-zone h4 {
  color: #2c3e50;
  margin-bottom: 5px;
}

.upload-zone p {
  color: #7f8c8d;
  font-size: 14px;
}

.selected-files h4 {
  color: #2c3e50;
  margin-bottom: 15px;
  font-size: 16px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 10px;
}

.file-name {
  flex: 1;
  font-weight: 500;
}

.file-size {
  color: #7f8c8d;
  font-size: 14px;
}

.remove-file {
  background: #e74c3c;
  color: white;
  border: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 30px;
  border-top: 1px solid #ecf0f1;
  padding-top: 20px;
}

/* Responsive */
@media (max-width: 768px) {
  .documents {
    padding: 20px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 20px;
    align-items: stretch;
  }
  
  .filters {
    flex-direction: column;
    gap: 15px;
  }
  
  .filter-group {
    min-width: auto;
  }
  
  .search-input {
    min-width: auto;
  }
  
  .view-controls {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
  
  .documents-grid {
    grid-template-columns: 1fr;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column-reverse;
  }
}

@media (max-width: 480px) {
  .documents-table {
    overflow-x: auto;
  }
  
  .table {
    min-width: 700px;
  }
}
</style>
