import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'

import { CODE_TAB_OPTIONS, getDocsNavItems, getImageFormatOptions } from '@/constants'
import { logger } from '@/utils/system/logger'
import { useTexts } from '@/composables/useTexts'
import {
  BACK_TO_TOP_THRESHOLD,
  SCROLL_OFFSET,
  SECTION_ACTIVE_OFFSET,
  COPY_SUCCESS_DURATION,
  DOM_INIT_DELAY,
  SCROLL_REINIT_DELAY,
} from '../constants'

const { $t } = useTexts()
export const navItems = computed(() => getDocsNavItems($t))
export const codeTabs = CODE_TAB_OPTIONS
export const supportedFormats = getImageFormatOptions($t).map((option) => option.value.toUpperCase())

export function useDocsPage() {
  const activeSection = ref('api-overview')
  const activeCodeTab = ref('curl')
  const copyStatus = ref({})
  const mobileMenuOpen = ref(false)
  const showBackToTop = ref(false)
  const readingProgress = ref(0)

  const currentDomain = computed(() => {
    if (typeof window !== 'undefined') {
      return window.location.origin
    }
    return 'https://your-domain.com'
  })

  const toggleMobileMenu = () => {
    mobileMenuOpen.value = !mobileMenuOpen.value
  }

  const handleMobileNavClick = (event) => {
    scrollToSection(event)
    mobileMenuOpen.value = false
  }

  const scrollToSection = (eventOrTarget) => {
    if (eventOrTarget && typeof eventOrTarget === 'object' && eventOrTarget.preventDefault) {
      eventOrTarget.preventDefault()
    }

    let targetId = ''

    if (typeof eventOrTarget === 'string') {
      targetId = eventOrTarget.replace('#', '')
    } else if (eventOrTarget && eventOrTarget.target) {
      targetId =
        eventOrTarget.target.getAttribute('href')?.substring(1) || eventOrTarget.currentTarget.getAttribute('href')?.substring(1)
    }

    if (targetId) {
      const element = document.getElementById(targetId)
      if (element) {
        const docsContent = document.querySelector('.docs-content')
        if (docsContent) {
          const elementTop = element.offsetTop
          docsContent.scrollTo({
            top: elementTop - SCROLL_OFFSET,
            behavior: 'smooth',
          })
        }

        activeSection.value = targetId
      }
    }
  }

  const backToTop = () => {
    const docsContent = document.querySelector('.docs-content')
    if (docsContent) {
      docsContent.scrollTo({
        top: 0,
        behavior: 'smooth',
      })
    }
  }

  const handleScroll = () => {
    const docsContent = document.querySelector('.docs-content')
    if (!docsContent) {
      return
    }

    const { scrollTop, scrollHeight, clientHeight } = docsContent
    const docHeight = scrollHeight - clientHeight

    if (docHeight > 0) {
      const progress = Math.min((scrollTop / docHeight) * 100, 100)
      readingProgress.value = progress
    } else {
      readingProgress.value = 0
    }

    showBackToTop.value = scrollTop > BACK_TO_TOP_THRESHOLD
    updateActiveSection()
  }

  const updateActiveSection = () => {
    const docsContent = document.querySelector('.docs-content')
    if (!docsContent) {
      return
    }

    const sections = document.querySelectorAll('[id^="api-"]')
    const { scrollTop } = docsContent

    let current = ''

    sections.forEach((section) => {
      const offsetTop = section.offsetTop - SECTION_ACTIVE_OFFSET
      if (scrollTop >= offsetTop) {
        current = section.id
      }
    })

    if (current && current !== activeSection.value) {
      activeSection.value = current
    }
  }

  const copyCode = (codeId, code) => {
    if (code) {
      navigator.clipboard
        .writeText(code)
        .then(() => {
          copyStatus.value[codeId] = true
          setTimeout(() => {
            copyStatus.value[codeId] = false
          }, COPY_SUCCESS_DURATION)
        })
        .catch((err) => {
          logger.error('Copy failed:', err)
        })
    }
  }

  const switchCodeTab = (tabId) => {
    activeCodeTab.value = tabId
    nextTick(() => {
      highlightCodeBlocks()
    })
  }

  const highlightCodeBlocks = () => {
    document.querySelectorAll('pre code:not([data-highlighted])').forEach((block) => {
      hljs.highlightElement(block)
    })
  }

  const handleClickOutside = (event) => {
    const mobileMenu = document.querySelector('.mobile-menu')
    const menuBtn = document.querySelector('.mobile-menu-btn')

    if (mobileMenuOpen.value && mobileMenu && !mobileMenu.contains(event.target) && !menuBtn.contains(event.target)) {
      mobileMenuOpen.value = false
    }
  }

  onMounted(() => {
    setTimeout(() => {
      const docsContent = document.querySelector('.docs-content')
      if (docsContent) {
        docsContent.addEventListener('scroll', handleScroll, { passive: true })
        handleScroll()
      }
    }, DOM_INIT_DELAY)

    document.addEventListener('click', handleClickOutside)

    nextTick(() => {
      highlightCodeBlocks()
      setTimeout(() => {
        handleScroll()
      }, SCROLL_REINIT_DELAY)
    })
  })

  onUnmounted(() => {
    const docsContent = document.querySelector('.docs-content')
    if (docsContent) {
      docsContent.removeEventListener('scroll', handleScroll)
    }
    document.removeEventListener('click', handleClickOutside)
  })

  return {
    activeSection,
    activeCodeTab,
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
    highlightCodeBlocks,
  }
}
