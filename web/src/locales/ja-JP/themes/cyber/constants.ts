/**
 * Constants Text - Cyber Style
 */
export const constants = {
  fileTypes: {
    image: 'データファイル',
    video: 'ビデオストリーム',
    document: 'ドキュメントデータ',
    archive: 'アーカイブパッケージ',
    audio: 'オーディオストリーム',
    other: '不明タイプ',
  },

  accessLevels: {
    private: 'プライベートモード',
    public: 'パブリックアクセス',
    protected: '保護',
  },

  sortOptions: {
    created_at: '作成時間',
    updated_at: '更新時間',
    size: 'データサイズ',
    name: 'ファイル名',
    views: 'アクセス数',
    downloads: 'ダウンロード数',
  },

  categorySortBy: {
    sort_order: '優先度順',
    name: '名前順',
    usage_count: '使用頻度順',
    created_at: '作成時間順',
  },

  sortOrder: {
    asc: '昇順',
    desc: '降順',
  },

  categoryType: {
    all: 'すべてのカテゴリ',
    hot: 'ホットタグ',
    normal: '通常タグ',
  },

  batchOperations: {
    delete: 'データパージ',
    move: 'データ転送',
    update_access: '権限変更',
    update_folder: 'フォルダに移動',
    download: 'データ抽出',
  },

  accessControl: {
    modes: {
      ipWhitelist: 'IPホワイトリスト',
      ipBlacklist: 'IPブラックリスト',
      domainWhitelist: 'ドメインホワイトリスト',
      domainBlacklist: 'ドメインブラックリスト',
    },
    restrictionModes: {
      strict: '厳格モード',
      moderate: '中程度モード',
      loose: '緩いモード',
    },
    blockActions: {
      block: '完全ブロック',
      redirect: 'リダイレクト',
      thumbnail: 'サムネイルを返す',
      watermark: '透かしを追加',
    },
    warnings: {
      block:
        '訪問者に403エラーを返し、データリソースへのアクセスを完全に拒否します。ユーザーエクスペリエンスに影響を与える可能性があります。',
      redirect: '訪問者を指定したURLにリダイレクトします。ターゲットURLがアクセス可能であることを確認してください。',
      thumbnail: '許可されていない訪問者にはサムネイルバージョンのみを表示し、オリジナルデータを盗難から保護します。',
      watermark: 'ファイルに自動的に透かしを追加します。所有権をマークできますが、オリジナルの外観が変わります。',
    },
    descriptions: {
      ipWhitelist: 'ホワイトリストモード: リスト内のIPのみアクセスを許可',
      ipBlacklist: 'ブラックリストモード: リスト内のIPからのアクセスをブロック',
      domainWhitelist: 'ホワイトリストモード: リスト内のドメインのみ参照を許可',
      domainBlacklist: 'ブラックリストモード: リスト内のドメインからの参照をブロック',
      strictMode: '厳格モード: ルールに完全に一致する必要がある',
      moderateMode: '中程度モード: 部分的な一致が受け入れられる',
      looseMode: '緩いモード: 最小限の制限',
    },
  },

  storageOptions: {
    unlimited: '無制限',
    storageTypes: {
      local: 'ローカルストレージ',
      s3: 'S3クラウド',
      aliyun: 'アリババOSS',
      tencent: 'テンセントCOS',
    },
    channelTypes: {
      oss: 'アリババOSS',
      cos: 'テンセントCOS',
      rainyun: 'RainYun',
      local: 'ローカルストレージ',
      unknown: '不明ノード',
    },
    storageClasses: {
      standard: '標準ストレージ',
      standardMultiAz: '標準ストレージ（マルチAZ）',
      infrequent: '低頻度ストレージ',
      infrequentMultiAz: '低頻度ストレージ（マルチAZ）',
      infrequentAccess: '低頻度アクセス',
      intelligentTiering: 'インテリジェント階層化',
      intelligentTieringMultiAz: 'インテリジェント階層化（マルチAZ）',
      archive: 'アーカイブストレージ',
      deepArchive: 'ディープアーカイブ',
      coldArchive: 'コールドアーカイブ',
      glacier: 'Glacierストレージ',
    },
  },

  filterOptions: {
    sortOptions: {
      newest: '最新インジェクション',
      oldest: '最古インジェクション',
      name: '名前',
      size: 'データサイズ',
      width: '幅',
      height: '高さ',
      quality: 'データ品質',
      nsfw_score: 'センシティビティスコア',
    },
    resolutionOptions: {
      low: '低解像度',
      '720p': '720p',
      '1080p': '1080p',
      '2k': '2k',
      '4k': '4K',
      '8k': '8K',
    },
    imageFormatOptions: {
      jpg: 'JPG',
      jpeg: 'JPEG',
      png: 'PNG',
      gif: 'GIF',
      webp: 'WebP',
      bmp: 'BMP',
      svg: 'SVG',
      ico: 'ICO',
    },
    fileTypeOptions: {
      all: 'すべてのデータタイプ',
      jpeg: 'JPEGデータ',
      png: 'PNGデータ',
      gif: 'GIFデータ',
      webp: 'WebPデータ',
      svg: 'SVGデータ',
      bmp: 'BMPデータ',
      ico: 'ICOアイコン',
    },
  },

  category: {
    sortByOptions: {
      sort_order: '順序順',
      name: '名前順',
      usage_count: '使用数順',
      created_at: '作成時間順',
    },
    sortOrderOptions: {
      asc: '昇順',
      desc: '降順',
    },
    typeOptions: {
      all: 'すべて',
      popular: '人気',
      normal: '通常',
    },
  },

  contentReview: {
    actions: {
      all: 'すべてのアクション',
      approve: '検証通過',
      reject: 'データ拒否',
    },
    sort: {
      newest: '最新インジェクション',
      oldest: '最古インジェクション',
      size: 'データサイズ',
      nsfw: 'リスクスコア',
    },
    nsfwLevels: {
      safe: '安全レベル',
      low: '低リスク',
      medium: '中リスク',
      high: '高リスク',
    },
    columns: {
      file: 'データファイル',
      fileInfo: 'データ情報',
      action: 'レビューアクション',
      softDelete: 'マークステータス',
      auditor: '監査者',
      auditTime: 'レビュー時間',
      reason: 'レビュー理由',
      nsfw: 'リスクスコア',
      operations: 'システム操作',
    },
    defaultReasons: {
      approve: 'データ検証通過',
      reject: 'システム仕様を満たしていません',
      batchApprove: '一括検証通過',
      batchReject: '一括データ拒否',
    },
    messages: {
      approveSuccess: 'データ検証通過',
      rejectSuccess: 'データ拒否済み',
      batchApproveSuccess: '一括検証成功、{count}件のデータアイテムを処理しました',
      batchRejectSuccess: '一括拒否成功、{count}件のデータアイテムを処理しました',
      restoreSuccess: 'データ復元済み',
      hardDeleteSuccess: 'データ完全パージ済み',
      batchRestoreSuccess: '一括復元完了: {success}件成功',
      batchRestoreWithFail: '一括復元完了: {success}件成功、{fail}件失敗',
      batchHardDeleteSuccess: '一括パージ完了: {success}件成功',
      batchHardDeleteWithFail: '一括パージ完了: {success}件成功、{fail}件失敗',
      loadError: 'データ読み込み失敗',
      operationError: '操作実行失敗',
      noSelection: '操作するデータアイテムを選択してください',
      queueEmpty: 'キュークリア済み',
      queueEmptyDesc: 'すべてのデータがレビューを完了しました',
      logsEmpty: 'レビューログなし',
      logsEmptyDesc: 'システムにまだレビュー記録がありません',
    },
    warnings: {
      hardDelete: '警告: この操作は元に戻せません',
      hardDeleteDesc: 'データはサーバーから完全にパージされ、復元できません。',
      batchOperation: 'この操作は複数のデータアイテムに影響します',
    },
    descriptions: {
      approveEffects: [
        'データはパブリックデータベースに表示されます',
        'データステータスが「検証済み」に変更されました',
        'ユーザーは通常通りデータをシェアおよび配布できます',
      ],
      restoreEffects: ['データはマーク削除ステータスから通常ステータスに復元されます'],
    },
  },

  share: {
    sortOptions: {
      date: '日付',
      name: '名前',
      size: 'サイズ',
    },
    layoutOptions: {
      grid: {
        label: 'グリッドレイアウト',
        description: '整然としたグリッド表示',
      },
      waterfall: {
        label: 'ウォーターフォール',
        description: 'データの高さによって自動配置',
      },
      masonry: {
        label: 'メイソンリーレイアウト',
        description: '左寄せウォーターフォール',
      },
      large: {
        label: '大型モード',
        description: 'より大きなサムネイルを表示',
      },
    },
  },

  navigation: {
    admin: {
      dashboard: 'コントロールパネル',
      files: 'データ管理',
      tags: 'タグシステム',
      categories: 'カテゴリシステム',
      shares: 'シェア管理',
      ai: 'インテリジェンスシステム',
      contentReview: 'コンテンツレビュー',
      tagging: '自動タグ付け',
      vectors: 'ベクトルシステム',
      users: 'ユーザー管理',
      channels: 'チャンネル管理',
      construction: 'サイト構築',
      announcements: 'お知らせシステム',
      settings: 'システム設定',
    },
    docs: {
      overview: 'API概要',
      authentication: '認証メカニズム',
      upload: 'アップロードAPI',
      limits: '制限説明',
      examples: 'コード例',
      tester: 'APIテスター',
      faq: 'FAQ',
    },
    main: {
      home: 'ホーム',
      explore: 'データベース',
      upload: 'アップロード',
      random: 'ランダム',
      docs: 'API',
    },
    user: {
      dashboard: 'コントロールパネル',
      myFiles: 'マイデータ',
      folders: 'フォルダ',
      tagManage: 'タグ管理',
      categoryManage: 'カテゴリ管理',
      shares: 'マイシェア',
      settings: '設定',
    },
    settings: {
      profile: '個人情報',
      security: 'セキュリティ設定',
      apikey: 'APIキー',
      storage: 'ストレージ設定',
    },
    adminSettings: {
      website: 'ウェブサイト設定',
      upload: 'アップロード設定',
      security: 'セキュリティ設定',
      registration: '登録設定',
      mail: 'メール設定',
      ai: 'AI設定',
      vector: 'ベクトル設定',
    },
  },

  api: {
    errors: {
      badRequest: 'リクエストパラメータ例外',
      unauthorized: '不正アクセス',
      forbidden: 'アクセス拒否',
      notFound: 'リソースが見つかりません',
      methodNotAllowed: 'メソッドが許可されていません',
      conflict: 'データ競合',
      payloadTooLarge: 'データパッケージが大きすぎます',
      unprocessableEntity: 'データ検証失敗',
      tooManyRequests: 'リクエスト頻度が高すぎます',
      unavailableForLegalReasons: 'IPアドレスが禁止されています',
      internalServerError: 'システム内部エラー',
      badGateway: 'ゲートウェイ例外',
      serviceUnavailable: 'サービス利用不可',
      gatewayTimeout: 'ゲートウェイタイムアウト',
    },
  },

  systemOptions: {
    sensitiveContent: {
      markOnly: 'マークのみ',
      pendingReview: 'システムレビュー',
      autoDelete: '自動パージ',
    },
    sizeUnits: {
      bytes: 'バイト（B）',
      kilobytes: 'キロバイト（KB）',
      megabytes: 'メガバイト（MB）',
      gigabytes: 'ギガバイト（GB）',
      terabytes: 'テラバイト（TB）',
    },
    userRoles: {
      normal: '通常ユーザー',
      advanced: '上級ユーザー',
      admin: 'システム管理者',
      superAdmin: 'スーパー管理者',
    },
    userStatus: {
      normal: '通常',
      disabled: '無効',
      pending: 'レビュー待ち',
    },
    imageStatus: {
      normal: '通常',
      hidden: '非表示',
      pending: 'レビュー待ち',
      deleted: '削除済み',
    },
    shareStatus: {
      normal: '通常',
      disabled: '無効',
      expired: '期限切れ',
    },
  },
}
