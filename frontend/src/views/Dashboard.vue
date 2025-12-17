<template>
  <Layout>
    <div class="page-container">
      <!-- Dashboard Header -->
      <div class="page-header">
        <div>
          <h2>Dashboard</h2>
          <div class="breadcrumb">Početna > Dashboard</div>
        </div>
      </div>
      
      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card projects">
          <div class="stat-content">
            <h3>{{ stats.projects }}</h3>
            <p>Aktivnih projekata</p>
          </div>
          <div class="stat-icon">📁</div>
        </div>
        
        <div class="stat-card documents">
          <div class="stat-content">
            <h3>{{ userDocuments }}</h3>
            <p>Mojih dokumenata</p>
          </div>
          <div class="stat-icon">📄</div>
        </div>
        
        <div class="stat-card tasks">
          <div class="stat-content">
            <h3>{{ stats.tasks }}</h3>
            <p>Aktivnih zadataka</p>
          </div>
          <div class="stat-icon">📋</div>
        </div>
        
        <div class="stat-card storage">
          <div class="stat-content">
            <h3>{{ userStorage.toFixed(2) }} MB</h3>
            <p>Zauzetost dokumenata</p>
          </div>
          <div class="stat-icon">💾</div>
        </div>
      </div>
      
      <!-- Content Grid -->
      <div class="content-grid">
        <!-- Recent Activities -->
        <div class="recent-activities card">
          <div class="card-header">
            <h3>Poslednje aktivnosti</h3>
          </div>
          
          <div class="activities-list">
            <div 
              v-for="activity in recentActivities" 
              :key="activity.id"
              class="activity-item"
            >
              <div class="activity-content">
                <div class="activity-title">{{ activity.title }}</div>
                <div class="activity-desc">{{ activity.description }}</div>
              </div>
              <div class="activity-time">{{ activity.time }}</div>
            </div>
          </div>
          
          <div v-if="recentActivities.length === 0" class="no-activities">
            <p>Nema poslednih aktivnosti</p>
          </div>
        </div>
        
        <!-- Quick Actions -->
        <div class="quick-actions-card card">
          <div class="card-header">
            <h3>Brze akcije</h3>
          </div>
          
          <div class="quick-actions">
            <router-link to="/projects" class="quick-btn btn-primary">
              <span class="btn-icon">📁</span>
              Novi projekat
            </router-link>
            
            <router-link to="/tasks" class="quick-btn btn-success">
              <span class="btn-icon">📋</span>
              Dodaj zadatak
            </router-link>
            
            <router-link to="/documents" class="quick-btn btn-warning">
              <span class="btn-icon">📄</span>
              Upload dokument
            </router-link>
            
            <router-link to="/documents/report" class="quick-btn btn-info">
              <span class="btn-icon">📈</span>
              PL/SQL Izveštaj
            </router-link>
            
            <router-link v-if="authStore.isAdmin" to="/users" class="quick-btn btn-secondary">
              <span class="btn-icon">👥</span>
              Upravljanje korisnicima
            </router-link>
          </div>
        </div>
      </div>
      
      <!-- Project Progress -->
      <div class="project-progress card">
        <div class="card-header">
          <h3>Status projekata</h3>
        </div>
        
        <div class="projects-list">
          <div 
            v-for="project in projectsProgress" 
            :key="project.id"
            class="project-item"
          >
            <div class="project-info">
              <div class="project-name">{{ project.name }}</div>
              <div class="project-desc">{{ project.description }}</div>
            </div>
            
            <div class="project-progress-bar">
              <div class="progress-container">
                <div 
                  class="progress-bar" 
                  :style="{ width: project.progress + '%' }"
                ></div>
              </div>
              <span class="progress-text">{{ project.progress }}%</span>
            </div>
            
            <div class="project-status">
              <span 
                class="status-badge" 
                :class="`status-${project.status.toLowerCase()}`"
              >
                {{ project.status }}
              </span>
            </div>
          </div>
        </div>
        
        <div v-if="projectsProgress.length === 0" class="no-projects">
          <p>Nema aktivnih projekata</p>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/auth'
import { GetUserDocumentSize, GetUserDocumentCount } from '../../wailsjs/go/main/App'

const authStore = useAuthStore()

// Reactive data
const stats = ref({
  projects: 12,
  documents: 158,
  tasks: 23,
  users: 8
})

const userStorage = ref(0)
const userDocuments = ref(0)

const recentActivities = ref([
  {
    id: 1,
    title: 'Kreiran novi projekat "AI Analiza"',
    description: 'Marko Petrović je kreirao novi projekat',
    time: 'Pre 2 sata'
  },
  {
    id: 2,
    title: 'Završen zadatak "Dokumentacija"',
    description: 'Ana Jovanović je završila zadatak',
    time: 'Pre 4 sata'
  },
  {
    id: 3,
    title: 'Dodato je 5 novih dokumenata',
    description: 'Stefan Nikolić je dodao dokumente',
    time: 'Juče'
  },
  {
    id: 4,
    title: 'Ažuriran projekat "Web Portal"',
    description: 'Milica Stojković je ažurirala status',
    time: 'Juče'
  }
])

const projectsProgress = ref([
  {
    id: 1,
    name: 'Web Portal Refactoring',
    description: 'Modernizacija postojećeg web portala',
    progress: 75,
    status: 'Active'
  },
  {
    id: 2,
    name: 'AI Analitika Modula',
    description: 'Implementacija AI za analizu podataka',
    progress: 45,
    status: 'Active'
  },
  {
    id: 3,
    name: 'Mobile Aplikacija',
    description: 'Razvoj mobilne aplikacije za iOS/Android',
    progress: 20,
    status: 'Active'
  },
  {
    id: 4,
    name: 'Sigurnosni Audit',
    description: 'Kompletan sigurnosni pregled sistema',
    progress: 90,
    status: 'Pending'
  }
async function loadDashboardData() {
  // Load user document storage size
  try {
    const size = await GetUserDocumentSize()
    userStorage.value = size || 0
  } catch (error) {
    console.error('Error loading user storage:', error)
    userStorage.value = 0
  }
  
  // In a real app, this would fetch other
// Methods
function loadDashboardData() {
  // In a real app, this would fetch data from the backend
  console.log('Loading dashboard data...')
}

// Lifecycle
onMounted(() => {
  loadDashboardData()
})
</script>
