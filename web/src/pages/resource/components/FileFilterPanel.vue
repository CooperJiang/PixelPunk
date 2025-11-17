<script setup lang="ts">
  import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
  import { getUserColorList, getUserTagList, getUserCategoryList } from '@/api/common'
  import { getFilterSortOptions, getResolutionOptions } from '@/constants'
  import { useTexts } from '@/composables/useTexts'

  const emit = defineEmits(['filter'])
  const { $t } = useTexts()

  let debounceTimer: number | null = null

  const filterForm = reactive({
    keyword: '',
    access_level: 'all',
    sort: 'newest',
    resolution: '',
    tags: '',
    categoryId: '',
    dominant_color: '',
    min_width: undefined,
    min_height: undefined,
  })

  const showAdvanced = ref(false)

  const selectedColors = ref<string[]>([])
  const selectedTags = ref<string[]>([])
  const selectedCategories = ref<string[]>([])

  watch(
    selectedColors,
    (newVal, oldVal) => {
      if (newVal.length > oldVal.length) {
        nextTick(() => {
          scrollToLatestSelection('.selected-tags')
        })
      }

      filterForm.dominant_color = selectedColors.value.join(',')
    },
    { deep: true }
  )

  watch(
    selectedTags,
    (newVal, oldVal) => {
      if (newVal.length > oldVal.length) {
        nextTick(() => {
          scrollToLatestSelection('.selected-tags')
        })
      }

      filterForm.tags = selectedTags.value.join(',')
    },
    { deep: true }
  )

  watch(
    selectedCategories,
    (newVal, oldVal) => {
      if (newVal.length > oldVal.length) {
        nextTick(() => {
          scrollToLatestSelection('.selected-tags')
        })
      }

      filterForm.categoryId = selectedCategories.value.join(',')
    },
    { deep: true }
  )

  const scrollToLatestSelection = (selector: string) => {
    const containers = document.querySelectorAll(selector)
    containers.forEach((container) => {
      if (container) {
        container.scrollLeft = container.scrollWidth
      }
    })
  }

  const _applyFilterWithDebounce = () => {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
    debounceTimer = window.setTimeout(() => {
      applyFilter()
    }, 300)
  }

  const accessLevelOptions = computed(() => [
    { label: $t('resource.filter.accessLevel.all'), value: 'all' },
    { label: $t('resource.filter.accessLevel.public'), value: 'public' },
    { label: $t('resource.filter.accessLevel.private'), value: 'private' },
    { label: $t('resource.filter.accessLevel.protected'), value: 'protected' },
  ])

  const _sortOptions = computed(() => getFilterSortOptions($t))

  const resolutionOptions = computed(() => getResolutionOptions($t))
  const tags = ref<{ id: number; name: string; count: number }[]>([])
  const tagOptions = computed(
    () =>
      tags.value?.map((tag) => ({
        label: `${tag.name}`,
        value: tag.name,
      })) || []
  )

  const colors = ref<string[]>([])
  const _colorOptions = computed(
    () =>
      colors.value?.map((color) => ({
        label: color,
        value: color,
        color,
      })) || []
  )

  const categories = ref<{ id: number; name: string; file_count: number }[]>([])
  const categoryOptions = computed(
    () =>
      categories.value?.map((category) => ({
        label: `${category.name} (${category.file_count})`,
        value: String(category.id),
      })) || []
  )

  const loadTags = async () => {
    try {
      const result = await getUserTagList()
      if (result.success && result.data) {
        tags.value = result.data.items || []
      }
    } catch (error) {
      console.error('Failed to load tags:', error)
    }
  }

  const loadColors = async () => {
    try {
      const result = await getUserColorList()
      if (result.success && result.data) {
        colors.value = result.data.items || []
      }
    } catch (error) {
      console.error('Failed to load colors:', error)
    }
  }

  const loadCategories = async () => {
    try {
      const result = await getUserCategoryList()
      if (result.success && result.data) {
        categories.value = result.data.items || []
      }
    } catch (error) {
      console.error('Failed to load categories:', error)
    }
  }

  const applyFilter = () => {
    emit('filter', { ...filterForm })
  }

  const resetFilter = () => {
    filterForm.keyword = ''
    filterForm.access_level = 'all'
    filterForm.sort = 'newest'
    filterForm.resolution = ''
    filterForm.tags = ''
    filterForm.categoryId = ''
    filterForm.dominant_color = ''
    filterForm.min_width = undefined
    filterForm.min_height = undefined
    selectedColors.value = []
    selectedTags.value = []
    selectedCategories.value = []
    emit('filter', { ...filterForm })
  }

  onMounted(() => {
    loadTags()
    loadColors()
    loadCategories()
  })
</script>

<template>
  <div class="mb-4 rounded-lg border border-subtle bg-background-600 p-3">
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-base font-bold text-content-heading">
        <i class="fas fa-filter mr-1.5 text-content" />{{ $t('resource.filter.title') }}
      </h3>
      <div class="flex space-x-1.5">
        <CyberButton type="secondary" class="px-2 py-1 text-xs" @click="resetFilter">
          <i class="fas fa-undo-alt mr-1" />{{ $t('resource.filter.reset') }}
        </CyberButton>
        <CyberButton type="primary" class="px-2 py-1 text-xs" @click="applyFilter">
          <i class="fas fa-search mr-1" />{{ $t('resource.filter.apply') }}
        </CyberButton>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-x-3 gap-y-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('resource.filter.accessLevelLabel') }}</label>
        <CyberDropdown
          v-model="filterForm.access_level"
          :options="accessLevelOptions"
          :placeholder="$t('resource.filter.accessLevel.all')"
          class="compact"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('resource.filter.categoryLabel') }}</label>
        <CyberDropdown
          v-model="selectedCategories"
          :options="categoryOptions"
          :placeholder="$t('resource.filter.categoryPlaceholder')"
          multiple
          searchable
          class="compact"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('resource.filter.tagLabel') }}</label>
        <CyberDropdown
          v-model="selectedTags"
          :options="tagOptions"
          :placeholder="$t('resource.filter.tagPlaceholder')"
          multiple
          searchable
          class="compact"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('resource.filter.resolutionLabel') }}</label>
        <CyberDropdown
          v-model="filterForm.resolution"
          :options="resolutionOptions"
          :placeholder="$t('resource.filter.resolutionPlaceholder')"
          class="compact"
        />
      </div>
    </div>

    <div v-if="showAdvanced" class="mt-2 grid grid-cols-1 gap-x-3 gap-y-2 sm:grid-cols-2 md:grid-cols-3"></div>
  </div>
</template>
