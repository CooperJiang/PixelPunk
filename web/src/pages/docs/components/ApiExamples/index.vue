<script setup lang="ts">
  import { codeTabs } from '../../composables/useDocsPage'
  import ApiSection from '../ApiSection/index.vue'
  import CodeTabs from '../CodeTabs/index.vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'ApiExamples',
  })

  const props = defineProps<Props>()

  interface Props {
    codeExamples: any
    copyStatus: Record<string, boolean>
    copyCode: (codeId: string, code: string) => void
    switchCodeTab: (tabId: string) => void
  }

  const onCopyCode = (codeId: string, code: string) => {
    props.copyCode(codeId, code)
  }

  const onSwitchCodeTab = (tabId: string) => {
    props.switchCodeTab(tabId)
  }
</script>

<template>
  <ApiSection id="api-examples" :title="$t('docs.examples.title')" icon-class="fas fa-code">
    <CodeTabs :tabs="codeTabs" default="curl" @update:tab="onSwitchCodeTab">
      <template #curl>
        <div class="code-example">
          <div class="code-header">
            <span class="code-title">{{ $t('docs.examples.singleUpload') }}</span>
            <button
              class="copy-btn"
              :class="{ copied: copyStatus['curl-single'] }"
              @click="onCopyCode('curl-single', codeExamples.curl.single)"
            >
              <i :class="copyStatus['curl-single'] ? 'fas fa-check' : 'fas fa-copy'" />
              <span>{{ copyStatus['curl-single'] ? $t('docs.examples.copied') : $t('docs.examples.copy') }}</span>
            </button>
          </div>
          <div class="code-container">
            <pre class="code-block"><code class="language-bash">{{ codeExamples.curl.single }}</code></pre>
          </div>
        </div>
      </template>

      <template #javascript>
        <div class="code-example">
          <div class="code-header">
            <span class="code-title">{{ $t('docs.examples.singleUpload') }}</span>
            <button
              class="copy-btn"
              :class="{ copied: copyStatus['js-single'] }"
              @click="onCopyCode('js-single', codeExamples.javascript.single)"
            >
              <i :class="copyStatus['js-single'] ? 'fas fa-check' : 'fas fa-copy'" />
              <span>{{ copyStatus['js-single'] ? $t('docs.examples.copied') : $t('docs.examples.copy') }}</span>
            </button>
          </div>
          <div class="code-container">
            <pre class="code-block"><code class="language-javascript">{{ codeExamples.javascript.single }}</code></pre>
          </div>
        </div>
      </template>

      <template #nodejs>
        <div class="code-example">
          <div class="code-header">
            <span class="code-title">{{ $t('docs.examples.singleUpload') }}</span>
            <button
              class="copy-btn"
              :class="{ copied: copyStatus['nodejs-single'] }"
              @click="onCopyCode('nodejs-single', codeExamples.nodejs.single)"
            >
              <i :class="copyStatus['nodejs-single'] ? 'fas fa-check' : 'fas fa-copy'" />
              <span>{{ copyStatus['nodejs-single'] ? $t('docs.examples.copied') : $t('docs.examples.copy') }}</span>
            </button>
          </div>
          <div class="code-container">
            <pre class="code-block"><code class="language-javascript">{{ codeExamples.nodejs.single }}</code></pre>
          </div>
        </div>
      </template>

      <template #python>
        <div class="code-example">
          <div class="code-header">
            <span class="code-title">{{ $t('docs.examples.basicUpload') }}</span>
            <button
              class="copy-btn"
              :class="{ copied: copyStatus['py-basic'] }"
              @click="onCopyCode('py-basic', codeExamples.python.basic)"
            >
              <i :class="copyStatus['py-basic'] ? 'fas fa-check' : 'fas fa-copy'" />
              <span>{{ copyStatus['py-basic'] ? $t('docs.examples.copied') : $t('docs.examples.copy') }}</span>
            </button>
          </div>
          <div class="code-container">
            <pre class="code-block"><code class="language-python">{{ codeExamples.python.basic }}</code></pre>
          </div>
        </div>
      </template>

      <template #java>
        <div class="code-example">
          <div class="code-header">
            <span class="code-title">{{ $t('docs.examples.basicUpload') }}</span>
            <button
              class="copy-btn"
              :class="{ copied: copyStatus['java-basic'] }"
              @click="onCopyCode('java-basic', codeExamples.java.basic)"
            >
              <i :class="copyStatus['java-basic'] ? 'fas fa-check' : 'fas fa-copy'" />
              <span>{{ copyStatus['java-basic'] ? $t('docs.examples.copied') : $t('docs.examples.copy') }}</span>
            </button>
          </div>
          <div class="code-container">
            <pre class="code-block"><code class="language-java">{{ codeExamples.java.basic }}</code></pre>
          </div>
        </div>
      </template>

      <template #go>
        <div class="code-example">
          <div class="code-header">
            <span class="code-title">{{ $t('docs.examples.singleUpload') }}</span>
            <button
              class="copy-btn"
              :class="{ copied: copyStatus['go-single'] }"
              @click="onCopyCode('go-single', codeExamples.go.single)"
            >
              <i :class="copyStatus['go-single'] ? 'fas fa-check' : 'fas fa-copy'" />
              <span>{{ copyStatus['go-single'] ? $t('docs.examples.copied') : $t('docs.examples.copy') }}</span>
            </button>
          </div>
          <div class="code-container">
            <pre class="code-block"><code class="language-go">{{ codeExamples.go.single }}</code></pre>
          </div>
        </div>
      </template>
    </CodeTabs>
  </ApiSection>
</template>

<style scoped>
  .code-example {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.9) 0%,
      rgba(var(--color-background-800-rgb), 0.95) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-top: none;
    border-bottom-left-radius: 8px;
    border-bottom-right-radius: 8px;
    margin-bottom: 1.5rem;
    overflow: hidden;
    position: relative;
  }

  .code-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.08) 0%,
      rgba(var(--color-background-800-rgb), 0.6) 100%
    );
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
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
</style>
