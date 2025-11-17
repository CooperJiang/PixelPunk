<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed } from 'vue'
  import { marked } from 'marked'
  import DOMPurify from 'dompurify'
  import type { AnnouncementDetail } from '@/api/types/announcement'
  import CyborgCharacter from '@/components/CyborgCharacter/index.vue'
  import { useTexts } from '@/composables/useTexts'

  /* 配置 marked */
  marked.setOptions({
    gfm: true,
    breaks: true,
  })

  const { $t } = useTexts()

  interface Props {
    modelValue: boolean
    announcement: AnnouncementDetail | null
    loading?: boolean
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'close'): void
    (e: 'neverShowAgain', announcementId: number): void
    (e: 'viewMore'): void
    (e: 'dismissTemporarily', announcementId: number): void
  }

  const props = withDefaults(defineProps<Props>(), {
    loading: false,
  })

  const emit = defineEmits<Emits>()

  const visible = computed({
    get: () => props.modelValue,
    set: (value: boolean) => {
      emit('update:modelValue', value)
      if (!value) {
        emit('close')
      }
    },
  })

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '-'
    try {
      const date = new Date(dateStr)
      if (isNaN(date.getTime())) return '-'
      return date.toLocaleDateString(getCurrentLocale(), {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      })
    } catch {
      return '-'
    }
  }

  const handleClose = () => {
    if (props.announcement) {
      emit('dismissTemporarily', props.announcement.id)
    }
    visible.value = false
  }

  const handleNeverShowAgain = () => {
    if (props.announcement) {
      emit('neverShowAgain', props.announcement.id)
    }
    visible.value = false
  }

  const handleViewMore = () => {
    emit('viewMore')
    visible.value = false
  }

  const handleMaskClick = () => {
    if (props.announcement) {
      emit('dismissTemporarily', props.announcement.id)
    }
    visible.value = false
  }

  const renderedContent = computed(() => {
    if (!props.announcement?.content) return ''
    try {
      const html = marked.parse(props.announcement.content) as string
      return DOMPurify.sanitize(html, {
        ALLOWED_TAGS: [
          'h1',
          'h2',
          'h3',
          'h4',
          'h5',
          'h6',
          'p',
          'br',
          'hr',
          'ul',
          'ol',
          'li',
          'a',
          'strong',
          'em',
          'u',
          's',
          'code',
          'pre',
          'blockquote',
          'img',
          'table',
          'thead',
          'tbody',
          'tr',
          'th',
          'td',
          'div',
          'span',
          'b',
          'i',
          'small',
          'mark',
          'del',
          'ins',
          'sub',
          'sup',
        ],
        ALLOWED_ATTR: ['href', 'src', 'alt', 'title', 'class', 'style', 'width', 'height', 'id', 'name'],
        ALLOW_DATA_ATTR: false,
      })
    } catch (error) {
      return props.announcement.content
    }
  })
</script>

<template>
  <Teleport to="body">
    <Transition name="announcement-dialog">
      <div v-if="visible" class="announcement-dialog-overlay" @click="handleMaskClick">
        <div class="announcement-dialog-container" @click.stop>
          <CyborgCharacter size="small" position="top-left" />

          <div class="corner-border top-left" />
          <div class="corner-border top-right" />
          <div class="corner-border bottom-left" />
          <div class="corner-border bottom-right" />

          <div class="corner-decoration top-left" />
          <div class="corner-decoration top-right" />
          <div class="corner-decoration bottom-left" />
          <div class="corner-decoration bottom-right" />

          <div class="scan-line" />

          <div class="dialog-border-top" />
          <div class="dialog-border-bottom" />

          <div class="dialog-header">
            <div class="header-icon">
              <i class="fas fa-bullhorn" />
            </div>
            <h2 class="title-text">{{ announcement?.title || $t('components.announcement.dialog.title') }}</h2>
            <button class="close-button" @click="handleClose">
              <i class="fas fa-times" />
            </button>
          </div>

          <div class="dialog-body">
            <div v-if="loading" class="loading-container">
              <div class="loading-spinner" />
              <p class="loading-text">{{ $t('components.announcement.dialog.loading') }}</p>
            </div>

            <div v-else-if="announcement" class="announcement-content">
              <div class="meta-info">
                <div class="meta-item">
                  <i class="fas fa-calendar-alt" />
                  <span
                    >{{ $t('components.announcement.dialog.meta.published') }}：{{ formatDate(announcement.created_at) }}</span
                  >
                </div>
                <div class="meta-item">
                  <i class="fas fa-edit" />
                  <span>{{ $t('components.announcement.dialog.meta.updated') }}：{{ formatDate(announcement.updated_at) }}</span>
                </div>
              </div>

              <div class="divider" />

              <div class="content-body markdown-body" v-html="renderedContent" />
            </div>
          </div>

          <div class="dialog-footer">
            <CyberButton v-if="announcement?.is_pinned" type="secondary" size="small" @click="handleNeverShowAgain">
              <i class="fas fa-eye-slash mr-1" />
              {{ $t('components.announcement.dialog.actions.neverShowAgain') }}
            </CyberButton>

            <CyberButton v-else type="secondary" size="small" @click="handleViewMore">
              <i class="fas fa-list mr-1" />
              {{ $t('components.announcement.dialog.actions.viewMore') }}
            </CyberButton>
            <CyberButton type="primary" size="small" @click="handleClose">
              <i class="fas fa-check mr-1" />
              {{ $t('components.announcement.dialog.actions.gotIt') }}
            </CyberButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .announcement-dialog-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(3px);
    z-index: 10000;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .announcement-dialog-container {
    position: relative;
    width: 100%;
    max-width: 800px;
    max-height: 85vh;
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.98), rgba(var(--color-background-900-rgb), 0.95));
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-md);
    box-shadow:
      0 25px 50px rgba(0, 0, 0, 0.5),
      0 0 40px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
    display: flex;
    flex-direction: column;
    overflow: visible;
  }

  .corner-border {
    position: absolute;
    width: 18px;
    height: 18px;
    z-index: 8;
    pointer-events: none;
  }

  .corner-border.top-left {
    top: -2px;
    left: -2px;
    border-top: 2px solid rgba(var(--color-warning-rgb), 0.6);
    border-left: 2px solid rgba(var(--color-warning-rgb), 0.6);
    border-top-left-radius: var(--radius-md);
  }

  .corner-border.top-right {
    top: -2px;
    right: -2px;
    border-top: 2px solid rgba(var(--color-warning-rgb), 0.6);
    border-right: 2px solid rgba(var(--color-warning-rgb), 0.6);
    border-top-right-radius: var(--radius-md);
  }

  .corner-border.bottom-left {
    bottom: -2px;
    left: -2px;
    border-bottom: 2px solid rgba(var(--color-warning-rgb), 0.6);
    border-left: 2px solid rgba(var(--color-warning-rgb), 0.6);
    border-bottom-left-radius: var(--radius-md);
  }

  .corner-border.bottom-right {
    bottom: -2px;
    right: -2px;
    border-bottom: 2px solid rgba(var(--color-warning-rgb), 0.6);
    border-right: 2px solid rgba(var(--color-warning-rgb), 0.6);
    border-bottom-right-radius: var(--radius-md);
  }

  .corner-decoration {
    position: absolute;
    opacity: 0.8;
    z-index: 10;
    pointer-events: none;
  }

  .corner-decoration::before,
  .corner-decoration::after {
    content: '';
    position: absolute;
    background: var(--color-brand-500);
  }

  .corner-decoration.top-left {
    top: -10px;
    left: -10px;
  }

  .corner-decoration.top-left::before {
    width: 16px;
    height: 2px;
    top: 0;
    left: 0;
  }

  .corner-decoration.top-left::after {
    width: 2px;
    height: 16px;
    top: 0;
    left: 0;
  }

  .corner-decoration.top-right {
    top: -10px;
    right: -10px;
  }

  .corner-decoration.top-right::before {
    width: 16px;
    height: 2px;
    top: 0;
    right: 0;
  }

  .corner-decoration.top-right::after {
    width: 2px;
    height: 16px;
    top: 0;
    right: 0;
  }

  .corner-decoration.bottom-left {
    bottom: -10px;
    left: -10px;
  }

  .corner-decoration.bottom-left::before {
    width: 16px;
    height: 2px;
    bottom: 0;
    left: 0;
  }

  .corner-decoration.bottom-left::after {
    width: 2px;
    height: 16px;
    bottom: 0;
    left: 0;
  }

  .corner-decoration.bottom-right {
    bottom: -10px;
    right: -10px;
  }

  .corner-decoration.bottom-right::before {
    width: 16px;
    height: 2px;
    bottom: 0;
    right: 0;
  }

  .corner-decoration.bottom-right::after {
    width: 2px;
    height: 16px;
    bottom: 0;
    right: 0;
  }

  .scan-line {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 2px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.8), transparent);
    opacity: 0.6;
    animation: scan-animation 3s linear infinite;
    pointer-events: none;
    z-index: 5;
  }

  @keyframes scan-animation {
    0% {
      transform: translateY(0);
    }
    100% {
      transform: translateY(85vh);
    }
  }

  .dialog-border-top,
  .dialog-border-bottom {
    position: absolute;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(
      90deg,
      transparent,
      var(--color-brand-500),
      var(--color-brand-400),
      var(--color-brand-500),
      transparent
    );
    opacity: 0.6;
    animation: border-flow 3s linear infinite;
  }

  .dialog-border-top {
    top: 0;
  }

  .dialog-border-bottom {
    bottom: 0;
    animation-delay: 1.5s;
  }

  @keyframes border-flow {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(100%);
    }
  }

  .dialog-header {
    position: relative;
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    padding: var(--space-sm) var(--space-md);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-background-700-rgb), 0.6));
  }

  .header-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-600-rgb), 0.1));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    border-radius: var(--radius-md);
    color: var(--color-brand-400);
    font-size: 14px;
    flex-shrink: 0;
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .header-icon i {
    filter: drop-shadow(0 0 4px rgba(var(--color-brand-500-rgb), 0.6));
  }

  .title-text {
    flex: 1;
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    color: var(--color-content-heading);
    line-height: var(--leading-tight);
    margin: 0;
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .close-button {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(var(--color-background-700-rgb), 0.6);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all var(--transition-normal) ease;
    flex-shrink: 0;
    font-size: 12px;
  }

  .close-button:hover {
    background: rgba(var(--color-error-rgb), 0.15);
    border-color: rgba(var(--color-error-rgb), 0.5);
    color: var(--color-error-400);
    box-shadow: 0 0 8px rgba(var(--color-error-rgb), 0.4);
  }

  .dialog-body {
    flex: 1;
    overflow-y: auto;
    padding: var(--space-md);
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.02) 0%,
      transparent 10%,
      transparent 90%,
      rgba(var(--color-brand-500-rgb), 0.02) 100%
    );
  }

  .dialog-body::-webkit-scrollbar {
    width: 6px;
  }

  .dialog-body::-webkit-scrollbar-track {
    background: rgba(var(--color-background-900-rgb), 0.4);
    border-radius: var(--radius-full);
  }

  .dialog-body::-webkit-scrollbar-thumb {
    background: linear-gradient(180deg, rgba(var(--color-brand-500-rgb), 0.6), rgba(var(--color-brand-600-rgb), 0.4));
    border-radius: var(--radius-full);
  }

  .dialog-body::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(180deg, rgba(var(--color-brand-500-rgb), 0.8), rgba(var(--color-brand-600-rgb), 0.6));
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-2xl);
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-top: 3px solid var(--color-brand-500);
    border-radius: var(--radius-full);
    animation: spin 1s linear infinite;
    margin-bottom: var(--space-md);
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .loading-text {
    color: var(--color-content-muted);
    font-size: var(--text-sm);
  }

  .announcement-content {
    animation: content-fade-in 0.4s ease-out;
  }

  @keyframes content-fade-in {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .meta-info {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-lg);
    margin-bottom: var(--space-md);
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    color: var(--color-content-muted);
    font-size: var(--text-xs);
  }

  .meta-item i {
    color: var(--color-brand-500);
    font-size: 12px;
  }

  .divider {
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.3), transparent);
    margin-bottom: var(--space-md);
  }

  .content-body {
    font-size: var(--text-sm);
    line-height: var(--leading-relaxed);
    color: var(--color-content-default);
  }

  .markdown-body {
    word-wrap: break-word;
  }

  .markdown-body > *:first-child {
    margin-top: 0 !important;
  }

  .markdown-body > *:last-child {
    margin-bottom: 0 !important;
  }

  .markdown-body :deep(h1),
  .markdown-body :deep(h2),
  .markdown-body :deep(h3),
  .markdown-body :deep(h4),
  .markdown-body :deep(h5),
  .markdown-body :deep(h6) {
    margin-top: var(--space-lg);
    margin-bottom: var(--space-sm);
    font-weight: var(--font-semibold);
    color: var(--color-content-heading);
    position: relative;
    padding-left: var(--space-md);
    line-height: 1.4;
  }

  .markdown-body :deep(h1) {
    font-size: var(--text-xl);
  }

  .markdown-body :deep(h2) {
    font-size: var(--text-lg);
  }

  .markdown-body :deep(h3) {
    font-size: var(--text-base);
  }

  .markdown-body :deep(h1)::before,
  .markdown-body :deep(h2)::before,
  .markdown-body :deep(h3)::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0.3em;
    bottom: 0.3em;
    width: 3px;
    background: linear-gradient(180deg, var(--color-brand-500), var(--color-brand-600));
    border-radius: var(--radius-full);
  }

  .markdown-body :deep(p) {
    margin-top: 0;
    margin-bottom: var(--space-sm);
  }

  .markdown-body :deep(ul),
  .markdown-body :deep(ol) {
    padding-left: var(--space-lg);
    margin-top: 0;
    margin-bottom: var(--space-sm);
  }

  .markdown-body :deep(li) {
    margin-bottom: var(--space-xs);
  }

  .markdown-body :deep(li > p) {
    margin-bottom: var(--space-xs);
  }

  .markdown-body :deep(a) {
    color: var(--color-brand-500);
    text-decoration: none;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    transition: all var(--transition-fast) ease;
  }

  .markdown-body :deep(a:hover) {
    color: var(--color-brand-400);
    border-bottom-color: var(--color-brand-400);
  }

  .markdown-body :deep(strong) {
    font-weight: var(--font-semibold);
    color: var(--color-content-heading);
  }

  .markdown-body :deep(em) {
    font-style: italic;
  }

  .markdown-body :deep(code) {
    padding: 2px var(--space-xs);
    background: rgba(var(--color-background-900-rgb), 0.8);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
    font-size: 0.9em;
    color: var(--color-brand-400);
  }

  .markdown-body :deep(pre) {
    padding: var(--space-md);
    background: rgba(var(--color-background-900-rgb), 0.9);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-left: 3px solid var(--color-brand-500);
    border-radius: var(--radius-md);
    overflow-x: auto;
    margin-top: 0;
    margin-bottom: var(--space-sm);
    box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.3);
  }

  .markdown-body :deep(pre code) {
    background: transparent;
    border: none;
    padding: 0;
    color: var(--color-content-default);
  }

  .markdown-body :deep(blockquote) {
    padding: var(--space-sm) var(--space-md);
    margin: var(--space-sm) 0;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-left: 3px solid var(--color-brand-500);
    border-radius: var(--radius-md);
    color: var(--color-content-muted);
  }

  .markdown-body :deep(blockquote > p) {
    margin: 0;
  }

  .markdown-body :deep(hr) {
    height: 1px;
    border: none;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.3), transparent);
    margin: var(--space-md) 0;
  }

  .markdown-body :deep(img) {
    max-width: 100%;
    height: auto;
    border-radius: var(--radius-md);
    margin: var(--space-sm) 0;
  }

  .markdown-body :deep(table) {
    width: 100%;
    border-collapse: collapse;
    margin: var(--space-sm) 0;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  .markdown-body :deep(th),
  .markdown-body :deep(td) {
    padding: var(--space-xs) var(--space-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .markdown-body :deep(th) {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    font-weight: var(--font-semibold);
    color: var(--color-content-heading);
  }

  .markdown-body :deep(tr:nth-child(even)) {
    background: rgba(var(--color-background-900-rgb), 0.3);
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: var(--space-sm);
    padding: var(--space-sm) var(--space-md);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.6), rgba(var(--color-brand-500-rgb), 0.03));
  }

  .announcement-dialog-enter-active,
  .announcement-dialog-leave-active {
    transition: opacity 0.3s ease;
  }

  .announcement-dialog-enter-active .announcement-dialog-container,
  .announcement-dialog-leave-active .announcement-dialog-container {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .announcement-dialog-enter-from,
  .announcement-dialog-leave-to {
    opacity: 0;
  }

  .announcement-dialog-enter-from .announcement-dialog-container,
  .announcement-dialog-leave-to .announcement-dialog-container {
    opacity: 0;
    transform: scale(0.9) translateY(-20px);
  }

  @media (max-width: 768px) {
    .announcement-dialog-overlay {
      padding: var(--space-sm);
    }

    .announcement-dialog-container {
      max-width: 100%;
      max-height: 92vh;
    }

    .dialog-header {
      padding: var(--space-xs) var(--space-sm);
      gap: var(--space-xs);
    }

    .header-icon {
      width: 24px;
      height: 24px;
      font-size: 12px;
    }

    .title-text {
      font-size: var(--text-sm);
    }

    .close-button {
      width: 20px;
      height: 20px;
      font-size: 10px;
    }

    .dialog-body {
      padding: var(--space-sm);
    }

    .dialog-footer {
      padding: var(--space-xs) var(--space-sm);
      flex-wrap: wrap;
    }

    .corner-border {
      width: 18px;
      height: 18px;
    }

    .corner-decoration::before {
      width: 15px !important;
    }

    .corner-decoration::after {
      height: 15px !important;
    }
  }
</style>
