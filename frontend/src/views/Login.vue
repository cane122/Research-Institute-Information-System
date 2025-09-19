<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo">
        <h1>Istraživačko razvojni centar</h1>
        <p>Sistem za upravljanje informacijama</p>
      </div>
      
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">Korisničko ime</label>
          <input 
            type="text" 
            id="username" 
            v-model="credentials.username" 
            placeholder="Unesite korisničko ime" 
            required
            :disabled="isLoading"
            autocomplete="username"
          >
        </div>
        
        <div class="form-group">
          <label for="password">Lozinka</label>
          <input 
            type="password" 
            id="password" 
            v-model="credentials.password" 
            placeholder="Unesite lozinku" 
            required
            :disabled="isLoading"
            autocomplete="current-password"
          >
        </div>
        
        <div class="checkbox-group">
          <input type="checkbox" id="remember" v-model="rememberMe" :disabled="isLoading">
          <label for="remember">Zapamti me</label>
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        
        <button type="submit" class="btn btn-primary login-btn" :disabled="isLoading">
          <span v-if="isLoading" class="spinner-small"></span>
          {{ isLoading ? 'Prijavljivanje...' : 'Prijavite se' }}
        </button>
      </form>
      
      <div class="demo-credentials">
        <h4>Demo kredencijali:</h4>
        <p><strong>Administrator:</strong> admin / password</p>
        <p><strong>Rukovodilac:</strong> manager1 / password</p>
        <p><strong>Istraživač:</strong> researcher1 / password</p>
      </div>
    </div>
    
    <!-- First Time Setup Modal -->
    <FirstTimeSetup
      v-if="showFirstTimeSetup"
      :username="credentials.username"
      @setup-complete="handleSetupComplete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import FirstTimeSetup from '../components/FirstTimeSetup.vue'

const router = useRouter()
const authStore = useAuthStore()

// Reactive data
const credentials = ref({
  username: '',
  password: ''
})

const rememberMe = ref(false)
const loginError = ref('')
const showFirstTimeSetup = ref(false)

// Computed properties
const isLoading = computed(() => authStore.isLoading)
const error = computed(() => loginError.value || authStore.error)

// Methods
async function handleLogin() {
  loginError.value = ''
  
  try {
    const result = await authStore.login(credentials.value)
    
    if (result.success) {
      console.log('Login successful:', result)
      
      // Handle first-time login
      if (result.isFirstTime) {
        console.log('First time login detected - showing setup modal')
        showFirstTimeSetup.value = true
      } else {
        router.push('/dashboard')
      }
    } else {
      loginError.value = result.error || 'Neuspešna prijava'
    }
  } catch (err) {
    console.error('Login error:', err)
    loginError.value = 'Greška prilikom komunikacije sa serverom'
  }
}

function handleSetupComplete() {
  showFirstTimeSetup.value = false
  // The FirstTimeSetup component will handle the redirect to dashboard
}

// Test backend connection on component mount
onMounted(async () => {
  try {
    const connectionTest = await authStore.testBackendConnection()
    console.log('Backend connection test result:', connectionTest)
  } catch (err) {
    console.error('Backend connection test failed:', err)
    loginError.value = 'Greška: Nema veze sa serverom'
  }
})
</script>