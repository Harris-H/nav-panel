<template>
  <div class="search-container" v-if="store.settings.search.enabled">
    <div class="search-box">
      <div class="search-engine-selector" @click="showEngineSelector = !showEngineSelector">
        <img
          v-if="currentEngine?.icon"
          :src="currentEngine.icon"
          :alt="currentEngine.name"
          class="engine-icon"
        />
        <div v-else class="engine-icon engine-icon-fallback">
          {{ currentEngine?.name?.charAt(0) || 'G' }}
        </div>

        <svg
          class="dropdown-icon"
          :class="{ rotated: showEngineSelector }"
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <polyline points="6,9 12,15 18,9"></polyline>
        </svg>
      </div>

      <input
        v-model="searchInput"
        type="text"
        :placeholder="currentEngine?.placeholder || '搜索...'"
        class="search-input"
        @keyup.enter="handleSearch"
        @focus="showSiteSearch = false"
        @input="handleSearchInput"
      />

      <button class="search-button" @click="handleSearch">
        <svg
          width="18"
          height="18"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.35-4.35"></path>
        </svg>
      </button>
    </div>

    <!-- 横向搜索引擎选择器 - 改为非绝对定位 -->
    <Transition name="slide-down">
      <div v-if="showEngineSelector" class="engines-horizontal-selector">
        <div class="engines-title">选择搜索引擎</div>
        <div class="engines-grid">
          <div
            v-for="engine in store.settings.search.engines"
            :key="engine.id"
            class="engine-item"
            :class="{ active: currentEngine?.id === engine.id }"
            @click="selectSearchEngine(engine.id)"
            :title="engine.name"
          >
            <img
              v-if="engine.icon"
              :src="engine.icon"
              :alt="engine.name"
              class="engine-item-icon"
            />
            <div v-else class="engine-item-icon engine-icon-fallback">
              {{ engine.name.charAt(0) }}
            </div>
            <span class="engine-item-name">{{ engine.name }}</span>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 网站搜索建议 - 保持绝对定位 -->
    <Transition name="fade">
      <div v-if="showSiteSearch && filteredSitesForSearch.length > 0" class="site-suggestions">
        <div class="suggestion-header">网站建议</div>
        <div
          v-for="site in filteredSitesForSearch.slice(0, 5)"
          :key="site.id"
          class="suggestion-item"
          @click="openSite(site.url)"
        >
          <div class="suggestion-icon">
            <img v-if="site.icon" :src="site.icon" :alt="site.name" />
            <span v-else>{{ site.name.charAt(0).toUpperCase() }}</span>
          </div>
          <div class="suggestion-content">
            <div class="suggestion-name">{{ site.name }}</div>
            <div class="suggestion-url">{{ site.url }}</div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useAppStore } from '@/stores/app'

const store = useAppStore()
const searchInput = ref('')
const showEngineSelector = ref(false)
const showSiteSearch = ref(false)
// 本地搜索查询，不影响主页面的网站显示
const localSearchQuery = ref('')

// 当前搜索引擎 - 使用更健壮的逻辑
const currentEngine = computed(() => {
  // 优先使用 store.currentSearchEngine
  if (store.currentSearchEngine) {
    // 验证当前搜索引擎是否在引擎列表中
    const isValidEngine = store.settings.search.engines.some(
      (engine) => engine.id === store.currentSearchEngine?.id,
    )
    if (isValidEngine) {
      return store.currentSearchEngine
    }
  }

  // 如果当前搜索引擎无效，使用默认搜索引擎
  if (store.defaultSearchEngine) {
    return store.defaultSearchEngine
  }

  // 最后的备选方案：返回第一个可用的搜索引擎
  if (store.settings.search.engines.length > 0) {
    return store.settings.search.engines[0]
  }

  // 如果都没有，返回内置默认搜索引擎
  return store.defaultSearchEngines[0]
})

const openSite = (url: string) => {
  window.open(url, '_blank')
}

// 处理搜索输入
const handleSearchInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = target.value

  // 只更新本地搜索查询，不影响主页面
  localSearchQuery.value = value

  // 显示/隐藏网站搜索建议
  showSiteSearch.value = value.length > 0
}

// 执行搜索
const handleSearch = () => {
  if (searchInput.value.trim()) {
    store.performSearch(searchInput.value, currentEngine.value?.id)
    searchInput.value = ''
    localSearchQuery.value = ''
    showSiteSearch.value = false
  }
}

// 选择搜索引擎
const selectSearchEngine = (engineId: string) => {
  console.log('🎯 SearchBox: 选择搜索引擎', engineId)
  console.log(
    '🎯 可用引擎:',
    store.settings.search.engines.map((e) => ({ id: e.id, name: e.name })),
  )

  store.setCurrentSearchEngine(engineId)
  showEngineSelector.value = false
}

// 过滤网站用于搜索建议 - 使用本地搜索查询
const filteredSitesForSearch = computed(() => {
  if (!localSearchQuery.value || !showSiteSearch.value) {
    return []
  }

  return store.sites.filter(
    (site) =>
      site.name.toLowerCase().includes(localSearchQuery.value.toLowerCase()) ||
      site.url.toLowerCase().includes(localSearchQuery.value.toLowerCase()) ||
      (site.description &&
        site.description.toLowerCase().includes(localSearchQuery.value.toLowerCase())),
  )
})

// 点击外部关闭下拉菜单
const handleClickOutside = (event: Event) => {
  const target = event.target as Element
  if (!target.closest('.search-engine-selector')) {
    showEngineSelector.value = false
  }
  if (!target.closest('.search-container')) {
    showSiteSearch.value = false
  }
}

// 添加调试 watch
watch(
  () => showEngineSelector.value,
  (isOpen) => {
    if (isOpen) {
      console.log('🔍 SearchBox: 下拉菜单打开')
      console.log('📊 store.settings.search.engines:', store.settings.search.engines)
      console.log('📊 engines 数量:', store.settings.search.engines?.length || 0)
      console.log('📊 当前搜索引擎:', store.currentSearchEngine)
      console.log('📊 默认搜索引擎:', store.defaultSearchEngine)
    }
  },
)

onMounted(() => {
  // 添加全局点击事件监听
  document.addEventListener('click', handleClickOutside)

  // 确保组件挂载时有正确的当前搜索引擎
  if (!store.currentSearchEngine && store.settings.search.engines.length > 0) {
    const defaultEngine =
      store.settings.search.engines.find(
        (engine) => engine.id === store.settings.search.defaultEngineId,
      ) || store.settings.search.engines[0]

    store.setCurrentSearchEngine(defaultEngine.id)
  }
})

onUnmounted(() => {
  // 清理事件监听
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.search-container {
  max-width: 500px;
  margin: 0 auto 40px;
  position: relative;
}

.search-box {
  display: flex;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.search-box:hover {
  background: rgba(255, 255, 255, 1);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.search-engine-selector {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 12px;
  cursor: pointer;
  border-right: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(0, 0, 0, 0.02);
  transition: all 0.3s ease;
  min-width: 60px;
}

.search-engine-selector:hover {
  background: rgba(0, 0, 0, 0.05);
}

.engine-icon {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  object-fit: cover;
}

.engine-icon-fallback {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
}

.dropdown-icon {
  color: #666;
  transition: transform 0.3s ease;
  margin-left: 4px;
}

.dropdown-icon.rotated {
  transform: rotate(180deg);
}

/* 横向搜索引擎选择器 - 改为占用实际空间 */
.engines-horizontal-selector {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-top: 12px;
  padding: 16px;
  overflow: hidden;
}

/* Vue Transition 动画 */
.slide-down-enter-active {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
}

.slide-down-leave-active {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-20px);
  max-height: 0;
  padding-top: 0;
  padding-bottom: 0;
  margin-top: 0;
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-20px);
  max-height: 0;
  padding-top: 0;
  padding-bottom: 0;
  margin-top: 0;
}

.slide-down-enter-to,
.slide-down-leave-from {
  opacity: 1;
  transform: translateY(0);
  max-height: 200px;
}

.engines-title {
  font-size: 12px;
  font-weight: 600;
  color: #666;
  margin-bottom: 12px;
  text-align: center;
}

.engines-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(80px, 1fr));
  gap: 12px;
  max-width: 400px;
  margin: 0 auto;
}

.engine-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  cursor: pointer;
  border-radius: 12px;
  transition: all 0.2s ease;
  background: rgba(255, 255, 255, 0.5);
  border: 2px solid transparent;
}

.engine-item:hover {
  background: rgba(102, 126, 234, 0.1);
  transform: translateY(-2px);
}

.engine-item.active {
  background: rgba(102, 126, 234, 0.15);
  border-color: rgba(102, 126, 234, 0.4);
  transform: translateY(-2px);
}

.engine-item-icon {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  object-fit: cover;
  margin-bottom: 6px;
}

.engine-item-name {
  font-size: 11px;
  font-weight: 500;
  color: #333;
  text-align: center;
  line-height: 1.2;
}

.search-input {
  flex: 1;
  padding: 16px 20px;
  border: none;
  font-size: 16px;
  outline: none;
  background: transparent;
  color: #333;
}

.search-input::placeholder {
  color: #999;
}

.search-button {
  padding: 12px 16px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: #666;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-button:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

/* 网站搜索建议 - 保持绝对定位，因为它是临时提示 */
.site-suggestions {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  z-index: 999;
  margin-top: 8px;
  overflow: hidden;
}

/* 淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.suggestion-header {
  padding: 12px 20px 8px;
  font-size: 12px;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.suggestion-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  border-bottom: 1px solid rgba(0, 0, 0, 0.03);
}

.suggestion-item:last-child {
  border-bottom: none;
}

.suggestion-item:hover {
  background: rgba(102, 126, 234, 0.05);
}

.suggestion-icon {
  width: 32px;
  height: 32px;
  margin-right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 8px;
  overflow: hidden;
}

.suggestion-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.suggestion-icon span {
  font-size: 14px;
  font-weight: 600;
  color: #667eea;
}

.suggestion-content {
  flex: 1;
  min-width: 0;
}

.suggestion-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.suggestion-url {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .search-container {
    padding: 0 16px;
    margin-bottom: 24px;
  }

  .search-engine-selector {
    min-width: 50px;
    padding: 6px 10px;
  }

  .engines-grid {
    grid-template-columns: repeat(auto-fit, minmax(70px, 1fr));
    gap: 10px;
  }

  .engine-item {
    padding: 10px 6px;
  }

  .engine-item-icon {
    width: 20px;
    height: 20px;
  }

  .engine-item-name {
    font-size: 10px;
  }

  .search-input {
    padding: 14px 16px;
    font-size: 16px;
  }

  .search-button {
    padding: 10px 14px;
  }

  .slide-down-enter-to,
  .slide-down-leave-from {
    max-height: 180px;
  }
}
</style>
