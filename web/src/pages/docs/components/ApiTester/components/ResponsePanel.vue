<script setup lang="ts">
  import { ref, watch, nextTick } from 'vue'
  import hljs from 'highlight.js'
  import { useTexts } from '@/composables/useTexts'
  import type { ApiResponse, UploadedImage } from '../types'

  defineOptions({
    name: 'ResponsePanel',
  })

  interface Props {
    response: ApiResponse | null
    isLoading: boolean
    responseCopied: boolean
    selectedFiles: File[]
    totalFileSize: number
    formattedResponse: string
    responseStatusClass: string
    uploadedImages: UploadedImage[]
    hasErrors: boolean
    oversizedFiles: string[]
    unsupportedFiles: string[]
    invalidFiles: string[]
    uploadErrors: string[]
    sizeLimit: string
    formatFileSize: (bytes: number) => string
  }

  const props = defineProps<Props>()

  defineEmits<{
    'copy-response': []
    'copy-image-url': [url: string]
    'open-image': [url: string]
  }>()

  const { $t } = useTexts()

  /* 代码块引用，用于语法高亮 */
  const jsonCodeElement = ref<HTMLElement | null>(null)

  /* 监听响应变化，高亮JSON代码 */
  watch(
    () => props.formattedResponse,
    () => {
      nextTick(() => {
        if (jsonCodeElement.value && props.formattedResponse) {
          delete jsonCodeElement.value.dataset.highlighted
          hljs.highlightElement(jsonCodeElement.value)
        }
      })
    }
  )

  const handleImageError = (event: Event) => {
    const target = event.target as HTMLImageElement
    target.style.display = 'none'
  }
</script>

<template>
  <div class="response-cyber-panel">
    <div class="section-header">
      <i class="fas fa-terminal"></i>
      <h4>{{ $t('docs.tester.response') }}</h4>
      <div v-if="response" class="status-indicator" :class="responseStatusClass">
        <div class="status-dot"></div>
        <span>{{ response.status }} {{ response.statusText }}</span>
      </div>
      <div class="section-line"></div>
    </div>

    <div class="response-content">
      <div v-if="!response && !isLoading" class="empty-state">
        <div class="empty-container">
          <div class="empty-icon-wrapper">
            <i class="fas fa-flask"></i>
          </div>
          <div class="empty-text-content">
            <h3 class="empty-title">{{ $t('docs.response.emptyTitle') }}</h3>
            <p class="empty-subtitle">{{ $t('docs.response.emptySubtitle') }}</p>
          </div>
        </div>

        <div class="test-guide">
          <div class="guide-steps">
            <div class="guide-step">
              <span class="step-num">1</span>
              <span class="step-text">{{ $t('docs.response.step1') }}</span>
            </div>
            <div class="step-arrow">→</div>
            <div class="guide-step">
              <span class="step-num">2</span>
              <span class="step-text">{{ $t('docs.response.step2') }}</span>
            </div>
            <div class="step-arrow">→</div>
            <div class="guide-step">
              <span class="step-num">3</span>
              <span class="step-text">{{ $t('docs.response.step3') }}</span>
            </div>
          </div>
          <p class="guide-note">{{ $t('docs.tester.guide') }}</p>
        </div>
      </div>

      <div v-else-if="isLoading" class="loading-state">
        <div class="loading-animation">
          <div class="loading-rings">
            <div class="ring ring-1"></div>
            <div class="ring ring-2"></div>
            <div class="ring-3 ring"></div>
          </div>
        </div>
        <p class="loading-text">{{ $t('docs.response.processing') }}</p>
      </div>

      <div v-else class="response-data">
        <div class="request-summary">
          <div class="summary-item">
            <span class="label">{{ $t('docs.response.fileCount') }}</span>
            <span class="value">{{ selectedFiles.length }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ $t('docs.response.totalSize') }}</span>
            <span class="value">{{ formatFileSize(totalFileSize) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ $t('docs.response.responseTime') }}</span>
            <span class="value">{{ response.duration }}ms</span>
          </div>
        </div>

        <div class="json-response">
          <div class="json-header">
            <span>{{ $t('docs.response.responseData') }}</span>
            <button class="copy-json-btn" @click="$emit('copy-response')" :class="{ copied: responseCopied }">
              <i :class="responseCopied ? 'fas fa-check' : 'fas fa-copy'"></i>
              <span>{{ responseCopied ? $t('docs.tester.copied') : $t('docs.tester.copy') }}</span>
            </button>
          </div>
          <div class="json-container">
            <pre class="json-code"><code ref="jsonCodeElement" class="language-json">{{ formattedResponse }}</code></pre>
          </div>
        </div>

        <div v-if="uploadedImages.length > 0" class="upload-results">
          <h5 class="results-title">
            <i class="fas fa-check-circle"></i>
            {{ $t('docs.response.uploadSuccess', { count: uploadedImages.length }) }}
          </h5>
          <div class="images-grid">
            <div v-for="image in uploadedImages" :key="image.id" class="image-result-card">
              <div class="image-preview">
                <img :src="image.thumb_url" :alt="image.original_name" @error="handleImageError" />
                <div class="image-overlay">
                  <button
                    @click="$emit('copy-image-url', image.url)"
                    class="copy-url-btn"
                    :title="$t('docs.tester.copyImageUrl')"
                  >
                    <i class="fas fa-link"></i>
                  </button>
                  <button
                    @click="$emit('open-image', image.url)"
                    class="view-original-btn"
                    :title="$t('docs.tester.viewOriginal')"
                  >
                    <i class="fas fa-external-link-alt"></i>
                  </button>
                </div>
              </div>
              <div class="image-details">
                <div class="image-header">
                  <p class="image-name">{{ image.display_name || image.original_name }}</p>
                  <span class="image-id">{{ image.id }}</span>
                </div>
                <div class="image-meta-grid">
                  <div class="meta-item">
                    <span class="meta-label">{{ $t('docs.response.dimensions') }}</span>
                    <span class="meta-value">{{ image.width }}×{{ image.height }}</span>
                  </div>
                  <div class="meta-item">
                    <span class="meta-label">{{ $t('docs.response.size') }}</span>
                    <span class="meta-value">{{ formatFileSize(image.size) }}</span>
                  </div>
                  <div class="meta-item">
                    <span class="meta-label">{{ $t('docs.response.format') }}</span>
                    <span class="meta-value format-badge">{{ image.format.toUpperCase() }}</span>
                  </div>
                  <div class="meta-item">
                    <span class="meta-label">{{ $t('docs.response.accessLevel') }}</span>
                    <span class="meta-value access-level-badge" :class="image.access_level">
                      {{
                        image.access_level === 'public'
                          ? $t('docs.params.public')
                          : image.access_level === 'private'
                            ? $t('docs.params.private')
                            : image.access_level === 'protected'
                              ? $t('docs.params.protected')
                              : image.access_level
                      }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="hasErrors" class="error-results">
          <h5 class="results-title error">
            <i class="fas fa-exclamation-triangle"></i>
            {{ $t('docs.response.processingError') }}
          </h5>
          <div v-if="oversizedFiles.length > 0" class="error-section">
            <h6>{{ $t('docs.response.oversizedFiles') }}</h6>
            <ul class="error-list">
              <li v-for="file in oversizedFiles" :key="file">{{ file }}</li>
            </ul>
            <p class="size-limit">{{ $t('docs.response.sizeLimit', { limit: sizeLimit }) }}</p>
          </div>
          <div v-if="unsupportedFiles.length > 0" class="error-section">
            <h6>{{ $t('docs.response.unsupportedFiles') }}</h6>
            <ul class="error-list">
              <li v-for="file in unsupportedFiles" :key="file">{{ file }}</li>
            </ul>
          </div>
          <div v-if="invalidFiles.length > 0" class="error-section">
            <h6>{{ $t('docs.response.invalidFiles') }}</h6>
            <ul class="error-list">
              <li v-for="file in invalidFiles" :key="file">{{ file }}</li>
            </ul>
          </div>
          <div v-if="uploadErrors.length > 0" class="error-section">
            <h6>{{ $t('docs.response.uploadFailed') }}</h6>
            <ul class="error-list">
              <li v-for="error in uploadErrors" :key="error">{{ error }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .response-cyber-panel {
    padding: 1rem;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-700-rgb), 0.6) 0%,
      rgba(var(--color-background-800-rgb), 0.8) 100%
    );
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.75rem;
    padding-bottom: 0.4rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .section-header i {
    color: var(--color-brand-500);
    font-size: 1rem;
  }

  .section-header h4 {
    color: var(--color-text-content);
    font-size: 1rem;
    font-weight: 600;
    margin: 0;
  }

  .section-line {
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.3), transparent);
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1.5rem;
    text-align: center;
    color: var(--color-slate-500);
    height: 100%;
    gap: 1.5rem;
  }

  .empty-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .empty-icon-wrapper {
    font-size: 2rem;
    color: var(--color-brand-500);
    opacity: 0.7;
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 0.7;
    }
    50% {
      opacity: 1;
    }
  }

  .empty-text-content {
    max-width: 350px;
  }

  .empty-title {
    font-size: 1.1rem;
    color: var(--color-text-content);
    margin: 0 0 0.5rem 0;
    font-weight: 600;
  }

  .empty-subtitle {
    font-size: 0.85rem;
    color: var(--color-slate-500);
    margin: 0;
    line-height: 1.4;
  }

  .test-guide {
    width: 100%;
    max-width: 500px;
    margin: 1.5rem auto;
    text-align: center;
  }

  .guide-steps {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .guide-step {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
  }

  .step-num {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.5rem;
    height: 1.5rem;
    background: var(--color-brand-500);
    color: var(--color-text-on-brand);
    border-radius: var(--radius-full);
    font-weight: 700;
    font-size: 0.75rem;
  }

  .step-text {
    color: var(--color-content-default);
    font-size: 0.875rem;
    font-weight: 500;
  }

  .step-arrow {
    color: var(--color-brand-500);
    font-size: 1.25rem;
    font-weight: bold;
  }

  .guide-note {
    color: var(--color-content-muted);
    font-size: 0.8rem;
    margin: 0;
    line-height: 1.5;
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem 2rem;
    text-align: center;
    height: 100%;
  }

  .loading-animation {
    margin-bottom: 1.5rem;
  }

  .loading-rings {
    display: flex;
    gap: 0.5rem;
  }

  .ring {
    width: 12px;
    height: 12px;
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-top: 2px solid var(--color-brand-500);
    border-radius: var(--radius-full);
    animation: spin 1s linear infinite;
  }

  .ring-2 {
    animation-delay: 0.1s;
  }

  .ring-3 {
    animation-delay: 0.2s;
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
    color: var(--color-gray-200);
    font-size: 1rem;
    margin: 0;
  }

  .response-data {
    height: 100%;
  }

  .request-summary {
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 1rem;
    margin-bottom: 1.5rem;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
  }

  .summary-item {
    text-align: center;
  }

  .label {
    display: block;
    color: var(--color-slate-500);
    font-size: 0.75rem;
    margin-bottom: 0.25rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .value {
    color: var(--color-brand-500);
    font-size: 1rem;
    font-weight: 600;
  }

  .json-response {
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-md);
    overflow: hidden;
    margin-bottom: 1.5rem;
  }

  .json-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .json-header span {
    color: var(--color-brand-500);
    font-size: 0.875rem;
    font-weight: 600;
  }

  .copy-json-btn {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 0.75rem;
    display: flex;
    align-items: center;
    gap: 0.25rem;
    transition: all 0.3s ease;
  }

  .copy-json-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .copy-json-btn.copied {
    background: rgba(var(--color-success-rgb), 0.15);
    border-color: rgba(var(--color-success-rgb), 0.4);
    color: var(--color-success-500);
  }

  .json-container {
    max-height: 300px;
    overflow-y: auto;
  }

  .json-code {
    background:
      linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.015) 1px, transparent 1px),
      linear-gradient(rgba(var(--color-brand-500-rgb), 0.015) 1px, transparent 1px), var(--color-background-900);
    background-size:
      14px 14px,
      14px 14px;
    color: var(--color-text-content);
    padding: 1rem;
    margin: 0;
    font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
    font-size: 0.8rem;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-word;
    position: relative;
  }

  @keyframes json-grid-flow {
    0% {
      background-position:
        0 0,
        0 0;
    }
    100% {
      background-position:
        14px 14px,
        14px 14px;
    }
  }

  .upload-results {
    margin-bottom: 1.5rem;
  }

  .results-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--color-success-500);
    font-size: 0.9rem;
    font-weight: 600;
    margin: 0 0 1rem 0;
    position: sticky;
    top: 0;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-700-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.95) 100%
    );
    padding: 0.5rem 0;
    z-index: 1;
    border-bottom: 1px solid var(--color-border-subtle);
  }

  .results-title.error {
    color: var(--color-error-500);
  }

  .images-grid {
    display: grid;
    gap: 1rem;
    padding-right: 0.5rem;
  }

  .image-result-card {
    display: flex;
    gap: 0.75rem;
    padding: 0.75rem;
    background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.08) 0%, rgba(var(--color-background-700-rgb), 0.9) 100%);
    border: 1px solid rgba(var(--color-success-rgb), 0.2);
    border-radius: var(--radius-sm);
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
    min-height: 120px;
  }

  @keyframes grid-subtle {
    0% {
      background-position:
        0 0,
        0 0;
    }
    100% {
      background-position:
        20px 20px,
        20px 20px;
    }
  }

  .image-result-card:hover {
    background: linear-gradient(
      135deg,
      rgba(var(--color-success-rgb), 0.12) 0%,
      rgba(var(--color-background-700-rgb), 0.95) 100%
    );
    border-color: rgba(var(--color-success-rgb), 0.4);
    transform: translateY(-1px);
    box-shadow:
      0 4px 15px var(--color-border-subtle),
      0 0 12px rgba(var(--color-success-rgb), 0.06);
  }

  .image-preview {
    width: 70px;
    height: 70px;
    border-radius: var(--radius-sm);
    overflow: hidden;
    background: rgba(var(--color-background-900-rgb), 0.3);
    position: relative;
    flex-shrink: 0;
    border: 1px solid rgba(var(--color-success-rgb), 0.2);
    transition: all 0.3s ease;
  }

  .image-preview:hover {
    border-color: rgba(var(--color-success-rgb), 0.4);
    transform: scale(1.05);
  }

  .image-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  .image-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.8) 0%,
      rgba(var(--color-background-700-rgb), 0.9) 100%
    );
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.25rem;
    opacity: 0;
    transition: all 0.3s ease;
    backdrop-filter: blur(3px);
  }

  .image-preview:hover .image-overlay {
    opacity: 1;
  }

  .copy-url-btn,
  .view-original-btn {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-brand-500);
    width: 24px;
    height: 24px;
    border-radius: var(--radius-sm);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.7rem;
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
  }

  .view-original-btn {
    background: rgba(var(--color-success-rgb), 0.2);
    border-color: rgba(var(--color-success-rgb), 0.4);
    color: var(--color-success-500);
  }

  .copy-url-btn:hover,
  .view-original-btn:hover {
    transform: scale(1.1);
    box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .view-original-btn:hover {
    box-shadow: 0 2px 8px rgba(var(--color-success-rgb), 0.3);
  }

  .image-details {
    flex: 1;
    position: relative;
    z-index: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }

  .image-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 0.5rem;
    padding-bottom: 0.25rem;
    border-bottom: 1px solid var(--color-border-subtle);
  }

  .image-name {
    color: var(--color-text-content);
    font-size: 0.85rem;
    font-weight: 600;
    margin: 0;
    word-break: break-all;
    line-height: 1.2;
    flex: 1;
    background: linear-gradient(135deg, var(--color-text-content) 0%, var(--color-success-500) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    max-height: 2.4rem;
    overflow: hidden;
  }

  .image-id {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    padding: 0.15rem 0.3rem;
    border-radius: var(--radius-sm);
    font-size: 0.6rem;
    font-family: 'SF Mono', Monaco, monospace;
    font-weight: 600;
    margin-left: 0.5rem;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    white-space: nowrap;
    align-self: flex-start;
  }

  .image-meta-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 0.4rem;
    flex: 1;
    align-items: center;
  }

  .meta-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: 1px solid rgba(var(--color-success-rgb), 0.15);
    border-radius: var(--radius-sm);
    padding: 0.25rem 0.3rem;
    transition: all 0.3s ease;
    text-align: center;
    min-height: 38px;
  }

  .meta-item:hover {
    background: rgba(var(--color-success-rgb), 0.03);
    border-color: rgba(var(--color-success-rgb), 0.25);
    transform: translateY(-1px);
  }

  .meta-label {
    color: var(--color-slate-500);
    font-size: 0.6rem;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.3px;
    margin-bottom: 0.15rem;
  }

  .meta-value {
    color: var(--color-gray-200);
    font-size: 0.65rem;
    font-weight: 600;
    text-align: center;
    line-height: 1.2;
  }

  .format-badge {
    background: transparent;
    color: var(--color-error-500);
    padding: 0.1rem 0.25rem;
    border-radius: var(--radius-sm);
    font-size: 0.6rem !important;
    font-weight: 600 !important;
    border: 1px solid rgba(var(--color-badge-accent-text-rgb), 0.3);
  }

  .access-level-badge {
    background: transparent;
    padding: 0.1rem 0.25rem;
    border-radius: var(--radius-sm);
    font-size: 0.6rem !important;
    font-weight: 600 !important;
    border: 1px solid;
  }

  .access-level-badge.public {
    color: var(--color-success-500);
    border-color: rgba(var(--color-success-rgb), 0.3);
  }

  .access-level-badge.private {
    color: var(--color-error-500);
    border-color: rgba(var(--color-badge-accent-text-rgb), 0.3);
  }

  .access-level-badge.protected {
    color: var(--color-warning-400);
    border-color: rgba(255, 159, 64, 0.3);
  }

  .error-results {
    margin-bottom: 1.5rem;
  }

  .error-section {
    background: rgba(var(--color-badge-accent-text-rgb), 0.05);
    border: 1px solid rgba(var(--color-badge-accent-text-rgb), 0.15);
    border-radius: var(--radius-sm);
    padding: 1rem;
    margin-bottom: 1rem;
  }

  .error-section h6 {
    color: var(--color-error-500);
    font-size: 0.875rem;
    font-weight: 600;
    margin: 0 0 0.5rem 0;
  }

  .error-list {
    margin: 0;
    padding-left: 1.25rem;
    color: var(--color-gray-200);
    font-size: 0.8rem;
  }

  .size-limit {
    color: var(--color-slate-500);
    font-size: 0.75rem;
    margin: 0.5rem 0 0 0;
    font-style: italic;
  }

  .status-indicator {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
  }

  .status-indicator .status-dot {
    width: 8px;
    height: 8px;
    border-radius: var(--radius-full);
  }

  .status-indicator.success {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-500);
  }

  .status-indicator.success .status-dot {
    background: var(--color-success-500);
  }

  .status-indicator.error {
    background: rgba(var(--color-badge-accent-text-rgb), 0.15);
    color: var(--color-error-500);
  }

  .status-indicator.error .status-dot {
    background: var(--color-error-500);
  }

  .status-indicator.server-error {
    background: rgba(248, 113, 113, 0.15);
    color: var(--color-error-500);
  }

  .status-indicator.server-error .status-dot {
    background: var(--color-error-500);
  }

  .json-code :deep(.hljs-attr) {
    color: var(--color-error-500) !important;
  }

  .json-code :deep(.hljs-string) {
    color: var(--color-success-500) !important;
  }

  .json-code :deep(.hljs-number) {
    color: var(--color-purple-500) !important;
  }

  .json-code :deep(.hljs-literal) {
    color: var(--color-brand-500) !important;
  }

  .json-code :deep(.hljs-punctuation) {
    color: var(--color-content-heading) !important;
  }

  .json-code :deep(.hljs-keyword) {
    color: var(--color-brand-500) !important;
    font-weight: 600;
  }

  @media (max-width: 768px) {
    .quick-tips {
      gap: 1rem;
    }

    .tip-item {
      min-width: 70px;
      padding: 0.6rem;
    }

    .request-summary {
      grid-template-columns: 1fr;
      gap: 0.75rem;
    }
  }

  .json-container::-webkit-scrollbar,
  .upload-results::-webkit-scrollbar,
  .image-meta-grid::-webkit-scrollbar {
    width: 6px;
  }

  .json-container::-webkit-scrollbar-track,
  .upload-results::-webkit-scrollbar-track,
  .image-meta-grid::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.05);
    border-radius: var(--radius-sm);
  }

  .json-container::-webkit-scrollbar-thumb,
  .upload-results::-webkit-scrollbar-thumb,
  .image-meta-grid::-webkit-scrollbar-thumb {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.6) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.6) 100%
    );
    border-radius: var(--radius-sm);
  }

  .json-container::-webkit-scrollbar-thumb:hover,
  .upload-results::-webkit-scrollbar-thumb:hover,
  .image-meta-grid::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.8) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.8) 100%
    );
  }
</style>
