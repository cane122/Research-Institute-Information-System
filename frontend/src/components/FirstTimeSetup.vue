<template>
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h3 class="modal-title">Prva prijava - Postavljanje lozinke</h3>
      </div>
      
      <div class="first-time-message">
        <p>Dobrodošli u sistem! Ovo je Vaša prva prijava.</p>
        <p>Molimo Vas da postavite novu lozinku za Vaš nalog.</p>
      </div>
      
      <form @submit.prevent="handlePasswordSetup">
        <div class="form-group">
          <label for="newPassword">Nova lozinka</label>
          <input 
            type="password" 
            id="newPassword" 
            v-model="passwords.newPassword" 
            placeholder="Unesite novu lozinku (minimum 8 karaktera)" 
            required
            minlength="8"
            :disabled="isLoading"
          >
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">Potvrda lozinke</label>
          <input 
            type="password" 
            id="confirmPassword" 
            v-model="passwords.confirmPassword" 
            placeholder="Potvrdite novu lozinku" 
            required
            :disabled="isLoading"
          >
        </div>
        
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        
        <div class="form-actions">
          <button 
            type="submit" 
            class="btn btn-primary" 
            :disabled="isLoading || !isFormValid"
          >
            <span v-if="isLoading" class="spinner-small"></span>
            {{ isLoading ? 'Postavljanje...' : 'Postavi lozinku' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const emit = defineEmits(['setup-complete'])

const props = defineProps({
  username: {
    type: String,
    required: true
  }
})

const authStore = useAuthStore()
const router = useRouter()

const passwords = ref({
  newPassword: '',
  confirmPassword: ''
})

const isLoading = ref(false)
const error = ref('')

const isFormValid = computed(() => {
  return passwords.value.newPassword.length >= 8 && 
         passwords.value.newPassword === passwords.value.confirmPassword
})

async function handlePasswordSetup() {
  error.value = ''
  
  if (passwords.value.newPassword !== passwords.value.confirmPassword) {
    error.value = 'Lozinke se ne poklapaju'
    return
  }
  
  if (passwords.value.newPassword.length < 8) {
    error.value = 'Lozinka mora imati najmanje 8 karaktera'
    return
  }
  
  isLoading.value = true
  
  try {
    const result = await authStore.completeFirstTimeSetup(props.username, passwords.value.newPassword)
    
    if (result.success) {
      emit('setup-complete')
      router.push('/dashboard')
    } else {
      error.value = result.error || 'Greška pri postavljanju lozinke'
    }
  } catch (err) {
    console.error('Password setup error:', err)
    error.value = 'Greška prilikom komuniciranja sa serverom'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.first-time-message {
  background: #e3f2fd;
  border: 1px solid #bbdefb;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 25px;
}

.first-time-message p {
  color: #1565c0;
  margin-bottom: 8px;
  line-height: 1.5;
}

.first-time-message p:last-child {
  margin-bottom: 0;
}

.error-message {
  background: #ffebee;
  border: 1px solid #ffcdd2;
  border-radius: 4px;
  padding: 12px;
  color: #c62828;
  font-size: 14px;
  margin-bottom: 20px;
}

.spinner-small {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid #ffffff40;
  border-top: 2px solid #ffffff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-right: 8px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>