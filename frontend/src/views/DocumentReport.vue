<template>
  <Layout>
    <div class="page-container">
      <!-- Page Header -->
      <div class="page-header">
        <div>
          <h2>📊 Kompleksan Izveštaj - Korisnici i Dokumenti</h2>
          <div class="breadcrumb">Početna > Izveštaji > Dokumenti</div>
        </div>
        <button class="btn-primary" @click="loadReport" :disabled="loading">
          <span v-if="loading">⏳ Učitavanje...</span>
          <span v-else>🔄 Osveži</span>
        </button>
      </div>

      <!-- Error Message -->
      <div v-if="error" class="error-banner">
        <span>⚠️ {{ error }}</span>
        <button @click="error = null" class="close-btn">✕</button>
      </div>

      <!-- Report Content -->
      <div v-if="!loading && reportData.length > 0" class="report-container">
        <!-- Summary Stats -->
        <div class="summary-stats">
          <div class="stat-box">
            <div class="stat-value">{{ reportData.length }}</div>
            <div class="stat-label">Aktivnih korisnika</div>
          </div>
          <div class="stat-box">
            <div class="stat-value">{{ totalDocuments }}</div>
            <div class="stat-label">Ukupno dokumenata</div>
          </div>
          <div class="stat-box">
            <div class="stat-value">{{ totalVersions }}</div>
            <div class="stat-label">Ukupno verzija</div>
          </div>
          <div class="stat-box">
            <div class="stat-value">{{ totalSize.toFixed(2) }} MB</div>
            <div class="stat-label">Ukupna zauzetost</div>
          </div>
          <div class="stat-box">
            <div class="stat-value">{{ avgDocuments.toFixed(2) }}</div>
            <div class="stat-label">Prosečno dok/korisnik</div>
          </div>
          <div class="stat-box">
            <div class="stat-value">{{ avgSize.toFixed(2) }} MB</div>
            <div class="stat-label">Prosečna zauzetost</div>
          </div>
        </div>

        <!-- Detailed Table -->
        <div class="table-container">
          <table class="report-table">
            <thead>
              <tr>
                <th>#</th>
                <th>Korisnik</th>
                <th>Username</th>
                <th>Dokumenta</th>
                <th>Verzije</th>
                <th>Ukupna veličina</th>
                <th>Prosečna veličina</th>
                <th>Projekti</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(user, index) in reportData" :key="user.korisnik_id">
                <td>{{ index + 1 }}</td>
                <td class="user-name">
                  <div class="user-icon">👤</div>
                  <span>{{ user.puno_ime }}</span>
                </td>
                <td class="username">{{ user.korisnicko_ime }}</td>
                <td class="number">{{ user.broj_kreiranih_dokumenata }}</td>
                <td class="number">{{ user.broj_postavljenih_verzija }}</td>
                <td class="size">{{ user.ukupna_velicina_mb.toFixed(2) }} MB</td>
                <td class="size">{{ user.prosecna_velicina_mb.toFixed(2) }} MB</td>
                <td class="number">{{ user.broj_projekata }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Loading State -->
      <div v-else-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Učitavam izveštaj...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="!loading && reportData.length === 0" class="empty-state">
        <div class="empty-icon">📊</div>
        <h3>Nema podataka</h3>
        <p>Trenutno nema korisnika sa dokumentima u sistemu.</p>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { GetDocumentReportResults } from '../../wailsjs/go/main/App'
import Layout from '../components/Layout.vue'

const loading = ref(true)
const error = ref(null)
const reportData = ref([])

// Computed properties for summary statistics
const totalDocuments = computed(() => {
  return reportData.value.reduce((sum, user) => sum + user.broj_kreiranih_dokumenata, 0)
})

const totalVersions = computed(() => {
  return reportData.value.reduce((sum, user) => sum + user.broj_postavljenih_verzija, 0)
})

const totalSize = computed(() => {
  return reportData.value.reduce((sum, user) => sum + user.ukupna_velicina_mb, 0)
})

const avgDocuments = computed(() => {
  if (reportData.value.length === 0) return 0
  return totalDocuments.value / reportData.value.length
})

const avgSize = computed(() => {
  if (reportData.value.length === 0) return 0
  return totalSize.value / reportData.value.length
})

async function loadReport() {
  loading.value = true
  error.value = null
  
  try {
    const data = await GetDocumentReportResults()
    reportData.value = data || []
  } catch (err) {
    console.error('Error loading document report:', err)
    error.value = 'Greška pri učitavanju izveštaja: ' + err
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadReport()
})
</script>

<style scoped>
.page-container {
  padding: 30px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-header h2 {
  font-size: 28px;
  color: #2c3e50;
  margin-bottom: 5px;
}

.breadcrumb {
  color: #7f8c8d;
  font-size: 14px;
}

.btn-primary {
  background: #3498db;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.3s ease;
}

.btn-primary:hover:not(:disabled) {
  background: #2980b9;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-banner {
  background: #e74c3c;
  color: white;
  padding: 15px 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.close-btn {
  background: transparent;
  border: none;
  color: white;
  font-size: 20px;
  cursor: pointer;
  padding: 0 5px;
}

.report-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.summary-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 20px;
  padding: 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-box {
  text-align: center;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

.table-container {
  padding: 30px;
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
}

.report-table thead {
  background: #f8f9fa;
}

.report-table th {
  padding: 15px;
  text-align: left;
  font-weight: 600;
  color: #2c3e50;
  border-bottom: 2px solid #e1e8ed;
}

.report-table td {
  padding: 15px;
  border-bottom: 1px solid #e1e8ed;
}

.report-table tbody tr:hover {
  background: #f8f9fa;
}

.user-name {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 500;
}

.user-icon {
  width: 36px;
  height: 36px;
  background: #3498db;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.username {
  color: #7f8c8d;
  font-family: monospace;
}

.number {
  text-align: center;
  font-weight: 500;
  color: #2c3e50;
}

.size {
  text-align: right;
  font-family: monospace;
  color: #e74c3c;
  font-weight: 500;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #e1e8ed;
  border-top-color: #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 24px;
  color: #2c3e50;
  margin-bottom: 10px;
}

.empty-state p {
  color: #7f8c8d;
  font-size: 16px;
}
</style>
