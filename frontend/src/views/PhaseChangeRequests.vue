<template>
  <Layout>
    <div class="phase-change-requests">
      <!-- Header -->
      <div class="page-header">
        <div>
          <h2>Lista Zahteva</h2>
          <div class="breadcrumb">Početna > Zadaci > Lista Zahteva</div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" @click="goBack">
            ← Nazad na Zadatke
          </button>
          <button class="btn btn-primary" @click="loadRequests">
            🔄 Osveži
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state card">
        <p>Učitavanje zahteva...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="requests.length === 0" class="empty-state card">
        <div class="empty-icon">📋</div>
        <h3>Nema aktivnih zahteva</h3>
        <p>Trenutno nemate zahteva za promenu faze koji čekaju na pregled.</p>
      </div>

      <!-- Requests List -->
      <div v-else class="requests-list">
        <div 
          v-for="request in requests" 
          :key="request.zahtev_id"
          class="request-card card"
        >
          <div class="request-header">
            <div class="request-title-section">
              <h3 class="request-title">{{ request.naziv_zadatka }}</h3>
              <div class="request-meta">
                <span class="meta-item">
                  <span class="meta-icon">📁</span>
                  {{ request.naziv_projekta }}
                </span>
                <span class="meta-item">
                  <span class="meta-icon">👤</span>
                  Podnosilac: {{ request.podnosilac_ime }}
                </span>
                <span class="meta-item">
                  <span class="meta-icon">📅</span>
                  {{ formatDate(request.datum_kreiranja) }}
                </span>
              </div>
            </div>
          </div>

          <div class="request-body">
            <div class="request-info">
              <div class="info-item">
                <label>Zahtevana faza:</label>
                <span class="phase-badge">{{ request.naziv_faze }}</span>
              </div>
              <div v-if="request.komentar" class="info-item full-width">
                <label>Opis zahteva:</label>
                <p class="request-description">{{ request.komentar }}</p>
              </div>
              <div v-else class="info-item full-width">
                <p class="no-description">Podnosilac nije dodao opis uz zahtev.</p>
              </div>
            </div>
          </div>

          <div class="request-actions">
            <button 
              class="btn btn-success"
              @click="approveRequest(request)"
              :disabled="processingRequest === request.zahtev_id"
            >
              <span v-if="processingRequest === request.zahtev_id">⏳ Obrađujem...</span>
              <span v-else>✓ Prihvati</span>
            </button>
            <button 
              class="btn btn-danger"
              @click="rejectRequest(request)"
              :disabled="processingRequest === request.zahtev_id"
            >
              <span v-if="processingRequest === request.zahtev_id">⏳ Obrađujem...</span>
              <span v-else>✗ Odbij</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/auth.js'

const router = useRouter()
const authStore = useAuthStore()

// State
const requests = ref([])
const loading = ref(false)
const processingRequest = ref(null)

// Methods
async function loadRequests() {
  loading.value = true
  try {
    const { GetManagerPhaseChangeRequests } = window.go.main.App
    const result = await GetManagerPhaseChangeRequests()
    requests.value = result || []
  } catch (error) {
    console.error('Error loading requests:', error)
    alert('Greška pri učitavanju zahteva: ' + (error.message || error))
  } finally {
    loading.value = false
  }
}

async function approveRequest(request) {
  if (!confirm(`Da li ste sigurni da želite da prihvatite zahtev za zadatak "${request.naziv_zadatka}"?\n\nZadatak će biti premešten u fazu: ${request.naziv_faze}`)) {
    return
  }

  processingRequest.value = request.zahtev_id
  try {
    const { ApprovePhaseChangeRequest } = window.go.main.App
    await ApprovePhaseChangeRequest(request.zahtev_id)
    
    alert(`Zahtev je uspešno prihvaćen! Zadatak "${request.naziv_zadatka}" je premešten u fazu "${request.naziv_faze}".`)
    
    // Remove from list
    requests.value = requests.value.filter(r => r.zahtev_id !== request.zahtev_id)
  } catch (error) {
    console.error('Error approving request:', error)
    alert('Greška pri prihvatanju zahteva: ' + (error.message || error))
  } finally {
    processingRequest.value = null
  }
}

async function rejectRequest(request) {
  if (!confirm(`Da li ste sigurni da želite da odbijete zahtev za zadatak "${request.naziv_zadatka}"?`)) {
    return
  }

  processingRequest.value = request.zahtev_id
  try {
    const { RejectPhaseChangeRequest } = window.go.main.App
    await RejectPhaseChangeRequest(request.zahtev_id)
    
    alert(`Zahtev je odbijen.`)
    
    // Remove from list
    requests.value = requests.value.filter(r => r.zahtev_id !== request.zahtev_id)
  } catch (error) {
    console.error('Error rejecting request:', error)
    alert('Greška pri odbijanju zahteva: ' + (error.message || error))
  } finally {
    processingRequest.value = null
  }
}

function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('sr-RS', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function goBack() {
  router.push('/tasks')
}

// Lifecycle
onMounted(async () => {
  // Check if user is a manager
  if (!authStore.user || (authStore.user.naziv_uloge !== 'Rukovodilac projekta' && authStore.user.naziv_uloge !== 'Administrator')) {
    alert('Samo rukovodioci projekata mogu pristupiti listi zahteva.')
    router.push('/tasks')
    return
  }

  await loadRequests()
})
</script>

<style scoped>
.phase-change-requests {
  padding: 30px;
  max-width: 1200px;
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

/* Loading and Empty States */
.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #7f8c8d;
}

.loading-state p {
  margin: 0;
  font-size: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.empty-icon {
  font-size: 64px;
  opacity: 0.3;
}

.empty-state h3 {
  color: #2c3e50;
  margin: 0;
  font-size: 20px;
}

.empty-state p {
  color: #7f8c8d;
  margin: 0;
  font-size: 14px;
}

/* Requests List */
.requests-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.request-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s ease;
  border-left: 4px solid #3498db;
}

.request-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.request-header {
  padding: 20px;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.request-title-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.request-title {
  color: #2c3e50;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.request-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  font-size: 13px;
  color: #7f8c8d;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-icon {
  font-size: 14px;
}

.request-body {
  padding: 20px;
}

.request-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-item label {
  font-size: 12px;
  font-weight: 600;
  color: #7f8c8d;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.phase-badge {
  display: inline-block;
  padding: 6px 12px;
  background: #3498db;
  color: white;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
}

.request-description {
  margin: 0;
  color: #2c3e50;
  line-height: 1.6;
  font-size: 14px;
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
  border-left: 3px solid #3498db;
}

.no-description {
  margin: 0;
  color: #95a5a6;
  font-style: italic;
  font-size: 14px;
}

.request-actions {
  padding: 15px 20px;
  background: #f8f9fa;
  border-top: 1px solid #e9ecef;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-success {
  background: #2ecc71;
  color: white;
  border: none;
  min-width: 120px;
}

.btn-success:hover:not(:disabled) {
  background: #27ae60;
}

.btn-success:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}

.btn-danger {
  background: #e74c3c;
  color: white;
  border: none;
  min-width: 120px;
}

.btn-danger:hover:not(:disabled) {
  background: #c0392b;
}

.btn-danger:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}

/* Responsive */
@media (max-width: 768px) {
  .phase-change-requests {
    padding: 20px;
  }

  .page-header {
    flex-direction: column;
    gap: 20px;
    align-items: stretch;
  }

  .header-actions {
    flex-direction: column;
  }

  .request-meta {
    flex-direction: column;
    gap: 8px;
  }

  .request-info {
    grid-template-columns: 1fr;
  }

  .request-actions {
    flex-direction: column;
  }

  .btn-success,
  .btn-danger {
    width: 100%;
  }
}
</style>
