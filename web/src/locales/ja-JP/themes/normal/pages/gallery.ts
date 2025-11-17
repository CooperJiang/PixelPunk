/**
 * Gallery browsing texts
 */
export const gallery = {
  pages: {
    gallery: 'ギャラリー',
    upload: 'ファイルをアップロード',
    login: 'ログイン',
    register: '登録',
    settings: '設定',
    profile: 'プロフィール',
    folders: 'フォルダ',
    dashboard: 'ダッシュボード',
    admin: '管理者パネル',
  },
  viewer: {
    previous: '前へ',
    next: '次へ',
    fitMode: 'フィットモード',
    fillMode: 'フィルモード',
    fullscreen: '全画面',
    exitFullscreen: '全画面を終了',
    originalSize: 'オリジナルサイズ',
    imageInfo: 'ファイル情報',
    similarImages: '類似ファイル',
  },
  waterfallLayout: {
    loadFailed: '読み込み失敗',
    collapse: '折りたたむ',
  },
  file: {
    loadFailed: '読み込み失敗',
    nsfwWarning: '不適切なコンテンツ',
  },
  enhancedFilePreview: {
    defaultImageName: 'ファイル',
    fillMode: 'フィルモード',
    fitMode: 'フィットモード',
    switchToFit: 'フィットに切り替え',
    switchToFill: 'フィルに切り替え',
    wheelZoom: 'ホイールズーム',
    dragMove: 'ドラッグ移動',
    spaceKey: 'スペース',
    exitPreview: 'プレビューを終了',
    exitFullscreen: '全画面を終了',
    exitPreviewEsc: 'プレビューを終了（ESC）',
    exitFullscreenEsc: '全画面を終了（ESC）',
    enterFullscreen: 'ブラウザ全画面に入る',
    fullscreen: '全画面',
    preview: 'プレビュー',
    mode: 'モード',
  },
  fileDetailModal: {
    title: 'ファイル詳細',
    untitled: '無題',

    basicInfo: {
      title: '基本情報',
      fileId: 'ファイルID',
      originalName: 'オリジナル名',
      fileFormat: 'ファイルフォーマット',
      fileSize: 'ファイルサイズ',
      imageDimensions: 'ファイルサイズ',
      visibility: '可視性',
      uploader: 'アップロード者',
      uploadTime: 'アップロード日時',
      updateTime: '更新日時',
      fileHash: 'ファイルハッシュ',
      viewCount: '閲覧数',
      status: 'ステータス',
      times: '回',
      count: '項目',
      unknownResolution: '不明な解像度',
      none: 'なし',
    },

    statusLabels: {
      duplicate: '重複',
      recommended: '推奨',
      nsfw: 'NSFW',
    },

    visibility: {
      public: '公開',
      private: 'プライベート',
      protected: '保護',
      link: 'リンク可視',
    },

    qualityLevels: {
      high: '高',
      medium: '中',
      low: '低',
      excellent: '優秀',
      normal: '通常',
      poor: '低い',
    },

    contentSafety: {
      title: 'コンテンツセーフティ',
      nsfwStatus: 'NSFWステータス',
      inappropriate: '不適切',
      safe: '安全',
      nsfwScore: 'NSFWスコア',
      high: '高',
      medium: '中',
      low: '低',
    },

    visualAnalysis: {
      title: '視覚分析',
      dominantColor: '主な色',
      imageDimensions: 'ファイルサイズ',
      aspectRatio: 'アスペクト比',
      resolution: '解像度',
      colorPalette: 'カラーパレット',
      composition: '構成',
      objectsCount: 'オブジェクト数',
    },

    ai: {
      tags: 'AIタグ',
      description: 'AI説明',
    },

    actions: {
      deleteImage: 'ファイルを削除',
      cancelRecommend: '推奨をキャンセル',
      setRecommend: '推奨に設定',
      close: '閉じる',
    },
  },
  fileExpiryTag: {
    statuses: {
      expired: '期限切れ',
      expiringSoon: 'まもなく期限切れ',
      active: '時間制限',
    },
    labels: {
      storageDuration: '保存期間',
      expiresAt: '有効期限',
      expiredAtPrefix: '期限切れ日',
      expiredAtSuffix: '期限切れ',
    },
    units: {
      day: '日',
      hour: '時間',
      minute: '分',
    },
  },
  fileLoading: {
    loading: '読み込み中',
  },
  fileActions: {
    copyLink: {
      protectedNotSupported: '保護されたファイルは外部シェアをサポートしていません。アプリケーション内で表示してください',
      publicLinkCopied: '公開リンクをクリップボードにコピーしました',
      privateLinkCopied: 'プライベートリンクをクリップボードにコピーしました',
      copyFailed: 'リンクのコピーに失敗しました',
    },
    delete: {
      success: 'ファイルが正常に削除されました',
      failed: '削除に失敗しました。もう一度お試しください',
    },
    accessLevel: {
      public: '公開',
      private: 'プライベート',
      protected: '保護',
      switched: 'ファイルを{level}に設定しました',
      switchFailed: '切り替えに失敗しました。もう一度お試しください',
    },
    batch: {
      noPublicFiles: '選択したファイルにコピー可能な公開ファイルがありません',
      filteredFiles: '{count}件のプライベート/保護ファイルを除外しました',
      linksCopied: '{count}件のファイルリンクをコピーしました',
      copyFailed: '一括コピーに失敗しました',
      deleteSuccess: '{count}件のリソースを正常に削除しました',
      deletePartial: '{successCount}/{totalCount}件のリソースを正常に削除しました',
      deleteFailed: '削除に失敗しました。もう一度お試しください',
    },
  },
}
