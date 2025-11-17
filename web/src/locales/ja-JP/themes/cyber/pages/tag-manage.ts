/**
 * Tag Management Page Text - Cyber Style
 */
export const tagManage = {
  title: 'タグ管理',
  subtitle: 'タグを管理し、データユニットを整理・並べ替え',
  search: 'タグを検索...',
  create: 'タグを作成',
  batchMode: '一括モード',
  exitBatchMode: '一括モードを終了',
  loading: 'タグを読み込み中...',
  batch: {
    viewFiles: '選択したタグのファイルを表示',
    delete: '一括削除（{count}）',
    merge: 'タグをマージ',
    cancel: '選択をクリア',
    modeHint: '一括操作モード',
    modeDescription: 'タグをクリックして選択/選択解除、{count}個のタグを選択',
  },
  multiSelect: {
    modeHint: '複数選択モード',
    modeDescription: '{modifier} / Altキーを押しながらタグをクリックして複数選択、{count}個のタグを選択',
  },
  normalMode: {
    hint: 'ヒント: タグをクリックして迅速に検索、{modifier} / Altキーを押しながら複数選択',
  },
  actions: {
    edit: 'タグを編集',
    delete: 'タグを削除',
  },
  dialog: {
    create: {
      title: 'タグを作成',
      submit: '作成',
    },
    edit: {
      title: 'タグを編集',
      submit: '更新',
    },
    merge: {
      title: 'タグをマージ',
      targetLabel: 'ターゲットタグ',
      targetPlaceholder: 'ターゲットタグを選択してください',
    },
    delete: {
      title: 'タグを削除',
      message:
        'タグ「{name}」を削除してもよろしいですか？このタグは{count}回使用されており、削除するとすべてのデータユニットからこのタグが削除されます。',
      confirmText: '削除を確認',
    },
    batchDelete: {
      title: 'タグを一括削除',
      message: '選択した{count}個のタグを削除してもよろしいですか？この操作は元に戻せません。',
      confirmText: '削除を確認',
    },
    mergeConfirm: {
      title: 'タグをマージ',
      message: [
        '{count}個のタグを「{targetName}」にマージしてもよろしいですか？',
        'ソースタグは削除され、それらのファイル関連付けはターゲットタグに転送されます。',
        'この操作は元に戻せません。',
      ],
      confirmText: 'マージを確認',
    },
  },
  form: {
    name: {
      label: 'タグ名',
      placeholder: 'タグ名を入力してください',
      hint: '最大50文字、短く覚えやすい名前の使用を推奨',
    },
    cancel: 'キャンセル',
  },
  merge: {
    description:
      '選択したタグをターゲットタグにマージします。ソースタグは削除され、それらのファイル関連付けはターゲットタグに転送されます。',
    selectedInfo:
      '{count}個のタグを選択しました。ターゲットタグを選択すると、他の{remaining}個のタグがターゲットタグにマージされます。',
    confirm: 'マージを確認',
    cancel: 'キャンセル',
  },
  empty: {
    title: '一致するタグが見つかりません',
    description: 'フィルター条件を調整するか、新しいタグを作成してください',
  },
  tagCloud: {
    groupCount: '{count}個のタグ',
    groupUsage: '{usage}回使用',
    letterNav: '{letter}（{count}）',
    letterEmpty: '{letter}（なし）',
    usageCount: '{count}回使用',
  },
  toast: {
    nameRequired: 'タグ名を入力してください',
    updateSuccess: 'タグを更新しました',
    createSuccess: 'タグを作成しました',
    deleteSuccess: 'タグを削除しました',
    selectRequired: 'まず削除するタグを選択してください',
    mergeMinRequired: 'マージするには少なくとも2つのタグを選択する必要があります',
    mergeTargetRequired: 'ターゲットタグを選択してください',
    mergeTargetInvalid: 'ターゲットタグは選択したタグの中から選択する必要があります',
    mergeSuccess: 'タグをマージしました',
    batchSelectRequired: 'まずタグを選択してください',
  },
  source: {
    manual: '手動作成',
    ai: 'AI認識',
    system: 'システムインポート',
    unknown: '不明',
  },
}
