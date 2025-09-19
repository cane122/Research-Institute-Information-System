import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { Login, Logout, GetCurrentUser, TestConnection, CompleteFirstTimeSetup } from '../../wailsjs/go/main/App.js'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const isLoading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.uloga === 'admin')
  const userName = computed(() => {
    if (!user.value) return ''
    return user.value.ime && user.value.prezime 
      ? `${user.value.ime} ${user.value.prezime}` 
      : user.value.korisnickoIme || ''
  })

  // Actions
  async function login(credentials) {
    isLoading.value = true
    error.value = null
    
    try {
      console.log('Attempting login with:', credentials.username)
      
      // Call Wails backend login function
      const response = await Login(credentials.username, credentials.password)
      
      console.log('Login response:', response)
      
      if (response && response.success) {
        user.value = response.user
        
        // Store user data in localStorage for persistence
        localStorage.setItem('user', JSON.stringify(user.value))
        
        return { 
          success: true, 
          message: response.message,
          isFirstTime: response.message === 'FIRST_TIME_LOGIN'
        }
      } else {
        error.value = response?.message || 'Neispravna prijava'
        return { 
          success: false, 
          error: error.value 
        }
      }
    } catch (err) {
      console.error('Login error:', err)
      error.value = 'Greška prilikom komuniciranja sa serverom'
      return { 
        success: false, 
        error: error.value 
      }
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    try {
      // Call backend logout
      await Logout()
      user.value = null
      localStorage.removeItem('user')
    } catch (err) {
      console.error('Logout error:', err)
      // Even if backend call fails, clear local state
      user.value = null
      localStorage.removeItem('user')
    }
  }

  async function getCurrentUser() {
    try {
      const currentUser = await GetCurrentUser()
      if (currentUser) {
        user.value = currentUser
        localStorage.setItem('user', JSON.stringify(currentUser))
      }
      return currentUser
    } catch (err) {
      console.error('Error getting current user:', err)
      return null
    }
  }

  async function testBackendConnection() {
    try {
      const result = await TestConnection()
      console.log('Backend connection test:', result)
      return result
    } catch (err) {
      console.error('Backend connection failed:', err)
      return { error: 'Connection failed' }
    }
  }

  function initializeAuth() {
    // First try to get user from localStorage
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch (err) {
        console.error('Error parsing saved user:', err)
        localStorage.removeItem('user')
      }
    }

    // Then try to get current user from backend
    getCurrentUser().catch(err => {
      console.error('Error initializing auth from backend:', err)
    })

    // Test backend connection on initialization
    testBackendConnection().catch(err => {
      console.error('Backend connection test failed on init:', err)
    })
  }

  async function completeFirstTimeSetup(username, newPassword) {
    try {
      const result = await CompleteFirstTimeSetup(username, newPassword)
      
      if (result && result.success) {
        return { 
          success: true, 
          message: result.message 
        }
      } else {
        return { 
          success: false, 
          error: result?.message || 'Greška pri postavljanju lozinke' 
        }
      }
    } catch (err) {
      console.error('First time setup error:', err)
      return { 
        success: false, 
        error: 'Greška prilikom komuniciranja sa serverom' 
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
    getCurrentUser,
    testBackendConnection,
    completeFirstTimeSetup,
    initializeAuth
  }
})
