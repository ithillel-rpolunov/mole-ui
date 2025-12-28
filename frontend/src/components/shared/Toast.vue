<script setup>
import { ref, onMounted } from 'vue'

const toasts = ref([])
let nextId = 0

onMounted(() => {
  window.addEventListener('show-toast', (e) => {
    const { message, type = 'info', duration = 5000 } = e.detail
    const id = nextId++

    toasts.value.push({ id, message, type })

    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, duration)
  })
})

function close(id) {
  toasts.value = toasts.value.filter(t => t.id !== id)
}
</script>

<template>
  <div class="toast-container">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      :class="['toast', toast.type]"
      @click="close(toast.id)"
    >
      {{ toast.message }}
    </div>
  </div>
</template>

<style scoped>
.toast-container {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.toast {
  padding: 1rem 1.5rem;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  cursor: pointer;
  max-width: 400px;
  animation: slideIn 0.3s;
}

.toast.error {
  background: #ef4444;
  color: white;
}

.toast.info {
  background: #3b82f6;
  color: white;
}

.toast.success {
  background: #10b981;
  color: white;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}
</style>
