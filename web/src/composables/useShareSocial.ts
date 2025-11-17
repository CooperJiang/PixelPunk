import { computed } from 'vue'
import { canUseShareAPI, copyShareLink, shareToQQ, shareToSystem, shareToWeChat, shareToWeibo } from '@/utils/business/shareUtils'
import { useTexts } from '@/composables/useTexts'

/**
 * 分享社交功能
 * 处理链接分享、社交媒体分享等功能
 */
export function useShareSocial(shareData: any, onSuccess?: (message: string) => void, onError?: (message: string) => void) {
  const { $t } = useTexts()

  const canShare = computed(() => canUseShareAPI())

  const handleCopyShareLink = () => {
    const shareUrl = window.location.href

    copyShareLink(
      shareUrl,
      () => onSuccess?.($t('share.social.linkCopied')),
      () => onError?.($t('share.social.copyFailed'))
    )
  }

  const handleShareToSystem = async () => {
    const shareTitle = shareData?.share?.name || $t('share.social.defaultTitle')
    const shareText = shareData?.share?.description || $t('share.social.defaultText')
    const shareUrl = window.location.href

    shareToSystem(
      { title: shareTitle, text: shareText, url: shareUrl },
      () => onSuccess?.($t('share.social.shareSuccess')),
      () => onError?.($t('share.social.shareFailed'))
    )
  }

  const handleShareToWeChat = () => {
    const shareUrl = window.location.href
    const shareTitle = shareData?.share?.name || $t('share.social.defaultTitle')

    shareToWeChat(shareUrl, shareTitle)
    onSuccess?.($t('share.social.wechatQRCode'))
  }

  const handleShareToWeibo = () => {
    const shareUrl = window.location.href
    const shareTitle = shareData?.share?.name || $t('share.social.defaultTitle')
    const shareDesc = shareData?.share?.description || ''

    shareToWeibo(shareUrl, shareTitle, shareDesc)
  }

  const handleShareToQQ = () => {
    const shareUrl = window.location.href
    const shareTitle = shareData?.share?.name || $t('share.social.defaultTitle')
    const shareDesc = shareData?.share?.description || ''

    shareToQQ(shareUrl, shareTitle, shareDesc)
  }

  return {
    canShare,
    handleCopyShareLink,
    handleShareToSystem,
    handleShareToWeChat,
    handleShareToWeibo,
    handleShareToQQ,
  }
}
