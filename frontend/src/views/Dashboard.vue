<template>
  <Layout>
    <div class="dashboard">
      <!-- Dashboard Header -->
      <div class="dashboard-header">
        <h2>Dashboard</h2>
        <div class="breadcrumb">Početna > Dashboard</div>
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
            <h3>{{ stats.documents }}</h3>
            <p>Dokumenata u sistemu</p>
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
        
        <div class="stat-card users">
          <div class="stat-content">
            <h3>{{ stats.users }}</h3>
            <p>Korisnika u sistemu</p>
          </div>
          <div class="stat-icon">👥</div>
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

const authStore = useAuthStore()

// Reactive data
const stats = ref({
  projects: 12,
  documents: 158,
  tasks: 23,
  users: 8
})

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
])

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

<style scoped>
.dashboard {
  padding: 30px;
  max-width: 1400px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: 30px;
}

.dashboard-header h2 {
  color: #2c3e50;
  margin-bottom: 8px;
  font-size: 28px;
}

.breadcrumb {
  color: #7f8c8d;
  font-size: 14px;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: transform 0.3s, box-shadow 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-card.projects {
  border-left: 4px solid #3498db;
}

.stat-card.documents {
  border-left: 4px solid #2ecc71;
}

.stat-card.tasks {
  border-left: 4px solid #f39c12;
}

.stat-card.users {
  border-left: 4px solid #9b59b6;
}

.stat-content h3 {
  color: #2c3e50;
  font-size: 32px;
  margin-bottom: 5px;
  font-weight: 700;
}

.stat-content p {
  color: #7f8c8d;
  font-size: 14px;
  font-weight: 500;
}

.stat-icon {
  font-size: 32px;
  opacity: 0.8;
}

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 30px;
  margin-bottom: 30px;
}

/* Recent Activities */
.recent-activities {
  min-height: 400px;
}

.activities-list {
  max-height: 350px;
  overflow-y: auto;
}

.activity-item {
  padding: 15px 0;
  border-bottom: 1px solid #ecf0f1;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-content {
  flex: 1;
  padding-right: 15px;
}

.activity-title {
  color: #2c3e50;
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
}

.activity-desc {
  color: #7f8c8d;
  font-size: 13px;
}

.activity-time {
  color: #95a5a6;
  font-size: 12px;
  white-space: nowrap;
}

.no-activities {
  text-align: center;
  padding: 40px 0;
  color: #7f8c8d;
}

/* Quick Actions */
.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.quick-btn {
  display: flex;
  align-items: center;
  padding: 15px 20px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-icon {
  margin-right: 12px;
  font-size: 16px;
}

.quick-btn:hover {
  transform: translateX(5px);
}

/* Project Progress */
.project-progress {
  width: 100%;
}

.projects-list {
  max-height: 400px;
  overflow-y: auto;
}

.project-item {
  display: grid;
  grid-template-columns: 1fr auto auto;
  align-items: center;
  gap: 20px;
  padding: 20px 0;
  border-bottom: 1px solid #ecf0f1;
}

.project-item:last-child {
  border-bottom: none;
}

.project-info {
  min-width: 0;
}

.project-name {
  color: #2c3e50;
  font-weight: 600;
  font-size: 15px;
  margin-bottom: 4px;
}

.project-desc {
  color: #7f8c8d;
  font-size: 13px;
}

.project-progress-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 150px;
}

.progress-container {
  flex: 1;
  height: 8px;
  background: #ecf0f1;
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #3498db, #2ecc71);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #7f8c8d;
  font-weight: 600;
  min-width: 35px;
}

.project-status {
  min-width: 80px;
  text-align: right;
}

.no-projects {
  text-align: center;
  padding: 40px 0;
  color: #7f8c8d;
}

/* Responsive */
@media (max-width: 1024px) {
  .content-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .quick-actions-card {
    order: -1;
  }
  
  .quick-actions {
    flex-direction: row;
    flex-wrap: wrap;
  }
  
  .quick-btn {
    flex: 1;
    min-width: 140px;
  }
}

@media (max-width: 768px) {
  .dashboard {
    padding: 20px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 15px;
  }
  
  .stat-card {
    padding: 20px;
  }
  
  .stat-content h3 {
    font-size: 24px;
  }
  
  .project-item {
    grid-template-columns: 1fr;
    gap: 10px;
  }
  
  .project-progress-bar {
    min-width: auto;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .quick-actions {
    flex-direction: column;
  }
  
  .activity-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .activity-content {
    padding-right: 0;
  }
}
</style>
