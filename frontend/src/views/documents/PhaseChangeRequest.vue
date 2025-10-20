<template>
  <div class="phase-request-page">
    <div class="phase-request-container">
      <h2>Zahtevaj promenu faze</h2>
      
      <div class="document-info">
        <div class="info-row">
          <span class="label">Dokument:</span>
          <span class="value">{{ document?.naziv_dokumenta || 'Učitavanje...' }}</span>
        </div>
        <div class="info-row">
          <span class="label">Trenutna verzija:</span>
          <span class="value">{{ latestVersion?.verzija_oznaka || 'N/A' }}</span>
        </div>
      </div>

      <div class="phase-selector">
        <label for="newPhase">Izaberite novu fazu:</label>
        <select id="newPhase" v-model="selectedPhaseId">
          <option :value="null" disabled>-- Izaberite fazu --</option>
          <option v-for="phase in phases" :key="phase.faza_id" :value="phase.faza_id">
            {{ phase.naziv_faze }}
          </option>
        </select>
      </div>

      <div class="phase-transition" v-if="selectedPhaseId && currentPhase">
        <div class="phase-box current">
          <span class="phase-label">Trenutna faza</span>
          <span class="phase-name">{{ currentPhase.naziv_faze }}</span>
        </div>
        <div class="arrow">→</div>
        <div class="phase-box new">
          <span class="phase-label">Nova faza</span>
          <span class="phase-name">{{ selectedPhase?.naziv_faze }}</span>
        </div>
      </div>

      <div class="actions">
        <button @click="submitRequest" class="btn-submit" :disabled="!selectedPhaseId">
          Pošalji zahtev
        </button>
        <button @click="goBack" class="btn-cancel">
          Otkaži
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const documentId = ref(Number(route.params.id))
const document = ref(null)
const versions = ref([])
const phases = ref([])
const selectedPhaseId = ref(null)
const pendingRequests = ref([])
const currentUser = ref(null)

const latestVersion = computed(() => {
  if (!versions.value || versions.value.length === 0) return null
  return versions.value[0] // Assuming sorted by date desc
})

const currentPhase = computed(() => {
  if (!document.value || !document.value.trenutna_faza_id) return null
  return phases.value.find(p => p.faza_id === document.value.trenutna_faza_id)
})

const selectedPhase = computed(() => {
  if (!selectedPhaseId.value) return null
  return phases.value.find(p => p.faza_id === selectedPhaseId.value)
})

async function loadDocument() {
  try {
    document.value = await window.go?.main?.App?.GetDocumentByID(documentId.value)
  } catch (e) {
    console.error('Error loading document:', e)
    alert('Greška pri učitavanju dokumenta')
  }
}

async function loadVersions() {
  try {
    versions.value = await window.go?.main?.App?.GetDocumentVersions(documentId.value) || []
    // Sort by date descending
    versions.value.sort((a, b) => new Date(b.datuma_postavke) - new Date(a.datuma_postavke))
  } catch (e) {
    console.error('Error loading versions:', e)
  }
}

async function loadPhases() {
  try {
    if (!document.value || !document.value.radni_tok_id) return
    phases.value = await window.go?.main?.App?.GetWorkflowPhases(document.value.radni_tok_id) || []
  } catch (e) {
    console.error('Error loading phases:', e)
  }
}

async function loadPendingRequests() {
  try {
    const allRequests = await window.go?.main?.App?.ListPhaseChangeRequestsByDocument(documentId.value) || []
    // Filter to show only current user's requests
    if (currentUser.value) {
      pendingRequests.value = allRequests.filter(r => r.podnosilac_zahteva_id === currentUser.value.korisnik_id)
    }
  } catch (e) {
    console.error('Error loading requests:', e)
  }
}

async function loadCurrentUser() {
  try {
    currentUser.value = await window.go?.main?.App?.GetCurrentUser()
  } catch (e) {
    console.error('Error loading current user:', e)
  }
}

async function submitRequest() {
  if (!selectedPhaseId.value || !currentPhase.value) {
    alert('Molimo izaberite fazu')
    return
  }

  if (selectedPhaseId.value === currentPhase.value.faza_id) {
    alert('Izabrana faza je već trenutna faza dokumenta')
    return
  }

  try {
    const request = {
      dokument_id: documentId.value,
      zadatak_id: null,
      podnosilac_zahteva_id: null, // backend sets this
      zahtevana_faza_id: selectedPhaseId.value,
      status: 'Na čekanju',
      komentar: null
    }
    await window.go?.main?.App?.CreatePhaseChangeRequest(request)
    alert('Zahtev za promenu faze je uspešno poslat!')
    // Reload requests and reset selection
    await loadPendingRequests()
    selectedPhaseId.value = null
  } catch (e) {
    console.error('Error submitting request:', e)
    alert('Greška pri slanju zahteva: ' + (e?.message || 'Nepoznata greška'))
  }
}

function formatRequestText(req) {
  if (!currentUser.value) return ''
  const oldPhase = phases.value.find(p => p.faza_id === req.prethodna_faza_id)
  const newPhase = phases.value.find(p => p.faza_id === req.nova_faza_id)
  
  return `${currentUser.value.ime} ${currentUser.value.prezime} želi da promeni fazu iz ${oldPhase?.naziv_faze || 'N/A'} u ${newPhase?.naziv_faze || 'N/A'}`
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('sr-RS', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function goBack() {
  router.back()
}

onMounted(async () => {
  await loadCurrentUser()
  await loadDocument()
  await loadVersions()
  await loadPhases()
  await loadPendingRequests()
})
</script>

<style scoped>
.phase-request-page {
  display: flex;
  justify-content: center;
  padding: 2rem;
  background: #f8fafc;
  min-height: 100vh;
}

.phase-request-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  padding: 2rem;
  max-width: 700px;
  width: 100%;
}

h2 {
  color: #1e293b;
  margin-bottom: 1.5rem;
  font-size: 1.75rem;
}

.document-info {
  background: #f1f5f9;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.info-row {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.info-row:last-child {
  margin-bottom: 0;
}

.label {
  font-weight: 600;
  color: #475569;
}

.value {
  color: #1e293b;
}

.phase-selector {
  margin-bottom: 1.5rem;
}

.phase-selector label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #334155;
}

.phase-selector select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 1rem;
  background: white;
  cursor: pointer;
}

.phase-selector select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.phase-transition {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1.5rem;
  margin: 2rem 0;
  padding: 1.5rem;
  background: #f8fafc;
  border-radius: 8px;
}

.phase-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  min-width: 150px;
}

.phase-box.current {
  background: #dbeafe;
  border: 2px solid #3b82f6;
}

.phase-box.new {
  background: #dcfce7;
  border: 2px solid #22c55e;
}

.phase-label {
  font-size: 0.875rem;
  color: #64748b;
  margin-bottom: 0.25rem;
}

.phase-name {
  font-weight: 600;
  font-size: 1.125rem;
  color: #1e293b;
}

.arrow {
  font-size: 2rem;
  color: #64748b;
  font-weight: bold;
}

.actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}

.btn-submit, .btn-cancel {
  padding: 0.75rem 1.5rem;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  font-size: 1rem;
  transition: all 0.2s;
}

.btn-submit {
  background: #3b82f6;
  color: white;
  flex: 1;
}

.btn-submit:hover:not(:disabled) {
  background: #2563eb;
}

.btn-submit:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
}

.btn-cancel {
  background: #e2e8f0;
  color: #475569;
}

.btn-cancel:hover {
  background: #cbd5e1;
}

.pending-requests {
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 2px solid #e2e8f0;
}

.pending-requests h3 {
  color: #1e293b;
  margin-bottom: 1rem;
  font-size: 1.25rem;
}

.request-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.request-item {
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: white;
}

.request-item.status-na-čekanju {
  border-left: 4px solid #f59e0b;
}

.request-item.status-odobren {
  border-left: 4px solid #22c55e;
}

.request-item.status-odbijen {
  border-left: 4px solid #ef4444;
}

.request-text {
  color: #1e293b;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.request-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-weight: 600;
  font-size: 0.75rem;
}

.status-na-čekanju .status-badge {
  background: #fef3c7;
  color: #92400e;
}

.status-odobren .status-badge {
  background: #dcfce7;
  color: #166534;
}

.status-odbijen .status-badge {
  background: #fee2e2;
  color: #991b1b;
}

.date {
  color: #64748b;
}
</style>
