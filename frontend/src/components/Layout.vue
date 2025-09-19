<template>
  <div class="app-layout">
    <!-- Header -->
    <div class="header">
      <div class="header-left">
        <button class="sidebar-toggle" @click="toggleSidebar">
          <span class="toggle-icon">{{ isCollapsed ? '☰' : '✕' }}</span>
        </button>
        <h1>Istraživačko razvojni centar</h1>
      </div>
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
      <div class="sidebar" :class="{ collapsed: isCollapsed }">
        <nav class="nav-menu">
          <router-link 
            v-for="item in navigationItems" 
            :key="item.path"
            :to="item.path" 
            class="nav-item"
            :class="{ active: $route.path === item.path }"
            :title="isCollapsed ? item.text : ''"
          >
            <span class="nav-icon">{{ item.icon }}</span>
            <span class="nav-text" :class="{ hidden: isCollapsed }">{{ item.text }}</span>
          </router-link>
        </nav>
      </div>
      
      <!-- Main Content -->
      <div class="main-content" :class="{ expanded: isCollapsed }">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// Reactive state for sidebar collapse
const isCollapsed = ref(false)

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
function toggleSidebar() {
  isCollapsed.value = !isCollapsed.value
}

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

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.sidebar-toggle {
  background: transparent;
  border: 2px solid #3498db;
  color: #3498db;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  font-size: 18px;
}

.sidebar-toggle:hover {
  background: #3498db;
  color: white;
  transform: scale(1.05);
}

.toggle-icon {
  transition: transform 0.3s ease;
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
  transition: width 0.3s ease;
}

.sidebar.collapsed {
  width: 80px;
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
  position: relative;
}

.sidebar.collapsed .nav-item {
  justify-content: center;
  padding: 15px 10px;
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
  transition: margin 0.3s ease;
}

.sidebar.collapsed .nav-icon {
  margin-right: 0;
}

.nav-text {
  font-size: 14px;
  transition: opacity 0.3s ease;
  white-space: nowrap;
}

.nav-text.hidden {
  opacity: 0;
  width: 0;
  overflow: hidden;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  background: #f8f9fa;
  transition: margin-left 0.3s ease;
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
    width: 80px;
  }
  
  .sidebar.collapsed {
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
  .sidebar-toggle {
    width: 35px;
    height: 35px;
    font-size: 16px;
  }
  
  .sidebar {
    position: fixed;
    left: -260px;
    width: 260px;
    height: 100%;
    z-index: 1000;
    transition: left 0.3s, width 0.3s;
  }
  
  .sidebar.collapsed {
    left: -80px;
    width: 80px;
  }
  
  .nav-text {
    display: block;
  }
  
  .nav-text.hidden {
    display: none;
  }
  
  .nav-item {
    justify-content: flex-start;
    padding: 15px 25px;
  }
  
  .sidebar.collapsed .nav-item {
    justify-content: center;
    padding: 15px 10px;
  }
  
  .nav-icon {
    margin-right: 12px;
  }
  
  .sidebar.collapsed .nav-icon {
    margin-right: 0;
  }
}
</style>
