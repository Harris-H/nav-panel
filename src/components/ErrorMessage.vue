<template>
  <div v-if="store.error" class="error-message">
    <div class="error-content">
      <svg class="error-icon" viewBox="0 0 24 24" width="20" height="20">
        <path
          fill="currentColor"
          d="M13,13H11V7H13M13,17H11V15H13M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z"
        />
      </svg>
      <span class="error-text">{{ store.error }}</span>
      <button @click="store.clearError" class="error-close">
        <svg viewBox="0 0 24 24" width="16" height="16">
          <path
            fill="currentColor"
            d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"
          />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAppStore } from '@/stores/app'

const store = useAppStore()
</script>

<style scoped>
.error-message {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 9999;
  max-width: 400px;
}

.error-content {
  background: #fee;
  border: 1px solid #fcc;
  border-radius: 8px;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.15);
}

.error-icon {
  color: #dc2626;
  flex-shrink: 0;
}

.error-text {
  color: #dc2626;
  font-size: 14px;
  flex: 1;
}

.error-close {
  background: none;
  border: none;
  color: #dc2626;
  cursor: pointer;
  padding: 2px;
  border-radius: 4px;
  flex-shrink: 0;
}

.error-close:hover {
  background: rgba(220, 38, 38, 0.1);
}
</style>
