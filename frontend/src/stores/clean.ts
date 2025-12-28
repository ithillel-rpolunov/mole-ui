import { defineStore } from 'pinia'
import { ref } from 'vue'
import { CleanScanTargets, CleanExecute } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { handleError } from '../utils/errorHandler'

export const useCleanStore = defineStore('clean', () => {
  const categories = ref([])
  const loading = ref(false)
  const cleaning = ref(false)
  const progress = ref(0)
  const progressMessage = ref('')
  const result = ref(null)
  const error = ref(null)

  async function scanTargets() {
    loading.value = true
    error.value = null
    try {
      const data = await CleanScanTargets()
      categories.value = data || []
    } catch (err) {
      handleError(err, 'Clean scan')
      error.value = 'Failed to scan'
    } finally {
      loading.value = false
    }
  }

  async function executeClean(selectedCategories: string[], dryRun = false) {
    if (selectedCategories.length === 0) {
      error.value = 'Please select at least one category'
      return
    }

    cleaning.value = true
    progress.value = 0
    result.value = null
    error.value = null

    try {
      await CleanExecute(selectedCategories, dryRun)
    } catch (err) {
      handleError(err, 'Clean execution')
      error.value = 'Clean failed'
      cleaning.value = false
    }
  }

  function setupEventListeners() {
    EventsOn('clean:progress', (data) => {
      progress.value = data.percent
      progressMessage.value = data.message
    })

    EventsOn('clean:complete', (data) => {
      cleaning.value = false
      result.value = data
    })
  }

  return {
    categories,
    loading,
    cleaning,
    progress,
    progressMessage,
    result,
    error,
    scanTargets,
    executeClean,
    setupEventListeners,
  }
})
