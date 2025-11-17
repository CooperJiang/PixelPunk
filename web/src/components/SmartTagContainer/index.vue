<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import type { SmartTagContainerProps, SmartTagContainerEmits } from './types'

  defineOptions({
    name: 'SmartTagContainer',
  })

  const props = withDefaults(defineProps<SmartTagContainerProps>(), {
    variant: 'primary',
    size: 'small',
    truncate: true,
    minTags: 1,
    moreVariant: 'secondary',
  })

  const emit = defineEmits<SmartTagContainerEmits>()

  const containerRef = ref<HTMLElement>()
  const measureRef = ref<HTMLElement>()
  const visibleTagCount = ref(props.tags.length)

  const calculateVisibleTags = async () => {
    if (!containerRef.value || !measureRef.value || props.tags.length === 0) {
      return
    }

    await nextTick()

    const containerWidth = containerRef.value.offsetWidth
    const gap = 4 // gap-1 = 4px

    if (props.tags.length === 0) {
      visibleTagCount.value = 0
      return
    }

    const measureTags = measureRef.value.querySelectorAll('.cyber-tag')

    if (measureTags.length === 0) {
      visibleTagCount.value = props.tags.length
      return
    }

    let allTagsWidth = 0
    for (let i = 0; i < measureTags.length; i++) {
      const tagElement = measureTags[i] as HTMLElement
      if (tagElement) {
        allTagsWidth += tagElement.offsetWidth
        if (i > 0) allTagsWidth += gap
      }
    }

    if (allTagsWidth <= containerWidth) {
      visibleTagCount.value = props.tags.length
      return
    }

    const estimatedMoreTagWidth = 50 // 估算+xx按钮宽度

    let totalTagsWidth = 0
    let visibleCount = 0
    const _canAddTruncatedTag = false

    for (let i = 0; i < measureTags.length; i++) {
      const tagElement = measureTags[i] as HTMLElement
      if (!tagElement) continue

      const tagWidth = tagElement.offsetWidth
      const gapWidth = visibleCount > 0 ? gap : 0
      const newTotalWidth = totalTagsWidth + gapWidth + tagWidth
      const totalWithMoreButton = newTotalWidth + gap + estimatedMoreTagWidth

      if (totalWithMoreButton <= containerWidth) {
        totalTagsWidth = newTotalWidth
        visibleCount++
      } else {
        const minTruncatedTagWidth = 30
        const truncatedTotalWidth = totalTagsWidth + gapWidth + minTruncatedTagWidth + gap + estimatedMoreTagWidth

        if (truncatedTotalWidth <= containerWidth && props.truncate) {
          visibleCount++
          const _canAddTruncatedTag = true
        }
        break
      }
    }

    const finalVisibleCount = Math.max(visibleCount, props.minTags)
    visibleTagCount.value = finalVisibleCount
  }

  const visibleTags = computed(() => {
    return props.tags.slice(0, visibleTagCount.value)
  })

  const remainingCount = computed(() => {
    return Math.max(0, props.tags.length - visibleTagCount.value)
  })

  const showMoreButton = computed(() => {
    return remainingCount.value > 0
  })

  const handleTagClick = (tag: string, index: number) => {
    emit('tag-click', tag, index)
  }

  watch(
    () => props.tags,
    () => {
      visibleTagCount.value = props.tags.length
      nextTick(() => {
        calculateVisibleTags()
      })
    },
    { immediate: true }
  )

  let resizeObserver: ResizeObserver | null = null

  onMounted(() => {
    setTimeout(() => {
      calculateVisibleTags()
    }, 100)

    if (containerRef.value && ResizeObserver) {
      resizeObserver = new ResizeObserver(() => {
        setTimeout(() => {
          calculateVisibleTags()
        }, 10)
      })
      resizeObserver.observe(containerRef.value)
    }
  })

  onUnmounted(() => {
    if (resizeObserver && containerRef.value) {
      resizeObserver.unobserve(containerRef.value)
      resizeObserver = null
    }
  })
</script>

<template>
  <div ref="containerRef" class="smart-tag-container flex w-full items-center gap-1 overflow-hidden">
    <template v-for="(tag, index) in visibleTags" :key="`visible-${index}`">
      <CyberTag
        :variant="variant"
        :size="size"
        :truncate="truncate && index === visibleTagCount - 1 && showMoreButton"
        class="flex-shrink-0"
        :class="{
          'truncated-tag': truncate && index === visibleTagCount - 1 && showMoreButton,
        }"
        @click="handleTagClick(tag, index)"
      >
        {{ tag }}
      </CyberTag>
    </template>

    <!-- 更多标签按钮 - 紧跟在最后一个标签后面 -->
    <CyberTag v-if="showMoreButton" :variant="moreVariant" :size="size" class="flex-shrink-0"> +{{ remainingCount }} </CyberTag>

    <!-- 隐藏的标签用于测量 -->
    <div ref="measureRef" class="pointer-events-none invisible absolute left-0 top-0">
      <div class="flex items-center gap-1 whitespace-nowrap">
        <CyberTag v-for="(tag, index) in props.tags" :key="`measure-${index}`" :variant="variant" :size="size" :truncate="false">
          {{ tag }}
        </CyberTag>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .smart-tag-container {
    min-height: 20px;
  }

  .smart-tag-container .cyber-tag.truncated-tag {
    flex-shrink: 1;
    min-width: 30px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .smart-tag-container .cyber-tag:not(.truncated-tag) {
    flex-shrink: 0;
  }
</style>
