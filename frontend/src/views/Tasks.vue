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
          <button class="btn btn-primary" @click="showCreateModal = true">
            <span class="btn-icon">➕</span>
            Novi zadatak
          </button>
        </div>
      </div>
      
      <!-- Filters -->
      <div class="filters card">
        <div class="filter-group">
          <label>Projekat:</label>
          <select v-model="filters.project">
            <option value="">Svi projekti</option>
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
          <select v-model="filters.priority">
            <option value="">Svi prioriteti</option>
            <option value="high">Visok</option>
            <option value="medium">Srednji</option>
            <option value="low">Nizak</option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>Dodeljeno:</label>
          <select v-model="filters.assignee">
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
          >
        </div>
      </div>
      
      <!-- Kanban Board -->
      <div v-if="viewMode === 'kanban'" class="kanban-board">
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
                  class="btn-icon-small" 
                  @click.stop="editTask(task)"
                  title="Uredi"
                >
                  ✏️
                </button>
              </div>
              
              <div class="task-description">
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
                  {{ getProjectName(task.projectId) }}
                </span>
              </div>
              
              <div class="task-footer">
                <div class="task-assignee">
                  <div class="avatar">
                    {{ getAssigneeName(task.assigneeId).charAt(0) }}
                  </div>
                  <span>{{ getAssigneeName(task.assigneeId) }}</span>
                </div>
                
                <div class="task-deadline">
                  {{ formatDate(task.deadline) }}
                </div>
              </div>
            </div>
            
            <div v-if="getColumnTasks(column.id).length === 0" class="empty-column">
              <p>Nema zadataka</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- List View -->
      <div v-else class="tasks-list">
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
                    <div class="task-desc">{{ task.description }}</div>
                  </div>
                </td>
                <td>{{ getProjectName(task.projectId) }}</td>
                <td>
                  <span 
                    class="priority-badge" 
                    :class="`priority-${task.priority}`"
                  >
                    {{ getPriorityText(task.priority) }}
                  </span>
                </td>
                <td>
                  <span 
                    class="status-badge" 
                    :class="`status-${task.status}`"
                  >
                    {{ getStatusText(task.status) }}
                  </span>
                </td>
                <td>{{ getAssigneeName(task.assigneeId) }}</td>
                <td>{{ formatDate(task.deadline) }}</td>
                <td>
                  <div class="action-buttons">
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
              <label>Naziv zadatka</label>
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
                <label>Projekat</label>
                <select v-model="taskForm.projectId" required>
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
              <div class="form-group">
                <label>Dodeli korisniku</label>
                <select v-model="taskForm.assigneeId" required>
                  <option value="">Izaberite korisnika</option>
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
                <select v-model="taskForm.status">
                  <option value="todo">Za rad</option>
                  <option value="in-progress">U toku</option>
                  <option value="review">Na proveri</option>
                  <option value="done">Završeno</option>
                </select>
              </div>
            </div>
            
            <div class="form-group">
              <label>Deadline</label>
              <input 
                type="date" 
                v-model="taskForm.deadline"
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
const viewMode = ref('kanban')
const showCreateModal = ref(false)
const showEditModal = ref(false)
const selectedTask = ref(null)
const draggedTask = ref(null)
const listSort = ref('title')

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
  assigneeId: '',
  status: 'todo',
  deadline: ''
})

const kanbanColumns = ref([
  { id: 'todo', title: 'Za rad', color: '#95a5a6' },
  { id: 'in-progress', title: 'U toku', color: '#3498db' },
  { id: 'review', title: 'Na proveri', color: '#f39c12' },
  { id: 'done', title: 'Završeno', color: '#2ecc71' }
])

// Mock data
const tasks = ref([
  {
    id: 1,
    title: 'Implementacija login sistema',
    description: 'Kreiranje sigurnog login sistema sa 2FA',
    projectId: 1,
    priority: 'high',
    status: 'in-progress',
    assigneeId: 1,
    deadline: '2024-09-15',
    created: '2024-09-01'
  },
  {
    id: 2,
    title: 'Dizajn korisničkog interfejsa',
    description: 'Kreiranje modernog UI/UX dizajna',
    projectId: 1,
    priority: 'medium',
    status: 'todo',
    assigneeId: 2,
    deadline: '2024-09-20',
    created: '2024-09-02'
  },
  {
    id: 3,
    title: 'Testiranje API funkcionalnosti',
    description: 'Unit i integration testovi za API',
    projectId: 2,
    priority: 'high',
    status: 'review',
    assigneeId: 3,
    deadline: '2024-09-12',
    created: '2024-08-25'
  },
  {
    id: 4,
    title: 'Optimizacija baze podataka',
    description: 'Poboljšanje performansi database upita',
    projectId: 2,
    priority: 'medium',
    status: 'done',
    assigneeId: 1,
    deadline: '2024-09-08',
    created: '2024-08-20'
  }
])

const projects = ref([
  { id: 1, name: 'Web Portal Refactoring' },
  { id: 2, name: 'AI Analitika Modula' },
  { id: 3, name: 'Mobile Aplikacija' }
])

const users = ref([
  { id: 1, name: 'Marko Petrović' },
  { id: 2, name: 'Ana Jovanović' },
  { id: 3, name: 'Stefan Nikolić' },
  { id: 4, name: 'Milica Stojković' }
])

// Computed
const filteredTasks = computed(() => {
  let filtered = tasks.value

  // Apply filters
  if (filters.value.project) {
    filtered = filtered.filter(t => t.projectId == filters.value.project)
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
      t.description.toLowerCase().includes(search)
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
          return new Date(a.deadline) - new Date(b.deadline)
        case 'created':
          return new Date(b.created) - new Date(a.created)
        default:
          return 0
      }
    })
  }

  return filtered
})

// Methods
function toggleView() {
  viewMode.value = viewMode.value === 'kanban' ? 'list' : 'kanban'
}

function getColumnTasks(columnId) {
  return filteredTasks.value.filter(task => task.status === columnId)
}

function getProjectName(projectId) {
  return projects.value.find(p => p.id === projectId)?.name || 'Nepoznat projekat'
}

function getAssigneeName(assigneeId) {
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

function getStatusText(status) {
  const statusMap = {
    'todo': 'Za rad',
    'in-progress': 'U toku',
    'review': 'Na proveri',
    'done': 'Završeno'
  }
  return statusMap[status] || status
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString('sr-RS')
}

function selectTask(task) {
  selectedTask.value = task
  console.log('Selected task:', task)
}

function editTask(task) {
  selectedTask.value = task
  taskForm.value = {
    title: task.title,
    description: task.description,
    projectId: task.projectId,
    priority: task.priority,
    assigneeId: task.assigneeId,
    status: task.status,
    deadline: task.deadline
  }
  showEditModal.value = true
}

function deleteTask(task) {
  if (confirm(`Da li ste sigurni da želite da obrišete zadatak "${task.title}"?`)) {
    const index = tasks.value.findIndex(t => t.id === task.id)
    if (index > -1) {
      tasks.value.splice(index, 1)
    }
  }
}

function saveTask() {
  if (showCreateModal.value) {
    // Create new task
    const newTask = {
      id: Date.now(),
      ...taskForm.value,
      created: new Date().toISOString().split('T')[0]
    }
    tasks.value.push(newTask)
  } else if (showEditModal.value && selectedTask.value) {
    // Update existing task
    const index = tasks.value.findIndex(t => t.id === selectedTask.value.id)
    if (index > -1) {
      tasks.value[index] = {
        ...tasks.value[index],
        ...taskForm.value
      }
    }
  }
  
  closeModals()
}

function closeModals() {
  showCreateModal.value = false
  showEditModal.value = false
  selectedTask.value = null
  taskForm.value = {
    title: '',
    description: '',
    projectId: '',
    priority: 'medium',
    assigneeId: '',
    status: 'todo',
    deadline: ''
  }
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

function handleDrop(event, newStatus) {
  event.preventDefault()
  
  if (draggedTask.value && draggedTask.value.status !== newStatus) {
    const taskIndex = tasks.value.findIndex(t => t.id === draggedTask.value.id)
    if (taskIndex > -1) {
      tasks.value[taskIndex].status = newStatus
    }
  }
  
  draggedTask.value = null
}

// Lifecycle
onMounted(() => {
  console.log('Tasks loaded')
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
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  overflow-x: auto;
  min-height: 500px;
}

.kanban-column {
  background: #f8f9fa;
  border-radius: 8px;
  min-width: 300px;
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
}

.empty-column {
  text-align: center;
  padding: 40px 20px;
  color: #95a5a6;
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

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 30px;
  border-top: 1px solid #ecf0f1;
  padding-top: 20px;
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
</style>
