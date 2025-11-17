<script lang="ts">
  import { defineComponent } from 'vue'

  export default defineComponent({
    name: 'ApiSection',
    props: {
      id: {
        type: String,
        required: true,
      },
      title: {
        type: String,
        required: true,
      },
      iconClass: {
        type: String,
        required: true,
      },
    },
  })
</script>

<template>
  <div :id="id" class="cyber-panel api-section mb-6 rounded-lg p-6">
    <h3 class="mb-4 flex items-center text-xl font-bold"><i class="mr-2" :class="[iconClass]" />{{ title }}</h3>

    <slot />
  </div>
</template>

<style scoped>
  .api-section {
    width: 100%;
    max-width: 100%;
    overflow-x: auto;
  }

  pre,
  code {
    max-width: 100%;
    white-space: pre-wrap;
    word-break: break-word;
  }

  @media (max-width: 768px) {
    .api-section {
      padding: 1rem;
    }
  }
</style>
