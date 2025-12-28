<script setup>
import { ref, onMounted, computed } from 'vue'
import { OptimizeGetTasks, OptimizeExecute } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime/runtime'

const tasks = ref([])
const loading = ref(false)
const optimizing = ref(false)
const progress = ref(0)
const progressMessage = ref('')
const result = ref(null)

const selectedTasks = computed(() => {
  return tasks.value.filter(task => task.selected)
})

onMounted(async () => {
  await getTasks()

  EventsOn('optimize:progress', (data) => {
    progress.value = data.percent
    progressMessage.value = data.message
  })

  EventsOn('optimize:complete', (data) => {
    optimizing.value = false
    result.value = data
  })
})

async function getTasks() {
  loading.value = true
  try {
    const data = await OptimizeGetTasks()
    tasks.value = data.map(task => ({ ...task, selected: true }))
  } catch (error) {
    console.error('Failed to get tasks:', error)
    alert('Failed to get optimization tasks: ' + error)
  } finally {
    loading.value = false
  }
}

async function optimize() {
  const selected = selectedTasks.value
  if (selected.length === 0) {
    alert('Please select at least one task')
    return
  }

  if (!confirm(`Run ${selected.length} optimization tasks?`)) {
    return
  }

  optimizing.value = true
  const taskIDs = selected.map(task => task.id)

  try {
    await OptimizeExecute(taskIDs)
  } catch (error) {
    console.error('Optimization failed:', error)
    alert('Optimization failed: ' + error)
    optimizing.value = false
  }
}

function toggleTask(task) {
  task.selected = !task.selected
}
</script>

<template>
  <div class="optimize-tab">
    <h1>System Optimization</h1>
    <p class="subtitle">Optimize your Mac for better performance</p>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Loading optimization tasks...</p>
    </div>

    <div v-else-if="optimizing" class="optimizing">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <p class="progress-message">{{ progressMessage }}</p>
      <p class="progress-percent">{{ progress }}%</p>
    </div>

    <div v-else-if="result" class="result">
      <h2>Optimization Complete</h2>
      <p class="result-info">Successfully optimized {{ result.tasksCompleted }} tasks</p>
      <button @click="result = null" class="btn-primary">Done</button>
    </div>

    <div v-else class="task-list">
      <div class="task-items">
        <div
          v-for="task in tasks"
          :key="task.id"
          class="task-item"
          :class="{ selected: task.selected }"
          @click="toggleTask(task)"
        >
          <input type="checkbox" :checked="task.selected" @click.stop />
          <div class="task-info">
            <h3>{{ task.name }}</h3>
            <p>{{ task.description }}</p>
          </div>
        </div>
      </div>

      <div class="actions">
        <button @click="optimize" class="btn-primary" :disabled="selectedTasks.length === 0">
          Optimize Selected ({{ selectedTasks.length }})
        </button>
        <button @click="getTasks" class="btn-secondary">Refresh</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.optimize-tab {
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
.optimizing {
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

.result-info {
  font-size: 1.5rem;
  color: #d1d5db;
  margin: 0 0 2rem 0;
}

.task-items {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.task-item {
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

.task-item:hover {
  border-color: #8b5cf6;
}

.task-item.selected {
  border-color: #8b5cf6;
  background: #283448;
}

.task-info {
  flex: 1;
}

.task-info h3 {
  margin: 0 0 0.25rem 0;
  font-size: 1rem;
  color: #f3f4f6;
}

.task-info p {
  margin: 0;
  font-size: 0.875rem;
  color: #9ca3af;
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
