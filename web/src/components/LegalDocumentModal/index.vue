<script setup lang="ts">
  import { computed } from 'vue'
  import { marked } from 'marked'
  import DOMPurify from 'dompurify'
  import { useTexts } from '@/composables/useTexts'

  /* 配置 Markdown 渲染选项 */
  marked.setOptions({
    gfm: true, // GitHub Flavored Markdown
    breaks: true, // 将换行符转换为 <br>
  })

  defineOptions({
    name: 'LegalDocumentModal',
  })

  const { $t } = useTexts()

  interface Props {
    modelValue: boolean
    type: 'privacy' | 'terms'
    content: string
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    content: '',
  })

  const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
  }>()

  /* 渲染 Markdown 内容并清理 HTML */
  const renderedContent = computed(() => {
    if (!props.content) {
      return `<p class="empty-notice">${$t('components.legalDocumentModal.noContent')}</p>`
    }
    try {
      const html = marked.parse(props.content) as string
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
          'strong',
          'em',
          'u',
          's',
          'a',
          'ul',
          'ol',
          'li',
          'blockquote',
          'code',
          'pre',
          'img',
          'table',
          'thead',
          'tbody',
          'tr',
          'th',
          'td',
        ],
        ALLOWED_ATTR: ['href', 'src', 'alt', 'title', 'class'],
      })
    } catch (error) {
      console.error('Markdown render failed:', error)
      return props.content
    }
  })

  const close = () => {
    emit('update:modelValue', false)
  }

  const handleBackdropClick = (event: MouseEvent) => {
    if (event.target === event.currentTarget) {
      close()
    }
  }
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="modelValue" class="legal-document-modal-backdrop" @click="handleBackdropClick">
        <div class="legal-document-modal" @click.stop>
          <button type="button" class="close-button" @click="close">
            <i class="fas fa-times" />
          </button>

          <div class="modal-body">
            <div class="document-content" v-html="renderedContent" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
  .legal-document-modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: var(--color-overlay-dark);
    backdrop-filter: var(--backdrop-blur-md);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    padding: var(--space-xl);
  }

  .legal-document-modal {
    position: relative;
    background: var(--color-background-800);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    box-shadow:
      var(--shadow-2xl),
      0 0 40px rgba(var(--color-brand-500-rgb), 0.2);
    max-width: 900px;
    width: 100%;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .close-button {
    position: absolute;
    top: var(--space-md);
    right: var(--space-md);
    width: 32px;
    height: 32px;
    border-radius: var(--radius-md);
    background: var(--color-background-700);
    border: 1px solid var(--color-border-subtle);
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all var(--transition-fast);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10;
    backdrop-filter: var(--backdrop-blur-sm);

    &:hover {
      background: var(--color-background-600);
      border-color: var(--color-brand-500);
      color: var(--color-brand-400);
      transform: scale(1.05);
    }

    i {
      font-size: var(--text-base);
    }
  }

  .modal-body {
    flex: 1;
    overflow-y: auto;
    padding: var(--space-2xl) var(--space-xl);
  }

  .document-content {
    color: var(--color-content-default);
    line-height: 1.8;
    font-size: var(--text-base);

    :deep(h1),
    :deep(h2),
    :deep(h3),
    :deep(h4),
    :deep(h5),
    :deep(h6) {
      color: var(--color-content-heading);
      font-weight: var(--font-semibold);
      margin-top: var(--space-lg);
      margin-bottom: var(--space-md);
      line-height: 1.5;
      position: relative;
      padding-left: var(--space-md);
    }

    :deep(h1) {
      font-size: var(--text-xl);
      margin-top: 0;
      padding: var(--space-sm) var(--space-md);
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-brand-500-rgb), 0.05));
      border-left: 4px solid var(--color-brand-500);
      border-radius: 0 var(--radius-md) var(--radius-md) 0;
      box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.1);
      text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.2);
    }

    :deep(h2) {
      font-size: var(--text-lg);
      padding-left: var(--space-md);
      border-left: 3px solid var(--color-brand-400);
      background: linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.08), transparent);
      padding: var(--space-xs) var(--space-md);
      border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
    }

    :deep(h3) {
      font-size: var(--text-base);
      color: var(--color-brand-400);
      border-left: 2px solid var(--color-brand-500);
      padding-left: var(--space-sm);
    }

    :deep(h4) {
      font-size: var(--text-base);
      color: var(--color-content-default);
      font-weight: var(--font-medium);
    }

    :deep(h5),
    :deep(h6) {
      font-size: var(--text-sm);
      color: var(--color-content-muted);
      font-weight: var(--font-medium);
    }

    :deep(p) {
      margin-bottom: var(--space-md);
      color: var(--color-content-default);
      line-height: 1.9;
      text-align: justify;
    }

    :deep(p:first-of-type) {
      margin-top: var(--space-sm);
    }

    :deep(a) {
      color: var(--color-brand-400);
      text-decoration: none;
      position: relative;
      padding-bottom: 2px;
      border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
      transition: all var(--transition-fast);

      &:hover {
        color: var(--color-brand-300);
        border-bottom-color: var(--color-brand-400);
        text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.6);
      }
    }

    :deep(ul),
    :deep(ol) {
      margin-bottom: var(--space-lg);
      padding-left: var(--space-xl);
      line-height: 1.9;
      background: rgba(var(--color-background-700-rgb), 0.3);
      padding: var(--space-md) var(--space-md) var(--space-md) var(--space-2xl);
      border-radius: var(--radius-md);
      border-left: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
    }

    :deep(ul) {
      list-style-type: none;
    }

    :deep(ul > li)::before {
      content: '▹';
      color: var(--color-brand-400);
      font-weight: bold;
      display: inline-block;
      width: 1em;
      margin-left: -1em;
      text-shadow: 0 0 4px rgba(var(--color-brand-500-rgb), 0.5);
    }

    :deep(ol) {
      list-style-type: none;
      counter-reset: item;
    }

    :deep(ol > li) {
      counter-increment: item;
      position: relative;
    }

    :deep(ol > li)::before {
      content: counter(item) '.';
      color: var(--color-brand-400);
      font-weight: var(--font-semibold);
      display: inline-block;
      width: 2em;
      margin-left: -2em;
      text-align: right;
      padding-right: 0.5em;
    }

    :deep(li) {
      margin-bottom: var(--space-sm);
      color: var(--color-content-default);
      padding-left: var(--space-sm);
      position: relative;
    }

    :deep(ul ul),
    :deep(ol ol),
    :deep(ul ol),
    :deep(ol ul) {
      margin-top: var(--space-sm);
      margin-bottom: var(--space-sm);
      background: transparent;
      border-left: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    }

    :deep(strong),
    :deep(b) {
      color: var(--color-brand-400);
      font-weight: var(--font-bold);
      text-shadow: 0 0 6px rgba(var(--color-brand-500-rgb), 0.2);
    }

    :deep(em),
    :deep(i) {
      font-style: italic;
      color: var(--color-brand-300);
    }

    :deep(u) {
      text-decoration: underline;
      text-decoration-color: var(--color-brand-400);
      text-decoration-thickness: 2px;
      text-underline-offset: 3px;
    }

    :deep(s),
    :deep(strike) {
      text-decoration: line-through;
      color: var(--color-content-muted);
      opacity: 0.7;
    }

    :deep(code) {
      background: rgba(var(--color-brand-500-rgb), 0.1);
      padding: 3px 8px;
      border-radius: var(--radius-sm);
      font-family: var(--font-mono);
      font-size: 0.9em;
      color: var(--color-brand-400);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
      box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.1);
    }

    :deep(pre) {
      background: linear-gradient(135deg, var(--color-background-900), var(--color-background-800));
      padding: var(--space-lg);
      border-radius: var(--radius-md);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
      overflow-x: auto;
      margin-bottom: var(--space-lg);
      box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.3);
      position: relative;

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: linear-gradient(90deg, var(--color-brand-500), transparent);
      }

      code {
        background: transparent;
        padding: 0;
        border: none;
        font-size: var(--text-sm);
        color: var(--color-brand-300);
        box-shadow: none;
      }
    }

    :deep(blockquote) {
      border-left: 4px solid var(--color-brand-500);
      padding: var(--space-md) var(--space-lg);
      margin: var(--space-lg) 0;
      background: linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-background-700-rgb), 0.5));
      border-radius: 0 var(--radius-md) var(--radius-md) 0;
      color: var(--color-content-default);
      position: relative;
      box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.1);

      &::before {
        content: '"';
        position: absolute;
        top: var(--space-xs);
        left: var(--space-sm);
        font-size: 2rem;
        color: rgba(var(--color-brand-500-rgb), 0.3);
        font-family: Georgia, serif;
        line-height: 1;
      }

      p {
        margin-bottom: var(--space-sm);
        padding-left: var(--space-md);

        &:last-child {
          margin-bottom: 0;
        }
      }

      &:first-child {
        margin-top: 0;
      }

      &:last-child {
        margin-bottom: 0;
      }
    }

    :deep(hr) {
      border: none;
      height: 2px;
      background: linear-gradient(90deg, transparent, var(--color-brand-500), transparent);
      margin: var(--space-2xl) 0;
      position: relative;
      opacity: 0.3;

      &::after {
        content: '◆';
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
        background: var(--color-background-800);
        padding: 0 var(--space-sm);
        color: var(--color-brand-500);
        font-size: 0.75rem;
      }
    }

    :deep(table) {
      width: 100%;
      border-collapse: collapse;
      margin-bottom: var(--space-md);
      border: 1px solid var(--color-border-default);
      border-radius: var(--radius-md);
      overflow: hidden;
    }

    :deep(th),
    :deep(td) {
      padding: var(--space-sm) var(--space-md);
      border: 1px solid var(--color-border-default);
      text-align: left;
    }

    :deep(th) {
      background: var(--color-background-700);
      color: var(--color-content-heading);
      font-weight: var(--font-semibold);
    }

    :deep(tr:hover) {
      background: var(--color-background-700);
    }

    .empty-notice {
      text-align: center;
      color: var(--color-content-muted);
      font-style: italic;
      padding: var(--space-3xl);
      font-size: var(--text-base);
      background: rgba(var(--color-brand-500-rgb), 0.05);
      border-radius: var(--radius-sm);
      border: 1px dashed rgba(var(--color-brand-500-rgb), 0.2);
    }
  }

  .modal-fade-enter-active,
  .modal-fade-leave-active {
    transition: all var(--transition-normal);
  }

  .modal-fade-enter-from,
  .modal-fade-leave-to {
    opacity: 0;
  }

  .modal-fade-enter-active .legal-document-modal,
  .modal-fade-leave-active .legal-document-modal {
    transition: all var(--transition-normal);
  }

  .modal-fade-enter-from .legal-document-modal,
  .modal-fade-leave-to .legal-document-modal {
    transform: scale(0.96) translateY(-10px);
    opacity: 0;
  }

  .modal-body::-webkit-scrollbar {
    width: 8px;
  }

  .modal-body::-webkit-scrollbar-track {
    background: var(--color-background-900);
    border-radius: var(--radius-md);
    margin: var(--space-sm);
  }

  .modal-body::-webkit-scrollbar-thumb {
    background: var(--color-background-600);
    border-radius: var(--radius-md);
    border: 2px solid transparent;
    background-clip: padding-box;
    transition: all var(--transition-fast);
  }

  .modal-body::-webkit-scrollbar-thumb:hover {
    background: var(--color-brand-500);
    background-clip: padding-box;
  }

  @media (max-width: 768px) {
    .legal-document-modal-backdrop {
      padding: var(--space-md);
    }

    .legal-document-modal {
      max-width: 100%;
      max-height: 95vh;
    }

    .modal-body {
      padding: var(--space-xl) var(--space-md);
    }

    .close-button {
      top: var(--space-sm);
      right: var(--space-sm);
      width: 28px;
      height: 28px;

      i {
        font-size: var(--text-sm);
      }
    }
  }
</style>
