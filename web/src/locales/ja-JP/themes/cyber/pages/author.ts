/**
 * Author Page Text - Cyber Style
 */
export const author = {
  loading: 'データ読み込み中...',

  header: {
    joinedLabel: '接続時間',
    websiteLabel: 'ネットワークノード',
    stats: {
      views: 'アクセス数',
      shares: '共有数',
      images: 'データユニット',
    },
  },

  folders: {
    rootTitle: 'すべてのディレクトリ',
    rootSubtitle: '作成者別のすべてのデータディレクトリを閲覧',
    childTitle: 'サブディレクトリ',
    childSubtitle: '現在のディレクトリ内のサブディレクトリ',
    countLabel: '個',
    createdLabel: '作成日時',
    unknownTime: '不明',
    rootFolder: 'ルートディレクトリ',
  },

  shares: {
    sectionTitle: '公開共有',
    sectionSubtitle: '作成者の公開データ共有',
    statusBadge: 'ステータス',
    meta: {
      viewsLabel: 'アクセス',
      createdLabel: '作成日時',
      keyLabel: '共有キー',
    },
    empty: {
      title: '共有なし',
      description: 'この作成者はまだ公開共有を作成していません',
    },
  },

  toolbar: {
    searchPlaceholder: 'ディレクトリを検索...',
    refresh: '更新',
    fullscreenEnter: 'フルスクリーンモード',
    fullscreenExit: 'フルスクリーンを終了',
  },

  empty: {
    title: 'データなし',
    description: 'この作成者はまだデータをアップロードしていません',
    suggestion: '後でもう一度確認してください',
    retry: '再試行',
    folderEmpty: {
      title: 'ディレクトリが空です',
      description: 'このディレクトリには現在コンテンツがありません',
      suggestion: '戻るか、他のコンテンツを閲覧してください',
    },
  },

  error: {
    title: '読み込み失敗',
    retry: '再試行',
    invalidAuthor: '無効な作成者ID',
    loadFailed: '作成者情報の読み込みに失敗しました',
    loadFolderFailed: 'ディレクトリの読み込みに失敗しました',
  },

  toast: {
    copySuccess: 'コピー成功',
    copyFailed: 'コピー失敗',
  },
}
