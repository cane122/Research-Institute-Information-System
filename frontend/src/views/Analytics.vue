<template>
  <Layout>
    <div class="analytics">
      <!-- Header -->
      <div class="page-header">
        <div>
          <h2>Analitika Projekta {{ projectName }}</h2>
          <div class="breadcrumb">Početna > Projekti > Analitika</div>
        </div>
        <div class="header-actions">
          <button 
            class="btn btn-secondary"
            :class="{ active: statusFilter === 'active' }"
            @click="setStatusFilter('active')"
          >
            Aktivan
          </button>
          <button class="btn btn-secondary" @click="showFilters = !showFilters">
            Filteri
          </button>
          <button class="btn btn-secondary" @click="goBack">
            ← Nazad
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Učitavanje analitike...</p>
      </div>

      <!-- Analytics Content -->
      <div v-else class="analytics-content">
        <!-- Top Statistics Cards -->
        <div class="stats-grid">
          <div class="stat-card" style="background: linear-gradient(135deg, #a8d8ff 0%, #87ceeb 100%);">
            <div class="stat-label">Ukupno zadataka</div>
            <div class="stat-value">{{ analytics.totalTasks || 0 }}</div>
          </div>
          
          <div class="stat-card" style="background: linear-gradient(135deg, #90ee90 0%, #66bb6a 100%);">
            <div class="stat-label">Aktivni</div>
            <div class="stat-value">{{ analytics.activeTasks || 0 }}</div>
          </div>
          
          <div class="stat-card" style="background: linear-gradient(135deg, #ffb3ba 0%, #ff8a80 100%);">
            <div class="stat-label">Kasne</div>
            <div class="stat-value">{{ analytics.overdueTasks || 0 }}</div>
          </div>
          
          <div class="stat-card" style="background: linear-gradient(135deg, #b0c4de 0%, #87ceeb 100%);">
            <div class="stat-label">Završeni</div>
            <div class="stat-value">{{ analytics.completedTasks || 0 }}</div>
          </div>
        </div>

        <!-- Middle Section: Progress & Daily Average -->
        <div class="middle-section">
          <div class="progress-card">
            <div class="card-label">Stepen realizacije</div>
            <div class="progress-value">{{ analytics.completionRate || 0 }}%</div>
          </div>
          
          <div class="daily-card">
            <div class="card-label">Prosečno dana po fazi</div>
            <div class="daily-value">{{ analytics.averageDaysPerPhase || 0 }}</div>
          </div>
        </div>

        <!-- Charts Section -->
        <div class="charts-section">
          <!-- Tasks by Phase -->
          <div class="chart-container">
            <h3>Zadaci po fazama</h3>
            <canvas ref="tasksByPhaseChart"></canvas>
          </div>

          <!-- Project Progress Over Time -->
          <div class="chart-container">
            <h3>Napredak projekta</h3>
            <canvas ref="progressChart"></canvas>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Layout from '../components/Layout.vue'
import { GetProjectAnalytics } from '../../wailsjs/go/main/App.js'
import Chart from 'chart.js/auto'

const router = useRouter()
const route = useRoute()

const projectId = computed(() => parseInt(route.params.id))
const projectName = ref('')
const statusFilter = ref('active')
const showFilters = ref(false)
const loading = ref(false)
const analytics = ref({
  totalTasks: 0,
  activeTasks: 0,
  overdueTasks: 0,
  completedTasks: 0,
  completionRate: 0,
  averageDaysPerPhase: 0,
  tasksByPhase: [],
  progressOverTime: []
})

const tasksByPhaseChart = ref(null)
const progressChart = ref(null)
let tasksByPhaseChartInstance = null
let progressChartInstance = null

function setStatusFilter(status) {
  statusFilter.value = status
  loadAnalytics()
}

function goBack() {
  router.push('/projects')
}

async function loadAnalytics() {
  try {
    loading.value = true
    
    const data = await GetProjectAnalytics(projectId.value)
    
    console.log('Analytics data:', data)
    
    // Map backend data to frontend structure
    const totalTasks = data.total_tasks || 0
    const completedTasks = data.completed_tasks || 0
    const inProgressTasks = data.in_progress_tasks || 0
    const pendingTasks = data.pending_tasks || 0
    const overdueTasks = data.overdue_tasks || 0
    
    // Calculate active tasks (in progress + pending)
    const activeTasks = inProgressTasks + pendingTasks
    
    // Calculate completion rate
    const completionRate = data.completion_percentage || 0
    
    // Map tasks by phase
    const tasksByPhase = (data.tasks_by_phase || []).map(phase => ({
      phaseName: phase.phase_name,
      count: phase.task_count
    }))
    
    // Calculate average days per phase (mock for now)
    const averageDaysPerPhase = totalTasks > 0 ? Math.round(30 + Math.random() * 20) : 0
    
    // Generate mock progress over time data
    const progressOverTime = []
    const today = new Date()
    for (let i = 5; i >= 0; i--) {
      const date = new Date(today)
      date.setDate(date.getDate() - i * 7)
      const formattedDate = date.toLocaleDateString('sr-RS', { month: 'short', day: 'numeric' })
      const progress = Math.min(completionRate, (completionRate / 6) * (6 - i))
      progressOverTime.push({
        date: formattedDate,
        completionRate: Math.round(progress)
      })
    }
    
    // Update analytics data
    analytics.value = {
      totalTasks,
      activeTasks,
      overdueTasks,
      completedTasks,
      completionRate: Math.round(completionRate),
      averageDaysPerPhase,
      tasksByPhase,
      progressOverTime
    }
    
    projectName.value = data.project?.naziv_projekta || `Projekat #${projectId.value}`
    
    // Wait for DOM to update before rendering charts
    await nextTick()
    renderCharts()
  } catch (err) {
    console.error('Error loading analytics:', err)
    alert('Greška pri učitavanju analitike: ' + (err.message || err))
  } finally {
    loading.value = false
  }
}

function renderCharts() {
  renderTasksByPhaseChart()
  renderProgressChart()
}

function renderTasksByPhaseChart() {
  if (!tasksByPhaseChart.value) return
  
  // Destroy existing chart
  if (tasksByPhaseChartInstance) {
    tasksByPhaseChartInstance.destroy()
  }
  
  const ctx = tasksByPhaseChart.value.getContext('2d')
  
  // Prepare data
  const phases = analytics.value.tasksByPhase || []
  const labels = phases.map(p => p.phaseName || 'Nepoznato')
  const data = phases.map(p => p.count || 0)
  
  tasksByPhaseChartInstance = new Chart(ctx, {
    type: 'bar',
    data: {
      labels: labels,
      datasets: [{
        label: 'Broj zadataka',
        data: data,
        backgroundColor: [
          'rgba(168, 216, 255, 0.8)',
          'rgba(144, 238, 144, 0.8)',
          'rgba(255, 179, 186, 0.8)',
          'rgba(176, 196, 222, 0.8)',
          'rgba(255, 218, 185, 0.8)',
          'rgba(221, 160, 221, 0.8)'
        ],
        borderColor: [
          'rgba(135, 206, 235, 1)',
          'rgba(102, 187, 106, 1)',
          'rgba(255, 138, 128, 1)',
          'rgba(135, 206, 235, 1)',
          'rgba(255, 200, 150, 1)',
          'rgba(186, 85, 211, 1)'
        ],
        borderWidth: 2
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            stepSize: 1
          }
        }
      }
    }
  })
}

function renderProgressChart() {
  if (!progressChart.value) return
  
  // Destroy existing chart
  if (progressChartInstance) {
    progressChartInstance.destroy()
  }
  
  const ctx = progressChart.value.getContext('2d')
  
  // Prepare data
  const progress = analytics.value.progressOverTime || []
  const labels = progress.map(p => p.date || '')
  const data = progress.map(p => p.completionRate || 0)
  
  progressChartInstance = new Chart(ctx, {
    type: 'line',
    data: {
      labels: labels,
      datasets: [{
        label: 'Procenat završenosti',
        data: data,
        borderColor: 'rgba(135, 206, 235, 1)',
        backgroundColor: 'rgba(168, 216, 255, 0.2)',
        borderWidth: 3,
        fill: true,
        tension: 0.4
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          max: 100,
          ticks: {
            callback: function(value) {
              return value + '%'
            }
          }
        }
      }
    }
  })
}

onMounted(() => {
  loadAnalytics()
})
</script>

<style scoped>
.analytics {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.8rem;
  color: #1a1a1a;
}

.breadcrumb {
  color: #666;
  font-size: 0.9rem;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary {
  background: #f5f5f5;
  color: #333;
}

.btn-secondary:hover {
  background: #e0e0e0;
}

.btn-secondary.active {
  background: #90ee90;
  color: white;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #0066cc;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.analytics-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}

.stat-card {
  padding: 2rem 1.5rem;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 2px solid rgba(0, 0, 0, 0.1);
}

.stat-label {
  font-size: 0.95rem;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 0.75rem;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1a1a1a;
}

.middle-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.progress-card,
.daily-card {
  padding: 2rem 1.5rem;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 2px solid rgba(0, 0, 0, 0.1);
}

.progress-card {
  background: linear-gradient(135deg, #a8d8ff 0%, #87ceeb 100%);
}

.daily-card {
  background: linear-gradient(135deg, #dda0dd 0%, #ba55d3 100%);
}

.card-label {
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 0.75rem;
}

.progress-value,
.daily-value {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1a1a1a;
}

.charts-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

.chart-container {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 2px solid rgba(0, 0, 0, 0.05);
}

.chart-container h3 {
  margin: 0 0 1.5rem 0;
  font-size: 1.1rem;
  color: #1a1a1a;
  text-align: center;
}

.chart-container canvas {
  max-height: 300px;
}

/* Responsive */
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .charts-section {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .analytics {
    padding: 1rem;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .middle-section {
    grid-template-columns: 1fr;
  }
}
</style>
