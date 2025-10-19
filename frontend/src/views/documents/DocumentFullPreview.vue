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
          <div class="opis-text">{{ document?.opis || 'Nema opisa.' }}</div>
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

        <!-- Lista korisnika -->
        <div class="section korisnici">
          <label class="korisnici-label">Korisnici na dokumentu:</label>
          <div class="korisnici-list">
            <div v-for="korisnik in korisnici" :key="korisnik.korisnik_id" class="korisnik-item">
              {{ korisnik.ime || korisnik.korisnicko_ime }} {{ korisnik.prezime || '' }}
            </div>
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
              <button class="btn btn-secondary" disabled>Preuzmi</button>
              <button class="btn btn-secondary" disabled>Izmeni</button>
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
          <button class="btn btn-primary" disabled>Zahtevaj promenu faze</button>
          <div><b>Rok za završetak:</b> {{ formatDate(document?.rok) }}</div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import Layout from '../../components/Layout.vue'
import { GetDocumentByID, GetDocumentVersions, GetWorkflowPhases, GetProjectMembers, GetWorkflowByID } from '../../../wailsjs/go/main/App.js'

const route = useRoute()
const documentId = Number(route.params.id)

const document = ref(null)
const verzije = ref([])
const selectedVerzijaId = ref(null)
const selectedVersion = computed(() => verzije.value.find(v => v.verzija_id === selectedVerzijaId.value) || verzije.value[0] || null)
const faze = ref([])
const korisnici = ref([])
const workflow = ref(null)
const currentFaza = computed(() => faze.value.find(f => f.faza_id === document.value?.trenutna_faza_id) || null)

const selectedVersionUser = computed(() => {
  // Find user by postavio_korisnik_id
  if (!selectedVersion.value) return '-'
  const user = korisnici.value.find(u => u.korisnik_id === selectedVersion.value.postavio_korisnik_id)
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
function onVersionChange() {
  // No-op, computed will update selectedVersion
}

onMounted(async () => {
  // Load document
  document.value = await GetDocumentByID(documentId)
  // Load versions
  verzije.value = await GetDocumentVersions(documentId)
  selectedVerzijaId.value = verzije.value[0]?.verzija_id || null
  // Load workflow phases
  if (document.value?.radni_tok_id) {
    faze.value = await GetWorkflowPhases(document.value.radni_tok_id)
    workflow.value = await GetWorkflowByID(document.value.radni_tok_id)
  }
  // Load project members (users on document)
  if (document.value?.projekat_id) {
    korisnici.value = await GetProjectMembers(document.value.projekat_id)
  } else {
    korisnici.value = []
  }
})
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
.korisnik-item {
  background: #e0e7ff;
  border-radius: 12px;
  padding: 6px 12px;
  color: #3730a3;
  font-weight: 500;
}
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
</style>
/* Naziv dokumenta iznad opisa */
.naziv-dokumenta {
  margin-bottom: 0.5em;
  color: #1e293b;
  font-size: 1.6em;
  font-weight: 700;
}
