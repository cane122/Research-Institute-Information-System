<template>
  <Layout>
    <div class="document-analytics">
      <div class="page-header">
        <div>
          <h2>Analitika dokumenta</h2>
          <div class="breadcrumb">Početna > Projekti > Dokumenti > Analitika</div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" @click="goBack">⟵ Nazad</button>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <h3>Dokument #{{ docId }} — događaji</h3>
        </div>
        <div class="card-body">
          <div v-if="loading" class="loading">Učitavanje...</div>
          <div v-else-if="error" class="error">{{ error }}</div>
          <div v-else>
            <div v-if="events.length === 0" class="empty">Nema zabeleženih događaja za ovaj dokument.</div>
            <div v-else class="table-wrapper">
              <div class="table">
                <div class="thead">
                  <div>Datum i vreme</div>
                  <div>Tip događaja</div>
                  <div>Korisnik</div>
                  <div>Opis</div>
                </div>
                <div v-for="e in events" :key="e.log_id" class="trow">
                  <div>{{ formatDateTime(e.kreiran_datuma) }}</div>
                  <div>
                    <span :class="['tag', e.tip_aktivnosti]">{{ formatType(e.tip_aktivnosti) }}</span>
                  </div>
                  <div>{{ e.korisnik_ime || '—' }}</div>
                  <div>{{ e.opis || e.naziv_entiteta || '-' }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Faze trajanja -->
      <div class="card" style="margin-top: 24px;">
        <div class="card-header">
          <h3>Vreme provedeno u svakoj fazi</h3>
        </div>
        <div class="card-body">
          <div v-if="loadingPhases" class="loading">Učitavanje...</div>
          <div v-else-if="errorPhases" class="error">{{ errorPhases }}</div>
          <div v-else>
            <div v-if="!phaseDurations.length" class="empty">Nema podataka o fazama za ovaj dokument.</div>
            <div v-else class="phase-table">
              <div class="phase-thead">
                <div>Faza</div>
                <div>Trajanje</div>
              </div>
              <div v-for="(p, idx) in phaseDurations" :key="idx" class="phase-trow">
                <div>{{ p.naziv_faze || p.NazivFaze || '-' }}</div>
                <div>{{ formatDuration(p.trajanje_sekundi || p.TrajanjeSekundi) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>


<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Layout from '../../components/Layout.vue'

const route = useRoute()
const router = useRouter()
const projectId = Number(route.params.projectId)
const docId = Number(route.params.docId)

const events = ref([])
const loading = ref(true)
const error = ref('')

const phaseDurations = ref([])
const loadingPhases = ref(true)
const errorPhases = ref('')

function goBack() {
  // Go back to project documents list
  if (projectId) router.push(`/projects/${projectId}/documents`)
  else router.back()
}

function formatDateTime(s) {
  if (!s) return '-'
  const d = new Date(s)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('sr-RS')
}

function formatType(t) {
  if (!t) return '-'
  const key = String(t).toUpperCase()
  if (key === 'DOCUMENT_EDIT') return 'Izmena dokumenta'
  if (key === 'PHASE_CHANGE') return 'Promena faze'
  return key
}

onMounted(async () => {
  try {
    loading.value = true
    const res = await window.go?.main?.App?.GetDocumentActivity(docId)
    events.value = Array.isArray(res) ? res : []
  } catch (e) {
    error.value = e?.message || 'Greška pri učitavanju analitike'
  } finally {
    loading.value = false
  }
  // Fetch phase durations
  try {
    loadingPhases.value = true
    const res = await window.go?.main?.App?.GetDocumentPhaseDurations(docId)
    console.log('Phase durations response:', res)
    phaseDurations.value = Array.isArray(res) ? res : []
    console.log('Phase durations array:', phaseDurations.value)
  } catch (e) {
    console.error('Error loading phase durations:', e)
    errorPhases.value = e?.message || 'Greška pri učitavanju faza'
  } finally {
    loadingPhases.value = false
  }
})
function formatDuration(secs) {
  if (!secs || isNaN(secs)) return '-'
  const d = Math.floor(secs / 86400)
  const h = Math.floor((secs % 86400) / 3600)
  const m = Math.floor((secs % 3600) / 60)
  const s = secs % 60
  let out = []
  if (d) out.push(`${d}d`)
  if (h) out.push(`${h}h`)
  if (m) out.push(`${m}m`)
  if (s || out.length === 0) out.push(`${s}s`)
  return out.join(' ')
}
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.card { background: #fff; border: 1px solid #e5e7eb; border-radius: 8px; }
.card-header { padding: 12px 16px; border-bottom: 1px solid #e5e7eb; }
.card-body { padding: 16px; }

/* Table wrapper with max height and scrollbar for events */
.table-wrapper { max-height: 400px; overflow-y: auto; }

.table { display: grid; gap: 8px; }
.thead, .trow { display: grid; grid-template-columns: 1.2fr 1fr 1fr 2fr; gap: 12px; align-items: center; }
.thead { font-weight: 600; color: #374151; padding: 8px 0; border-bottom: 1px solid #e5e7eb; position: sticky; top: 0; background: #fff; z-index: 1; }
.trow { padding: 8px 0; border-bottom: 1px solid #f3f4f6; }
.empty, .loading, .error { padding: 12px; color: #6b7280; }
.error { color: #b91c1c; }
.tag { display:inline-block; padding: 2px 8px; border-radius: 999px; font-size: 12px; font-weight: 600; background:#eef2ff; color:#3730a3; }
.tag.PHASE_CHANGE { background:#ecfeff; color:#0e7490; }
.tag.DOCUMENT_EDIT { background:#eff6ff; color:#1d4ed8; }
.card + .card { margin-top: 32px; }

/* Phase duration table with 2 columns */
.phase-table { display: grid; gap: 8px; }
.phase-thead, .phase-trow { display: grid; grid-template-columns: 2fr 1fr; gap: 12px; align-items: center; }
.phase-thead { font-weight: 600; color: #374151; padding: 8px 0; border-bottom: 1px solid #e5e7eb; }
.phase-trow { padding: 8px 0; border-bottom: 1px solid #f3f4f6; }
</style>