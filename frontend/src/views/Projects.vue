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
  leaderId: '',
  status: 'active',
  deadline: ''
})

import { fetchUserProjects, createProject } from '../services/projectService.js'
const projects = ref([])

const router = useRouter()

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

function openDocumentation(project) {
  router.push(`/projects/${project.id}/documents`)
}

function saveProject() {
  if (showCreateModal.value) {
    // Create new project via backend
    createProject(projectForm.value.name, projectForm.value.description, projectForm.value.deadline)
      .then(async () => {
        const data = await fetchUserProjects()
        projects.value = data
        closeModals()
      })
      .catch(err => alert('Greška pri kreiranju projekta: ' + err))
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
  
  if (!showCreateModal.value) closeModals()
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
onMounted(async () => {
  try {
    const data = await fetchUserProjects()
    projects.value = data
  } catch (e) {
    console.error('Neuspešno učitavanje projekata:', e)
  }
})
</script>
