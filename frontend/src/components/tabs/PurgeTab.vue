<script setup>
import { ref, onMounted, computed } from 'vue'
import { PurgeScanProjects, PurgeExecute } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime/runtime'

const projects = ref([])
const loading = ref(false)
const purging = ref(false)
const progress = ref(0)
const progressMessage = ref('')
const searchQuery = ref('')

const filteredProjects = computed(() => {
  if (!searchQuery.value) return projects.value
  return projects.value.filter(project =>
    project.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    project.path.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const selectedProjects = computed(() => {
  return projects.value.filter(project => project.selected)
})

onMounted(async () => {
  await scanProjects()

  EventsOn('purge:progress', (data) => {
    progress.value = data.percent
    progressMessage.value = data.message
  })

  EventsOn('purge:complete', (data) => {
    purging.value = false
    alert(`Purged ${data.projectsCleaned} projects, freed ${formatSize(data.spaceFreed / 1024 / 1024)}`)
    scanProjects()
  })
})

async function scanProjects() {
  loading.value = true
  try {
    // Use empty string to scan HOME directory by default
    const data = await PurgeScanProjects('')
    // Auto-select old projects (isRecent === false)
    projects.value = data.map(project => ({
      ...project,
      selected: !project.isRecent
    }))
  } catch (error) {
    console.error('Scan failed:', error)
    alert('Failed to scan projects: ' + error)
  } finally {
    loading.value = false
  }
}

async function purge() {
  const selected = selectedProjects.value
  if (selected.length === 0) {
    alert('Please select at least one project')
    return
  }

  if (!confirm(`Purge ${selected.length} projects? This will remove build artifacts like node_modules, target, etc.`)) {
    return
  }

  purging.value = true
  const projectPaths = selected.map(project => project.path)

  try {
    await PurgeExecute(projectPaths)
  } catch (error) {
    console.error('Purge failed:', error)
    alert('Purge failed: ' + error)
    purging.value = false
  }
}

function toggleProject(project) {
  project.selected = !project.selected
}

function formatSize(mb) {
  if (mb >= 1024) {
    return (mb / 1024).toFixed(2) + ' GB'
  }
  return mb.toFixed(2) + ' MB'
}

function formatDate(dateStr) {
  if (!dateStr) return 'Unknown'
  const date = new Date(dateStr)
  return date.toLocaleDateString()
}
</script>

<template>
  <div class="purge-tab">
    <h1>Project Purge</h1>
    <p class="subtitle">Clean old project build artifacts to free up space</p>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Scanning projects...</p>
    </div>

    <div v-else-if="purging" class="purging">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <p class="progress-message">{{ progressMessage }}</p>
      <p class="progress-percent">{{ progress }}%</p>
    </div>

    <div v-else class="project-list">
      <div class="controls">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search projects..."
          class="search-input"
        />
        <button @click="scanProjects" class="btn-secondary">Refresh</button>
      </div>

      <div class="project-items">
        <div
          v-for="project in filteredProjects"
          :key="project.path"
          class="project-item"
          :class="{ selected: project.selected, recent: project.isRecent }"
          @click="toggleProject(project)"
        >
          <input type="checkbox" :checked="project.selected" @click.stop />
          <div class="project-info">
            <h3>{{ project.name }}</h3>
            <p class="project-path">{{ project.path }}</p>
            <p class="project-meta">
              {{ project.type }} · {{ formatSize(project.size) }} · Last modified: {{ formatDate(project.lastModified) }}
              <span v-if="project.isRecent" class="recent-badge">Recent</span>
            </p>
          </div>
        </div>
      </div>

      <div class="actions">
        <button @click="purge" class="btn-primary" :disabled="selectedProjects.length === 0">
          Purge Selected ({{ selectedProjects.length }})
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.purge-tab {
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

.loading,
.purging {
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

.project-items {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  max-height: 500px;
  overflow-y: auto;
  margin-bottom: 1.5rem;
}

.project-item {
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

.project-item:hover {
  border-color: #8b5cf6;
}

.project-item.selected {
  border-color: #8b5cf6;
  background: #283448;
}

.project-item.recent {
  border-color: #10b981;
}

.project-info {
  flex: 1;
}

.project-info h3 {
  margin: 0 0 0.25rem 0;
  font-size: 1rem;
  color: #f3f4f6;
}

.project-path {
  margin: 0;
  font-size: 0.75rem;
  color: #6b7280;
  font-family: monospace;
}

.project-meta {
  margin: 0.25rem 0 0 0;
  font-size: 0.875rem;
  color: #9ca3af;
}

.recent-badge {
  display: inline-block;
  padding: 0.125rem 0.5rem;
  background: #10b981;
  color: white;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  margin-left: 0.5rem;
}

.actions {
  display: flex;
  gap: 1rem;
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

.btn-secondary:hover {
  background: #4b5563;
}
</style>
