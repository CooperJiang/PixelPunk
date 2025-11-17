<script setup lang="ts">
  import type { FeatureItem } from './types'
  import ApiSection from './components/ApiSection/index.vue'
  import ApiTester from './components/ApiTester.vue'
  import DocsIntro from './components/DocsIntro/index.vue'
  import DocsSidebar from './components/DocsSidebar/index.vue'
  import MobileHeader from './components/MobileHeader/index.vue'
  import BackToTop from './components/BackToTop/index.vue'
  import ApiOverview from './components/ApiOverview/index.vue'
  import ApiExamples from './components/ApiExamples/index.vue'
  import AuthenticationSection from './components/AuthenticationSection.vue'
  import UploadApiSection from './components/UploadApiSection.vue'
  import ApiLimitsSection from './components/ApiLimitsSection.vue'
  import FaqSection from './components/FaqSection.vue'
  import ApiTesterFeatures from './components/ApiTesterFeatures/index.vue'
  import { getApiTesterFeatures } from './constants'
  import { navItems, supportedFormats, useDocsPage } from './composables/useDocsPage'
  import { useCodeExamples } from './composables/useCodeExamples'

  defineOptions({
    name: 'DocsPage',
  })

  const {
    activeSection,
    copyStatus,
    mobileMenuOpen,
    showBackToTop,
    readingProgress,
    currentDomain,
    toggleMobileMenu,
    handleMobileNavClick,
    scrollToSection,
    backToTop,
    copyCode,
    switchCodeTab,
  } = useDocsPage()

  const { codeExamples, jsonResponseExamples } = useCodeExamples(currentDomain)
  const apiTesterFeatures = getApiTesterFeatures()
</script>

<template>
  <div class="docs-page">
    <CyberParticleBackground theme="docs" :max-particles="60" />

    <DocsSidebar
      class="desktop-only"
      :nav-items="navItems"
      :active-section="activeSection"
      :reading-progress="readingProgress"
      :scroll-to-section="scrollToSection"
    />

    <main class="docs-content">
      <MobileHeader
        :nav-items="navItems"
        :active-section="activeSection"
        :mobile-menu-open="mobileMenuOpen"
        :toggle-mobile-menu="toggleMobileMenu"
        :handle-mobile-nav-click="handleMobileNavClick"
      />

      <DocsIntro :scroll-to-section="scrollToSection" />

      <ApiOverview :current-domain="currentDomain" />

      <AuthenticationSection />

      <UploadApiSection :supported-formats="supportedFormats" :json-response-examples="jsonResponseExamples" />

      <ApiLimitsSection :supported-formats="supportedFormats" />

      <ApiExamples
        :code-examples="codeExamples"
        :copy-status="copyStatus"
        :copy-code="copyCode"
        :switch-code-tab="switchCodeTab"
      />

      <ApiSection id="api-tester" :title="$t('docs.tester.title')" icon-class="fas fa-flask">
        <div class="section-intro">
          <p>{{ $t('docs.tester.description') }}</p>
        </div>

        <ApiTesterFeatures :features="apiTesterFeatures" />
        <ApiTester />
      </ApiSection>

      <FaqSection />
    </main>

    <BackToTop :show-back-to-top="showBackToTop" :back-to-top="backToTop" />
  </div>
</template>

<style scoped>
  .docs-page {
    @apply relative flex h-screen w-full overflow-hidden;
    color: var(--color-text-content);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Noto Sans', Helvetica, Arial, sans-serif;
    line-height: 1.6;
  }

  .docs-content {
    @apply relative z-10 h-screen overflow-y-auto bg-transparent p-8;
    margin-left: 280px;
    width: calc(100% - 280px);
  }

  .section-intro {
    @apply mb-5;
  }

  .section-intro p {
    @apply m-0 text-sm leading-6;
    color: var(--color-gray-200);
  }

  .desktop-only {
    @apply block;
  }

  .mobile-only {
    @apply hidden;
  }

  @media (max-width: 768px) {
    .desktop-only {
      @apply hidden;
    }

    .mobile-only {
      @apply block;
    }

    .docs-page {
      @apply flex-col;
    }

    .docs-content {
      @apply ml-0 w-full p-4;
      height: calc(100vh - 60px);
    }
  }

  .docs-content::-webkit-scrollbar {
    width: 6px;
  }

  .docs-content::-webkit-scrollbar-track {
    background: rgba(var(--color-background-700-rgb), 0.6);
    border-radius: var(--radius-sm);
  }

  .docs-content::-webkit-scrollbar-thumb {
    background: linear-gradient(180deg, rgba(var(--color-brand-500-rgb), 0.7) 0%, rgba(var(--color-error-rgb), 0.7) 100%);
    border-radius: var(--radius-sm);
    transition: background 0.3s ease;
  }

  .docs-content::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(180deg, rgba(var(--color-brand-500-rgb), 0.9) 0%, rgba(var(--color-error-rgb), 0.9) 100%);
  }
</style>
