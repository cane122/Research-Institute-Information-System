<template>
  <div class="app-layout">
    <!-- Header -->
    <div class="header">
      <h1>Istraživačko razvojni centar</h1>
      <div class="user-menu">
        <div class="user-info">
          <span class="user-icon">👤</span>
          <div class="user-details">
            <span class="user-name">{{ authStore.userName }}</span>
            <span class="user-role">{{ userRoleText }}</span>
          </div>
        </div>
        <button class="logout-btn" @click="handleLogout">Odjavi se</button>
      </div>
    </div>
    
    <!-- Main Layout -->
    <div class="layout">
      <!-- Sidebar -->
      <div class="sidebar">
        <nav class="nav-menu">
          <router-link 
            v-for="item in navigationItems" 
            :key="item.path"
            :to="item.path" 
            class="nav-item"
            :class="{ active: $route.path === item.path }"
          >
            <span class="nav-icon">{{ item.icon }}</span>
            <span class="nav-text">{{ item.text }}</span>
          </router-link>
        </nav>
      </div>
      
      <!-- Main Content -->
      <div class="main-content">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// Computed
const userRoleText = computed(() => {
  return authStore.isAdmin ? 'Administrator' : 'Korisnik'
})

const navigationItems = computed(() => {
  const items = [
    { path: '/dashboard', icon: '📊', text: 'Dashboard' },
    { path: '/projects', icon: '📁', text: 'Projekti' },
    { path: '/tasks', icon: '📋', text: 'Zadaci' },
    { path: '/documents', icon: '📄', text: 'Dokumenti' }
  ]
  
  // Add admin-only items
  if (authStore.isAdmin) {
    items.push(
      { path: '/users', icon: '👥', text: 'Korisnici' }
    )
  }
  
  return items
})

// Methods
function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Header */
.header {
  background: #2c3e50;
  color: white;
  padding: 15px 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 100;
}

.header h1 {
  font-size: 20px;
  font-weight: 600;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-icon {
  font-size: 24px;
}

.user-details {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.user-name {
  font-weight: 600;
  font-size: 14px;
}

.user-role {
  font-size: 12px;
  color: #bdc3c7;
}

.logout-btn {
  background: #e74c3c;
  color: white;
  border: none;
  padding: 8px 15px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.3s;
}

.logout-btn:hover {
  background: #c0392b;
}

/* Layout */
.layout {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 260px;
  background: white;
  border-right: 1px solid #e9ecef;
  overflow-y: auto;
  box-shadow: 2px 0 4px rgba(0, 0, 0, 0.05);
}

.nav-menu {
  padding: 20px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 15px 25px;
  cursor: pointer;
  border-left: 3px solid transparent;
  transition: all 0.3s;
  text-decoration: none;
  color: #2c3e50;
}

.nav-item:hover {
  background: #f8f9fa;
  border-left-color: #3498db;
}

.nav-item.active {
  background: #ecf0f1;
  border-left-color: #3498db;
  color: #3498db;
  font-weight: 600;
}

.nav-icon {
  font-size: 18px;
  margin-right: 12px;
  min-width: 20px;
}

.nav-text {
  font-size: 14px;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  background: #f8f9fa;
}

/* Responsive */
@media (max-width: 768px) {
  .header {
    padding: 12px 20px;
  }
  
  .header h1 {
    font-size: 18px;
  }
  
  .user-details {
    display: none;
  }
  
  .sidebar {
    width: 60px;
  }
  
  .nav-text {
    display: none;
  }
  
  .nav-item {
    justify-content: center;
    padding: 15px 10px;
  }
  
  .nav-icon {
    margin-right: 0;
  }
}

@media (max-width: 480px) {
  .sidebar {
    position: fixed;
    left: -260px;
    width: 260px;
    height: 100%;
    z-index: 1000;
    transition: left 0.3s;
  }
  
  .sidebar.open {
    left: 0;
  }
  
  .nav-text {
    display: block;
  }
  
  .nav-item {
    justify-content: flex-start;
    padding: 15px 25px;
  }
  
  .nav-icon {
    margin-right: 12px;
  }
}
</style>
