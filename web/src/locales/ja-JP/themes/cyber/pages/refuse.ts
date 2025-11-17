/**
 * Access Denied Page Text - Cyber Style
 */
export const refuse = {
  defaultMessage: 'このシステムにアクセスする権限がありません',
  unknown: '不明',
  protocolHeader: 'システムセキュリティプロトコル有効化 / システムセキュリティプロトコル',
  titles: {
    accountDisabled: 'アカウント無効',
    accessDenied: 'アクセス拒否',
  },
  analysisTitle: {
    account: 'アカウントステータス説明:',
    system: 'システム分析結果:',
  },
  messages: {
    ipNotInWhitelist: 'あなたのIPはアクセスホワイトリストにありません',
    ipInBlacklist: 'あなたのIPはシステムブラックリストに追加されました',
    domainNotInWhitelist: 'あなたのドメインはアクセスホワイトリストにありません',
    domainInBlacklist: 'あなたのドメインはシステムブラックリストに追加されました',
    accountDisabled: 'あなたのアカウントは無効化されました',
  },
  ip: {
    notInWhitelist: 'あなたのIPアドレスはシステムアクセスホワイトリストにありません',
    inBlacklist: 'または、あなたのIPはシステムブラックリストに追加されました',
    contactAdmin: 'システム管理者に連絡して、あなたのIPをホワイトリストに追加してください',
  },
  domain: {
    notInWhitelist: 'あなたのアクセスドメインはシステムホワイトリストにありません',
    inBlacklist: 'または、あなたのドメインはシステムブラックリストに追加されました',
    contactAdmin: 'システム管理者に連絡して、ドメインアクセス権限を設定してください',
  },
  status: {
    accountStatus: 'アカウントステータス:',
    currentIp: '現在のIPアドレス:',
    currentDomain: '現在のドメイン:',
    disabled: '無効',
    detectionTime: '検出時間:',
  },
  accountDetails: {
    disabled: 'あなたのアカウントはシステム管理者によって無効化されました',
    reason: 'アカウント無効化は、プラットフォームルール違反またはその他のセキュリティ上の理由による可能性があります',
    contact: '制限を解除するには、システム管理者に連絡して支援を求めてください',
  },
}
