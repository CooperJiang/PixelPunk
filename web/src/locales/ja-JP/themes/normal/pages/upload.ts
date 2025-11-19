/**
 * Upload page texts
 */
export const upload = {
  title: 'ファイルをアップロード',

  settings: {
    title: 'アップロード設定',
    advanced: '詳細設定',
    folder: 'フォルダ',
    folderHint: 'マルチレベルディレクトリをサポート、フォルダ名を検索可能',
    guestFolderInfo: 'ゲストがアップロードしたファイルはシステムフォルダに保存され、デフォルトで{level}アクセス',
    guestLimitInfo: 'ゲストモードは単一ファイルアップロードのみサポート（同時実行制限は1）',
    guestDisabled: 'ゲストアップロード機能は無効です。管理者に連絡してください',
    accessLevelLabel: 'アクセスレベル設定',
    accessLevelPlaceholder: 'アクセスレベルを選択',
    accessLevel: {
      public: '🌍 公開',
      publicDesc: '推奨・表示可能、リンクで誰でもアクセス可能',
      private: '🔗 プライベート',
      privateDesc: '推奨されませんが、シェアと外部参照は通常通り可能',
      protected: '🔒 保護',
      protectedDesc: 'ログイン後のみ表示可能、外部リンクはアクセス不可',
    },
    storageDurationLabel: '保存期間',
    storageDurationPlaceholder: '保存期間を選択してください',
    storageDurationHint: 'ゲストアップロードは保存期間を選択する必要があります。期限切れ後、ファイルは自動的に削除されます',
    required: '（必須）',
    autoOptimize: 'ファイルを自動最適化（推奨）',
    autoOptimizeHint: 'サムネイルを自動的に圧縮・最適化し、品質を大幅に低下させることなくファイルサイズを削減',
    autoRemove: 'アップロード成功後に自動削除',
    autoRemoveHint:
      'アップロードに成功したファイルは自動的にリストから削除され、バッチアップロード中にリストをクリーンに保つのに適しています',
    duration: {
      permanent: '永続',
    },
  },

  uploadProgress: {
    status: {
      uploadSuccess: 'アップロード成功',
      instantSuccess: '即時アップロード成功',
      uploadFailed: 'アップロード失敗',
      analyzing: 'ファイルを分析中...',
      paused: '不完全なアップロードが見つかりました。続行可能',
      retrying: '再試行中...',
      pending: 'アップロード待機中',
    },
    time: {
      seconds: '{count}秒',
      minutes: '{count}分',
      hoursMinutes: '{hours}時間{minutes}分',
    },
    remaining: '残り',
    actions: {
      openInNewTab: '新しいタブで開く',
      copyLink: 'ファイルリンクをコピー',
      copyLinkTooltip: '元のリンク',
      copyMarkdownLink: 'Markdown リンクをコピー',
      copyMarkdownTooltip: 'MD形式',
      copyHtmlLink: 'HTML タグをコピー',
      copyHtmlTooltip: 'HTML形式',
      copyThumbnailLink: 'サムネイルリンクをコピー',
      copyThumbnailTooltip: 'サムネイル',
      resume: 'アップロードを再開',
      retry: 'アップロードを再試行',
      remove: 'ファイルを削除',
    },
    toast: {
      linkCopied: 'リンクをコピーしました',
      markdownLinkCopied: 'Markdown リンクをコピーしました',
      htmlLinkCopied: 'HTML タグをコピーしました',
      thumbnailLinkCopied: 'サムネイルリンクをコピーしました',
      copyFailed: 'コピーに失敗しました',
      uploadSuccessWithRemainingCount: 'アップロード成功！残り{count}回のアップロード機会',
      watermarkApplyFailed: 'ファイルはアップロードされましたが、透かしの適用に失敗しました: {reason}',
      thumbnailGenerationFailed: 'ファイルはアップロードされましたが、サムネイルの生成に失敗しました: {reason}',
      fileTypeNotSupported: 'サポートされていないファイルタイプ: {name}',
      fileSizeExceeded: 'ファイルサイズ超過: {name} ({size})',
      instantUploadSuccess: '即時アップロード成功: {name}',
      cacheCleared: 'キャッシュをクリアしました',
    },
  },

  dropZone: {
    title: 'ファイルをここにドラッグしてアップロード',
    orText: 'または',
    clickToSelect: 'クリックしてファイルを選択',
    supportFolder: 'フォルダ全体のドラッグをサポート',
    features: {
      batch: 'バッチアップロード',
      autoOptimize: '自動最適化',
      secureStorage: 'セキュアストレージ',
    },
    shortcut: {
      label: 'ショートカット',
      paste: '貼り付けアップロード',
    },
  },

  actions: {
    startUpload: 'アップロード開始',
    selectFiles: 'ファイルを選択',
    selectFolder: 'フォルダを選択',
    copyLinks: 'リンクをコピー',
    copyMarkdownLinks: '全ての Markdown リンクをコピー',
    copyHtmlLinks: '全ての HTML タグをコピー',
    copyThumbnailLinks: '全てのサムネイルリンクをコピー',
    clearQueue: 'キューをクリア',
    cancelUpload: 'アップロードをキャンセル',
    clearCache: 'キャッシュをクリア',
  },

  queue: {
    title: 'アップロードキュー',
    filesCount: '{count}ファイル',
    concurrent: '{current}/{max}同時実行',
    stats: {
      pending: '待機中',
      uploading: 'アップロード中',
      success: '成功',
      failed: '失敗',
      total: '合計',
      totalSize: '合計サイズ',
    },
    empty: {
      title: 'キューは空です',
      desc: 'ファイルをドラッグ、クリックして選択、またはスクリーンショットを貼り付けてアップロードを開始',
      action: 'ファイルを選択',
    },
  },

  guide: {
    title: 'アップロードのヒント',
    multipleWays: {
      title: '複数のアップロード方法',
      desc: 'ドラッグ&ドロップ、クリックして選択、または{key}+Vでクリップボード画像を貼り付けをサポート',
    },
    smartChunking: {
      title: 'スマートチャンキング',
      desc: '大きなファイルは自動的にチャンクに分割してアップロード、再開可能なアップロードをサポート',
    },
    oneCopy: {
      title: 'ワンクリックコピー',
      desc: 'バッチアップロード後、ワンクリックですべてのファイルリンクをコピー',
    },
  },

  advanced: {
    title: '詳細設定',
    storageDuration: '保存期間',
    storageDurationDesc: 'サーバー上のファイル保存時間',
    watermark: '透かしを追加',
    watermarkDesc: 'アップロードした画像に透かしを追加',
    watermarkInfo: {
      type: 'タイプ:',
      content: '内容:',
      position: '位置:',
    },
    watermarkType: {
      text: 'テキスト',
      image: '画像',
    },
    loginPrompt: {
      title: 'より高度な機能',
      desc: 'ログイン後、永続ストレージ、カスタム透かしなどの機能を使用できます',
      action: '今すぐログイン',
    },
  },

  folder: {
    browserNotSupported: 'お使いのブラウザはフォルダアップロードをサポートしていません',
    noValidFiles: 'フォルダにアップロード可能なファイルが見つかりません',
    parsing: 'フォルダを解析中...',
    addSuccess: 'フォルダから{count}個のファイルを追加しました',
    dialog: {
      title: 'フォルダアップロード確認',
      totalFiles: 'ファイル総数',
      validFiles: 'アップロード可能ファイル',
      invalidFiles: '非サポートファイル',
      confirmMessage: '<strong class="highlight-text">{count}</strong>個のアップロード可能なファイルが見つかりました。アップロードキューに追加しますか？',
      warningMessage: '{count}個の非サポートファイルをフィルタしました',
      cancel: 'キャンセル',
      confirm: 'アップロードを確認',
    },
  },

  toast: {
    folderSelected: 'フォルダを選択しました: {path}',
    folderPreselected: 'フォルダが事前選択されました。直接アップロードを開始できます',
    sensitiveAutoDelete: '不適切なコンテンツは自動的に削除されます。ファイルが規制に準拠していることを確認してください',
    sensitivePendingReview: 'センシティブなコンテンツはレビューキューに入り、レビュー中は一時的にアクセスできません',
    defaultHint: 'アップロードしたコンテンツがプラットフォーム規制に準拠していることを確認してください',
    guestDurationRequired: 'ゲストアップロードは保存期間を選択する必要があります',
    uploadCancelled: 'アップロードがキャンセルされました',
    noFilesInQueue: 'アップロードするファイルがキューにありません',
    allLinksCopied: 'すべてのリンクをクリップボードにコピーしました',
    markdownLinksCopied: '全ての Markdown リンクをコピーしました',
    htmlLinksCopied: '全ての HTML タグをコピーしました',
    thumbnailLinksCopied: '全てのサムネイルリンクをコピーしました',
    noLinksToCopy: 'コピーするファイルリンクがありません',
  },

  duration: {
    permanent: '永続',
  },

  watermark: {
    position: {
      'top-left': '左上',
      'top-center': '上中央',
      'top-right': '右上',
      'middle-left': '左中央',
      'middle-center': '中央',
      'middle-right': '右中央',
      'bottom-left': '左下',
      'bottom-center': '下中央',
      'bottom-right': '右下',
      custom: 'カスタム位置',
    },
    configure: '透かしを設定',
  },

  form: {
    storageDuration: {
      placeholder: '保存期間を選択',
    },
  },

  smartUpload: {
    error: {
      fileTypeNotSupported: 'サポートされていないファイルタイプ: {name} ({type})',
      fileSizeExceeded: 'ファイルサイズ超過: {name} ({size})',
      chunkedUploadNotReady: 'チャンクアップロードが初期化されていません。後でもう一度お試しください',
      uploadFailed: 'アップロードに失敗しました',
      analysisFailed: 'ファイル分析に失敗しました',
      analysisTimeout: 'ファイル分析タイムアウト',
    },
    status: {
      preparing: 'アップロード準備中',
      uploading: 'アップロード中...',
      uploadingWithProgress: 'アップロード中 {progress}% - {speed} - 残り{remaining}',
      completed: 'アップロード完了',
      failed: 'アップロード失敗: {error}',
      paused: '一時停止',
      cancelled: 'キャンセル済み',
      retrying: '再試行準備中',
      userPaused: 'ユーザーが一時停止',
      userCancelled: 'ユーザーがキャンセル',
      analyzing: 'ファイルを分析中...',
      analyzingProgress: '分析中 {progress}%',
      checkingInstant: '即時アップロードを確認中...',
      checkingInstantProgress: '即時アップロード確認中 {progress}%',
      instantCompleted: '即時アップロード完了',
      foundResumable: '不完全なアップロードが見つかりました。続行可能',
      analyzingCompleted: '✓ 分析完了、アップロード準備完了',
      analyzingFailed: '分析失敗: {error}',
      retryingUpload: 'アップロードを再試行中...',
      resumingUpload: 'アップロードを再開中...',
      continuingUpload: 'アップロードを続行中...',
      restartingUpload: 'アップロードを再起動中...',
      preparingSession: 'アップロードセッションを準備中...',
      verifyingIntegrity: 'アップロードの整合性を確認中...',
      completingUpload: 'アップロードを完了中...',
      chunkUploadFailed: 'チャンク{chunk}のアップロードに失敗しました: {error}',
      chunkRetrying: 'チャンク{chunk}を再試行中... ({attempt}/{maxRetry}) - {delay}秒後に再試行',
    },
    time: {
      minutes: '{minutes}分{seconds}秒',
      seconds: '{seconds}秒',
    },
  },
}
