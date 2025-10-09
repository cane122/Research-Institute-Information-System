<template>
  <Layout>
    <div class="document-management">
      <!-- Header -->
      <div class="page-header">
        <h2>Document Management</h2>
        <div class="header-actions">
          <button 
            class="btn btn-secondary" 
            @click="$router.push('/documents/analytics')"
            style="margin-right: 10px;"
          >
            📊 Analytics
          </button>
          <button 
            class="btn btn-primary" 
            @click="goToAddDocument"
          >
            Upload Document
          </button>
        </div>
      </div>

      <!-- Main Content with Sidebar Layout -->
      <div class="main-content">
        <!-- Left Sidebar for Filters -->
        <div class="filter-sidebar">
          <h3>Filter Documents</h3>
          
          <!-- Search -->
          <div class="sidebar-section">
            <h4>Search</h4>
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Search Documents..."
              class="sidebar-input"
            >
          </div>

          <!-- Filter Author -->
          <div class="sidebar-section">
            <h4>Filter Author</h4>
            <input 
              v-model="authorFilter"
              type="text" 
              placeholder="Search by author"
              class="sidebar-input"
              @input="applyFilters"
            >
          </div>

          <!-- Date Range -->
          <div class="sidebar-section">
            <h4>Date Range</h4>
            <div class="date-inputs">
              <input 
                v-model="dateFrom" 
                type="date" 
                @change="applyFilters"
                placeholder="From"
              >
              <input 
                v-model="dateTo" 
                type="date" 
                @change="applyFilters"
                placeholder="To"
              >
            </div>
          </div>

          <!-- Tags -->
          <div class="sidebar-section">
            <h4>Tags</h4>
            <input 
              v-model="tagSearchQuery"
              type="text" 
              placeholder="Search tags..."
              class="sidebar-input"
              style="margin-bottom: 10px;"
            >
            <div class="tags-section">
              <div class="tag-list">
                <div v-for="tag in filteredTags" :key="tag.tag_id" class="tag-item">
                  <input 
                    type="checkbox" 
                    :id="`tag-${tag.tag_id}`" 
                    :value="tag.naziv_taga"
                    v-model="selectedTags"
                    @change="applyFilters"
                  >
                  <label :for="`tag-${tag.tag_id}`">{{ tag.naziv_taga }}</label>
                </div>
                <div v-if="filteredTags.length === 0" class="empty-tags">
                  No tags available
                </div>
              </div>
            </div>
          </div>

          <!-- Apply Filters Button -->
          <button class="apply-filters-btn" @click="applyFilters">
            Apply filters
          </button>
        </div>

        <!-- Right Content Area -->
        <div class="content-area">
          <!-- Content Header with Search and Controls -->
          <div class="content-header">
            <div class="search-box">
              <select v-model="selectedType" class="sidebar-select" style="width: auto;" @change="applyFilters">
                <option value="">All Types</option>
                <option v-for="type in availableTypes" :key="type" :value="type">{{ type }}</option>
              </select>
              <select v-model="selectedProject" class="sidebar-select" style="width: auto;" @change="applyFilters">
                <option value="">All Projects</option>
                <option v-for="project in availableProjects" :key="project" :value="project">{{ project }}</option>
              </select>
              <button class="search-box button" @click="applyFilters">Search</button>
            </div>
            <div class="view-toggles">
              <label style="margin-left: 15px;">
                <input type="checkbox" v-model="showFolderPreview">
                Folder preview
              </label>
            </div>
          </div>

          <!-- Document List -->
          <div class="documents-panel">
            <!-- Loading State -->
            <div v-if="loading" class="loading-state">
              <div class="loading-spinner">📄</div>
              <p>Loading documents...</p>
            </div>

            <!-- Error State -->
            <div v-else-if="error" class="error-state">
              <div class="error-icon">❌</div>
              <p>{{ error }}</p>
              <button class="btn btn-secondary" @click="loadDocuments">
                Try Again
              </button>
            </div>

            <!-- Documents Table -->
            <div v-else class="documents-table">
              <!-- Folder Preview View -->
              <template v-if="showFolderPreview">
                <div v-for="(docs, projectName) in groupedDocuments" :key="projectName" class="project-folder">
                  <!-- Folder Header (clickable to expand/collapse) -->
                  <div 
                    class="folder-header" 
                    @click="toggleProjectFolder(projectName)"
                    :class="{ expanded: expandedProjects.includes(projectName) }"
                  >
                    <span class="folder-icon">
                      {{ expandedProjects.includes(projectName) ? '📂' : '📁' }}
                    </span>
                    <span class="folder-name">{{ projectName }}</span>
                    <span class="folder-count">({{ docs.length }}) docs</span>
                    <span class="expand-arrow">
                      {{ expandedProjects.includes(projectName) ? '▼' : '▶' }}
                    </span>
                  </div>
                  
                  <!-- Folder Contents (shown when expanded) -->
                  <div v-if="expandedProjects.includes(projectName)" class="folder-contents">
                    <div class="table-header">
                      <div class="header-cell">
                        <input type="checkbox" @change="toggleSelectAllInGroup(docs)">
                      </div>
                      <div class="header-cell">Document Name</div>
                      <div class="header-cell">Author</div>
                      <div class="header-cell">Type</div>
                      <div class="header-cell">Modified</div>
                      <div class="header-cell">Actions</div>
                    </div>

                    <div class="table-body">
                      <div 
                        v-for="doc in docs" 
                        :key="doc.dokument_id"
                        class="table-row document-item"
                        :class="{ selected: selectedDocuments.includes(doc.dokument_id) }"
                      >
                        <div class="cell">
                          <input 
                            type="checkbox" 
                            :checked="selectedDocuments.includes(doc.dokument_id)"
                            @change="toggleDocumentSelection(doc.dokument_id)"
                          >
                        </div>
                        <div class="cell document-name">
                          <span class="file-icon">📄</span>
                          {{ doc.naziv_dokumenta }}
                        </div>
                        <div class="cell">{{ doc.ime_kreirao }}</div>
                        <div class="cell">
                          <span class="type-badge" :class="`type-${doc.tip_dokumenta?.toLowerCase() || 'document'}`">
                            {{ doc.tip_dokumenta || 'Document' }}
                          </span>
                        </div>
                        <div class="cell">{{ formatDate(doc.poslednja_izmena || doc.datuma_postavke) }}</div>
                        <div class="cell actions">
                          <button 
                            v-if="doc.permissions?.canRead"
                            class="action-btn view-btn" 
                            @click="previewDocument(doc)"
                            title="Preview"
                          >
                            View
                          </button>
                          <button 
                            v-if="doc.permissions?.canWrite"
                            class="action-btn edit-btn" 
                            @click="editDocument(doc)"
                            title="Edit"
                          >
                            Edit
                          </button>
                          <button 
                            v-if="doc.permissions?.canDelete"
                            class="action-btn delete-btn" 
                            @click="deleteDocument(doc.dokument_id)"
                            title="Delete"
                          >
                            Delete
                          </button>
                          <span v-if="!doc.permissions?.canRead && !doc.permissions?.canWrite && !doc.permissions?.canDelete" class="no-access-text">
                            No Access
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </template>

              <!-- Regular Table View -->
              <template v-else>
                <div class="table-header">
                  <div class="header-cell">
                    <input 
                      type="checkbox" 
                      @change="toggleSelectAll"
                      :checked="allSelected"
                    >
                  </div>
                  <div class="header-cell sortable" @click="sortBy('naziv_dokumenta')">
                    Document Name
                    <span class="sort-indicator" v-if="sortField === 'naziv_dokumenta'">
                      {{ sortDirection === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                  <div class="header-cell sortable" @click="sortBy('ime_kreirao')">
                    Author
                    <span class="sort-indicator" v-if="sortField === 'ime_kreirao'">
                      {{ sortDirection === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                  <div class="header-cell sortable" @click="sortBy('tip_dokumenta')">
                    Type
                    <span class="sort-indicator" v-if="sortField === 'tip_dokumenta'">
                      {{ sortDirection === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                  <div class="header-cell sortable" @click="sortBy('poslednja_izmena')">
                    Modified
                    <span class="sort-indicator" v-if="sortField === 'poslednja_izmena'">
                      {{ sortDirection === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                  <div class="header-cell sortable" @click="sortBy('naziv_projekta')">
                    Project
                    <span class="sort-indicator" v-if="sortField === 'naziv_projekta'">
                      {{ sortDirection === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                  <div class="header-cell">Actions</div>
                </div>

                <div class="table-body">
                  <div 
                    v-for="doc in paginatedDocuments" 
                    :key="doc.dokument_id"
                    class="table-row"
                    :class="{ selected: selectedDocuments.includes(doc.dokument_id) }"
                  >
                    <div class="cell">
                      <input 
                        type="checkbox" 
                        :checked="selectedDocuments.includes(doc.dokument_id)"
                        @change="toggleDocumentSelection(doc.dokument_id)"
                      >
                    </div>
                    <div class="cell document-name">
                      <span class="file-icon">📄</span>
                      {{ doc.naziv_dokumenta }}
                    </div>
                    <div class="cell">{{ doc.ime_kreirao }}</div>
                    <div class="cell">
                      <span class="type-badge" :class="`type-${doc.tip_dokumenta?.toLowerCase() || 'document'}`">
                        {{ doc.tip_dokumenta || 'Document' }}
                      </span>
                    </div>
                    <div class="cell">{{ formatDate(doc.poslednja_izmena || doc.datuma_postavke) }}</div>
                    <div class="cell">{{ doc.naziv_projekta || 'N/A' }}</div>
                    <div class="cell actions">
                      <button 
                        v-if="doc.permissions?.canRead"
                        class="action-btn view-btn" 
                        @click="previewDocument(doc)"
                        title="Preview"
                      >
                        View
                      </button>
                      <button 
                        v-if="doc.permissions?.canWrite"
                        class="action-btn edit-btn" 
                        @click="editDocument(doc)"
                        title="Edit"
                      >
                        Edit
                      </button>
                      <button 
                        v-if="doc.permissions?.canDelete"
                        class="action-btn delete-btn" 
                        @click="deleteDocument(doc.dokument_id)"
                        title="Delete"
                      >
                        Delete
                      </button>
                      <span v-if="!doc.permissions?.canRead && !doc.permissions?.canWrite && !doc.permissions?.canDelete" class="no-access-text">
                        No Access
                      </span>
                    </div>
                  </div>
                </div>
              </template>
            </div>
            <!-- End of documents table v-else -->

            <!-- Pagination (only for regular table view) -->
            <div class="pagination" v-if="!showFolderPreview && totalPages > 1">
              <button 
                class="page-btn" 
                :disabled="currentPage === 1"
                @click="currentPage = 1"
              >
                First
              </button>
              <button 
                class="page-btn" 
                :disabled="currentPage === 1"
                @click="currentPage--"
              >
                Previous
              </button>
              
              <span class="page-info">
                Page {{ currentPage }} of {{ totalPages }}
              </span>
              
              <button 
                class="page-btn" 
                :disabled="currentPage === totalPages"
                @click="currentPage++"
              >
                Next
              </button>
              <button 
                class="page-btn" 
                :disabled="currentPage === totalPages"
                @click="currentPage = totalPages"
              >
                Last
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Loading Overlay -->
      <div v-if="loading" class="loading-overlay">
        <div class="loading-spinner">
          <div class="spinner"></div>
          <p>Loading documents...</p>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '../../components/Layout.vue'
import { useAuthStore } from '../../stores/auth'
import DocumentService from '../../services/documentService.js'

const router = useRouter()
const authStore = useAuthStore()

function goToAddDocument() {
  router.push({ path: '/documents/add' })
}

// Reactive data
const documents = ref([])
const allTags = ref([])
const loading = ref(true)
const error = ref('')

// Filter and search states
const searchQuery = ref('')
const authorFilter = ref('')
const selectedType = ref('')
const selectedProject = ref('')
const dateFrom = ref('')
const dateTo = ref('')
const selectedTags = ref([])
const tagSearchQuery = ref('')
const showFolderPreview = ref(false)
const expandedProjects = ref([]) // Track which project folders are expanded
const selectedFolder = ref('')

// Table and pagination states
const selectedDocuments = ref([])
const sortField = ref('poslednja_izmena')
const sortDirection = ref('desc')
const currentPage = ref(1)
const itemsPerPage = ref(10)

// Modal states
const showUploadModal = ref(false)
const showPreviewModal = ref(false)

// Computed properties
const filteredTags = computed(() => {
  if (!tagSearchQuery.value) return allTags.value
  const query = tagSearchQuery.value.toLowerCase()
  return allTags.value.filter(tag => 
    tag.naziv_taga.toLowerCase().includes(query)
  )
})

const availableTypes = computed(() => {
  const types = new Set()
  documents.value.forEach(doc => {
    if (doc.tip_dokumenta) types.add(doc.tip_dokumenta)
  })
  return Array.from(types).sort()
})

const availableProjects = computed(() => {
  const projects = new Set()
  documents.value.forEach(doc => {
    if (doc.naziv_projekta) projects.add(doc.naziv_projekta)
  })
  return Array.from(projects).sort()
})

const filteredDocuments = computed(() => {
  let filtered = documents.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(doc =>
      doc.naziv_dokumenta.toLowerCase().includes(query) ||
      doc.ime_kreirao.toLowerCase().includes(query) ||
      (doc.opis && doc.opis.toLowerCase().includes(query))
    )
  }

  if (authorFilter.value) {
    const author = authorFilter.value.toLowerCase()
    filtered = filtered.filter(doc =>
      doc.ime_kreirao.toLowerCase().includes(author)
    )
  }

  if (selectedType.value) {
    filtered = filtered.filter(doc => doc.tip_dokumenta === selectedType.value)
  }

  if (selectedProject.value) {
    filtered = filtered.filter(doc => doc.naziv_projekta === selectedProject.value)
  }

  if (dateFrom.value) {
    const fromDate = new Date(dateFrom.value)
    filtered = filtered.filter(doc => {
      const docDate = new Date(doc.poslednja_izmena || doc.datuma_postavke)
      return docDate >= fromDate
    })
  }

  if (dateTo.value) {
    const toDate = new Date(dateTo.value)
    filtered = filtered.filter(doc => {
      const docDate = new Date(doc.poslednja_izmena || doc.datuma_postavke)
      return docDate <= toDate
    })
  }

  // Tag filtering - filter documents that have ANY of the selected tags
  if (selectedTags.value.length > 0) {
    filtered = filtered.filter(doc => {
      // If document has tags loaded and at least one matches selected tags
      return doc.tagovi && doc.tagovi.some(tag => 
        selectedTags.value.includes(tag.naziv_taga)
      )
    })
  }

  return filtered
})

const sortedDocuments = computed(() => {
  const docs = [...filteredDocuments.value]
  
  if (sortField.value) {
    docs.sort((a, b) => {
      let aVal = a[sortField.value]
      let bVal = b[sortField.value]
      
      // Handle null/undefined values
      if (aVal === null || aVal === undefined) aVal = ''
      if (bVal === null || bVal === undefined) bVal = ''
      
      // Handle date fields
      if (sortField.value === 'poslednja_izmena' || sortField.value === 'datuma_postavke') {
        aVal = new Date(aVal || 0).getTime()
        bVal = new Date(bVal || 0).getTime()
      }
      
      if (sortDirection.value === 'asc') {
        return aVal < bVal ? -1 : aVal > bVal ? 1 : 0
      } else {
        return aVal > bVal ? -1 : aVal < bVal ? 1 : 0
      }
    })
  }
  
  return docs
})

const groupedDocuments = computed(() => {
  if (!showFolderPreview.value) return {}
  
  const groups = {}
  sortedDocuments.value.forEach(doc => {
    const projectName = doc.naziv_projekta || 'Unassigned'
    if (!groups[projectName]) {
      groups[projectName] = []
    }
    groups[projectName].push(doc)
  })
  
  return groups
})

const paginatedDocuments = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return sortedDocuments.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(sortedDocuments.value.length / itemsPerPage.value)
})

const allSelected = computed(() => {
  return paginatedDocuments.value.length > 0 && 
         paginatedDocuments.value.every(doc => selectedDocuments.value.includes(doc.dokument_id))
})

// Methods
function searchDocuments() {
  // Trigger reactive filtering (already handled by computed properties)
}

function applyFilters() {
  // Trigger reactive filtering (already handled by computed properties)
  currentPage.value = 1
}

function clearFilters() {
  searchQuery.value = ''
  authorFilter.value = ''
  selectedType.value = ''
  selectedProject.value = ''
  dateFrom.value = ''
  dateTo.value = ''
  selectedTags.value = []
  currentPage.value = 1
}

function sortBy(field) {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDirection.value = 'asc'
  }
}

function toggleSelectAll() {
  if (allSelected.value) {
    selectedDocuments.value = selectedDocuments.value.filter(id => 
      !paginatedDocuments.value.some(doc => doc.id === id)
    )
  } else {
    const newSelections = paginatedDocuments.value.map(doc => doc.id)
    selectedDocuments.value = [...new Set([...selectedDocuments.value, ...newSelections])]
  }
}

function toggleDocumentSelection(docId) {
  const index = selectedDocuments.value.indexOf(docId)
  if (index > -1) {
    selectedDocuments.value.splice(index, 1)
  } else {
    selectedDocuments.value.push(docId)
  }
}

function toggleSelectAllInGroup(docs) {
  const allGroupSelected = docs.every(doc => selectedDocuments.value.includes(doc.dokument_id))
  
  if (allGroupSelected) {
    // Remove all docs from this group
    selectedDocuments.value = selectedDocuments.value.filter(id => 
      !docs.some(doc => doc.dokument_id === id)
    )
  } else {
    // Add all docs from this group
    const newSelections = docs.map(doc => doc.dokument_id)
    selectedDocuments.value = [...new Set([...selectedDocuments.value, ...newSelections])]
  }
}

function toggleProjectFolder(projectName) {
  const index = expandedProjects.value.indexOf(projectName)
  if (index > -1) {
    expandedProjects.value.splice(index, 1)
  } else {
    expandedProjects.value.push(projectName)
  }
}

function previewDocument(doc) {
  // Check if user has read permission
  if (!doc.permissions?.canRead) {
    alert('You do not have permission to view this document.')
    return
  }
  
  // Log view activity
  try {
    const { LogActivity } = window.go.main.App
    LogActivity({
      tip_aktivnosti: 'VIEW',
      entitet_tip: 'DOKUMENT',
      entitet_id: doc.dokument_id,
      naziv_entiteta: doc.naziv_dokumenta,
      opis: `Viewed document: ${doc.naziv_dokumenta}`,
      rezultat: 'SUCCESS'
    }).catch(err => console.error('Failed to log view activity:', err))
  } catch (err) {
    console.error('Failed to log view activity:', err)
  }
  
  // Navigate to document preview page with document ID
  router.push(`/documents/preview/${doc.dokument_id}`)
}

function editDocument(doc) {
  // Check if user has write permission
  if (!doc.permissions?.canWrite) {
    alert('You do not have permission to edit this document.')
    return
  }
  // Navigate to document add/edit page with document ID as query parameter
  router.push({
    path: '/documents/add',
    query: { id: doc.dokument_id }
  })
}

async function deleteDocument(docId) {
  // Find the document to check permissions
  const doc = documents.value.find(d => d.dokument_id === docId)
  
  if (!doc || !doc.permissions?.canDelete) {
    alert('You do not have permission to delete this document.')
    return
  }
  
  if (confirm('Are you sure you want to delete this document?')) {
    try {
      await DocumentService.deleteDocument(docId)
      
      // Log delete activity
      try {
        const { LogActivity } = window.go.main.App
        await LogActivity({
          tip_aktivnosti: 'DELETE',
          entitet_tip: 'DOKUMENT',
          entitet_id: docId,
          naziv_entiteta: doc.naziv_dokumenta,
          opis: `Deleted document: ${doc.naziv_dokumenta}`,
          rezultat: 'SUCCESS'
        })
      } catch (logErr) {
        console.error('Failed to log delete activity:', logErr)
      }
      
      // Refresh the document list
      await loadDocuments()
    } catch (err) {
      error.value = err.message
      console.error('Delete error:', err)
    }
  }
}

async function loadDocuments() {
  try {
    loading.value = true
    error.value = ''
    
    const { GetAllDocuments, GetAllTags, GetDocumentTags, CheckUserPermission } = window.go.main.App
    
    // Load documents and tags in parallel
    const [docs, tags] = await Promise.all([
      GetAllDocuments(),
      GetAllTags()
    ])
    
    // Load tags and permissions for each document
    if (docs && docs.length > 0) {
      const docsWithTagsAndPermissions = await Promise.all(
        docs.map(async (doc) => {
          try {
            // Load tags for this document
            const docTags = await GetDocumentTags(doc.dokument_id)
            
            // Load permissions for current user on this document
            const [canRead, canWrite, canDelete] = await Promise.all([
              CheckUserPermission(doc.dokument_id, 'read').catch(() => false),
              CheckUserPermission(doc.dokument_id, 'write').catch(() => false),
              CheckUserPermission(doc.dokument_id, 'delete').catch(() => false)
            ])
            
            return { 
              ...doc, 
              tagovi: docTags || [],
              permissions: {
                canRead,
                canWrite,
                canDelete
              }
            }
          } catch (err) {
            console.error(`Error loading data for document ${doc.dokument_id}:`, err)
            return { 
              ...doc, 
              tagovi: [],
              permissions: {
                canRead: false,
                canWrite: false,
                canDelete: false
              }
            }
          }
        })
      )
      documents.value = docsWithTagsAndPermissions
    } else {
      documents.value = []
    }
    
    allTags.value = tags || []
    
    console.log('Loaded documents with tags and permissions:', documents.value)
    console.log('Loaded tags:', allTags.value)
  } catch (err) {
    error.value = err.message
    console.error('Loading error:', err)
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleDateString()
}

// Load documents on component mount
onMounted(() => {
  loadDocuments()
})
</script>

<style scoped>
/* Override global styles for folder preview */
.documents-table {
  overflow: visible !important;
}

.project-folder {
  margin-bottom: 15px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: visible;
  background: white;
  width: 100%;
  max-width: none;
}

.folder-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 20px;
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s ease;
  width: 100%;
}

.folder-header:hover {
  background: linear-gradient(135deg, #5568d3 0%, #6a3f8f 100%);
}

.folder-header.expanded {
  border-bottom: 1px solid #dee2e6;
}

.folder-icon {
  font-size: 1.2em;
  margin-right: 10px;
}

.folder-name {
  font-size: 1em;
  font-weight: 600;
  flex: 1;
}

.folder-count {
  opacity: 0.9;
  font-size: 0.9em;
  margin-right: 15px;
}

.expand-arrow {
  font-size: 0.8em;
}

.folder-contents {
  background: #ffffff;
  width: 100%;
  max-height: none !important;
  overflow: visible !important;
}

.folder-contents .table-header {
  background: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

.folder-contents .table-body {
  max-height: none !important;
  overflow: visible !important;
  flex: none !important;
}

.document-item {
  border-left: 2px solid #e0e0e0;
}

.document-item:hover {
  background: #f8f9fa;
  border-left-color: #667eea;
}

.empty-tags {
  color: #999;
  font-style: italic;
  padding: 10px;
  text-align: center;
}

.view-toggles {
  display: flex;
  align-items: center;
}

.view-toggles label {
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.view-toggles input[type="checkbox"] {
  margin-right: 5px;
}

.no-access-text {
  color: #999;
  font-size: 0.85em;
  font-style: italic;
  padding: 4px 8px;
}
</style>