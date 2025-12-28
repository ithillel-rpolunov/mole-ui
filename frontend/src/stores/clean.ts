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
    console.log('executeClean called with:', selectedCategories, dryRun)

    if (selectedCategories.length === 0) {
      error.value = 'Please select at least one category'
      return
    }

    cleaning.value = true
    progress.value = 0
    result.value = null
    error.value = null

    try {
      console.log('Calling CleanExecute backend method')
      await CleanExecute(selectedCategories, dryRun)
      console.log('CleanExecute backend method completed')
      // Note: cleaning.value will be set to false by the 'clean:complete' event
      // But if for some reason the event doesn't fire, we set a timeout fallback
      setTimeout(() => {
        if (cleaning.value && !result.value) {
          console.warn('clean:complete event not received, resetting state')
          cleaning.value = false
          error.value = 'Cleanup may have completed but status is unclear'
        }
      }, 1000)
    } catch (err) {
      console.error('CleanExecute error:', err)
      handleError(err, 'Clean execution')
      error.value = 'Clean failed'
      cleaning.value = false
    }
  }

  function setupEventListeners() {
    console.log('Setting up clean event listeners')
    EventsOn('clean:progress', (data) => {
      console.log('clean:progress event:', data)
      progress.value = data.percent
      progressMessage.value = data.message
    })

    EventsOn('clean:complete', (data) => {
      console.log('clean:complete event:', data)
      cleaning.value = false
      result.value = data
    })
  }

  function resetResult() {
    result.value = null
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
    resetResult,
  }
})
