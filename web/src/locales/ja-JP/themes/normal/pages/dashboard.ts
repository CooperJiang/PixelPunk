/**
 * Dashboard page texts
 */
export const dashboard = {
  quickActions: {
    title: 'クイックアクション',
    dragHint: 'ドラッグして並び替え',
    actions: {
      myFiles: 'マイファイル',
      folders: 'フォルダ',
      tagManage: 'タグ管理',
      categoryManage: 'カテゴリ管理',
      openApi: 'オープンAPI',
      automation: '自動化タスク',
    },
    settings: {
      divider: '設定センター',
      api: 'API',
      profile: 'プロフィール',
      security: 'セキュリティ',
      accessControl: 'ホットリンク保護',
      preferences: '設定',
    },
  },
  stats: {
    title: 'データ概要',
    items: {
      images: 'リソース',
      storage: '使用ストレージ',
      views: '閲覧数',
      shares: 'シェア',
    },
    quotas: {
      title: 'クォータ概要',
      storage: {
        label: 'ストレージクォータ',
        used: '使用済み',
        total: '合計',
      },
      bandwidth: {
        label: '帯域幅クォータ',
        used: '使用済み',
        total: '合計',
      },
      status: {
        healthy: '十分',
        warning: '中程度',
        danger: '深刻',
      },
      description: {
        healthy: 'リソースは十分です。安心して使用できます',
        warning: 'リソース使用量は中程度です。定期的に不要なコンテンツをクリーンアップすることを検討してください',
        danger: 'リソース使用量が深刻です。コンテンツをクリーンアップするか、管理者に連絡して拡張してください',
      },
    },
  },
  messages: {
    title: '通知',
    unread: '未読',
    viewAll: 'すべて表示',
    loading: 'メッセージを読み込み中...',
    empty: 'メッセージなし',
    loadingMore: 'さらに読み込み中...',
    loadingMoreDesc: 'お待ちください',
    noMore: 'すべてのメッセージを読み込みました',
    noMoreDesc: '合計{count}件のメッセージ',
    scrollHint: 'スクロールしてさらに表示',
    markReadFailed: '既読にマークできませんでした。もう一度お試しください',
  },
  errors: {
    fetchStatsFailed: 'データの取得に失敗しました',
    networkFailed: 'ネットワークリクエストに失敗しました',
  },
  activityMonitor: {
    title: 'アクティビティモニター',
    todayUploads: '本日のアップロード',
    totalViews: '総閲覧数',
    recentActivity: '最近のアクティビティ',
    recordCount: '合計{count}件の記録',
    loadingActivities: 'アクティビティログを読み込み中...',
    loadingMore: 'さらに読み込み中...',
    loadingHint: 'お待ちください',
    allLoaded: 'すべてのアクティビティ記録を読み込みました',
    empty: 'アクティビティ記録なし',
    emptyHint: '使用を開始すると、ここにアクティビティが表示されます',
    status: {
      active: '監視中',
      warning: '警告',
      error: 'エラー',
      offline: 'オフライン',
    },
    toast: {
      fetchFailed: 'アクティビティログの取得に失敗しました',
      noData: 'リアルタイムアクティビティデータを一時的に取得できません',
    },
    timeAgo: {
      justNow: 'たった今',
      hoursAgo: '{hours}時間前',
      minutesAgo: '{minutes}分前',
      daysAgo: '{days}日前',
    },
  },
  recentUploads: {
    title: '最近のアップロード',
    loading: '読み込み中...',
    empty: 'アップロード記録なし',
    toast: {
      noImages: 'コピーする画像がありません',
      copySuccess: '{count}件の画像リンクをコピーしました',
      copyFailed: 'コピーに失敗しました',
      fetchFailed: 'データの取得に失敗しました',
    },
    actions: {
      copyAll: 'すべてのリンクをコピー',
      viewAll: 'すべて表示',
    },
  },
  uploadQueue: {
    title: 'アップロードキュー',
    statsTitle: 'アップロード統計',
    totalFiles: '合計{count}ファイル',
    overallProgress: '全体の進捗',
    queueTotal: 'キュー合計',
    uploading: 'アップロード中',
    pending: '待機中',
    completed: '完了',
    failed: '失敗',
    actions: {
      goToUpload: 'アップロードページに移動',
    },
  },
  uploadEntry: {
    title: {
      default: 'ファイルをアップロード',
      withCount: 'キューに{count}ファイル',
    },
  },
  folderDistribution: {
    title: 'フォルダ分布',
    description: 'ブロックサイズはストレージスペースを表し、色の深度は割合を表します',
    totalFolders: '合計フォルダ数: {count}',
    empty: 'フォルダなし',
    tooltipSize: 'サイズ: {size}',
    tooltipFiles: 'ファイル: {count}',
    tooltipPercent: '割合: {percent}%',
    unnamed: '無名',
  },
}
