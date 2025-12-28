import { defineStore } from 'pinia'
import { ref } from 'vue'
import { UninstallScanApps, UninstallApps, UninstallGetRelatedFiles } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export const useUninstallStore = defineStore('uninstall', () => {
  const apps = ref([])
  const loading = ref(false)
  const uninstalling = ref(false)
  const progress = ref(0)
  const progressMessage = ref('')
  const result = ref(null)
  const error = ref(null)
  const relatedFiles = ref([])

  async function scanApps(forceRescan = false) {
    loading.value = true
    error.value = null
    try {
      const data = await UninstallScanApps(forceRescan)
      apps.value = data || []
    } catch (err) {
      error.value = 'Failed to scan apps: ' + err
      console.error('Scan failed:', err)
    } finally {
      loading.value = false
    }
  }

  async function uninstallApps(bundleIDs: string[]) {
    if (bundleIDs.length === 0) {
      error.value = 'Please select at least one app'
      return
    }

    uninstalling.value = true
    progress.value = 0
    result.value = null
    error.value = null

    try {
      await UninstallApps(bundleIDs)
    } catch (err) {
      error.value = 'Uninstall failed: ' + err
      uninstalling.value = false
    }
  }

  async function getRelatedFiles(bundleID: string) {
    error.value = null
    try {
      const data = await UninstallGetRelatedFiles(bundleID)
      relatedFiles.value = data || []
      return data
    } catch (err) {
      error.value = 'Failed to get related files: ' + err
      console.error('Get related files failed:', err)
      return []
    }
  }

  function setupEventListeners() {
    EventsOn('uninstall:progress', (data) => {
      progress.value = data.percent
      progressMessage.value = data.message
    })

    EventsOn('uninstall:complete', (data) => {
      uninstalling.value = false
      result.value = data
    })
  }

  return {
    apps,
    loading,
    uninstalling,
    progress,
    progressMessage,
    result,
    error,
    relatedFiles,
    scanApps,
    uninstallApps,
    getRelatedFiles,
    setupEventListeners,
  }
})
