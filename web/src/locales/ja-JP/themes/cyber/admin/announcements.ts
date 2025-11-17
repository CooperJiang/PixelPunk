/**
 * Announcement Management Module
 */
export const announcements = {
  title: '█ お知らせ管理 █',
  subtitle: 'システムお知らせを管理、Markdown & HTMLをサポート',
  actions: {
    create: 'お知らせを作成',
    refresh: '更新',
    filter: 'フィルター',
    edit: '編集',
    delete: '削除',
    pin: 'ピン留め',
    unpin: 'ピン留め解除',
    cancel: 'キャンセル',
    update: '更新',
  },
  status: {
    draft: '下書き',
    published: '公開済み',
    archived: 'アーカイブ済み',
  },
  table: {
    title: 'タイトル',
    summary: '概要',
    status: 'ステータス',
    viewCount: '閲覧数',
    publishedAt: '公開日時',
    createdAt: '作成日時',
    actions: '操作',
    noSummary: '概要なし',
    notSet: '未設定',
    loadingText: '>>> 読み込み中...',
  },
  empty: {
    title: 'お知らせなし',
    description: '「お知らせを作成」をクリックして最初のお知らせを作成してください',
  },
  form: {
    title: {
      label: 'タイトル',
      placeholder: 'お知らせタイトルを入力',
    },
    summary: {
      label: '概要',
      placeholder: '概要を入力（オプション）',
    },
    content: {
      label: 'コンテンツ',
      placeholder: 'Markdownコンテンツを入力...',
      hint: 'Markdownをサポート、ライブプレビュー付き。画像をドラッグまたは貼り付けてアップロード',
    },
    status: {
      label: 'ステータス',
    },
    options: {
      label: 'オプション',
      pinned: 'トップにピン留め',
      pinnedHint: '初回訪問時に自動ポップアップ',
    },
  },
  dialog: {
    createTitle: '◢ お知らせを作成',
    editTitle: '◢ お知らせを編集',
  },
  messages: {
    fetchError: 'お知らせの読み込みに失敗しました: ',
    deleteConfirm: 'このお知らせを削除しますか？この操作は元に戻せません。',
    deleteSuccess: '✓ 削除成功',
    deleteError: '✗ 削除失敗: ',
    pinSuccess: '✓ {action}成功',
    pinError: '✗ {action}失敗: ',
    createSuccess: '✓ お知らせ作成済み',
    createError: '✗ 作成失敗: ',
    updateSuccess: '✓ お知らせ更新済み',
    updateError: '✗ 更新失敗: ',
    operationError: '✗ 操作失敗: ',
    uploadSuccess: '✓ {count}枚の画像をアップロードしました',
    uploadError: '✗ 画像アップロード失敗: ',
    uploadFailed: 'アップロード失敗',
  },
  validation: {
    titleRequired: 'タイトルは必須です',
    contentRequired: 'コンテンツは必須です',
  },
  filter: {
    status: {
      label: 'ステータス',
      all: 'すべて',
    },
    pinned: {
      label: 'ピン留め',
      all: 'すべて',
      onlyPinned: 'ピン留めのみ',
      notPinned: 'ピン留めなし',
    },
    search: {
      label: '検索',
      placeholder: 'タイトルまたは概要を検索...',
    },
    actions: {
      query: 'クエリ',
      reset: 'リセット',
    },
  },
  settings: {
    systemToggle: 'システムスイッチ',
    drawerPosition: {
      label: 'ドロワー位置',
      left: '左',
      right: '右',
    },
    displayLimit: '表示制限',
    autoShowDelay: 'ポップアップ遅延',
    saveButton: '設定を保存',
    messages: {
      fetchError: '設定の取得に失敗しました:',
      limitRange: '表示制限は1-20の間である必要があります',
      saveSuccess: '✓ 設定保存済み',
      saveError: '✗ 保存失敗: ',
    },
    units: {
      second: '秒',
    },
  },
}
