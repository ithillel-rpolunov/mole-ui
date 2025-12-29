<script setup>
import { onMounted, ref } from 'vue'
import { useCleanStore } from '../../stores/clean'
import { storeToRefs } from 'pinia'
import ConfirmDialog from '../shared/ConfirmDialog.vue'

const store = useCleanStore()
const { categories, loading, cleaning, progress, progressMessage, result, error } = storeToRefs(store)

const showConfirmDialog = ref(false)
const selectedCategories = ref([])

onMounted(async () => {
  await store.scanTargets()
  store.setupEventListeners()
})

async function startClean() {
  console.log('startClean called')
  selectedCategories.value = categories.value
    .filter(cat => cat.enabled)
    .map(cat => cat.id)

  console.log('Selected categories:', selectedCategories.value)

  if (selectedCategories.value.length === 0) {
    console.log('No categories selected')
    error.value = 'Please select at least one category'
    return
  }

  // Clear any previous errors
  error.value = null

  console.log('Showing confirmation dialog')
  showConfirmDialog.value = true
}

async function handleConfirm() {
  console.log('User confirmed, calling store.executeClean')
  await store.executeClean(selectedCategories.value, false)
  console.log('store.executeClean completed')
}

function handleCancel() {
  console.log('User cancelled cleanup')
}

function toggleCategory(category) {
  category.enabled = !category.enabled
}

function formatSize(mb) {
  if (mb >= 1024) {
    return (mb / 1024).toFixed(2) + ' GB'
  }
  return mb.toFixed(2) + ' MB'
}
</script>

<template>
  <div class="clean-tab">
    <h1>System Cleanup</h1>
    <p class="subtitle">Deep clean your Mac to reclaim disk space</p>

    <ConfirmDialog
      v-model:show="showConfirmDialog"
      title="Confirm Cleanup"
      message="This will clean the selected categories. Continue?"
      confirm-text="Start Cleanup"
      cancel-text="Cancel"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    />

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Scanning system...</p>
    </div>

    <div v-else-if="cleaning" class="cleaning">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <p class="progress-message">{{ progressMessage }}</p>
      <p class="progress-percent">{{ progress }}%</p>
    </div>

    <div v-else-if="result" class="result">
      <h2>✓ Cleanup Complete</h2>
      <p class="space-freed">Space freed: {{ formatSize(result.spaceFreed / 1024 / 1024) }}</p>
      <button @click="store.resetResult" class="btn-primary">Done</button>
    </div>

    <div v-else class="categories">
      <div
        v-for="category in categories"
        :key="category.id"
        class="category-item"
        :class="{ disabled: !category.enabled }"
        @click="toggleCategory(category)"
      >
        <input type="checkbox" :checked="category.enabled" @click.stop />
        <div class="category-info">
          <h3>{{ category.name }}</h3>
          <p>{{ category.description }}</p>
        </div>
        <div class="category-size">
          {{ formatSize(category.estimatedMB) }}
        </div>
      </div>

      <div class="actions">
        <button @click="startClean" class="btn-primary">Start Cleanup</button>
        <button @click="store.scanTargets" class="btn-secondary">Refresh</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.clean-tab {
  max-width: 900px;
}

h1 {
  font-size: 2rem;
  margin: 0 0 0.5rem 0;
}

.subtitle {
  color: #9ca3af;
  margin: 0 0 2rem 0;
}

.error-message {
  background: #7f1d1d;
  color: #fecaca;
  padding: 1rem;
  border-radius: 6px;
  margin-bottom: 1.5rem;
  border: 1px solid #991b1b;
}

.loading,
.cleaning {
  text-align: center;
  padding: 4rem 2rem;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #374151;
  border-top-color: #8b5cf6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #374151;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 1rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #8b5cf6, #ec4899);
  transition: width 0.3s;
}

.progress-message {
  color: #d1d5db;
  margin: 0.5rem 0;
}

.progress-percent {
  font-size: 2rem;
  font-weight: 700;
  color: #8b5cf6;
  margin: 1rem 0 0 0;
}

.result {
  text-align: center;
  padding: 4rem 2rem;
}

.result h2 {
  color: #10b981;
  font-size: 2.5rem;
  margin: 0 0 1rem 0;
}

.space-freed {
  font-size: 1.5rem;
  color: #d1d5db;
  margin: 0 0 2rem 0;
}

.categories {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem;
  background: #1f2937;
  border: 2px solid #374151;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.category-item:hover {
  border-color: #8b5cf6;
  background: #283448;
}

.category-item.disabled {
  opacity: 0.5;
}

.category-info {
  flex: 1;
}

.category-info h3 {
  margin: 0 0 0.25rem 0;
  font-size: 1rem;
  color: #f3f4f6;
}

.category-info p {
  margin: 0;
  font-size: 0.875rem;
  color: #9ca3af;
}

.category-size {
  font-size: 1.125rem;
  font-weight: 600;
  color: #8b5cf6;
}

.actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  color: white;
  flex: 1;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}

.btn-secondary {
  background: #374151;
  color: #d1d5db;
}

.btn-secondary:hover {
  background: #4b5563;
}
</style>
