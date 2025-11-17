export const users = {
  title: 'ユーザー管理',
  subtitle: 'すべてのユーザーアカウントを管理・監視',

  columns: {
    user: 'ユーザー',
    email: 'メールアドレス',
    registeredAt: '登録日',
    lastActivity: '最終アクティビティ',
    status: 'ステータス',
    role: 'ロール',
    storage: 'ストレージ/帯域幅',
    actions: '操作',
    userId: 'ID',
    neverActive: 'なし',
    storageLabel: 'ストレージ',
    bandwidthLabel: '帯域幅',
  },

  status: {
    all: 'すべて',
    active: 'アクティブ',
    disabled: '無効',
    suspended: '停止',
  },

  role: {
    all: 'すべて',
    superAdmin: 'スーパー管理者',
    admin: '管理者',
    user: 'ユーザー',
    guest: 'ゲスト',
    unknown: '不明',
  },

  roles: {
    all: 'すべて',
    superAdmin: 'スーパー管理者',
    admin: '管理者',
    normal: 'ユーザー',
    guest: 'ゲスト',
    unknown: '不明',
  },

  actions: {
    view: '表示',
    edit: '編集',
    disable: '無効化',
    enable: '有効化',
    delete: '削除',
    export: 'エクスポート',
    import: 'インポート',
    refresh: '更新',
  },

  buttons: {
    createUser: '新規ユーザー',
    create: '作成',
    batchOperation: '一括',
    filter: 'フィルター',
    resetFilter: 'リセット',
    applyFilter: '適用',
    cancel: 'キャンセル',
    confirm: '確認',
    save: '保存',
    close: '閉じる',
  },

  toast: {
    userCreated: 'ユーザー作成済み',
    userUpdated: 'ユーザー更新済み',
    userDeleted: 'ユーザー削除済み',
    statusChanged: 'ステータス変更済み',
    batchOperationSuccess: '一括操作完了',
    operationFailed: '操作失敗',
    pleaseSelectUsers: 'ユーザーを選択してください',
    userEnabled: 'ユーザー有効化済み',
    userDisabled: 'ユーザー無効化済み',
    userDataNotFound: 'ユーザーが見つかりません',
    success: '{operation}成功',
    error: '{operation}失敗',
    exportSuccess: 'ユーザーデータエクスポート済み',
    exportError: 'ユーザーデータエクスポート失敗',
  },

  confirm: {
    enableUser: 'このユーザーを有効化しますか？',
    disableUser: 'このユーザーを無効化しますか？',
    deleteUser: 'このユーザーを削除しますか？',
    delete: '選択したユーザーを削除しますか？',
    batchDelete: '選択した{count}人のユーザーを削除しますか？',
    batchDisable: '選択したユーザーを無効化しますか？',
    batchEnable: '選択した{count}人のユーザーを有効化しますか？',
    disableTitle: '{count}人のユーザーを無効化',
    disableMessage1: '無効化されたユーザーはログインできません',
    disableMessage2: 'アクティブなセッションは終了されます',
    disableMessage3: '共有リンクは無効になります',
  },

  disableDialog: {
    title: '無効化の確認',
    warning: '無効化しようとしています',
    consequences: 'これにより：',
    consequence1: 'ログインが阻止されます',
    consequence2: 'セッションが終了されます',
    consequence3: '共有リンクが無効になります',
    consequence4: 'ファイル/データへのアクセスがブロックされます',
    inputLabel: '確認のためユーザー名を入力',
    inputPlaceholder: 'ユーザー名を入力',
    inputHint: '確認のため「{username}」と入力してください',
    usernameNotMatch: 'ユーザー名が一致しません',
    confirmButton: '無効化',
    cancelButton: 'キャンセル',
  },

  filter: {
    title: 'ユーザーフィルター',
    keyword: 'キーワード',
    keywordPlaceholder: 'ユーザー名またはメールアドレスで検索',
    status: {
      label: 'ステータス',
      all: 'すべて',
      normal: 'アクティブ',
      disabled: '無効',
      deleted: '削除済み',
    },
    role: {
      label: 'ロール',
      all: 'すべて',
      superAdmin: 'スーパー管理者',
      admin: '管理者',
      user: 'ユーザー',
    },
    registrationTime: '登録',
    registrationDate: '登録日',
    dateRangePlaceholder: '期間を選択',
    startDatePlaceholder: '開始日',
    endDatePlaceholder: '終了日',
    startDate: '開始日',
    endDate: '終了日',
    sortBy: '並び替え',
    sort: {
      label: '並び替え',
      placeholder: '並び替えを選択',
      default: 'デフォルト',
      newest: '最新',
      oldest: '最古',
      lastActivityNewest: '最終アクティビティ（最新）',
      lastActivityOldest: '最終アクティビティ（最古）',
      usernameAZ: 'ユーザー名A-Z',
      usernameZA: 'ユーザー名Z-A',
    },
  },

  batch: {
    title: '一括',
    selected: '{count}件選択',
    selectAll: 'すべて選択',
    deselectAll: '選択解除',
    clearSelection: '選択クリア',
    enable: '有効化',
    disable: '無効化',
    delete: '削除',
    export: 'エクスポート',
    sendEmail: 'メール送信',
    role: 'ロール設定',
    apply: '適用',
    cancel: 'キャンセル',
    noSelection: 'ユーザーを選択してください',
    setRole: 'ロール設定',
    selectRole: 'ロールを選択',
    selectUsersToExport: 'エクスポートするユーザーを選択',
    confirmDisable: '無効化の確認',
  },

  form: {
    createTitle: 'ユーザー作成',
    editTitle: 'ユーザー編集',
    username: 'ユーザー名',
    usernamePlaceholder: 'ユーザー名を入力',
    usernameHint: '3-20文字；英数字、アンダースコア',
    email: 'メールアドレス',
    emailPlaceholder: 'メールアドレスを入力',
    emailHint: 'ログインと通知に使用',
    password: 'パスワード',
    passwordPlaceholder: 'パスワードを入力',
    passwordHint: '最低8文字',
    confirmPassword: 'パスワード確認',
    confirmPasswordPlaceholder: 'パスワードを再入力',
    role: 'ロール',
    rolePlaceholder: 'ロールを選択',
    roleUser: 'ユーザー',
    roleAdmin: '管理者',
    roleSuperAdmin: 'スーパー管理者',
    status: 'ステータス',
    statusActive: 'アクティブ',
    statusDisabled: '無効',
    storageConfig: 'ストレージ',
    storageLimit: 'ストレージ制限',
    storageLimitPlaceholder: '制限を入力',
    storageLimitDefault: 'デフォルト5GB',
    storageLimitHint: '単位：GB。0は無制限',
    bandwidthLimit: '帯域幅制限',
    bandwidthLimitPlaceholder: '制限を入力',
    bandwidthLimitDefault: 'デフォルト100GB',
    bandwidthLimitHint: '単位：GB。0は無制限',
    avatar: 'アバター',
    avatarUpload: 'アバターアップロード',
    usernameRequired: 'ユーザー名は必須です',
    usernameMinLength: '最低2文字',
    usernameMaxLength: '20文字以内',
    emailRequired: 'メールアドレスは必須です',
    emailInvalid: '無効なメールアドレス',
    passwordRequired: 'パスワードは必須です',
    passwordMinLength: '最低6文字',
    passwordMaxLength: '20文字以内',
    required: 'この項目は必須です',
    invalidEmail: '無効なメールアドレス',
    passwordTooShort: '最低8文字',
    passwordNotMatch: 'パスワードが一致しません',
  },

  detail: {
    title: 'ユーザー詳細',
    loading: '読み込み中...',
    noData: 'データがありません',
    basicInfo: '基本情報',
    accountInfo: 'アカウント情報',
    statisticsInfo: '統計',
    activityInfo: 'アクティビティ',
    usageStats: '使用統計',
    storageManagement: 'ストレージ管理',
    quickActions: 'クイック操作',
    userId: 'ユーザーID',
    username: 'ユーザー名',
    email: 'メールアドレス',
    role: 'ロール',
    registeredAt: '登録日',
    lastLogin: '最終ログイン',
    lastLoginAt: '最終ログイン',
    lastActivity: '最終アクティビティ',
    lastActivityAt: '最終アクティビティ',
    lastActivityIp: 'IPアドレス',
    uploadedFiles: 'アップロード',
    shareCount: '共有',
    totalViews: '総閲覧数',
    usedStorage: '使用ストレージ',
    storageLimit: 'ストレージ制限',
    monthlyBandwidth: '月間帯域幅',
    bandwidthLimit: '帯域幅制限',
    fileCount: 'ファイル',
    unlimited: '無制限',
    never: 'なし',
    currentSetting: '現在',
    enterValue: '値を入力',
    unit: '単位',
    updateStorageSettings: 'ストレージ更新',
    status: {
      label: 'ステータス',
      active: 'アクティブ',
      disabled: '無効',
      deleted: '削除済み',
      unknown: '不明',
    },
    roles: {
      superAdmin: 'スーパー管理者',
      admin: '管理者',
      user: 'ユーザー',
      unknown: '不明',
    },
    units: {
      mb: 'MB',
      gb: 'GB',
      tb: 'TB',
    },
    actions: {
      resetPassword: 'パスワードリセット',
      resetPasswordDesc: '新しいパスワードを生成',
      sendEmail: 'メール送信',
      sendEmailDesc: 'ユーザーにメールを送信',
      enableUser: 'ユーザー有効化',
      enableUserDesc: 'このユーザーのログインを許可',
      disableUser: 'ユーザー無効化',
      disableUserDesc: 'このユーザーのログインを阻止',
    },
    messages: {
      storageUpdateSuccess: 'ストレージ更新済み',
      passwordResetSuccess: 'パスワードリセット：{password}',
      emailSentSuccess: 'メール送信済み',
      userEnabled: 'ユーザー有効化済み',
      userDisabled: 'ユーザー無効化済み',
    },
    errors: {
      fetchFailed: 'ユーザー取得失敗',
      storageUpdateFailed: 'ストレージ更新失敗',
      passwordResetFailed: 'パスワードリセット失敗',
      emailSentFailed: 'メール送信失敗',
      operationFailed: '操作失敗',
    },
  },

  search: {
    placeholder: 'ユーザー名またはメールアドレスで検索...',
    button: '検索',
  },

  sendEmail: {
    title: 'メール送信',
    recipients: '宛先',
    recipientsHint: '{count}人のユーザーが選択されています',
    subject: '件名',
    subjectPlaceholder: '件名を入力',
    content: '内容',
    contentPlaceholder: '内容を入力',
    send: '送信',
    cancel: 'キャンセル',
    sending: '送信中...',
    sent: 'メール送信済み',
    failed: '失敗',
    templates: 'テンプレート',
    templateButtons: {
      welcome: 'ウェルカム',
      warning: '警告',
      notice: 'お知らせ',
      custom: 'カスタム',
    },
    errors: {
      subjectRequired: '件名は必須です',
      subjectTooLong: '件名は100文字以内',
      contentRequired: '内容は必須です',
      contentTooLong: '内容は5000文字以内',
    },
    messages: {
      sent: 'メールリクエスト送信済み',
      failed: '送信失敗、再試行してください',
    },
    templatesContent: {
      welcome: {
        subject: 'プラットフォームへようこそ',
        content: `ユーザー様へ

画像ホスティングプラットフォームへようこそ！

当サービスをご利用いただきありがとうございます。ご不明な点がございましたら、サポートチームまでお気軽にお問い合わせください。

敬具
運営チーム`,
      },
      warning: {
        subject: 'アカウントセキュリティのご案内',
        content: `ユーザー様へ

アカウントにセキュリティ上の問題が検出されました。以下をご確認ください：

1. 強力なパスワードを使用する
2. アカウント情報を共有しない
3. 不審な活動を検出した場合はお問い合わせください

敬具
セキュリティチーム`,
      },
      notice: {
        subject: 'システムお知らせ',
        content: `ユーザー様へ

重要な更新情報をお知らせします：

[詳細をここに挿入]

ありがとうございます。
運営チーム`,
      },
      custom: {
        content: 'ユーザー様へ\n\n[ここにメッセージを記入]\n\n敬具\n運営チーム',
      },
    },
  },

  empty: {
    noUsers: 'ユーザーがありません',
    noSearchResults: '一致するユーザーが見つかりません',
    createFirst: '最初のユーザーを作成',
    title: '一致するユーザーが見つかりません',
    description: 'フィルターを調整するか、新しいユーザーを作成してください',
    resetFilter: 'フィルターリセット',
  },

  pagination: {
    total: '合計{total}人のユーザー',
    pageSize: '{size}件/ページ',
  },

  statusText: {
    normal: 'アクティブ',
    disabled: '無効',
    deleted: '削除済み',
    unknown: '不明',
  },

  loading: {
    users: 'ユーザーを読み込み中...',
  },

  disableWarning: {
    title: '無効化の確認',
    warning: 'ユーザー「{username}」を無効化しますか？',
    consequence1: 'ユーザーは即座に権限を失います',
    consequence2: 'アクティブなセッションはログアウトされます',
    consequence3: '再有効化されるまでログインできません',
    confirmButton: '無効化',
  },

  editDialog: {
    title: 'ユーザー編集',
    saveChanges: '変更を保存',
  },

  export: {
    id: 'ユーザーID',
    username: 'ユーザー名',
    email: 'メールアドレス',
    status: 'ステータス',
    role: 'ロール',
    createdAt: '登録日',
    updatedAt: '更新日',
    filename: 'users',
  },
}
