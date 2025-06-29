import type { Website, AppSettings, SearchEngine } from '@/types'

const API_BASE_URL = 'http://localhost:8080/api'

class ApiClient {
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    try {
      const response = await fetch(url, config)

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      return data.data || data
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error)
      throw error
    }
  }

  // 多部分表单请求（用于文件上传）
  private async requestWithFile<T>(endpoint: string, formData: FormData): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    try {
      const response = await fetch(url, {
        method: 'POST',
        body: formData,
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      return data.data || data
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error)
      throw error
    }
  }

  // 将文件转换为DataURL
  private fileToDataURL(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result as string)
      reader.onerror = reject
      reader.readAsDataURL(file)
    })
  }

  // 网站相关API
  async getWebsites(): Promise<Website[]> {
    return this.request<Website[]>('/websites')
  }

  async createWebsite(website: Omit<Website, 'id' | 'createdAt' | 'updatedAt'>): Promise<Website> {
    return this.request<Website>('/websites', {
      method: 'POST',
      body: JSON.stringify(website),
    })
  }

  async updateWebsite(id: string, updates: Partial<Website>): Promise<Website> {
    return this.request<Website>(`/websites/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updates),
    })
  }

  async deleteWebsite(id: string): Promise<void> {
    return this.request<void>(`/websites/${id}`, {
      method: 'DELETE',
    })
  }

  async reorderWebsites(websiteIds: string[]): Promise<Website[]> {
    return this.request<Website[]>('/websites/reorder', {
      method: 'PUT',
      body: JSON.stringify({ websiteIds }),
    })
  }

  // 搜索引擎相关API
  async getSearchEngines(): Promise<SearchEngine[]> {
    return this.request<SearchEngine[]>('/search-engines')
  }

  async createSearchEngine(engine: Omit<SearchEngine, 'id'>): Promise<SearchEngine> {
    return this.request<SearchEngine>('/search-engines', {
      method: 'POST',
      body: JSON.stringify(engine),
    })
  }

  // 创建带图片上传的搜索引擎
  async createSearchEngineWithIcon(
    engine: { name: string; url: string; placeholder?: string; isDefault?: boolean },
    iconFile?: File,
  ): Promise<SearchEngine> {
    // 生成ID
    const id = engine.name.toLowerCase().replace(/[^a-z0-9]/g, '') + '_' + Date.now()

    let iconDataURL: string | undefined
    if (iconFile) {
      // 将文件转换为DataURL
      iconDataURL = await this.fileToDataURL(iconFile)
    }

    return this.request<SearchEngine>('/search-engines', {
      method: 'POST',
      body: JSON.stringify({
        id,
        name: engine.name,
        url: engine.url,
        icon: iconDataURL,
        placeholder: engine.placeholder,
        isDefault: engine.isDefault,
      }),
    })
  }

  async updateSearchEngine(id: string, updates: Partial<SearchEngine>): Promise<SearchEngine> {
    return this.request<SearchEngine>(`/search-engines/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updates),
    })
  }

  // 更新带图片上传的搜索引擎
  async updateSearchEngineWithIcon(
    id: string,
    updates: { name?: string; url?: string; placeholder?: string; isDefault?: boolean },
    iconFile?: File,
  ): Promise<SearchEngine> {
    let iconDataURL: string | undefined
    if (iconFile) {
      // 将文件转换为DataURL
      iconDataURL = await this.fileToDataURL(iconFile)
    }

    return this.request<SearchEngine>(`/search-engines/${id}`, {
      method: 'PUT',
      body: JSON.stringify({
        ...updates,
        icon: iconDataURL,
      }),
    })
  }

  async deleteSearchEngine(id: string): Promise<void> {
    return this.request<void>(`/search-engines/${id}`, {
      method: 'DELETE',
    })
  }

  // 设置相关API
  async getSettings(): Promise<AppSettings> {
    return this.request<AppSettings>('/settings')
  }

  async updateSettings(settings: Partial<AppSettings>): Promise<AppSettings> {
    return this.request<AppSettings>('/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    })
  }

  // 导入导出API
  async exportData(): Promise<{
    websites: Website[]
    searchEngines: SearchEngine[]
    settings: AppSettings
  }> {
    return this.request<{
      websites: Website[]
      searchEngines: SearchEngine[]
      settings: AppSettings
    }>('/export')
  }

  async importData(data: {
    websites?: Website[]
    searchEngines?: SearchEngine[]
    settings?: AppSettings
  }): Promise<void> {
    return this.request<void>('/import', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }
}

export const api = new ApiClient()
