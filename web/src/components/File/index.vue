<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import type { FileEmits, FileInfo, FileProps } from './types'
  import { useTexts } from '@/composables/useTexts'
  import { getFileTypeInfo, FILE_TYPES } from '@/constants/file'

  defineOptions({ name: 'CyberFile' })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<FileProps>(), {
    alt: '',
    width: '',
    height: '',
    fitMode: 'cover',
    rounded: false,
    borderRadius: '4px',
    errorText: '',
    className: '',
    placeholder: '',
    retryCount: 0,
    retryInterval: 2000,
    maintainAspectRatio: true,
    showDimensions: false,
    backgroundColor: 'var(--color-background-900)',
    backgroundPattern: 'cyber',
    isNsfw: false,
  })

  const emit = defineEmits<FileEmits>()

  /* 状态管理（方法内部不加注释） */
  const loadError = ref(false)
  const fileRef = ref<HTMLImageElement | null>(null)
  const retryAttempts = ref(0)
  const fileInfo = ref<FileInfo>({ width: 0, height: 0, aspectRatio: 1 })
  const showNsfwContent = ref(false)

  /* 文件类型信息 */
  const fileTypeInfo = computed(() => {
    const type = props.fileType || FILE_TYPES.IMAGE
    return getFileTypeInfo(type, $t)
  })

  /* 是否为图片文件 */
  const isImageFile = computed(() => {
    return !props.fileType || props.fileType === FILE_TYPES.IMAGE
  })

  /* 动态CSS类名 */
  const fileTypeClasses = computed(() => {
    return ['cyber-file', fileTypeInfo.value.className, `file-${fileTypeInfo.value.color}`]
  })

  /* 主题化的展示文案 */
  const displayErrorText = computed(() => props.errorText || $t('utils.file.loadFailed'))

  /* 查看NSFW文件 */
  const viewNsfwImage = () => {
    showNsfwContent.value = true
    emit('view-nsfw')
  }

  const getFileStyle = () => {
    const style: Record<string, string> = { objectFit: props.fitMode }
    if (props.maintainAspectRatio && fileInfo.value.aspectRatio && props.fitMode === 'contain') {
      style.maxWidth = '100%'
      style.maxHeight = '100%'
      style.width = 'auto'
      style.height = 'auto'
    }
    return style
  }

  const handleImageLoaded = () => {
    if (fileRef.value) {
      const { naturalWidth, naturalHeight } = fileRef.value
      fileInfo.value = { width: naturalWidth, height: naturalHeight, aspectRatio: naturalWidth / naturalHeight }
      emit('dimensions', naturalWidth, naturalHeight)
    }
    loadError.value = false
    emit('load')
  }

  const handleImageError = () => {
    if (retryAttempts.value < props.retryCount) {
      retryAttempts.value++
      setTimeout(() => {
        if (fileRef.value) {
          const timestamp = Date.now()
          const separator = props.src.includes('?') ? '&' : '?'
          fileRef.value.src = `${props.src}${separator}_retry=${timestamp}`
        }
      }, props.retryInterval)
      return
    }
    loadError.value = true
    emit('error')
  }

  const preloadImage = (src: string) => {
    loadError.value = false
    retryAttempts.value = 0
    if (!src) {
      loadError.value = true
    }
  }

  watch(
    () => props.src,
    (newSrc) => {
      if (newSrc) preloadImage(newSrc)
    }
  )

  onMounted(() => {
    if (props.src) preloadImage(props.src)
  })

  const getContainerStyle = computed(() => {
    const style: Record<string, string> = { width: props.width, height: props.height, borderRadius: props.borderRadius }
    return style
  })
</script>

<template>
  <div
    class="cyber-file-container"
    :class="[
      fileTypeClasses,
      `cyber-file-${fitMode}`,
      { 'cyber-file-rounded': rounded },
      { 'cyber-file-maintain-ratio': maintainAspectRatio },
      {
        [`cyber-file-bg-${backgroundPattern}`]: maintainAspectRatio && backgroundPattern !== 'none',
      },
      className,
    ]"
    :style="getContainerStyle"
  >
    <div v-show="loadError" class="cyber-file-placeholder">
      <div class="cyber-file-error">
        <i class="fas fa-exclamation-triangle text-error-500" />
        <span class="cyber-file-error-text">{{ displayErrorText }}</span>
      </div>
    </div>

    <img
      v-if="isImageFile && !loadError && (!isNsfw || showNsfwContent)"
      ref="fileRef"
      :src="src"
      :alt="alt"
      :style="getFileStyle()"
      @load="handleImageLoaded"
      @error="handleImageError"
    />

    <div v-else-if="!isImageFile && !loadError" class="cyber-file-icon-display">
      <div class="cyber-file-icon-wrapper">
        <i :class="fileTypeInfo.icon" class="cyber-file-type-icon" />
        <span class="cyber-file-type-name">{{ fileTypeInfo.name }}</span>
      </div>
      <div v-if="alt" class="cyber-file-name">{{ alt }}</div>
    </div>

    <div v-if="isNsfw && !showNsfwContent && !loadError" class="cyber-image-nsfw-overlay">
      <i class="fas fa-eye-slash cyber-image-nsfw-icon" />
      <p class="cyber-image-nsfw-text">{{ $t('utils.file.nsfwWarning') }}</p>
      <button class="cyber-image-nsfw-button" @click="viewNsfwImage">
        <i class="fas fa-eye cyber-image-nsfw-button-icon" />{{ $t('actions.view') }}
      </button>
    </div>

    <div v-if="$slots.overlay" class="cyber-image-overlay">
      <slot name="overlay" />
    </div>

    <div v-if="showDimensions && fileInfo.width && !loadError && isImageFile" class="cyber-file-dimensions">
      <span class="cyber-file-dimension-text">{{ fileInfo.width }} × {{ fileInfo.height }}</span>
    </div>
  </div>
</template>

<style scoped>
  .cyber-file-container {
    position: relative;
    overflow: hidden;
    background-color: transparent;
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 80px;
    min-height: 100px;
    transition: all var(--transition-slow) var(--ease-in-out);
  }

  .file-type-image .cyber-file-container {
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .file-type-video .cyber-file-container {
    border: 2px solid rgba(var(--color-error-rgb), 0.3);
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.05) 0%, transparent 100%);
  }

  .file-type-document .cyber-file-container {
    border: 2px solid rgba(var(--color-success-rgb), 0.3);
    background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.05) 0%, transparent 100%);
  }

  .file-type-archive .cyber-file-container {
    border: 2px solid rgba(var(--color-warning-rgb), 0.3);
    background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.05) 0%, transparent 100%);
  }

  .file-type-audio .cyber-file-container {
    border: 2px solid rgba(var(--color-error-rgb), 0.5);
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.08) 0%, transparent 100%);
  }

  .file-type-other .cyber-file-container {
    border: 2px solid rgba(var(--color-content-rgb), 0.2);
    background: linear-gradient(135deg, rgba(var(--color-content-rgb), 0.03) 0%, transparent 100%);
  }

  .cyber-file-icon-display {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    text-align: center;
    width: 100%;
    height: 100%;
  }

  .cyber-file-icon-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .cyber-file-type-icon {
    font-size: 2.5rem;
    margin-bottom: 0.5rem;
    opacity: 0.8;
  }

  .file-blue .cyber-file-type-icon {
    color: var(--color-brand-500);
  }
  .file-purple .cyber-file-type-icon {
    color: var(--color-error-500);
  }
  .file-green .cyber-file-type-icon {
    color: var(--color-success-500);
  }
  .file-orange .cyber-file-type-icon {
    color: var(--color-warning-500);
  }
  .file-pink .cyber-file-type-icon {
    color: var(--color-error-500);
  }
  .file-gray .cyber-file-type-icon {
    color: var(--color-content-muted);
  }

  .cyber-file-type-name {
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--color-content-default);
    opacity: 0.8;
    text-transform: uppercase;
    letter-spacing: var(--tracking-wide);
  }

  .cyber-file-name {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .cyber-file-bg-grid {
    background-image:
      linear-gradient(rgba(var(--color-brand-500-rgb), 0.1) 1px, transparent 1px),
      linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.1) 1px, transparent 1px);
    background-size: 20px 20px;
  }

  .cyber-file-bg-dots {
    background-image: radial-gradient(rgba(var(--color-brand-500-rgb), 0.15) 1px, transparent 1px);
    background-size: 10px 10px;
  }

  .cyber-file-bg-diagonal {
    background-image: repeating-linear-gradient(
      45deg,
      rgba(var(--color-brand-500-rgb), 0.05),
      rgba(var(--color-brand-500-rgb), 0.05) 5px,
      rgba(var(--color-error-rgb), 0.05) 5px,
      rgba(var(--color-error-rgb), 0.05) 10px
    );
  }

  .cyber-file-bg-cyber {
    background:
      linear-gradient(
        to bottom,
        rgba(var(--color-error-rgb), 0.05) 0%,
        transparent 50%,
        rgba(var(--color-brand-500-rgb), 0.05) 100%
      ),
      repeating-linear-gradient(
        to right,
        transparent 0,
        transparent 10px,
        rgba(var(--color-error-rgb), 0.1) 10px,
        rgba(var(--color-error-rgb), 0.1) 12px
      ),
      repeating-linear-gradient(
        to bottom,
        transparent 0,
        transparent 10px,
        rgba(var(--color-brand-500-rgb), 0.1) 10px,
        rgba(var(--color-brand-500-rgb), 0.1) 12px
      );
    background-size:
      100% 100%,
      20px 100%,
      100% 20px;
    background-position: center;
    position: relative;
  }

  .cyber-file-bg-cyber::after {
    content: '';
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    width: 1px;
    background: linear-gradient(to bottom, transparent, rgba(var(--color-error-rgb), 0.5), transparent);
    box-shadow: 0 0 10px rgba(var(--color-error-rgb), 0.5);
    animation: cyber-scan-vert 4s linear infinite;
  }

  .cyber-file-bg-circuit {
    background-color: var(--color-background-900);
    background-image:
      radial-gradient(circle at 25% 25%, rgba(var(--color-brand-500-rgb), 0.2) 1px, transparent 1px),
      radial-gradient(circle at 75% 75%, rgba(var(--color-error-rgb), 0.2) 1px, transparent 1px),
      linear-gradient(to right, rgba(var(--color-brand-500-rgb), 0.1) 1px, transparent 1px),
      linear-gradient(to bottom, rgba(var(--color-brand-500-rgb), 0.1) 1px, transparent 1px);
    background-size:
      20px 20px,
      20px 20px,
      10px 10px,
      10px 10px;
    position: relative;
    overflow: hidden;
  }

  .cyber-file-bg-glitch {
    background-color: var(--color-background-900);
    position: relative;
    overflow: hidden;
  }

  .cyber-file-bg-glitch::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: linear-gradient(
      90deg,
      rgba(var(--color-error-rgb), 0) 0%,
      rgba(var(--color-error-rgb), 0.2) 50%,
      rgba(var(--color-error-rgb), 0) 100%
    );
    opacity: 0.5;
    animation: glitch-scan 2s ease-in-out infinite alternate;
    pointer-events: none;
  }

  .cyber-image-container img {
    width: 100%;
    height: 100%;
    transition: opacity var(--transition-slow) var(--ease-in-out);
  }

  .cyber-image-maintain-ratio {
    display: flex;
    align-items: center;
    justify-content: center;
    backdrop-filter: blur(1px);
  }

  .cyber-image-maintain-ratio img {
    max-width: 100%;
    max-height: 100%;
    width: auto;
    height: auto;
    position: relative;
    z-index: 1;
    box-shadow:
      0 0 15px var(--color-overlay-medium),
      0 0 3px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .cyber-image-rounded {
    border-radius: var(--radius-md);
  }

  .cyber-image-cover img {
    object-fit: cover;
  }

  .cyber-image-contain img {
    object-fit: contain;
  }

  .cyber-image-fill img {
    object-fit: fill;
  }

  .cyber-image-placeholder {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    background-color: var(--color-overlay-medium);
    z-index: 1;
    width: 100%;
    height: 100%;
    box-sizing: border-box;
  }

  .cyber-image-error-text {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    margin-top: var(--space-sm);
  }

  .cyber-image-error {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--color-error-500);
    font-size: var(--text-base);
    gap: var(--space-sm);
  }

  .cyber-image-overlay {
    position: absolute;
    inset: 0;
    z-index: 2;
    pointer-events: none;
  }

  .cyber-image-nsfw-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background-color: var(--color-overlay-heavy);
    backdrop-filter: var(--backdrop-blur-sm);
    z-index: 3;
    padding: var(--space-md);
  }

  .cyber-image-nsfw-icon {
    color: var(--color-error-500);
    font-size: var(--text-2xl);
    margin-bottom: var(--space-sm);
  }

  .cyber-image-nsfw-text {
    color: var(--color-content-default);
    font-size: var(--text-sm);
    text-align: center;
    margin: 0 0 var(--space-sm);
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
    max-width: calc(100% - var(--space-sm));
  }

  .cyber-image-nsfw-button {
    background: linear-gradient(to right, var(--color-brand-500), var(--color-brand-500));
    color: var(--color-text-on-brand);
    border: none;
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-md);
    font-size: var(--text-xs);
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    position: relative;
    overflow: hidden;
    transition: all var(--transition-slow) var(--ease-in-out);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.5);
    min-width: 0;
    white-space: nowrap;
  }

  .cyber-image-nsfw-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.7);
  }

  .cyber-image-nsfw-button:active {
    transform: translateY(1px);
  }

  .cyber-image-nsfw-button-icon {
    font-size: var(--text-xs);
  }

  .cyber-image-dimensions {
    position: absolute;
    bottom: var(--space-sm);
    right: var(--space-sm);
    background-color: var(--color-overlay-heavy);
    color: var(--color-content-default);
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-md);
    font-size: var(--text-xs);
    z-index: 3;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.6);
    text-shadow: 0 0 3px var(--color-brand-500);
  }

  .cyber-image-dimension-text {
    font-family: var(--font-mono);
    letter-spacing: var(--tracking-wide);
  }

  @keyframes cyber-scan {
    0% {
      transform: translateY(-100%);
    }
    100% {
      transform: translateY(200%);
    }
  }

  @keyframes cyber-scan-vert {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(200%);
    }
  }

  @keyframes circuit-animate {
    0% {
      transform: translate(0, 0);
    }
    100% {
      transform: translate(20px, 20px);
    }
  }

  @keyframes glitch-animate {
    0% {
      transform: translateX(0);
    }
    5% {
      transform: translateX(-2px);
    }
    10% {
      transform: translateX(1px);
    }
    15% {
      transform: translateX(-1px);
    }
    20% {
      transform: translateX(0);
    }
    100% {
      transform: translateX(0);
    }
  }

  @keyframes glitch-scan {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(100%);
    }
  }

  @keyframes cyber-button-border {
    0% {
      background-position: 0% 50%;
    }
    50% {
      background-position: 100% 50%;
    }
    100% {
      background-position: 0% 50%;
    }
  }

  @media (max-width: 768px) {
    .cyber-file-type-icon {
      font-size: var(--text-3xl);
    }

    .cyber-file-type-name {
      font-size: var(--text-xs);
    }

    .cyber-file-name {
      font-size: 0.7rem;
    }
  }

  .cyber-file-container:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 20px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .file-type-video .cyber-file-container:hover {
    box-shadow: 0 4px 20px rgba(var(--color-error-rgb), 0.3);
  }

  .file-type-document .cyber-file-container:hover {
    box-shadow: 0 4px 20px rgba(var(--color-success-rgb), 0.3);
  }

  .file-type-archive .cyber-file-container:hover {
    box-shadow: 0 4px 20px rgba(var(--color-warning-rgb), 0.3);
  }

  .file-type-audio .cyber-file-container:hover {
    box-shadow: 0 4px 20px rgba(var(--color-error-rgb), 0.4);
  }

  .file-type-other .cyber-file-container:hover {
    box-shadow: 0 4px 20px rgba(var(--color-content-rgb), 0.15);
  }
</style>
