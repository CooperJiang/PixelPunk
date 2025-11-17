<script setup lang="ts">
  import { computed, defineEmits, defineProps, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps({
    current: {
      type: Number,
      required: true,
    },
    total: {
      type: Number,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    threshold: {
      type: Number,
      default: 0.5, // 元素露出多少比例时触发加载
    },
    throttle: {
      type: Number,
      default: 500, // 减少节流时间，提高响应性
    },
    externalLoading: {
      type: Boolean,
      default: false,
    },
    canLoadMore: {
      type: Boolean,
      default: true,
    },
  })

  const emit = defineEmits(['load-more'])
  const loaderRef = ref(null)
  const loading = ref(false)
  const observer = ref(null)
  const needsRecheck = ref(false) // 标记是否需要重新检查
  let lastLoadTime = 0

  const isComplete = computed(() => props.current >= props.total || !props.canLoadMore)

  const isInViewport = () => {
    if (!loaderRef.value) {
      return false
    }

    const rect = loaderRef.value.getBoundingClientRect()
    const windowHeight = window.innerHeight || document.documentElement.clientHeight
    const windowWidth = window.innerWidth || document.documentElement.clientWidth

    return (
      rect.top >= 0 &&
      rect.left >= 0 &&
      rect.bottom <= windowHeight + 200 && // 提前200px触发
      rect.right <= windowWidth
    )
  }

  const setupObserver = () => {
    if (isComplete.value || props.disabled) {
      destroyObserver()
      return
    }

    destroyObserver()

    observer.value = new IntersectionObserver(
      (entries) => {
        const entry = entries[0]
        if (entry.isIntersecting && !props.externalLoading && !loading.value && props.canLoadMore) {
          const now = Date.now()
          if (now - lastLoadTime > props.throttle) {
            loadMore()
            lastLoadTime = now
          }
        }
      },
      {
        root: null, // 使用viewport作为观察区域
        rootMargin: '200px', // 提前200px触发，增加触发范围
        threshold: [0, props.threshold, 1], // 多个阈值，提高检测精度
      }
    )

    if (loaderRef.value) {
      observer.value.observe(loaderRef.value)
    }
  }

  const destroyObserver = () => {
    if (observer.value) {
      observer.value.disconnect()
      observer.value = null
    }
  }

  const loadMore = () => {
    if (isComplete.value || props.externalLoading || loading.value || props.disabled || !props.canLoadMore) {
      return
    }

    loading.value = true
    needsRecheck.value = false
    emit('load-more')
  }

  const recheckLoadMore = () => {
    if (props.externalLoading || loading.value || isComplete.value || props.disabled) {
      return
    }

    nextTick(() => {
      setTimeout(() => {
        if (isInViewport() && props.canLoadMore && !isComplete.value) {
          loadMore()
        }
      }, 100)
    })
  }

  watch(
    () => props.externalLoading,
    (newLoading, oldLoading) => {
      if (oldLoading && !newLoading) {
        loading.value = false

        if (needsRecheck.value || isInViewport()) {
          recheckLoadMore()
        }
      }
    }
  )

  watch(
    () => [props.current, props.total, props.disabled, props.canLoadMore],
    () => {
      nextTick(() => {
        setupObserver()

        if (!props.externalLoading && !loading.value) {
          recheckLoadMore()
        }
      })
    }
  )

  const handleScroll = () => {
    if (props.externalLoading || loading.value || isComplete.value || props.disabled) {
      return
    }

    if (isInViewport()) {
      needsRecheck.value = true
      const now = Date.now()
      if (now - lastLoadTime > props.throttle) {
        loadMore()
        lastLoadTime = now
      }
    }
  }

  onMounted(() => {
    setupObserver()

    window.addEventListener('scroll', handleScroll, { passive: true })

    nextTick(() => {
      if (isInViewport() && props.canLoadMore && !isComplete.value) {
        setTimeout(() => {
          if (!props.externalLoading && !loading.value) {
            loadMore()
          }
        }, 100)
      }
    })
  })

  onUnmounted(() => {
    destroyObserver()
    window.removeEventListener('scroll', handleScroll)
  })
</script>

<template>
  <div ref="loaderRef" class="auto-loader">
    <div v-if="externalLoading || loading" class="loading-indicator">
      <i class="fas fa-circle-notch fa-spin" />
      <span>{{ $t('share.autoLoader.loading') }}</span>
    </div>
    <div v-else-if="!isComplete" class="load-more-info">
      <span class="count-info">{{ $t('share.autoLoader.loaded', { current, total }) }}</span>
    </div>
    <div v-else class="all-loaded">
      <span>{{ $t('share.autoLoader.allLoaded') }}</span>
    </div>
  </div>
</template>

<style scoped>
  .auto-loader {
    margin: var(--space-md) 0;
    text-align: center;
    height: 50px;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .loading-indicator {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    color: rgba(var(--color-content-rgb), 0.8);
    font-size: 0.9rem;
  }

  .loading-indicator i {
    color: rgba(var(--color-brand-500-rgb), 0.8);
    font-size: 1.1rem;
  }

  .load-more-info {
    padding: var(--space-md) 1rem;
    border-radius: var(--radius-md);
    background: rgba(var(--color-background-800-rgb), 0.4);
    color: var(--color-content-default);
    font-size: 0.85rem;
  }

  .count-info {
    color: var(--color-content-default);
  }

  .all-loaded {
    color: var(--color-content-default);
    font-size: 0.85rem;
    padding: var(--space-md) 1rem;
    border-top: 1px dashed var(--color-content-default);
    border-bottom: 1px dashed var(--color-content-default);
  }
</style>
