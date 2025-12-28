import { defineStore } from 'pinia'
import { ref } from 'vue'
import { OptimizeGetTasks, OptimizeExecute } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { handleError } from '../utils/errorHandler'

export const useOptimizeStore = defineStore('optimize', () => {
  const tasks = ref([])
  const loading = ref(false)
  const optimizing = ref(false)
  const progress = ref(0)
  const progressMessage = ref('')
  const result = ref(null)
  const error = ref(null)

  async function getTasks() {
    loading.value = true
    error.value = null
    try {
      const data = await OptimizeGetTasks()
      tasks.value = data || []
    } catch (err) {
      handleError(err, 'Get optimization tasks')
      error.value = 'Failed to get tasks'
    } finally {
      loading.value = false
    }
  }

  async function executeOptimize(taskIDs: string[]) {
    if (taskIDs.length === 0) {
      error.value = 'Please select at least one task'
      return
    }

    optimizing.value = true
    progress.value = 0
    result.value = null
    error.value = null

    try {
      await OptimizeExecute(taskIDs)
    } catch (err) {
      handleError(err, 'Optimize execution')
      error.value = 'Optimization failed'
      optimizing.value = false
    }
  }

  function setupEventListeners() {
    EventsOn('optimize:progress', (data) => {
      progress.value = data.percent
      progressMessage.value = data.message
    })

    EventsOn('optimize:complete', (data) => {
      optimizing.value = false
      result.value = data
    })
  }

  return {
    tasks,
    loading,
    optimizing,
    progress,
    progressMessage,
    result,
    error,
    getTasks,
    executeOptimize,
    setupEventListeners,
  }
})
