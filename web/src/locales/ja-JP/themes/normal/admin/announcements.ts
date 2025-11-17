export const announcements = {
  title: 'お知らせ管理',
  subtitle: 'システムお知らせの管理。MarkdownとHTMLをサポート。',

  actions: {
    create: 'お知らせ作成',
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
    published: '公開',
    archived: 'アーカイブ',
  },

  table: {
    title: 'タイトル',
    summary: '概要',
    status: 'ステータス',
    viewCount: '閲覧数',
    publishedAt: '公開日',
    createdAt: '作成日',
    actions: '操作',
    noSummary: '概要なし',
    notSet: '未設定',
    loadingText: '読み込み中...',
  },

  empty: {
    title: 'お知らせがありません',
    description: '「お知らせ作成」をクリックして最初のお知らせを追加',
  },

  form: {
    title: {
      label: 'タイトル',
      placeholder: 'お知らせタイトルを入力',
    },
    summary: {
      label: '概要',
      placeholder: 'お知らせ概要を入力（オプション）',
    },
    content: {
      label: '内容',
      placeholder: 'Markdownで内容を入力...',
      hint: 'Markdownをサポート、ライブプレビュー付き。画像をドラッグ/貼り付けしてアップロード。',
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
    createTitle: 'お知らせ作成',
    editTitle: 'お知らせ編集',
  },

  messages: {
    fetchError: 'お知らせ読み込み失敗：',
    deleteConfirm: 'このお知らせを削除しますか？この操作は元に戻せません。',
    deleteSuccess: '削除成功',
    deleteError: '削除失敗：',
    pinSuccess: '{action}成功',
    pinError: '{action}失敗：',
    createSuccess: 'お知らせ作成済み',
    createError: '作成失敗：',
    updateSuccess: 'お知らせ更新済み',
    updateError: '更新失敗：',
    operationError: '操作失敗：',
    uploadSuccess: '{count}件の画像をアップロード',
    uploadError: '画像アップロード失敗：',
    uploadFailed: 'アップロード失敗',
  },

  validation: {
    titleRequired: 'タイトルを入力してください',
    contentRequired: '内容を入力してください',
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
      placeholder: 'タイトルまたは概要で検索...',
    },
    actions: {
      query: '検索',
      reset: 'リセット',
    },
  },

  settings: {
    systemToggle: 'システム切り替え',
    drawerPosition: {
      label: 'ドロワー位置',
      left: '左',
      right: '右',
    },
    displayLimit: '表示数',
    autoShowDelay: 'ポップアップ遅延',
    saveButton: '設定保存',
    messages: {
      fetchError: '設定読み込み失敗：',
      limitRange: '表示数は1から20の間である必要があります',
      saveSuccess: '設定保存済み',
      saveError: '保存失敗：',
    },
    units: {
      second: '秒',
    },
  },
}
