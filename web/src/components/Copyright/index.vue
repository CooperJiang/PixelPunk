<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { useSettingsStore } from '@/store/settings'
  import type { CopyrightProps } from './types'

  defineOptions({
    name: 'Copyright',
  })

  const { $t } = useTexts()
  const settingsStore = useSettingsStore()

  const props = withDefaults(defineProps<CopyrightProps>(), {
    theme: 'dark',
  })

  const currentYear = computed(() => new Date().getFullYear())

  const handleEmailClick = (event: MouseEvent) => {
    const target = event.target as HTMLAnchorElement
    const email = settingsStore.contactEmail

    if (!email) {
      event.preventDefault()
      return
    }

    try {
      if (!target.href || target.href === '') {
        event.preventDefault()
        window.location.href = `mailto:${email}`
      }
    } catch (_error) {}
  }
</script>

<template>
  <div class="copyright-section" :class="{ dark: props.theme === 'dark', transparent: props.theme === 'transparent' }">
    <div class="copyright-container">
      <div class="copyright-content">
        <div class="desktop-layout">
          <div class="copyright-info">
            <div class="copyright-text">
              <span class="copyright-basic">{{
                settingsStore.copyrightText ||
                `© ${currentYear} ${settingsStore.siteName || 'PixelPunk'}. ${$t('components.copyright.allRightsReserved')}`
              }}</span>
              <template v-if="settingsStore.icpNumber">
                <span class="info-divider">|</span>
                <span class="icp-info">
                  <i class="fas fa-shield-alt" />
                  <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener" class="icp-link">{{
                    settingsStore.icpNumber
                  }}</a>
                </span>
              </template>
              <template v-if="settingsStore.footerCustomText">
                <span class="info-divider">|</span>
                <span class="custom-text">{{ settingsStore.footerCustomText }}</span>
              </template>
              <template v-if="settingsStore.contactEmail">
                <span class="info-divider">|</span>
                <span class="contact-info">
                  <i class="fas fa-envelope" />
                  <a
                    :href="`mailto:${settingsStore.contactEmail}`"
                    class="contact-link"
                    @click="handleEmailClick"
                    rel="noopener noreferrer"
                    >{{ settingsStore.contactEmail }}</a
                  >
                </span>
              </template>
            </div>
          </div>

          <div class="status-indicators">
            <div class="status-item online">
              <div class="status-dot" />
              <span>{{ $t('components.copyright.systemOnline') }}</span>
            </div>
            <div class="status-item secure">
              <i class="fas fa-shield-alt" />
              <span>{{ $t('components.copyright.secureRunning') }}</span>
            </div>
          </div>
        </div>

        <div class="mobile-layout">
          <div class="mobile-content">
            <div class="copyright-text-mobile">
              <p>
                {{
                  settingsStore.copyrightText ||
                  `© ${currentYear} ${settingsStore.siteName || 'PixelPunk'}. ${$t('components.copyright.allRightsReserved')}`
                }}
                <span v-if="settingsStore.icpNumber" class="info-divider-mobile">|</span>
                <span v-if="settingsStore.icpNumber" class="icp-info-mobile">
                  <i class="fas fa-shield-alt" />
                  <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener" class="icp-link-mobile">{{
                    settingsStore.icpNumber
                  }}</a>
                </span>
                <span v-if="settingsStore.footerCustomText" class="info-divider-mobile">|</span>
                <span v-if="settingsStore.footerCustomText" class="custom-text-mobile">{{ settingsStore.footerCustomText }}</span>
                <span v-if="settingsStore.contactEmail" class="info-divider-mobile">|</span>
                <span v-if="settingsStore.contactEmail" class="contact-info-mobile">
                  <i class="fas fa-envelope" />
                  <a
                    :href="`mailto:${settingsStore.contactEmail}`"
                    class="contact-link-mobile"
                    @click="handleEmailClick"
                    rel="noopener noreferrer"
                    >{{ settingsStore.contactEmail }}</a
                  >
                </span>
              </p>
            </div>

            <div class="status-indicators-mobile">
              <div class="status-item online">
                <div class="status-dot" />
                <span>{{ $t('components.copyright.online') }}</span>
              </div>
              <div class="status-item secure">
                <i class="fas fa-shield-alt" />
                <span>{{ $t('components.copyright.secure') }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="copyright-decoration" />
  </div>
</template>

<style scoped>
  .copyright-section {
    @apply relative py-4;
    background: rgba(var(--color-background-900-rgb), 0.5);
    border-top: 1px solid var(--color-border-subtle);
  }

  .copyright-section.transparent {
    @apply backdrop-blur-[10px];
    background: rgba(var(--color-background-900-rgb), 0.3);
  }

  .copyright-container {
    @apply mx-auto max-w-[1400px] px-8;
  }

  .copyright-content {
    @apply relative;
    z-index: 5;
  }

  .desktop-layout {
    @apply flex flex-wrap items-center justify-between gap-4;
  }

  .mobile-layout {
    @apply hidden;
  }

  .copyright-info {
    @apply min-w-0 flex-1;
  }

  .copyright-text {
    @apply m-0 text-sm text-content;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .info-divider {
    @apply text-content-muted;
    opacity: 0.6;
    flex-shrink: 0;
  }

  .copyright-basic,
  .icp-info,
  .custom-text,
  .contact-info {
    white-space: nowrap;
    flex-shrink: 0;
  }

  .icp-info {
    @apply inline-flex items-center gap-1;
  }

  .icp-info i {
    @apply text-xs;
    color: var(--color-brand-500);
    opacity: 0.8;
  }

  .icp-link,
  .contact-link {
    color: var(--color-brand-500);
    text-decoration: none;
    transition: all 0.3s ease;
    position: relative;
    z-index: 10;
    cursor: pointer;
    padding: 2px 4px;
    border-radius: var(--radius-sm);
  }

  .icp-link:hover,
  .contact-link:hover {
    color: var(--color-error-500);
    text-shadow: 0 0 8px rgba(var(--color-error-rgb), 0.6);
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .icp-link:active,
  .contact-link:active {
    background: rgba(var(--color-error-rgb), 0.2);
    transform: scale(0.98);
  }

  .custom-text {
    color: var(--color-success-500);
    opacity: 0.9;
  }

  .contact-info {
    @apply inline-flex items-center gap-1;
  }

  .contact-info i {
    @apply text-xs;
    color: var(--color-success-500);
    opacity: 0.8;
  }

  .status-indicators {
    @apply flex flex-shrink-0 gap-6;
  }

  .status-item {
    @apply flex items-center gap-2 text-xs;
  }

  .status-item.online {
    color: var(--color-success-500);
  }

  .status-item.secure {
    color: var(--color-brand-500);
  }

  .status-dot {
    @apply h-1.5 w-1.5 rounded-full;
    background: var(--color-success-500);
    box-shadow: 0 0 8px rgba(var(--color-success-rgb), 0.5);
    animation: pulse 2s infinite;
  }

  .status-item i {
    @apply text-xs opacity-80;
  }

  .icp-info {
    @apply mt-1 text-xs text-content-muted;
  }

  .icp-info a {
    color: var(--color-brand-500);
    text-decoration: none;
  }

  .icp-info a:hover {
    text-decoration: underline;
  }

  .custom-text {
    @apply mt-1 text-xs text-content-muted;
  }

  .contact-info {
    @apply mt-1 flex items-center gap-1 text-xs text-content-muted;
  }

  .contact-info a {
    color: var(--color-brand-500);
    text-decoration: none;
  }

  .contact-info a:hover {
    text-decoration: underline;
  }

  .contact-info i {
    @apply text-xs opacity-80;
  }

  .icp-info-mobile {
    @apply mt-1 text-xs text-content-muted;
  }

  .icp-info-mobile a {
    color: var(--color-brand-500);
    text-decoration: none;
  }

  .icp-info-mobile a:hover {
    text-decoration: underline;
  }

  .custom-text-mobile {
    @apply mt-1 text-xs text-content-muted;
  }

  .contact-info-mobile {
    @apply mt-1 flex items-center gap-1 text-xs text-content-muted;
  }

  .contact-info-mobile a {
    color: var(--color-brand-500);
    text-decoration: none;
  }

  .contact-info-mobile a:hover {
    text-decoration: underline;
  }

  .contact-info-mobile i {
    @apply text-xs opacity-80;
  }

  .copyright-decoration {
    @apply absolute bottom-0 left-1/2 h-px w-3/5 -translate-x-1/2;
    background: linear-gradient(90deg, transparent, var(--color-brand-500), transparent);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.6);
    pointer-events: none;
    z-index: 1;
  }

  @media (max-width: 768px) {
    .desktop-layout {
      @apply hidden;
    }

    .mobile-layout {
      @apply block;
    }

    .copyright-container {
      @apply px-6;
    }

    .mobile-content {
      @apply flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between;
    }

    .copyright-text-mobile {
      @apply text-center sm:text-left;
    }

    .copyright-text-mobile p {
      @apply m-0 text-sm leading-relaxed text-content;
    }

    .info-divider-mobile {
      @apply mx-1 text-content-muted;
      opacity: 0.6;
    }

    .icp-info-mobile {
      @apply inline-flex items-center gap-1;
    }

    .icp-info-mobile i {
      @apply text-xs;
      color: var(--color-brand-500);
      opacity: 0.8;
    }

    .icp-link-mobile,
    .contact-link-mobile {
      color: var(--color-brand-500);
      text-decoration: none;
      transition: all 0.3s ease;
      position: relative;
      z-index: 10;
      cursor: pointer;
      padding: 2px 4px;
      border-radius: var(--radius-sm);
      -webkit-tap-highlight-color: transparent;
    }

    .icp-link-mobile:hover,
    .contact-link-mobile:hover {
      color: var(--color-error-500);
      text-shadow: 0 0 8px rgba(var(--color-error-rgb), 0.6);
      background: rgba(var(--color-brand-500-rgb), 0.1);
    }

    .icp-link-mobile:active,
    .contact-link-mobile:active {
      background: rgba(var(--color-error-rgb), 0.2);
      transform: scale(0.98);
    }

    .custom-text-mobile {
      color: var(--color-success-500);
      opacity: 0.9;
    }

    .contact-info-mobile {
      @apply inline-flex items-center gap-1;
    }

    .contact-info-mobile i {
      @apply text-xs;
      color: var(--color-success-500);
      opacity: 0.8;
    }

    .status-indicators-mobile {
      @apply flex flex-shrink-0 justify-center gap-3 sm:justify-end;
    }

    .status-indicators-mobile .status-item {
      @apply gap-1 text-xs;
    }

    .status-indicators-mobile .status-item span {
      @apply text-[0.65rem];
    }

    .status-indicators-mobile .status-dot {
      @apply h-[5px] w-[5px];
    }

    .status-indicators-mobile .status-item i {
      @apply text-[0.65rem];
    }
  }

  @media (max-width: 480px) {
    .copyright-container {
      @apply px-5;
    }

    .mobile-content {
      @apply gap-2;
    }

    .status-indicators-mobile {
      @apply gap-2;
    }

    .copyright-text-mobile p {
      @apply text-xs;
    }
  }

  @media (max-width: 360px) {
    .copyright-container {
      @apply px-4;
    }

    .mobile-content {
      @apply flex-col items-center gap-2;
    }

    .status-indicators-mobile .status-item span {
      @apply hidden;
    }

    .status-indicators-mobile .status-item {
      @apply min-w-[20px] justify-center gap-0;
    }
  }
</style>
