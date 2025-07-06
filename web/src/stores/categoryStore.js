import { defineStore } from 'pinia'
import { listCategories } from '@/api/categories'

export const useCategoryStore = defineStore('category', {
  state: () => ({
    categoriesList: [],
    total: 0,
    loaded: false
  }),
  actions: {
    async fetchCategories(force = false) {
      if (this.loaded && !force) return this.categoriesList
      try {
        const res = await listCategories({ page: 1, pageSize: 99 })
        this.categoriesList = res.data.data.categoriesList || []
        this.total = res.data.data.total || 0
        this.loaded = true
        return this.categoriesList
      } catch (e) {
        this.loaded = false
        throw e
      }
    },
    async getCategories() {
      if (this.categoriesList.length > 0) {
        return this.categoriesList
      } else {
        return await this.fetchCategories()
      }
    }
  }
}) 