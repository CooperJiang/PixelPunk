/**
 * Category Management Page Text - Cyber Style
 */
export const category = {
  page: {
    title: 'カテゴリ管理',
    subtitle: 'データカテゴリを管理し、リソースノードを整理・並べ替え',
    create: 'カテゴリを作成',
    refresh: '更新',
    filter: 'フィルター',
    hideFilter: 'フィルターを非表示',
  },
  table: {
    columns: {
      name: 'カテゴリ名',
      description: '説明',
      source: 'ソース',
      fileCount: 'データユニット数',
      sortOrder: '並び順',
      status: 'ステータス',
      createdAt: '作成日時',
      actions: 'アクション',
    },
    noDescription: '説明なし',
    clickToView: 'クリックしてこのカテゴリのすべてのデータユニットを表示',
    loading: 'カテゴリデータを読み込み中...',
  },
  source: {
    system: 'システム',
    user: 'ユーザー',
    imported: 'インポート済み',
  },
  status: {
    active: '有効',
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
    description: 'まだカテゴリが作成されていません。「カテゴリを作成」をクリックして最初のカテゴリを作成',
    action: 'カテゴリを作成',
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
      title: 'カテゴリを作成',
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
      label: '並び順値',
      placeholder: '並び順値を入力してください',
      hint: '並び順値が小さいほどリストの先頭に表示されます',
    },
    status: {
      label: 'ステータス',
      active: '有効',
      archived: 'アーカイブ済み',
    },
  },
  filter: {
    title: 'フィルター条件',
    sortBy: '並び順方法',
    status: {
      label: 'ステータス',
      all: 'すべてのステータス',
      active: '有効化',
      archived: 'アーカイブ済み',
    },
    source: {
      label: 'ソース',
      all: 'すべて',
      system: 'システム',
      user: 'ユーザー',
      imported: 'インポート済み',
    },
    search: {
      placeholder: 'カテゴリ名または説明を検索...',
    },
    apply: '適用',
    reset: 'リセット',
  },
  batch: {
    selected: '{count}件を選択',
    clearSelection: '選択をクリア',
    clear: '選択をクリア',
    delete: '一括削除',
    archive: '一括アーカイブ',
    activate: '一括有効化',
    deleteConfirm: '選択した{count}個のカテゴリを削除してもよろしいですか？',
    archiveConfirm: '選択した{count}個のカテゴリをアーカイブしてもよろしいですか？',
    activateConfirm: '選択した{count}個のカテゴリを有効化してもよろしいですか？',
  },
  toast: {
    createSuccess: 'カテゴリを作成しました',
    createError: 'カテゴリの作成に失敗しました',
    createWarning: 'カテゴリ名を入力してください',
    updateSuccess: 'カテゴリを更新しました',
    updateError: 'カテゴリの更新に失敗しました',
    deleteSuccess: 'カテゴリを削除しました',
    deleteError: 'カテゴリの削除に失敗しました',
    archiveSuccess: 'アーカイブしました',
    activateSuccess: '有効化しました',
    archiveError: 'アーカイブに失敗しました',
    activateError: '有効化に失敗しました',
    batchDeleteSuccess: '{count}個のカテゴリを削除しました',
    batchDeleteError: '一括削除に失敗しました',
    batchArchiveSuccess: '一括アーカイブ成功',
    batchActivateSuccess: '一括有効化成功',
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
    active: '有効化',
    archived: 'アーカイブ済み',
  },
  placeholders: {
    search: 'カテゴリ名または説明を検索...',
    status: 'ステータス',
    name: 'カテゴリ名を入力してください',
    description: 'カテゴリ説明を入力してください（オプション）',
    sortOrder: '値が小さいほど先頭に表示',
  },
}
