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
        <p><strong>Admin:</strong> admin / admin</p>
        <p><strong>Korisnik:</strong> user / user</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// Reactive data
const credentials = ref({
  username: '',
  password: ''
})

const rememberMe = ref(false)

// Computed properties
const isLoading = authStore.isLoading
const error = authStore.error

// Methods
async function handleLogin() {
  const result = await authStore.login(credentials.value)
  
  if (result.success) {
    router.push('/dashboard')
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  width: 100%;
  max-width: 400px;
}

.logo {
  text-align: center;
  margin-bottom: 30px;
}

.logo h1 {
  color: #2c3e50;
  font-size: 24px;
  margin-bottom: 5px;
}

.logo p {
  color: #7f8c8d;
  font-size: 14px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #2c3e50;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #bdc3c7;
  border-radius: 5px;
  font-size: 14px;
  transition: all 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 5px rgba(52, 152, 219, 0.3);
}

.form-group input:disabled {
  background: #f8f9fa;
  color: #6c757d;
}

.checkbox-group {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.checkbox-group input[type="checkbox"] {
  margin-right: 8px;
  width: auto;
}

.checkbox-group label {
  margin-bottom: 0;
  font-weight: normal;
  cursor: pointer;
}

.login-btn {
  width: 100%;
  padding: 12px;
  font-size: 16px;
  position: relative;
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.spinner-small {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid #ffffff;
  border-top: 2px solid transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-right: 8px;
}

.error-message {
  background: #f8d7da;
  color: #721c24;
  padding: 12px;
  border-radius: 5px;
  margin-bottom: 20px;
  text-align: center;
  font-size: 14px;
}

.demo-credentials {
  margin-top: 30px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  text-align: center;
}

.demo-credentials h4 {
  color: #2c3e50;
  margin-bottom: 10px;
  font-size: 14px;
}

.demo-credentials p {
  margin: 5px 0;
  font-size: 12px;
  color: #6c757d;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (max-width: 480px) {
  .login-card {
    padding: 30px 20px;
  }
  
  .logo h1 {
    font-size: 20px;
  }
}
</style>
