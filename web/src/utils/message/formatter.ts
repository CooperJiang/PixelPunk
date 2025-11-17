import type { Message } from '@/api/message'
import type { TranslationFunction } from '@/composables/useTexts'

/**
 * Format message title and content based on type and data
 */
export function formatMessage(
  message: Message,
  $t: TranslationFunction
): {
  title: string
  content: string
} {
  const { type, data } = message

  try {
    if (data.title && data.content) {
      return {
        title: String(data.title),
        content: String(data.content),
      }
    }

    switch (type) {
      case 'system.maintenance':
        return {
          title: $t('message.system.maintenance.title'),
          content: $t('message.system.maintenance.content', data),
        }

      case 'system.announcement':
        return {
          title: $t('message.system.announcement.title'),
          content: $t('message.system.announcement.content', data),
        }

      case 'account.register':
        return {
          title: $t('message.account.register.title'),
          content: $t('message.account.register.content', data),
        }

      case 'account.storage_granted':
        return {
          title: $t('message.account.storageGranted.title'),
          content: $t('message.account.storageGranted.content', data),
        }

      case 'account.bandwidth_granted':
        return {
          title: $t('message.account.bandwidthGranted.title'),
          content: $t('message.account.bandwidthGranted.content', data),
        }

      case 'content.review_approved':
        return {
          title: $t('message.content.reviewApproved.title'),
          content: $t('message.content.reviewApproved.content', data),
        }

      case 'content.review_rejected':
        return {
          title: $t('message.content.reviewRejected.title'),
          content: $t('message.content.reviewRejected.content', data),
        }

      case 'content.review_pending':
        return {
          title: $t('message.content.reviewPending.title'),
          content: $t('message.content.reviewPending.content', data),
        }

      case 'storage.quota_warning':
        return {
          title: $t('message.storage.quotaWarning.title'),
          content: $t('message.storage.quotaWarning.content', data),
        }

      case 'storage.quota_increased':
        return {
          title: $t('message.storage.quotaIncreased.title'),
          content: $t('message.storage.quotaIncreased.content', data),
        }

      case 'storage.quota_decreased':
        return {
          title: $t('message.storage.quotaDecreased.title'),
          content: $t('message.storage.quotaDecreased.content', data),
        }

      case 'file.deleted_by_admin':
        return {
          title: $t('message.file.deletedByAdmin.title'),
          content: $t('message.file.deletedByAdmin.content', data),
        }

      case 'file.batch_deleted_by_admin':
        return {
          title: $t('message.file.batchDeletedByAdmin.title'),
          content: $t('message.file.batchDeletedByAdmin.content', data),
        }

      case 'file.hard_deleted_by_admin':
        return {
          title: $t('message.file.hardDeletedByAdmin.title'),
          content: $t('message.file.hardDeletedByAdmin.content', data),
        }

      case 'file.expiry_warning':
        return {
          title: $t('message.file.expiryWarning.title', { count: data.file_count || 1 }),
          content: $t('message.file.expiryWarning.content', data),
        }

      case 'file.thumbnail_failed':
        return {
          title: $t('message.file.thumbnailFailed.title'),
          content: $t('message.file.thumbnailFailed.content', data),
        }

      case 'security.login_alert':
        return {
          title: $t('message.security.loginAlert.title'),
          content: $t('message.security.loginAlert.content', data),
        }

      case 'apikey.created':
        return {
          title: $t('message.apikey.created.title'),
          content: $t('message.apikey.created.content', data),
        }

      case 'apikey.deleted':
        return {
          title: $t('message.apikey.deleted.title'),
          content: $t('message.apikey.deleted.content', data),
        }

      case 'apikey.regenerated':
        return {
          title: $t('message.apikey.regenerated.title'),
          content: $t('message.apikey.regenerated.content', data),
        }

      case 'apikey.disabled':
        return {
          title: $t('message.apikey.disabled.title'),
          content: $t('message.apikey.disabled.content', data),
        }

      case 'apikey.enabled':
        return {
          title: $t('message.apikey.enabled.title'),
          content: $t('message.apikey.enabled.content', data),
        }

      case 'random_api.created':
        return {
          title: $t('message.randomApi.created.title'),
          content: $t('message.randomApi.created.content', data),
        }

      case 'random_api.deleted':
        return {
          title: $t('message.randomApi.deleted.title'),
          content: $t('message.randomApi.deleted.content', data),
        }

      case 'random_api.disabled':
        return {
          title: $t('message.randomApi.disabled.title'),
          content: $t('message.randomApi.disabled.content', data),
        }

      case 'random_api.enabled':
        return {
          title: $t('message.randomApi.enabled.title'),
          content: $t('message.randomApi.enabled.content', data),
        }

      case 'share.expiry_warning':
        return {
          title: $t('message.share.expiryWarning.title'),
          content: $t('message.share.expiryWarning.content', data),
        }

      default:
        return {
          title: `Message: ${type}`,
          content: JSON.stringify(data),
        }
    }
  } catch (error) {
    return {
      title: `Message: ${type}`,
      content: JSON.stringify(data),
    }
  }
}
