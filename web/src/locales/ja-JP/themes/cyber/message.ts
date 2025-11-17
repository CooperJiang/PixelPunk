/**
 * Message Translations - Japanese Cyber Theme
 */
export const message = {
  system: {
    maintenance: {
      title: 'システムメンテナンスプロトコル',
      content: 'ニューラルネットワークは {time} にメンテナンスモードに入ります。所要時間は {duration} を予定しております。メンテナンス中は一部ノードがオフラインになる可能性がありますので、事前にデータを同期してください。',
    },
    announcement: {
      title: 'システムブロードキャスト',
      content: '{content}',
    },
  },
  account: {
    register: {
      title: 'PixelPunk ニューラルネットワークへようこそ',
      content: '接続が確立されました！ノードが正常に登録されました。システムは {storage} のデータ容量を割り当てました。データ管理の旅を始めましょう！',
    },
    storageGranted: {
      title: 'データ容量が拡張されました',
      content: 'システムアップグレード！データ容量が {size} 拡張され、現在の総容量は {total} です。',
    },
    bandwidthGranted: {
      title: '帯域幅が強化されました',
      content: 'ネットワーク強化！毎月の帯域幅が {size} 増加し、現在の月間帯域幅は {total} です。',
    },
  },
  content: {
    reviewApproved: {
      title: 'データ審査が承認されました',
      content: 'データ「{content_name}」が承認され、ネットワーク上で表示されています。',
    },
    reviewRejected: {
      title: 'データ審査が失敗しました',
      content: '申し訳ございません。データ「{content_name}」は承認されませんでした。理由：{reason}',
    },
    reviewPending: {
      title: 'データ審査中',
      content: 'データ「{content_name}」の審査が提出されました。システムが処理中です。',
    },
  },
  storage: {
    quotaWarning: {
      title: 'データ容量アラート',
      content: 'データ容量の {used_percent}%（{used}/{total}）を使用しています。クリーンアップまたは容量のアップグレードを推奨します。',
    },
    quotaIncreased: {
      title: 'データ容量がアップグレードされました',
      content: 'データ容量が {old_size} から {new_size} にアップグレードされました。',
    },
    quotaDecreased: {
      title: 'データ容量が調整されました',
      content: 'データ容量が {old_size} から {new_size} に調整されました。',
    },
  },
  file: {
    deletedByAdmin: {
      title: 'データが削除されました',
      content: 'データユニット「{file_name}」がネットワークプロトコル違反により管理者に削除されました。',
    },
    batchDeletedByAdmin: {
      title: 'データが一括削除されました',
      content: '{count} 個のデータユニットがネットワークプロトコル違反により管理者に削除されました。',
    },
    hardDeletedByAdmin: {
      title: 'データが完全に削除されました',
      content: 'データユニット「{file_name}」が管理者により完全に削除され、復元できません。',
    },
    expiryWarning: {
      title: '{count} 個のデータユニットが間もなく期限切れになります',
      content: '以下のデータユニットが間もなく期限切れになります。バックアップをお取りください：\n{file_list}\n\n期限切れのデータは自動的にパージされ、復元できません。',
    },
    thumbnailFailed: {
      title: 'プレビュー生成に失敗しました',
      content: 'データユニット「{file_name}」のプレビュー生成に失敗しましたが、データ自体には影響ありません。',
    },
  },
  security: {
    loginAlert: {
      title: '新しいノード接続アラート',
      content: '{time} に {ip}（{location}）からアカウントが接続されました。ご本人でない場合は、すぐにアクセスキーを変更してください。',
    },
  },
  apikey: {
    created: {
      title: 'API キーが生成されました',
      content: '新しい API キー「{key_name}」を生成しました。アクセス認証情報を大切に保管してください。',
    },
    deleted: {
      title: 'API キーが破棄されました',
      content: 'API キー「{key_name}」が破棄されました。',
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
      title: 'API キーがアクティブ化されました',
      content: 'API キー「{key_name}」がアクティブ化されました。',
    },
  },
  randomApi: {
    created: {
      title: 'ランダムデータ API が作成されました',
      content: 'ランダムデータ API「{api_name}」を作成しました。',
    },
    deleted: {
      title: 'ランダムデータ API が削除されました',
      content: 'ランダムデータ API「{api_name}」が削除されました。',
    },
    disabled: {
      title: 'ランダムデータ API が無効化されました',
      content: 'ランダムデータ API「{api_name}」が無効化されました。',
    },
    enabled: {
      title: 'ランダムデータ API がアクティブ化されました',
      content: 'ランダムデータ API「{api_name}」がアクティブ化されました。',
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
