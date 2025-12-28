<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime'
import {
  StatusGetMetrics,
  StatusStartMonitoring,
  StatusStopMonitoring
} from '../../../wailsjs/go/main/App'
import { handleError } from '../../utils/errorHandler'

// State
const metrics = ref(null)
const monitoring = ref(false)
const loading = ref(false)

// Format helpers
function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatPercent(value) {
  return value.toFixed(1) + '%'
}

function formatRate(bytesPerSec) {
  return formatBytes(bytesPerSec) + '/s'
}

function formatTemp(celsius) {
  return celsius.toFixed(1) + '°C'
}

// Get color based on percentage
function getHealthColor(percent, invert = false) {
  if (invert) {
    // For things like disk space where high is bad
    if (percent >= 90) return '#ef4444' // red
    if (percent >= 75) return '#f59e0b' // yellow
    return '#10b981' // green
  } else {
    // For things like battery where high is good
    if (percent >= 80) return '#10b981' // green
    if (percent >= 50) return '#f59e0b' // yellow
    return '#ef4444' // red
  }
}

// Calculate overall health score (0-100)
function calculateHealthScore() {
  if (!metrics.value) return 0

  let score = 100
  const m = metrics.value

  // CPU usage (max -30 points)
  if (m.cpu && m.cpu.usagePercent) {
    if (m.cpu.usagePercent > 90) score -= 30
    else if (m.cpu.usagePercent > 70) score -= 20
    else if (m.cpu.usagePercent > 50) score -= 10
  }

  // Memory usage (max -25 points)
  if (m.memory && m.memory.usedPercent) {
    if (m.memory.usedPercent > 90) score -= 25
    else if (m.memory.usedPercent > 75) score -= 15
    else if (m.memory.usedPercent > 60) score -= 5
  }

  // Disk usage (max -20 points)
  if (m.disk && m.disk.usedPercent) {
    if (m.disk.usedPercent > 95) score -= 20
    else if (m.disk.usedPercent > 85) score -= 10
    else if (m.disk.usedPercent > 75) score -= 5
  }

  // CPU temperature (max -15 points)
  if (m.cpu && m.cpu.temperature) {
    if (m.cpu.temperature > 85) score -= 15
    else if (m.cpu.temperature > 75) score -= 10
    else if (m.cpu.temperature > 65) score -= 5
  }

  // Battery (max -10 points, only if present)
  if (m.battery && m.battery.percent !== undefined) {
    if (m.battery.percent < 20) score -= 10
    else if (m.battery.percent < 40) score -= 5
  }

  return Math.max(0, Math.min(100, score))
}

function getHealthScoreLabel(score) {
  if (score >= 90) return 'Excellent'
  if (score >= 75) return 'Good'
  if (score >= 60) return 'Fair'
  if (score >= 40) return 'Poor'
  return 'Critical'
}

// Start monitoring
async function startMonitoring() {
  loading.value = true
  try {
    await StatusStartMonitoring(2) // 2 second interval
    monitoring.value = true

    // Get initial metrics
    await fetchMetrics()
  } catch (error) {
    handleError(error, 'Start Monitoring')
  } finally {
    loading.value = false
  }
}

// Stop monitoring
async function stopMonitoring() {
  loading.value = true
  try {
    await StatusStopMonitoring()
    monitoring.value = false
  } catch (error) {
    handleError(error, 'Stop Monitoring')
  } finally {
    loading.value = false
  }
}

// Fetch current metrics
async function fetchMetrics() {
  try {
    const data = await StatusGetMetrics()
    metrics.value = data
  } catch (error) {
    handleError(error, 'Fetch Metrics')
  }
}

// Setup
onMounted(() => {
  // Listen for metric updates
  EventsOn('status:update', (data) => {
    metrics.value = data
  })

  // Auto-start monitoring
  startMonitoring()
})

// Cleanup
onUnmounted(() => {
  EventsOff('status:update')
  if (monitoring.value) {
    StatusStopMonitoring()
  }
})
</script>

<template>
  <div class="status-tab">
    <h1>System Status</h1>
    <p class="subtitle">Real-time system health monitoring</p>

    <!-- Control Buttons -->
    <div class="controls">
      <button
        v-if="!monitoring"
        @click="startMonitoring"
        class="btn-primary"
        :disabled="loading"
      >
        Start Monitoring
      </button>
      <button
        v-else
        @click="stopMonitoring"
        class="btn-stop"
        :disabled="loading"
      >
        Stop Monitoring
      </button>
      <button
        v-if="monitoring"
        @click="fetchMetrics"
        class="btn-secondary"
        :disabled="loading"
      >
        Refresh
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !metrics" class="loading">
      <div class="spinner"></div>
      <p>Loading metrics...</p>
    </div>

    <!-- Metrics Display -->
    <div v-else-if="metrics" class="metrics">
      <!-- Health Score Card -->
      <div class="health-score-card">
        <h2>System Health</h2>
        <div class="health-score-circle">
          <div class="score-value">{{ calculateHealthScore() }}</div>
          <div class="score-label">{{ getHealthScoreLabel(calculateHealthScore()) }}</div>
        </div>
      </div>

      <!-- Metrics Grid -->
      <div class="metrics-grid">
        <!-- CPU Card -->
        <div class="metric-card" v-if="metrics.cpu">
          <div class="card-header">
            <span class="card-icon">🖥️</span>
            <h3>CPU</h3>
          </div>
          <div class="metric-value">
            {{ formatPercent(metrics.cpu.usagePercent || 0) }}
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: (metrics.cpu.usagePercent || 0) + '%',
                backgroundColor: getHealthColor(metrics.cpu.usagePercent || 0, true)
              }"
            ></div>
          </div>
          <div class="metric-details">
            <div v-if="metrics.cpu.temperature">
              <span class="detail-label">Temperature:</span>
              <span class="detail-value">{{ formatTemp(metrics.cpu.temperature) }}</span>
            </div>
            <div v-if="metrics.cpu.cores">
              <span class="detail-label">Cores:</span>
              <span class="detail-value">{{ metrics.cpu.cores }}</span>
            </div>
          </div>
        </div>

        <!-- Memory Card -->
        <div class="metric-card" v-if="metrics.memory">
          <div class="card-header">
            <span class="card-icon">💾</span>
            <h3>Memory</h3>
          </div>
          <div class="metric-value">
            {{ formatPercent(metrics.memory.usedPercent || 0) }}
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: (metrics.memory.usedPercent || 0) + '%',
                backgroundColor: getHealthColor(metrics.memory.usedPercent || 0, true)
              }"
            ></div>
          </div>
          <div class="metric-details">
            <div>
              <span class="detail-label">Used:</span>
              <span class="detail-value">{{ formatBytes(metrics.memory.used || 0) }}</span>
            </div>
            <div>
              <span class="detail-label">Total:</span>
              <span class="detail-value">{{ formatBytes(metrics.memory.total || 0) }}</span>
            </div>
          </div>
        </div>

        <!-- Disk Card -->
        <div class="metric-card" v-if="metrics.disk">
          <div class="card-header">
            <span class="card-icon">💿</span>
            <h3>Disk</h3>
          </div>
          <div class="metric-value">
            {{ formatPercent(metrics.disk.usedPercent || 0) }}
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: (metrics.disk.usedPercent || 0) + '%',
                backgroundColor: getHealthColor(metrics.disk.usedPercent || 0, true)
              }"
            ></div>
          </div>
          <div class="metric-details">
            <div>
              <span class="detail-label">Free:</span>
              <span class="detail-value">{{ formatBytes(metrics.disk.free || 0) }}</span>
            </div>
            <div>
              <span class="detail-label">Total:</span>
              <span class="detail-value">{{ formatBytes(metrics.disk.total || 0) }}</span>
            </div>
          </div>
        </div>

        <!-- Network Card -->
        <div class="metric-card" v-if="metrics.network">
          <div class="card-header">
            <span class="card-icon">🌐</span>
            <h3>Network</h3>
          </div>
          <div class="metric-stats">
            <div class="stat-item">
              <span class="stat-icon">⬇️</span>
              <div class="stat-info">
                <div class="stat-label">Download</div>
                <div class="stat-value">{{ formatRate(metrics.network.downloadRate || 0) }}</div>
              </div>
            </div>
            <div class="stat-item">
              <span class="stat-icon">⬆️</span>
              <div class="stat-info">
                <div class="stat-label">Upload</div>
                <div class="stat-value">{{ formatRate(metrics.network.uploadRate || 0) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Battery Card -->
        <div class="metric-card" v-if="metrics.battery && metrics.battery.percent !== undefined">
          <div class="card-header">
            <span class="card-icon">🔋</span>
            <h3>Battery</h3>
          </div>
          <div class="metric-value">
            {{ formatPercent(metrics.battery.percent || 0) }}
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: (metrics.battery.percent || 0) + '%',
                backgroundColor: getHealthColor(metrics.battery.percent || 0, false)
              }"
            ></div>
          </div>
          <div class="metric-details">
            <div v-if="metrics.battery.status">
              <span class="detail-label">Status:</span>
              <span class="detail-value">{{ metrics.battery.status }}</span>
            </div>
            <div v-if="metrics.battery.health !== undefined">
              <span class="detail-label">Health:</span>
              <span class="detail-value">{{ formatPercent(metrics.battery.health) }}</span>
            </div>
          </div>
        </div>

        <!-- GPU Card -->
        <div class="metric-card" v-if="metrics.gpu && metrics.gpu.usagePercent !== undefined">
          <div class="card-header">
            <span class="card-icon">🎮</span>
            <h3>GPU</h3>
          </div>
          <div class="metric-value">
            {{ formatPercent(metrics.gpu.usagePercent || 0) }}
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: (metrics.gpu.usagePercent || 0) + '%',
                backgroundColor: getHealthColor(metrics.gpu.usagePercent || 0, true)
              }"
            ></div>
          </div>
          <div class="metric-details">
            <div v-if="metrics.gpu.temperature">
              <span class="detail-label">Temperature:</span>
              <span class="detail-value">{{ formatTemp(metrics.gpu.temperature) }}</span>
            </div>
            <div v-if="metrics.gpu.memory">
              <span class="detail-label">Memory:</span>
              <span class="detail-value">{{ formatBytes(metrics.gpu.memory) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Monitoring Indicator -->
      <div v-if="monitoring" class="monitoring-indicator">
        <div class="pulse-dot"></div>
        <span>Live monitoring active (2s interval)</span>
      </div>
    </div>

    <!-- Initial State -->
    <div v-else class="initial-state">
      <p>Click "Start Monitoring" to view real-time system metrics</p>
    </div>
  </div>
</template>

<style scoped>
.status-tab {
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

/* Controls */
.controls {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.btn-primary,
.btn-stop,
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
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}

.btn-stop {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.btn-stop:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
}

.btn-secondary {
  background: #374151;
  color: #d1d5db;
}

.btn-secondary:hover:not(:disabled) {
  background: #4b5563;
}

.btn-primary:disabled,
.btn-stop:disabled,
.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Loading */
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

/* Health Score */
.health-score-card {
  background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
  border: 2px solid #374151;
  border-radius: 12px;
  padding: 2rem;
  margin-bottom: 2rem;
  text-align: center;
}

.health-score-card h2 {
  margin: 0 0 1.5rem 0;
  font-size: 1.5rem;
  color: #f3f4f6;
}

.health-score-circle {
  display: inline-block;
}

.score-value {
  font-size: 4rem;
  font-weight: 700;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
}

.score-label {
  margin-top: 0.5rem;
  font-size: 1.25rem;
  color: #9ca3af;
  font-weight: 600;
}

/* Metrics Grid */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.metric-card {
  background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
  border: 2px solid #374151;
  border-radius: 8px;
  padding: 1.5rem;
  transition: all 0.2s;
}

.metric-card:hover {
  border-color: #8b5cf6;
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.card-icon {
  font-size: 1.5rem;
}

.card-header h3 {
  margin: 0;
  font-size: 1.125rem;
  color: #f3f4f6;
}

.metric-value {
  font-size: 2.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 1rem;
}

/* Progress Bar */
.progress-bar {
  width: 100%;
  height: 8px;
  background: #1f2937;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 1rem;
}

.progress-fill {
  height: 100%;
  transition: width 0.3s, background-color 0.3s;
}

/* Metric Details */
.metric-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.metric-details > div {
  display: flex;
  justify-content: space-between;
  font-size: 0.875rem;
}

.detail-label {
  color: #9ca3af;
}

.detail-value {
  color: #d1d5db;
  font-weight: 600;
}

/* Metric Stats (for Network card) */
.metric-stats {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #1f2937;
  border-radius: 6px;
}

.stat-icon {
  font-size: 1.5rem;
}

.stat-info {
  flex: 1;
}

.stat-label {
  color: #9ca3af;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.stat-value {
  color: #f3f4f6;
  font-size: 1.125rem;
  font-weight: 600;
}

/* Monitoring Indicator */
.monitoring-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 1rem;
  background: #1f2937;
  border-radius: 6px;
  color: #10b981;
  font-size: 0.875rem;
  font-weight: 600;
}

.pulse-dot {
  width: 10px;
  height: 10px;
  background: #10b981;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.5;
    transform: scale(1.2);
  }
}

/* Initial State */
.initial-state {
  text-align: center;
  padding: 4rem 2rem;
  background: #1f2937;
  border-radius: 8px;
  border: 2px dashed #374151;
}

.initial-state p {
  color: #9ca3af;
  font-size: 1.125rem;
}
</style>
