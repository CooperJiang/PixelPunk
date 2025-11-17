export const shares = {
  title: '共有管理',
  subtitle: 'ユーザー共有リンクと訪問者分析の管理',
  tabs: {
    list: '共有リスト',
    visitors: '訪問者',
    stats: '統計',
  },
  shareList: {
    buttons: {
      refresh: '更新',
      filter: 'フィルター',
    },
    columns: {
      name: '名前',
      shareKey: '共有キー',
      expiredAt: '有効期限',
      views: '閲覧数',
      createdAt: '作成日',
      status: 'ステータス',
      actions: '操作',
    },
    actions: {
      copyLink: 'リンクコピー',
      viewContent: 'コンテンツ表示',
      disable: '無効化',
      enable: '有効化',
      delete: '削除',
      viewIpInfo: 'IP情報表示',
    },
    status: {
      normal: '通常',
      expired: '期限切れ',
      deleted: '削除済み',
      disabled: '無効',
      unknown: '不明',
    },
    format: {
      unlimited: '無制限',
      infiniteViews: '∞',
      viewsFormat: '{current} / {max}',
    },
    messages: {
      fetchFailed: '共有取得失敗',
      updateStatusSuccess: '共有ステータス更新済み',
      updateStatusFailed: 'ステータス更新失敗',
      deleteSuccess: '共有削除済み',
      deleteFailed: '削除失敗',
      copySuccess: '共有リンクコピー済み',
      copySuccessNormal: 'コピー済み（通常リンク、パスワードが必要な場合があります）',
      copyFailed: 'コピー失敗、手動でコピーしてください',
      getAccessFailed: '管理者アクセス取得失敗、パスワードが必要な場合があります',
    },
    loading: {
      text: '共有を読み込み中...',
    },
    empty: {
      title: '共有がありません',
      description: 'まだ共有リンクがありません。新しく作成された共有がここに表示されます。',
      resetButton: 'フィルターリセット',
    },
    statusDialog: {
      titleEnable: '共有有効化',
      titleUpdate: 'ステータス更新',
      message: 'ステータスを変更',
      messageHighlight: '{status}',
      reasonLabel: '理由（オプション）',
      reasonPlaceholder: '理由を入力...',
      cancel: 'キャンセル',
      confirm: '確認',
    },
    deleteDialog: {
      title: '削除の確認',
      message: '共有を削除',
      messageHighlight: '{name}',
      forceDelete: '完全削除（元に戻せません）',
      cancel: 'キャンセル',
      confirm: '削除',
    },
  },
  visitorList: {
    ipInfoTitle: 'IP情報表示',
    filter: {
      shareIdPlaceholder: '共有IDでフィルター',
      keywordPlaceholder: '名前、メールまたはIPで検索...',
      searchButton: '検索',
    },
    columns: {
      id: 'ID',
      visitorName: '訪問者',
      visitorEmail: 'メール',
      createdAt: '訪問日時',
      ipAddress: 'IPアドレス',
      shareKey: '共有',
      visitCount: '訪問数',
      actions: '操作',
    },
    actions: {
      delete: '削除',
    },
    format: {
      unknown: '不明',
      empty: '-',
    },
    messages: {
      fetchFailed: '訪問者取得失敗',
      deleteSuccess: '訪問者削除済み',
      deleteFailed: '削除失敗',
    },
    loading: {
      text: '訪問者を読み込み中...',
    },
    empty: {
      title: '訪問者がありません',
      description: 'まだ訪問者レコードがありません。訪問がここに表示されます。',
    },
    deleteDialog: {
      title: '削除の確認',
      message: 'この訪問者レコードを削除しますか？',
      visitorName: '名前：',
      email: 'メール：',
      visitTime: '訪問日時：',
      cancel: 'キャンセル',
      confirm: '削除',
    },
  },
  statistics: {
    cards: {
      totalShares: '総共有数',
      activeShares: 'アクティブ共有',
      viewsToday: '本日の閲覧数',
    },
    charts: {
      trend: {
        title: 'トラフィック傾向',
        viewsBar: '閲覧数',
        sharesBar: '新規共有',
        viewsTooltip: '{count}回閲覧',
        sharesTooltip: '{count}件の新規共有',
      },
      popular: {
        title: '人気の共有',
        refreshButton: '更新',
      },
    },
    periodOptions: {
      last7Days: '過去7日',
      last30Days: '過去30日',
      last90Days: '過去90日',
    },
    legend: {
      views: '閲覧数',
      shares: '新規共有',
    },
    messages: {
      fetchFailed: '統計取得失敗',
    },
    loading: {
      text: '読み込み中...',
    },
    empty: {
      noData: 'データがありません',
      noPopular: 'まだ人気の共有がありません',
    },
  },
}
