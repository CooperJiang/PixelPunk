<template>
  <div class="api-tester-cyber">
    <TesterHeader
      :currentDomain="currentDomain"
      :canSend="selectedFiles.length > 0 && apiKey.trim().length > 0"
      :isLoading="isLoading"
      @send-request="sendRequest"
    />

    <div class="cyber-tester-body">
      <ParamsPanel
        v-model:apiKey="apiKey"
        v-model:selectedFiles="selectedFiles"
        v-model:params="params"
        v-model:showValidation="showValidation"
        v-model:isDragOver="isDragOver"
        :totalFileSize="totalFileSize"
      />

      <ResponsePanel
        :response="response"
        :isLoading="isLoading"
        :responseCopied="responseCopied"
        :selectedFiles="selectedFiles"
        :totalFileSize="totalFileSize"
        :formattedResponse="formattedResponse"
        :responseStatusClass="responseStatusClass"
        :uploadedImages="uploadedImages"
        :hasErrors="hasErrors"
        :oversizedFiles="oversizedFiles"
        :unsupportedFiles="unsupportedFiles"
        :invalidFiles="invalidFiles"
        :uploadErrors="uploadErrors"
        :sizeLimit="sizeLimit"
        :formatFileSize="formatFileSize"
        @copy-response="copyResponse"
        @copy-image-url="copyImageUrl"
        @open-image="openImageInNewTab"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
  import TesterHeader from './ApiTester/components/TesterHeader.vue'
  import ParamsPanel from './ApiTester/components/ParamsPanel.vue'
  import ResponsePanel from './ApiTester/components/ResponsePanel.vue'
  import { useApiTest } from './ApiTester/composables/useApiTest'
  import { useFileUpload } from './ApiTester/composables/useFileUpload'
  import { useJsonFormatter } from './ApiTester/composables/useJsonFormatter'

  defineOptions({
    name: 'ApiTester',
  })

  /* 使用组合式函数 */
  const apiTestComposable = useApiTest()
  const {
    apiKey,
    selectedFiles,
    params,
    response,
    isLoading,
    responseCopied,
    showValidation,
    isDragOver,
    currentDomain,
    totalFileSize,
    responseStatusClass,
    uploadedImages,
    hasErrors,
    oversizedFiles,
    unsupportedFiles,
    invalidFiles,
    uploadErrors,
    sizeLimit,
    sendRequest,
    copyResponse,
    copyImageUrl,
    openImageInNewTab,
  } = apiTestComposable

  /* 使用文件上传功能 */
  const { formatFileSize } = useFileUpload(selectedFiles, showValidation, isDragOver)

  /* 使用JSON格式化功能 */
  const { formattedResponse } = useJsonFormatter(response)
</script>

<style scoped>
  .api-tester-cyber {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.8) 0%,
      rgba(var(--color-background-800-rgb), 0.9) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    overflow: hidden;
    margin-bottom: 2rem;
    box-shadow:
      0 8px 32px var(--color-background-700),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-tester-body {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0;
    min-height: 400px;
  }

  @media (max-width: 768px) {
    .cyber-tester-body {
      grid-template-columns: 1fr;
    }
  }
</style>
