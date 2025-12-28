<script setup>
import { ref, onMounted, computed } from 'vue'
import { UninstallScanApps, UninstallApps, UninstallGetRelatedFiles } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime/runtime'

const apps = ref([])
const loading = ref(false)
const uninstalling = ref(false)
const progress = ref(0)
const progressMessage = ref('')
const searchQuery = ref('')
const showRelatedFiles = ref(false)
const relatedFiles = ref([])
const selectedApp = ref(null)

const filteredApps = computed(() => {
  if (!searchQuery.value) return apps.value
  return apps.value.filter(app =>
    app.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const selectedApps = computed(() => {
  return apps.value.filter(app => app.selected)
})

onMounted(async () => {
  await scanApps(false)

  EventsOn('uninstall:progress', (data) => {
    progress.value = data.percent
    progressMessage.value = data.message
  })

  EventsOn('uninstall:complete', (data) => {
    uninstalling.value = false
    alert(`Uninstalled ${data.appsRemoved} apps, freed ${formatSize(data.spaceFreed / 1024 / 1024)}`)
    scanApps(true)
  })
})

async function scanApps(forceRescan) {
  loading.value = true
  try {
    const data = await UninstallScanApps(forceRescan)
    apps.value = data.map(app => ({ ...app, selected: false }))
  } catch (error) {
    console.error('Scan failed:', error)
    alert('Failed to scan apps: ' + error)
  } finally {
    loading.value = false
  }
}

async function showRelated(app) {
  selectedApp.value = app
  try {
    const files = await UninstallGetRelatedFiles(app.bundleId)
    relatedFiles.value = files
    showRelatedFiles.value = true
  } catch (error) {
    console.error('Failed to get related files:', error)
  }
}

async function uninstall() {
  const selected = selectedApps.value
  if (selected.length === 0) {
    alert('Please select at least one app')
    return
  }

  if (!confirm(`Uninstall ${selected.length} apps?`)) {
    return
  }

  uninstalling.value = true
  const bundleIds = selected.map(app => app.bundleId)

  try {
    await UninstallApps(bundleIds)
  } catch (error) {
    console.error('Uninstall failed:', error)
    alert('Uninstall failed: ' + error)
    uninstalling.value = false
  }
}

function toggleApp(app) {
  app.selected = !app.selected
}

function formatSize(bytes) {
  if (bytes >= 1024 * 1024 * 1024) {
    return (bytes / 1024 / 1024 / 1024).toFixed(2) + ' GB'
  }
  return (bytes / 1024 / 1024).toFixed(2) + ' MB'
}
</script>

<template>
  <div class="uninstall-tab">
    <h1>Uninstall Applications</h1>
    <p class="subtitle">Remove apps and their associated files</p>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Scanning applications...</p>
    </div>

    <div v-else-if="uninstalling" class="uninstalling">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <p class="progress-message">{{ progressMessage }}</p>
      <p class="progress-percent">{{ progress }}%</p>
    </div>

    <div v-else class="app-list">
      <div class="controls">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search apps..."
          class="search-input"
        />
        <button @click="scanApps(true)" class="btn-secondary">Refresh</button>
      </div>

      <div class="app-items">
        <div
          v-for="app in filteredApps"
          :key="app.bundleId"
          class="app-item"
          :class="{ selected: app.selected }"
          @click="toggleApp(app)"
        >
          <input type="checkbox" :checked="app.selected" @click.stop />
          <div class="app-info">
            <h3>{{ app.name }}</h3>
            <p>{{ app.path }}</p>
            <p class="app-meta">{{ formatSize(app.size) }} · {{ app.age }}</p>
          </div>
          <button @click.stop="showRelated(app)" class="btn-related">
            Show Files
          </button>
        </div>
      </div>

      <div class="actions">
        <button @click="uninstall" class="btn-primary" :disabled="selectedApps.length === 0">
          Uninstall Selected ({{ selectedApps.length }})
        </button>
      </div>
    </div>

    <div v-if="showRelatedFiles" class="modal" @click="showRelatedFiles = false">
      <div class="modal-content" @click.stop>
        <h3>Related Files for {{ selectedApp?.name }}</h3>
        <ul class="file-list">
          <li v-for="file in relatedFiles" :key="file">{{ file }}</li>
        </ul>
        <button @click="showRelatedFiles = false" class="btn-primary">Close</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.uninstall-tab { max-width: 900px; }
h1 { font-size: 2rem; margin: 0 0 0.5rem 0; }
.subtitle { color: #9ca3af; margin: 0 0 2rem 0; }

.loading, .uninstalling {
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

.controls {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.search-input {
  flex: 1;
  padding: 0.75rem;
  background: #1f2937;
  border: 1px solid #374151;
  border-radius: 6px;
  color: #f3f4f6;
  font-size: 1rem;
}

.app-items {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  max-height: 500px;
  overflow-y: auto;
}

.app-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #1f2937;
  border: 2px solid #374151;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.app-item:hover {
  border-color: #8b5cf6;
}

.app-item.selected {
  border-color: #8b5cf6;
  background: #283448;
}

.app-info {
  flex: 1;
}

.app-info h3 {
  margin: 0 0 0.25rem 0;
  font-size: 1rem;
}

.app-info p {
  margin: 0;
  font-size: 0.875rem;
  color: #9ca3af;
}

.app-meta {
  margin-top: 0.25rem !important;
  font-size: 0.75rem !important;
}

.btn-related {
  padding: 0.5rem 1rem;
  background: #374151;
  border: none;
  border-radius: 6px;
  color: #d1d5db;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-related:hover {
  background: #4b5563;
}

.actions {
  margin-top: 1.5rem;
}

.btn-primary, .btn-secondary {
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
  width: 100%;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #374151;
  color: #d1d5db;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: #1f2937;
  padding: 2rem;
  border-radius: 8px;
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
}

.file-list {
  list-style: none;
  padding: 0;
  margin: 1rem 0;
}

.file-list li {
  padding: 0.5rem;
  background: #111827;
  margin-bottom: 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: monospace;
}
</style>
