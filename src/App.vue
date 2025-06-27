<template>
  <div id="app">
    <RouterView />
    <ErrorMessage />
    <LoadingSpinner />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app'
import ErrorMessage from '@/components/ErrorMessage.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'

const store = useAppStore()

// 在应用启动时从API加载数据
onMounted(async () => {
  try {
    await store.loadData()
  } catch (error) {
    console.error('Failed to load initial data:', error)
  }
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family:
    -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  line-height: 1.6;
  color: #333;
  overflow-x: hidden;
}

#app {
  min-height: 100vh;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.5);
}

/* 选择文本样式 */
::selection {
  background: rgba(102, 126, 234, 0.3);
}

/* 焦点样式 */
button:focus,
input:focus,
textarea:focus,
select:focus {
  outline: 2px solid rgba(102, 126, 234, 0.5);
  outline-offset: 2px;
}

/* 按钮基础样式 */
button {
  font-family: inherit;
  cursor: pointer;
  border: none;
  background: none;
  transition: all 0.2s ease;
}

/* 输入框基础样式 */
input,
textarea,
select {
  font-family: inherit;
  font-size: inherit;
}

/* 链接样式 */
a {
  color: inherit;
  text-decoration: none;
}

/* 图片样式 */
img {
  max-width: 100%;
  height: auto;
}

/* 响应式图片 */
.responsive-img {
  width: 100%;
  height: auto;
  object-fit: cover;
}

/* 工具类 */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

.text-center {
  text-align: center;
}

.text-left {
  text-align: left;
}

.text-right {
  text-align: right;
}

.flex {
  display: flex;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.grid {
  display: grid;
}

.hidden {
  display: none;
}

.visible {
  visibility: visible;
}

.invisible {
  visibility: hidden;
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

@keyframes slideDown {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.fade-in {
  animation: fadeIn 0.3s ease;
}

.slide-up {
  animation: slideUp 0.3s ease;
}

.slide-down {
  animation: slideDown 0.3s ease;
}

/* 打印样式 */
@media print {
  * {
    background: transparent !important;
    color: black !important;
    box-shadow: none !important;
    text-shadow: none !important;
  }

  a,
  a:visited {
    text-decoration: underline;
  }

  a[href]:after {
    content: ' (' attr(href) ')';
  }

  abbr[title]:after {
    content: ' (' attr(title) ')';
  }

  .no-print {
    display: none !important;
  }
}
</style>
