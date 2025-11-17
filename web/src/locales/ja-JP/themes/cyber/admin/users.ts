export const users = {
  title: 'ノード管理',
  subtitle: 'すべての接続ノードを管理および監視',

  columns: {
    user: 'ノード',
    email: '連絡先アドレス',
    registeredAt: 'アクセス時間',
    lastActivity: '最終アクティビティ',
    status: 'ステータス',
    role: '権限',
    storage: 'ストレージ/帯域幅',
    actions: '操作',
    userId: 'ID',
    neverActive: 'アクティブなし',
    storageLabel: 'ストレージ',
    bandwidthLabel: '帯域幅',
  },

  status: {
    all: 'すべて',
    active: 'オンライン',
    disabled: '無効',
    suspended: '停止',
  },

  role: {
    all: 'すべて',
    superAdmin: 'スーパー管理者',
    admin: '管理者ノード',
    user: '通常ノード',
    guest: 'ゲストノード',
    unknown: '不明',
  },

  roles: {
    all: 'すべて',
    superAdmin: 'スーパー管理者',
    admin: '管理者ノード',
    normal: '通常ノード',
    guest: 'ゲストノード',
    unknown: '不明',
  },

  actions: {
    view: '詳細を表示',
    edit: 'ノードを編集',
    disable: 'ノードを無効化',
    enable: 'ノードを有効化',
    delete: 'ノードを削除',
    export: 'エクスポート',
    import: 'インポート',
    refresh: '更新',
  },

  buttons: {
    createUser: '新規ノード',
    create: 'ノードを作成',
    batchOperation: '一括操作',
    filter: 'フィルター',
    resetFilter: 'フィルターをリセット',
    applyFilter: 'フィルターを適用',
    cancel: 'キャンセル',
    confirm: '確認',
    save: '保存',
    close: '閉じる',
  },

  toast: {
    userCreated: 'ノード作成成功',
    userUpdated: 'ノード情報更新済み',
    userDeleted: 'ノード削除成功',
    statusChanged: 'ステータス変更済み',
    batchOperationSuccess: '一括操作完了',
    operationFailed: '操作失敗',
    pleaseSelectUsers: '操作するノードを選択してください',
    userEnabled: 'ノード有効化済み',
    userDisabled: 'ノード無効化済み',
    userDataNotFound: 'ノードデータが見つかりませんでした',
    success: '{operation}成功',
    error: '{operation}失敗',
    exportSuccess: 'ノードデータエクスポート成功',
    exportError: 'ノードデータエクスポート失敗',
  },

  confirm: {
    enableUser: 'このノードを有効にしますか？',
    disableUser: 'このノードを無効にしますか？',
    deleteUser: 'このノードを削除しますか？',
    delete: '選択したノードを削除しますか？',
    batchDelete: '選択した{count}件のノードを削除しますか？',
    batchDisable: '選択したノードを無効にしますか？',
    batchEnable: '選択した{count}件のノードを有効にしますか？',
    disableTitle: '{count}件のノードを無効にしようとしています',
    disableMessage1: '無効化後、これらのノードはシステムにアクセスできなくなります',
    disableMessage2: '現在オンラインのノードセッションは即座に終了します',
    disableMessage3: 'すべての関連プロトコルリンクが無効になります',
  },

  disableDialog: {
    title: 'ノード無効化確認',
    warning: 'ノードを無効にしようとしています',
    consequences: 'この操作により以下が発生します:',
    consequence1: 'ノードはシステムにログインできません',
    consequence2: 'すべてのノードセッションが即座に終了します',
    consequence3: 'すべてのノードプロトコルリンクが無効になります',
    consequence4: 'ノードはデータにアクセスできません',
    inputLabel: '確認のためノード名を入力',
    inputPlaceholder: 'ノード名を入力',
    inputHint: '"{username}"を入力して無効化を確認',
    usernameNotMatch: 'ノード名が一致しません',
    confirmButton: '無効化を確認',
    cancelButton: 'キャンセル',
  },

  filter: {
    title: 'ノードをフィルター',
    keyword: 'キーワード',
    keywordPlaceholder: 'ノード名または連絡先アドレスを検索',
    status: {
      label: 'ノードステータス',
      all: 'すべてのノード',
      normal: '通常ノード',
      disabled: '無効ノード',
      deleted: '削除済みノード',
    },
    role: {
      label: 'ノード権限',
      all: 'すべての権限',
      superAdmin: 'スーパー管理者ノード',
      admin: '管理者ノード',
      user: '通常ノード',
    },
    registrationTime: 'アクセス時間',
    registrationDate: 'アクセス日',
    dateRangePlaceholder: '時間範囲を選択',
    startDatePlaceholder: '開始日',
    endDatePlaceholder: '終了日',
    startDate: '開始日',
    endDate: '終了日',
    sortBy: '並び順',
    sort: {
      label: '並び順',
      placeholder: '並び順方法を選択',
      default: 'デフォルト',
      newest: '最新アクセス',
      oldest: '最古アクセス',
      lastActivityNewest: '最近アクティブ',
      lastActivityOldest: '最古アクティブ',
      usernameAZ: 'ノード名 A-Z',
      usernameZA: 'ノード名 Z-A',
    },
  },

  batch: {
    title: '一括操作',
    selected: '{count}件のノードを選択',
    selectAll: 'すべて選択',
    deselectAll: 'すべて解除',
    clearSelection: '選択をクリア',
    enable: '有効化',
    disable: '無効化',
    delete: '削除',
    export: 'エクスポート',
    sendEmail: '連絡を送信',
    role: '権限を設定',
    apply: '適用',
    cancel: 'キャンセル',
    noSelection: '操作するノードを選択してください',
    setRole: '権限を設定',
    selectRole: '権限を選択してください',
    selectUsersToExport: 'エクスポートするノードを選択してください',
    confirmDisable: '無効化を確認',
  },

  form: {
    createTitle: '新しいノードを作成',
    editTitle: 'ノードを編集',
    username: 'ノード名',
    usernamePlaceholder: 'ノード名を入力',
    usernameHint: '3-20文字、文字、数字、アンダースコアのみ',
    email: '連絡先アドレス',
    emailPlaceholder: '連絡先アドレスを入力',
    emailHint: 'ログインおよび通知受信に使用',
    password: 'パスワード',
    passwordPlaceholder: 'パスワードを入力',
    passwordHint: '少なくとも8文字',
    confirmPassword: 'パスワードを確認',
    confirmPasswordPlaceholder: 'パスワードを再度入力',
    role: '権限',
    rolePlaceholder: 'ノード権限を選択',
    roleUser: '通常ノード',
    roleAdmin: '管理者ノード',
    roleSuperAdmin: 'スーパー管理者ノード',
    status: 'ステータス',
    statusActive: '通常ノード',
    statusDisabled: '無効ノード',
    storageConfig: 'ストレージ設定',
    storageLimit: 'ストレージ制限',
    storageLimitPlaceholder: 'ストレージ制限を入力',
    storageLimitDefault: 'デフォルト5GB',
    storageLimitHint: '単位: GB、0は無制限',
    bandwidthLimit: '帯域幅制限',
    bandwidthLimitPlaceholder: '帯域幅制限を入力',
    bandwidthLimitDefault: 'デフォルト100GB',
    bandwidthLimitHint: '単位: GB、0は無制限',
    avatar: 'アバター',
    avatarUpload: 'アバターをアップロード',
    usernameRequired: 'ノード名は空にできません',
    usernameMinLength: 'ノード名は少なくとも2文字',
    usernameMaxLength: 'ノード名は20文字を超えることはできません',
    emailRequired: '連絡先アドレスは空にできません',
    emailInvalid: '連絡先アドレス形式が正しくありません',
    passwordRequired: 'パスワードは空にできません',
    passwordMinLength: 'パスワードは少なくとも6文字',
    passwordMaxLength: 'パスワードは20文字を超えることはできません',
    required: 'このフィールドは必須です',
    invalidEmail: '連絡先アドレス形式が正しくありません',
    passwordTooShort: 'パスワードは少なくとも8文字必要です',
    passwordNotMatch: 'パスワードが一致しません',
  },

  detail: {
    title: 'ノード詳細',
    loading: 'データを読み込み中...',
    noData: 'データなし',
    basicInfo: '基本情報',
    accountInfo: 'アカウント情報',
    statisticsInfo: '統計情報',
    activityInfo: 'アクティビティ情報',
    usageStats: '使用統計',
    storageManagement: 'ストレージ管理',
    quickActions: 'クイックアクション',
    userId: 'ノードID',
    username: 'ノード名',
    email: '連絡先アドレス',
    role: '権限',
    registeredAt: 'アクセス時間',
    lastLogin: '最終ログイン',
    lastLoginAt: '最終ログイン',
    lastActivity: '最終アクティビティ',
    lastActivityAt: '最終アクティビティ',
    lastActivityIp: 'IPアドレス',
    uploadedFiles: 'アップロードデータ数',
    shareCount: 'プロトコル数',
    totalViews: '総閲覧数',
    usedStorage: '使用ストレージ',
    storageLimit: 'ストレージ制限',
    monthlyBandwidth: '月間帯域幅',
    bandwidthLimit: '帯域幅制限',
    fileCount: 'データ数',
    unlimited: '無制限',
    never: 'なし',
    currentSetting: '現在の設定',
    enterValue: '値を入力',
    unit: '単位',
    updateStorageSettings: 'ストレージ設定を更新',
    status: {
      label: 'ステータス',
      active: 'オンライン',
      disabled: '無効',
      deleted: '削除済み',
      unknown: '不明',
    },
    roles: {
      superAdmin: 'スーパー管理者ノード',
      admin: '管理者ノード',
      user: '通常ノード',
      unknown: '不明',
    },
    units: {
      mb: 'MB',
      gb: 'GB',
      tb: 'TB',
    },
    actions: {
      resetPassword: 'パスワードをリセット',
      resetPasswordDesc: 'このノードの新しいパスワードを生成',
      sendEmail: 'メールを送信',
      sendEmailDesc: 'このノードにメールを送信',
      enableUser: 'ノードを有効化',
      enableUserDesc: 'このノードのログインを許可',
      disableUser: 'ノードを無効化',
      disableUserDesc: 'このノードのログインを防止',
    },
    messages: {
      storageUpdateSuccess: 'ストレージ設定更新済み',
      passwordResetSuccess: 'パスワードをリセットしました: {password}',
      emailSentSuccess: 'メール送信成功',
      userEnabled: 'ノード有効化済み',
      userDisabled: 'ノード無効化済み',
    },
    errors: {
      fetchFailed: 'ノード情報の取得に失敗しました',
      storageUpdateFailed: 'ストレージ設定の更新に失敗しました',
      passwordResetFailed: 'パスワードリセット失敗',
      emailSentFailed: 'メール送信失敗',
      operationFailed: '操作失敗',
    },
  },

  search: {
    placeholder: 'ノード名または連絡先アドレスを検索...',
    button: '検索',
  },

  sendEmail: {
    title: '連絡を送信',
    recipients: '受信者ノード',
    recipientsHint: '{count}件のノードを選択',
    subject: '件名',
    subjectPlaceholder: '連絡件名を入力',
    content: 'コンテンツ',
    contentPlaceholder: '連絡コンテンツを入力',
    send: '送信',
    cancel: 'キャンセル',
    sending: '送信中...',
    sent: '連絡送信済み',
    failed: '送信失敗',
    templates: 'クイックテンプレート',
    templateButtons: {
      welcome: 'ウェルカムテンプレート',
      warning: '警告テンプレート',
      notice: '通知テンプレート',
      custom: 'カスタムテンプレート',
    },
    errors: {
      subjectRequired: '連絡件名は空にできません',
      subjectTooLong: '連絡件名の長さは100文字を超えることはできません',
      contentRequired: '連絡コンテンツは空にできません',
      contentTooLong: '連絡コンテンツの長さは5000文字を超えることはできません',
    },
    messages: {
      sent: '連絡送信リクエストを提出しました',
      failed: '提出失敗、もう一度お試しください',
    },
    templatesContent: {
      welcome: {
        subject: '当プラットフォームへようこそ',
        content: `ユーザー様へ

当画像ホスティングプラットフォームへようこそ！

当サービスをお選びいただき、ありがとうございます。使用中に問題が発生した場合は、お気軽にカスタマーサービスチームまでお問い合わせください。

よろしくお願いいたします。

敬具
管理チーム`,
      },
      warning: {
        subject: 'アカウントセキュリティリマインダー',
        content: `ユーザー様へ

アカウントに注意が必要なセキュリティ問題があることを確認しました:

1. パスワードが十分に安全であることを確認してください
2. アカウント情報を他人と共有しないでください
3. 異常なアクティビティが検出された場合は、すぐにご連絡ください

ご質問がある場合は、カスタマーサービスまでお問い合わせください。

敬具
セキュリティチーム`,
      },
      notice: {
        subject: 'システム通知',
        content: `ユーザー様へ

以下の重要な情報をお知らせします:

[ここに具体的な通知内容を記入してください]

ご注意いただき、ありがとうございます。

敬具
管理チーム`,
      },
      custom: {
        content: 'ユーザー様へ\n\n[ここにメール内容を記入してください]\n\n敬具\n管理チーム',
      },
    },
  },

  empty: {
    noUsers: 'ノードなし',
    noSearchResults: '条件に一致するノードが見つかりませんでした',
    createFirst: '最初のノードを作成',
    title: '条件に一致するノードが見つかりませんでした',
    description: 'フィルター条件を調整するか新しいノードを作成してください',
    resetFilter: 'フィルター条件をリセット',
  },

  pagination: {
    total: '合計{total}件のノード',
    pageSize: '{size}件/ページ',
  },

  statusText: {
    normal: 'オンライン',
    disabled: '無効',
    deleted: '削除済み',
    unknown: '不明',
  },

  loading: {
    users: 'ノードデータを読み込み中...',
  },

  disableWarning: {
    title: 'ノード無効化確認',
    warning: 'ノード "{username}" を無効にしてもよろしいですか？',
    consequence1: '無効化後、このノードは即座にすべての操作権限を失います',
    consequence2: 'ノードが現在オンラインの場合、システムから強制退出されます',
    consequence3: 'ノードは再有効化を待ってから再度ログインする必要があります',
    confirmButton: '無効化を確認',
  },

  editDialog: {
    title: 'ノードを編集',
    saveChanges: '変更を保存',
  },

  export: {
    id: 'ノードID',
    username: 'ノード名',
    email: '連絡先アドレス',
    status: 'ステータス',
    role: '権限',
    createdAt: 'アクセス時間',
    updatedAt: '最終更新',
    filename: 'ノードデータ',
  },
}
