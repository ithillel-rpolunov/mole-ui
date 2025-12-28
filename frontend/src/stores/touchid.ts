import { defineStore } from 'pinia'
import { ref } from 'vue'
import { TouchIDGetStatus, TouchIDEnable, TouchIDDisable } from '../../wailsjs/go/main/App'

export const useTouchIDStore = defineStore('touchid', () => {
  const status = ref(null)
  const loading = ref(false)
  const error = ref(null)

  async function getStatus() {
    loading.value = true
    error.value = null
    try {
      const data = await TouchIDGetStatus()
      status.value = data
      return data
    } catch (err) {
      error.value = 'Failed to get TouchID status: ' + err
      console.error('Get status failed:', err)
      return null
    } finally {
      loading.value = false
    }
  }

  async function enable() {
    loading.value = true
    error.value = null
    try {
      await TouchIDEnable()
      await getStatus()
    } catch (err) {
      error.value = 'Failed to enable TouchID: ' + err
      console.error('Enable failed:', err)
    } finally {
      loading.value = false
    }
  }

  async function disable() {
    loading.value = true
    error.value = null
    try {
      await TouchIDDisable()
      await getStatus()
    } catch (err) {
      error.value = 'Failed to disable TouchID: ' + err
      console.error('Disable failed:', err)
    } finally {
      loading.value = false
    }
  }

  return {
    status,
    loading,
    error,
    getStatus,
    enable,
    disable,
  }
})
