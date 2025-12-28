import { defineStore } from 'pinia'
import { ref } from 'vue'
import { PurgeScanProjects, PurgeExecute } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { handleError } from '../utils/errorHandler'

export const usePurgeStore = defineStore('purge', () => {
  const projects = ref([])
  const loading = ref(false)
  const purging = ref(false)
  const progress = ref(0)
  const progressMessage = ref('')
  const result = ref(null)
  const error = ref(null)

  async function scanProjects(searchPath: string) {
    if (!searchPath) {
      error.value = 'Please provide a search path'
      return
    }

    loading.value = true
    error.value = null
    try {
      const data = await PurgeScanProjects(searchPath)
      projects.value = data || []
    } catch (err) {
      handleError(err, 'Purge scan')
      error.value = 'Failed to scan projects'
    } finally {
      loading.value = false
    }
  }

  async function executePurge(projectPaths: string[]) {
    if (projectPaths.length === 0) {
      error.value = 'Please select at least one project'
      return
    }

    purging.value = true
    progress.value = 0
    result.value = null
    error.value = null

    try {
      await PurgeExecute(projectPaths)
    } catch (err) {
      handleError(err, 'Purge execution')
      error.value = 'Purge failed'
      purging.value = false
    }
  }

  function setupEventListeners() {
    EventsOn('purge:progress', (data) => {
      progress.value = data.percent
      progressMessage.value = data.message
    })

    EventsOn('purge:complete', (data) => {
      purging.value = false
      result.value = data
    })
  }

  return {
    projects,
    loading,
    purging,
    progress,
    progressMessage,
    result,
    error,
    scanProjects,
    executePurge,
    setupEventListeners,
  }
})
