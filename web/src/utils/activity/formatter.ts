import type { TranslationFunction } from '@/composables/useTexts'
import type { ActivityLog } from '@/api/activity'
import { ActivityAction } from '@/constants/activity'
import { formatFileSize } from '@/utils/formatting/format'

/**
 * 格式化活动描述
 * 根据活动类型和数据动态生成多语言描述
 */
export function formatActivityDescription(activity: ActivityLog, $t: TranslationFunction): string {
  const { type, data } = activity

  try {
    switch (type) {
      // ==================== 用户相关 ====================
      case 'user_login':
        return data.ip_address
          ? $t('activity.user.loginWithIp', { ip: data.ip_address })
          : $t('activity.user.login')

      case 'user_register':
        return $t('activity.user.register')

      // ==================== 文件相关 ====================
      case 'batch_image_upload':
        return data.total_size
          ? $t('activity.file.batchUploadWithSize', {
              count: data.file_count || data.count || 0,
              size: formatFileSize(data.total_size),
            })
          : $t('activity.file.batchUpload', { count: data.file_count || data.count || 0 })

      case 'file_delete':
        return $t('activity.file.delete', { fileName: data.file_name || 'Unknown' })

      case 'file_rename':
        return $t('activity.file.rename', {
          oldName: data.old_name || 'Unknown',
          newName: data.new_name || 'Unknown',
        })

      case 'file_move':
        return $t('activity.file.move', {
          fileName: data.file_name || 'Unknown',
          oldFolder: data.old_folder_name || $t('activity.common.rootFolder'),
          newFolder: data.new_folder_name || $t('activity.common.rootFolder'),
        })

      case 'image_access_level_change':
      case 'file_access_level_change': {
        const level = $t(`activity.accessLevel.${data.new_access_level || 'private'}`)
        return $t('activity.file.accessLevelChange', {
          fileName: data.file_name || 'Unknown',
          level,
        })
      }

      case 'image_expired':
      case 'file_expired':
        if (data.expired_count && data.expired_count > 1) {
          return $t('activity.file.expiredBatch', { count: data.expired_count })
        }
        return $t('activity.file.expired', { fileName: data.file_name || 'Unknown' })

      // ==================== 文件夹相关 ====================
      case 'folder_create':
        return data.parent_name
          ? $t('activity.folder.createInParent', {
              folderName: data.folder_name || 'Unknown',
              parentName: data.parent_name,
            })
          : $t('activity.folder.create', { folderName: data.folder_name || 'Unknown' })

      case 'folder_rename':
        return $t('activity.folder.rename', {
          oldName: data.old_name || 'Unknown',
          newName: data.new_name || 'Unknown',
        })

      case 'folder_delete':
        return data.file_count && data.file_count > 0
          ? $t('activity.folder.deleteWithFiles', {
              folderName: data.folder_name || 'Unknown',
              count: data.file_count,
            })
          : $t('activity.folder.delete', { folderName: data.folder_name || 'Unknown' })

      case 'folder_access_level_change': {
        const oldLevel = $t(`activity.accessLevel.${data.old_access_level || 'private'}`)
        const newLevel = $t(`activity.accessLevel.${data.new_access_level || 'public'}`)
        return $t('activity.folder.accessLevelChange', {
          folderName: data.folder_name || 'Unknown',
          oldLevel,
          newLevel,
        })
      }

      // ==================== 分享相关 ====================
      case 'share_create':
        return $t('activity.share.create')

      case 'share_delete':
        return $t('activity.share.delete')

      case 'share_milestone':
        return $t('activity.share.milestone', { count: data.visit_count || data.milestone || 0 })

      // ==================== API密钥相关 ====================
      case 'apikey_create':
        return $t('activity.apikey.create', { keyName: data.key_name || 'API Key' })

      case 'apikey_delete':
        return $t('activity.apikey.delete', { keyName: data.key_name || 'API Key' })

      case 'apikey_toggle_status':
        if (data.action === ActivityAction.ENABLE || data.new_status === 1) {
          return $t('activity.apikey.toggleEnable', { keyName: data.key_name || 'API Key' })
        } else {
          return $t('activity.apikey.toggleDisable', { keyName: data.key_name || 'API Key' })
        }

      case 'apikey_regenerate':
        return $t('activity.apikey.regenerate', { keyName: data.key_name || 'API Key' })

      // ==================== 随机API相关 ====================
      case 'random_api_create':
        return $t('activity.randomApi.create', { apiName: data.api_name || 'Random API' })

      case 'random_api_delete':
        return $t('activity.randomApi.delete', { apiName: data.api_name || 'Random API' })

      case 'random_api_toggle_status': {
        const apiName = data.api_name || 'Random API'
        const folderName = data.folder_name

        if (data.action === ActivityAction.ENABLE || data.new_status === 1) {
          return folderName
            ? $t('activity.randomApi.toggleEnableWithFolder', { apiName, folderName })
            : $t('activity.randomApi.toggleEnable', { apiName })
        } else {
          return folderName
            ? $t('activity.randomApi.toggleDisableWithFolder', { apiName, folderName })
            : $t('activity.randomApi.toggleDisable', { apiName })
        }
      }

      // ==================== 个人设置 ====================
      case 'profile_update': {
        const fieldType = data.field_type || 'profile'
        const key = `activity.profile.${fieldType}Update`
        // 尝试获取特定字段的翻译，如果不存在则使用默认
        try {
          return $t(key)
        } catch {
          return $t('activity.profile.update')
        }
      }

      case 'password_change':
        return $t('activity.profile.passwordChange')

      case 'email_change':
        return $t('activity.profile.emailChange', { email: data.new_email || '' })

      case 'hotlink_protection_change':
        if (data.action === ActivityAction.ENABLE) {
          return $t('activity.security.hotlinkProtectionEnable')
        } else if (data.action === ActivityAction.DISABLE) {
          return $t('activity.security.hotlinkProtectionDisable')
        }
        return $t('activity.security.hotlinkProtectionChange')

      // ==================== 系统相关 ====================
      case 'batch_delete':
        return data.folder_name
          ? $t('activity.system.batchDeleteWithFolder', {
              count: data.file_count || data.count || 0,
              folderName: data.folder_name,
            })
          : $t('activity.system.batchDelete', { count: data.file_count || data.count || 0 })

      case 'admin_delete':
        return $t('activity.system.adminDelete')

      case 'system_cleanup':
        return $t('activity.system.cleanup', { count: data.cleaned_count || 0, type: data.cleanup_type || '' })

      case 'guest_file_expired':
      case 'guest_image_expired':
        return $t('activity.system.guestImageExpired', { count: data.expired_count || 0 })

      // ==================== 未知类型 ====================
      default:
        console.warn(`Unknown activity type: ${type}`, data)
        return `Activity: ${type}`
    }
  } catch (error) {
    console.error('Error formatting activity description:', error, activity)
    return `Activity: ${type}`
  }
}
