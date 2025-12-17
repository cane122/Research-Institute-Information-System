<template>
  <Layout>
    <div class="page-container">
      <div class="page-header">
        <h2>{{ isEditMode ? 'Edit Document' : 'Upload Document' }}</h2>
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
              <label>Select File{{ isEditMode ? ' (optional - leave empty to keep current file)' : '' }}</label>
              <div class="file-select-container">
                <input 
                  type="text" 
                  class="form-input"
                  :value="selectedFiles.length > 0 ? selectedFiles[0].name : (isEditMode ? 'No new file selected - keeping current file' : 'Choose file (PDF, DOC, XLS, etc. - max 80MB)')"
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

            <!-- Keywords (Free-form) -->
            <div class="form-group">
              <label>
                Keywords
                <span class="label-hint" title="Add custom keywords">
                  💡 Type and press Enter or click Add
                </span>
              </label>
              <div class="keywords-input-group">
                <input 
                  type="text" 
                  class="form-input" 
                  v-model="currentKeyword" 
                  placeholder="Type keyword and press Enter" 
                  @keypress.enter.prevent="addKeyword"
                />
                <button 
                  type="button" 
                  class="btn btn-sm btn-add-tag" 
                  @click="addKeyword"
                  title="Add keyword"
                >
                  + Add
                </button>
                <button 
                  type="button"
                  class="btn btn-secondary btn-sm btn-generate-keywords" 
                  @click="generateKeywords"
                  :disabled="isGeneratingKeywords || !documentInfo.name"
                  :title="!documentInfo.name ? 'Please enter document name first' : 'Generate keywords using AI'"
                >
                  {{ isGeneratingKeywords ? '🤖 Generating...' : '🤖 Generate Keywords' }}
                </button>
              </div>
              
              <!-- Keywords Display -->
              <div v-if="keywords.length > 0" class="tags-container">
                <div 
                  v-for="(keyword, index) in keywords" 
                  :key="index" 
                  class="tag-label"
                >
                  <span class="tag-text">{{ keyword }}</span>
                  <button 
                    class="tag-remove" 
                    @click="removeKeyword(index)"
                    type="button"
                    title="Remove keyword"
                  >
                    ×
                  </button>
                </div>
              </div>
            </div>

            <!-- Tags (From database) -->
            <div class="form-group">
              <label>
                Tags (From Database)
                <span class="label-hint" title="Select existing tags or use AI to suggest from available tags">
                  🏷️ Only existing tags
                </span>
              </label>
              
              <!-- Available tags as buttons -->
              <div class="available-tags-list">
                <div 
                  v-for="tag in availableTags" 
                  :key="tag.tag_id" 
                  class="tag-item"
                >
                  <span class="tag-name">{{ tag.naziv_taga }}</span>
                  <button 
                    type="button"
                    class="btn btn-xs btn-add-tag-from-db" 
                    @click="addTagFromDb(tag.tag_id)"
                    :disabled="selectedTags.includes(tag.tag_id)"
                    :title="selectedTags.includes(tag.tag_id) ? 'Already added' : 'Add this tag'"
                  >
                    {{ selectedTags.includes(tag.tag_id) ? '✓' : '+' }}
                  </button>
                </div>
              </div>
              
              <button 
                type="button"
                class="btn btn-secondary btn-sm btn-generate-tags" 
                @click="generateTags"
                :disabled="isGeneratingTags || !documentInfo.name || availableTags.length === 0"
                :title="!documentInfo.name ? 'Please enter document name first' : 'AI will select from existing tags'"
                style="margin-top: 8px;"
              >
                {{ isGeneratingTags ? '🤖 Selecting...' : '🤖 AI Select Tags' }}
              </button>
              
              <!-- Selected Tags Display -->
              <div v-if="selectedTags.length > 0" class="selected-tags-display">
                <div 
                  v-for="tagId in selectedTags" 
                  :key="tagId" 
                  class="tag-label"
                >
                  <span class="tag-text">{{ getTagName(tagId) }}</span>
                  <button 
                    class="tag-remove" 
                    @click="removeSelectedTag(tagId)"
                    type="button"
                    title="Remove tag"
                  >
                    ×
                  </button>
                </div>
              </div>
            </div>

            <!-- Description -->
            <div class="form-group">
              <label>
                Description
                <span class="label-hint" title="Use AI to generate description automatically">
                  💡 Try AI generation →
                </span>
              </label>
              <div class="description-group">
                <textarea 
                  class="form-textarea" 
                  v-model="documentInfo.description" 
                  rows="3" 
                  placeholder="Enter document description or use AI to generate"
                ></textarea>
                <button 
                  type="button"
                  class="btn btn-secondary btn-sm btn-generate-description" 
                  @click="generateDescription"
                  :disabled="isGeneratingDescription || !documentInfo.name"
                  :title="!documentInfo.name ? 'Please enter document name first' : 'Generate description using AI'"
                >
                  {{ isGeneratingDescription ? '🤖 Generating...' : '🤖 Generate Description' }}
                </button>
              </div>
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
                <label>ISO Number</label>
                <input type="text" class="form-input" v-model="metadata.isoNumber" placeholder="ISO standard number" />
              </div>
              <div class="form-group">
                <label>Source URL</label>
                <input type="url" class="form-input" v-model="metadata.sourceUrl" placeholder="https://..." />
              </div>
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
              {{ isUploading ? (isEditMode ? 'Updating...' : 'Uploading...') : (isEditMode ? 'Update Document' : 'Upload Document') }}
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
import { useRouter, useRoute } from 'vue-router'
import Layout from '../../components/Layout.vue'
import DocumentService from '../../services/documentService.js'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// Edit mode detection
const isEditMode = ref(false)
const documentId = ref(null)

// Reactive data
const selectedFiles = ref([])
const documentInfo = ref({
  name: '',
  author: '',
  description: '',
  type: 'PDF',
  language: 'Serbian'
})

// Keywords and Tags
const currentKeyword = ref('')
const keywords = ref([])
const availableTags = ref([])
const selectedTags = ref([])

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
const isGeneratingDescription = ref(false)
const isGeneratingKeywords = ref(false)
const isGeneratingTags = ref(false)

// Users with permissions
const users = ref([])

// Computed
const filteredUsers = computed(() => {
  if (!userSearch.value) return users.value
  return users.value.filter(user => 
    user.name.toLowerCase().includes(userSearch.value.toLowerCase())
  )
})

const canUpload = computed(() => {
  // In edit mode, allow save even without new file (just updating metadata)
  if (isEditMode.value) {
    return documentInfo.value.name && documentInfo.value.author
  }
  // In upload mode, require a file
  return selectedFiles.value.length > 0 && documentInfo.value.name && documentInfo.value.author
})

// Load available tags from database on component mount
onMounted(async () => {
  try {
    // Check if we're in edit mode
    if (route.query.id) {
      isEditMode.value = true
      documentId.value = parseInt(route.query.id)
      
      // Load the document details
      const { GetDocumentByID, GetDocumentTags } = window.go.main.App
      const doc = await GetDocumentByID(documentId.value)
      
      // Populate the form with existing document data
      documentInfo.value = {
        name: doc.naziv_dokumenta || '',
        author: doc.ime_kreirao || '',
        description: doc.opis || '',
        type: doc.tip_dokumenta || 'PDF',
        language: doc.jezik || 'Serbian'
      }
      
      // Load document tags
      const docTags = await GetDocumentTags(documentId.value)
      selectedTags.value = (docTags || []).map(tag => tag.tag_id)
      
      // Load keywords (assuming they're stored in opis or a separate field)
      // For now, we'll leave keywords empty in edit mode
      keywords.value = []
    }
    
    // Load tags
    const { GetAllTags } = window.go.main.App
    const tags = await GetAllTags()
    availableTags.value = tags || []
    console.log('Loaded tags from database:', tags)
    
    // Load all users for permissions
    const { GetAllUsers } = window.go.main.App
    const allUsers = await GetAllUsers()
    
    // Transform users into the format we need with permissions
    users.value = (allUsers || []).map(user => ({
      id: user.korisnik_id,
      name: user.korisnicko_ime || `${user.ime || ''} ${user.prezime || ''}`.trim(),
      permissions: {
        read: false,
        write: false,
        delete: false
      }
    }))
    
    // If in edit mode, load existing permissions
    if (isEditMode.value) {
      const { GetDocumentPermissions } = window.go.main.App
      const permissions = await GetDocumentPermissions(documentId.value)
      
      // Update user permissions based on loaded data
      permissions.forEach(perm => {
        const user = users.value.find(u => u.id === perm.korisnik_id)
        if (user) {
          user.permissions.read = perm.moze_citati
          user.permissions.write = perm.moze_menjati
          user.permissions.delete = perm.moze_brisati
        }
      })
    }
    
    console.log('Loaded users:', users.value)
  } catch (err) {
    console.error('Error loading data:', err)
  }
})

// Keywords management
function addKeyword() {
  const keyword = currentKeyword.value.trim()
  if (keyword && !keywords.value.includes(keyword)) {
    keywords.value.push(keyword)
    currentKeyword.value = ''
  }
}

function removeKeyword(index) {
  keywords.value.splice(index, 1)
}

// Tags management
function addTagFromDb(tagId) {
  if (!selectedTags.value.includes(tagId)) {
    selectedTags.value.push(tagId)
  }
}

function removeSelectedTag(tagId) {
  selectedTags.value = selectedTags.value.filter(id => id !== tagId)
}

// Helper function to get tag name by ID
function getTagName(tagId) {
  const tag = availableTags.value.find(t => t.tag_id === tagId)
  return tag ? tag.naziv_taga : ''
}

// Methods
function handleFileSelect(event) {
  const files = Array.from(event.target.files)
  selectedFiles.value = files
  
  // Auto-fill document name from first file
  if (files.length > 0 && !documentInfo.value.name) {
    documentInfo.value.name = files[0].name.replace(/\.[^/.]+$/, "")
  }
}

async function generateDescription() {
  if (!documentInfo.value.name) {
    alert('Please enter a document name first!')
    return
  }

  isGeneratingDescription.value = true
  
  try {
    // Import GenerateDescriptionFromText from Wails
    const { GenerateDescriptionFromText } = window.go.main.App
    
    const fileName = selectedFiles.value.length > 0 ? selectedFiles.value[0].name : documentInfo.value.name
    
    const description = await GenerateDescriptionFromText(
      documentInfo.value.name,
      documentInfo.value.type,
      fileName
    )
    
    // Update description field with generated text
    if (description) {
      documentInfo.value.description = description
      alert('✅ AI Description generated successfully!')
    } else {
      alert('No description was generated. Please check your OpenAI API key configuration.')
    }
  } catch (err) {
    console.error('Error generating description:', err)
    const errMsg = err?.message || String(err)
    if (errMsg.includes('API key not configured')) {
      alert('⚠️ OpenAI API key not configured. Please set OPENAI_API_KEY environment variable and restart the application.')
    } else {
      alert('Failed to generate description: ' + errMsg)
    }
  } finally {
    isGeneratingDescription.value = false
  }
}

async function generateKeywords() {
  if (!documentInfo.value.name) {
    alert('Please enter a document name first!')
    return
  }

  isGeneratingKeywords.value = true
  
  try {
    const { GenerateTagsFromText } = window.go.main.App
    
    // Generate tags using AI (as free-form keywords)
    const tags = await GenerateTagsFromText(
      documentInfo.value.name,
      documentInfo.value.description || '',
      documentInfo.value.type,
      10 // More keywords since this is free-form
    )
    
    if (tags && tags.length > 0) {
      // Add generated keywords to the array
      tags.forEach(tag => {
        if (!keywords.value.includes(tag)) {
          keywords.value.push(tag)
        }
      })
      alert('✅ AI Keywords generated successfully!')
    } else {
      alert('No keywords were generated. Please check your OpenAI API key configuration.')
    }
  } catch (err) {
    console.error('Error generating keywords:', err)
    const errMsg = err?.message || String(err)
    if (errMsg.includes('API key not configured')) {
      alert('⚠️ OpenAI API key not configured. Please set OPENAI_API_KEY environment variable and restart the application.')
    } else {
      alert('Failed to generate keywords: ' + errMsg)
    }
  } finally {
    isGeneratingKeywords.value = false
  }
}

async function generateTags() {
  if (!documentInfo.value.name) {
    alert('Please enter a document name first!')
    return
  }

  if (availableTags.value.length === 0) {
    alert('No tags available in database. Please add tags first.')
    return
  }

  isGeneratingTags.value = true
  
  try {
    const { GenerateTagsFromText } = window.go.main.App
    
    // Create a text that includes available tags for AI context
    const availableTagNames = availableTags.value.map(t => t.naziv_taga).join(', ')
    const contextText = `Document: ${documentInfo.value.name}
Description: ${documentInfo.value.description || 'N/A'}
Type: ${documentInfo.value.type}

Available tags in database: ${availableTagNames}

Please select ONLY from the available tags listed above.`
    
    // Generate tags using AI
    const aiSuggestedTags = await GenerateTagsFromText(
      contextText,
      '',
      documentInfo.value.type,
      5
    )
    
    if (aiSuggestedTags && aiSuggestedTags.length > 0) {
      // Match AI suggestions with existing tags (case-insensitive)
      const matchedTagIds = []
      
      aiSuggestedTags.forEach(suggestedTag => {
        const matchedTag = availableTags.value.find(
          t => t.naziv_taga.toLowerCase().trim() === suggestedTag.toLowerCase().trim()
        )
        if (matchedTag && !selectedTags.value.includes(matchedTag.tag_id)) {
          matchedTagIds.push(matchedTag.tag_id)
        }
      })
      
      if (matchedTagIds.length > 0) {
        selectedTags.value = [...selectedTags.value, ...matchedTagIds]
        alert(`✅ AI selected ${matchedTagIds.length} tag(s) from database!`)
      } else {
        alert('⚠️ AI could not match suggested tags with existing database tags. Please select manually.')
      }
    } else {
      alert('No tags were generated. Please check your OpenAI API key configuration.')
    }
  } catch (err) {
    console.error('Error generating tags:', err)
    const errMsg = err?.message || String(err)
    if (errMsg.includes('API key not configured')) {
      alert('⚠️ OpenAI API key not configured. Please set OPENAI_API_KEY environment variable and restart the application.')
    } else {
      alert('Failed to generate tags: ' + errMsg)
    }
  } finally {
    isGeneratingTags.value = false
  }
}

async function handleSave() {
  // In edit mode, file selection is optional (only update if new file is selected)
  if (!isEditMode.value && selectedFiles.value.length === 0) {
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
    // Convert selected tag IDs to tag names for backend
    const tagNames = selectedTags.value.map(tagId => {
      const tag = availableTags.value.find(t => t.tag_id === tagId)
      return tag ? tag.naziv_taga : null
    }).filter(name => name !== null)

    // Prepare document data
    const documentData = {
      name: documentInfo.value.name,
      description: documentInfo.value.description,
      type: documentInfo.value.type,
      language: documentInfo.value.language,
      tags: tagNames,  // Tags from database (as names)
      keywords: keywords.value.join(', ') || '',  // Keywords as comma-separated string
      projectId: null, // TODO: Add project selection
      folderId: null   // TODO: Add folder selection
    }

    // Simulate progress for better UX
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += Math.random() * 20
      }
    }, 200)

    let docId = documentId.value

    if (isEditMode.value) {
      // Update existing document
      const { UpdateDocument } = window.go.main.App
      
      // Create request object matching backend expectations
      const updateRequest = {
        naziv_dokumenta: documentData.name,
        opis: documentData.description,
        tip_dokumenta: documentData.type,
        jezik: documentData.language,
        kljucne_reci: documentData.keywords,
        tagovi: documentData.tags,
        projekat_id: documentData.projectId,
        folder_id: documentData.folderId
      }
      
      // If new file is selected, include it
      if (selectedFiles.value.length > 0) {
        const file = selectedFiles.value[0]
        updateRequest.fajl_sadrzaj = await file.arrayBuffer()
        updateRequest.putanja_fajla = file.name
      }
      
      await UpdateDocument(docId, updateRequest)
      console.log('Document updated with ID:', docId)
    } else {
      // Upload new document
      const file = selectedFiles.value[0]
      docId = await DocumentService.uploadDocument(documentData, file)
      console.log('Document uploaded with ID:', docId)
    }

    // Save permissions for users who have any permission set
    const { SetDocumentPermission, RemoveDocumentPermission } = window.go.main.App
    
    for (const user of users.value) {
      // Check if user has any permissions enabled
      const hasPermissions = user.permissions.read || user.permissions.write || user.permissions.delete
      
      if (hasPermissions) {
        try {
          await SetDocumentPermission({
            dokument_id: docId,
            korisnik_id: user.id,
            moze_citati: user.permissions.read,
            moze_menjati: user.permissions.write,
            moze_brisati: user.permissions.delete
          })
          console.log(`Saved permissions for user ${user.name}:`, user.permissions)
        } catch (permErr) {
          console.error(`Error saving permissions for user ${user.name}:`, permErr)
        }
      } else if (isEditMode.value) {
        // In edit mode, if no permissions are set, remove existing permissions
        try {
          await RemoveDocumentPermission(docId, user.id)
          console.log(`Removed permissions for user ${user.name}`)
        } catch (permErr) {
          // Ignore errors when removing non-existent permissions
          console.log(`No permissions to remove for user ${user.name}`)
        }
      }
    }
    
    clearInterval(progressInterval)
    uploadProgress.value = 100

    // Log activity
    try {
      const { LogActivity } = window.go.main.App
      await LogActivity({
        tip_aktivnosti: isEditMode.value ? 'EDIT' : 'UPLOAD',
        entitet_tip: 'DOKUMENT',
        entitet_id: docId,
        naziv_entiteta: documentInfo.value.name,
        opis: `${isEditMode.value ? 'Updated' : 'Uploaded'} document: ${documentInfo.value.name}`,
        rezultat: 'SUCCESS'
      })
    } catch (logErr) {
      console.error('Failed to log activity:', logErr)
      // Don't fail the whole operation if logging fails
    }

    setTimeout(() => {
      isUploading.value = false
      alert(isEditMode.value ? 'Document updated successfully!' : 'Document uploaded successfully!')
      router.push('/documents')
    }, 500)

  } catch (err) {
    isUploading.value = false
    uploadProgress.value = 0
    error.value = err.message
    alert((isEditMode.value ? 'Update' : 'Upload') + ' failed: ' + err.message)
    console.error((isEditMode.value ? 'Update' : 'Upload') + ' error:', err)
  }
}
</script>

<style scoped>
.label-hint {
  font-size: 0.85em;
  color: #6c757d;
  font-weight: normal;
  margin-left: 8px;
  cursor: help;
}

.label-hint:hover {
  color: #007bff;
}

.btn-sm:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-sm:not(:disabled):hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 123, 255, 0.3);
}

/* Keywords Input Group */
.keywords-input-group {
  display: flex;
  gap: 8px;
  align-items: stretch;
}

.keywords-input-group .form-input {
  flex: 1;
}

/* Description Group */
.description-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn-generate-description {
  align-self: flex-end;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
}

.btn-generate-description:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-generate-description:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-add-tag {
  min-width: 80px;
  white-space: nowrap;
  background: #28a745;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-add-tag:hover {
  background: #218838;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(40, 167, 69, 0.3);
}

/* Tags Container */
.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

/* Tag Label */
.tag-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 500;
  box-shadow: 0 2px 4px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
  animation: tagAppear 0.3s ease;
}

.tag-label:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(102, 126, 234, 0.4);
}

.tag-text {
  user-select: none;
}

/* Tag Remove Button */
.tag-remove {
  background: rgba(255, 255, 255, 0.3);
  border: none;
  color: white;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
  padding: 0;
  transition: all 0.2s ease;
}

.tag-remove:hover {
  background: rgba(255, 255, 255, 0.5);
  transform: rotate(90deg);
}

/* Tag Appear Animation */
@keyframes tagAppear {
  0% {
    opacity: 0;
    transform: scale(0.8) translateY(-10px);
  }
  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

/* Empty state */
.tags-container:empty {
  display: none;
}

/* Tags Select Container */
.tags-select-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tags-multiselect {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  font-size: 14px;
}

.tags-multiselect option {
  padding: 6px 8px;
}

.tags-multiselect option:checked {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-generate-keywords,
.btn-generate-tags {
  align-self: flex-end;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
}

.btn-generate-keywords:hover:not(:disabled),
.btn-generate-tags:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(240, 147, 251, 0.4);
}

.btn-generate-keywords:disabled,
.btn-generate-tags:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Selected Tags Display */
.selected-tags-display {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 4px;
  min-height: 40px;
}

/* Available Tags List */
.available-tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
  max-height: 200px;
  overflow-y: auto;
}

.tag-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 20px;
  transition: all 0.2s ease;
}

.tag-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.tag-name {
  font-size: 14px;
  color: #495057;
  font-weight: 500;
}

.btn-xs {
  padding: 2px 8px;
  font-size: 12px;
  min-width: 30px;
  height: 24px;
  border-radius: 12px;
}

.btn-add-tag-from-db {
  background: #28a745;
  color: white;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: bold;
}

.btn-add-tag-from-db:hover:not(:disabled) {
  background: #218838;
  transform: scale(1.1);
}

.btn-add-tag-from-db:disabled {
  background: #6c757d;
  cursor: not-allowed;
  opacity: 0.6;
}
</style>

