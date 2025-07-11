import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { ElNotification, ElMessageBox } from 'element-plus'
import type { Website, AppSettings, SearchEngine, Group, GroupWithWebsites } from '@/types'
import { api } from '@/utils/api'

export const useAppStore = defineStore('app', () => {
  // 状态
  const sites = ref<Website[]>([])
  const groups = ref<Group[]>([])
  const groupsWithWebsites = ref<GroupWithWebsites[]>([])
  const searchQuery = ref('')
  const isAddCardModalOpen = ref(false)
  const isSettingsModalOpen = ref(false)
  const isGroupModalOpen = ref(false)
  const editingCard = ref<Website | null>(null)
  const editingGroup = ref<Group | null>(null)
  const currentSearchEngine = ref<SearchEngine | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 默认搜索引擎配置
  const defaultSearchEngines: SearchEngine[] = [
    {
      id: 'google',
      name: 'Google',
      url: 'https://www.google.com/search?q={}',
      icon: 'https://www.google.com/favicon.ico',
      placeholder: '使用 Google 搜索...',
      isDefault: true,
    },
    {
      id: 'baidu',
      name: '百度',
      url: 'https://www.baidu.com/s?wd={}',
      icon: 'https://www.baidu.com/favicon.ico',
      placeholder: '使用百度搜索...',
    },
    {
      id: 'bing',
      name: 'Bing',
      url: 'https://www.bing.com/search?q={}',
      icon: 'https://www.bing.com/favicon.ico',
      placeholder: '使用 Bing 搜索...',
    },
    {
      id: 'github',
      name: 'GitHub',
      url: 'https://github.com/search?q={}',
      icon: 'https://github.com/favicon.ico',
      placeholder: '在 GitHub 中搜索...',
    },
  ]

  const settings = ref<AppSettings>({
    theme: 'light',
    layout: {
      columns: 6,
      cardSize: 'medium',
      showLabels: true,
      gap: 20,
    },
    background: {
      type: 'gradient',
      value: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    },
    cardStyle: {
      borderRadius: 12,
      opacity: 0.9,
      shadow: true,
    },
    search: {
      enabled: true,
      engines: defaultSearchEngines,
      defaultEngineId: 'google',
      openInNewTab: true,
    },
  })

  // 计算属性
  const filteredSites = computed(() => {
    if (!sites.value || !Array.isArray(sites.value)) {
      return []
    }

    if (!searchQuery.value) {
      return sites.value
    }

    return sites.value.filter(
      (site) =>
        site.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        site.url.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        (site.description &&
          site.description.toLowerCase().includes(searchQuery.value.toLowerCase())),
    )
  })

  const defaultSearchEngine = computed(() => {
    return (
      settings.value.search.engines.find(
        (engine) => engine.id === settings.value.search.defaultEngineId,
      ) || settings.value.search.engines[0]
    )
  })

  // 错误处理
  const handleError = (err: unknown, message: string) => {
    console.error(message, err)
    error.value = err instanceof Error ? err.message : message
  }

  const clearError = () => {
    error.value = null
  }

  // 数据加载方法
  const loadData = async () => {
    try {
      loading.value = true
      clearError()
      console.log('Loading data from API...')

      // 并行加载所有数据
      const [websitesData, settingsData, searchEnginesData, groupsData, groupsWithWebsitesData] =
        await Promise.all([
          api.getWebsites().catch((err) => {
            console.warn('Failed to load websites:', err)
            return []
          }),
          api.getSettings().catch((err) => {
            console.warn('Failed to load settings:', err)
            return null
          }),
          api.getSearchEngines().catch((err) => {
            console.warn('Failed to load search engines:', err)
            return []
          }),
          api.getGroups().catch((err) => {
            console.warn('Failed to load groups:', err)
            return []
          }),
          api.getGroupsWithWebsites().catch((err) => {
            console.warn('Failed to load groups with websites:', err)
            return []
          }),
        ])

      // 设置网站数据
      sites.value = Array.isArray(websitesData) ? websitesData : []

      // 设置分组数据
      groups.value = Array.isArray(groupsData) ? groupsData : []
      groupsWithWebsites.value = Array.isArray(groupsWithWebsitesData) ? groupsWithWebsitesData : []

      // 处理搜索引擎数据
      let searchEngines = Array.isArray(searchEnginesData) ? searchEnginesData : []

      // 如果后端没有搜索引擎数据，使用默认数据并保存到后端
      if (searchEngines.length === 0) {
        console.log('No search engines found in backend, creating defaults...')
        try {
          const createdEngines = []
          for (const engine of defaultSearchEngines) {
            try {
              const created = await api.createSearchEngine(engine)
              createdEngines.push(created)
            } catch (err) {
              console.warn('Failed to create default search engine:', engine.name, err)
              // 如果创建失败，使用本地默认数据
              createdEngines.push(engine)
            }
          }
          searchEngines = createdEngines
        } catch {
          console.warn('Failed to initialize default search engines, using local defaults')
          searchEngines = defaultSearchEngines
        }
      }

      // 合并设置数据
      if (settingsData) {
        settings.value = {
          ...settings.value,
          ...settingsData,
          search: {
            ...settings.value.search,
            ...settingsData.search,
            engines: searchEngines,
          },
        }
      } else {
        settings.value.search.engines = searchEngines
      }

      // 确保默认搜索引擎 ID 有效
      const validEngineIds = settings.value.search.engines.map((e) => e.id)
      if (!validEngineIds.includes(settings.value.search.defaultEngineId)) {
        settings.value.search.defaultEngineId = settings.value.search.engines[0]?.id || 'google'
      }

      // 设置当前搜索引擎
      currentSearchEngine.value = defaultSearchEngine.value

      console.log('Loaded sites:', sites.value.length)
      console.log('Loaded groups:', groups.value.length)
      console.log('Loaded search engines:', settings.value.search.engines.length)
      console.log('Loaded settings:', settings.value)
      console.log('Current search engine:', currentSearchEngine.value)
    } catch (err) {
      handleError(err, 'Error loading data from API')
      // 确保即使出错也有基本的搜索引擎配置
      settings.value.search.engines = defaultSearchEngines
      currentSearchEngine.value = defaultSearchEngines[0]
    } finally {
      loading.value = false
    }
  }

  // 网站管理方法
  const addSite = async (site: Omit<Website, 'id' | 'createdAt' | 'updatedAt'>) => {
    try {
      loading.value = true
      clearError()

      const newSite = await api.createWebsite(site)
      sites.value.push(newSite)

      // 更新 groupsWithWebsites 中的网站数据
      updateWebsiteInGroups(newSite)

      showSuccess('网站添加成功')
      console.log('Site added successfully:', newSite)
    } catch (err) {
      handleError(err, 'Error adding site')
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateSite = async (updatedSite: Website) => {
    try {
      loading.value = true
      clearError()

      const updated = await api.updateWebsite(updatedSite.id, updatedSite)
      const index = sites.value.findIndex((site) => site.id === updatedSite.id)
      if (index !== -1) {
        sites.value[index] = updated
      }

      // 更新 groupsWithWebsites 中的网站数据
      updateWebsiteInGroups(updated)

      showSuccess('网站更新成功')
      console.log('Site updated successfully:', updated)
    } catch (err) {
      handleError(err, 'Error updating site')
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteSite = async (id: string) => {
    try {
      loading.value = true
      clearError()

      await api.deleteWebsite(id)
      const index = sites.value.findIndex((site) => site.id === id)
      if (index !== -1) {
        sites.value.splice(index, 1)
      }

      // 从所有分组中移除该网站
      groupsWithWebsites.value.forEach((group) => {
        group.websites = group.websites.filter((site) => site.id !== id)
      })

      showSuccess('网站删除成功')
      console.log('Site deleted successfully:', id)
    } catch (err) {
      handleError(err, 'Error deleting site')
      throw err
    } finally {
      loading.value = false
    }
  }

  const reorderSites = async (fromIndex: number, toIndex: number) => {
    if (fromIndex === toIndex) return

    try {
      // 先在前端更新，提供即时反馈
      const sitesArray = [...sites.value]
      const [movedSite] = sitesArray.splice(fromIndex, 1)
      sitesArray.splice(toIndex, 0, movedSite)

      sites.value = sitesArray

      // 发送新的顺序到后端
      const websiteIds = sitesArray.map((site) => site.id)
      await api.reorderWebsites(websiteIds)

      console.log('Sites reordered successfully:', { fromIndex, toIndex })
    } catch (err) {
      console.error('Failed to reorder sites:', err)

      // 如果后端更新失败，重新加载数据以恢复正确的顺序
      await loadData()

      handleError(err, 'Error reordering sites')
    }
  }

  // 模态框控制
  const openAddCardModal = () => {
    editingCard.value = null
    isAddCardModalOpen.value = true
  }

  const openEditCardModal = (site: Website) => {
    editingCard.value = site
    isAddCardModalOpen.value = true
  }

  const closeAddCardModal = () => {
    editingCard.value = null
    isAddCardModalOpen.value = false
  }

  const openSettingsModal = () => {
    isSettingsModalOpen.value = true
  }

  const closeSettingsModal = () => {
    isSettingsModalOpen.value = false
  }

  // 设置管理方法
  const updateSettings = async (newSettings: Partial<AppSettings>) => {
    try {
      loading.value = true
      clearError()

      const updatedSettings = await api.updateSettings(newSettings)

      // 更新本地设置，确保搜索引擎数据正确同步
      settings.value = {
        ...settings.value,
        ...updatedSettings,
      }

      // 确保搜索引擎数组存在
      if (!settings.value.search?.engines || !Array.isArray(settings.value.search.engines)) {
        // 如果后端返回的搜索引擎数据无效，从API重新获取
        try {
          const searchEngines = await api.getSearchEngines()
          settings.value.search = {
            ...settings.value.search,
            engines: searchEngines.length > 0 ? searchEngines : defaultSearchEngines,
          }
        } catch {
          console.warn('Failed to load search engines, using defaults')
          settings.value.search = {
            ...settings.value.search,
            engines: defaultSearchEngines,
          }
        }
      }

      // 确保默认搜索引擎 ID 有效
      const validEngineIds = settings.value.search.engines.map((e) => e.id)
      if (!validEngineIds.includes(settings.value.search.defaultEngineId)) {
        settings.value.search.defaultEngineId = settings.value.search.engines[0]?.id || 'google'
      }

      // 只有在当前搜索引擎无效时才更新
      if (currentSearchEngine.value) {
        // 检查当前搜索引擎是否还在引擎列表中
        const currentEngineStillExists = settings.value.search.engines.find(
          (engine) => engine.id === currentSearchEngine.value?.id,
        )

        if (currentEngineStillExists) {
          // 如果当前搜索引擎仍然存在，更新为最新数据
          currentSearchEngine.value = currentEngineStillExists
        } else {
          // 如果当前搜索引擎已被删除，切换到默认搜索引擎
          currentSearchEngine.value = defaultSearchEngine.value
        }
      } else {
        // 如果没有当前搜索引擎，设置为默认搜索引擎
        currentSearchEngine.value = defaultSearchEngine.value
      }
    } catch (err) {
      handleError(err, 'Error updating settings')
      throw err
    } finally {
      loading.value = false
    }
  }

  const resetSettings = async () => {
    try {
      loading.value = true
      clearError()

      // 重置为默认设置
      const defaultSettings: Partial<AppSettings> = {
        theme: 'light',
        layout: {
          columns: 6,
          cardSize: 'medium',
          showLabels: true,
          gap: 20,
        },
        background: {
          type: 'gradient',
          value: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        },
        cardStyle: {
          borderRadius: 12,
          opacity: 0.9,
          shadow: true,
        },
        search: {
          enabled: true,
          engines: defaultSearchEngines,
          defaultEngineId: 'google',
          openInNewTab: true,
        },
      }

      await updateSettings(defaultSettings)

      console.log('Settings reset to default')
    } catch (err) {
      handleError(err, 'Error resetting settings')
      throw err
    }
  }

  // 搜索功能
  const performSearch = (query: string, engineId?: string) => {
    if (!query.trim()) return

    const engine = engineId
      ? settings.value.search.engines.find((e) => e.id === engineId)
      : currentSearchEngine.value

    if (!engine) {
      console.error('No search engine found')
      return
    }

    const searchUrl = engine.url.replace('{}', encodeURIComponent(query))

    if (settings.value.search.openInNewTab) {
      window.open(searchUrl, '_blank')
    } else {
      window.location.href = searchUrl
    }
  }

  // 临时切换当前搜索引擎（仅UI状态，不保存到后端）
  const setCurrentSearchEngine = (engineId: string) => {
    const engine = settings.value.search.engines.find((e) => e.id === engineId)
    if (engine) {
      currentSearchEngine.value = engine
      console.log('Current search engine switched to:', engine.name)
    }
  }

  // 设置默认搜索引擎（仅在设置页面使用，会保存到后端）
  const setDefaultSearchEngine = async (engineId: string) => {
    const engine = settings.value.search.engines.find((e) => e.id === engineId)
    if (engine) {
      // 更新默认搜索引擎设置
      await updateSettings({
        search: {
          ...settings.value.search,
          defaultEngineId: engineId,
        },
      })

      // 同时更新当前搜索引擎
      currentSearchEngine.value = engine
      console.log('Default search engine saved:', engine.name)
    }
  }

  // 搜索引擎管理方法
  const addSearchEngine = async (engine: Omit<SearchEngine, 'id'>) => {
    try {
      loading.value = true
      clearError()

      const newEngine = await api.createSearchEngine(engine)

      // 更新本地设置中的搜索引擎列表
      settings.value.search.engines.push(newEngine)

      // 如果是第一个搜索引擎，设为默认
      if (settings.value.search.engines.length === 1) {
        settings.value.search.defaultEngineId = newEngine.id
        currentSearchEngine.value = newEngine
      }

      return newEngine
    } catch (err) {
      handleError(err, 'Error adding search engine')
      throw err
    } finally {
      loading.value = false
    }
  }

  // 添加带图片上传的搜索引擎
  const addSearchEngineWithIcon = async (
    engine: { name: string; url: string; placeholder?: string; isDefault?: boolean },
    iconFile?: File,
  ) => {
    try {
      loading.value = true
      clearError()

      const newEngine = await api.createSearchEngineWithIcon(engine, iconFile)

      // 更新本地设置中的搜索引擎列表
      settings.value.search.engines.push(newEngine)

      // 如果是第一个搜索引擎，设为默认
      if (settings.value.search.engines.length === 1) {
        settings.value.search.defaultEngineId = newEngine.id
        currentSearchEngine.value = newEngine
      }

      console.log('Search engine with icon added successfully:', newEngine)
      return newEngine
    } catch (err) {
      handleError(err, 'Error adding search engine with icon')
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateSearchEngine = async (id: string, updates: Partial<SearchEngine>) => {
    try {
      loading.value = true
      clearError()

      const updatedEngine = await api.updateSearchEngine(id, updates)

      // 更新本地设置中的搜索引擎
      const index = settings.value.search.engines.findIndex((engine) => engine.id === id)
      if (index !== -1) {
        settings.value.search.engines[index] = updatedEngine
      }

      // 如果更新的是当前搜索引擎，同步更新
      if (currentSearchEngine.value?.id === id) {
        currentSearchEngine.value = updatedEngine
      }

      console.log('Search engine updated successfully:', updatedEngine)
      return updatedEngine
    } catch (err) {
      handleError(err, 'Error updating search engine')
      throw err
    } finally {
      loading.value = false
    }
  }

  // 更新带图片上传的搜索引擎
  const updateSearchEngineWithIcon = async (
    id: string,
    updates: { name?: string; url?: string; placeholder?: string; isDefault?: boolean },
    iconFile?: File,
  ) => {
    try {
      loading.value = true
      clearError()

      const updatedEngine = await api.updateSearchEngineWithIcon(id, updates, iconFile)

      // 更新本地设置中的搜索引擎
      const index = settings.value.search.engines.findIndex((engine) => engine.id === id)
      if (index !== -1) {
        settings.value.search.engines[index] = updatedEngine
      }

      // 如果更新的是当前搜索引擎，同步更新
      if (currentSearchEngine.value?.id === id) {
        currentSearchEngine.value = updatedEngine
      }

      console.log('Search engine with icon updated successfully:', updatedEngine)
      return updatedEngine
    } catch (err) {
      handleError(err, 'Error updating search engine with icon')
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteSearchEngine = async (id: string) => {
    try {
      loading.value = true
      clearError()

      // 检查是否是最后一个搜索引擎
      if (settings.value.search.engines.length <= 1) {
        throw new Error('不能删除最后一个搜索引擎')
      }

      await api.deleteSearchEngine(id)

      // 从本地设置中移除
      const index = settings.value.search.engines.findIndex((engine) => engine.id === id)
      if (index !== -1) {
        settings.value.search.engines.splice(index, 1)
      }

      // 如果删除的是默认搜索引擎，设置新的默认搜索引擎
      if (
        settings.value.search.defaultEngineId === id &&
        settings.value.search.engines.length > 0
      ) {
        settings.value.search.defaultEngineId = settings.value.search.engines[0].id
        currentSearchEngine.value = settings.value.search.engines[0]
      }

      // 如果删除的是当前搜索引擎，切换到默认搜索引擎
      if (currentSearchEngine.value?.id === id) {
        currentSearchEngine.value = defaultSearchEngine.value
      }

      console.log('Search engine deleted successfully:', id)
    } catch (err) {
      handleError(err, 'Error deleting search engine')
      throw err
    } finally {
      loading.value = false
    }
  }

  // 导入导出功能
  const exportData = async () => {
    try {
      loading.value = true
      clearError()

      const data = await api.exportData()

      // 创建下载链接
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'nav-panel-backup.json'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)

      console.log('Data exported successfully')
    } catch (err) {
      handleError(err, 'Error exporting data')
      throw err
    } finally {
      loading.value = false
    }
  }

  const importData = async (data: {
    websites?: Website[]
    searchEngines?: SearchEngine[]
    settings?: AppSettings
  }) => {
    try {
      loading.value = true
      clearError()

      await api.importData(data)

      // 重新加载数据
      await loadData()

      console.log('Data imported successfully')
    } catch (err) {
      handleError(err, 'Error importing data')
      throw err
    } finally {
      loading.value = false
    }
  }

  // 通知方法
  const showNotification = (
    message: string,
    type: 'success' | 'error' | 'warning' | 'info' = 'info',
  ) => {
    ElNotification({
      title:
        type === 'success'
          ? '成功'
          : type === 'error'
            ? '错误'
            : type === 'warning'
              ? '警告'
              : '提示',
      message,
      type,
      duration: 4000,
      position: 'top-right',
    })
  }

  const showSuccess = (message: string) => showNotification(message, 'success')
  const showError = (message: string) => showNotification(message, 'error')
  const showWarning = (message: string) => showNotification(message, 'warning')
  const showInfo = (message: string) => showNotification(message, 'info')

  // 确认对话框方法
  const showConfirm = async (message: string, title: string = '确认') => {
    try {
      await ElMessageBox.confirm(message, title, {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
      return true
    } catch {
      return false
    }
  }

  // 更新分组中的网站数据的辅助函数
  const updateWebsiteInGroups = (updatedSite: Website) => {
    // 首先从所有分组中移除该网站
    groupsWithWebsites.value.forEach((group) => {
      group.websites = group.websites.filter((site) => site.id !== updatedSite.id)
    })

    // 如果网站有有效的 groupId，将其添加到对应的分组中
    if (updatedSite.groupId && updatedSite.groupId !== '' && updatedSite.groupId !== null) {
      const targetGroup = groupsWithWebsites.value.find((group) => group.id === updatedSite.groupId)
      if (targetGroup) {
        targetGroup.websites.push(updatedSite)
        // 按 sortOrder 排序网站
        targetGroup.websites.sort((a, b) => (a.sortOrder || 0) - (b.sortOrder || 0))
      }
    }
  }

  // 分组相关方法
  const loadGroups = async () => {
    try {
      const [groupsData, groupsWithWebsitesData] = await Promise.all([
        api.getGroups(),
        api.getGroupsWithWebsites(),
      ])

      groups.value = groupsData
      groupsWithWebsites.value = groupsWithWebsitesData

      console.log('Groups loaded:', groups.value.length)
    } catch (err) {
      handleError(err, 'Error loading groups')
      throw err
    }
  }

  const createGroup = async (groupData: { name: string; color?: string; icon?: string }) => {
    try {
      loading.value = true
      clearError()

      const newGroup = await api.createGroup(groupData)
      groups.value.push(newGroup)

      // 重新加载所有数据以确保同步
      await loadData()

      showSuccess('分组创建成功')
      return newGroup
    } catch (err) {
      handleError(err, 'Error creating group')
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateGroup = async (
    id: string,
    updates: { name?: string; color?: string; icon?: string; isCollapsed?: boolean },
  ) => {
    try {
      loading.value = true
      clearError()

      const updatedGroup = await api.updateGroup(id, updates)

      // 更新本地状态
      const index = groups.value.findIndex((g) => g.id === id)
      if (index !== -1) {
        groups.value[index] = updatedGroup
      }

      // 更新groupsWithWebsites中的分组信息
      const groupWithWebsitesIndex = groupsWithWebsites.value.findIndex((g) => g.id === id)
      if (groupWithWebsitesIndex !== -1) {
        groupsWithWebsites.value[groupWithWebsitesIndex] = {
          ...updatedGroup,
          websites: groupsWithWebsites.value[groupWithWebsitesIndex].websites,
        }
      }

      showSuccess('分组更新成功')
      return updatedGroup
    } catch (err) {
      handleError(err, 'Error updating group')
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteGroup = async (id: string) => {
    try {
      loading.value = true
      clearError()

      await api.deleteGroup(id)

      // 从本地状态中移除
      groups.value = groups.value.filter((g) => g.id !== id)
      groupsWithWebsites.value = groupsWithWebsites.value.filter((g) => g.id !== id)

      // 重新加载网站数据，因为分组删除后网站会被移出分组
      const websitesData = await api.getWebsites()
      sites.value = Array.isArray(websitesData) ? websitesData : []

      showSuccess('分组删除成功')
    } catch (err) {
      handleError(err, 'Error deleting group')
      throw err
    } finally {
      loading.value = false
    }
  }

  const toggleGroupCollapse = (groupId: string) => {
    // 更新本地状态中的折叠状态
    const groupIndex = groups.value.findIndex((g) => g.id === groupId)
    if (groupIndex !== -1) {
      groups.value[groupIndex].isCollapsed = !groups.value[groupIndex].isCollapsed
    }

    // 同时更新groupsWithWebsites中的状态
    const groupWithWebsitesIndex = groupsWithWebsites.value.findIndex((g) => g.id === groupId)
    if (groupWithWebsitesIndex !== -1) {
      groupsWithWebsites.value[groupWithWebsitesIndex].isCollapsed =
        !groupsWithWebsites.value[groupWithWebsitesIndex].isCollapsed
    }
  }

  const moveWebsiteToGroup = async (websiteId: string, groupId?: string, position?: number) => {
    try {
      loading.value = true
      clearError()

      await api.moveWebsiteToGroup(websiteId, groupId, position)

      // 重新加载数据
      await Promise.all([loadGroups(), loadData()])

      showSuccess('网站移动成功')
    } catch (err) {
      handleError(err, 'Error moving website')
      throw err
    } finally {
      loading.value = false
    }
  }

  const reorderGroups = async (groupIds: string[]) => {
    try {
      loading.value = true
      clearError()

      await api.reorderGroups(groupIds)

      // 重新排序本地状态
      const orderedGroups = groupIds
        .map((id) => groups.value.find((g) => g.id === id)!)
        .filter(Boolean)
      groups.value = orderedGroups

      showSuccess('分组排序成功')
    } catch (err) {
      handleError(err, 'Error reordering groups')
      throw err
    } finally {
      loading.value = false
    }
  }

  // 分组模态框控制
  const openGroupModal = () => {
    isGroupModalOpen.value = true
    editingGroup.value = null
  }

  const openEditGroupModal = (group: Group) => {
    isGroupModalOpen.value = true
    editingGroup.value = group
  }

  const closeGroupModal = () => {
    isGroupModalOpen.value = false
    editingGroup.value = null
  }

  return {
    // 状态
    sites,
    groups,
    groupsWithWebsites,
    searchQuery,
    isAddCardModalOpen,
    isSettingsModalOpen,
    isGroupModalOpen,
    editingCard,
    editingGroup,
    currentSearchEngine,
    settings,
    loading,
    error,

    // 计算属性
    filteredSites,
    defaultSearchEngine,

    // 默认数据
    defaultSearchEngines,

    // 方法
    loadData,
    addSite,
    updateSite,
    deleteSite,
    reorderSites,
    openAddCardModal,
    openEditCardModal,
    closeAddCardModal,
    openSettingsModal,
    closeSettingsModal,
    updateSettings,
    resetSettings,
    performSearch,
    setCurrentSearchEngine,
    setDefaultSearchEngine,
    exportData,
    importData,
    clearError,
    addSearchEngine,
    updateSearchEngine,
    deleteSearchEngine,
    addSearchEngineWithIcon,
    updateSearchEngineWithIcon,

    // 分组方法
    loadGroups,
    createGroup,
    updateGroup,
    deleteGroup,
    toggleGroupCollapse,
    moveWebsiteToGroup,
    reorderGroups,
    openGroupModal,
    openEditGroupModal,
    closeGroupModal,

    // 通知方法
    showNotification,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    showConfirm,
  }
})
