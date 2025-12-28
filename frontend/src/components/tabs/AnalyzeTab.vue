<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime'
import {
  AnalyzeScanDirectory,
  AnalyzeGetLargeFiles,
  AnalyzeDeletePath,
  AnalyzeOpenInFinder
} from '../../../wailsjs/go/main/App'
import { validatePath } from '../../utils/validation'
import { handleError } from '../../utils/errorHandler'

// State
const scanPath = ref('')
const scanning = ref(false)
const scanProgress = ref('')
const scanResult = ref(null)
const largeFiles = ref([])
const loading = ref(false)

// Format bytes to human readable
function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Event listener cleanup function
let unsubscribeProgress = null

// Initialize with home directory
onMounted(() => {
  console.log('[AnalyzeTab] Component mounted')

  // Set default scan path (use a safe default instead of process.env)
  scanPath.value = '/Users'
  console.log('[AnalyzeTab] Default scan path set to:', scanPath.value)

  // Listen for scan progress events
  unsubscribeProgress = EventsOn('analyze:progress', (data) => {
    scanProgress.value = data.message || 'Scanning...'
  })

  console.log('[AnalyzeTab] Event listener registered successfully')
})

// Cleanup event listeners on unmount
onBeforeUnmount(() => {
  if (unsubscribeProgress) {
    EventsOff('analyze:progress')
  }
})

// Scan directory
async function scan() {
  console.log('[AnalyzeTab] Starting scan for path:', scanPath.value)

  const validation = validatePath(scanPath.value, true) // Pass true for scan mode
  if (!validation.valid) {
    console.error('[AnalyzeTab] Path validation failed:', validation.error)
    handleError(new Error(validation.error), 'Path Validation')
    return
  }

  console.log('[AnalyzeTab] Setting scanning to true...')
  scanning.value = true
  scanProgress.value = 'Starting scan...'
  scanResult.value = null
  largeFiles.value = []

  console.log('[AnalyzeTab] State after reset - scanning:', scanning.value, 'scanResult:', scanResult.value)

  try {
    console.log('[AnalyzeTab] Calling AnalyzeScanDirectory...')

    // Scan directory
    const result = await AnalyzeScanDirectory(scanPath.value)
    console.log('[AnalyzeTab] Scan result:', result)

    if (!result) {
      throw new Error('Scan returned no results')
    }

    scanResult.value = result

    // Get large files
    console.log('[AnalyzeTab] Fetching large files...')
    const files = await AnalyzeGetLargeFiles(scanPath.value, 20)
    console.log('[AnalyzeTab] Large files:', files)
    largeFiles.value = files || []

    scanProgress.value = ''
    console.log('[AnalyzeTab] Scan completed successfully')
  } catch (error) {
    console.error('[AnalyzeTab] Scan failed:', error)
    handleError(error, 'Disk Analysis')
    scanResult.value = null
  } finally {
    console.log('[AnalyzeTab] Setting scanning to false...')
    scanning.value = false
    console.log('[AnalyzeTab] Final state - scanning:', scanning.value, 'scanResult:', scanResult.value)
  }
}

// Delete file or directory
async function deleteItem(path) {
  const validation = validatePath(path)
  if (!validation.valid) {
    handleError(new Error(validation.error), 'Path Validation')
    return
  }

  if (!confirm(`Are you sure you want to delete:\n${path}\n\nThis action cannot be undone!`)) {
    return
  }

  loading.value = true
  try {
    await AnalyzeDeletePath(path)

    // Show success
    window.dispatchEvent(new CustomEvent('show-toast', {
      detail: {
        message: 'Item deleted successfully',
        type: 'success'
      }
    }))

    // Rescan
    await scan()
  } catch (error) {
    handleError(error, 'Delete')
  } finally {
    loading.value = false
  }
}

// Open in Finder
async function openInFinder(path) {
  try {
    await AnalyzeOpenInFinder(path)
  } catch (error) {
    handleError(error, 'Open in Finder')
  }
}
</script>

<template>
  <div class="analyze-tab">
    <h1>Disk Space Analyzer</h1>
    <p class="subtitle">Visualize and analyze disk usage</p>

    <!-- DEBUG INFO - Remove after testing -->
    <div style="background: #374151; padding: 1rem; border-radius: 4px; margin-bottom: 1rem; font-size: 0.875rem;">
      <strong>DEBUG:</strong>
      scanning={{ scanning }} |
      scanResult={{ scanResult ? 'SET' : 'NULL' }} |
      largeFiles.length={{ largeFiles.length }}
    </div>

    <!-- Scan Input Section -->
    <div class="scan-section">
      <div class="input-group">
        <input
          v-model="scanPath"
          type="text"
          placeholder="Enter directory path to analyze"
          class="path-input"
          :disabled="scanning"
        />
        <button @click="scan" class="btn-primary" :disabled="scanning">
          {{ scanning ? 'Scanning...' : 'Scan Directory' }}
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="scanning" class="loading">
      <div class="spinner"></div>
      <p class="progress-message">{{ scanProgress }}</p>
    </div>

    <!-- Results Section -->
    <div v-else-if="scanResult" class="results">
      <!-- Summary Cards -->
      <div class="summary-cards">
        <div class="card">
          <div class="card-label">Total Size</div>
          <div class="card-value">{{ formatBytes(scanResult.totalSize || 0) }}</div>
        </div>
        <div class="card">
          <div class="card-label">Total Items</div>
          <div class="card-value">{{ (scanResult.totalItems || 0).toLocaleString() }}</div>
        </div>
        <div class="card">
          <div class="card-label">Scanned Path</div>
          <div class="card-value small">{{ scanPath }}</div>
        </div>
      </div>

      <!-- Large Files Table -->
      <div v-if="largeFiles.length > 0" class="files-section">
        <h2>Largest Files (Top {{ largeFiles.length }})</h2>
        <div class="files-table">
          <div
            v-for="(file, index) in largeFiles"
            :key="index"
            class="file-row"
          >
            <div class="file-info">
              <div class="file-name">{{ file.name }}</div>
              <div class="file-path">{{ file.path }}</div>
            </div>
            <div class="file-size">{{ formatBytes(file.size) }}</div>
            <div class="file-actions">
              <button @click="openInFinder(file.path)" class="btn-action" title="Open in Finder">
                📁
              </button>
              <button
                @click="deleteItem(file.path)"
                class="btn-action btn-danger"
                title="Delete"
                :disabled="loading"
              >
                🗑️
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="no-files">
        <p>No large files found in this directory</p>
      </div>
    </div>

    <!-- Initial State -->
    <div v-else class="initial-state">
      <div class="initial-message">
        <p>Enter a directory path above and click "Scan Directory" to analyze disk usage</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.analyze-tab {
  max-width: 1200px;
}

h1 {
  font-size: 2rem;
  margin: 0 0 0.5rem 0;
}

.subtitle {
  color: #9ca3af;
  margin: 0 0 2rem 0;
}

/* Scan Section */
.scan-section {
  margin-bottom: 2rem;
}

.input-group {
  display: flex;
  gap: 1rem;
}

.path-input {
  flex: 1;
  padding: 0.875rem 1rem;
  background: #1f2937;
  border: 2px solid #374151;
  border-radius: 6px;
  color: #f3f4f6;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.path-input:focus {
  outline: none;
  border-color: #8b5cf6;
}

.path-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Loading State */
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

.progress-message {
  color: #d1d5db;
  margin-top: 1rem;
}

/* Summary Cards */
.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.card {
  background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
  border: 2px solid #374151;
  border-radius: 8px;
  padding: 1.5rem;
  transition: all 0.2s;
}

.card:hover {
  border-color: #8b5cf6;
  transform: translateY(-2px);
}

.card-label {
  color: #9ca3af;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.card-value {
  color: #f3f4f6;
  font-size: 1.75rem;
  font-weight: 700;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.card-value.small {
  font-size: 1rem;
  word-break: break-all;
}

/* Files Section */
.files-section {
  margin-top: 2rem;
}

.files-section h2 {
  font-size: 1.25rem;
  margin: 0 0 1rem 0;
  color: #f3f4f6;
}

.files-table {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.file-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #1f2937;
  border: 2px solid #374151;
  border-radius: 6px;
  transition: all 0.2s;
}

.file-row:hover {
  border-color: #8b5cf6;
  background: #283448;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  color: #f3f4f6;
  font-weight: 600;
  margin-bottom: 0.25rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-path {
  color: #9ca3af;
  font-size: 0.875rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  color: #8b5cf6;
  font-weight: 600;
  font-size: 1rem;
  min-width: 100px;
  text-align: right;
}

.file-actions {
  display: flex;
  gap: 0.5rem;
}

/* Buttons */
.btn-primary {
  padding: 0.875rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  color: white;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-action {
  padding: 0.5rem 0.75rem;
  border: none;
  border-radius: 4px;
  font-size: 1.25rem;
  cursor: pointer;
  background: #374151;
  transition: all 0.2s;
}

.btn-action:hover:not(:disabled) {
  background: #4b5563;
  transform: scale(1.1);
}

.btn-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-danger:hover:not(:disabled) {
  background: #991b1b;
}

/* States */
.initial-state {
  text-align: center;
  padding: 4rem 2rem;
  background: #1f2937;
  border-radius: 8px;
  border: 2px dashed #374151;
}

.initial-message p {
  color: #9ca3af;
  font-size: 1.125rem;
}

.no-files {
  text-align: center;
  padding: 3rem 2rem;
  background: #1f2937;
  border-radius: 8px;
  border: 2px dashed #374151;
}

.no-files p {
  color: #9ca3af;
  font-size: 1rem;
}
</style>
