<template>
  <div v-if="show" class="modal-overlay" @click.self="closeModal">
    <div class="modal-container">
      <div class="modal-header">
        <h3>Uslovi za fazu: {{ phaseName }}</h3>
        <button class="btn-close" @click="closeModal">&times;</button>
      </div>

      <div class="modal-body">
        <!-- Loading State -->
        <div v-if="loading" class="loading-state">
          <p>Učitavam uslove...</p>
        </div>

        <!-- Content when not loading -->
        <div v-else>
          <!-- Empty State -->
          <div v-if="conditions.length === 0 && !showForm" class="empty-state">
            <div class="empty-icon">📋</div>
            <h4>Nema definisanih uslova</h4>
            <p class="hint">Ova faza nema definisane uslove koji moraju biti ispunjeni pre prelaska u sledeću fazu.</p>
            <button class="btn btn-primary btn-large" @click="openForm">
              ➕ Kreiraj prvi uslov
            </button>
          </div>

          <!-- Conditions Table -->
          <div v-if="conditions.length > 0" class="conditions-table">
            <table>
              <thead>
                <tr>
                  <th>Opis</th>
                  <th>Kriterijum</th>
                  <th class="actions-col">Akcije</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="condition in conditions" :key="condition.uslov_id">
                  <td>{{ condition.opis }}</td>
                  <td class="kriterijum-cell">{{ condition.kriterijum }}</td>
                  <td class="actions-cell">
                    <button 
                      class="btn-icon-small" 
                      @click="editCondition(condition)"
                      title="Izmeni uslov"
                    >
                      ✏️
                    </button>
                    <button 
                      class="btn-icon-small danger" 
                      @click="deleteCondition(condition.uslov_id)"
                      title="Obriši uslov"
                    >
                      🗑️
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Add/Edit Condition Form -->
          <div v-if="showForm" class="condition-form">
            <h4>{{ editingCondition ? 'Izmeni uslov' : 'Dodaj novi uslov' }}</h4>
            
            <div class="form-group">
              <label>Opis uslova *</label>
              <input 
                v-model="formData.opis" 
                type="text" 
                placeholder="npr. Završiti dokumentaciju"
                class="form-control"
              />
            </div>

            <div class="form-group">
              <label>Kriterijum *</label>
              <textarea 
                v-model="formData.kriterijum" 
                placeholder="Detaljno opišite šta mora biti ispunjeno..."
                rows="3"
                class="form-control"
              ></textarea>
            </div>

            <div class="form-actions">
              <button class="btn btn-secondary" @click="cancelForm">Otkaži</button>
              <button class="btn btn-primary" @click="saveCondition">
                {{ editingCondition ? 'Sačuvaj izmene' : 'Dodaj uslov' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button v-if="!showForm" class="btn btn-primary" @click="openForm">
          ➕ Dodaj uslov
        </button>
        <button class="btn btn-secondary" @click="closeModal">Zatvori</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { GetConditionsByPhase, CreateCondition, UpdateCondition, DeleteCondition } from '../../wailsjs/go/main/App'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  phaseId: {
    type: Number,
    default: null
  },
  phaseName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'updated'])

const loading = ref(false)
const conditions = ref([])
const showForm = ref(false)
const editingCondition = ref(null)

const formData = ref({
  opis: '',
  kriterijum: ''
})

// Watch for phase changes
watch(() => props.phaseId, (newPhaseId) => {
  if (newPhaseId && props.show) {
    loadConditions()
  }
}, { immediate: true })

// Watch for modal visibility
watch(() => props.show, (newShow) => {
  if (newShow && props.phaseId) {
    showForm.value = false // Reset form state when opening modal
    loadConditions()
  }
})

async function loadConditions() {
  if (!props.phaseId) return

  loading.value = true
  try {
    const result = await GetConditionsByPhase(props.phaseId)
    // Ensure conditions is always an array, never null/undefined
    conditions.value = Array.isArray(result) ? result : []
  } catch (error) {
    console.error('Greška pri učitavanju uslova:', error)
    conditions.value = [] // Set to empty array on error
    alert('Greška pri učitavanju uslova')
  } finally {
    loading.value = false
  }
}

function openForm() {
  showForm.value = true
  editingCondition.value = null
  formData.value = {
    opis: '',
    kriterijum: ''
  }
}

function editCondition(condition) {
  showForm.value = true
  editingCondition.value = condition
  formData.value = {
    opis: condition.opis,
    kriterijum: condition.kriterijum
  }
}

function cancelForm() {
  showForm.value = false
  editingCondition.value = null
  formData.value = {
    opis: '',
    kriterijum: ''
  }
}

async function saveCondition() {
  if (!formData.value.opis.trim()) {
    alert('Opis uslova je obavezan')
    return
  }

  if (!formData.value.kriterijum.trim()) {
    alert('Kriterijum je obavezan')
    return
  }

  loading.value = true
  try {
    const conditionData = {
      faza_id: props.phaseId,
      opis: formData.value.opis.trim(),
      kriterijum: formData.value.kriterijum.trim()
    }

    if (editingCondition.value) {
      // Update existing condition
      conditionData.uslov_id = editingCondition.value.uslov_id
      await UpdateCondition(conditionData)
    } else {
      // Create new condition
      await CreateCondition(conditionData)
    }

    await loadConditions()
    cancelForm()
    emit('updated')
  } catch (error) {
    console.error('Greška pri čuvanju uslova:', error)
    alert('Greška pri čuvanju uslova: ' + error)
  } finally {
    loading.value = false
  }
}

async function deleteCondition(conditionId) {
  if (!confirm('Da li ste sigurni da želite da obrišete ovaj uslov?')) {
    return
  }

  loading.value = true
  try {
    await DeleteCondition(conditionId)
    await loadConditions()
    emit('updated')
  } catch (error) {
    console.error('Greška pri brisanju uslova:', error)
    alert('Greška pri brisanju uslova: ' + error)
  } finally {
    loading.value = false
  }
}

function resetForm() {
  showForm.value = false
  editingCondition.value = null
  formData.value = {
    opis: '',
    kriterijum: ''
  }
}

function closeModal() {
  resetForm()
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 700px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 30px;
  border-bottom: 1px solid #ecf0f1;
}

.modal-header h3 {
  margin: 0;
  color: #2c3e50;
  font-size: 20px;
}

.btn-close {
  background: none;
  border: none;
  font-size: 32px;
  color: #95a5a6;
  cursor: pointer;
  line-height: 1;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.3s;
}

.btn-close:hover {
  background: #ecf0f1;
  color: #2c3e50;
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 30px;
}

.loading-state {
  text-align: center;
  padding: 40px 20px;
  color: #7f8c8d;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #7f8c8d;
}

.empty-state .empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
  opacity: 0.3;
}

.empty-state h4 {
  color: #2c3e50;
  font-size: 20px;
  margin: 0 0 10px 0;
  font-weight: 600;
}

.empty-state .hint {
  font-size: 14px;
  color: #95a5a6;
  margin: 10px 0 30px 0;
  line-height: 1.6;
}

.empty-state .btn-large {
  padding: 12px 30px;
  font-size: 16px;
  font-weight: 600;
}

.conditions-table {
  margin-bottom: 20px;
}

.conditions-table table {
  width: 100%;
  border-collapse: collapse;
}

.conditions-table th {
  background: #f8f9fa;
  padding: 12px;
  text-align: left;
  font-weight: 600;
  color: #2c3e50;
  border-bottom: 2px solid #ecf0f1;
}

.conditions-table td {
  padding: 12px;
  border-bottom: 1px solid #ecf0f1;
}

.kriterijum-cell {
  color: #7f8c8d;
  font-size: 14px;
  max-width: 300px;
}

.actions-col {
  width: 100px;
  text-align: center;
}

.actions-cell {
  text-align: center;
}

.btn-icon-small {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  padding: 6px;
  border-radius: 4px;
  transition: background 0.3s;
  margin: 0 4px;
}

.btn-icon-small:hover {
  background: #f8f9fa;
}

.btn-icon-small.danger:hover {
  background: #f8d7da;
}

.condition-form {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-top: 20px;
}

.condition-form h4 {
  margin: 0 0 20px 0;
  color: #2c3e50;
  font-size: 16px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #2c3e50;
  font-size: 14px;
}

.form-control {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #bdc3c7;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
  transition: border-color 0.3s;
}

.form-control:focus {
  outline: none;
  border-color: #3498db;
}

textarea.form-control {
  resize: vertical;
  min-height: 80px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}

.modal-footer {
  padding: 20px 30px;
  border-top: 1px solid #ecf0f1;
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary {
  background: #3498db;
  color: white;
}

.btn-primary:hover {
  background: #2980b9;
}

.btn-secondary {
  background: #ecf0f1;
  color: #2c3e50;
}

.btn-secondary:hover {
  background: #bdc3c7;
}

@media (max-width: 768px) {
  .modal-container {
    width: 95%;
    max-height: 95vh;
  }

  .modal-header,
  .modal-body,
  .modal-footer {
    padding: 15px 20px;
  }

  .conditions-table {
    overflow-x: auto;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}
</style>
