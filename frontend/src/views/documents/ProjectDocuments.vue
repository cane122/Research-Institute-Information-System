<template>
  <Layout>
    <div class="project-documents">
      <div class="page-header">
        <div>
          <h2>Dokumentacija projekta</h2>
          <div class="breadcrumb">Početna > Projekti > Dokumentacija</div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" @click="goBack">⟵ Nazad na projekte</button>
          <button v-if="canAdd" class="btn btn-primary" @click="goToUpload">📤 Dodaj dokument</button>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <h3>Dokumenti za projekat #{{ projectId }}</h3>
        </div>
        <div class="card-body">
          <div v-if="loading" class="loading">Učitavanje dokumenata...</div>
          <div v-else-if="error" class="error">{{ error }}</div>
          <div v-else>
            <div v-if="docs.length === 0" class="empty">Nema dokumenata za ovaj projekat.</div>
            <div v-else class="docs-list">
              <div class="docs-header">
                <div>Naziv</div>
                <div>Autor</div>
                <div>Tip</div>
                <div>Poslednja izmena</div>
                <div style="min-width:120px">Progres</div>
                <div>Akcije</div>
              </div>
              <div v-for="d in docs" :key="d.dokument_id" class="docs-row">
                <div class="name">📄 {{ d.naziv_dokumenta }}</div>
                <div>{{ d.ime_kreirao }}</div>
                <div>{{ d.tip_dokumenta || 'Dokument' }}</div>
                <div>{{ formatDate(d.poslednja_izmena || d.datuma_postavke) }}</div>
                <div>
                  <div v-if="d.progres !== undefined && d.progres !== null" class="progress-bar-outer">
                    <div class="progress-bar-inner" :style="{ width: d.progres + '%' }"></div>
                    <span class="progress-label">{{ d.progres }}%</span>
                  </div>
                  <span v-else style="color:#aaa">N/A</span>
                </div>
                <div class="actions">
                  <button class="btn btn-small" @click="fullPreview(d)">Pregled</button>
                </div>
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
const projectId = Number(route.params.id)

const docs = ref([])
const loading = ref(true)
const error = ref('')
const canAdd = ref(false)

function formatDate(s) {
  if (!s) return '-'
  const d = new Date(s)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('sr-RS')
}

function goBack() {
  router.push('/projects')
}

function goToUpload() {
  router.push(`/documents/create?projectId=${projectId}`)
}

function fullPreview(doc) {
  router.push(`/documents/full-preview/${doc.dokument_id}`)
}

onMounted(async () => {
  try {
    loading.value = true
    const res = await window.go?.main?.App?.GetProjectDocuments(projectId)
    docs.value = res || []
    // ask backend if current user can add a document to this project
    if (window.go?.main?.App?.CanAddProjectDocument) {
      try {
        canAdd.value = await window.go.main.App.CanAddProjectDocument(projectId)
      } catch (e) {
        canAdd.value = false
      }
    }
  } catch (e) {
    error.value = e?.message || 'Greška pri učitavanju dokumenata'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.card { background: #fff; border: 1px solid #e5e7eb; border-radius: 8px; }
.card-header { padding: 12px 16px; border-bottom: 1px solid #e5e7eb; }
.card-body { padding: 16px; }
.docs-header, .docs-row { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr 1.2fr 1fr; gap: 12px; align-items: center; }
.progress-bar-outer {
  position: relative;
  background: #f3f4f6;
  border-radius: 6px;
  height: 18px;
  width: 100%;
  min-width: 90px;
  margin-right: 4px;
  overflow: hidden;
  display: flex;
  align-items: center;
}
.progress-bar-inner {
  background: linear-gradient(90deg, #3b82f6 0%, #06b6d4 100%);
  height: 100%;
  transition: width 0.3s;
}
.progress-label {
  position: absolute;
  left: 50%;
  top: 0;
  transform: translateX(-50%);
  font-size: 12px;
  color: #222;
  font-weight: 600;
  z-index: 2;
}
.docs-header { font-weight: 600; color: #374151; padding: 8px 0; border-bottom: 1px solid #e5e7eb; }
.docs-row { padding: 10px 0; border-bottom: 1px solid #f3f4f6; }
.name { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.empty, .loading, .error { padding: 16px; color: #6b7280; }
.error { color: #b91c1c; }
.btn-small { font-size: 12px; padding: 6px 10px; }
</style>