<template>
  <Layout>
    <div class="tasks">
      <!-- Header -->
      <div class="page-header">
        <div>
          <h2>Zadaci</h2>
          <div class="breadcrumb">Početna > Zadaci</div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" @click="toggleView">
            {{ viewMode === 'kanban' ? '📋 Lista' : '📊 Kanban' }}
          </button>
          <button 
            v-if="authStore.isManager"
            class="btn btn-info" 
            @click="goToRequests"
          >
            <span class="btn-icon">📋</span>
            Lista zahteva
          </button>
          <button 
            v-if="canManageProject"
            class="btn btn-primary" 
            @click="openCreateModal" 
            :disabled="!filters.project"
          >
            <span class="btn-icon">➕</span>
            Novi zadatak
          </button>
        </div>
      </div>
      
      <!-- Filters -->
            <!-- Filters -->
      <div class="filters card">
        <div class="filter-group">
          <label>Projekat:</label>
          <select v-model="filters.project" :disabled="loading">
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
        
        <div class="filter-group">
          <label>Prioritet:</label>
          <select v-model="filters.priority" :disabled="!filters.project">
            <option value="">Svi prioriteti</option>
            <option value="high">Visok</option>
            <option value="medium">Srednji</option>
            <option value="low">Nizak</option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>Dodeljeno:</label>
          <select v-model="filters.assignee" :disabled="!filters.project">
            <option value="">Svi korisnici</option>
            <option 
              v-for="user in users" 
              :key="user.id"
              :value="user.id"
            >
              {{ user.name }}
            </option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>Pretraga:</label>
          <input 
            type="text" 
            v-model="filters.search" 
            placeholder="Pretraži zadatke..."
            class="search-input"
            :disabled="!filters.project"
          >
        </div>
      </div>
      
      <!-- Kanban Board -->
      <div v-if="viewMode === 'kanban'" class="kanban-board">
        <div v-if="loading" class="loading-state">
          <p>Učitavanje...</p>
        </div>
        
        <div v-else-if="!filters.project" class="empty-state">
          <p>👆 Izaberite projekat da biste videli zadatke</p>
        </div>
        
        <div v-else-if="kanbanColumns.length === 0" class="empty-state">
          <p>📋 Ovaj projekat nema definisan radni tok</p>
        </div>
        
        <template v-else>
          <div 
            v-for="column in kanbanColumns" 
            :key="column.id"
            class="kanban-column"
          >
            <div class="column-header">
              <h3>{{ column.title }}</h3>
              <span class="task-count">{{ getColumnTasks(column.id).length }}</span>
            </div>
            
            <div 
              class="column-content"
              @drop="handleDrop($event, column.id)"
              @dragover="handleDragOver"
            >
              <div 
                v-for="task in getColumnTasks(column.id)" 
                :key="task.id"
                class="task-card"
                :draggable="true"
                @dragstart="handleDragStart($event, task)"
                @click="selectTask(task)"
              >
                <div class="task-header">
                  <h4>{{ task.title }}</h4>
                  <button 
                    v-if="canManageProject"
                    class="btn-icon-small" 
                    @click.stop="editTask(task)"
                    title="Uredi"
                  >
                    ✏️
                  </button>
                </div>
                
                <div class="task-description" v-if="task.description">
                  {{ task.description }}
                </div>
                
                <div class="task-meta">
                  <span 
                    class="priority-badge" 
                    :class="`priority-${task.priority}`"
                  >
                    {{ getPriorityText(task.priority) }}
                  </span>
                  
                  <span class="task-project">
                    {{ task.projectName || getProjectName(task.projectId) }}
                  </span>
                </div>
                
                <div class="task-footer">
                  <div class="task-assignee">
                    <div class="avatar" v-if="task.assigneeName || task.assigneeId">
                      {{ (task.assigneeName || getAssigneeName(task.assigneeId)).charAt(0) }}
                    </div>
                    <span>{{ task.assigneeName || getAssigneeName(task.assigneeId) }}</span>
                  </div>
                  
                  <div class="task-actions" v-if="canManageProject">
                    <button 
                      class="btn-icon-small danger" 
                      @click.stop="deleteTask(task)"
                      title="Obriši"
                    >
                      🗑️
                    </button>
                  </div>
                </div>
                
                <div class="task-deadline" v-if="task.deadline">
                  📅 {{ formatDate(task.deadline) }}
                </div>
              </div>
              
              <div v-if="getColumnTasks(column.id).length === 0" class="empty-column">
                <p>Nema zadataka</p>
              </div>
              
              <!-- Add Task Button in Column -->
              <button 
                v-if="canManageProject"
                class="btn-add-task-column" 
                @click="openCreateModalForPhase(column.id)"
                title="Dodaj zadatak u ovu fazu"
              >
                Dodaj Zadatak
              </button>
            </div>
          </div>
        </template>
      </div>
      
      <!-- List View -->
      <div v-else class="tasks-list">
        <div v-if="loading" class="loading-state">
          <p>Učitavanje...</p>
        </div>
        
        <div v-else-if="!filters.project" class="empty-state">
          <p>👆 Izaberite projekat da biste videli zadatke</p>
        </div>
        
        <template v-else>
          <div class="list-header">
            <div class="sort-options">
              <label>Sortiranje:</label>
              <select v-model="listSort">
                <option value="title">Naziv</option>
                <option value="priority">Prioritet</option>
                <option value="deadline">Deadline</option>
                <option value="created">Datum kreiranja</option>
              </select>
            </div>
          </div>
          
          <div class="tasks-table">
            <table class="table">
              <thead>
                <tr>
                  <th>Zadatak</th>
                  <th>Projekat</th>
                  <th>Prioritet</th>
                  <th>Status</th>
                  <th>Dodeljeno</th>
                  <th>Deadline</th>
                  <th>Akcije</th>
                </tr>
              </thead>
              <tbody>
                <tr 
                  v-for="task in filteredTasks" 
                  :key="task.id"
                  @click="selectTask(task)"
                  class="task-row"
                >
                  <td>
                    <div class="task-cell">
                      <strong>{{ task.title }}</strong>
                      <div class="task-desc" v-if="task.description">{{ task.description }}</div>
                    </div>
                  </td>
                  <td>{{ task.projectName || getProjectName(task.projectId) }}</td>
                  <td>
                    <span 
                      class="priority-badge" 
                      :class="`priority-${task.priority}`"
                    >
                      {{ getPriorityText(task.priority) }}
                    </span>
                  </td>
                  <td>
                    <span class="status-badge">
                      {{ task.phaseName || getStatusText(task.phaseId) }}
                    </span>
                  </td>
                  <td>{{ task.assigneeName || getAssigneeName(task.assigneeId) }}</td>
                  <td>{{ formatDate(task.deadline) }}</td>
                  <td>
                    <div class="action-buttons" v-if="canManageProject">
                      <button 
                        class="btn-icon-small" 
                        @click.stop="editTask(task)"
                        title="Uredi"
                      >
                        ✏️
                      </button>
                      <button 
                        class="btn-icon-small danger" 
                        @click.stop="deleteTask(task)"
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
        </template>
      </div>
      
      <!-- View Task Details Modal -->
      <div v-if="showDetailsModal" class="modal-overlay" @click="closeModals">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h3 class="modal-title">Detalji zadatka</h3>
            <button class="modal-close" @click="closeModals">×</button>
          </div>
          
          <div class="form-group">
            <label>Naziv zadatka</label>
            <input 
              type="text" 
              :value="selectedTask?.title || ''" 
              readonly
              disabled
            >
          </div>
          
          <div class="form-group">
            <label>Opis</label>
            <textarea 
              :value="selectedTask?.description || ''" 
              rows="3"
              readonly
              disabled
            ></textarea>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>Projekat</label>
              <input 
                type="text" 
                :value="selectedTask?.projectName || getProjectName(selectedTask?.projectId)" 
                readonly
                disabled
              >
            </div>
            
            <div class="form-group">
              <label>Prioritet</label>
              <input 
                type="text" 
                :value="getPriorityText(selectedTask?.priority)" 
                readonly
                disabled
              >
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>Faza</label>
              <input 
                type="text" 
                :value="selectedTask?.phaseName || getStatusText(selectedTask?.phaseId)" 
                readonly
                disabled
              >
            </div>
            
            <div class="form-group">
              <label>Dodeljeno korisniku</label>
              <input 
                type="text" 
                :value="selectedTask?.assigneeName || getAssigneeName(selectedTask?.assigneeId)" 
                readonly
                disabled
              >
            </div>
          </div>
          
          <!-- Conditions Section -->
          <div v-if="phaseConditions.length > 0" class="form-group">
            <label>Uslovi za ovu fazu</label>
            <div class="conditions-list">
              <div 
                v-for="condition in phaseConditions" 
                :key="condition.id"
                class="condition-item"
              >
                <div class="condition-checkbox">
                  <input 
                    type="checkbox" 
                    :id="`condition-${condition.id}`"
                    :checked="condition.fulfilled"
                    @change="toggleConditionFulfillment(condition)"
                    class="condition-check"
                  >
                  <label :for="`condition-${condition.id}`" class="condition-label">
                    <span class="condition-description">{{ condition.description }}</span>
                    <span class="condition-criteria">{{ condition.criteria }}</span>
                  </label>
                </div>
              </div>
            </div>
          </div>
          
          <div v-else-if="selectedTask?.phaseId" class="form-group">
            <label>Uslovi za ovu fazu</label>
            <p class="no-conditions">Nema definisanih uslova za ovu fazu</p>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>Deadline</label>
              <input 
                type="date" 
                :value="selectedTask?.deadline" 
                readonly
                disabled
              >
            </div>
            
            <div class="form-group">
              <label>Progres</label>
              <input 
                type="text" 
                :value="`${selectedTask?.progress || 0}%`" 
                readonly
                disabled
              >
            </div>
          </div>
          
          <div class="form-group">
            <label>Resursi</label>
            <textarea 
              :value="selectedTask?.resources || ''" 
              rows="2"
              readonly
              disabled
            ></textarea>
          </div>
          
          <div class="form-group">
            <label>Datum kreiranja</label>
            <input 
              type="text" 
              :value="formatDate(selectedTask?.created)" 
              readonly
              disabled
            >
          </div>
          
          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="closeModals">
              Zatvori
            </button>
            <button 
              v-if="!canManageProject && areAllConditionsFulfilled" 
              type="button" 
              class="btn btn-success" 
              @click="openPhaseChangeRequestModal"
              :disabled="loadingConditions"
            >
              🔄 Zahtev za promenu faze
            </button>
            <button 
              v-if="canManageProject" 
              type="button" 
              class="btn btn-primary" 
              @click="editTaskFromDetails"
            >
              Uredi zadatak
            </button>
          </div>
        </div>
      </div>
      
      <!-- Phase Change Request Modal -->
      <div v-if="showPhaseChangeModal" class="modal-overlay" @click="closeModals">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h3 class="modal-title">Zahtev za promenu faze</h3>
            <button class="modal-close" @click="closeModals">×</button>
          </div>
          
          <div class="form-group">
            <label>Zadatak</label>
            <input 
              type="text" 
              :value="selectedTask?.title || ''" 
              readonly
              disabled
            >
          </div>
          
          <div class="form-group">
            <label>Trenutna faza</label>
            <input 
              type="text" 
              :value="selectedTask?.phaseName || getStatusText(selectedTask?.phaseId)" 
              readonly
              disabled
            >
          </div>
          
          <div class="form-group">
            <label>Zahtevana faza (sledeća u redosledu)</label>
            <input 
              type="text" 
              :value="nextPhase?.naziv_faze || nextPhase?.name || 'Učitavanje...'" 
              readonly
              disabled
            >
          </div>
          
          <div class="form-group">
            <label>Opis zahteva</label>
            <textarea 
              v-model="phaseChangeForm.description" 
              rows="4"
              placeholder="Opišite razlog za promenu faze..."
              required
            ></textarea>
          </div>
          
          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="closeModals">
              Otkaži
            </button>
            <button 
              type="button" 
              class="btn btn-primary" 
              @click="submitPhaseChangeRequest"
              :disabled="!phaseChangeForm.description || loadingPhaseChange"
            >
              {{ loadingPhaseChange ? 'Slanje...' : 'Pošalji zahtev' }}
            </button>
          </div>
        </div>
      </div>
      
      <!-- Create/Edit Task Modal -->
      <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click="closeModals">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h3 class="modal-title">
              {{ showCreateModal ? 'Kreiranje novog zadatka' : 'Uređivanje zadatka' }}
            </h3>
            <button class="modal-close" @click="closeModals">×</button>
          </div>
          
          <form @submit.prevent="saveTask">
            <div class="form-group">
              <label>Naziv zadatka *</label>
              <input 
                type="text" 
                v-model="taskForm.title" 
                placeholder="Unesite naziv zadatka" 
                required
              >
            </div>
            
            <div class="form-group">
              <label>Opis</label>
              <textarea 
                v-model="taskForm.description" 
                placeholder="Unesite opis zadatka"
                rows="3"
              ></textarea>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label>Projekat *</label>
                <select v-model="taskForm.projectId" required :disabled="showEditModal">
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
                <label>Prioritet</label>
                <select v-model="taskForm.priority">
                  <option value="low">Nizak</option>
                  <option value="medium">Srednji</option>
                  <option value="high">Visok</option>
                </select>
              </div>
            </div>
            
            <div class="form-row">
              <div class="form-group form-group-full">
                <label>Faza {{ showEditModal ? '*' : '' }}</label>
                <select v-model="taskForm.phaseId" :required="showEditModal">
                  <option value="">{{ phases.length > 0 ? 'Izaberite fazu (opciono)' : 'Nema dostupnih faza' }}</option>
                  <option 
                    v-for="phase in phases" 
                    :key="phase.id" 
                    :value="phase.id"
                  >
                    {{ phase.name }}
                  </option>
                </select>
                <small v-if="showCreateModal" style="color: #7f8c8d; font-size: 12px; margin-top: 4px; display: block;">
                  Ako ne izaberete fazu, zadatak će biti dodeljen prvoj fazi projekta
                </small>
              </div>
            </div>
            
            <!-- Conditions Section in Edit Modal -->
            <div v-if="showEditModal && phaseConditions.length > 0" class="form-group">
              <label>Uslovi za ovu fazu</label>
              <div class="conditions-list">
                <div 
                  v-for="condition in phaseConditions" 
                  :key="condition.id"
                  class="condition-item"
                >
                  <div class="condition-checkbox">
                    <input 
                      type="checkbox" 
                      :id="`condition-edit-${condition.id}`"
                      :checked="condition.fulfilled"
                      @change="toggleConditionFulfillment(condition)"
                      class="condition-check"
                    >
                    <label :for="`condition-edit-${condition.id}`" class="condition-label">
                      <span class="condition-description">{{ condition.description }}</span>
                      <span class="condition-criteria">{{ condition.criteria }}</span>
                    </label>
                  </div>
                </div>
              </div>
            </div>
            
            <div v-else-if="showEditModal && taskForm.phaseId && phaseConditions.length === 0" class="form-group">
              <label>Uslovi za ovu fazu</label>
              <p class="no-conditions">Nema definisanih uslova za ovu fazu</p>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label>Dodeli korisniku</label>
                <select v-model="taskForm.assigneeId">
                  <option :value="null">Nedodeljeno</option>
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
                <label>Deadline</label>
                <input 
                  type="date" 
                  v-model="taskForm.deadline"
                >
              </div>
            </div>
            
            <div class="form-group">
              <label>Resursi</label>
              <input 
                type="text" 
                v-model="taskForm.resources" 
                placeholder="Unesite potrebne resurse (npr. 2 programera, server, oprema)"
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/auth.js'
import {
  fetchUserProjects,
  fetchWorkflowPhases,
  fetchProjectTasks,
  fetchProjectMembers,
  createTask,
  updateTask,
  deleteTask as deleteTaskService,
  moveTaskToPhase,
  fetchConditionsByPhase,
  fetchConditionAssessmentsByTask,
  updateConditionAssessment
} from '../services/taskService.js'

// Router
const router = useRouter()

// Auth store
const authStore = useAuthStore()

// Reactive data
const viewMode = ref('kanban')
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDetailsModal = ref(false)
const selectedTask = ref(null)
const draggedTask = ref(null)
const listSort = ref('title')
const loading = ref(false)
const selectedProjectId = ref(null)

const filters = ref({
  project: '',
  priority: '',
  assignee: '',
  search: ''
})

const taskForm = ref({
  title: '',
  description: '',
  projectId: '',
  priority: 'medium',
  assigneeId: null,
  deadline: '',
  phaseId: null,
  resources: ''
})

// Data from backend
const projects = ref([])
const phases = ref([])
const tasks = ref([])
const users = ref([])
const phaseConditions = ref([]) // Conditions for the selected phase
const taskConditionAssessments = ref([]) // Assessments for the selected task
const showPhaseChangeModal = ref(false)
const nextPhase = ref(null)
const loadingConditions = ref(false)
const loadingPhaseChange = ref(false)

const phaseChangeForm = ref({
  description: ''
})

// Check if all conditions are fulfilled
const areAllConditionsFulfilled = computed(() => {
  if (phaseConditions.value.length === 0) {
    return true // No conditions means they're all fulfilled
  }
  return phaseConditions.value.every(condition => condition.fulfilled === true)
})

// Check if user can manage the selected project (must be the project manager)
const canManageProject = computed(() => {
  if (!authStore.user || !selectedProjectId.value) return false
  
  // Find the selected project
  const project = projects.value.find(p => p.id === selectedProjectId.value)
  if (!project) return false
  
  // Check if logged-in user ID matches project's manager ID
  const currentUserId = authStore.user.korisnik_id || authStore.user.korisnikID
  const projectLeaderId = project.leaderId || project.rukovodilac_id
  
  console.log('🔐 Permission check:', {
    currentUserId,
    projectLeaderId,
    projectName: project.name,
    canManage: currentUserId === projectLeaderId
  })
  
  return currentUserId === projectLeaderId
})

// Computed kanban columns based on phases
const kanbanColumns = computed(() => {
  if (!phases.value || phases.value.length === 0) {
    return [
      { id: 1, title: 'Za rad', color: '#95a5a6' },
      { id: 2, title: 'U toku', color: '#3498db' },
      { id: 3, title: 'Na proveri', color: '#f39c12' },
      { id: 4, title: 'Završeno', color: '#2ecc71' }
    ]
  }
  
  const colors = ['#95a5a6', '#3498db', '#f39c12', '#2ecc71', '#9b59b6', '#e74c3c']
  return phases.value.map((phase, index) => ({
    id: phase.id,
    title: phase.name,
    color: colors[index % colors.length]
  }))
})

// Computed
const filteredTasks = computed(() => {
  let filtered = tasks.value

  // Apply filters
  if (filters.value.project) {
    filtered = filtered.filter(t => t.projectId == filters.value.project)
  }
  
  // If user is a researcher, only show tasks assigned to them
  if (authStore.isResearcher && authStore.user) {
    const currentUserId = authStore.user.korisnik_id || authStore.user.korisnikID
    console.log('🔍 Researcher filter active:', {
      userId: currentUserId,
      naziv_uloge: authStore.user.naziv_uloge,
      totalTasks: filtered.length
    })
    filtered = filtered.filter(t => {
      const isAssigned = t.assigneeId == currentUserId
      if (!isAssigned) {
        console.log('  ❌ Filtered out task:', t.title, 'assigned to:', t.assigneeId)
      }
      return isAssigned
    })
    console.log('  ✅ Tasks after researcher filter:', filtered.length)
  }
  
  if (filters.value.priority) {
    filtered = filtered.filter(t => t.priority === filters.value.priority)
  }
  
  if (filters.value.assignee) {
    filtered = filtered.filter(t => t.assigneeId == filters.value.assignee)
  }
  
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(t => 
      t.title.toLowerCase().includes(search) || 
      (t.description && t.description.toLowerCase().includes(search))
    )
  }

  // Sort for list view
  if (viewMode.value === 'list') {
    filtered.sort((a, b) => {
      switch (listSort.value) {
        case 'title':
          return a.title.localeCompare(b.title)
        case 'priority':
          const priorityOrder = { high: 3, medium: 2, low: 1 }
          return priorityOrder[b.priority] - priorityOrder[a.priority]
        case 'deadline':
          if (!a.deadline) return 1
          if (!b.deadline) return -1
          return new Date(a.deadline) - new Date(b.deadline)
        case 'created':
          if (!a.created) return 1
          if (!b.created) return -1
          return new Date(b.created) - new Date(a.created)
        default:
          return 0
      }
    })
  }

  return filtered
})

// Watch for project selection changes
watch(() => filters.value.project, async (newProjectId) => {
  if (newProjectId) {
    selectedProjectId.value = newProjectId
    await loadProjectData(newProjectId)
  } else {
    selectedProjectId.value = null
    phases.value = []
    tasks.value = []
    users.value = []
  }
})

// Watch for phase changes in task form
watch(() => taskForm.value.phaseId, async (newPhaseId) => {
  if (newPhaseId && (showEditModal.value || showDetailsModal.value)) {
    await loadPhaseConditions(newPhaseId)
  }
})

// Methods
async function loadProjects() {
  loading.value = true
  try {
    projects.value = await fetchUserProjects()
    
    // Auto-select first project if available
    if (projects.value.length > 0 && !filters.value.project) {
      filters.value.project = projects.value[0].id
    }
  } catch (error) {
    console.error('Error loading projects:', error)
  } finally {
    loading.value = false
  }
}

async function loadProjectData(projectId) {
  loading.value = true
  try {
    // Find project to get workflow ID
    const project = projects.value.find(p => p.id === projectId)
    
    if (project && project.workflowId) {
      // Load phases for the workflow
      phases.value = await fetchWorkflowPhases(project.workflowId)
    } else {
      phases.value = []
    }
    
    // Load tasks for the project
    tasks.value = await fetchProjectTasks(projectId)
    
    // Load project members
    users.value = await fetchProjectMembers(projectId)
  } catch (error) {
    console.error('Error loading project data:', error)
  } finally {
    loading.value = false
  }
}

function toggleView() {
  viewMode.value = viewMode.value === 'kanban' ? 'list' : 'kanban'
}

function goToRequests() {
  router.push('/phase-change-requests')
}

function getColumnTasks(columnId) {
  return filteredTasks.value.filter(task => task.phaseId === columnId)
}

function getProjectName(projectId) {
  return projects.value.find(p => p.id === projectId)?.name || 'Nepoznat projekat'
}

function getAssigneeName(assigneeId) {
  if (!assigneeId) return 'Nedodeljeno'
  return users.value.find(u => u.id === assigneeId)?.name || 'Nepoznat korisnik'
}

function getPriorityText(priority) {
  const priorityMap = {
    high: 'Visok',
    medium: 'Srednji',
    low: 'Nizak'
  }
  return priorityMap[priority] || priority
}

function getStatusText(phaseId) {
  const phase = phases.value.find(p => p.id === phaseId)
  return phase ? phase.name : 'Nepoznata faza'
}

function formatDate(dateString) {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('sr-RS')
}

function selectTask(task) {
  selectedTask.value = task
  showDetailsModal.value = true
  console.log('Selected task:', task)
  
  // Load conditions for the task's phase
  if (task.phaseId) {
    loadPhaseConditionsForTask(task.id, task.phaseId)
  }
}

async function loadPhaseConditions(phaseId) {
  try {
    phaseConditions.value = await fetchConditionsByPhase(phaseId)
    console.log('Loaded conditions for phase:', phaseId, phaseConditions.value)
  } catch (error) {
    console.error('Error loading phase conditions:', error)
    phaseConditions.value = []
  }
}

async function loadPhaseConditionsForTask(taskId, phaseId) {
  try {
    // Load conditions for the phase
    phaseConditions.value = await fetchConditionsByPhase(phaseId)
    
    // Load assessments for the task
    const assessments = await fetchConditionAssessmentsByTask(taskId)
    
    // Create a map of assessments by condition ID
    const assessmentMap = {}
    assessments.forEach(a => {
      assessmentMap[a.conditionId] = a
    })
    
    // Merge conditions with assessments
    phaseConditions.value = phaseConditions.value.map(condition => ({
      ...condition,
      fulfilled: assessmentMap[condition.id]?.fulfilled || false,
      assessmentId: assessmentMap[condition.id]?.id || null,
      note: assessmentMap[condition.id]?.note || ''
    }))
    
    console.log('Loaded conditions with assessments:', phaseConditions.value)
  } catch (error) {
    console.error('Error loading conditions for task:', error)
    phaseConditions.value = []
  }
}

async function toggleConditionFulfillment(condition) {
  if (!selectedTask.value) return
  
  const newFulfilledStatus = !condition.fulfilled
  
  loading.value = true
  try {
    await updateConditionAssessment(
      selectedTask.value.id,
      condition.id,
      newFulfilledStatus,
      condition.note || ''
    )
    
    // Update local state
    condition.fulfilled = newFulfilledStatus
    
    console.log(`Condition ${condition.id} updated to ${newFulfilledStatus}`)
  } catch (error) {
    console.error('Error updating condition assessment:', error)
    alert('Greška pri ažuriranju uslova')
  } finally {
    loading.value = false
  }
}

function editTaskFromDetails() {
  // Close details modal and open edit modal
  showDetailsModal.value = false
  editTask(selectedTask.value)
}

function editTask(task) {
  // Check permission before allowing edit
  if (!canManageProject.value) {
    alert('Nemate dozvolu za uređivanje zadataka. Samo rukovodilac projekta može vršiti izmene.')
    return
  }
  
  selectedTask.value = task
  taskForm.value = {
    title: task.title,
    description: task.description || '',
    projectId: task.projectId,
    priority: task.priority,
    assigneeId: task.assigneeId,
    deadline: task.deadline || '',
    phaseId: task.phaseId,
    resources: task.resources || ''
  }
  showEditModal.value = true
  
  // Load conditions for the task's phase
  if (task.phaseId) {
    loadPhaseConditionsForTask(task.id, task.phaseId)
  }
}

async function deleteTask(task) {
  // Check permission before allowing delete
  if (!canManageProject.value) {
    alert('Nemate dozvolu za brisanje zadataka. Samo rukovodilac projekta može vršiti izmene.')
    return
  }
  
  if (confirm(`Da li ste sigurni da želite da obrišete zadatak "${task.title}"?`)) {
    loading.value = true
    try {
      const result = await deleteTaskService(task.id)
      if (result.success) {
        // Reload tasks
        if (selectedProjectId.value) {
          await loadProjectData(selectedProjectId.value)
        }
      } else {
        alert('Greška pri brisanju zadatka: ' + (result.error || 'Nepoznata greška'))
      }
    } catch (error) {
      console.error('Error deleting task:', error)
      alert('Greška pri brisanju zadatka')
    } finally {
      loading.value = false
    }
  }
}

async function saveTask() {
  if (!taskForm.value.title) {
    alert('Naziv zadatka je obavezan')
    return
  }
  
  if (!taskForm.value.projectId) {
    alert('Projekat je obavezan')
    return
  }
  
  if (showEditModal.value && !taskForm.value.phaseId) {
    alert('Faza je obavezna prilikom uređivanja zadatka')
    return
  }
  
  loading.value = true
  try {
    if (showCreateModal.value) {
      // Create new task
      const result = await createTask(taskForm.value)
      if (result.success) {
        closeModals()
        // Reload tasks
        if (selectedProjectId.value) {
          await loadProjectData(selectedProjectId.value)
        }
      } else {
        alert('Greška pri kreiranju zadatka: ' + (result.error || 'Nepoznata greška'))
      }
    } else if (showEditModal.value && selectedTask.value) {
      // Update existing task
      const result = await updateTask(selectedTask.value.id, taskForm.value)
      if (result.success) {
        closeModals()
        // Reload tasks
        if (selectedProjectId.value) {
          await loadProjectData(selectedProjectId.value)
        }
      } else {
        alert('Greška pri ažuriranju zadatka: ' + (result.error || 'Nepoznata greška'))
      }
    }
  } catch (error) {
    console.error('Error saving task:', error)
    alert('Greška pri čuvanju zadatka')
  } finally {
    loading.value = false
  }
}

function closeModals() {
  showCreateModal.value = false
  showEditModal.value = false
  showDetailsModal.value = false
  showPhaseChangeModal.value = false
  selectedTask.value = null
  phaseConditions.value = []
  nextPhase.value = null
  phaseChangeForm.value = {
    description: ''
  }
  taskForm.value = {
    title: '',
    description: '',
    projectId: filters.value.project || '',
    priority: 'medium',
    assigneeId: null,
    deadline: '',
    phaseId: null,
    resources: ''
  }
}

async function openPhaseChangeRequestModal() {
  if (!selectedTask.value) return
  
  loadingConditions.value = true
  try {
    const { CheckTaskConditionsFulfilled, GetNextPhaseForTask, HasPendingPhaseChangeRequest } = window.go.main.App
    
    // Check if there's already a pending request
    const hasPending = await HasPendingPhaseChangeRequest(selectedTask.value.id)
    if (hasPending) {
      alert('Već postoji aktivni zahtev za promenu faze za ovaj zadatak. Molimo sačekajte da rukovodilac projekta pregleda postojeći zahtev.')
      return
    }
    
    // Check if all conditions are fulfilled
    const fulfilled = await CheckTaskConditionsFulfilled(selectedTask.value.id)
    if (!fulfilled) {
      alert('Svi uslovi za trenutnu fazu moraju biti ispunjeni pre nego što možete zatražiti promenu faze.')
      return
    }
    
    // Get next phase
    nextPhase.value = await GetNextPhaseForTask(selectedTask.value.id)
    
    if (!nextPhase.value) {
      alert('Nema sledeće faze u radnom toku.')
      return
    }
    
    // Open modal
    showDetailsModal.value = false
    showPhaseChangeModal.value = true
  } catch (error) {
    console.error('Error opening phase change modal:', error)
    alert('Greška: ' + (error.message || error))
  } finally {
    loadingConditions.value = false
  }
}

async function submitPhaseChangeRequest() {
  if (!selectedTask.value || !nextPhase.value || !phaseChangeForm.value.description) {
    return
  }
  
  loadingPhaseChange.value = true
  try {
    const { RequestPhaseChange } = window.go.main.App
    
    await RequestPhaseChange(
      selectedTask.value.id,
      nextPhase.value.faza_id || nextPhase.value.id,
      phaseChangeForm.value.description
    )
    
    alert('Zahtev za promenu faze je uspešno poslat! Rukovodilac projekta će ga pregledati.')
    closeModals()
    
    // Reload tasks
    if (selectedProjectId.value) {
      await loadProjectData(selectedProjectId.value)
    }
  } catch (error) {
    console.error('Error submitting phase change request:', error)
    alert('Greška pri slanju zahteva: ' + (error.message || error))
  } finally {
    loadingPhaseChange.value = false
  }
}

function openCreateModal() {
  // Check permission before allowing create
  if (!canManageProject.value) {
    alert('Nemate dozvolu za kreiranje zadataka. Samo rukovodilac projekta može kreirati zadatke.')
    return
  }
  
  if (!filters.value.project) {
    alert('Molimo prvo izaberite projekat')
    return
  }
  taskForm.value.projectId = filters.value.project
  showCreateModal.value = true
}

function openCreateModalForPhase(phaseId) {
  // Check permission before allowing create
  if (!canManageProject.value) {
    alert('Nemate dozvolu za kreiranje zadataka. Samo rukovodilac projekta može kreirati zadatke.')
    return
  }
  
  if (!filters.value.project) {
    alert('Molimo prvo izaberite projekat')
    return
  }
  taskForm.value.projectId = filters.value.project
  taskForm.value.phaseId = phaseId
  showCreateModal.value = true
}

// Drag and Drop handlers
function handleDragStart(event, task) {
  draggedTask.value = task
  event.dataTransfer.effectAllowed = 'move'
}

function handleDragOver(event) {
  event.preventDefault()
  event.dataTransfer.dropEffect = 'move'
}

async function handleDrop(event, newPhaseId) {
  event.preventDefault()
  
  if (draggedTask.value && draggedTask.value.phaseId !== newPhaseId) {
    loading.value = true
    try {
      const result = await moveTaskToPhase(draggedTask.value.id, newPhaseId)
      if (result.success) {
        // Update local task
        const taskIndex = tasks.value.findIndex(t => t.id === draggedTask.value.id)
        if (taskIndex > -1) {
          tasks.value[taskIndex].phaseId = newPhaseId
          // Update phase name
          const phase = phases.value.find(p => p.id === newPhaseId)
          if (phase) {
            tasks.value[taskIndex].phaseName = phase.name
          }
        }
      } else {
        alert('Greška pri premeštanju zadatka: ' + (result.error || 'Nepoznata greška'))
      }
    } catch (error) {
      console.error('Error moving task:', error)
      alert('Greška pri premeštanju zadatka')
    } finally {
      loading.value = false
    }
  }
  
  draggedTask.value = null
}

// Lifecycle
onMounted(async () => {
  await loadProjects()
})
</script>

<style scoped>
.tasks {
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

.header-actions {
  display: flex;
  gap: 15px;
}

.btn-icon {
  margin-right: 8px;
}

/* Loading and Empty States */
.loading-state,
.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #7f8c8d;
  font-size: 16px;
  grid-column: 1 / -1;
}

.loading-state p,
.empty-state p {
  margin: 0;
  font-size: 18px;
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

/* Kanban Board */
.kanban-board {
  display: flex;
  gap: 20px;
  overflow-x: auto;
  overflow-y: hidden;
  min-height: 500px;
  padding-bottom: 10px;
}

/* Custom scrollbar for better UX */
.kanban-board::-webkit-scrollbar {
  height: 8px;
}

.kanban-board::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.kanban-board::-webkit-scrollbar-thumb {
  background: #bdc3c7;
  border-radius: 4px;
}

.kanban-board::-webkit-scrollbar-thumb:hover {
  background: #95a5a6;
}

.kanban-column {
  background: #f8f9fa;
  border-radius: 8px;
  min-width: 300px;
  max-width: 300px;
  flex-shrink: 0;
}

.column-header {
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.column-header h3 {
  color: #2c3e50;
  margin: 0;
  font-size: 16px;
}

.task-count {
  background: #95a5a6;
  color: white;
  border-radius: 12px;
  padding: 2px 8px;
  font-size: 12px;
  font-weight: 600;
}

.column-content {
  padding: 20px;
  min-height: 400px;
}

.task-card {
  background: white;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
}

.task-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.task-card:active {
  transform: scale(0.95);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.task-header h4 {
  color: #2c3e50;
  margin: 0;
  font-size: 14px;
  flex: 1;
}

.task-description {
  color: #7f8c8d;
  font-size: 13px;
  line-height: 1.4;
  margin-bottom: 15px;
}

.task-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.task-project {
  font-size: 11px;
  color: #95a5a6;
  background: #ecf0f1;
  padding: 2px 6px;
  border-radius: 4px;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #7f8c8d;
}

.task-assignee {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.task-actions {
  display: flex;
  gap: 4px;
}

.avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #3498db;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 600;
}

.task-deadline {
  font-size: 11px;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #ecf0f1;
}

.empty-column {
  text-align: center;
  padding: 40px 20px;
  color: #95a5a6;
}

/* Add Task Button in Column */
.btn-add-task-column {
  width: 100%;
  padding: 12px;
  background: transparent;
  border: 2px dashed #bdc3c7;
  border-radius: 8px;
  color: #7f8c8d;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 10px;
}

.btn-add-task-column:hover {
  background: #f8f9fa;
  border-color: #3498db;
  color: #3498db;
}

.btn-add-task-column:active {
  transform: scale(0.98);
}

/* List View */
.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
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

.tasks-table {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.task-row {
  cursor: pointer;
  transition: background 0.3s;
}

.task-row:hover {
  background: #f8f9fa;
}

.task-cell {
  max-width: 200px;
}

.task-desc {
  font-size: 12px;
  color: #7f8c8d;
  margin-top: 4px;
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

/* Priority badges */
.priority-badge {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.priority-high {
  background: #f8d7da;
  color: #721c24;
}

.priority-medium {
  background: #fff3cd;
  color: #856404;
}

.priority-low {
  background: #d1ecf1;
  color: #0c5460;
}

/* Status badges */
.status-badge {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.status-todo {
  background: #f8f9fa;
  color: #495057;
}

.status-in-progress {
  background: #d1ecf1;
  color: #0c5460;
}

.status-review {
  background: #fff3cd;
  color: #856404;
}

.status-done {
  background: #d4edda;
  color: #155724;
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

.form-group-full {
  grid-column: 1 / -1;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 30px;
  border-top: 1px solid #ecf0f1;
  padding-top: 20px;
}

/* Conditions List */
.conditions-list {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 15px;
  margin-top: 8px;
}

.condition-item {
  padding: 12px 0;
  border-bottom: 1px solid #e9ecef;
}

.condition-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.condition-item:first-child {
  padding-top: 0;
}

.condition-checkbox {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.condition-check {
  width: 20px;
  height: 20px;
  margin-top: 2px;
  cursor: pointer;
  flex-shrink: 0;
  accent-color: #2ecc71;
}

.condition-label {
  flex: 1;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.condition-description {
  font-size: 14px;
  color: #2c3e50;
  font-weight: 500;
}

.condition-criteria {
  font-size: 12px;
  color: #7f8c8d;
  font-style: italic;
}

.no-conditions {
  color: #95a5a6;
  font-size: 14px;
  font-style: italic;
  margin: 10px 0 0 0;
  padding: 15px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  text-align: center;
}


/* Responsive */
@media (max-width: 1024px) {
  .kanban-board {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .tasks {
    padding: 20px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 20px;
    align-items: stretch;
  }
  
  .header-actions {
    justify-content: space-between;
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
  
  .kanban-board {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .kanban-column {
    min-width: auto;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column-reverse;
  }
}

@media (max-width: 480px) {
  .tasks-table {
    overflow-x: auto;
  }
  
  .table {
    min-width: 600px;
  }
}

/* Success button */
.btn-success {
  background: #2ecc71;
  color: white;
  border: none;
}

.btn-success:hover {
  background: #27ae60;
}

.btn-success:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}
</style>

