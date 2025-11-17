<script lang="ts">
  import { defineComponent, ref } from 'vue'
  import ApiSection from './ApiSection/index.vue'
  import { useTexts } from '@/composables/useTexts'

  export default defineComponent({
    name: 'FaqSection',
    components: {
      ApiSection,
    },
    setup() {
      const { $t } = useTexts()
      const openFaqIndex = ref<number | null>(null)

      const faqList = $t('docs.faq.items')

      const toggleFaq = (index: number) => {
        openFaqIndex.value = openFaqIndex.value === index ? null : index
      }

      return {
        openFaqIndex,
        faqList,
        toggleFaq,
      }
    },
  })
</script>

<template>
  <ApiSection id="api-faq" :title="$t('docs.faq.title')" icon-class="fas fa-question-circle">
    <div class="faq-accordion">
      <div v-for="(faq, index) in faqList" :key="index" class="faq-accordion-item" :class="{ active: openFaqIndex === index }">
        <div class="faq-header" @click="toggleFaq(index)">
          <div class="faq-title-wrapper">
            <i :class="faq.icon" class="faq-icon" />
            <h4 class="faq-title">{{ faq.question }}</h4>
          </div>
          <div class="faq-toggle">
            <i class="fas fa-chevron-down" :class="{ rotated: openFaqIndex === index }" />
          </div>
        </div>

        <Transition name="faq-content">
          <div v-show="openFaqIndex === index" class="faq-content">
            <div class="faq-answer">
              <p>{{ faq.answer }}</p>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </ApiSection>
</template>

<style scoped>
  .faq-accordion {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .faq-accordion-item {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    transform: translateZ(0);
    backface-visibility: hidden;
  }

  .faq-accordion-item:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-2px) translateZ(0);
    box-shadow: 0 4px 20px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .faq-accordion-item.active {
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    background: rgba(var(--color-background-800-rgb), 0.8);
    box-shadow: 0 6px 25px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .faq-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
    padding: 0.25rem 0;
    user-select: none;
    transition: all 0.2s ease;
  }

  .faq-header:hover .faq-title {
    color: var(--color-brand-500);
  }

  .faq-header:hover .faq-icon {
    transform: scale(1.1);
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.6);
  }

  .faq-title-wrapper {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex: 1;
  }

  .faq-icon {
    color: var(--color-brand-500);
    font-size: 1.1rem;
    transition: all 0.3s ease;
    opacity: 0.8;
  }

  .faq-accordion-item.active .faq-icon {
    opacity: 1;
    transform: scale(1.05);
    filter: drop-shadow(0 0 8px rgba(var(--color-brand-500-rgb), 0.4));
  }

  .faq-title {
    margin: 0;
    font-size: 0.95rem;
    color: var(--color-text-content);
    line-height: 1.4;
    font-weight: 600;
    transition: color 0.3s ease;
  }

  .faq-accordion-item.active .faq-title {
    color: var(--color-brand-500);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .faq-toggle {
    color: var(--color-content-muted);
    font-size: 1rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: var(--radius-full);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    flex-shrink: 0;
  }

  .faq-toggle:hover {
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.2);
    transform: scale(1.1);
  }

  .faq-toggle .fas.rotated {
    transform: rotate(180deg);
    color: var(--color-brand-500);
  }

  .faq-accordion-item.active .faq-toggle {
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .faq-content {
    overflow: hidden;
    will-change: transform, opacity;
  }

  .faq-answer {
    padding-top: 0.75rem;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    margin-top: 0.75rem;
    transform: translateZ(0);
  }

  .faq-answer p {
    color: var(--color-gray-200);
    font-size: 0.875rem;
    line-height: 1.6;
    margin: 0;
    opacity: 0.9;
    text-align: justify;
  }

  .faq-content-enter-active {
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    overflow: hidden;
  }

  .faq-content-leave-active {
    transition: all 0.25s cubic-bezier(0.25, 0.8, 0.25, 1);
    overflow: hidden;
  }

  .faq-content-enter-from {
    opacity: 0;
    transform: translateY(-8px) scaleY(0.95);
  }

  .faq-content-leave-to {
    opacity: 0;
    transform: translateY(-8px) scaleY(0.95);
  }

  .faq-content-enter-to,
  .faq-content-leave-from {
    opacity: 1;
    transform: translateY(0) scaleY(1);
  }

  @media (max-width: 768px) {
    .faq-accordion {
      flex-direction: column;
    }

    .faq-accordion-item {
      padding: 1rem;
    }

    .faq-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .faq-toggle {
      margin-left: 0;
    }

    .faq-content {
      margin-top: 0.5rem;
    }
  }
</style>
