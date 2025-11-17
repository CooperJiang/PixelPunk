/**
 * Shares Management - Cyberpunk Style
 */
export const shares = {
  title: 'シェアネットワーク管理',
  subtitle: 'データシェアプロトコル管理およびノードアクセス監視',
  tabs: {
    list: 'プロトコルリスト',
    visitors: 'ノードアクセス記録',
    stats: 'ネットワークデータ統計',
  },
  shareList: {
    buttons: {
      refresh: '再同期',
      filter: 'データフィルター',
    },
    columns: {
      name: 'プロトコル名',
      shareKey: 'アクセスキー',
      expiredAt: 'プロトコル期限',
      views: '接続数',
      createdAt: 'プロトコル確立',
      status: 'プロトコルステータス',
      actions: 'コマンド操作',
    },
    actions: {
      copyLink: 'アクセスリンクをコピー',
      viewContent: 'プロトコルコンテンツを表示',
      disable: 'プロトコルを無効化',
      enable: 'プロトコルを有効化',
      delete: 'プロトコルを終了',
      viewIpInfo: 'ノード情報を表示',
    },
    status: {
      normal: 'オンライン',
      expired: '期限切れ',
      deleted: '終了済み',
      disabled: '無効',
      unknown: '不明なステータス',
    },
    format: {
      unlimited: '永続有効',
      infiniteViews: '∞',
      viewsFormat: '{current} / {max}',
    },
    messages: {
      fetchFailed: 'プロトコルリスト同期失敗',
      updateStatusSuccess: 'プロトコルステータス更新成功',
      updateStatusFailed: 'プロトコルステータス更新失敗',
      deleteSuccess: 'プロトコル終了成功',
      deleteFailed: 'プロトコル終了失敗',
      copySuccess: 'アクセスリンクをデータバッファにコピーしました',
      copySuccessNormal: 'アクセスリンクをデータバッファにコピーしました（キー検証が必要）',
      copyFailed: 'データコピー失敗、手動で実行してください',
      getAccessFailed: '管理者アクセストークン取得失敗、キー検証が必要',
    },
    loading: {
      text: 'プロトコルデータ同期中...',
    },
    empty: {
      title: 'プロトコルリスト空',
      description: '現在シェアプロトコルがありません。ノード確立プロトコルがここに表示されます',
      resetButton: 'フィルターコマンドをリセット',
    },
    statusDialog: {
      titleEnable: 'プロトコルを有効化',
      titleUpdate: 'プロトコルステータスを更新',
      message: 'このプロトコルステータスを変更してもよろしいですか',
      messageHighlight: '{status}',
      reasonLabel: 'ステータス変更理由（オプション）',
      reasonPlaceholder: 'ステータス変更理由を入力...',
      cancel: '操作を中止',
      confirm: '実行を確認',
    },
    deleteDialog: {
      title: '終了を確認',
      message: 'シェアプロトコルを終了してもよろしいですか',
      messageHighlight: '{name}',
      forceDelete: '強制パージ（復元不可）',
      cancel: '操作を中止',
      confirm: '終了を確認',
    },
  },
  visitorList: {
    ipInfoTitle: 'IP情報を表示',
    filter: {
      shareIdPlaceholder: 'プロトコルIDフィルター',
      keywordPlaceholder: 'ノード名、連絡先メールまたはノードアドレスを検索...',
      searchButton: 'データ検索',
    },
    columns: {
      id: 'ノードID',
      visitorName: 'ノード名',
      visitorEmail: '連絡先メール',
      createdAt: '初回接続',
      ipAddress: 'ノードアドレス',
      shareKey: 'プロトコルに属する',
      visitCount: '接続数',
      actions: 'コマンド操作',
    },
    actions: {
      delete: 'データパージ',
    },
    format: {
      unknown: '不明なノード',
      empty: '-',
    },
    messages: {
      fetchFailed: 'ノードアクセス記録同期失敗',
      deleteSuccess: 'アクセス記録パージ成功',
      deleteFailed: 'アクセス記録パージ失敗',
    },
    loading: {
      text: 'ノードアクセスデータ同期中...',
    },
    empty: {
      title: 'アクセス記録空',
      description: '現在ノードアクセス記録がありません。ノードプロトコル接続がここに表示されます',
    },
    deleteDialog: {
      title: 'パージを確認',
      message: 'このノードアクセス記録をパージしてもよろしいですか？',
      visitorName: 'ノード名:',
      email: '連絡先メール:',
      visitTime: '接続時間:',
      cancel: '操作を中止',
      confirm: 'パージを確認',
    },
  },
  statistics: {
    cards: {
      totalShares: '総プロトコル数',
      activeShares: 'オンラインプロトコル',
      viewsToday: '本日接続',
    },
    charts: {
      trend: {
        title: '接続傾向分析',
        viewsBar: '接続数',
        sharesBar: '新規プロトコル',
        viewsTooltip: '{count}接続',
        sharesTooltip: '{count}新規プロトコル',
      },
      popular: {
        title: '人気プロトコルランキング',
        refreshButton: '再同期',
      },
    },
    periodOptions: {
      last7Days: '過去7日間サイクル',
      last30Days: '過去30日間サイクル',
      last90Days: '過去90日間サイクル',
    },
    legend: {
      views: '接続数',
      shares: '新規プロトコル',
    },
    messages: {
      fetchFailed: '統計データ同期失敗',
    },
    loading: {
      text: 'データ同期中...',
    },
    empty: {
      noData: 'データギャップ',
      noPopular: '人気プロトコルデータギャップ',
    },
  },
}
