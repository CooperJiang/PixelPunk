/**
 * Author Info Component
 */
export const author = {
  loading: 'データを読み込み中...',
  header: {
    joinedLabel: '{days}日前に接続',
    websiteLabel: 'ネットワークアドレス',
    stats: {
      views: 'アクセス数',
      shares: 'シェアデータ',
      images: 'ファイルストレージ',
    },
  },
  folders: {
    rootTitle: 'データコンテナ',
    rootSubtitle: 'ユーザーデータストレージ構造を閲覧',
    childTitle: 'サブコンテナ',
    childSubtitle: '現在のコンテナ内のサブレベルストレージユニット',
    countLabel: ' データユニット',
    createdLabel: '{time}',
  },
  images: {
    sectionTitle: 'ファイルデータベース',
    sectionSubtitle: 'ユーザーファイルデータセットを閲覧',
    countLabel: ' リソースデータ',
    loadingText: 'ファイルデータを読み込み中...',
    paginationInfo: '{total}セクション中{current}セクション',
    actions: {
      preview: {
        label: 'データプレビュー',
        tooltip: 'ファイルプレビューモードを起動',
      },
      download: {
        label: 'データダウンロード',
        tooltip: 'オリジナルファイルデータを取得',
      },
      detail: {
        label: '詳細',
        tooltip: 'データリンクをコピー',
      },
    },
  },
  shares: {
    sectionTitle: 'シェアネットワーク',
    sectionSubtitle: 'ユーザーパブリックデータシェアを閲覧',
    statusBadge: '接続ステータス',
    meta: {
      viewsLabel: 'アクセス数',
      createdLabel: '{time}',
      keyLabel: 'アクセスキー',
    },
    empty: {
      title: 'ネットワークスペース空',
      description: 'ユーザーはパブリックデータ接続を確立していません',
    },
  },
  empty: {
    title: 'データベース空',
    description: 'ユーザーはデータコンテンツをアップロードしていません',
    suggestion: '後で再スキャン',
    retry: '再接続',
  },
  error: {
    title: '接続失敗',
    retry: '再試行',
    invalidAuthor: '無効なユーザーID',
    loadFailed: 'ユーザープロフィール読み込み失敗',
    loadFolderFailed: 'データコンテナ読み込み失敗',
  },
  toolbar: {
    searchPlaceholder: 'データ検索...',
    refresh: '再同期',
    gridView: 'マトリックス表示',
    listView: 'リスト表示',
    fullscreenEnter: 'フルスクリーンモード',
    fullscreenExit: 'フルスクリーンを終了',
  },
  breadcrumb: {
    root: 'ルートディレクトリ',
  },
  toast: {
    copySuccess: 'データリンクをバッファにコピーしました',
    copyFailed: 'データコピー操作失敗',
  },
}
