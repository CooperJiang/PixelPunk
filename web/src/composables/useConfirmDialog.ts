/**
 * 确认对话框 Composable
 * 用于在代码中便捷地调用确认对话框

 */
import { ref } from 'vue'
import { useTexts } from '@/composables/useTexts'

interface ConfirmOptions {
  title?: string
  message: string | string[]
  type?: 'warning' | 'danger' | 'info' | 'success'
  confirmText?: string
  cancelText?: string
  width?: string
  requireInput?: boolean
  inputLabel?: string
  inputPlaceholder?: string
  expectedInput?: string
}

interface ConfirmDialogState {
  isVisible: boolean
  config: ConfirmOptions
  resolve: ((value: boolean | string) => void) | null
}

const dialogState = ref<ConfirmDialogState>({
  isVisible: false,
  config: {
    message: '',
    type: 'warning',
  },
  resolve: null,
})

/**
 * 显示确认对话框
 * @param options 对话框选项
 * @returns Promise<boolean | string> - true/false 表示确认/取消，string 表示输入的内容
 */
export function showConfirmDialog(options: ConfirmOptions): Promise<boolean | string> {
  return new Promise((resolve) => {
    dialogState.value = {
      isVisible: true,
      config: {
        type: 'warning',
        ...options,
      },
      resolve,
    }
  })
}

export function handleConfirm(inputValue?: string) {
  if (dialogState.value.resolve) {
    if (dialogState.value.config.requireInput && inputValue) {
      dialogState.value.resolve(inputValue)
    } else {
      dialogState.value.resolve(true)
    }
  }
  dialogState.value.isVisible = false
  dialogState.value.resolve = null
}

export function handleCancel() {
  if (dialogState.value.resolve) {
    dialogState.value.resolve(false)
  }
  dialogState.value.isVisible = false
  dialogState.value.resolve = null
}

export function useConfirmDialog() {
  return {
    dialogState,
    showConfirmDialog,
    handleConfirm,
    handleCancel,
  }
}

export function confirmWarning(message: string, title?: string): Promise<boolean | string> {
  const { $t } = useTexts()
  return showConfirmDialog({
    message,
    title: title || $t('common.confirmDialog.titles.warning'),
    type: 'warning',
  })
}

export function confirmDanger(message: string, title?: string): Promise<boolean | string> {
  const { $t } = useTexts()
  return showConfirmDialog({
    message,
    title: title || $t('common.confirmDialog.titles.danger'),
    type: 'danger',
  })
}

export function confirmInfo(message: string, title?: string): Promise<boolean | string> {
  const { $t } = useTexts()
  return showConfirmDialog({
    message,
    title: title || $t('common.confirmDialog.titles.info'),
    type: 'info',
  })
}

export function confirmSuccess(message: string, title?: string): Promise<boolean | string> {
  const { $t } = useTexts()
  return showConfirmDialog({
    message,
    title: title || $t('common.confirmDialog.titles.success'),
    type: 'success',
  })
}

export default useConfirmDialog
