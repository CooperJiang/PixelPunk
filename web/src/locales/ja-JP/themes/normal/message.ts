/**
 * Message Translations - Japanese Normal Theme
 */
export const message = {
  system: {
    maintenance: {
      title: 'システムメンテナンスのお知らせ',
      content: 'システムは {time} にメンテナンスを実施します。所要時間は {duration} を予定しております。メンテナンス中は一部機能がご利用いただけない場合がございますので、事前にご準備ください。',
    },
    announcement: {
      title: 'システムのお知らせ',
      content: '{content}',
    },
  },
  account: {
    register: {
      title: 'PixelPunk へようこそ',
      content: 'PixelPunk にご登録いただきありがとうございます！アカウントが正常に作成されました。{storage} のストレージ容量をご用意しました。画像管理の旅を始めましょう！',
    },
    storageGranted: {
      title: 'ストレージ容量が増加しました',
      content: 'おめでとうございます！ストレージ容量が {size} 増加し、現在の総容量は {total} です。',
    },
    bandwidthGranted: {
      title: '帯域幅が増加しました',
      content: 'おめでとうございます！毎月の帯域幅が {size} 増加し、現在の月間帯域幅は {total} です。',
    },
  },
  content: {
    reviewApproved: {
      title: 'コンテンツ審査が承認されました',
      content: 'コンテンツ「{content_name}」が承認され、公開されました。',
    },
    reviewRejected: {
      title: 'コンテンツ審査が承認されませんでした',
      content: '申し訳ございません。コンテンツ「{content_name}」は承認されませんでした。理由：{reason}',
    },
    reviewPending: {
      title: 'コンテンツ審査中',
      content: 'コンテンツ「{content_name}」の審査が提出されました。できるだけ早く処理いたします。',
    },
  },
  storage: {
    quotaWarning: {
      title: 'ストレージ容量不足',
      content: 'ストレージ容量の {used_percent}%（{used}/{total}）を使用しています。容量の整理またはアップグレードを行ってください。',
    },
    quotaIncreased: {
      title: 'ストレージ容量がアップグレードされました',
      content: 'ストレージ容量が {old_size} から {new_size} にアップグレードされました。',
    },
    quotaDecreased: {
      title: 'ストレージ容量が調整されました',
      content: 'ストレージ容量が {old_size} から {new_size} に調整されました。',
    },
  },
  file: {
    deletedByAdmin: {
      title: 'ファイルが削除されました',
      content: 'ファイル「{file_name}」がコミュニティ規則違反により管理者に削除されました。',
    },
    batchDeletedByAdmin: {
      title: 'ファイルが一括削除されました',
      content: '{count} 個のファイルがコミュニティ規則違反により管理者に削除されました。',
    },
    hardDeletedByAdmin: {
      title: 'ファイルが完全に削除されました',
      content: 'ファイル「{file_name}」が管理者により完全に削除され、復元できません。',
    },
    expiryWarning: {
      title: '{count} 個のファイルが間もなく期限切れになります',
      content: '以下のファイルが間もなく期限切れになります。バックアップをお取りください：\n{file_list}\n\n期限切れのファイルは自動的に削除され、復元できません。',
    },
    thumbnailFailed: {
      title: 'サムネイル生成に失敗しました',
      content: 'ファイル「{file_name}」のサムネイル生成に失敗しましたが、ファイル自体には影響ありません。',
    },
  },
  security: {
    loginAlert: {
      title: '新しいデバイスからのログイン通知',
      content: '{time} に {ip}（{location}）からアカウントにログインされました。ご本人でない場合は、すぐにパスワードを変更してください。',
    },
  },
  apikey: {
    created: {
      title: 'API キーが作成されました',
      content: '新しい API キー「{key_name}」を作成しました。キーを大切に保管してください。',
    },
    deleted: {
      title: 'API キーが削除されました',
      content: 'API キー「{key_name}」が削除されました。',
    },
    regenerated: {
      title: 'API キーが再生成されました',
      content: 'API キー「{key_name}」が再生成されました。古いキーは即座に無効になります。',
    },
    disabled: {
      title: 'API キーが無効化されました',
      content: 'API キー「{key_name}」が無効化されました。',
    },
    enabled: {
      title: 'API キーが有効化されました',
      content: 'API キー「{key_name}」が有効化されました。',
    },
  },
  randomApi: {
    created: {
      title: 'ランダム画像 API が作成されました',
      content: 'ランダム画像 API「{api_name}」を作成しました。',
    },
    deleted: {
      title: 'ランダム画像 API が削除されました',
      content: 'ランダム画像 API「{api_name}」が削除されました。',
    },
    disabled: {
      title: 'ランダム画像 API が無効化されました',
      content: 'ランダム画像 API「{api_name}」が無効化されました。',
    },
    enabled: {
      title: 'ランダム画像 API が有効化されました',
      content: 'ランダム画像 API「{api_name}」が有効化されました。',
    },
  },
  share: {
    expiryWarning: {
      title: '共有リンクが間もなく期限切れになります',
      content: '共有リンク「{share_name}」は {expires_at} に期限切れになります。',
    },
  },
  messageList: {
    title: 'メッセージセンター',
    total: '合計',
    unread: '未読',
    allStatus: 'すべてのステータス',
    allTypes: 'すべてのタイプ',
    markAllRead: 'すべて既読にする',
    emptyTitle: 'メッセージがありません',
    emptyDescription: 'まだメッセージを受信していません',
    loadingText: '読み込み中...',
    columns: {
      type: 'タイプ',
      title: 'タイトル',
      content: '内容',
      status: 'ステータス',
      priority: '優先度',
      time: '時間',
      actions: '操作',
    },
    statusLabels: {
      unread: '未読',
      read: '既読',
    },
    priorityLabels: {
      high: '高',
      normal: '通常',
      low: '低',
    },
    actions: {
      markRead: '既読にする',
      delete: '削除',
    },
    dialog: {
      deleteTitle: '削除の確認',
      deleteMessage: 'このメッセージを削除してもよろしいですか？この操作は元に戻せません。',
      cancel: 'キャンセル',
      confirmDelete: '削除',
    },
    toasts: {
      deleted: 'メッセージを削除しました',
      deleteFailed: '削除に失敗しました。もう一度お試しください',
      markReadSuccess: '既読にしました',
      markReadFailed: '既読にできませんでした。もう一度お試しください',
      markAllReadSuccess: 'すべてのメッセージを既読にしました',
      markAllReadFailed: '一括既読に失敗しました。もう一度お試しください',
    },
  },
}
