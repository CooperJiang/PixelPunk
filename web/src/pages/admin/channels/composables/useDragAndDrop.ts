import { onMounted, onUnmounted, ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'

/**
 * 拖拽上传 Composable
 * 负责拖拽文件上传功能
 */
export function useDragAndDrop(
  containerSelector: string,
  onFileDrop: (file: File) => void,
  options?: { isDisabled?: () => boolean }
) {
  const toast = useToast()
  const { $t } = useTexts()

  const dragOverCounter = ref(0)
  const isDragOver = ref(false)

  const handleDragEnter = (event: DragEvent) => {
    event.preventDefault()
    dragOverCounter.value++
    isDragOver.value = true
  }

  const handleDragOver = (event: DragEvent) => {
    event.preventDefault()
  }

  const handleDragLeave = (event: DragEvent) => {
    event.preventDefault()
    dragOverCounter.value--
    if (dragOverCounter.value === 0) {
      isDragOver.value = false
    }
  }

  const handleDrop = async (event: DragEvent) => {
    event.preventDefault()
    dragOverCounter.value = 0
    isDragOver.value = false

    const files = event.dataTransfer?.files
    if (!files || files.length === 0) {
      return
    }

    const file = files[0]

    if (!file.name.toLowerCase().endsWith('.json') && file.type !== 'application/json') {
      toast.error($t('admin.channels.messages.dragJsonOnly'))
      return
    }

    onFileDrop(file)
  }

  const isDialogOpen = () => !!options?.isDisabled && options.isDisabled()

  const isEditingField = (target: EventTarget | null): boolean => {
    const el = target as HTMLElement | null
    if (!el) return false
    const tag = (el.tagName || '').toLowerCase()
    if (tag === 'input' || tag === 'textarea') return true
    let cur: HTMLElement | null = el
    while (cur) {
      if (cur.getAttribute && cur.getAttribute('contenteditable') === 'true') return true
      cur = cur.parentElement
    }
    return false
  }

  const handlePaste = async (event: ClipboardEvent) => {
    if (isDialogOpen() || isEditingField(event.target)) {
      return
    }
    const items = event.clipboardData?.items
    if (!items) {
      return
    }

    for (let i = 0; i < items.length; i++) {
      const item = items[i]

      if (item.kind === 'file' && item.type === 'application/json') {
        event.preventDefault()
        const file = item.getAsFile()
        if (file) {
          onFileDrop(file)
        }
        return
      }

      if (item.kind === 'string' && item.type === 'text/plain') {
        event.preventDefault()
        item.getAsString(async (text) => {
          try {
            JSON.parse(text)

            const blob = new Blob([text], { type: 'application/json' })
            const file = new File([blob], 'pasted-config.json', { type: 'application/json' })

            onFileDrop(file)
          } catch {
            toast.error($t('admin.channels.messages.invalidJsonPaste'))
          }
        })
        return
      }
    }
  }

  const initEventListeners = () => {
    document.addEventListener('paste', handlePaste)

    const container = document.querySelector(containerSelector)
    if (container) {
      container.addEventListener('dragenter', handleDragEnter)
      container.addEventListener('dragover', handleDragOver)
      container.addEventListener('dragleave', handleDragLeave)
      container.addEventListener('drop', handleDrop)
    }
  }

  const cleanupEventListeners = () => {
    document.removeEventListener('paste', handlePaste)

    const container = document.querySelector(containerSelector)
    if (container) {
      container.removeEventListener('dragenter', handleDragEnter)
      container.removeEventListener('dragover', handleDragOver)
      container.removeEventListener('dragleave', handleDragLeave)
      container.removeEventListener('drop', handleDrop)
    }
  }

  onMounted(() => {
    initEventListeners()
  })

  onUnmounted(() => {
    cleanupEventListeners()
  })

  return {
    isDragOver,

    initEventListeners,
    cleanupEventListeners,
  }
}
