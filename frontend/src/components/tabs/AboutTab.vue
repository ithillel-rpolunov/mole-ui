<script setup>
import { ref } from 'vue'

// App information
const appInfo = {
  name: 'MoleUI - Deep Clean & Optimize',
  version: '1.0.0',
  description: 'Visual interface for the Mole CLI - Advanced macOS system maintenance and optimization tool',
}

// Features
const features = [
  {
    icon: '🧹',
    title: 'System Cleanup',
    description: 'Remove cache, logs, and temporary files'
  },
  {
    icon: '🗑️',
    title: 'Application Management',
    description: 'Uninstall and optimize applications'
  },
  {
    icon: '📊',
    title: 'Disk Space Analysis',
    description: 'Visualize and analyze disk usage'
  },
  {
    icon: '📈',
    title: 'System Monitoring',
    description: 'Real-time CPU, Memory, Disk, and Network metrics'
  },
  {
    icon: '🔐',
    title: 'Touch ID Management',
    description: 'Manage Touch ID authentication settings'
  }
]

// Technology stack
const techStack = [
  { label: 'Frontend', value: 'Vue 3 + Vite' },
  { label: 'Backend', value: 'Go + Wails v2' },
  { label: 'Platform', value: 'macOS (Apple Silicon & Intel)' }
]

// Credits
const credits = [
  { label: 'UI Developer', value: 'Claude Code' },
  { label: 'Framework', value: 'Wails' },
  { label: 'License', value: 'MIT' }
]

// Core project attribution
const coreAttribution = {
  title: 'Based on Mole CLI',
  description: 'This application uses the Mole project as its core engine, providing a visual interface to the powerful command-line tool.',
  repoUrl: 'https://github.com/tw93/Mole',
  author: 'tw93'
}
</script>

<template>
  <div class="about-tab">
    <h1>About MoleUI</h1>
    <p class="subtitle">Visual interface for macOS system maintenance and optimization</p>

    <!-- App Info Card -->
    <div class="info-card">
      <div class="app-header">
        <div class="app-icon">🦡</div>
        <div class="app-details">
          <h2>{{ appInfo.name }}</h2>
          <div class="version">Version {{ appInfo.version }}</div>
        </div>
      </div>
      <p class="app-description">{{ appInfo.description }}</p>
    </div>

    <!-- Features Section -->
    <div class="section">
      <h3 class="section-title">Features</h3>
      <div class="features-grid">
        <div v-for="feature in features" :key="feature.title" class="feature-card">
          <div class="feature-icon">{{ feature.icon }}</div>
          <h4>{{ feature.title }}</h4>
          <p>{{ feature.description }}</p>
        </div>
      </div>
    </div>

    <!-- Technology Stack Section -->
    <div class="section">
      <h3 class="section-title">Technology Stack</h3>
      <div class="tech-list">
        <div v-for="tech in techStack" :key="tech.label" class="tech-item">
          <span class="tech-label">{{ tech.label }}:</span>
          <span class="tech-value">{{ tech.value }}</span>
        </div>
      </div>
    </div>

    <!-- Credits Section -->
    <div class="section">
      <h3 class="section-title">Credits</h3>
      <div class="credits-list">
        <div v-for="credit in credits" :key="credit.label" class="credit-item">
          <span class="credit-label">{{ credit.label }}:</span>
          <span class="credit-value">{{ credit.value }}</span>
        </div>
      </div>
    </div>

    <!-- Core Attribution Section -->
    <div class="section">
      <h3 class="section-title">{{ coreAttribution.title }}</h3>
      <div class="attribution-card">
        <p class="attribution-description">{{ coreAttribution.description }}</p>
        <div class="attribution-details">
          <div class="attribution-item">
            <span class="attribution-label">Original Project:</span>
            <a :href="coreAttribution.repoUrl" target="_blank" rel="noopener noreferrer" class="attribution-link">
              {{ coreAttribution.repoUrl }}
            </a>
          </div>
          <div class="attribution-item">
            <span class="attribution-label">Author:</span>
            <span class="attribution-value">{{ coreAttribution.author }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="footer">
      <p>Built with care for macOS users</p>
      <p class="copyright">&copy; 2025 MoleUI. All rights reserved.</p>
    </div>
  </div>
</template>

<style scoped>
.about-tab {
  max-width: 900px;
  margin: 0 auto;
}

h1 {
  font-size: 2rem;
  margin: 0 0 0.5rem 0;
}

.subtitle {
  color: #9ca3af;
  margin: 0 0 2rem 0;
}

/* App Info Card */
.info-card {
  background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
  border: 2px solid #374151;
  border-radius: 12px;
  padding: 2rem;
  margin-bottom: 2rem;
  transition: all 0.2s;
}

.info-card:hover {
  border-color: #8b5cf6;
  transform: translateY(-2px);
}

.app-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.app-icon {
  font-size: 4rem;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  border-radius: 16px;
}

.app-details h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.75rem;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.version {
  color: #9ca3af;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.app-description {
  color: #d1d5db;
  font-size: 1.125rem;
  line-height: 1.6;
  margin: 0;
}

/* Section */
.section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 1rem 0;
  color: #f3f4f6;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #374151;
}

/* Features Grid */
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.feature-card {
  background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
  border: 2px solid #374151;
  border-radius: 8px;
  padding: 1.5rem;
  transition: all 0.2s;
}

.feature-card:hover {
  border-color: #8b5cf6;
  transform: translateY(-2px);
}

.feature-icon {
  font-size: 2rem;
  margin-bottom: 0.75rem;
}

.feature-card h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  color: #f3f4f6;
}

.feature-card p {
  margin: 0;
  color: #9ca3af;
  font-size: 0.875rem;
  line-height: 1.5;
}

/* Tech List */
.tech-list {
  background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
  border: 2px solid #374151;
  border-radius: 8px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.tech-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.tech-label {
  color: #9ca3af;
  font-size: 0.875rem;
  font-weight: 600;
  min-width: 100px;
}

.tech-value {
  color: #f3f4f6;
  font-size: 1rem;
  font-weight: 600;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Credits List */
.credits-list {
  background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
  border: 2px solid #374151;
  border-radius: 8px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.credit-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.credit-label {
  color: #9ca3af;
  font-size: 0.875rem;
  font-weight: 600;
  min-width: 120px;
}

.credit-value {
  color: #d1d5db;
  font-size: 1rem;
  font-weight: 600;
}

/* Attribution Card */
.attribution-card {
  background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
  border: 2px solid #374151;
  border-radius: 8px;
  padding: 1.5rem;
  transition: all 0.2s;
}

.attribution-card:hover {
  border-color: #8b5cf6;
}

.attribution-description {
  color: #d1d5db;
  font-size: 1rem;
  line-height: 1.6;
  margin: 0 0 1.5rem 0;
}

.attribution-details {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.attribution-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.attribution-label {
  color: #9ca3af;
  font-size: 0.875rem;
  font-weight: 600;
  min-width: 120px;
}

.attribution-value {
  color: #d1d5db;
  font-size: 1rem;
  font-weight: 600;
}

.attribution-link {
  color: #8b5cf6;
  font-size: 0.875rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s;
  word-break: break-all;
}

.attribution-link:hover {
  color: #ec4899;
  text-decoration: underline;
}

/* Footer */
.footer {
  text-align: center;
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 1px solid #374151;
}

.footer p {
  color: #9ca3af;
  margin: 0.5rem 0;
}

.copyright {
  font-size: 0.875rem;
}
</style>
