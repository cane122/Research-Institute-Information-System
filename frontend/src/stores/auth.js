import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const isLoading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const userName = computed(() => user.value?.name || '')

  // Actions
  async function login(credentials) {
    isLoading.value = true
    error.value = null
    
    try {
      // Mock login - replace with actual Wails backend call
      if (credentials.username === 'admin' && credentials.password === 'admin') {
        user.value = {
          id: 1,
          username: 'admin',
          name: 'Marko Petrović',
          email: 'admin@institut.rs',
          role: 'admin'
        }
        localStorage.setItem('user', JSON.stringify(user.value))
        return { success: true }
      } else if (credentials.username === 'user' && credentials.password === 'user') {
        user.value = {
          id: 2,
          username: 'user',
          name: 'Ana Jovanović',
          email: 'ana@institut.rs',
          role: 'user'
        }
        localStorage.setItem('user', JSON.stringify(user.value))
        return { success: true }
      } else {
        error.value = 'Neispravno korisničko ime ili lozinka'
        return { success: false, error: error.value }
      }
    } catch (err) {
      error.value = 'Greška prilikom prijave'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  function logout() {
    user.value = null
    localStorage.removeItem('user')
  }

  function initializeAuth() {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch {
        localStorage.removeItem('user')
      }
    }
  }

  return {
    user,
    isLoading,
    error,
    isAuthenticated,
    isAdmin,
    userName,
    login,
    logout,
    initializeAuth
  }
})
