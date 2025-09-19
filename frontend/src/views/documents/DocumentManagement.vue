<template>
  <Layout>
    <div class="document-management">
      <!-- Header -->
      <div class="page-header">
        <h2>Document Management</h2>
        <div class="header-actions">
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
              type="text" 
              placeholder="Search"
              class="sidebar-input"
              style="margin-bottom: 10px;"
            >
            <div class="tags-section">
              <div class="tag-list">
                <div class="tag-item">
                  <input type="checkbox" id="machine-learning" v-model="selectedTags" value="machine-learning">
                  <label for="machine-learning">Machine learning</label>
                </div>
                <div class="tag-item">
                  <input type="checkbox" id="climate" v-model="selectedTags" value="climate">
                  <label for="climate">Climate</label>
                </div>
                <div class="tag-item">
                  <input type="checkbox" id="quantum-physics" v-model="selectedTags" value="quantum-physics">
                  <label for="quantum-physics">Quantum Physics</label>
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
              <select v-model="selectedType" class="sidebar-select" style="width: auto;">
                <option value="">All Types</option>
                <option value="PDF">PDF</option>
                <option value="DOC">DOC</option>
                <option value="CSV">CSV</option>
                <option value="Research Paper">Research Paper</option>
              </select>
              <select v-model="selectedProject" class="sidebar-select" style="width: auto;">
                <option value="">All Projects</option>
                <option value="AI Research">AI Research</option>
                <option value="Climate">Climate</option>
                <option value="Project A">Project A</option>
              </select>
              <button class="search-box button">Search</button>
            </div>
            <div class="view-toggles">
              <label>
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
              <div class="table-header">
                <div class="header-cell">
                  <input 
                    type="checkbox" 
                    @change="toggleSelectAll"
                    :checked="allSelected"
                  >
                </div>
                <div class="header-cell sortable" @click="sortBy('name')">
                  Document Name
                  <span class="sort-indicator" v-if="sortField === 'name'">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </div>
                <div class="header-cell sortable" @click="sortBy('author')">
                  Author
                  <span class="sort-indicator" v-if="sortField === 'author'">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </div>
                <div class="header-cell sortable" @click="sortBy('type')">
                  Type
                  <span class="sort-indicator" v-if="sortField === 'type'">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </div>
                <div class="header-cell sortable" @click="sortBy('modified')">
                  Modified
                  <span class="sort-indicator" v-if="sortField === 'modified'">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </div>
                <div class="header-cell sortable" @click="sortBy('project')">
                  Project
                  <span class="sort-indicator" v-if="sortField === 'project'">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </div>
                <div class="header-cell">Actions</div>
              </div>

              <div class="table-body">
                <div 
                  v-for="doc in paginatedDocuments" 
                  :key="doc.id"
                  class="table-row"
                  :class="{ selected: selectedDocuments.includes(doc.id) }"
                >
                  <div class="cell">
                    <input 
                      type="checkbox" 
                      :checked="selectedDocuments.includes(doc.id)"
                      @change="toggleDocumentSelection(doc.id)"
                    >
                  </div>
                  <div class="cell document-name">
                    <span class="file-icon">📄</span>
                    {{ doc.name }}
                  </div>
                  <div class="cell">{{ doc.author }}</div>
                  <div class="cell">
                    <span class="type-badge" :class="`type-${doc.type.toLowerCase()}`">
                      {{ doc.type }}
                    </span>
                  </div>
                  <div class="cell">{{ formatDate(doc.modified) }}</div>
                  <div class="cell">{{ doc.project }}</div>
                  <div class="cell actions">
                    <button 
                      class="action-btn view-btn" 
                      @click="previewDocument(doc)"
                      title="Preview"
                    >
                      View
                    </button>
                    <button 
                      class="action-btn edit-btn" 
                      @click="editDocument(doc)"
                      title="Edit"
                    >
                      Edit
                    </button>
                    <button 
                      class="action-btn delete-btn" 
                      @click="deleteDocument(doc.id)"
                      title="Delete"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <!-- End of documents table v-else -->

            <!-- Pagination -->
            <div class="pagination" v-if="totalPages > 1">
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
const showFolderPreview = ref(false)
const selectedFolder = ref('')

// Table and pagination states
const selectedDocuments = ref([])
const sortField = ref('modified')
const sortDirection = ref('desc')
const currentPage = ref(1)
const itemsPerPage = ref(10)

// Modal states
const showUploadModal = ref(false)
const showPreviewModal = ref(false)

// Computed properties
const filteredDocuments = computed(() => {
  let filtered = documents.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(doc =>
      doc.name.toLowerCase().includes(query) ||
      doc.author.toLowerCase().includes(query) ||
      doc.description.toLowerCase().includes(query)
    )
  }

  if (authorFilter.value) {
    const author = authorFilter.value.toLowerCase()
    filtered = filtered.filter(doc =>
      doc.author.toLowerCase().includes(author)
    )
  }

  if (selectedType.value) {
    filtered = filtered.filter(doc => doc.type === selectedType.value)
  }

  if (selectedProject.value) {
    filtered = filtered.filter(doc => doc.project === selectedProject.value)
  }

  if (dateFrom.value) {
    filtered = filtered.filter(doc => doc.modified >= dateFrom.value)
  }

  if (dateTo.value) {
    filtered = filtered.filter(doc => doc.modified <= dateTo.value)
  }

  return filtered
})

const sortedDocuments = computed(() => {
  const docs = [...filteredDocuments.value]
  
  if (sortField.value) {
    docs.sort((a, b) => {
      const aVal = a[sortField.value]
      const bVal = b[sortField.value]
      
      if (sortDirection.value === 'asc') {
        return aVal < bVal ? -1 : aVal > bVal ? 1 : 0
      } else {
        return aVal > bVal ? -1 : aVal < bVal ? 1 : 0
      }
    })
  }
  
  return docs
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
         paginatedDocuments.value.every(doc => selectedDocuments.value.includes(doc.id))
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

function previewDocument(doc) {
  // Navigate to document preview page with document ID
  router.push(`/documents/preview/${doc.id}`)
}

function editDocument(doc) {
  // Navigate to document add/edit page with document ID as query parameter
  router.push({
    path: '/documents/add',
    query: { id: doc.id }
  })
}

async function deleteDocument(docId) {
  if (confirm('Are you sure you want to delete this document?')) {
    try {
      await DocumentService.deleteDocument(docId)
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
    const docs = await DocumentService.getAllDocuments()
    documents.value = docs
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