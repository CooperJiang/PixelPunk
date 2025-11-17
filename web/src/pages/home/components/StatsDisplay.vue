<script setup lang="ts">
  import { computed, nextTick, onMounted, ref } from 'vue'
  import { useCountAnimation } from '@/composables/useCountAnimation'
  import { useSettingsStore } from '@/store/settings'
  import { getPbData } from '@/api/common'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'StatsDisplay',
  })

  const settingsStore = useSettingsStore()
  const { $t } = useTexts()

  const imageStats = ref({
    total: 0,
  })

  const storageStats = ref({
    formattedStorage: '',
    totalStorage: 0,
  })

  const statsNumber = ref<HTMLElement>()
  const _storageNumber = ref<HTMLElement>()

  const shouldShowImageCount = computed(() => settingsStore.showImageCount)
  const shouldShowStorageUsage = computed(() => settingsStore.showStorageUsage)

  const shouldShowStats = computed(() => shouldShowImageCount.value || shouldShowStorageUsage.value)

  const { displayValue, observeElement } = useCountAnimation(
    computed(() => imageStats.value.total),
    {
      duration: 2500,
      startValue: 0,
      separator: ',',
      startOnVisible: true,
      autoStart: false,
    }
  )

  const fetchStats = async () => {
    try {
      const response = await getPbData()

      imageStats.value.total = response.data?.total || 0

      if (response.data?.storage) {
        storageStats.value.formattedStorage = response.data.storage.formatted_storage || ''
        storageStats.value.totalStorage = response.data.storage.total_storage || 0
      } else {
        if (shouldShowStorageUsage.value) {
          storageStats.value.formattedStorage = $t('home.stats.calculating')
        }
      }
    } catch (_error) {
      imageStats.value.total = 0
      storageStats.value.formattedStorage = '0 B'
    }
  }

  const dynamicStatsText = computed(() => {
    const showImageCount = shouldShowImageCount.value
    const showStorageUsage = shouldShowStorageUsage.value
    const hasStorageData =
      storageStats.value.formattedStorage && storageStats.value.formattedStorage !== $t('home.stats.calculating')
    if (showImageCount && showStorageUsage && hasStorageData) {
      return $t('home.stats.countAndStorage', { count: displayValue.value, storage: storageStats.value.formattedStorage })
    } else if (showImageCount) {
      return $t('home.stats.countOnly', { count: displayValue.value })
    } else if (showStorageUsage && hasStorageData) {
      return $t('home.stats.storageOnly', { storage: storageStats.value.formattedStorage })
    }
    return showImageCount ? $t('home.stats.countOnly', { count: displayValue.value }) : ''
  })

  onMounted(async () => {
    if (shouldShowStats.value) {
      await fetchStats()
      nextTick(() => {
        if (statsNumber.value) {
          observeElement(statsNumber.value)
        }
      })
    }
  })
</script>

<template>
  <div v-if="shouldShowStats" class="stats-simple mb-8" style="margin-top: 15px">
    <div class="stats-badge">
      <span ref="statsNumber">{{ dynamicStatsText }}</span>
    </div>
  </div>
</template>

<style scoped>
  .stats-simple {
    display: flex;
    justify-content: center;
  }

  .stats-badge {
    @apply inline-flex items-center text-sm font-medium;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    color: var(--color-content-default);
    background: transparent;
    padding: 0.5rem 1.5rem;
    border-radius: var(--radius-full);
    position: relative;
    overflow: hidden;
  }

  .stats-badge:hover {
    color: var(--color-content-heading);
    transform: translateY(-2px);
  }

  .stats-badge span {
    @apply font-medium;
    color: var(--color-content-heading);
    font-family: 'JetBrains Mono', 'Orbitron', 'Fira Code', monospace;
    font-size: 1em;
    letter-spacing: 0.5px;
    position: relative;
    z-index: 1;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  @media (max-width: 768px) {
    .stats-badge {
      @apply px-2 py-1 text-sm;
    }
  }
</style>
