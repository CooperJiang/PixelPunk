/**
 * Automation tasks page texts
 */
export const automation = {
  page: {
    title: '自動化タスク',
    subtitle: 'AIタグ付けとベクトル化タスクの進捗をリアルタイムで監視',
    refresh: '更新',
    tip: 'キューは手動介入なしでアップロードした画像を自動的に処理します。システムはすべてのタスクを順次完了します。',
  },
  tagging: {
    title: 'AIタグ付けタスク',
    subtitle: '画像コンテンツを自動識別してタグを生成',
    statusRunning: '実行中',
    statusPaused: '一時停止',
    progress: {
      label: '全体の進捗',
      completed: '{done} / {total} 完了',
    },
    queue: {
      position: 'あなたの画像はキュー内の{position}番目の位置にあります',
      pending: '{count}件の画像がキューで待機中',
    },
    status: {
      none: '未処理',
      pending: 'キュー待ち',
      processing: '処理中',
      done: '完了',
      failed: '失敗',
      ignored: '無視',
    },
    empty: {
      title: 'AIタグ付けタスクなし',
      subtitle: '画像をアップロードすると、AI分析が自動的に開始されます',
    },
  },
  vector: {
    title: 'ベクトル化タスク',
    subtitle: '類似検索のために画像をベクトルに変換',
    statusRunning: '実行中',
    statusPaused: '一時停止',
    progress: {
      label: '全体の進捗',
      completed: '{done} / {total} 完了',
    },
    queue: {
      position: 'あなたのベクトルタスクはキュー内の{position}番目の位置にあります',
      pending: '{count}件のベクトルタスクがキューで待機中',
    },
    status: {
      pending: '待機中',
      reset: 'リセット',
      processing: '処理中',
      completed: '完了',
      failed: '失敗',
    },
    empty: {
      title: 'ベクトル化タスクなし',
      subtitle: 'AIタグ付けが完了すると、ベクトル化が自動的に開始されます',
    },
  },
  taggingHistory: {
    filter: {
      statusPlaceholder: 'ステータスを選択',
    },
  },
  history: {
    title: 'タグ付け履歴',
    recordCount: '（合計{count}件の記録）',
    refresh: '更新',
    empty: 'タグ付け記録なし',
    columns: {
      preview: 'プレビュー',
      filename: 'ファイル名',
      status: 'ステータス',
      format: 'フォーマット',
      size: 'サイズ',
      resolution: '解像度',
      createdAt: '作成日時',
      updatedAt: '更新日時',
      error: 'エラーメッセージ',
    },
  },
  toast: {
    loadHistoryError: 'タグ付け記録の読み込みに失敗しました',
    loadDataError: 'タスクデータの読み込みに失敗しました',
  },
}
