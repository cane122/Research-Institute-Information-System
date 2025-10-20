<template>
  <Layout>
    <div class="doc-full-preview">
      <div class="main-content">
        <!-- Naziv dokumenta -->
        <div class="section naziv">
          <h2 class="naziv-dokumenta">{{ document?.naziv_dokumenta || '-' }}</h2>
        </div>
        <!-- Opis dokumenta -->
        <div class="section opis">
          <label class="opis-label">Opis:</label>
          <template v-if="isLeader">
            <textarea class="opis-input" v-model="editableOpis" rows="4" placeholder="Unesite opis dokumenta..."></textarea>
            <div class="opis-actions">
              <button class="btn btn-secondary" @click="resetOpis" :disabled="savingOpis">Poništi</button>
              <button class="btn btn-primary" @click="saveOpis" :disabled="savingOpis">
                {{ savingOpis ? 'Čuvanje…' : 'Sačuvaj opis' }}
              </button>
            </div>
          </template>
          <template v-else>
            <div class="opis-text">{{ document?.opis || 'Nema opisa.' }}</div>
          </template>
          <!-- Checklist zadacici -->
          <div class="zadacici-section">
            <label class="zadacici-label">Zadaci (checklista):</label>
            <ul class="zadacici-list">
              <li v-for="z in zadacici" :key="z.zadacic_id" :class="{ 'zadacic-izvrsen': z.izvrsen }">
                <input type="checkbox" :checked="z.izvrsen" :disabled="!isLeader" @change="toggleZadacic(z)" />
                <span>{{ z.opis }}</span>
                <button v-if="isLeader" class="zadacic-delete" @click="deleteZadacic(z.zadacic_id)">🗑️</button>
              </li>
            </ul>
            <div v-if="isLeader" class="zadacic-add-row">
              <input v-model="noviZadacic" placeholder="Novi zadatak..." @keyup.enter="addZadacic" />
              <button class="btn btn-primary" @click="addZadacic" :disabled="!noviZadacic.trim()">Dodaj</button>
            </div>
          </div>
        </div>

        <!-- Faze dokumenta -->
        <div class="section faze">
          <label class="faze-label">Faze dokumenta:</label>
          <div class="faze-list">
            <div v-for="faza in faze" :key="faza.faza_id"
                 :class="['faza-item', { current: faza.faza_id === document?.trenutna_faza_id }]">
              {{ faza.naziv_faze }}
            </div>
          </div>
        </div>

        <!-- Lista korisnika (dozvole) -->
        <div class="section korisnici">
          <div class="korisnici-header">
            <label class="korisnici-label">Korisnici na dokumentu:</label>
            <button v-if="isLeader" class="btn btn-secondary" @click="openAddUser">+ Dodaj korisnika</button>
          </div>
          <div class="korisnici-list">
            <div v-for="u in docUsers" :key="u.korisnik_id" class="korisnik-item">
              <span>
                {{ (u.ime ? u.ime + ' ' + (u.prezime || '') : u.korisnicko_ime) }}
              </span>
              <button v-if="isLeader" class="x-btn" title="Ukloni korisnika" @click="removeUser(u.korisnik_id)">✕</button>
            </div>
            <div v-if="docUsers.length === 0" class="hint">Nema korisnika sa dozvolama za ovaj dokument.</div>
          </div>
        </div>

        <!-- Verzije dokumenta -->
        <div class="section verzije">
          <label class="verzije-label">Verzije dokumenta:</label>
          <div class="verzije-panel">
            <div class="file-info">
              <span class="file-name">{{ selectedVersion?.putanja_do_fajla?.split('/').pop() || document?.naziv_dokumenta }}</span>
              <select v-model="selectedVerzijaId" @change="onVersionChange">
                <option v-for="v in verzije" :key="v.verzija_id" :value="v.verzija_id">
                  {{ v.verzija_oznaka || ('v' + v.verzija_id) }}
                </option>
              </select>
              <button class="btn btn-secondary" @click="downloadFile" :disabled="!selectedVersion">Preuzmi</button>
              <button class="btn btn-secondary" @click="openEditModal">Izmeni</button>
            </div>
            <div class="version-meta">
              <span>Izmenio: {{ selectedVersionUser }}</span>
              <span>Datum izmene: {{ formatDateTime(selectedVersion?.datuma_postavke) }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="side-panel">
        <div class="side-section">
          <div><b>Radni tok:</b> {{ workflow?.naziv || '-' }}</div>
          <div><b>Trenutna faza:</b> {{ currentFaza?.naziv_faze || '-' }}</div>
          <button v-if="isLeader" class="btn btn-primary" @click="openPhaseChangeModal">Promeni fazu</button>
          <button v-else class="btn btn-primary" @click="goToPhaseChangeRequest">Zahtevaj promenu faze</button>
          <!-- Phase change requests notification -->
          <div v-if="phaseRequests.length > 0" class="phase-requests-box">
            <div class="phase-requests-title">Zahtevi za promenu faze:</div>
            <div v-for="req in phaseRequests" :key="req.zahtev_id" class="phase-request-item">
              <div class="phase-request-info">
                <span class="phase-request-status" :class="'status-' + (req.status || '').toLowerCase().replace(/\\s+/g, '-')">{{ req.status }}</span>
                <span class="phase-request-text">{{ formatPhaseRequestText(req) }}</span>
                <span class="phase-request-date">{{ formatDateTime(req.datum_kreiranja) }}</span>
              </div>
              <div v-if="isLeader && isStatusPending(req.status)" class="phase-request-actions">
                <button class="btn-approve" @click="approveRequest(req.zahtev_id)">Odobri</button>
                <button class="btn-reject" @click="rejectRequest(req.zahtev_id)">Odbij</button>
              </div>
            </div>
          </div>
          <div><b>Rok za završetak:</b> {{ formatDate(document?.rok) }}</div>
        </div>
      </div>
    </div>
  </Layout>

  <!-- Edit (New Version) Modal -->
  <div v-if="showEdit" class="modal-overlay" @click="closeEdit">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h3>Nova verzija dokumenta</h3>
        <button class="modal-close" @click="closeEdit">×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>Naziv dokumenta</label>
          <input type="text" :value="document?.naziv_dokumenta || ''" disabled />
        </div>
        <div class="form-group">
          <label>Oznaka verzije (opciono)</label>
          <input type="text" v-model="newVersionLabel" placeholder="npr. v2.0 ili Rev A" />
        </div>
        <div class="form-group">
          <label>Novi fajl</label>
          <input type="file" @change="onFilePick" />
          <div v-if="pickedFileName" class="hint">Izabrano: {{ pickedFileName }}</div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="closeEdit">Otkaži</button>
        <button class="btn btn-primary" :disabled="!pickedFile" @click="saveNewVersion">Sačuvaj verziju</button>
      </div>
    </div>
  </div>

  <!-- Add User Modal -->
  <div v-if="showAddUser" class="modal-overlay" @click="closeAddUser">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h3>Dodaj korisnika na dokument</h3>
        <button class="modal-close" @click="closeAddUser">×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>Izaberite korisnika</label>
          <select v-model="selectedUserToAdd">
            <option value="">-- Izaberite --</option>
            <option v-for="m in availableMembers" :key="m.korisnik_id" :value="m.korisnik_id">
              {{ m.ime || m.korisnicko_ime }} {{ m.prezime || '' }}
            </option>
          </select>
          <div v-if="availableMembers.length === 0" class="hint">Nema dostupnih članova tima za dodavanje.</div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="closeAddUser">Otkaži</button>
        <button class="btn btn-primary" :disabled="!selectedUserToAdd" @click="confirmAddUser">Dodaj</button>
      </div>
    </div>
  </div>

  <!-- Phase Change Modal (for leaders) -->
  <div v-if="showPhaseChange" class="modal-overlay" @click="closePhaseChange">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h3>Promeni fazu dokumenta</h3>
        <button class="modal-close" @click="closePhaseChange">×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>Trenutna faza</label>
          <input type="text" :value="currentFaza?.naziv_faze || '-'" disabled />
        </div>
        <div class="form-group">
          <label>Nova faza</label>
          <select v-model="selectedNewPhaseId">
            <option value="">-- Izaberite fazu --</option>
            <option v-for="f in faze" :key="f.faza_id" :value="f.faza_id" :disabled="f.faza_id === document?.trenutna_faza_id">
              {{ f.naziv_faze }}
            </option>
          </select>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="closePhaseChange">Otkaži</button>
        <button class="btn btn-primary" :disabled="!selectedNewPhaseId" @click="confirmPhaseChange">Promeni fazu</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Layout from '../../components/Layout.vue'
import { GetDocumentByID, GetDocumentVersions, GetWorkflowPhases, GetProjectMembers, GetWorkflowByID, DownloadDocumentVersion, GetCurrentUser, GetDocumentPermissions, SetDocumentPermission, RemoveDocumentPermission, UpdateDocument, ListZadaciciByDocument, AddZadacic, UpdateZadacicStatus, DeleteZadacic } from '../../../wailsjs/go/main/App.js'

const route = useRoute()
const router = useRouter()
const documentId = Number(route.params.id)

const document = ref(null)
const verzije = ref([])
const selectedVerzijaId = ref(null)
const selectedVersion = computed(() => verzije.value.find(v => v.verzija_id === selectedVerzijaId.value) || verzije.value[0] || null)
const faze = ref([])
// Users & permissions
const projectMembers = ref([])
const docUsers = ref([]) // from GetDocumentPermissions
// Current user & role check
const currentUser = ref(null)
const isLeader = computed(() => (currentUser.value?.uloga_id === 2) || (currentUser.value?.naziv_uloge?.toLowerCase?.() === 'rukovodilac projekta'))
// Editable opis
const editableOpis = ref('')
const savingOpis = ref(false)
const workflow = ref(null)
const currentFaza = computed(() => faze.value.find(f => f.faza_id === document.value?.trenutna_faza_id) || null)

// Zadacici (checklist)
const zadacici = ref([])
const noviZadacic = ref('')
async function loadZadacici() {
  try {
    zadacici.value = await ListZadaciciByDocument(documentId) || []
  } catch {
    zadacici.value = []
  }
}
async function addZadacic() {
  if (!noviZadacic.value.trim()) return
  try {
    await AddZadacic({ dokument_id: documentId, opis: noviZadacic.value.trim() })
    noviZadacic.value = ''
    await loadZadacici()
  } catch (e) {
    alert('Greška pri dodavanju zadatka: ' + (e?.message || e))
  }
}
async function toggleZadacic(z) {
  try {
    await UpdateZadacicStatus(z.zadacic_id, !z.izvrsen)
    await loadZadacici()
  } catch (e) {
    alert('Greška pri izmeni statusa: ' + (e?.message || e))
  }
}
async function deleteZadacic(id) {
  if (!confirm('Obrisati zadatak?')) return
  try {
    await DeleteZadacic(id)
    await loadZadacici()
  } catch (e) {
    alert('Greška pri brisanju zadatka: ' + (e?.message || e))
  }
}

const showEdit = ref(false)
const newVersionLabel = ref('')
const pickedFile = ref(null)
const pickedFileName = ref('')

const selectedVersionUser = computed(() => {
  if (!selectedVersion.value) return '-'
  const user = (docUsers.value || []).find(u => u.korisnik_id === selectedVersion.value.postavio_korisnik_id) || (projectMembers.value || []).find(u => u.korisnik_id === selectedVersion.value.postavio_korisnik_id)
  return user ? (user.ime || user.korisnicko_ime) + ' ' + (user.prezime || '') : '-'
})

function formatDateTime(dt) {
  if (!dt) return '-'
  const d = new Date(dt)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('sr-RS') + ' ' + d.toLocaleTimeString('sr-RS', { hour: '2-digit', minute: '2-digit' })
}
function formatDate(dt) {
  if (!dt) return '-'
  const d = new Date(dt)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('sr-RS')
}
function onVersionChange() {}

async function downloadFile() {
  if (!selectedVersion.value) {
    alert('Molimo izaberite verziju za preuzimanje.')
    return
  }
  try {
    const fileName = selectedVersion.value.putanja_do_fajla?.split('/').pop() ||
                     selectedVersion.value.putanja_do_fajla?.split('\\').pop() ||
                     document.value?.naziv_dokumenta ||
                     'dokument'
    const fileData = await DownloadDocumentVersion(selectedVersion.value.verzija_id)
    const uint8Array = new Uint8Array(fileData)
    const blob = new Blob([uint8Array])
    const url = window.URL.createObjectURL(blob)
    const link = window.document.createElement('a')
    link.href = url
    link.download = fileName
    window.document.body.appendChild(link)
    link.click()
    setTimeout(() => {
      window.URL.revokeObjectURL(url)
      window.document.body.removeChild(link)
    }, 100)
  } catch (error) {
    console.error('Greška pri preuzimanju:', error)
    alert('Greška pri preuzimanju fajla: ' + (error.message || error))
  }
}

function goToPhaseChangeRequest() {
  if (!documentId) return
  router.push({ name: 'PhaseChangeRequest', params: { id: documentId } })
}

// Phase change requests for this document (current user)
const phaseRequests = ref([])

async function loadPhaseRequests() {
  try {
    phaseRequests.value = await window.go?.main?.App?.ListPhaseChangeRequestsByDocument(documentId) || []
  } catch (e) {
    phaseRequests.value = []
  }
}

function formatPhaseRequestText(req) {
  const user = docUsers.value.find(u => u.korisnik_id === req.podnosilac_zahteva_id)
  const newPhase = faze.value.find(f => f.faza_id === req.zahtevana_faza_id)
  const userName = user ? (user.ime || user.korisnicko_ime) + ' ' + (user.prezime || '') : `Korisnik #${req.podnosilac_zahteva_id}`
  return `${userName.trim()} zahteva promenu u fazu: ${newPhase?.naziv_faze || 'N/A'}`
}

function isStatusPending(status) {
  if (!status) return false
  const normalized = status.toLowerCase().trim()
  return normalized === 'na cekanju' || normalized === 'na čekanju' || normalized.includes('ceka') || normalized.includes('čeka')
}

async function approveRequest(requestId) {
  try {
    await window.go?.main?.App?.UpdatePhaseChangeRequestStatus(requestId, 'Odobren', null)
    alert('Zahtev je odobren i faza dokumenta je promenjena.')
    // Reload document to show updated phase
    document.value = await GetDocumentByID(documentId)
    await loadPhaseRequests()
  } catch (e) {
    console.error('Error approving request:', e)
    alert('Greška pri odobravanju zahteva: ' + (e?.message || e))
  }
}

async function rejectRequest(requestId) {
  try {
    await window.go?.main?.App?.UpdatePhaseChangeRequestStatus(requestId, 'Odbijen', null)
    alert('Zahtev je odbijen.')
    await loadPhaseRequests()
  } catch (e) {
    console.error('Error rejecting request:', e)
    alert('Greška pri odbijanju zahteva: ' + (e?.message || e))
  }
}

onMounted(async () => {
  currentUser.value = await GetCurrentUser()
  document.value = await GetDocumentByID(documentId)
  editableOpis.value = document.value?.opis || ''
  verzije.value = await GetDocumentVersions(documentId)
  selectedVerzijaId.value = verzije.value[0]?.verzija_id || null
  if (document.value?.radni_tok_id) {
    faze.value = await GetWorkflowPhases(document.value.radni_tok_id)
    workflow.value = await GetWorkflowByID(document.value.radni_tok_id)
  }
  projectMembers.value = document.value?.projekat_id ? (await GetProjectMembers(document.value.projekat_id)) : []
  await loadDocUsers()
  await loadPhaseRequests()
  await loadZadacici()
})

async function loadDocUsers() {
  try {
    const perms = await GetDocumentPermissions(documentId)
    docUsers.value = perms || []
  } catch (e) {
    docUsers.value = []
  }
}

function openEditModal() { showEdit.value = true }
function closeEdit() {
  showEdit.value = false
  newVersionLabel.value = ''
  pickedFile.value = null
  pickedFileName.value = ''
}
function onFilePick(e) {
  const f = e.target.files && e.target.files[0]
  pickedFile.value = f || null
  pickedFileName.value = f ? f.name : ''
}

// Leader: opis save/reset
function resetOpis() {
  editableOpis.value = document.value?.opis || ''
}
async function saveOpis() {
  if (!document.value) return
  try {
    savingOpis.value = true
    const req = {
      naziv_dokumenta: document.value.naziv_dokumenta,
      projekat_id: document.value.projekat_id || null,
      folder_id: document.value.folder_id || null,
      opis: editableOpis.value || '',
      tip_dokumenta: document.value.tip_dokumenta || 'Document',
      jezik_dokumenta: document.value.jezik_dokumenta || 'Serbian',
      tagovi: []
    }
    await UpdateDocument(document.value.dokument_id, req)
    document.value = await GetDocumentByID(documentId)
    editableOpis.value = document.value?.opis || ''
  } catch (e) {
    alert('Greška pri čuvanju opisa: ' + (e.message || e))
  } finally {
    savingOpis.value = false
  }
}

// Leader: manage users
const showAddUser = ref(false)
const selectedUserToAdd = ref('')
const availableMembers = computed(() => {
  const existing = new Set((docUsers.value || []).map(u => u.korisnik_id))
  return (projectMembers.value || []).filter(m => !existing.has(m.korisnik_id))
})
function openAddUser() { showAddUser.value = true }
function closeAddUser() { showAddUser.value = false; selectedUserToAdd.value = '' }
async function confirmAddUser() {
  try {
    await SetDocumentPermission({
      dokument_id: documentId,
      korisnik_id: Number(selectedUserToAdd.value),
      moze_citati: true,
      moze_menjati: true,
      moze_brisati: false
    })
    await loadDocUsers()
    closeAddUser()
  } catch (e) {
    alert('Greška pri dodavanju korisnika: ' + (e.message || e))
  }
}

// Leader: direct phase change
const showPhaseChange = ref(false)
const selectedNewPhaseId = ref('')
function openPhaseChangeModal() { showPhaseChange.value = true }
function closePhaseChange() { 
  showPhaseChange.value = false
  selectedNewPhaseId.value = ''
}
async function confirmPhaseChange() {
  if (!selectedNewPhaseId.value) return
  try {
    await window.go?.main?.App?.ChangeDocumentPhase(documentId, Number(selectedNewPhaseId.value))
    alert('Faza dokumenta je uspešno promenjena.')
    // Reload document to show updated phase
    document.value = await GetDocumentByID(documentId)
    await loadPhaseRequests()
    closePhaseChange()
  } catch (e) {
    console.error('Error changing phase:', e)
    alert('Greška pri promeni faze: ' + (e?.message || e))
  }
}

async function removeUser(userId) {
  try {
    if (!confirm('Ukloniti korisnika sa dokumenata?')) return
    await RemoveDocumentPermission(documentId, userId)
    await loadDocUsers()
  } catch (e) {
    alert('Greška pri uklanjanju korisnika: ' + (e.message || e))
  }
}
async function saveNewVersion() {
  if (!pickedFile.value || !document.value) return
  try {
    const buf = await fileToByteArray(pickedFile.value)
    const label = newVersionLabel.value ? newVersionLabel.value : null
    await window.go.main.App.UploadDocumentVersion(document.value.dokument_id, label, buf, pickedFile.value.name)
    verzije.value = await GetDocumentVersions(document.value.dokument_id)
    selectedVerzijaId.value = verzije.value[0]?.verzija_id || selectedVerzijaId.value
    closeEdit()
    alert('Nova verzija je uspešno sačuvana.')
  } catch (err) {
    console.error('Greška pri čuvanju verzije:', err)
    alert('Greška pri čuvanju verzije: ' + (err.message || err))
  }
}
function fileToByteArray(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (ev) => {
      const arr = new Uint8Array(ev.target.result)
      resolve(Array.from(arr))
    }
    reader.onerror = (e) => reject(e)
    reader.readAsArrayBuffer(file)
  })
}
</script>

<style scoped>
.doc-full-preview {
  display: flex;
  gap: 32px;
  align-items: flex-start;
}
.main-content {
  flex: 2;
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin: 16px 16px
}
.side-panel {
  flex: 1;
  background: #f3f4f6;
  border-radius: 8px;
  padding: 24px 16px;
  min-width: 260px;
  max-width: 320px;
  box-shadow: 0 2px 8px #0001;
}
.section {
  margin-bottom: 12px;
}
.opis-label, .faze-label, .korisnici-label, .verzije-label {
  font-weight: 600;
  margin-bottom: 6px;
  display: block;
}
.opis-text {
  background: #f9fafb;
  border-radius: 6px;
  padding: 10px 14px;
  color: #374151;
}
.opis-input {
  width: 100%;
  min-height: 100px;
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  resize: vertical;
}
.opis-actions { display: flex; gap: 10px; margin-top: 8px; }
.faze-list {
  display: flex;
  gap: 10px;
  margin-top: 6px;
  flex-wrap: wrap;
}
.faza-item {
  padding: 6px 16px;
  border-radius: 16px;
  background: #e5e7eb;
  color: #374151;
  font-weight: 500;
}
.faza-item.current {
  background: #2563eb;
  color: #fff;
}
.korisnici-list {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-top: 6px;
}
.korisnici-header { display:flex; align-items:center; justify-content: space-between; }
.korisnik-item {
  background: #e0e7ff;
  border-radius: 12px;
  padding: 6px 12px;
  color: #3730a3;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}
.korisnik-item .x-btn { background: transparent; border: none; color: #7f1d1d; cursor: pointer; font-size: 14px; }
.verzije-panel {
  margin-top: 8px;
  background: #f9fafb;
  border-radius: 8px;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.file-info {
  display: flex;
  align-items: center;
  gap: 12px;
}
.file-name {
  font-weight: 600;
  color: #2563eb;
}
.version-meta {
  font-size: 0.95em;
  color: #374151;
  margin-top: 4px;
}
.side-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.btn[disabled] {
  opacity: 0.6;
  cursor: not-allowed;
}
/* Simple modal styles */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.35); display:flex; align-items:center; justify-content:center; z-index: 50; }
.modal { background: #fff; border-radius: 8px; width: 520px; max-width: 90vw; box-shadow: 0 10px 30px rgba(0,0,0,0.25); }
.modal-header { display:flex; align-items:center; justify-content: space-between; padding: 12px 16px; border-bottom: 1px solid #e5e7eb; }
.modal-body { padding: 16px; display: flex; flex-direction: column; gap: 12px; }
.modal-footer { padding: 12px 16px; border-top: 1px solid #e5e7eb; display:flex; gap: 10px; justify-content: flex-end; }
.modal-close { background: transparent; border: none; font-size: 20px; cursor: pointer; }
.form-group { display:flex; flex-direction: column; gap: 6px; }
.form-group input, .form-group select { padding: 8px 10px; border: 1px solid #cbd5e1; border-radius: 6px; }
.hint { color: #64748b; font-size: 12px; }

/* Naziv dokumenta iznad opisa */
.naziv-dokumenta {
  margin-bottom: 0.5em;
  color: #1e293b;
  font-size: 1.6em;
  font-weight: 700;
}

.phase-requests-box {
  background: #f1f5f9;
  border-radius: 8px;
  padding: 12px 10px 10px 10px;
  margin: 16px 0 8px 0;
}
.phase-requests-title {
  font-weight: 600;
  margin-bottom: 6px;
  color: #1e293b;
}
.phase-request-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e5e7eb;
}
.phase-request-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.phase-request-actions {
  display: flex;
  gap: 6px;
  margin-top: 4px;
}
.btn-approve, .btn-reject {
  padding: 4px 10px;
  font-size: 0.85em;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}
.btn-approve {
  background: #16a34a;
  color: white;
}
.btn-approve:hover {
  background: #15803d;
}
.btn-reject {
  background: #dc2626;
  color: white;
}
.btn-reject:hover {
  background: #b91c1c;
}
.phase-request-status {
  font-size: 0.95em;
  font-weight: 500;
  color: #2563eb;
}
.phase-request-status.status-odobren {
  color: #16a34a;
}
.phase-request-status.status-odbijen {
  color: #dc2626;
}
.phase-request-status.status-na,
.phase-request-status.status-na-cekanju,
.phase-request-status.status-na-čekanju {
  color: #eab308;
}
.phase-request-text {
  font-size: 0.97em;
  color: #374151;
}
.phase-request-date {
  font-size: 0.92em;
  color: #64748b;
}
.zadacici-section { margin-top: 1.5em; }
.zadacici-label { font-weight: bold; }
.zadacici-list { list-style: none; padding: 0; }
.zadacici-list li { display: flex; align-items: center; gap: 0.5em; margin-bottom: 0.3em; }
.zadacic-izvrsen span { text-decoration: line-through; color: #888; }
.zadacic-delete { background: none; border: none; color: #dc2626; cursor: pointer; font-size: 1.1em; }
.zadacic-add-row { display: flex; gap: 0.5em; margin-top: 0.5em; }
</style>