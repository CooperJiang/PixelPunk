<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
  import type { ExpiryStatusConfig, ImageExpiryTagProps } from './types'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'ImageExpiryTag',
  })

  const props = withDefaults(defineProps<ImageExpiryTagProps>(), {
    expiresAt: null,
    storageDuration: '',
    isTimeLimited: false,
    position: 'top-left',
    mode: 'both',
    showIcon: true,
    className: '',
  })

  const remainingTime = ref<number>(0)
  /* 在浏览器环境中使用 number 计时器类型，避免 NodeJS.Timeout 的类型误用 */
  const timer = ref<number | null>(null)

  const { $t } = useTexts()

  const isVisible = computed(() => props.isTimeLimited && props.expiresAt)

  const positionClass = computed(() => {
    switch (props.position) {
      case 'top-left':
        return 'pos-top-left'
      case 'top-right':
        return 'pos-top-right'
      case 'bottom-left':
        return 'pos-bottom-left'
      case 'bottom-right':
        return 'pos-bottom-right'
      default:
        return 'pos-top-left'
    }
  })

  const expiryStatus = computed(() => {
    if (!remainingTime.value) {
      return 'expired'
    }

    const oneDay = 24 * 60 * 60 * 1000
    const oneHour = 60 * 60 * 1000

    if (remainingTime.value <= 0) {
      return 'expired'
    } else if (remainingTime.value <= oneHour) {
      return 'expiring-soon'
    } else if (remainingTime.value <= oneDay) {
      return 'expiring-soon'
    }
    return 'active'
  })

  const statusConfig = computed<ExpiryStatusConfig>(() => {
    const configs = {
      expired: {
        color: 'var(--color-error-500)',
        bgColor: 'rgba(var(--color-error-rgb), 0.15)',
        borderColor: 'rgba(var(--color-error-rgb), 0.3)',
        icon: 'fas fa-exclamation-triangle',
        text: $t('components.fileExpiryTag.statuses.expired'),
      },
      'expiring-soon': {
        color: 'var(--color-warning-500)',
        bgColor: 'rgba(var(--color-warning-rgb), 0.15)',
        borderColor: 'rgba(var(--color-warning-rgb), 0.3)',
        icon: 'fas fa-clock',
        text: $t('components.fileExpiryTag.statuses.expiringSoon'),
      },
      active: {
        color: 'var(--color-brand-500)',
        bgColor: 'rgba(var(--color-brand-500-rgb), 0.15)',
        borderColor: 'rgba(var(--color-brand-500-rgb), 0.3)',
        icon: 'fas fa-hourglass-half',
        text: $t('components.fileExpiryTag.statuses.active'),
      },
    } as const
    return configs[expiryStatus.value as keyof typeof configs]
  })

  const statusClass = computed(() => `status-${expiryStatus.value}`)

  const remainingText = computed(() => {
    if (!remainingTime.value || remainingTime.value <= 0) {
      return $t('components.fileExpiryTag.statuses.expired')
    }

    const dayUnit = $t('components.fileExpiryTag.units.day')
    const hourUnit = $t('components.fileExpiryTag.units.hour')
    const minuteUnit = $t('components.fileExpiryTag.units.minute')

    const days = Math.floor(remainingTime.value / (24 * 60 * 60 * 1000))
    const hours = Math.floor((remainingTime.value % (24 * 60 * 60 * 1000)) / (60 * 60 * 1000))
    const minutes = Math.floor((remainingTime.value % (60 * 60 * 1000)) / (60 * 1000))

    if (days > 0) {
      return `${days}${dayUnit}${hours}${hourUnit}`
    } else if (hours > 0) {
      return `${hours}${hourUnit}${minutes}${minuteUnit}`
    }
    return `${minutes}${minuteUnit}`
  })

  const durationText = computed(() => {
    if (!props.storageDuration) {
      return ''
    }

    const unitDay = $t('components.fileExpiryTag.units.day')
    const unitHour = $t('components.fileExpiryTag.units.hour')
    const unitMinute = $t('components.fileExpiryTag.units.minute')

    const duration = props.storageDuration.toLowerCase()
    if (duration.includes('d')) {
      return duration.replace('d', unitDay)
    } else if (duration.includes('h')) {
      return duration.replace('h', unitHour)
    } else if (duration.includes('m')) {
      return duration.replace('m', unitMinute)
    }
    return props.storageDuration
  })

  const displayText = computed(() => {
    if (props.mode === 'countdown') {
      return remainingText.value
    } else if (props.mode === 'duration') {
      return durationText.value || remainingText.value || statusConfig.value.text
    }
    if (durationText.value && expiryStatus.value !== 'expired') {
      return `${durationText.value} (${remainingText.value})`
    }
    return remainingText.value
  })

  const tooltipText = computed(() => {
    if (!props.expiresAt) {
      return ''
    }

    const expiryDate = new Date(props.expiresAt)
    const formattedDate = expiryDate.toLocaleString(getCurrentLocale())

    const expiredPrefix = $t('components.fileExpiryTag.labels.expiredAtPrefix')
    const expiredSuffix = $t('components.fileExpiryTag.labels.expiredAtSuffix')
    const storageDurationLabel = $t('components.fileExpiryTag.labels.storageDuration')
    const expiresAtLabel = $t('components.fileExpiryTag.labels.expiresAt')

    if (expiryStatus.value === 'expired') {
      return `${expiredPrefix} ${formattedDate} ${expiredSuffix}`
    } else if (durationText.value) {
      return `${storageDurationLabel}: ${durationText.value}\n${expiresAtLabel}: ${formattedDate}`
    }
    return `${expiresAtLabel}: ${formattedDate}`
  })

  const updateRemainingTime = () => {
    if (!props.expiresAt) {
      remainingTime.value = 0
      return
    }

    const expiryDate = new Date(props.expiresAt).getTime()
    const now = Date.now()
    remainingTime.value = expiryDate - now
  }

  const startTimer = () => {
    if (!props.expiresAt) {
      return
    }

    updateRemainingTime()

    if (remainingTime.value > 0) {
      timer.value = window.setInterval(updateRemainingTime, 60 * 1000)
    }
  }

  const clearTimer = () => {
    if (timer.value) {
      window.clearInterval(timer.value)
      timer.value = null
    }
  }

  onMounted(() => {
    startTimer()
  })

  onUnmounted(() => {
    clearTimer()
  })

  watch(
    [() => props.expiresAt, () => props.isTimeLimited],
    () => {
      clearTimer()
      if (isVisible.value) {
        startTimer()
      }
    },
    { immediate: true }
  )
</script>

<template>
  <div v-if="isVisible" class="cyber-file-expiry-tag" :class="[positionClass, statusClass, props.className]" :title="tooltipText">
    <i v-if="showIcon" :class="statusConfig.icon" class="expiry-icon" />
    <span class="expiry-text">{{ displayText }}</span>
  </div>
</template>

<style scoped>
  .cyber-file-expiry-tag {
    @apply absolute z-[5] inline-flex h-5 min-w-0 max-w-[140px] items-center gap-1 overflow-hidden text-ellipsis whitespace-nowrap rounded border px-1.5 py-0.5 text-[10px] font-semibold leading-none tracking-[0.2px] backdrop-blur transition-all duration-300 ease-in-out;
    box-shadow: 0 2px 8px var(--color-overlay-light);
  }

  .cyber-file-expiry-tag:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px var(--color-overlay-medium);
  }

  .pos-top-left {
    @apply left-0.5 top-0.5;
  }

  .pos-top-right {
    @apply right-2 top-2;
  }

  .pos-bottom-left {
    @apply bottom-2 left-2;
  }

  .pos-bottom-right {
    @apply bottom-2 right-2;
  }

  .status-active {
    color: var(--color-brand-500);
    background-color: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .status-active:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .status-expiring-soon {
    color: var(--color-warning-500);
    background-color: rgba(var(--color-warning-rgb), 0.15);
    border-color: rgba(var(--color-warning-rgb), 0.3);
    animation: pulse-warning 2s ease-in-out infinite;
  }

  .status-expiring-soon:hover {
    background-color: rgba(var(--color-warning-rgb), 0.2);
    border-color: rgba(var(--color-warning-rgb), 0.4);
  }

  .status-expired {
    color: var(--color-error-500);
    background-color: rgba(var(--color-error-rgb), 0.15);
    border-color: rgba(var(--color-error-rgb), 0.3);
  }

  .status-expired:hover {
    background-color: rgba(var(--color-error-rgb), 0.2);
    border-color: rgba(var(--color-error-rgb), 0.4);
  }

  .expiry-icon {
    @apply text-[10px];
  }

  @media (max-width: 640px) {
    .cyber-file-expiry-tag {
      font-size: 10px;
      height: 18px;
      padding: 1px 4px;
    }

    .expiry-icon {
      font-size: 9px;
    }
  }
</style>
