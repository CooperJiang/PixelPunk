/**
 * Automation Tasks Page Text - Cyber Style
 */
export const automation = {
  page: {
    title: '自動化コマンド',
    subtitle: 'AIタグ付けとベクトル化タスクの実行ステータスをリアルタイムで監視',
    refresh: '更新',
    tip: 'キューは手動介入なしで注入されたデータユニットを自動的に処理します。システムはすべてのコマンドを順次実行します。',
  },
  tagging: {
    title: 'AIタグ付けタスク',
    subtitle: 'データユニットコンテンツを自動識別してタグを生成',
    statusRunning: '実行中',
    statusPaused: '一時停止',
    progress: {
      label: '実行進捗',
      completed: '{done} / {total} 完了',
    },
    queue: {
      position: 'あなたのデータユニットはキュー内の{position}番目の位置にあります',
      pending: '{count}個のデータユニットがキューで待機中',
    },
    status: {
      none: '未処理',
      pending: 'キュー待機中',
      processing: '実行中',
      done: '完了',
      failed: '失敗',
      ignored: '無視',
    },
    empty: {
      title: 'AIタグ付けタスクなし',
      subtitle: 'データユニットを注入すると、AI分析が自動的に開始されます',
    },
  },
  vector: {
    title: 'ベクトル化タスク',
    subtitle: 'データユニットをベクトルに変換して類似検索を実行',
    statusRunning: '実行中',
    statusPaused: '一時停止',
    progress: {
      label: '実行進捗',
      completed: '{done} / {total} 完了',
    },
    queue: {
      position: 'あなたのベクトルタスクはキュー内の{position}番目の位置にあります',
      pending: '{count}個のベクトルタスクがキューで待機中',
    },
    status: {
      pending: '待機中',
      reset: 'リセット',
      processing: '実行中',
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
