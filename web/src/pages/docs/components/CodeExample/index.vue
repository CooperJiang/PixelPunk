<script setup lang="ts">
  import { computed, onMounted, ref, useSlots } from 'vue'
  import type { CodeExampleProps } from './types'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'CodeExample',
  })

  const props = withDefaults(defineProps<CodeExampleProps>(), {
    language: 'javascript',
    marginBottom: false,
  })

  const slots = useSlots()
  const codeElement = ref()
  const rawCode = ref('')
  const copied = ref(false)

  onMounted(() => {
    /* 获取slot中的原始代码 */
    if (slots.default) {
      const slotContent = slots.default()
      if (slotContent && slotContent[0] && slotContent[0].children) {
        rawCode.value = slotContent[0].children.toString()
      }
    }
  })

  const cleanCode = computed(() => {
    if (!rawCode.value) {
      return ''
    }

    /* 清理代码：移除首尾空行，保持缩进 */
    const lines = rawCode.value.split('\n')

    /* 移除开头的空行 */
    while (lines.length > 0 && lines[0].trim() === '') {
      lines.shift()
    }

    /* 移除结尾的空行 */
    while (lines.length > 0 && lines[lines.length - 1].trim() === '') {
      lines.pop()
    }

    let code = lines.join('\n')

    /* 如果是JSON格式，尝试格式化 */
    if (props.language === 'json') {
      try {
        const parsed = JSON.parse(code)
        code = JSON.stringify(parsed, null, 2)
      } catch (_e) {
        /* 如果解析失败，保持原样 */
        /* JSON解析失败，保持原始格式 */
      }
    }

    return code
  })

  const highlightedCode = computed(() => {
    const code = cleanCode.value
    if (!code) {
      return ''
    }

    /* 先进行HTML转义 */
    const escaped = code.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

    /* 根据语言类型应用语法高亮 */
    switch (props.language) {
      case 'json':
        return highlightJson(escaped)
      case 'javascript':
        return highlightJavaScript(escaped)
      case 'bash':
        return highlightBash(escaped)
      case 'python':
        return highlightPython(escaped)
      case 'go':
        return highlightGo(escaped)
      default:
        return escaped
    }
  })

  const highlightJson = (code: string) =>
    code
      .replace(/"([^"]+)":/g, '<span class="token-key">"$1":</span>')
      .replace(/:\s*"([^"]*)"/g, ': <span class="token-string">"$1"</span>')
      .replace(/:\s*(\d+)/g, ': <span class="token-number">$1</span>')
      .replace(/:\s*(true|false|null)/g, ': <span class="token-boolean">$1</span>')

  const highlightJavaScript = (code: string) =>
    code
      .replace(
        /\b(async|function|const|let|var|if|else|for|while|return|try|catch|throw|import|export|from|class|new|this|super|await)\b/g,
        '<span class="token-keyword">$1</span>'
      )
      .replace(/'([^']*)'/g, '<span class="token-string">\'$1\'</span>')
      .replace(/"([^"]*)"/g, '<span class="token-string">"$1"</span>')
      .replace(/`([^`]*)`/g, '<span class="token-template">`$1`</span>')
      .replace(/\/\/.*$/gm, '<span class="token-comment">$&</span>')
      .replace(/\b(\d+)\b/g, '<span class="token-number">$1</span>')
      .replace(/\b(console|document|window|navigator)\b/g, '<span class="token-builtin">$1</span>')

  const highlightBash = (code: string) =>
    code
      .replace(/^#.*$/gm, '<span class="token-comment">$&</span>')
      .replace(
        /\b(curl|npm|git|cd|ls|mkdir|rm|cp|mv|cat|grep|find|chmod|sudo|wget|tar|zip|unzip)\b/g,
        '<span class="token-command">$1</span>'
      )
      .replace(/-{1,2}[a-zA-Z-]+/g, '<span class="token-flag">$&</span>')
      .replace(/"([^"]*)"/g, '<span class="token-string">"$1"</span>')
      .replace(/'([^']*)'/g, '<span class="token-string">\'$1\'</span>')

  const highlightPython = (code: string) =>
    code
      .replace(
        /\b(def|class|import|from|if|elif|else|for|while|try|except|finally|with|as|return|yield|break|continue|pass|lambda|and|or|not|in|is|None|True|False|self)\b/g,
        '<span class="token-keyword">$1</span>'
      )
      .replace(/#.*$/gm, '<span class="token-comment">$&</span>')
      .replace(/"""[\s\S]*?"""/g, '<span class="token-docstring">$&</span>')
      .replace(/"([^"]*)"/g, '<span class="token-string">"$1"</span>')
      .replace(/'([^']*)'/g, '<span class="token-string">\'$1\'</span>')
      .replace(/\b(\d+\.?\d*)\b/g, '<span class="token-number">$1</span>')
      .replace(
        /\b(print|len|range|enumerate|zip|map|filter|str|int|float|list|dict|set|tuple)\b/g,
        '<span class="token-builtin">$1</span>'
      )

  const highlightGo = (code: string) =>
    code
      .replace(
        /\b(package|import|func|var|const|type|struct|interface|if|else|for|range|switch|case|default|return|defer|go|chan|select|break|continue|fallthrough|goto)\b/g,
        '<span class="token-keyword">$1</span>'
      )
      .replace(/\/\/.*$/gm, '<span class="token-comment">$&</span>')
      .replace(/\/\*[\s\S]*?\*\//g, '<span class="token-comment">$&</span>')
      .replace(/"([^"]*)"/g, '<span class="token-string">"$1"</span>')
      .replace(/`([^`]*)`/g, '<span class="token-string">`$1`</span>')
      .replace(/\b(\d+)\b/g, '<span class="token-number">$1</span>')
      .replace(
        /\b(string|int|int64|float64|bool|byte|rune|error|fmt|http|json|time|os|io)\b/g,
        '<span class="token-type">$1</span>'
      )

  const copyCode = () => {
    const code = cleanCode.value
    if (code) {
      navigator.clipboard
        .writeText(code)
        .then(() => {
          copied.value = true
          /* 代码已复制到剪贴板 */

          /* 2秒后重置状态 */
          setTimeout(() => {
            copied.value = false
          }, 2000)
        })
        .catch((err) => {
          console.error('Copy failed:', err)
        })
    }
  }
</script>

<template>
  <div class="code-example">
    <div class="code-header">
      <span class="code-title">{{ title }}</span>
      <button class="copy-btn" :title="$t('docs.examples.copy')" :class="{ copied: copied }" @click="copyCode">
        <i :class="copied ? 'fas fa-check' : 'fas fa-copy'" />
        <span class="copy-text">{{ copied ? $t('docs.examples.copied') : $t('docs.examples.copy') }}</span>
      </button>
    </div>
    <div class="code-container">
      <pre
        class="code-block"
        :class="{ 'mb-4': marginBottom }"
      ><code ref="codeElement" :class="`language-${language}`" v-html="highlightedCode"/></pre>
    </div>
  </div>
</template>

<style scoped>
  .code-example {
    background: linear-gradient(
      145deg,
      rgba(var(--color-background-900-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.98) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    margin-bottom: 1.5rem;
    overflow: hidden;
    box-shadow: none;
    transition: all 0.3s ease;
  }

  .code-example:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow:
      0 12px 32px rgba(0, 0, 0, 0.5),
      0 0 25px rgba(var(--color-brand-500-rgb), 0.3),
      inset 0 1px 0 rgba(var(--color-content-rgb), 0.1);
    transform: translateY(-2px);
  }

  .code-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    background: linear-gradient(
      90deg,
      rgba(var(--color-brand-500-rgb), 0.08) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.04) 100%
    );
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .code-title {
    font-weight: 600;
    color: var(--color-brand-500);
    font-size: 0.875rem;
  }

  .copy-btn {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    padding: 0.375rem 0.75rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .copy-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.25);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .copy-btn.copied {
    background: rgba(var(--color-success-rgb), 0.2);
    border-color: rgba(var(--color-success-rgb), 0.4);
    color: var(--color-success-500);
  }

  .code-container {
    max-height: 400px;
    overflow: auto;
    background: rgba(var(--color-background-800-rgb), 0.4);
  }

  .code-block {
    background: transparent;
    color: var(--color-text-content);
    padding: 1.5rem;
    margin: 0;
    font-family: 'SF Mono', Monaco, monospace;
    font-size: 0.85rem;
    line-height: 1.6;
    white-space: pre;
    overflow: auto;
  }

  :deep(.token-keyword) {
    color: var(--color-error-400);
    font-weight: 600;
  }

  :deep(.token-string) {
    color: var(--color-info-300);
  }

  :deep(.token-template) {
    color: var(--color-info-300);
  }

  :deep(.token-number) {
    color: var(--color-info-400);
    font-weight: 500;
  }

  :deep(.token-boolean) {
    color: var(--color-info-400);
    font-weight: 500;
  }

  :deep(.token-comment) {
    color: var(--color-slate-500);
    font-style: italic;
    opacity: 0.8;
  }

  :deep(.token-key) {
    color: var(--color-info-400);
    font-weight: 600;
  }

  :deep(.token-builtin) {
    color: var(--color-brand-300);
    font-weight: 500;
  }

  :deep(.token-command) {
    color: var(--color-success-500);
    font-weight: 600;
  }

  :deep(.token-flag) {
    color: var(--color-warning-300);
    font-weight: 500;
  }

  :deep(.token-type) {
    color: var(--color-brand-300);
    font-weight: 500;
  }

  :deep(.token-docstring) {
    color: var(--color-info-300);
    font-style: italic;
  }

  .code-container::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }

  .code-container::-webkit-scrollbar-track {
    background: var(--color-hover-bg-neutral);
    border-radius: var(--radius-sm);
  }

  .code-container::-webkit-scrollbar-thumb {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.4) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.4) 100%
    );
    border-radius: var(--radius-sm);
    transition: background 0.3s ease;
  }

  .code-container::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.6) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.6) 100%
    );
  }

  @media (max-width: 768px) {
    .code-block {
      font-size: 0.8rem;
      padding: 1rem;
    }

    .code-header {
      padding: 0.6rem 1rem;
    }

    .code-title {
      font-size: 0.8rem;
    }

    .copy-btn {
      font-size: 0.7rem;
      padding: 0.3rem 0.6rem;
    }
  }
</style>
