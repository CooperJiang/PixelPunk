/**
 * Upload Page Text - Cyber Style
 */
export const upload = {
  title: 'データ注入',

  settings: {
    title: '注入設定',
    advanced: '高度な設定',
    folder: 'データフォルダ',
    folderHint: 'マルチレベルディレクトリをサポート、フォルダ名で検索可能',
    guestFolderInfo: 'ゲスト注入データユニットはシステムフォルダに保存され、デフォルト{level}アクセス',
    guestLimitInfo: 'ゲストモードは単一データユニット注入のみサポート（同時実行制限1）',
    guestDisabled: 'ゲスト注入機能が無効です。管理者に連絡してください',
    accessLevelLabel: 'アクセス権限設定',
    accessLevelPlaceholder: 'アクセス権限を選択',
    accessLevel: {
      public: '🌍 公開',
      publicDesc: '推奨および表示可能、リンク経由で誰でもアクセス可能',
      private: '🔗 プライベート',
      privateDesc: '推奨されませんが、共有および外部参照は通常通り可能',
      protected: '🔒 保護',
      protectedDesc: 'ログイン後にのみ表示可能、外部リンクはアクセス不可',
    },
    storageDurationLabel: '保存期間',
    storageDurationPlaceholder: '保存期間を選択してください',
    storageDurationHint: 'ゲスト注入は保存期間を選択する必要があります。データユニットは期限後に自動削除されます',
    required: '（必須）',
    autoOptimize: 'データユニットを自動最適化（推奨）',
    autoOptimizeHint: 'サムネイルを自動圧縮および最適化し、品質を大幅に低下させることなくデータユニットサイズを削減',
    autoRemove: '成功後に自動削除',
    autoRemoveHint:
      '正常に注入されたデータユニットはリストから自動的に削除され、一括注入中にリストをクリーンに保つのに適しています',
    duration: {
      permanent: '永続',
    },
  },

  uploadProgress: {
    status: {
      uploadSuccess: '注入成功',
      instantSuccess: '即時注入成功',
      uploadFailed: '注入失敗',
      analyzing: 'データユニットを分析中...',
      paused: '未完了の注入が見つかりました。続行可能',
      retrying: '注入を再試行中...',
      pending: '注入待機中',
    },
    time: {
      seconds: '{count}秒',
      minutes: '{count}分',
      hoursMinutes: '{hours}時間{minutes}分',
    },
    remaining: '残り',
    actions: {
      openInNewTab: '新しいタブで開く',
      copyLink: 'データユニットリンクをコピー',
      copyLinkTooltip: '元のリンク',
      copyMarkdownLink: 'Markdown リンクをコピー',
      copyMarkdownTooltip: 'MD形式',
      copyHtmlLink: 'HTML タグをコピー',
      copyHtmlTooltip: 'HTML形式',
      copyThumbnailLink: 'サムネイルリンクをコピー',
      copyThumbnailTooltip: 'サムネイル',
      resume: '注入を続行',
      retry: '注入を再試行',
      remove: 'データユニットを削除',
    },
    toast: {
      linkCopied: 'リンクをコピーしました',
      markdownLinkCopied: 'Markdown リンクをコピーしました',
      htmlLinkCopied: 'HTML タグをコピーしました',
      thumbnailLinkCopied: 'サムネイルリンクをコピーしました',
      copyFailed: 'コピー失敗',
      uploadSuccessWithRemainingCount: 'データ注入成功！{count}回の注入機会が残っています',
      watermarkApplyFailed: 'データユニットを注入しましたが、ウォーターマークの適用に失敗しました: {reason}',
      thumbnailGenerationFailed: 'データユニットを注入しましたが、サムネイル生成に失敗しました: {reason}',
      fileTypeNotSupported: 'データタイプがサポートされていません: {name}',
      fileSizeExceeded: 'データサイズ超過: {name}（{size}）',
      instantUploadSuccess: '即時アップロード成功: {name}',
      cacheCleared: 'キャッシュをクリアしました',
    },
  },

  dropZone: {
    title: 'データユニットをここにドラッグして注入',
    orText: 'または',
    clickToSelect: 'クリックしてデータを選択',
    supportFolder: 'データディレクトリ全体のドラッグをサポート',
    features: {
      batch: '一括注入',
      autoOptimize: '自動最適化',
      secureStorage: 'セキュアストレージ',
    },
    shortcut: {
      label: 'ショートカット',
      paste: '貼り付けて注入',
    },
  },

  actions: {
    startUpload: '注入を開始',
    selectFiles: 'データを選択',
    selectFolder: 'データディレクトリを選択',
    copyLinks: 'リンクをコピー',
    copyMarkdownLinks: '全ての Markdown リンクをコピー',
    copyHtmlLinks: '全ての HTML タグをコピー',
    copyThumbnailLinks: '全てのサムネイルリンクをコピー',
    clearQueue: 'キューをクリア',
    cancelUpload: '注入をキャンセル',
    clearCache: 'キャッシュをクリア',
  },

  queue: {
    title: '注入キュー',
    filesCount: '{count}個のデータユニット',
    concurrent: '{current}/{max}同時実行',
    stats: {
      pending: '待機中',
      uploading: '注入中',
      success: '成功',
      failed: '失敗',
      total: '合計',
      totalSize: '合計サイズ',
    },
    empty: {
      title: 'キューが空です',
      desc: 'データをドラッグ、クリックして選択、またはスクリーンショットを貼り付けて注入を開始',
      action: 'データを選択',
    },
  },

  guide: {
    title: '注入のヒント',
    multipleWays: {
      title: '複数の注入方法',
      desc: 'ドラッグアンドドロップ、クリックして選択、または{key}+Vを使用してクリップボードデータを貼り付けをサポート',
    },
    smartChunking: {
      title: 'スマートチャンキング',
      desc: '大きなデータユニットは転送用に自動的にチャンク化され、再開可能なアップロードをサポート',
    },
    oneCopy: {
      title: 'ワンクリックコピー',
      desc: '一括注入後、ワンクリックですべてのデータリンクをコピーできます',
    },
  },

  advanced: {
    title: '高度な設定',
    storageDuration: '保存期間',
    storageDurationDesc: 'データユニットがサーバーに保持される時間',
    watermark: 'データウォーターマークを追加',
    watermarkDesc: '注入されたデータにウォーターマーク識別子を追加',
    watermarkInfo: {
      type: 'タイプ:',
      content: 'コンテンツ:',
      position: '位置:',
    },
    watermarkType: {
      text: 'テキスト',
      image: '画像',
    },
    loginPrompt: {
      title: 'より高度な機能',
      desc: 'システムにアクセスして永続ストレージ、カスタムウォーターマーク、その他の機能を使用',
      action: '今すぐアクセス',
    },
  },

  folder: {
    browserNotSupported: 'お使いのブラウザはデータディレクトリ注入をサポートしていません',
    noValidFiles: 'ディレクトリに注入可能なデータユニットが見つかりません',
    parsing: 'データディレクトリを解析中...',
    addSuccess: 'ディレクトリから{count}個のデータユニットを追加しました',
    dialog: {
      title: 'データディレクトリ注入確認',
      totalFiles: 'データユニット総数',
      validFiles: '注入可能ユニット',
      invalidFiles: '非サポートユニット',
      confirmMessage: '<strong class="highlight-text">{count}</strong>個の注入可能なデータユニットが見つかりました。注入キューに追加しますか？',
      warningMessage: '{count}個の非サポートデータユニットをフィルタしました',
      cancel: 'キャンセル',
      confirm: '注入を確認',
    },
  },

  toast: {
    folderSelected: 'データディレクトリを選択しました: {path}',
    folderPreselected: 'データディレクトリが事前選択されました。データ注入を直接開始できます',
    sensitiveAutoDelete: '違反データは自動的にクリアされます。データユニットがプロトコルに準拠していることを確認してください',
    sensitivePendingReview: 'センシティブデータはレビューキューに入り、レビュー中は一時的にアクセス不可',
    defaultHint: '注入されたデータがプラットフォームプロトコルに準拠していることを確認してください',
    guestDurationRequired: 'ゲスト注入は保存期間を選択する必要があります',
    uploadCancelled: 'データ注入がキャンセルされました',
    noFilesInQueue: 'キューに注入するデータユニットがありません',
    allLinksCopied: 'すべてのリンクをクリップボードにコピーしました',
    markdownLinksCopied: '全ての Markdown リンクをコピーしました',
    htmlLinksCopied: '全ての HTML タグをコピーしました',
    thumbnailLinksCopied: '全てのサムネイルリンクをコピーしました',
    noLinksToCopy: 'まだコピーするデータリンクがありません',
  },

  duration: {
    permanent: '永続',
  },

  watermark: {
    position: {
      'top-left': '左上',
      'top-center': '上中央',
      'top-right': '右上',
      'middle-left': '中央左',
      'middle-center': '中央',
      'middle-right': '中央右',
      'bottom-left': '左下',
      'bottom-center': '下中央',
      'bottom-right': '右下',
      custom: 'カスタム位置',
    },
    configure: 'データウォーターマークを設定',
  },

  form: {
    storageDuration: {
      placeholder: '保存期間を選択',
    },
  },

  smartUpload: {
    error: {
      fileTypeNotSupported: 'データタイプがサポートされていません: {name}（{type}）',
      fileSizeExceeded: 'データユニットが制限を超えました: {name}（{size}）',
      chunkedUploadNotReady: 'チャンク転送が初期化されていません。後でもう一度お試しください',
      uploadFailed: '転送失敗',
      analysisFailed: 'データ分析失敗',
      analysisTimeout: 'データ分析タイムアウト',
    },
    status: {
      preparing: '注入を準備中',
      uploading: '注入中...',
      uploadingWithProgress: '注入中{progress}% - {speed} - 残り{remaining}',
      completed: '注入完了',
      failed: '注入失敗: {error}',
      paused: '一時停止',
      cancelled: 'キャンセル',
      retrying: '再試行を準備中',
      userPaused: 'ユーザー一時停止',
      userCancelled: 'ユーザーキャンセル',
      analyzing: 'データを分析中...',
      analyzingProgress: 'データを分析中{progress}%',
      checkingInstant: '即時アップロードを確認中...',
      checkingInstantProgress: '即時アップロードを確認中{progress}%',
      instantCompleted: '即時アップロード完了',
      foundResumable: '未完了の転送が見つかりました。続行可能',
      analyzingCompleted: '✓ 分析完了、注入準備完了',
      analyzingFailed: '分析失敗: {error}',
      retryingUpload: '注入を再試行中...',
      resumingUpload: '注入を再開中...',
      continuingUpload: '注入を続行中...',
      restartingUpload: '注入を再起動中...',
      preparingSession: '転送セッションを準備中...',
      verifyingIntegrity: 'データ整合性を確認中...',
      completingUpload: '注入を完了中...',
      chunkUploadFailed: 'チャンク{chunk}転送失敗: {error}',
      chunkRetrying: 'チャンク{chunk}再試行中...（{attempt}/{maxRetry}）- {delay}秒後に再試行',
    },
    time: {
      minutes: '{minutes}分{seconds}秒',
      seconds: '{seconds}秒',
    },
  },
}
