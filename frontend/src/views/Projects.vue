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
      
      <!-- Projects Grid -->
      <div class="projects-grid">
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
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h3 class="modal-title">
              {{ showCreateModal ? 'Kreiranje novog projekta' : 'Uređivanje projekta' }}
            </h3>
            <button class="modal-close" @click="closeModals">×</button>
          </div>
          
          <form @submit.prevent="saveProject">
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
              <label>Opis</label>
              <textarea 
                v-model="projectForm.description" 
                placeholder="Unesite opis projekta"
                rows="3"
              ></textarea>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label>Vođa projekta</label>
                <select v-model="projectForm.leaderId" required>
                  <option value="">Izaberite vođu</option>
                  <option 
                    v-for="user in users" 
                    :key="user.id" 
                    :value="user.id"
                  >
                    {{ user.name }}
                  </option>
                </select>
              </div>
              
              <div class="form-group">
                <label>Status</label>
                <select v-model="projectForm.status">
                  <option value="active">Aktivni</option>
                  <option value="on-hold">Na čekanju</option>
                  <option value="completed">Završen</option>
                  <option value="cancelled">Otkazan</option>
                </select>
              </div>
            </div>
            
            <div class="form-group">
              <label>Deadline</label>
              <input 
                type="date" 
                v-model="projectForm.deadline"
                required
              >
            </div>
            
            <div class="form-actions">
              <button type="button" class="btn btn-secondary" @click="closeModals">
                Otkaži
              </button>
              <button type="submit" class="btn btn-primary">
                {{ showCreateModal ? 'Kreiraj' : 'Sačuvaj' }}
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
  leaderId: '',
  status: 'active',
  deadline: ''
})

// Mock data
const projects = ref([
  {
    id: 1,
    name: 'Web Portal Refactoring',
    description: 'Modernizacija postojećeg web portala institucije',
    progress: 75,
    status: 'active',
    leader: 'Marko Petrović',
    leaderId: 1,
    deadline: '2024-12-31',
    created: '2024-01-15',
    updated: '2024-09-01',
    team: [
      { id: 1, name: 'Marko Petrović' },
      { id: 2, name: 'Ana Jovanović' },
      { id: 3, name: 'Stefan Nikolić' }
    ]
  },
  {
    id: 2,
    name: 'AI Analitika Modula',
    description: 'Implementacija AI za analizu istraživačkih podataka',
    progress: 45,
    status: 'active',
    leader: 'Ana Jovanović',
    leaderId: 2,
    deadline: '2025-03-15',
    created: '2024-03-01',
    updated: '2024-08-20',
    team: [
      { id: 2, name: 'Ana Jovanović' },
      { id: 4, name: 'Milica Stojković' }
    ]
  },
  {
    id: 3,
    name: 'Mobile Aplikacija',
    description: 'Razvoj mobilne aplikacije za iOS i Android',
    progress: 20,
    status: 'on-hold',
    leader: 'Stefan Nikolić',
    leaderId: 3,
    deadline: '2025-06-01',
    created: '2024-05-10',
    updated: '2024-07-15',
    team: [
      { id: 3, name: 'Stefan Nikolić' },
      { id: 5, name: 'Jovana Mitrović' }
    ]
  }
])

const users = ref([
  { id: 1, name: 'Marko Petrović' },
  { id: 2, name: 'Ana Jovanović' },
  { id: 3, name: 'Stefan Nikolić' },
  { id: 4, name: 'Milica Stojković' },
  { id: 5, name: 'Jovana Mitrović' }
])

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
  return new Date(dateString).toLocaleDateString('sr-RS')
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
    leaderId: project.leaderId,
    status: project.status,
    deadline: project.deadline
  }
  showEditModal.value = true
}

function deleteProject(project) {
  if (confirm(`Da li ste sigurni da želite da obrišete projekat "${project.name}"?`)) {
    const index = projects.value.findIndex(p => p.id === project.id)
    if (index > -1) {
      projects.value.splice(index, 1)
    }
  }
}

function saveProject() {
  if (showCreateModal.value) {
    // Create new project
    const newProject = {
      id: Date.now(),
      ...projectForm.value,
      progress: 0,
      created: new Date().toISOString().split('T')[0],
      updated: new Date().toISOString().split('T')[0],
      leader: users.value.find(u => u.id === projectForm.value.leaderId)?.name || '',
      team: [users.value.find(u => u.id === projectForm.value.leaderId)].filter(Boolean)
    }
    projects.value.push(newProject)
  } else if (showEditModal.value && selectedProject.value) {
    // Update existing project
    const index = projects.value.findIndex(p => p.id === selectedProject.value.id)
    if (index > -1) {
      projects.value[index] = {
        ...projects.value[index],
        ...projectForm.value,
        leader: users.value.find(u => u.id === projectForm.value.leaderId)?.name || '',
        updated: new Date().toISOString().split('T')[0]
      }
    }
  }
  
  closeModals()
}

function closeModals() {
  showCreateModal.value = false
  showEditModal.value = false
  selectedProject.value = null
  projectForm.value = {
    name: '',
    description: '',
    leaderId: '',
    status: 'active',
    deadline: ''
  }
}

// Lifecycle
onMounted(() => {
  // Load projects data
  console.log('Projects loaded')
})
</script>

<style scoped>
.projects {
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
  margin-bottom: 30px;
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

/* Projects Grid */
.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.project-card {
  background: white;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
}

.project-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.project-header h3 {
  color: #2c3e50;
  font-size: 18px;
  margin: 0;
  flex: 1;
}

.project-actions {
  display: flex;
  gap: 8px;
}

.btn-icon-small {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  padding: 5px;
  border-radius: 4px;
  transition: background 0.3s;
}

.btn-icon-small:hover {
  background: #f8f9fa;
}

.btn-icon-small.danger:hover {
  background: #f8d7da;
}

.project-description {
  color: #7f8c8d;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 20px;
}

.project-progress {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.progress-container {
  flex: 1;
  height: 8px;
  background: #ecf0f1;
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #3498db, #2ecc71);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #7f8c8d;
  font-weight: 600;
  min-width: 35px;
}

.project-info {
  margin-bottom: 20px;
}

.project-status {
  margin-bottom: 15px;
}

.project-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.project-meta > div {
  display: flex;
  align-items: center;
  font-size: 13px;
}

.label {
  font-weight: 600;
  color: #2c3e50;
  min-width: 60px;
  margin-right: 8px;
}

.project-team {
  border-top: 1px solid #ecf0f1;
  padding-top: 15px;
}

.team-avatars {
  display: flex;
  gap: 8px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #3498db;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.avatar.more {
  background: #95a5a6;
  font-size: 10px;
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

/* Modal */
.modal {
  max-width: 600px;
  width: 90%;
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
  .projects {
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
  
  .projects-grid {
    grid-template-columns: 1fr;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column-reverse;
  }
}
</style>
