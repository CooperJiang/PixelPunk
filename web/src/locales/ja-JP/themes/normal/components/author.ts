/**
 * Author info component
 */
export const author = {
  loading: '読み込み中...',
  header: {
    joinedLabel: '{days}日前に参加',
    websiteLabel: 'ウェブサイト',
    stats: {
      views: '閲覧数',
      shares: '共有',
      images: 'ファイル',
    },
  },
  folders: {
    rootTitle: 'フォルダ',
    rootSubtitle: '著者のフォルダを閲覧',
    childTitle: 'サブフォルダ',
    childSubtitle: '現在のフォルダ内のサブフォルダ',
    countLabel: 'ファイル',
    createdLabel: '{time}',
  },
  images: {
    sectionTitle: 'ファイル',
    sectionSubtitle: '著者のファイル作品を閲覧',
    countLabel: 'リソース',
    loadingText: 'ファイルを読み込み中...',
    paginationInfo: '{total}ページ中{current}ページ',
    actions: {
      preview: {
        label: 'プレビュー',
        tooltip: 'クリックしてファイルをプレビュー',
      },
      download: {
        label: 'ダウンロード',
        tooltip: 'オリジナルをダウンロード',
      },
      detail: {
        label: '詳細',
        tooltip: 'ファイルリンクをコピー',
      },
    },
  },
  shares: {
    sectionTitle: '共有',
    sectionSubtitle: '著者の公開共有を閲覧',
    statusBadge: 'ステータス',
    meta: {
      viewsLabel: '閲覧数',
      createdLabel: '{time}',
      keyLabel: '共有キー',
    },
    empty: {
      title: '共有がありません',
      description: 'この著者はまだ公開共有を作成していません',
    },
  },
  empty: {
    title: 'コンテンツがありません',
    description: 'この著者はまだコンテンツをアップロードしていません',
    suggestion: '後でもう一度確認してください',
    retry: '再読み込み',
  },
  error: {
    title: '読み込み失敗',
    retry: '再試行',
    invalidAuthor: '無効な著者ID',
    loadFailed: '著者情報の読み込みに失敗しました',
    loadFolderFailed: 'フォルダコンテンツの読み込みに失敗しました',
  },
  toolbar: {
    searchPlaceholder: 'コンテンツを検索...',
    refresh: '更新',
    gridView: 'グリッド表示',
    listView: 'リスト表示',
    fullscreenEnter: '全画面',
    fullscreenExit: '全画面終了',
  },
  breadcrumb: {
    root: 'ルート',
  },
  toast: {
    copySuccess: 'リンクをクリップボードにコピー',
    copyFailed: 'コピー失敗',
  },
}

