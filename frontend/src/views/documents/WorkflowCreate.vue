<template>
  <Layout>
    <div class="workflow-create-page">
      <div class="workflow-create-main">
        <h2>Kreiraj novi radni tok</h2>
        <div class="form-section">
          <label>Naziv radnog toka</label>
          <input v-model="naziv" type="text" placeholder="Unesite naziv radnog toka" />
        </div>
        <div class="form-section">
          <label>Faze radnog toka</label>
          <div v-for="(faza, idx) in faze" :key="idx" class="faza-row">
            <input v-model="faza.naziv" type="text" :placeholder="`Naziv faze #${idx+1}`" />
            <button v-if="idx === faze.length-1" class="btn btn-small plus-btn" @click="addFaza" title="Dodaj fazu">+</button>
            <button v-if="faze.length > 1" class="btn btn-small remove-btn" @click="removeFaza(idx)" title="Ukloni fazu">×</button>
          </div>
          <div v-if="faze.length < 2" class="error">Radni tok mora imati bar 2 faze.</div>
        </div>
        <div class="form-section">
          <button class="btn btn-primary" :disabled="!canSubmit" @click="submitWorkflow">Dodaj radni tok</button>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed } from 'vue'
import Layout from '../../components/Layout.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const naziv = ref('')
const faze = ref([{ naziv: '' }])

function addFaza() {
  faze.value.push({ naziv: '' })
}

function removeFaza(idx) {
  if (faze.value.length > 1) {
    faze.value.splice(idx, 1)
  }
}

const canSubmit = computed(() => {
  if (!naziv.value.trim()) return false
  if (faze.value.length < 2) return false
  return faze.value.every(f => f.naziv.trim())
})

async function submitWorkflow() {
  if (!canSubmit.value) return
  // Backend expects: { naziv, tip_toka: 'DOKUMENTACIJA', faze: [naziv1, naziv2, ...] }
  try {
    const wf = {
      naziv: naziv.value.trim(),
      tip_toka: 'DOKUMENTACIJA',
      faze: faze.value.map(f => f.naziv.trim())
    }
    await window.go?.main?.App?.CreateWorkflowWithPhases(wf)
    alert('Radni tok uspešno dodat!')
    router.push('/documents/create')
  } catch (e) {
    alert(e?.message || 'Greška pri dodavanju radnog toka')
  }
}
</script>

<style scoped>
.workflow-create-page {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: 100vh;
  background: #f8fafc;
}
.workflow-create-main {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  padding: 32px 40px;
  margin-top: 32px;
  width: 480px;
}
.form-section {
  margin-bottom: 24px;
}
.form-section label {
  font-weight: 600;
  display: block;
  margin-bottom: 8px;
}
input[type="text"] {
  width: 100%;
  padding: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 15px;
  margin-bottom: 4px;
}
.faza-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.btn {
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 15px;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-small {
  font-size: 13px;
  padding: 4px 10px;
}
.plus-btn {
  background: #22c55e;
}
.remove-btn {
  background: #ef4444;
}
.error {
  color: #b91c1c;
  font-size: 14px;
  margin-top: 4px;
}
</style>
