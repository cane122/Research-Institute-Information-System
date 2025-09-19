import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from './stores/auth'

// Import views
import Login from './views/Login.vue'
import Dashboard from './views/Dashboard.vue'
import Projects from './views/Projects.vue'
import DocumentManagement from './views/documents/DocumentManagement.vue'
import DocumentAdd from './views/documents/DocumentAdd.vue'
import DocumentPreview from './views/documents/DocumentPreview.vue'
import Tasks from './views/Tasks.vue'
import Users from './views/Users.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/projects',
    name: 'Projects',
    component: Projects,
    meta: { requiresAuth: true }
  },
  {
    path: '/documents',
    name: 'DocumentManagement',
    component: DocumentManagement,
    meta: { requiresAuth: true }
  },
  {
    path: '/documents/add',
    name: 'DocumentAdd',
    component: DocumentAdd,
    meta: { requiresAuth: true }
  },
  {
    path: '/documents/preview/:id',
    name: 'DocumentPreview',
    component: DocumentPreview,
    meta: { requiresAuth: true },
    props: true
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: Tasks,
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: Users,
    meta: { requiresAuth: true, requiresAdmin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Check if route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }
  
  // Check if route requires guest (not authenticated)
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }
  
  // Check if route requires admin privileges
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next('/dashboard')
    return
  }
  
  next()
})

export default router
