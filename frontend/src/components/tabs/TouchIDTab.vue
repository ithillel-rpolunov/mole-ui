<script setup>
import { ref, onMounted } from 'vue'
import { TouchIDGetStatus, TouchIDEnable, TouchIDDisable } from '../../../wailsjs/go/main/App'

const loading = ref(false)
const status = ref(null)
const message = ref(null)
const messageType = ref('success') // 'success' or 'error'

onMounted(async () => {
  await getStatus()
})

async function getStatus() {
  loading.value = true
  try {
    const data = await TouchIDGetStatus()
    status.value = data
    message.value = null
  } catch (error) {
    console.error('Failed to get status:', error)
    message.value = 'Failed to get TouchID status: ' + error
    messageType.value = 'error'
  } finally {
    loading.value = false
  }
}

async function enable() {
  if (!confirm('Enable TouchID for sudo commands? This will modify your PAM configuration.')) {
    return
  }

  loading.value = true
  message.value = null

  try {
    await TouchIDEnable()
    message.value = 'TouchID enabled successfully!'
    messageType.value = 'success'
    await getStatus()
  } catch (error) {
    console.error('Failed to enable:', error)
    message.value = 'Failed to enable TouchID: ' + error
    messageType.value = 'error'
  } finally {
    loading.value = false
  }
}

async function disable() {
  if (!confirm('Disable TouchID for sudo commands?')) {
    return
  }

  loading.value = true
  message.value = null

  try {
    await TouchIDDisable()
    message.value = 'TouchID disabled successfully!'
    messageType.value = 'success'
    await getStatus()
  } catch (error) {
    console.error('Failed to disable:', error)
    message.value = 'Failed to disable TouchID: ' + error
    messageType.value = 'error'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="touchid-tab">
    <h1>Touch ID Configuration</h1>
    <p class="subtitle">Enable Touch ID for sudo commands instead of typing passwords</p>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Loading...</p>
    </div>

    <div v-else class="content">
      <div v-if="message" class="message" :class="messageType">
        {{ message }}
      </div>

      <div v-if="status" class="status-card">
        <div class="status-header">
          <h2>Current Status</h2>
          <div class="status-badge" :class="{ enabled: status.enabled, disabled: !status.enabled }">
            {{ status.enabled ? 'Enabled' : 'Disabled' }}
          </div>
        </div>

        <div class="status-details">
          <div class="detail-item">
            <span class="detail-label">TouchID Available:</span>
            <span class="detail-value">{{ status.available ? 'Yes' : 'No' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">PAM Module:</span>
            <span class="detail-value">{{ status.pamModulePath || 'Not found' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Configuration:</span>
            <span class="detail-value">{{ status.configPath || 'Not configured' }}</span>
          </div>
        </div>

        <div class="info-box">
          <h3>How it works</h3>
          <p>
            When enabled, this feature allows you to use Touch ID instead of typing your password
            for sudo commands in the terminal. This works by configuring the PAM (Pluggable
            Authentication Module) system on macOS.
          </p>
        </div>

        <div class="actions">
          <button
            v-if="!status.enabled && status.available"
            @click="enable"
            class="btn-primary"
            :disabled="loading"
          >
            Enable Touch ID
          </button>
          <button
            v-if="status.enabled"
            @click="disable"
            class="btn-danger"
            :disabled="loading"
          >
            Disable Touch ID
          </button>
          <button
            v-if="!status.available"
            class="btn-secondary"
            disabled
          >
            Touch ID Not Available
          </button>
          <button @click="getStatus" class="btn-secondary" :disabled="loading">
            Refresh Status
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.touchid-tab {
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

.loading {
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

.message {
  padding: 1rem;
  border-radius: 6px;
  margin-bottom: 1.5rem;
  font-weight: 500;
}

.message.success {
  background: #065f46;
  color: #d1fae5;
  border: 1px solid #10b981;
}

.message.error {
  background: #7f1d1d;
  color: #fecaca;
  border: 1px solid #991b1b;
}

.status-card {
  background: #1f2937;
  border: 2px solid #374151;
  border-radius: 8px;
  padding: 2rem;
}

.status-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.status-header h2 {
  margin: 0;
  font-size: 1.5rem;
  color: #f3f4f6;
}

.status-badge {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-weight: 600;
  font-size: 0.875rem;
}

.status-badge.enabled {
  background: #065f46;
  color: #d1fae5;
  border: 1px solid #10b981;
}

.status-badge.disabled {
  background: #7f1d1d;
  color: #fecaca;
  border: 1px solid #991b1b;
}

.status-details {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #111827;
  border-radius: 6px;
}

.detail-label {
  color: #9ca3af;
  font-weight: 500;
}

.detail-value {
  color: #f3f4f6;
  font-family: monospace;
  font-size: 0.875rem;
}

.info-box {
  background: #111827;
  border: 1px solid #374151;
  border-radius: 6px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.info-box h3 {
  margin: 0 0 0.75rem 0;
  font-size: 1rem;
  color: #8b5cf6;
}

.info-box p {
  margin: 0;
  color: #9ca3af;
  line-height: 1.6;
  font-size: 0.875rem;
}

.actions {
  display: flex;
  gap: 1rem;
}

.btn-primary,
.btn-secondary,
.btn-danger {
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

.btn-danger {
  background: linear-gradient(135deg, #dc2626, #991b1b);
  color: white;
  flex: 1;
}

.btn-danger:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.4);
}

.btn-danger:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #374151;
  color: #d1d5db;
}

.btn-secondary:hover:not(:disabled) {
  background: #4b5563;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
