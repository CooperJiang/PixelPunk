/**
 * Dashboard Page Text - Cyber Style
 */
export const dashboard = {
  quickActions: {
    title: 'クイックアクション',
    dragHint: 'ドラッグして並べ替え',
    actions: {
      myFiles: 'データユニット',
      folders: 'フォルダ',
      tagManage: 'タグ管理',
      categoryManage: 'カテゴリ管理',
      openApi: 'Open API',
      automation: '自動タスク',
    },
    settings: {
      divider: '設定センター',
      api: 'API',
      profile: '個人情報',
      security: 'セキュリティ',
      accessControl: 'ホットリンク保護',
      preferences: '設定',
    },
  },
  stats: {
    title: 'データ概要',
    items: {
      images: 'リソース数',
      storage: '使用ストレージ',
      views: 'アクセス数',
      shares: '共有数',
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
        danger: '重大',
      },
      description: {
        healthy: 'リソースが十分です。安心してご利用ください',
        warning: 'リソース使用量が中程度です。不要なコンテンツを定期的にクリーンアップすることをお勧めします',
        danger: 'リソース使用量が重大です。コンテンツをクリーンアップするか、管理者に連絡して拡張してください',
      },
    },
  },
  messages: {
    title: 'メッセージ通知',
    unread: '未読',
    viewAll: 'すべて表示',
    loading: 'メッセージを読み込み中...',
    empty: 'メッセージなし',
    loadingMore: 'さらに読み込み中...',
    loadingMoreDesc: 'お待ちください',
    noMore: 'すべてのメッセージを読み込み済み',
    noMoreDesc: '合計{count}件のメッセージ',
    scrollHint: 'スクロールしてさらに表示',
    markReadFailed: '既読にマークできませんでした。再試行してください',
  },
  errors: {
    fetchStatsFailed: 'データの取得に失敗しました',
    networkFailed: 'ネットワークリクエスト失敗',
  },
  activityMonitor: {
    title: 'アクティビティモニター',
    todayUploads: '今日のアップロード',
    totalViews: '合計ビュー',
    recentActivity: '最近のアクティビティ',
    recordCount: '合計{count}件の記録',
    loadingActivities: 'アクティビティログを読み込み中...',
    loadingMore: 'さらに読み込み中...',
    loadingHint: 'お待ちください',
    allLoaded: 'すべてのアクティビティ記録を読み込み済み',
    empty: 'アクティビティ記録なし',
    emptyHint: '使用を開始すると、ここにアクティビティが表示されます',
    status: {
      active: '正常監視',
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
      noImages: 'コピーする画像なし',
      copySuccess: '{count}個のデータユニットリンクをコピーしました',
      copyFailed: 'コピー失敗',
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
    overallProgress: '全体進捗',
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
