<template>
  <Layout>
    <div class="projects">
      <!-- Header -->
      <div class="page-header">
        <div>
          <h2>Projekti</h2>
          <div class="breadcrumb">Početna > Projekti</div>
        </div>
        <button class="btn btn-primary" @click="showCreateModal = true">
          <span class="btn-icon">➕</span>
          Novi projekat
        </button>
      </div>
      
      <!-- Filters -->
      <div class="filters card">
        <div class="filter-group">
          <label>Status:</label>
          <select v-model="filters.status">
            <option value="">Svi statusi</option>
            <option value="active">Aktivni</option>
            <option value="completed">Završeni</option>
            <option value="on-hold">Na čekanju</option>
            <option value="cancelled">Otkazani</option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>Pretraga:</label>
          <input 
            type="text" 
            v-model="filters.search" 
            placeholder="Pretraži projekte..."
            class="search-input"
          >
        </div>
        
        <div class="filter-group">
          <label>Sortiranje:</label>
          <select v-model="filters.sortBy">
            <option value="name">Naziv</option>
            <option value="created">Datum kreiranja</option>
            <option value="updated">Poslednja izmena</option>
            <option value="progress">Napredak</option>
          </select>
        </div>
      </div>
      
      <!-- Loading State -->
      <div v-if="loading && projects.length === 0" class="loading-state">
        <div class="spinner"></div>
        <p>Učitavanje projekata...</p>
      </div>

      <!-- Error State -->
      <div v-if="error && !loading" class="error-state">
        <div class="error-icon">⚠️</div>
        <h3>Greška pri učitavanju</h3>
        <p>{{ error }}</p>
        <button class="btn btn-primary" @click="loadProjects">
          Pokušaj ponovo
        </button>
      </div>

      <!-- Projects Grid -->
      <div v-if="!loading || projects.length > 0" class="projects-grid">
        <div 
          v-for="project in filteredProjects" 
          :key="project.id"
          class="project-card"
          @click="selectProject(project)"
        >
          <div class="project-header">
            <h3>{{ project.name }}</h3>
            <div class="project-actions">
              <button 
                class="btn-icon-small" 
                @click.stop="editProject(project)"
                title="Uredi"
              >
                ✏️
              </button>
              <button 
                class="btn-icon-small danger" 
                @click.stop="deleteProject(project)"
                title="Obriši"
              >
                🗑️
              </button>
            </div>
          </div>
          
          <div class="project-description">
            {{ project.description }}
          </div>
          
          <div class="project-progress">
            <div class="progress-container">
              <div 
                class="progress-bar" 
                :style="{ width: project.progress + '%' }"
              ></div>
            </div>
            <span class="progress-text">{{ project.progress }}%</span>
          </div>
          
          <div class="project-info">
            <div class="project-status">
              <span 
                class="status-badge" 
                :class="`status-${project.status}`"
              >
                {{ getStatusText(project.status) }}
              </span>
            </div>
            
            <div class="project-meta">
              <div class="project-leader">
                <span class="label">Vođa:</span>
                <span>{{ project.leader }}</span>
              </div>
              <div class="project-deadline">
                <span class="label">Deadline:</span>
                <span>{{ formatDate(project.deadline) }}</span>
              </div>
            </div>
          </div>

          <div class="project-footer">
            <button class="btn btn-secondary" @click.stop="openDocumentation(project)">
              📚 Dokumentacija
            </button>
          </div>
          
          <div class="project-team">
            <div class="team-avatars">
              <div 
                v-for="member in project.team.slice(0, 3)" 
                :key="member.id"
                class="avatar"
                :title="member.name"
              >
                {{ member.name.charAt(0) }}
              </div>
              <div v-if="project.team.length > 3" class="avatar more">
                +{{ project.team.length - 3 }}
              </div>
            </div>
          </div>
        </div>
        
        <!-- Empty State -->
        <div v-if="filteredProjects.length === 0" class="empty-state">
          <div class="empty-icon">📁</div>
          <h3>Nema projekata</h3>
          <p>Kreirajte prvi projekat klikom na dugme "Novi projekat"</p>
        </div>
      </div>
      
      <!-- Create/Edit Project Modal -->
      <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click="closeModals">
        <div class="modal modal-large" @click.stop>
          <div class="modal-header">
            <h3 class="modal-title">
              {{ showCreateModal ? 'Kreiranje novog projekta' : 'Uređivanje projekta' }}
            </h3>
            <button class="modal-close" @click="closeModals">×</button>
          </div>
          
          <form @submit.prevent="saveProject">
            <div class="form-row-two-cols">
              <!-- Left Column -->
              <div class="form-column">
                <div class="form-group">
                  <label>Naziv projekta</label>
                  <input 
                    type="text" 
                    v-model="projectForm.name" 
                    placeholder="Unesite naziv projekta" 
                    required
                  >
                </div>
                
                <div class="form-group">
                  <label>Opis projekta</label>
                  <textarea 
                    v-model="projectForm.description" 
                    placeholder="Unesite opis projekta"
                    rows="4"
                  ></textarea>
                </div>
                
                <div class="form-row">
                  <div class="form-group">
                    <label>Datum početka</label>
                    <input 
                      type="date" 
                      v-model="projectForm.startDate"
                    >
                  </div>
                  
                  <div class="form-group">
                    <label>Datum završetka</label>
                    <input 
                      type="date" 
                      v-model="projectForm.endDate"
                    >
                  </div>
                </div>
                
                <div class="form-group">
                  <label>Radni tok</label>
                  <div class="workflow-section">
                    <div class="workflow-selector">
                      <select v-model="projectForm.workflowId" @change="onWorkflowChange">
                        <option value="">Izaberite radni tok</option>
                        <option 
                          v-for="workflow in workflows" 
                          :key="workflow.radni_tok_id" 
                          :value="workflow.radni_tok_id"
                        >
                          {{ workflow.naziv }}
                        </option>
                      </select>
                      <button 
                        type="button" 
                        class="btn-create-workflow"
                        @click="createNewWorkflow"
                        title="Kreiraj novi radni tok"
                      >
                        ➕ Novi tok
                      </button>
                    </div>
                    <div class="workflow-actions">
                      <button 
                        type="button" 
                        class="btn-add-phase"
                        :disabled="!projectForm.workflowId"
                        @click="addPhase"
                      >
                        Dodaj Fazu
                      </button>
                    </div>
                    
                    <div v-if="projectForm.phases.length > 0" class="phases-list">
                      <div class="section-label">Faze projekta:</div>
                      <div 
                        v-for="(phase, index) in projectForm.phases" 
                        :key="index"
                        class="phase-item"
                      >
                        <span class="phase-name">{{ phase.naziv_faze }}</span>
                        <button 
                          type="button" 
                          class="btn-remove"
                          @click="removePhase(index)"
                        >
                          Ukloni
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
                
                <div class="form-group">
                  <label>Resursi</label>
                  <input 
                    type="text" 
                    v-model="projectForm.resources" 
                    placeholder="Unesite potrebne resurse"
                  >
                </div>
              </div>
              
              <!-- Right Column -->
              <div class="form-column">
                <div class="form-group">
                  <label>Članovi tima</label>
                  <div class="team-section">
                    <div class="team-dropdown">
                      <select 
                        v-model="selectedUserId"
                        @change="addTeamMemberFromDropdown"
                        class="team-select"
                      >
                        <option value="">Pretraži i dodaj članove tima</option>
                        <option 
                          v-for="user in availableUsers" 
                          :key="user.id"
                          :value="user.id"
                        >
                          {{ user.name }} ({{ user.role }})
                        </option>
                      </select>
                    </div>
                    
                    <div class="section-label">Izabrani članovi:</div>
                    <div v-if="projectForm.teamMembers.length === 0" class="empty-team">
                      Nema izabranih članova
                    </div>
                    <div v-else class="selected-team-members">
                      <div 
                        v-for="member in projectForm.teamMembers" 
                        :key="member.id"
                        class="team-member-item"
                      >
                        <div class="user-info">
                          <div class="user-avatar">{{ member.name.charAt(0) }}</div>
                          <div>
                            <div class="user-name">{{ member.name }}</div>
                            <div class="user-role">{{ member.role }}</div>
                          </div>
                        </div>
                        <button 
                          type="button" 
                          class="btn-remove"
                          @click="removeTeamMember(member.id)"
                        >
                          Ukloni
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <div class="form-actions">
              <button type="button" class="btn btn-secondary" @click="closeModals">
                Otkaži
              </button>
              <button type="submit" class="btn btn-primary">
                {{ showCreateModal ? 'Kreiraj projekat' : 'Sačuvaj izmene' }}
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
import { useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'

// Reactive data
const showCreateModal = ref(false)
const showEditModal = ref(false)
const selectedProject = ref(null)

const filters = ref({
  status: '',
  search: '',
  sortBy: 'name'
})

const projectForm = ref({
  name: '',
  description: '',
  startDate: '',
  endDate: '',
  status: 'aktivan',
  workflowId: '',
  phases: [],
  teamMembers: [],
  resources: ''
})

const selectedUserId = ref('')

import { 
  fetchUserProjects, 
  createProject,
  updateProject,
  deleteProject as deleteProjectAPI,
  fetchProjectDetails
} from '../services/projectService.js'
import { 
  GetAllUsers, 
  GetAvailableTeamMembers,
  GetAllWorkflows, 
  GetWorkflowPhases, 
  CreateNewProject,
  CreatePhase,
  CreateWorkflow,
  DeletePhase
} from '../../wailsjs/go/main/App.js'

const projects = ref([])
const users = ref([])
const workflows = ref([])
const loading = ref(false)
const error = ref(null)

const router = useRouter()

// Computed
const filteredProjects = computed(() => {
  let filtered = projects.value

  // Filter by status
  if (filters.value.status) {
    filtered = filtered.filter(p => p.status === filters.value.status)
  }

  // Filter by search
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(p => 
      p.name.toLowerCase().includes(search) || 
      p.description.toLowerCase().includes(search)
    )
  }

  // Sort
  filtered.sort((a, b) => {
    switch (filters.value.sortBy) {
      case 'name':
        return a.name.localeCompare(b.name)
      case 'created':
        return new Date(b.created) - new Date(a.created)
      case 'updated':
        return new Date(b.updated) - new Date(a.updated)
      case 'progress':
        return b.progress - a.progress
      default:
        return 0
    }
  })

  return filtered
})

// Available users (exclude already selected team members)
const availableUsers = computed(() => {
  const selectedIds = new Set(projectForm.value.teamMembers.map(m => m.id))
  return users.value.filter(user => !selectedIds.has(user.id))
})

// Methods
function getStatusText(status) {
  const statusMap = {
    'active': 'Aktivni',
    'completed': 'Završen',
    'on-hold': 'Na čekanju',
    'cancelled': 'Otkazan'
  }
  return statusMap[status] || status
}

function formatDate(dateString) {
  if (!dateString) return '-'
  const d = new Date(dateString)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('sr-RS')
}

function selectProject(project) {
  selectedProject.value = project
  // Navigate to project details or show details modal
  console.log('Selected project:', project)
}

function editProject(project) {
  selectedProject.value = project
  projectForm.value = {
    name: project.name,
    description: project.description,
    startDate: project.startDate || '',
    endDate: project.deadline || '',
    status: project.status,
    workflowId: project.workflowId || '',
    phases: [],
    teamMembers: project.team || [],
    resources: project.resources || ''
  }
  showEditModal.value = true
}

// Team management methods
function addTeamMemberFromDropdown() {
  if (!selectedUserId.value) return
  
  const user = users.value.find(u => u.id === parseInt(selectedUserId.value))
  if (user && !projectForm.value.teamMembers.find(m => m.id === user.id)) {
    projectForm.value.teamMembers.push({ ...user })
  }
  
  // Reset dropdown
  selectedUserId.value = ''
}

function removeTeamMember(userId) {
  projectForm.value.teamMembers = projectForm.value.teamMembers.filter(m => m.id !== userId)
}

// Workflow and phase management
async function onWorkflowChange() {
  if (!projectForm.value.workflowId) {
    projectForm.value.phases = []
    return
  }
  
  try {
    const phases = await GetWorkflowPhases(projectForm.value.workflowId)
    projectForm.value.phases = phases || []
  } catch (err) {
    console.error('Error loading workflow phases:', err)
    projectForm.value.phases = []
  }
}

async function addPhase() {
  if (!projectForm.value.workflowId) {
    alert('Molimo prvo izaberite radni tok')
    return
  }
  
  const phaseName = prompt('Unesite naziv faze:')
  if (!phaseName || !phaseName.trim()) {
    return
  }
  
  try {
    loading.value = true
    
    // Create phase in database
    const newPhase = {
      radni_tok_id: parseInt(projectForm.value.workflowId),
      naziv_faze: phaseName.trim(),
      redosled: projectForm.value.phases.length + 1
    }
    
    await CreatePhase(newPhase)
    
    // Reload phases to get the new phase with its ID
    await onWorkflowChange()
    
    alert('Faza uspešno kreirana!')
  } catch (err) {
    console.error('Error creating phase:', err)
    alert('Greška pri kreiranju faze: ' + err.message)
  } finally {
    loading.value = false
  }
}

async function removePhase(index) {
  const phase = projectForm.value.phases[index]
  
  if (!confirm('Da li ste sigurni da želite da uklonite ovu fazu?')) {
    return
  }
  
  if (phase.faza_id) {
    // Phase exists in database - delete it
    try {
      loading.value = true
      
      await DeletePhase(phase.faza_id)
      
      // Reload phases to reflect database state
      await onWorkflowChange()
      
      alert('Faza uspešno uklonjena!')
    } catch (err) {
      console.error('Error deleting phase:', err)
      alert('Greška pri brisanju faze: ' + (err.message || err))
    } finally {
      loading.value = false
    }
  } else {
    // Phase doesn't exist in database yet, just remove from UI
    projectForm.value.phases.splice(index, 1)
    // Update order
    projectForm.value.phases.forEach((phase, idx) => {
      phase.redosled = idx + 1
    })
  }
}

async function createNewWorkflow() {
  const workflowName = prompt('Unesite naziv novog radnog toka:')
  if (!workflowName || !workflowName.trim()) {
    return
  }
  
  // Use a select-style prompt or validate input
  const workflowTypeInput = prompt('Izaberite tip radnog toka:\n1 - PROJEKAT\n2 - DOKUMENTACIJA\n\nUnesite broj (1 ili 2):', '1')
  if (!workflowTypeInput) {
    return
  }
  
  // Map user input to valid database values
  let workflowType
  if (workflowTypeInput === '1' || workflowTypeInput.toUpperCase() === 'PROJEKAT') {
    workflowType = 'PROJEKAT'
  } else if (workflowTypeInput === '2' || workflowTypeInput.toUpperCase() === 'DOKUMENTACIJA') {
    workflowType = 'DOKUMENTACIJA'
  } else {
    alert('Nevažeći tip radnog toka. Koristite "PROJEKAT" ili "DOKUMENTACIJA".')
    return
  }
  
  try {
    loading.value = true
    
    // Create workflow in database
    const newWorkflow = {
      radni_tok_id: 0, // Will be assigned by database
      naziv: workflowName.trim(),
      tip_toka: workflowType,
      opis: '',
      da_li_je_sablon: false
    }
    
    await CreateWorkflow(newWorkflow)
    
    // Reload workflows
    await loadWorkflows()
    
    // Select the newly created workflow (find by name)
    const createdWorkflow = workflows.value.find(w => w.naziv === workflowName.trim())
    if (createdWorkflow) {
      projectForm.value.workflowId = createdWorkflow.radni_tok_id.toString()
      await onWorkflowChange()
    }
    
    alert('Radni tok uspešno kreiran!')
  } catch (err) {
    console.error('Error creating workflow:', err)
    alert('Greška pri kreiranju radnog toka: ' + (err.message || err))
  } finally {
    loading.value = false
  }
}

async function deleteProject(project) {
  if (confirm(`Da li ste sigurni da želite da obrišete projekat "${project.name}"?`)) {
    try {
      loading.value = true
      await deleteProjectAPI(project.id)
      // Refresh projects list
      await loadProjects()
    } catch (err) {
      console.error('Error deleting project:', err)
      alert('Greška pri brisanju projekta: ' + err.message)
    } finally {
      loading.value = false
    }
  }
}

function openDocumentation(project) {
  router.push(`/projects/${project.id}/documents`)
}

async function saveProject() {
  try {
    loading.value = true
    error.value = null
    
    // Validation
    if (!projectForm.value.name || !projectForm.value.name.trim()) {
      alert('Naziv projekta je obavezan!')
      return
    }
    
    if (showCreateModal.value) {
      // Create new project via backend using Wails
      const projectData = {
        naziv_projekta: projectForm.value.name.trim(),
        opis: projectForm.value.description || '',
        datum_pocetka: projectForm.value.startDate ? new Date(projectForm.value.startDate).toISOString() : null,
        datum_zavrsetka: projectForm.value.endDate ? new Date(projectForm.value.endDate).toISOString() : null,
        radni_tok_id: projectForm.value.workflowId ? parseInt(projectForm.value.workflowId) : null,
        clanovi_tima: projectForm.value.teamMembers.map(m => m.id)
      }
      
      console.log('Creating project with data:', projectData)
      
      await CreateNewProject(projectData)
      
      alert(`Projekat "${projectForm.value.name}" uspešno kreiran!\n` +
            `Radni tok: ${projectForm.value.workflowId ? 'Da' : 'Nema'}\n` +
            `Broj članova tima: ${projectForm.value.teamMembers.length}`)
    } else if (showEditModal.value && selectedProject.value) {
      // Update existing project
      await updateProject(selectedProject.value.id, {
        name: projectForm.value.name,
        description: projectForm.value.description,
        startDate: projectForm.value.startDate,
        endDate: projectForm.value.endDate,
        status: projectForm.value.status,
        workflowId: projectForm.value.workflowId,
        teamMembers: projectForm.value.teamMembers.map(m => m.id),
        resources: projectForm.value.resources
      })
      
      alert('Projekat uspešno ažuriran!')
    }
    
    // Refresh projects list
    await loadProjects()
    closeModals()
  } catch (err) {
    console.error('Error saving project:', err)
    error.value = err.message
    alert('Greška pri čuvanju projekta: ' + (err.message || err))
  } finally {
    loading.value = false
  }
}

function closeModals() {
  showCreateModal.value = false
  showEditModal.value = false
  selectedProject.value = null
  selectedUserId.value = ''
  projectForm.value = {
    name: '',
    description: '',
    startDate: '',
    endDate: '',
    status: 'aktivan',
    workflowId: '',
    phases: [],
    teamMembers: [],
    resources: ''
  }
}

// Load projects from backend
async function loadProjects() {
  try {
    loading.value = true
    error.value = null
    const data = await fetchUserProjects()
    projects.value = data
    
    // Load team members for each project
    for (const project of projects.value) {
      try {
        const details = await fetchProjectDetails(project.id)
        project.team = details.team
      } catch (err) {
        console.error(`Error loading team for project ${project.id}:`, err)
        project.team = []
      }
    }
  } catch (e) {
    console.error('Neuspešno učitavanje projekata:', e)
    error.value = e.message
  } finally {
    loading.value = false
  }
}

// Load users for dropdown
async function loadUsers() {
  try {
    // Get available team members (researchers only)
    const userData = await GetAvailableTeamMembers()
    
    users.value = (userData || []).map(u => ({
      id: u.korisnik_id,
      name: `${u.ime || ''} ${u.prezime || ''}`.trim() || u.korisnicko_ime,
      username: u.korisnicko_ime,
      role: u.naziv_uloge
    }))
  } catch (e) {
    console.error('Neuspešno učitavanje istraživača:', e)
  }
}

// Load workflows for dropdown
async function loadWorkflows() {
  try {
    const workflowData = await GetAllWorkflows()
    workflows.value = workflowData || []
  } catch (e) {
    console.error('Neuspešno učitavanje radnih tokova:', e)
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    loadProjects(),
    loadUsers(),
    loadWorkflows()
  ])
})
</script>

<style scoped>
.projects {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.8rem;
  color: #1a1a1a;
}

.breadcrumb {
  color: #666;
  font-size: 0.9rem;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  margin-bottom: 2rem;
  padding: 1.5rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1 1 200px;
  min-width: 200px;
}

.filter-group label {
  font-size: 0.875rem;
  font-weight: 600;
  color: #555;
}

.filter-group select,
.filter-group input {
  padding: 0.5rem 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 0.9rem;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.project-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  height: fit-content;
}

.project-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  gap: 1rem;
}

.project-header h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #1a1a1a;
  flex: 1;
  line-height: 1.4;
}

.project-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-icon-small {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  font-size: 1rem;
  transition: transform 0.2s;
}

.btn-icon-small:hover {
  transform: scale(1.2);
}

.btn-icon-small.danger:hover {
  filter: brightness(0.8);
}

.project-description {
  font-size: 0.9rem;
  color: #666;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 2.7em;
}

.project-progress {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.progress-container {
  flex: 1;
  height: 8px;
  background: #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #4CAF50, #66BB6A);
  transition: width 0.3s ease;
  border-radius: 4px;
}

.progress-text {
  font-size: 0.85rem;
  font-weight: 600;
  color: #4CAF50;
  min-width: 38px;
}

.project-info {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.project-status {
  display: flex;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-active {
  background: #e8f5e9;
  color: #2e7d32;
}

.status-completed {
  background: #e3f2fd;
  color: #1565c0;
}

.status-on-hold {
  background: #fff3e0;
  color: #e65100;
}

.status-cancelled {
  background: #ffebee;
  color: #c62828;
}

.project-meta {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  font-size: 0.85rem;
}

.project-leader,
.project-deadline {
  display: flex;
  gap: 0.5rem;
}

.project-meta .label {
  font-weight: 600;
  color: #666;
  min-width: 70px;
}

.project-footer {
  padding-top: 0.5rem;
  border-top: 1px solid #f0f0f0;
}

.project-footer .btn {
  width: 100%;
  justify-content: center;
}

.project-team {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding-top: 0.5rem;
  border-top: 1px solid #f0f0f0;
}

.team-avatars {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.avatar.more {
  background: #f0f0f0;
  color: #666;
  font-size: 0.7rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-primary {
  background: #0066cc;
  color: white;
}

.btn-primary:hover {
  background: #0052a3;
}

.btn-secondary {
  background: #f5f5f5;
  color: #333;
}

.btn-secondary:hover {
  background: #e0e0e0;
}

.btn-icon {
  font-size: 1rem;
}

.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #0066cc;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-state p,
.error-state p {
  color: #666;
  margin-top: 0.5rem;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.error-state h3 {
  color: #d32f2f;
  margin-bottom: 0.5rem;
}

.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h3 {
  margin-bottom: 0.5rem;
  color: #333;
}

.empty-state p {
  color: #666;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.modal-large {
  max-width: 1000px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e0e0e0;
}

.modal-title {
  margin: 0;
  font-size: 1.4rem;
  color: #1a1a1a;
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.modal-close:hover {
  background: #f5f5f5;
  color: #333;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #333;
  font-size: 0.9rem;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 0.95rem;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #0066cc;
  box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e0e0e0;
}

/* Responsive */
@media (max-width: 768px) {
  .projects {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .filters {
    flex-direction: column;
    gap: 1rem;
  }

  .projects-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .modal {
    padding: 1.5rem;
  }
  
  .form-row-two-cols {
    grid-template-columns: 1fr;
  }
}

/* Two column form layout */
.form-row-two-cols {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 1rem;
}

.form-column {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

/* Workflow section */
.workflow-section {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.workflow-selector {
  display: flex;
  gap: 0.5rem;
}

.workflow-selector select {
  flex: 1;
}

.workflow-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.btn-create-workflow {
  padding: 0.5rem 1rem;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.2s;
}

.btn-create-workflow:hover {
  background: #218838;
}

.btn-add-phase {
  padding: 0.5rem 1rem;
  background: #0066cc;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.2s;
  flex: 1;
}

.btn-add-phase:hover:not(:disabled) {
  background: #0052a3;
}

.btn-add-phase:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.phases-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.phase-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #f5f5f5;
  border-radius: 6px;
  border: 1px solid #e0e0e0;
}

.phase-name {
  font-size: 0.9rem;
  color: #333;
  font-weight: 500;
}

.btn-remove {
  padding: 0.25rem 0.75rem;
  background: #ff4444;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-remove:hover {
  background: #cc0000;
}

/* Team section */
.team-section {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.team-dropdown {
  width: 100%;
}

.team-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 0.9rem;
  background: white;
  cursor: pointer;
  transition: border-color 0.2s;
}

.team-select:hover {
  border-color: #0066cc;
}

.team-select:focus {
  outline: none;
  border-color: #0066cc;
  box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
}

.section-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #555;
  margin-top: 0.5rem;
}

.empty-team {
  padding: 1rem;
  text-align: center;
  color: #999;
  font-size: 0.9rem;
  background: #f9f9f9;
  border-radius: 6px;
  border: 1px dashed #ddd;
}

.selected-team-members {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.team-member-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #f5f5f5;
  border-radius: 6px;
  border: 1px solid #e0e0e0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  font-weight: 600;
  text-transform: uppercase;
}

.user-name {
  font-size: 0.9rem;
  font-weight: 600;
  color: #333;
}

.user-role {
  font-size: 0.75rem;
  color: #666;
}
</style>

