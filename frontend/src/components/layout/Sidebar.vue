<script setup>
defineProps({
  activeTab: String
})

const emit = defineEmits(['change-tab'])

const tabs = [
  { id: 'clean', icon: '🧹', label: 'Clean' },
  { id: 'uninstall', icon: '🗑️', label: 'Uninstall' },
  { id: 'optimize', icon: '⚡', label: 'Optimize' },
  { id: 'analyze', icon: '📊', label: 'Analyze' },
  { id: 'status', icon: '📈', label: 'Status' },
  { id: 'touchid', icon: '🔐', label: 'Touch ID' },
]
</script>

<template>
  <aside class="sidebar">
    <div class="logo">
      <h1>MoleUI</h1>
      <p class="tagline">Deep Clean & Optimize</p>
    </div>

    <nav class="nav">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        class="nav-item"
        :class="{ active: activeTab === tab.id }"
        @click="emit('change-tab', tab.id)"
      >
        <span class="icon">{{ tab.icon }}</span>
        <span class="label">{{ tab.label }}</span>
      </button>
    </nav>

    <div class="about">
      <button
        class="nav-item"
        :class="{ active: activeTab === 'about' }"
        @click="emit('change-tab', 'about')"
      >
        <span class="icon">ℹ️</span>
        <span class="label">About</span>
      </button>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 220px;
  background: #1f2937;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #374151;
}

.logo {
  padding: 2rem 1.5rem;
  border-bottom: 1px solid #374151;
}

.logo h1 {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(135deg, #8b5cf6 0%, #ec4899 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.tagline {
  font-size: 0.75rem;
  color: #9ca3af;
  margin: 0.25rem 0 0 0;
}

.nav {
  flex: 1;
  padding: 1rem 0;
  overflow-y: auto;
}

.nav-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1.5rem;
  background: transparent;
  border: none;
  color: #d1d5db;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
  font-size: 0.9375rem;
}

.nav-item:hover {
  background: #374151;
  color: #f3f4f6;
}

.nav-item.active {
  background: #8b5cf6;
  color: #ffffff;
  font-weight: 600;
}

.icon {
  font-size: 1.25rem;
  width: 1.5rem;
  text-align: center;
}

.label {
  flex: 1;
}

.about {
  border-top: 1px solid #374151;
  padding: 1rem 0;
}
</style>
