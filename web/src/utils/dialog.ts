/**
 * 确认对话框工具函数
 * 使用自定义 CyberConfirmDialog 组件替代原生对话框
 * @deprecated 建议直接使用 useConfirmDialog composable
 * 此文件保留用于向后兼容

 */
import { showConfirmDialog } from '@/composables/useConfirmDialog'

/**
 * 显示确认对话框（兼容旧版本）
 * @param message 确认消息
 * @returns 用户是否确认
 * @deprecated 建议使用 showConfirmDialog 或 confirmWarning/confirmDanger 等方法
 */
export function showConfirm(message: string): Promise<boolean> {
  return showConfirmDialog({
    message,
    type: 'warning',
  }).then((result) => result === true)
}

export async function confirmAsync(message: string): Promise<boolean> {
  return showConfirm(message)
}

export function showAlert(message: string): void {
  import('@/components/Toast/useToast')
    .then(({ useToast }) => {
      const toast = useToast()
      toast.warning(message)
    })
    .catch(() => {
      /* eslint-disable-next-line no-alert */
      window.alert(message)
    })
}
