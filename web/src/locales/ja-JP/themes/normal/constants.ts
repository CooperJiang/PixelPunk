/**
 * Constants
 */
export const constants = {
  fileTypes: {
    image: 'ファイル',
    video: '動画',
    document: 'ドキュメント',
    archive: 'アーカイブ',
    audio: '音声',
    other: 'その他',
  },

  accessLevels: {
    private: 'プライベート',
    public: '公開',
    protected: '保護',
  },

  sortOptions: {
    created_at: '作成日時',
    updated_at: '更新日時',
    size: 'ファイルサイズ',
    name: 'ファイル名',
    views: '閲覧数',
    downloads: 'ダウンロード数',
  },

  categorySortBy: {
    sort_order: '並び順',
    name: '名前',
    usage_count: '使用回数',
    created_at: '作成日時',
  },

  sortOrder: {
    asc: '昇順',
    desc: '降順',
  },

  categoryType: {
    all: 'すべて',
    hot: '人気',
    normal: '通常',
  },

  batchOperations: {
    delete: '削除',
    move: '移動',
    update_access: '権限更新',
    update_folder: 'フォルダに移動',
    download: 'ダウンロード',
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
      loose: '緩和モード',
    },
    blockActions: {
      block: '完全ブロック',
      redirect: 'リダイレクト',
      thumbnail: 'サムネイルを返す',
      watermark: 'ウォーターマークを追加',
    },
    warnings: {
      block: '訪問者に403エラーを返し、ファイルリソースへのアクセスを完全に拒否します。ユーザー体験に影響する可能性があります。',
      redirect: '訪問者を指定したURLにリダイレクトします。ターゲットURLがアクセス可能であることを確認してください。',
      thumbnail: '未承認の訪問者にはファイルのサムネイルバージョンのみを表示し、オリジナル画像の盗用を防ぎます。',
      watermark: 'ファイルに自動的にウォーターマークを追加します。所有権を識別できますが、オリジナル画像の外観が変更されます。',
    },
    descriptions: {
      ipWhitelist: 'ホワイトリストモード：リスト内のIPのみアクセスを許可',
      ipBlacklist: 'ブラックリストモード：リスト内のIPのアクセスをブロック',
      domainWhitelist: 'ホワイトリストモード：リスト内のドメインのみ参照を許可',
      domainBlacklist: 'ブラックリストモード：リスト内のドメインの参照をブロック',
      strictMode: '厳格モード：ルールに完全一致する必要があります',
      moderateMode: '中程度モード：部分一致で十分です',
      looseMode: '緩和モード：最小限の制限',
    },
  },

  storageOptions: {
    unlimited: '無制限',
    storageTypes: {
      local: 'ローカルストレージ',
      s3: 'S3ストレージ',
      aliyun: 'アリババクラウドOSS',
      tencent: 'テンセントクラウドCOS',
    },
    channelTypes: {
      oss: 'アリババクラウドOSS',
      cos: 'テンセントクラウドCOS',
      rainyun: 'RainYun',
      local: 'ローカルストレージ',
      unknown: '不明',
    },
    storageClasses: {
      standard: '標準ストレージ',
      standardMultiAz: '標準ストレージ（マルチAZ）',
      infrequent: '低頻度ストレージ',
      infrequentMultiAz: '低頻度ストレージ（マルチAZ）',
      infrequentAccess: '低頻度アクセス',
      intelligentTiering: 'インテリジェント階層ストレージ',
      intelligentTieringMultiAz: 'インテリジェント階層ストレージ（マルチAZ）',
      archive: 'アーカイブストレージ',
      deepArchive: 'ディープアーカイブストレージ',
      coldArchive: 'コールドアーカイブストレージ',
      glacier: 'Glacierストレージ',
    },
  },

  filterOptions: {
    sortOptions: {
      newest: '最新アップロード',
      oldest: '最古アップロード',
      name: '名前',
      size: 'ファイルサイズ',
      width: '幅',
      height: '高さ',
      quality: 'ファイル品質',
      nsfw_score: 'NSFWスコア',
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
      all: 'すべてのファイルタイプ',
      jpeg: 'JPEGファイル',
      png: 'PNGファイル',
      gif: 'GIFファイル',
      webp: 'WebPファイル',
      svg: 'SVGファイル',
      bmp: 'BMPファイル',
      ico: 'ICOアイコン',
    },
  },

  category: {
    sortByOptions: {
      sort_order: '並び順',
      name: '名前',
      usage_count: '使用回数',
      created_at: '作成日時',
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
      all: 'すべての操作',
      approve: '承認',
      reject: '拒否',
    },
    sort: {
      newest: '最新アップロード',
      oldest: '最古アップロード',
      size: 'ファイルサイズ',
      nsfw: 'NSFWスコア',
    },
    nsfwLevels: {
      safe: '安全',
      low: '低',
      medium: '中',
      high: '高',
    },
    columns: {
      file: 'ファイル',
      fileInfo: 'ファイル情報',
      action: 'レビュー操作',
      softDelete: '削除ステータス',
      auditor: '監査者',
      auditTime: '監査日時',
      reason: 'レビュー理由',
      nsfw: 'NSFWスコア',
      operations: '操作',
    },
    defaultReasons: {
      approve: 'ファイルレビュー承認',
      reject: 'プラットフォーム基準を満たしていません',
      batchApprove: '一括承認',
      batchReject: '一括拒否',
    },
    messages: {
      approveSuccess: 'ファイル承認済み',
      rejectSuccess: 'ファイル拒否済み',
      batchApproveSuccess: '一括承認成功、{count}件のリソースを処理',
      batchRejectSuccess: '一括拒否成功、{count}件のリソースを処理',
      restoreSuccess: 'ファイル復元済み',
      hardDeleteSuccess: 'ファイル完全削除済み',
      batchRestoreSuccess: '一括復元完了：{success}件成功',
      batchRestoreWithFail: '一括復元完了：{success}件成功、{fail}件失敗',
      batchHardDeleteSuccess: '一括削除完了：{success}件成功',
      batchHardDeleteWithFail: '一括削除完了：{success}件成功、{fail}件失敗',
      loadError: 'データの読み込みに失敗しました',
      operationError: '操作に失敗しました',
      noSelection: '操作する項目を選択してください',
      queueEmpty: 'キューが空です',
      queueEmptyDesc: 'すべてのファイルがレビュー済みです',
      logsEmpty: 'レビューレコードがありません',
      logsEmptyDesc: 'まだレビュー操作レコードがありません',
    },
    warnings: {
      hardDelete: '警告：この操作は元に戻せません',
      hardDeleteDesc: 'ファイルはサーバーから完全に削除され、復元できません。',
      batchOperation: 'この操作は複数の項目に影響します',
    },
    descriptions: {
      approveEffects: [
        'ファイルが公開ギャラリーに表示されます',
        'ファイルステータスが「レビュー済み」に変更されます',
        'ユーザーは通常通り共有・配布できます',
      ],
      restoreEffects: ['ファイルがソフト削除ステータスから通常ステータスに復元されます'],
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
        description: 'ファイルの高さで自動配置',
      },
      masonry: {
        label: 'メイソンリーレイアウト',
        description: '左揃えウォーターフォール',
      },
      large: {
        label: '大画像モード',
        description: '大きなサムネイルを表示',
      },
    },
  },

  navigation: {
    admin: {
      dashboard: 'ダッシュボード',
      files: 'ファイル管理',
      tags: 'タグ管理',
      categories: 'カテゴリ管理',
      shares: '共有管理',
      ai: 'AI管理',
      contentReview: 'コンテンツレビュー',
      tagging: 'タグ付け管理',
      vectors: 'ベクトル管理',
      users: 'ユーザー管理',
      channels: 'チャンネル管理',
      construction: 'ウェブサイト構築',
      announcements: 'お知らせ管理',
      settings: 'グローバル設定',
    },
    docs: {
      overview: 'API概要',
      authentication: '認証',
      upload: 'アップロードAPI',
      limits: '制限',
      examples: 'コード例',
      tester: 'APIテスター',
      faq: 'FAQ',
    },
    main: {
      home: 'ホーム',
      explore: 'ギャラリー',
      upload: 'アップロード',
      random: 'ランダム',
      docs: 'API',
    },
    user: {
      dashboard: 'ダッシュボード',
      myFiles: 'マイファイル',
      folders: 'フォルダ',
      tagManage: 'タグ管理',
      categoryManage: 'カテゴリ管理',
      shares: 'マイ共有',
      settings: '設定',
    },
    settings: {
      profile: 'プロフィール',
      security: 'セキュリティ設定',
      apikey: 'APIキー',
      storage: 'ストレージ設定',
    },
    adminSettings: {
      website: 'サイト設定',
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
      badRequest: '不正なリクエスト',
      unauthorized: '認証が必要です',
      forbidden: 'アクセス拒否',
      notFound: 'リソースが見つかりません',
      methodNotAllowed: 'メソッドが許可されていません',
      conflict: 'リソース競合',
      payloadTooLarge: 'ペイロードが大きすぎます',
      unprocessableEntity: '検証に失敗しました',
      tooManyRequests: 'リクエストが多すぎます',
      unavailableForLegalReasons: 'IPアドレスが無効化されています',
      internalServerError: '内部サーバーエラー',
      badGateway: '不正なゲートウェイ',
      serviceUnavailable: 'サービス利用不可',
      gatewayTimeout: 'ゲートウェイタイムアウト',
    },
  },

  systemOptions: {
    sensitiveContent: {
      markOnly: 'マークのみ',
      pendingReview: '管理者レビュー',
      autoDelete: '自動削除',
    },
    sizeUnits: {
      bytes: 'バイト（B）',
      kilobytes: 'キロバイト（KB）',
      megabytes: 'メガバイト（MB）',
      gigabytes: 'ギガバイト（GB）',
      terabytes: 'テラバイト（TB）',
    },
    userRoles: {
      normal: '一般ユーザー',
      advanced: '上級ユーザー',
      admin: '管理者',
      superAdmin: 'スーパー管理者',
    },
    userStatus: {
      normal: '通常',
      disabled: '無効',
      pending: '保留中',
    },
    imageStatus: {
      normal: '通常',
      hidden: '非表示',
      pending: '保留中',
      deleted: '削除済み',
    },
    shareStatus: {
      normal: '通常',
      disabled: '無効',
      expired: '期限切れ',
    },
  },
}
