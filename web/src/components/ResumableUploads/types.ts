/* ResumableUploads 组件类型定义 */
import type { UploadSessionData } from '@/utils/storage/uploadStorage'

export interface ResumableUploadsProps {
  modelValue: boolean
}

export interface ResumableUploadsEmits {
  'update:modelValue': [value: boolean]
  resume: [sessions: UploadSessionData[]]
}
