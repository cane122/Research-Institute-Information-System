<template>
  <div class="document-preview-fullscreen">
    <div class="document-preview-page">
      <div class="page-header">
        <div class="header-left">
          <button class="btn btn-secondary" @click="goBack">← Back to Documents</button>
          <h2>{{ document?.name || 'Document Preview' }}</h2>
        </div>
        <div class="header-actions">
          <button class="btn btn-primary" @click="downloadDocument" v-if="document">
            📥 Download
          </button>
          <button class="btn btn-warning" @click="editDocument" v-if="canEdit && document">
            ✏️ Edit
          </button>
          <button class="btn btn-danger" @click="confirmDelete" v-if="canDelete && document">
            🗑️ Delete
          </button>
        </div>
      </div>
      
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner">📄</div>
        <p>Loading document...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <div class="error-icon">❌</div>
        <p>{{ error }}</p>
        <button class="btn btn-secondary" @click="loadDocument">
          Try Again
        </button>
      </div>
      
      <div class="preview-content" v-else-if="document">
        <div class="document-info">
          <div class="breadcrumb">
            Početna > {{ document?.project || 'Nerazvrstano' }} > Dokumenti > {{ document?.name }}
          </div>
          
          <div class="info-sections">
            <div class="document-details">
              <h4>Document Information</h4>
              
              <div class="info-item">
                <strong>Title:</strong> {{ document?.name }}
              </div>
              
              <div class="info-item">
                <strong>Author:</strong> {{ document?.author }}
              </div>
              
              <div class="info-item">
                <strong>Project:</strong> {{ document?.project || 'N/A' }}
              </div>
              
              <div class="info-item">
                <strong>Type:</strong> 
                <span class="type-badge" :class="`type-${document?.type?.toLowerCase()?.replace(' ', '-')}`">
                  {{ document?.type }}
                </span>
              </div>
              
              <div class="info-item">
                <strong>Language:</strong> {{ document?.language }}
              </div>
              
              <div class="info-item">
                <strong>Size:</strong> {{ document?.size }}
              </div>
              
              <div class="info-item">
                <strong>Last Modified:</strong> {{ formatDate(document?.modified) }}
              </div>
              
              <div class="info-item">
                <strong>Created:</strong> {{ formatDate(document?.created) }}
              </div>
              
              <div class="info-item" v-if="document?.description">
                <strong>Description:</strong> 
                <div class="description-text">{{ document.description }}</div>
              </div>

              <div class="info-item" v-if="document?.tags?.length">
                <strong>Tags:</strong>
                <div class="tags-container">
                  <span v-for="tag in document.tags" :key="tag" class="tag-badge">
                    {{ tag }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Separate Version History & Logs Section -->
          <div class="version-logs-section">
            <div class="version-history">
              <h4>📋 Version History & Logs</h4>
              
              <div v-if="versions.length" class="versions-list">
                <div 
                  v-for="version in versions" 
                  :key="version.id"
                  class="version-item"
                  :class="{ 'current-version': version.isCurrent }"
                >
                  <div class="version-info">
                    <div class="version-header">
                      <strong>{{ version.version }} {{ version.isCurrent ? '(Current)' : '' }}</strong>
                      <span class="version-size">{{ version.size }}</span>
                    </div>
                    <div class="version-meta">
                      {{ version.author }} - {{ formatDate(version.date) }}
                    </div>
                    <div class="version-note" v-if="version.note">
                      {{ version.note }}
                    </div>
                    <div class="version-actions">
                      <button 
                        class="btn-sm btn-secondary" 
                        @click="downloadVersion(version)"
                        title="Download this version"
                      >
                        📥 Download
                      </button>
                      <button 
                        v-if="!version.isCurrent && canEdit"
                        class="btn-sm btn-warning" 
                        @click="restoreVersion(version)"
                        title="Restore this version"
                      >
                        🔄 Restore
                      </button>
                    </div>
                  </div>
                </div>
              </div>
              
              <div v-else class="no-versions">
                <p>No version history available</p>
              </div>
            </div>

            <!-- Permissions Section -->
            <div class="permissions-section" v-if="showPermissions">
              <h4>Access Permissions</h4>
              <div class="permissions-list">
                <div class="permission-item">
                  <strong>Read:</strong> {{ permissions.read ? 'Yes' : 'No' }}
                </div>
                <div class="permission-item">
                  <strong>Write:</strong> {{ permissions.write ? 'Yes' : 'No' }}
                </div>
                <div class="permission-item">
                  <strong>Delete:</strong> {{ permissions.delete ? 'Yes' : 'No' }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <div class="footer-actions">
          <button class="btn btn-secondary" @click="closeModal">Close</button>
          <button 
            v-if="canEdit" 
            class="btn btn-primary" 
            @click="editDocument"
          >
            Edit Document
          </button>
          <button 
            v-if="canDelete" 
            class="btn btn-danger" 
            @click="confirmDelete"
          >
            Delete
          </button>
        </div>
      </div>
      
      <!-- Loading State -->
      <div v-if="loading" class="loading-overlay">
        <div class="loading-spinner">
          <div class="spinner"></div>
          <p>Loading document...</p>
        </div>
      </div>
      
      <!-- Error State -->
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
      
      <!-- Empty State -->
      <div v-if="!document && !loading" class="empty-state">
        <h3>Document not found</h3>
        <p>The requested document could not be found.</p>
        <button class="btn btn-primary" @click="goBack">Go Back</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import DocumentService from '../../services/documentService.js'

const router = useRouter()
const route = useRoute()

// Reactive data
const document = ref(null)
const versions = ref([])
const tags = ref([])
const permissions = ref({
  read: true,
  write: false,
  delete: false
})
const loading = ref(true)
const error = ref('')
const versionsLoading = ref(false)

// Computed properties
const canEdit = computed(() => permissions.value.write)
const canDelete = computed(() => permissions.value.delete)

// Methods
async function loadDocument() {
  try {
    loading.value = true
    error.value = ''
    
    const docId = parseInt(route.params.id)
    const doc = await DocumentService.getDocumentById(docId)
    document.value = doc
    
    // Load related data
    await Promise.all([
      loadVersionHistory(),
      loadTags(),
      loadPermissions()
    ])
    
  } catch (err) {
    error.value = err.message
    console.error('Loading document error:', err)
  } finally {
    loading.value = false
  }
}

async function loadVersionHistory() {
  if (!document.value) return
  
  try {
    versionsLoading.value = true
    const docVersions = await DocumentService.getDocumentVersions(document.value.id)
    versions.value = docVersions
  } catch (err) {
    console.error('Loading versions error:', err)
    versions.value = []
  } finally {
    versionsLoading.value = false
  }
}

async function loadTags() {
  if (!document.value) return
  
  try {
    const docTags = await DocumentService.getDocumentTags(document.value.id)
    tags.value = docTags
    // Update document tags for display
    if (document.value) {
      document.value.tags = docTags.map(tag => tag.name)
    }
  } catch (err) {
    console.error('Loading tags error:', err)
    tags.value = []
  }
}

function loadPermissions() {
  // In a real app, this would load user permissions from the backend
  // For now, setting default permissions based on user role
  permissions.value = {
    read: true,
    write: true, // This should be determined by user role and document permissions
    delete: false // This should be determined by user role and document permissions
  }
}

function goBack() {
  router.push('/documents')
}

function downloadDocument() {
  if (document.value) {
    // In a real app, this would trigger a download
    alert(`Downloading: ${document.value.name}`)
    console.log('Download document:', document.value)
  }
}

function editDocument() {
  if (document.value) {
    router.push({
      path: '/documents/add',
      query: { id: document.value.id }
    })
  }
}

function confirmDelete() {
  if (document.value) {
    if (confirm(`Da li ste sigurni da želite da obrišete dokument "${document.value.name}"?`)) {
      deleteDocument()
    }
  }
}

async function deleteDocument() {
  try {
    loading.value = true
    await DocumentService.deleteDocument(document.value.id)
    router.push('/documents')
  } catch (err) {
    error.value = err.message
    console.error('Delete error:', err)
    loading.value = false
  }
}

function downloadVersion(version) {
  // Download specific version
  console.log('Downloading version:', version)
}

// Load document on component mount
onMounted(() => {
  loadDocument()
})

function restoreVersion(version) {
  if (confirm(`Da li ste sigurni da želite da vratite verziju "${version.version}"?`)) {
    // Restore version logic
    console.log('Restoring version:', version)
  }
}

function formatDate(dateString) {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString('sr-RS', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Load document on component mount
onMounted(() => {
  loadDocument()
})
</script>