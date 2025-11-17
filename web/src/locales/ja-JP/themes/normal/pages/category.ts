/**
 * Category management page texts
 */
export const category = {
  page: {
    title: 'カテゴリ管理',
    subtitle: 'ファイルカテゴリを管理し、リソースを整理してソート',
    create: '新しいカテゴリ',
    refresh: '更新',
    filter: 'フィルター',
    hideFilter: 'フィルターを非表示',
  },
  table: {
    columns: {
      name: 'カテゴリ名',
      description: '説明',
      source: 'ソース',
      fileCount: 'ファイル数',
      sortOrder: 'ソート順',
      status: 'ステータス',
      createdAt: '作成日時',
      actions: 'アクション',
    },
    noDescription: '説明なし',
    clickToView: 'クリックしてこのカテゴリのすべてのファイルを表示',
    loading: 'カテゴリデータを読み込み中...',
  },
  source: {
    system: 'システム',
    user: 'ユーザー',
    imported: 'インポート',
  },
  status: {
    active: 'アクティブ',
    archived: 'アーカイブ済み',
  },
  actions: {
    edit: 'カテゴリを編集',
    archive: 'カテゴリをアーカイブ',
    activate: 'カテゴリを有効化',
    delete: 'カテゴリを削除',
  },
  empty: {
    title: 'カテゴリデータなし',
    description: 'まだカテゴリが作成されていません。「新しいカテゴリ」をクリックして最初のカテゴリを作成してください',
    action: '新しいカテゴリ',
  },
  dialog: {
    edit: {
      title: 'カテゴリを編集',
      cancel: 'キャンセル',
      save: '変更を保存',
    },
    delete: {
      title: 'カテゴリを削除',
      message: 'このカテゴリを削除してもよろしいですか？',
    },
    create: {
      title: '新しいカテゴリ',
      cancel: 'キャンセル',
      create: 'カテゴリを作成',
    },
  },
  form: {
    name: {
      label: 'カテゴリ名',
      placeholder: 'カテゴリ名を入力してください',
      required: 'カテゴリ名を入力してください',
    },
    description: {
      label: 'カテゴリ説明',
      placeholder: 'カテゴリ説明を入力してください（オプション）',
    },
    sortOrder: {
      label: 'ソート順',
      placeholder: 'ソート順を入力してください',
      hint: '小さい値がリストの最初に表示されます',
    },
    status: {
      label: 'ステータス',
      active: 'アクティブ',
      archived: 'アーカイブ済み',
    },
  },
  filter: {
    title: 'フィルター条件',
    sortBy: 'ソート',
    status: {
      label: 'ステータス',
      all: 'すべてのステータス',
      active: 'アクティブ',
      archived: 'アーカイブ済み',
    },
    source: {
      label: 'ソース',
      all: 'すべて',
      system: 'システム',
      user: 'ユーザー',
      imported: 'インポート',
    },
    search: {
      placeholder: 'カテゴリ名または説明を検索...',
    },
    apply: '適用',
    reset: 'リセット',
  },
  batch: {
    selected: '{count}項目を選択',
    clearSelection: '選択をクリア',
    clear: '選択をクリア',
    delete: '一括削除',
    archive: '一括アーカイブ',
    activate: '一括有効化',
    deleteConfirm: '選択した{count}件のカテゴリを削除してもよろしいですか？',
    archiveConfirm: '選択した{count}件のカテゴリをアーカイブしてもよろしいですか？',
    activateConfirm: '選択した{count}件のカテゴリを有効化してもよろしいですか？',
  },
  toast: {
    createSuccess: 'カテゴリが正常に作成されました',
    createError: 'カテゴリの作成に失敗しました',
    createWarning: 'カテゴリ名を入力してください',
    updateSuccess: 'カテゴリが正常に更新されました',
    updateError: 'カテゴリの更新に失敗しました',
    deleteSuccess: 'カテゴリが正常に削除されました',
    deleteError: 'カテゴリの削除に失敗しました',
    archiveSuccess: 'アーカイブに成功しました',
    activateSuccess: '有効化に成功しました',
    archiveError: 'アーカイブに失敗しました',
    activateError: '有効化に失敗しました',
    batchDeleteSuccess: '{count}件のカテゴリを正常に削除しました',
    batchDeleteError: '一括削除に失敗しました',
    batchArchiveSuccess: '一括アーカイブに成功しました',
    batchActivateSuccess: '一括有効化に成功しました',
    fetchError: 'カテゴリリストの取得に失敗しました',
    selectWarning: 'まず削除するカテゴリを選択してください',
  },
  sourceInfo: {
    manual: '手動作成',
    aiSuggestion: 'AI提案',
    system: 'システムテンプレート',
    imported: '一括インポート',
  },
  statusInfo: {
    active: 'アクティブ',
    archived: 'アーカイブ済み',
  },
  placeholders: {
    search: 'カテゴリ名または説明を検索...',
    status: 'ステータス',
    name: 'カテゴリ名を入力してください',
    description: 'カテゴリ説明を入力してください（オプション）',
    sortOrder: '小さい値が最初に表示されます',
  },
}
