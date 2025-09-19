<template>
  <Layout>
    <div class="page-container">
      <div class="page-header">
        <h2>Upload Document</h2>
        <button class="btn btn-secondary" @click="$router.push('/documents')">
          ← Back to Documents
        </button>
      </div>

      <div class="card">
        <div class="upload-content">
          <!-- Left Panel -->
          <div class="upload-left-panel">
            <!-- File Selection -->
            <div class="form-group">
              <label>Select File</label>
              <div class="file-select-container">
                <input 
                  type="text" 
                  class="form-input"
                  :value="selectedFiles.length > 0 ? selectedFiles[0].name : 'Choose file (PDF, DOC, XLS, etc. - max 80MB)'"
                  readonly
                />
                <input 
                  ref="fileInput"
                  type="file" 
                  multiple 
                  @change="handleFileSelect"
                  accept=".pdf,.doc,.docx,.xls,.xlsx,.csv,.txt,.png,.jpg,.jpeg"
                  style="display: none"
                >
                <button class="btn btn-secondary" @click="$refs.fileInput.click()">Browse...</button>
              </div>
            </div>

            <!-- Document Name -->
            <div class="form-group">
              <label>Document Name *</label>
              <input type="text" class="form-input" v-model="documentInfo.name" placeholder="Enter document name" />
            </div>

            <!-- Author -->
            <div class="form-group">
              <label>Author *</label>
              <input type="text" class="form-input" v-model="documentInfo.author" placeholder="Enter author name" />
            </div>

            <!-- Keywords -->
            <div class="form-group">
              <label>Keywords</label>
              <input type="text" class="form-input" v-model="documentInfo.keywords" placeholder="Separate keywords with commas" />
            </div>

            <!-- Description -->
            <div class="form-group">
              <label>Description</label>
              <textarea class="form-textarea" v-model="documentInfo.description" rows="3" placeholder="Enter document description"></textarea>
            </div>

            <!-- Bottom Controls -->
            <div class="form-row">
              <div class="form-group">
                <label>Document Type</label>
                <select class="form-select" v-model="documentInfo.type">
                  <option>Research Paper</option>
                  <option>PDF</option>
                  <option>CSV</option>
                  <option>DOC</option>
                  <option>Report</option>
                  <option>Presentation</option>
                </select>
              </div>
              <div class="form-group">
                <label>Language</label>
                <select class="form-select" v-model="documentInfo.language">
                  <option>Serbian</option>
                  <option>English</option>
                  <option>German</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Right Panel -->
          <div class="upload-right-panel">
            <!-- Access Management -->
            <div class="section">
              <h4>Access Management</h4>
              <div class="access-level-group">
                <label class="section-subtitle">Access Level</label>
                <div class="radio-group">
                  <label class="radio-option">
                    <input type="radio" name="access" v-model="accessLevel" value="team">
                    Team
                  </label>
                  <label class="radio-option">
                    <input type="radio" name="access" v-model="accessLevel" value="private">
                    Private
                  </label>
                  <label class="radio-option">
                    <input type="radio" name="access" v-model="accessLevel" value="public">
                    Public
                  </label>
                </div>
              </div>
            </div>

            <!-- Team Permissions -->
            <div class="section">
              <h4>Team Permissions</h4>
              <div class="search-box">
                <input type="text" class="search-input" placeholder="Search users..." v-model="userSearch">
                <span class="search-icon">🔍</span>
              </div>
              
              <div class="users-permissions">
                <div class="permissions-header">
                  <span class="user-col">User</span>
                  <span class="perm-col">R</span>
                  <span class="perm-col">W</span>
                  <span class="perm-col">D</span>
                </div>
                <div class="users-list">
                  <div class="user-permission-item" v-for="user in filteredUsers" :key="user.id">
                    <div class="user-info">
                      <span class="user-name">{{ user.name }}</span>
                    </div>
                    <div class="permission-checkboxes">
                      <label class="permission-checkbox" :title="'Read access for ' + user.name">
                        <input type="checkbox" v-model="user.permissions.read">
                        <span class="checkbox-label">R</span>
                      </label>
                      <label class="permission-checkbox" :title="'Write access for ' + user.name">
                        <input type="checkbox" v-model="user.permissions.write">
                        <span class="checkbox-label">W</span>
                      </label>
                      <label class="permission-checkbox" :title="'Delete access for ' + user.name">
                        <input type="checkbox" v-model="user.permissions.delete">
                        <span class="checkbox-label">D</span>
                      </label>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Additional Metadata -->
            <div class="section">
              <h4>Additional Metadata</h4>
              <div class="form-group">
                <label>Label</label>
                <input type="text" class="form-input" v-model="metadata.label" placeholder="Document label" />
              </div>
              <div class="form-group">
                <label>ISO Number</label>
                <input type="text" class="form-input" v-model="metadata.isoNumber" placeholder="ISO standard number" />
              </div>
              <div class="form-group">
                <label>Source URL</label>
                <input type="url" class="form-input" v-model="metadata.sourceUrl" placeholder="https://..." />
              </div>
              <button class="btn btn-secondary btn-sm" @click="generateTags">Generate Tags</button>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="card-footer">
          <div class="footer-actions">
            <button class="btn btn-secondary" @click="$router.push('/documents')" :disabled="isUploading">
              Cancel
            </button>
            <button class="btn btn-primary" @click="handleSave" :disabled="!canUpload || isUploading">
              {{ isUploading ? 'Uploading...' : 'Upload Document' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Upload Progress -->
      <div v-if="isUploading" class="upload-progress-overlay">
        <div class="upload-progress-card">
          <h4>Uploading Document</h4>
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
          </div>
          <div class="progress-text">
            {{ uploadProgress }}% complete
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '../../components/Layout.vue'
import DocumentService from '../../services/documentService.js'

const router = useRouter()

// Reactive data
const selectedFiles = ref([])
const documentInfo = ref({
  name: '',
  author: '',
  keywords: '',
  description: '',
  type: 'PDF',
  language: 'Serbian'
})

const accessLevel = ref('team')
const userSearch = ref('')
const metadata = ref({
  label: '',
  isoNumber: '',
  sourceUrl: ''
})

const isUploading = ref(false)
const uploadProgress = ref(0)
const error = ref('')

// Users with permissions
const users = ref([
  {
    id: 1,
    name: 'Marko Perić',
    permissions: {
      read: true,
      write: false,
      delete: false
    }
  },
  {
    id: 2,
    name: 'Milica Vujić',
    permissions: {
      read: true,
      write: true,
      delete: false
    }
  },
  {
    id: 3,
    name: 'Elon Musk',
    permissions: {
      read: true,
      write: false,
      delete: false
    }
  },
  {
    id: 4,
    name: 'Ana Nikolić',
    permissions: {
      read: true,
      write: true,
      delete: true
    }
  },
  {
    id: 5,
    name: 'Petar Jovanović',
    permissions: {
      read: true,
      write: false,
      delete: false
    }
  }
])

// Computed
const filteredUsers = computed(() => {
  if (!userSearch.value) return users.value
  return users.value.filter(user => 
    user.name.toLowerCase().includes(userSearch.value.toLowerCase())
  )
})

const canUpload = computed(() => {
  return selectedFiles.value.length > 0 && documentInfo.value.name && documentInfo.value.author
})

// Methods
function handleFileSelect(event) {
  const files = Array.from(event.target.files)
  selectedFiles.value = files
  
  // Auto-fill document name from first file
  if (files.length > 0 && !documentInfo.value.name) {
    documentInfo.value.name = files[0].name.replace(/\.[^/.]+$/, "")
  }
}

function generateTags() {
  // Mock tag generation
  alert('Tags generated based on document content!')
}

async function handleSave() {
  if (selectedFiles.value.length === 0) {
    alert('Please select a file to upload.')
    return
  }

  if (!documentInfo.value.name || !documentInfo.value.author) {
    alert('Please enter document name and author.')
    return
  }

  isUploading.value = true
  uploadProgress.value = 0
  error.value = ''

  try {
    // Prepare document data
    const documentData = {
      name: documentInfo.value.name,
      description: documentInfo.value.description,
      type: documentInfo.value.type,
      language: documentInfo.value.language,
      tags: documentInfo.value.keywords ? documentInfo.value.keywords.split(',').map(tag => tag.trim()) : [],
      projectId: null, // TODO: Add project selection
      folderId: null   // TODO: Add folder selection
    }

    // Simulate progress for better UX
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += Math.random() * 20
      }
    }, 200)

    // Upload the first file (for now)
    const file = selectedFiles.value[0]
    await DocumentService.uploadDocument(documentData, file)
    
    clearInterval(progressInterval)
    uploadProgress.value = 100

    setTimeout(() => {
      isUploading.value = false
      alert('Document uploaded successfully!')
      router.push('/documents')
    }, 500)

  } catch (err) {
    isUploading.value = false
    uploadProgress.value = 0
    error.value = err.message
    alert('Upload failed: ' + err.message)
    console.error('Upload error:', err)
  }
}
</script>