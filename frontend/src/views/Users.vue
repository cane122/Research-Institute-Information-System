<template>
  <Layout>
    <div class="users">
      <!-- Header -->
      <div class="page-header">
        <div>
          <h2>Korisnici</h2>
          <div class="breadcrumb">Početna > Korisnici</div>
        </div>
        <button class="btn btn-primary" @click="showCreateModal = true">
          <span class="btn-icon">➕</span>
          Novi korisnik
        </button>
      </div>
      
      <!-- Filters and Stats -->
      <div class="users-overview">
        <div class="stats-cards">
          <div class="stat-card total">
            <div class="stat-content">
              <h3>{{ totalUsers }}</h3>
              <p>Ukupno korisnika</p>
            </div>
            <div class="stat-icon">👥</div>
          </div>
          
          <div class="stat-card active">
            <div class="stat-content">
              <h3>{{ activeUsers }}</h3>
              <p>Aktivni korisnici</p>
            </div>
            <div class="stat-icon">✅</div>
          </div>
          
          <div class="stat-card admin">
            <div class="stat-content">
              <h3>{{ adminUsers }}</h3>
              <p>Administratori</p>
            </div>
            <div class="stat-icon">🔧</div>
          </div>
        </div>
        
        <div class="filters">
          <div class="filter-group">
            <label>Uloga:</label>
            <select v-model="filters.role">
              <option value="">Sve uloge</option>
              <option value="admin">Administrator</option>
              <option value="manager">Menadžer</option>
              <option value="researcher">Istraživač</option>
              <option value="user">Korisnik</option>
            </select>
          </div>
          
          <div class="filter-group">
            <label>Status:</label>
            <select v-model="filters.status">
              <option value="">Svi statusi</option>
              <option value="active">Aktivni</option>
              <option value="inactive">Neaktivni</option>
              <option value="pending">Na čekanju</option>
            </select>
          </div>
          
          <div class="filter-group">
            <label>Pretraga:</label>
            <input 
              type="text" 
              v-model="filters.search" 
              placeholder="Pretraži korisnike..."
              class="search-input"
            >
          </div>
        </div>
      </div>
      
      <!-- Users Table -->
      <div class="users-table card">
        <div class="table-header">
          <h3>Lista korisnika</h3>
          <div class="table-actions">
            <button class="btn btn-secondary" @click="exportUsers">
              📊 Export
            </button>
          </div>
        </div>
        
        <div class="table-wrapper">
          <table class="table">
            <thead>
              <tr>
                <th>
                  <input 
                    type="checkbox" 
                    v-model="selectAll"
                    @change="toggleSelectAll"
                  >
                </th>
                <th>Korisnik</th>
                <th>Email</th>
                <th>Uloga</th>
                <th>Status</th>
                <th>Poslednja prijava</th>
                <th>Kreiran</th>
                <th>Akcije</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="user in filteredUsers" 
                :key="user.id"
                class="user-row"
              >
                <td>
                  <input 
                    type="checkbox" 
                    v-model="selectedUsers"
                    :value="user.id"
                  >
                </td>
                <td>
                  <div class="user-info">
                    <div class="user-avatar">
                      {{ user.firstName.charAt(0) + user.lastName.charAt(0) }}
                    </div>
                    <div class="user-details">
                      <div class="user-name">
                        {{ user.firstName }} {{ user.lastName }}
                      </div>
                      <div class="user-username">@{{ user.username }}</div>
                    </div>
                  </div>
                </td>
                <td>{{ user.email }}</td>
                <td>
                  <span 
                    class="role-badge" 
                    :class="`role-${user.role}`"
                  >
                    {{ getRoleText(user.role) }}
                  </span>
                </td>
                <td>
                  <span 
                    class="status-badge" 
                    :class="`status-${user.status}`"
                  >
                    {{ getStatusText(user.status) }}
                  </span>
                </td>
                <td>{{ formatDateTime(user.lastLogin) }}</td>
                <td>{{ formatDate(user.created) }}</td>
                <td>
                  <div class="action-buttons">
                    <button 
                      class="btn-icon-small" 
                      @click="editUser(user)"
                      title="Uredi"
                    >
                      ✏️
                    </button>
                    <button 
                      class="btn-icon-small" 
                      @click="toggleUserStatus(user)"
                      :title="user.status === 'active' ? 'Deaktiviraj' : 'Aktiviraj'"
                    >
                      {{ user.status === 'active' ? '🔒' : '🔓' }}
                    </button>
                    <button 
                      class="btn-icon-small" 
                      @click="resetPassword(user)"
                      title="Resetuj lozinku"
                    >
                      🔑
                    </button>
                    <button 
                      class="btn-icon-small danger" 
                      @click="deleteUser(user)"
                      title="Obriši"
                    >
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          
          <!-- Empty State -->
          <div v-if="filteredUsers.length === 0" class="empty-state">
            <div class="empty-icon">👤</div>
            <h3>Nema korisnika</h3>
            <p>Nema korisnika koji odgovaraju filterima</p>
          </div>
        </div>
        
        <!-- Bulk Actions -->
        <div v-if="selectedUsers.length > 0" class="bulk-actions">
          <div class="selected-count">
            {{ selectedUsers.length }} korisnik(a) izabrano
          </div>
          <div class="bulk-buttons">
            <button class="btn btn-secondary" @click="bulkActivate">
              Aktiviraj
            </button>
            <button class="btn btn-secondary" @click="bulkDeactivate">
              Deaktiviraj
            </button>
            <button class="btn btn-danger" @click="bulkDelete">
              Obriši
            </button>
          </div>
        </div>
      </div>
      
      <!-- Create/Edit User Modal -->
      <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click="closeModals">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h3 class="modal-title">
              {{ showCreateModal ? 'Kreiranje novog korisnika' : 'Uređivanje korisnika' }}
            </h3>
            <button class="modal-close" @click="closeModals">×</button>
          </div>
          
          <form @submit.prevent="saveUser">
            <div class="form-row">
              <div class="form-group">
                <label>Ime</label>
                <input 
                  type="text" 
                  v-model="userForm.firstName" 
                  placeholder="Unesite ime" 
                  required
                >
              </div>
              
              <div class="form-group">
                <label>Prezime</label>
                <input 
                  type="text" 
                  v-model="userForm.lastName" 
                  placeholder="Unesite prezime" 
                  required
                >
              </div>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label>Korisničko ime</label>
                <input 
                  type="text" 
                  v-model="userForm.username" 
                  placeholder="Unesite korisničko ime" 
                  required
                >
              </div>
              
              <div class="form-group">
                <label>Email</label>
                <input 
                  type="email" 
                  v-model="userForm.email" 
                  placeholder="Unesite email" 
                  required
                >
              </div>
            </div>
            
            <div v-if="showCreateModal" class="form-group">
              <label>Lozinka</label>
              <input 
                type="password" 
                v-model="userForm.password" 
                placeholder="Unesite lozinku" 
                :required="showCreateModal"
              >
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label>Uloga</label>
                <select v-model="userForm.role" required>
                  <option value="user">Korisnik</option>
                  <option value="researcher">Istraživač</option>
                  <option value="manager">Menadžer</option>
                  <option value="admin">Administrator</option>
                </select>
              </div>
              
              <div class="form-group">
                <label>Status</label>
                <select v-model="userForm.status">
                  <option value="active">Aktivni</option>
                  <option value="inactive">Neaktivni</option>
                  <option value="pending">Na čekanju</option>
                </select>
              </div>
            </div>
            
            <div class="form-group">
              <label>Telefon (opciono)</label>
              <input 
                type="tel" 
                v-model="userForm.phone" 
                placeholder="Unesite broj telefona"
              >
            </div>
            
            <div class="form-group">
              <label>Departman (opciono)</label>
              <input 
                type="text" 
                v-model="userForm.department" 
                placeholder="Unesite departman"
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
const selectedUser = ref(null)
const selectedUsers = ref([])
const selectAll = ref(false)

const filters = ref({
  role: '',
  status: '',
  search: ''
})

const userForm = ref({
  firstName: '',
  lastName: '',
  username: '',
  email: '',
  password: '',
  role: 'user',
  status: 'active',
  phone: '',
  department: ''
})

// Mock data
const users = ref([
  {
    id: 1,
    firstName: 'Marko',
    lastName: 'Petrović',
    username: 'marko.petrovic',
    email: 'marko@institut.rs',
    role: 'admin',
    status: 'active',
    phone: '+381 11 123 456',
    department: 'IT Departman',
    lastLogin: '2024-09-09T10:30:00',
    created: '2024-01-15'
  },
  {
    id: 2,
    firstName: 'Ana',
    lastName: 'Jovanović',
    username: 'ana.jovanovic',
    email: 'ana@institut.rs',
    role: 'researcher',
    status: 'active',
    phone: '+381 11 234 567',
    department: 'Istraživanje',
    lastLogin: '2024-09-08T16:45:00',
    created: '2024-02-01'
  },
  {
    id: 3,
    firstName: 'Stefan',
    lastName: 'Nikolić',
    username: 'stefan.nikolic',
    email: 'stefan@institut.rs',
    role: 'manager',
    status: 'active',
    phone: '+381 11 345 678',
    department: 'Upravljanje projektima',
    lastLogin: '2024-09-07T09:15:00',
    created: '2024-01-20'
  },
  {
    id: 4,
    firstName: 'Milica',
    lastName: 'Stojković',
    username: 'milica.stojkovic',
    email: 'milica@institut.rs',
    role: 'user',
    status: 'inactive',
    phone: '+381 11 456 789',
    department: 'Marketing',
    lastLogin: '2024-08-25T14:20:00',
    created: '2024-03-10'
  },
  {
    id: 5,
    firstName: 'Jovana',
    lastName: 'Mitrović',
    username: 'jovana.mitrovic',
    email: 'jovana@institut.rs',
    role: 'researcher',
    status: 'pending',
    phone: '+381 11 567 890',
    department: 'Istraživanje',
    lastLogin: null,
    created: '2024-09-05'
  }
])

// Computed
const filteredUsers = computed(() => {
  let filtered = users.value

  // Apply filters
  if (filters.value.role) {
    filtered = filtered.filter(u => u.role === filters.value.role)
  }
  
  if (filters.value.status) {
    filtered = filtered.filter(u => u.status === filters.value.status)
  }
  
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(u => 
      u.firstName.toLowerCase().includes(search) ||
      u.lastName.toLowerCase().includes(search) ||
      u.username.toLowerCase().includes(search) ||
      u.email.toLowerCase().includes(search)
    )
  }

  return filtered
})

const totalUsers = computed(() => users.value.length)
const activeUsers = computed(() => users.value.filter(u => u.status === 'active').length)
const adminUsers = computed(() => users.value.filter(u => u.role === 'admin').length)

// Methods
function getRoleText(role) {
  const roleMap = {
    admin: 'Administrator',
    manager: 'Menadžer',
    researcher: 'Istraživač',
    user: 'Korisnik'
  }
  return roleMap[role] || role
}

function getStatusText(status) {
  const statusMap = {
    active: 'Aktivni',
    inactive: 'Neaktivni',
    pending: 'Na čekanju'
  }
  return statusMap[status] || status
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString('sr-RS')
}

function formatDateTime(dateString) {
  if (!dateString) return 'Nikad'
  return new Date(dateString).toLocaleString('sr-RS')
}

function editUser(user) {
  selectedUser.value = user
  userForm.value = {
    firstName: user.firstName,
    lastName: user.lastName,
    username: user.username,
    email: user.email,
    password: '',
    role: user.role,
    status: user.status,
    phone: user.phone || '',
    department: user.department || ''
  }
  showEditModal.value = true
}

function deleteUser(user) {
  if (confirm(`Da li ste sigurni da želite da obrišete korisnika "${user.firstName} ${user.lastName}"?`)) {
    const index = users.value.findIndex(u => u.id === user.id)
    if (index > -1) {
      users.value.splice(index, 1)
    }
  }
}

function toggleUserStatus(user) {
  const newStatus = user.status === 'active' ? 'inactive' : 'active'
  const index = users.value.findIndex(u => u.id === user.id)
  if (index > -1) {
    users.value[index].status = newStatus
  }
}

function resetPassword(user) {
  if (confirm(`Da li ste sigurni da želite da resetujete lozinku za korisnika "${user.firstName} ${user.lastName}"?`)) {
    // Here you would call the backend API to reset password
    alert('Nova lozinka je poslata na email korisnika.')
  }
}

function saveUser() {
  if (showCreateModal.value) {
    // Create new user
    const newUser = {
      id: Date.now(),
      ...userForm.value,
      lastLogin: null,
      created: new Date().toISOString().split('T')[0]
    }
    users.value.push(newUser)
  } else if (showEditModal.value && selectedUser.value) {
    // Update existing user
    const index = users.value.findIndex(u => u.id === selectedUser.value.id)
    if (index > -1) {
      users.value[index] = {
        ...users.value[index],
        ...userForm.value
      }
    }
  }
  
  closeModals()
}

function closeModals() {
  showCreateModal.value = false
  showEditModal.value = false
  selectedUser.value = null
  userForm.value = {
    firstName: '',
    lastName: '',
    username: '',
    email: '',
    password: '',
    role: 'user',
    status: 'active',
    phone: '',
    department: ''
  }
}

function toggleSelectAll() {
  if (selectAll.value) {
    selectedUsers.value = filteredUsers.value.map(u => u.id)
  } else {
    selectedUsers.value = []
  }
}

function bulkActivate() {
  if (confirm(`Da li ste sigurni da želite da aktivirate ${selectedUsers.value.length} korisnika?`)) {
    selectedUsers.value.forEach(userId => {
      const index = users.value.findIndex(u => u.id === userId)
      if (index > -1) {
        users.value[index].status = 'active'
      }
    })
    selectedUsers.value = []
    selectAll.value = false
  }
}

function bulkDeactivate() {
  if (confirm(`Da li ste sigurni da želite da deaktivirate ${selectedUsers.value.length} korisnika?`)) {
    selectedUsers.value.forEach(userId => {
      const index = users.value.findIndex(u => u.id === userId)
      if (index > -1) {
        users.value[index].status = 'inactive'
      }
    })
    selectedUsers.value = []
    selectAll.value = false
  }
}

function bulkDelete() {
  if (confirm(`Da li ste sigurni da želite da obrišete ${selectedUsers.value.length} korisnika?`)) {
    users.value = users.value.filter(u => !selectedUsers.value.includes(u.id))
    selectedUsers.value = []
    selectAll.value = false
  }
}

function exportUsers() {
  console.log('Exporting users...')
  // Here you would implement export functionality
}

// Lifecycle
onMounted(() => {
  console.log('Users loaded')
})
</script>

<style scoped>
.users {
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

/* Users Overview */
.users-overview {
  margin-bottom: 30px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-card.total {
  border-left: 4px solid #3498db;
}

.stat-card.active {
  border-left: 4px solid #2ecc71;
}

.stat-card.admin {
  border-left: 4px solid #9b59b6;
}

.stat-content h3 {
  color: #2c3e50;
  font-size: 28px;
  margin-bottom: 5px;
  font-weight: 700;
}

.stat-content p {
  color: #7f8c8d;
  font-size: 14px;
  font-weight: 500;
}

.stat-icon {
  font-size: 28px;
  opacity: 0.8;
}

.filters {
  display: flex;
  gap: 20px;
  align-items: flex-end;
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

/* Users Table */
.users-table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #ecf0f1;
}

.table-header h3 {
  color: #2c3e50;
  font-size: 18px;
}

.table-actions {
  display: flex;
  gap: 15px;
}

.table-wrapper {
  overflow-x: auto;
}

.user-row {
  transition: background 0.3s;
}

.user-row:hover {
  background: #f8f9fa;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3498db, #2ecc71);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

.user-name {
  font-weight: 600;
  color: #2c3e50;
  font-size: 14px;
}

.user-username {
  font-size: 12px;
  color: #7f8c8d;
}

/* Role badges */
.role-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.role-admin {
  background: #f8d7da;
  color: #721c24;
}

.role-manager {
  background: #fff3cd;
  color: #856404;
}

.role-researcher {
  background: #d1ecf1;
  color: #0c5460;
}

.role-user {
  background: #e9ecef;
  color: #495057;
}

/* Status badges */
.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
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
  padding: 6px;
  border-radius: 4px;
  transition: background 0.3s;
}

.btn-icon-small:hover {
  background: #f8f9fa;
}

.btn-icon-small.danger:hover {
  background: #f8d7da;
}

/* Empty State */
.empty-state {
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

/* Bulk Actions */
.bulk-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: #f8f9fa;
  border-top: 1px solid #ecf0f1;
}

.selected-count {
  font-weight: 600;
  color: #2c3e50;
}

.bulk-buttons {
  display: flex;
  gap: 15px;
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
  .users {
    padding: 20px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 20px;
    align-items: stretch;
  }
  
  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
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
  
  .table-header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
  
  .bulk-actions {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
  
  .bulk-buttons {
    justify-content: center;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column-reverse;
  }
}

@media (max-width: 480px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }
  
  .table {
    min-width: 700px;
  }
}
</style>
