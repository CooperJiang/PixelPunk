<script setup lang="ts">
  import { computed } from 'vue'
  import { useSettingsStore } from '@/store/settings'
  import { useTexts } from '@/composables/useTexts'

  const settingsStore = useSettingsStore()
  const { $t } = useTexts()

  const footerText = computed(() => {
    const template = $t('share.footer.poweredBy')
    return template.replace('{siteName}', `<a href="/" class="brand-link">${settingsStore.siteName}</a>`)
  })
</script>

<template>
  <div class="share-footer">
    <div class="copyright" v-html="footerText"></div>
  </div>
</template>

<style scoped lang="scss">
  .share-footer {
    text-align: center;
    padding: var(--space-lg) 0;
    margin-top: var(--space-xl);
    color: var(--color-content-muted);
    font-size: var(--text-sm);
  }

  .copyright {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-xs);
  }

  .brand-link {
    color: var(--color-brand-500);
    font-weight: 600;
    text-decoration: none;
    transition: all var(--transition-fast) var(--ease-in-out);
    padding: 0 var(--space-xs);

    &:hover {
      color: var(--color-content-heading);
      text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.8);
    }
  }
</style>
