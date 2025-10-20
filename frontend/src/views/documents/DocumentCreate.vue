<template>
  <Layout>
    <div class="doc-create-page">
      <div class="doc-create-main">
        <h2>Kreiraj novi dokument</h2>
        <div class="form-section">
          <label>Fajl dokumenta</label>
          <input type="file" @change="onFileChange" />
        </div>
        <div class="form-section">
          <label>Naziv dokumenta</label>
          <input v-model="naziv" type="text" placeholder="Unesite naziv dokumenta" />
        </div>
        <div class="form-section">
          <label>Radni tok</label>
          <select v-model="radniTokId" @change="loadFaze">
            <option v-for="tok in radniTokovi" :key="tok.radni_tok_id" :value="tok.radni_tok_id">
              {{ tok.naziv }}
            </option>
          </select>
          <button class="btn btn-secondary ml-2" @click="goToWorkflowCreate">Dodaj novi</button>
        </div>
        <div class="form-section">
          <label>Faze izabranog radnog toka</label>
          <ul class="faze-list">
            <li v-for="faza in faze" :key="faza.faza_id">{{ faza.naziv_faze }}</li>
          </ul>
        </div>
        <div class="form-section">
          <label>Korisnici sa dozvolama</label>
          <div class="user-picker">
            <select v-model="selectedUserId">
              <option v-for="user in availableUsers" :key="user.korisnik_id" :value="user.korisnik_id">
                {{ user.ime }} {{ user.prezime }} ({{ user.korisnicko_ime }})
              </option>
            </select>
            <button class="btn btn-small" @click="addUser">Dodaj</button>
            <button v-if="projectMembers.length > 0" class="btn btn-small btn-project" @click="addAllProjectMembers">Dodaj sve članove projekta</button>
          </div>
          <ul class="user-list">
            <li v-for="u in docUsers" :key="u.korisnik_id">
              {{ u.ime }} {{ u.prezime }} ({{ u.korisnicko_ime }})
              <span class="perm-label">[Čitanje, Izmena, Brisanje]</span>
              <button class="btn btn-small" @click="removeUser(u.korisnik_id)">Ukloni</button>
            </li>
          </ul>
        </div>
        <div class="form-section">
          <label>Rok za završetak</label>
          <input type="date" v-model="rok" />
        </div>
        <div class="form-section">
          <label>Opis dokumenta</label>
          <textarea v-model="opis" rows="3" placeholder="Unesite opis dokumenta"></textarea>
        </div>
        <div class="form-section">
          <button class="btn btn-primary" @click="submitDoc">Dodaj dokument</button>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../../components/Layout.vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const projectId = ref(Number(route.query.projectId) || null)
const naziv = ref('')
const radniTokId = ref(null)
const radniTokovi = ref([])
const faze = ref([])
const users = ref([])
const projectMembers = ref([])
const availableUsers = ref([])
const selectedUserId = ref(null)
const docUsers = ref([])
const rok = ref('')
const opis = ref('')
const file = ref(null)

function onFileChange(e) {
  file.value = e.target.files[0] || null
}

function goToWorkflowCreate() {
  router.push('/workflows/create')
}

async function loadProjectMembers() {
  if (!projectId.value) {
    projectMembers.value = []
    availableUsers.value = users.value
    return
  }
  try {
    projectMembers.value = await window.go?.main?.App?.GetProjectMembers(projectId.value) || []
    // Show project members in the dropdown if available
    availableUsers.value = projectMembers.value.length > 0 ? projectMembers.value : users.value
  } catch (e) {
    console.error('Error loading project members:', e)
    availableUsers.value = users.value
  }
}

function addUser() {
  const user = availableUsers.value.find(u => u.korisnik_id === selectedUserId.value)
  if (user && !docUsers.value.some(u => u.korisnik_id === user.korisnik_id)) {
    docUsers.value.push(user)
  }
}

function addAllProjectMembers() {
  projectMembers.value.forEach(member => {
    if (!docUsers.value.some(u => u.korisnik_id === member.korisnik_id)) {
      docUsers.value.push(member)
    }
  })
}

function removeUser(id) {
  docUsers.value = docUsers.value.filter(u => u.korisnik_id !== id)
}

async function loadRadniTokovi() {
  // TODO: Replace with backend call
  const allTokovi = await window.go?.main?.App?.GetAllWorkflows() || []
  radniTokovi.value = allTokovi.filter(t => t.tip_toka === 'DOKUMENTACIJA')
  if (radniTokovi.value.length) {
    radniTokId.value = radniTokovi.value[0].radni_tok_id
    await loadFaze()
  }
}

async function loadFaze() {
  if (!radniTokId.value) return
  faze.value = await window.go?.main?.App?.GetWorkflowPhases(radniTokId.value) || []
}

async function loadUsers() {
  try {
    users.value = await window.go?.main?.App?.GetAllUsersForDocuments() || []
    availableUsers.value = users.value
  } catch (e) {
    console.error('Error loading users:', e)
    users.value = []
    availableUsers.value = []
  }
}

async function submitDoc() {
  // Validation
  if (!naziv.value.trim()) {
    alert('Naziv dokumenta je obavezan')
    return
  }
  if (!file.value) {
    alert('Morate izabrati fajl')
    return
  }
  if (!radniTokId.value) {
    alert('Morate izabrati radni tok')
    return
  }
  if (!projectId.value) {
    alert('Dokument mora biti vezan za projekat')
    return
  }

  try {
    // Read file as bytes
    const fileBytes = await readFileAsBytes(file.value)
    
    // Prepare document data
    const docData = {
      naziv_dokumenta: naziv.value.trim(),
      projekat_id: projectId.value,
      radni_tok_id: radniTokId.value,
      opis: opis.value.trim() || undefined,
      rok: rok.value || undefined,
      korisnici_dozvole: docUsers.value.map(u => u.korisnik_id)
    }

    console.log('Sending document data:', docData)
    console.log('File bytes length:', fileBytes.length)
    console.log('File name:', file.value.name)

    // Call backend
    const documentId = await window.go?.main?.App?.CreateDocumentWithPermissions(
      docData,
      fileBytes,
      file.value.name
    )

    alert('Dokument uspešno dodat!')
    // Navigate back to project documents
    router.push(`/projects/${projectId.value}/documents`)
  } catch (e) {
    console.error('Error adding document:', e)
    alert(e?.message || 'Greška pri dodavanju dokumenta')
  }
}

function readFileAsBytes(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const arrayBuffer = reader.result
      const bytes = new Uint8Array(arrayBuffer)
      resolve(Array.from(bytes))
    }
    reader.onerror = reject
    reader.readAsArrayBuffer(file)
  })
}

onMounted(async () => {
  await loadRadniTokovi()
  await loadUsers()
  // Load project members if we have a project context
  if (projectId.value) {
    await loadProjectMembers()
  }
})
</script>

<style scoped>
.doc-create-page {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: 100vh;
  background: #f8fafc;
}
.doc-create-main {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  padding: 32px 40px;
  margin-top: 32px;
  width: 520px;
}
.form-section {
  margin-bottom: 24px;
}
.form-section label {
  font-weight: 600;
  display: block;
  margin-bottom: 8px;
}
input[type="text"], input[type="date"], select, textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 15px;
  margin-bottom: 4px;
}
textarea {
  resize: vertical;
}
.faze-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.faze-list li {
  background: #f3f4f6;
  border-radius: 4px;
  padding: 6px 12px;
  margin-bottom: 4px;
  font-size: 14px;
}
.user-picker {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}
.user-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.user-list li {
  background: #f3f4f6;
  border-radius: 4px;
  padding: 6px 12px;
  margin-bottom: 4px;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.perm-label {
  color: #64748b;
  font-size: 12px;
  margin-left: 8px;
}
.btn {
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 15px;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-secondary {
  background: #64748b;
}
.btn-small {
  font-size: 13px;
  padding: 4px 10px;
  margin-left: 8px;
}
.btn-project {
  background: #059669;
}
.ml-2 {
  margin-left: 8px;
}
</style>
