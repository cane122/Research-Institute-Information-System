<template>
  <Layout>
    <div class="document-management">
      <!-- Header -->
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

      <!-- Filters and Search -->
      <div class="filters-section">
        <div class="search-bar">
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Search Documents..."
            class="search-input"
          >
          <button class="search-btn" @click="searchDocuments">
            🔍
          </button>
        </div>

        <div class="filters">
          <div class="filter-group">
            <label>Type:</label>
            <select v-model="selectedType" @change="applyFilters">
              <option value="">All Types</option>
              <option value="PDF">PDF</option>
              <option value="DOC">DOC</option>
              <option value="CSV">CSV</option>
              <option value="Research Paper">Research Paper</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Project:</label>
            <select v-model="selectedProject" @change="applyFilters">
              <option value="">All Projects</option>
              <option value="AI Research">AI Research</option>
              <option value="Climate">Climate</option>
              <option value="Project A">Project A</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Date Range:</label>
            <input 
              v-model="dateFrom" 
              type="date" 
              @change="applyFilters"
            >
            <span>to</span>
            <input 
              v-model="dateTo" 
              type="date" 
              @change="applyFilters"
            >
          </div>

          <button class="btn btn-sm btn-secondary" @click="clearFilters">
            Clear Filters
          </button>
        </div>

        <div class="view-options">
          <label class="view-toggle">
            <input 
              v-model="showFolderPreview" 
              type="checkbox"
              @change="handleViewToggle"
            > 
            <span class="toggle-text">
              📂 Folder Preview
              <small v-if="showFolderPreview">(Tree View)</small>
              <small v-else>(Table View)</small>
            </span>
          </label>
        </div>
      </div>

      <!-- Main Content -->
      <div class="main-content" :class="{ 'with-folder-preview': showFolderPreview }">
        <!-- Folder Preview -->
        <div v-if="showFolderPreview" class="folder-panel">
          <div class="folder-header">
            <h3>Folder Structure</h3>
            <div class="folder-actions">
              <button class="btn btn-sm" @click="createFolder">Create Folder</button>
              <button class="btn btn-sm" @click="moveFiles">Move</button>
              <button class="btn btn-sm btn-danger" @click="deleteSelected">Delete</button>
            </div>
          </div>
          
          <div class="folder-tree">
            <!-- Project Folders -->
            <div v-for="project in folderStructure" :key="project.id" class="folder-item project">
              <div 
                class="folder-header" 
                :class="{ expanded: project.expanded }"
                @click="toggleFolder(project)"
              >
                <span class="expand-icon">{{ project.expanded ? '�' : '�📁' }}</span>
                <span class="folder-name">{{ project.name }}</span>
                <span class="document-count">({{ getTotalDocumentCount(project) }})</span>
              </div>
              
              <!-- Project Documents -->
              <div v-show="project.expanded" class="folder-content">
                <div v-for="doc in project.documents" :key="doc.id" class="document-item">
                  <span class="doc-icon">📄</span>
                  <span class="doc-name">{{ doc.name }}</span>
                  <span class="doc-type">{{ doc.type }}</span>
                </div>
                
                <!-- Subfolders -->
                <div v-for="subfolder in project.children" :key="subfolder.id" class="folder-item subfolder">
                  <div 
                    class="folder-header"
                    :class="{ expanded: subfolder.expanded }"
                    @click="toggleFolder(subfolder)"
                  >
                    <span class="expand-icon">{{ subfolder.expanded ? '📂' : '📁' }}</span>
                    <span class="folder-name">{{ subfolder.name }}</span>
                    <span class="document-count">({{ subfolder.documents.length }})</span>
                  </div>
                  
                  <!-- Subfolder Documents -->
                  <div v-show="subfolder.expanded" class="folder-content">
                    <div v-for="doc in subfolder.documents" :key="doc.id" class="document-item">
                      <span class="doc-icon">�</span>
                      <span class="doc-name">{{ doc.name }}</span>
                      <span class="doc-type">{{ doc.type }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Document List -->
        <div class="documents-panel">
          <div class="documents-table">
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
                    title="Preview Document"
                  >
                    VIEW
                  </button>
                  <button 
                    class="action-btn edit-btn" 
                    @click="editDocument(doc)"
                    title="Edit Document"
                  >
                    EDIT
                  </button>
                  <button 
                    class="action-btn delete-btn" 
                    @click="deleteDocument(doc.id)"
                    title="Delete Document"
                  >
                    DELETE
                  </button>
                </div>
              </div>
            </div>
          </div>

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
       <!-- Loading Overlay -->
       <div v-if="loading" class="loading-overlay">
         <div class="loading-spinner">
           <div class="spinner"></div>
           <p>Loading documents...</p>
         </div>
       </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

function goToAddDocument() {
  router.push({ path: '/documents/add' })
}

// Reactive data
const documents = ref([
  {
    id: 1,
    name: 'Research Methodology v3.2',
    author: 'Dr. Petar Nikolić',
    type: 'Research Paper',
    modified: '2025-08-30',
    project: 'AI Research',
    size: '2.4 MB',
    description: 'Comprehensive methodology for AI research projects.'
  },
  {
    id: 2,
    name: 'Research Analayse v4.2',
    author: 'Dr. Petar Nikolić',
    type: 'CSV',
    modified: '2025-04-27',
    project: 'AI Research',
    size: '1.2 MB',
    description: 'Analysis results from recent research.'
  },
  {
    id: 3,
    name: 'Climate Change Final',
    author: 'Dr. Nikoline Marković',
    type: 'Research Paper',
    modified: '2024-06-30',
    project: 'Climate',
    size: '3.8 MB',
    description: 'Final report on climate change research.'
  }
])


// Modal states
const showUploadModal = ref(false)
const showPreviewModal = ref(false)

// Form data
const uploadForm = ref({
  name: '',
  author: '',
  description: '',
  type: 'Research Paper',
  language: 'Serbian',
  project: '',
  tags: ''
})

// File handling
const selectedFiles = ref([])
const isDragOver = ref(false)
const uploading = ref(false)

// Filters and search
const searchQuery = ref('')
const selectedType = ref('')
const selectedProject = ref('')
const dateFrom = ref('')
const dateTo = ref('')
const showFolderPreview = ref(route.query.view === 'folders')
const selectedFolder = ref('')

// Folder tree structure with projects as root folders
const folderStructure = ref([
  {
    id: 'ai-research',
    name: 'AI Research',
    type: 'project',
    expanded: false,
    documents: [
      { id: 1, name: 'Research Methodology v3.2', type: 'Research Paper', author: 'Dr. Petar Nikolić' },
      { id: 2, name: 'Research Analayse v4.2', type: 'CSV', author: 'Dr. Petar Nikolić' }
    ],
    children: [
      {
        id: 'ai-research-models',
        name: 'AI Models',
        type: 'folder',
        expanded: false,
        documents: [
          { id: 4, name: 'Neural Network Architecture', type: 'PDF', author: 'Dr. Petar Nikolić' },
          { id: 5, name: 'Training Data Set', type: 'CSV', author: 'Dr. Ana Jovanović' }
        ],
        children: []
      },
      {
        id: 'ai-research-reports',
        name: 'Reports',
        type: 'folder',
        expanded: false,
        documents: [
          { id: 6, name: 'Monthly Progress Report', type: 'DOC', author: 'Dr. Petar Nikolić' }
        ],
        children: []
      }
    ]
  },
  {
    id: 'climate-research',
    name: 'Climate Research',
    type: 'project',
    expanded: false,
    documents: [
      { id: 3, name: 'Climate Change Final', type: 'Research Paper', author: 'Dr. Nikoline Marković' }
    ],
    children: [
      {
        id: 'climate-data',
        name: 'Data Collection',
        type: 'folder',
        expanded: false,
        documents: [
          { id: 7, name: 'Temperature Measurements', type: 'CSV', author: 'Dr. Nikoline Marković' },
          { id: 8, name: 'Humidity Data', type: 'CSV', author: 'Dr. Milan Stojanović' }
        ],
        children: []
      }
    ]
  },
  {
    id: 'blockchain-project',
    name: 'Blockchain Research',
    type: 'project',
    expanded: false,
    documents: [
      { id: 9, name: 'Blockchain Architecture', type: 'PDF', author: 'Dr. Marko Petrović' }
    ],
    children: [
      {
        id: 'blockchain-security',
        name: 'Security Analysis',
        type: 'folder',
        expanded: false,
        documents: [
          { id: 10, name: 'Security Audit Report', type: 'PDF', author: 'Dr. Marko Petrović' }
        ],
        children: []
      }
    ]
  }
])

// Selection
const selectedDocuments = ref([])

// Sorting
const sortField = ref('modified')
const sortDirection = ref('desc')

// Pagination
const currentPage = ref(1)
const itemsPerPage = ref(10)
const loading = ref(false)

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

  if (selectedType.value) {
    filtered = filtered.filter(doc => doc.type === selectedType.value)
  }

  if (selectedProject.value) {
    filtered = filtered.filter(doc => doc.project === selectedProject.value)
  }

  if (dateFrom.value || dateTo.value) {
    filtered = filtered.filter(doc => {
      const docDate = new Date(doc.modified)
      const from = dateFrom.value ? new Date(dateFrom.value) : new Date('1900-01-01')
      const to = dateTo.value ? new Date(dateTo.value) : new Date('2100-12-31')
      return docDate >= from && docDate <= to
    })
  }

  // Sort
  filtered.sort((a, b) => {
    let aVal = a[sortField.value]
    let bVal = b[sortField.value]

    if (sortField.value === 'modified') {
      aVal = new Date(aVal)
      bVal = new Date(bVal)
    }

    if (aVal < bVal) return sortDirection.value === 'asc' ? -1 : 1
    if (aVal > bVal) return sortDirection.value === 'asc' ? 1 : -1
    return 0
  })

  return filtered
})

const paginatedDocuments = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredDocuments.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(filteredDocuments.value.length / itemsPerPage.value)
})

const allSelected = computed(() => {
  return paginatedDocuments.value.length > 0 && 
    selectedDocuments.value.length === paginatedDocuments.value.length
})

// Methods
function searchDocuments() {
  currentPage.value = 1
}

function applyFilters() {
  currentPage.value = 1
}

function clearFilters() {
  searchQuery.value = ''
  selectedType.value = ''
  selectedProject.value = ''
  dateFrom.value = ''
  dateTo.value = ''
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
    selectedDocuments.value = []
  } else {
    selectedDocuments.value = paginatedDocuments.value.map(doc => doc.id)
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

function selectFolder(folderId) {
  selectedFolder.value = folderId
}

function handleViewToggle() {
  // Update the URL query parameter based on checkbox state
  const newQuery = { ...route.query }
  
  if (showFolderPreview.value) {
    newQuery.view = 'folders'
  } else {
    delete newQuery.view
  }
  
  // Navigate to new URL with updated query parameters
  router.push({ 
    path: route.path, 
    query: newQuery 
  })
}

function toggleFolder(folder) {
  folder.expanded = !folder.expanded
}

function getTotalDocumentCount(project) {
  let count = project.documents.length
  if (project.children) {
    project.children.forEach(child => {
      count += child.documents.length
    })
  }
  return count
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
      // await DeleteDocument(docId)
      documents.value = documents.value.filter(doc => doc.id !== docId)
    } catch (error) {
      console.error('Error deleting document:', error)
    }
  }
}

function handleFileDrop(e) {
  e.preventDefault()
  isDragOver.value = false
  const files = Array.from(e.dataTransfer.files)
  selectedFiles.value.push(...files)
}

function handleFileSelect(e) {
  const files = Array.from(e.target.files)
  selectedFiles.value.push(...files)
}

function removeFile(index) {
  selectedFiles.value.splice(index, 1)
}

async function uploadDocuments() {
  if (selectedFiles.value.length === 0) return

  uploading.value = true
  try {
    // In a real Wails app, you would call the backend function
    // await UploadDocument(selectedFiles.value, uploadForm.value)
    
    // Simulate upload delay
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // Add to local state for demo
    selectedFiles.value.forEach((file, index) => {
      documents.value.push({
        id: Date.now() + index,
        name: uploadForm.value.name || file.name,
        author: uploadForm.value.author || 'Current User',
        type: uploadForm.value.type,
        modified: new Date().toISOString().split('T')[0],
        project: uploadForm.value.project || 'Unassigned',
        size: formatFileSize(file.size),
        description: uploadForm.value.description
      })
    })
    
    closeModal()
  } catch (error) {
    console.error('Error uploading documents:', error)
  } finally {
    uploading.value = false
  }
}

function generateReport() {
  // Generate report functionality
  console.log('Generating report...')
}

function createFolder() {
  const folderName = prompt('Enter folder name:')
  if (folderName) {
    // Create folder logic
    console.log('Creating folder:', folderName)
  }
}

function moveFiles() {
  if (selectedDocuments.value.length === 0) {
    alert('Please select documents to move')
    return
  }
  // Move files logic
  console.log('Moving files:', selectedDocuments.value)
}

function deleteSelected() {
  if (selectedDocuments.value.length === 0) {
    alert('Please select documents to delete')
    return
  }
  
  if (confirm(`Are you sure you want to delete ${selectedDocuments.value.length} document(s)?`)) {
    documents.value = documents.value.filter(doc => !selectedDocuments.value.includes(doc.id))
    selectedDocuments.value = []
  }
}

function closeModal() {
  showUploadModal.value = false
  selectedFiles.value = []
  uploadForm.value = {
    name: '',
    author: '',
    description: '',
    type: 'Research Paper',
    language: 'Serbian',
    project: '',
    tags: ''
  }
}

function closePreviewModal() {
  showPreviewModal.value = false
  previewDocument.value = null
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString('sr-RS')
}

function formatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Lifecycle
onMounted(async () => {
  loading.value = true
  try {
    // In a real Wails app:
    // const docs = await GetDocuments()
    // documents.value = docs
  } catch (error) {
    console.error('Error loading documents:', error)
  } finally {
    loading.value = false
  }
})

// Watch for search changes
watch(searchQuery, () => {
  if (searchQuery.value.length > 2 || searchQuery.value.length === 0) {
    searchDocuments()
  }
})

// Watch for route query changes (browser navigation)
watch(() => route.query.view, (newView) => {
  showFolderPreview.value = newView === 'folders'
})
</script>