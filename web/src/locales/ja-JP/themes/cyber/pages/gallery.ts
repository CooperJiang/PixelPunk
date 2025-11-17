/**
 * Gallery Browsing Text - Cyber Style
 */
export const gallery = {
  pages: {
    gallery: 'データベース',
    upload: 'データアップロードセンター',
    login: '本人確認',
    register: 'システム登録',
    settings: 'システム設定',
    profile: 'ユーザープロフィール',
    folders: 'データコンテナ',
    dashboard: 'コントロールパネル',
    admin: '管理者ターミナル',
  },
  viewer: {
    previous: '前のフレーム',
    next: '次のフレーム',
    fitMode: 'ビューポートにフィット',
    fillMode: 'ビューポートをフィル',
    fullscreen: '没入モード',
    exitFullscreen: '没入を終了',
    originalSize: '元の解像度',
    imageInfo: 'データプロパティ',
    similarImages: '類似データ',
  },
  waterfallLayout: {
    loadFailed: 'データ破損',
    collapse: '折りたたむ',
  },
  file: {
    loadFailed: 'データ読み込み失敗',
    nsfwWarning: '制限コンテンツ',
  },
  enhancedFilePreview: {
    defaultImageName: 'データファイル',
    fillMode: 'フィルモード',
    fitMode: 'フィットモード',
    switchToFit: 'フィットに切り替え',
    switchToFill: 'フィルに切り替え',
    wheelZoom: 'ホイールズーム',
    dragMove: 'ドラッグ移動',
    spaceKey: 'スペース',
    exitPreview: 'プレビューを終了',
    exitFullscreen: 'フルスクリーンを終了',
    exitPreviewEsc: 'プレビューを終了（ESC）',
    exitFullscreenEsc: 'フルスクリーンを終了（ESC）',
    enterFullscreen: 'ブラウザフルスクリーンに入る',
    fullscreen: 'フルスクリーン',
    preview: 'プレビュー',
    mode: 'モード',
  },
  fileDetailModal: {
    title: 'データ詳細',
    untitled: '無名データ',

    basicInfo: {
      title: '基本情報',
      fileId: 'データID',
      originalName: 'オリジナル名',
      fileFormat: 'データフォーマット',
      fileSize: 'データサイズ',
      imageDimensions: 'データサイズ',
      visibility: '可視性',
      uploader: 'アップローダー',
      uploadTime: 'アップロード時間',
      updateTime: '更新時間',
      fileHash: 'データハッシュ',
      viewCount: 'ビュー数',
      status: 'ステータス',
      times: '回',
      count: '件',
      unknownResolution: '不明な解像度',
      none: 'なし',
    },

    statusLabels: {
      duplicate: '重複データ',
      recommended: '推奨',
      nsfw: 'NSFW',
    },

    visibility: {
      public: '公開プロトコル',
      private: 'プライベートアーカイブ',
      protected: 'セキュリティシールド',
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
      title: 'コンテンツ安全性',
      nsfwStatus: 'NSFWステータス',
      inappropriate: '不適切',
      safe: '安全',
      nsfwScore: 'NSFWスコア',
      high: '高リスク',
      medium: '中リスク',
      low: '低リスク',
    },

    visualAnalysis: {
      title: '視覚分析',
      dominantColor: '主要色',
      imageDimensions: 'データサイズ',
      aspectRatio: 'アスペクト比',
      resolution: '解像度',
      colorPalette: 'カラーパレット',
      composition: '構図',
      objectsCount: 'オブジェクト数',
    },

    ai: {
      tags: 'AIタグ',
      description: 'AI説明',
    },

    actions: {
      deleteImage: 'データをクリア',
      cancelRecommend: '推奨をキャンセル',
      setRecommend: '推奨に設定',
      close: '切断',
    },
  },
  fileExpiryTag: {
    statuses: {
      expired: 'データ期限切れ',
      expiringSoon: 'まもなく期限切れ',
      active: '期間限定保存',
    },
    labels: {
      storageDuration: '保存期間',
      expiresAt: '有効期限',
      expiredAtPrefix: '期限切れ日時',
      expiredAtSuffix: '',
    },
    units: {
      day: '日',
      hour: '時間',
      minute: '分',
    },
  },
  fileLoading: {
    loading: 'データを転送中',
  },
  fileActions: {
    copyLink: {
      protectedNotSupported: '保護されたデータは外部転送をサポートしていません。ターミナル内で閲覧してください',
      publicLinkCopied: '公開リンクをクリップボードにコピーしました',
      privateLinkCopied: 'プライベートリンクをクリップボードにコピーしました',
      copyFailed: 'リンクコピー失敗',
    },
    delete: {
      success: 'データを削除しました',
      failed: '削除失敗、再試行してください',
    },
    accessLevel: {
      public: '公開プロトコル',
      private: 'プライベートアーカイブ',
      protected: 'セキュリティシールド',
      switched: 'データを{level}に設定しました',
      switchFailed: '切り替え失敗、再試行してください',
    },
    batch: {
      noPublicFiles: '選択したデータにコピーする公開リソースがありません',
      filteredFiles: '{count}個のプライベート/保護データをフィルタリングしました',
      linksCopied: '{count}個のデータリンクをコピーしました',
      copyFailed: '一括コピー失敗',
      deleteSuccess: '{count}個のリソースを削除しました',
      deletePartial: '{successCount}/{totalCount}個のリソースを削除しました',
      deleteFailed: '削除失敗、再試行してください',
    },
  },
}
