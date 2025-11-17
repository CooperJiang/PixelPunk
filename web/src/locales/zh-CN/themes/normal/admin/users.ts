export const users = {
  title: '用户管理',
  subtitle: '管理和监控所有用户账号',

  columns: {
    user: '用户',
    email: '电子邮箱',
    registeredAt: '注册时间',
    lastActivity: '最后操作',
    status: '状态',
    role: '角色',
    storage: '存储/流量',
    actions: '操作',
    userId: 'ID',
    neverActive: '从未操作',
    storageLabel: '存储',
    bandwidthLabel: '流量',
  },

  status: {
    all: '全部',
    active: '正常',
    disabled: '已禁用',
    suspended: '已暂停',
  },

  role: {
    all: '全部',
    superAdmin: '超级管理员',
    admin: '管理员',
    user: '普通用户',
    guest: '游客',
    unknown: '未知',
  },

  roles: {
    all: '全部',
    superAdmin: '超级管理员',
    admin: '管理员',
    normal: '普通用户',
    guest: '游客',
    unknown: '未知',
  },

  actions: {
    view: '查看详情',
    edit: '编辑用户',
    disable: '禁用用户',
    enable: '启用用户',
    delete: '删除用户',
    export: '导出',
    import: '导入',
    refresh: '刷新',
  },

  buttons: {
    createUser: '新增用户',
    create: '创建',
    batchOperation: '批量操作',
    filter: '筛选',
    resetFilter: '重置筛选',
    applyFilter: '应用筛选',
    cancel: '取消',
    confirm: '确认',
    save: '保存',
    close: '关闭',
  },

  toast: {
    userCreated: '用户创建成功',
    userUpdated: '用户信息已更新',
    userDeleted: '用户删除成功',
    statusChanged: '状态已更改',
    batchOperationSuccess: '批量操作完成',
    operationFailed: '操作失败',
    pleaseSelectUsers: '请选择要操作的用户',
    userEnabled: '用户已启用',
    userDisabled: '用户已禁用',
    userDataNotFound: '未找到用户数据',
    success: '{operation}成功',
    error: '{operation}失败',
    exportSuccess: '用户数据导出成功',
    exportError: '用户数据导出失败',
  },

  confirm: {
    enableUser: '确定要启用此用户吗？',
    disableUser: '确定要禁用此用户吗？',
    deleteUser: '确定要删除此用户吗？',
    delete: '确定要删除选中的用户吗？',
    batchDelete: '确定要删除选中的 {count} 个用户吗？',
    batchDisable: '确定要禁用选中的用户吗？',
    batchEnable: '确定要启用选中的 {count} 个用户吗？',
    disableTitle: '即将禁用 {count} 个用户',
    disableMessage1: '禁用后，这些用户将无法登录系统',
    disableMessage2: '当前在线的用户会话将立即终止',
    disableMessage3: '所有相关的分享链接将失效',
  },

  disableDialog: {
    title: '禁用用户确认',
    warning: '您即将禁用用户',
    consequences: '此操作将导致:',
    consequence1: '用户无法登录系统',
    consequence2: '用户的所有会话将被立即终止',
    consequence3: '用户的所有分享链接将失效',
    consequence4: '用户无法访问其文件和数据',
    inputLabel: '请输入用户名以确认',
    inputPlaceholder: '输入用户名',
    inputHint: '请输入 "{username}" 以确认禁用',
    usernameNotMatch: '用户名不匹配',
    confirmButton: '确认禁用',
    cancelButton: '取消',
  },

  filter: {
    title: '筛选用户',
    keyword: '关键词',
    keywordPlaceholder: '搜索用户名或邮箱',
    status: {
      label: '状态',
      all: '全部',
      normal: '正常',
      disabled: '已禁用',
      deleted: '已删除',
    },
    role: {
      label: '角色',
      all: '全部',
      superAdmin: '超级管理员',
      admin: '管理员',
      user: '普通用户',
    },
    registrationTime: '注册时间',
    registrationDate: '注册日期',
    dateRangePlaceholder: '选择日期范围',
    startDatePlaceholder: '开始日期',
    endDatePlaceholder: '结束日期',
    startDate: '开始日期',
    endDate: '结束日期',
    sortBy: '排序方式',
    sort: {
      label: '排序方式',
      placeholder: '选择排序方式',
      default: '默认',
      newest: '最新注册',
      oldest: '最早注册',
      lastActivityNewest: '最近活跃',
      lastActivityOldest: '最早活跃',
      usernameAZ: '用户名A-Z',
      usernameZA: '用户名Z-A',
    },
  },

  batch: {
    title: '批量操作',
    selected: '已选择 {count} 个用户',
    selectAll: '全选',
    deselectAll: '取消全选',
    clearSelection: '清除选择',
    enable: '启用',
    disable: '禁用',
    delete: '删除',
    export: '导出',
    sendEmail: '发送邮件',
    role: '设置角色',
    apply: '应用',
    cancel: '取消',
    noSelection: '请选择要操作的用户',
    setRole: '设置角色',
    selectRole: '请选择角色',
    selectUsersToExport: '请选择要导出的用户',
    confirmDisable: '确认禁用',
  },

  form: {
    createTitle: '创建新用户',
    editTitle: '编辑用户',
    username: '用户名',
    usernamePlaceholder: '请输入用户名',
    usernameHint: '3-20个字符，只能包含字母、数字和下划线',
    email: '电子邮箱',
    emailPlaceholder: '请输入邮箱地址',
    emailHint: '用于登录和接收通知',
    password: '密码',
    passwordPlaceholder: '请输入密码',
    passwordHint: '至少8个字符',
    confirmPassword: '确认密码',
    confirmPasswordPlaceholder: '请再次输入密码',
    role: '角色',
    rolePlaceholder: '选择用户角色',
    roleUser: '普通用户',
    roleAdmin: '管理员',
    roleSuperAdmin: '超级管理员',
    status: '状态',
    statusActive: '正常',
    statusDisabled: '已禁用',
    storageConfig: '存储配置',
    storageLimit: '存储限制',
    storageLimitPlaceholder: '请输入存储限制',
    storageLimitDefault: '默认5GB',
    storageLimitHint: '单位: GB, 0表示无限制',
    bandwidthLimit: '流量限制',
    bandwidthLimitPlaceholder: '请输入流量限制',
    bandwidthLimitDefault: '默认100GB',
    bandwidthLimitHint: '单位: GB, 0表示无限制',
    avatar: '头像',
    avatarUpload: '上传头像',
    usernameRequired: '用户名不能为空',
    usernameMinLength: '用户名至少2个字符',
    usernameMaxLength: '用户名不能超过20个字符',
    emailRequired: '邮箱不能为空',
    emailInvalid: '邮箱格式不正确',
    passwordRequired: '密码不能为空',
    passwordMinLength: '密码至少6个字符',
    passwordMaxLength: '密码不能超过20个字符',
    required: '此项为必填',
    invalidEmail: '邮箱格式不正确',
    passwordTooShort: '密码至少需要8个字符',
    passwordNotMatch: '两次输入的密码不一致',
  },

  detail: {
    title: '用户详情',
    loading: '加载中...',
    noData: '暂无数据',
    basicInfo: '基本信息',
    accountInfo: '账号信息',
    statisticsInfo: '统计信息',
    activityInfo: '活动信息',
    usageStats: '使用统计',
    storageManagement: '存储管理',
    quickActions: '快捷操作',
    userId: '用户ID',
    username: '用户名',
    email: '邮箱',
    role: '角色',
    registeredAt: '注册时间',
    lastLogin: '最后登录',
    lastLoginAt: '最后登录',
    lastActivity: '最后活动',
    lastActivityAt: '最后活动',
    lastActivityIp: 'IP地址',
    uploadedFiles: '上传文件数',
    shareCount: '分享数量',
    totalViews: '总浏览量',
    usedStorage: '已用存储',
    storageLimit: '存储限制',
    monthlyBandwidth: '月流量',
    bandwidthLimit: '流量限制',
    fileCount: '文件数量',
    unlimited: '无限制',
    never: '从未',
    currentSetting: '当前设置',
    enterValue: '输入数值',
    unit: '单位',
    updateStorageSettings: '更新存储设置',
    status: {
      label: '状态',
      active: '正常',
      disabled: '已禁用',
      deleted: '已删除',
      unknown: '未知',
    },
    roles: {
      superAdmin: '超级管理员',
      admin: '管理员',
      user: '普通用户',
      unknown: '未知',
    },
    units: {
      mb: 'MB',
      gb: 'GB',
      tb: 'TB',
    },
    actions: {
      resetPassword: '重置密码',
      resetPasswordDesc: '为该用户生成新密码',
      sendEmail: '发送邮件',
      sendEmailDesc: '向该用户发送邮件',
      enableUser: '启用用户',
      enableUserDesc: '允许该用户登录',
      disableUser: '禁用用户',
      disableUserDesc: '阻止该用户登录',
    },
    messages: {
      storageUpdateSuccess: '存储设置已更新',
      passwordResetSuccess: '密码已重置为: {password}',
      emailSentSuccess: '邮件发送成功',
      userEnabled: '用户已启用',
      userDisabled: '用户已禁用',
    },
    errors: {
      fetchFailed: '获取用户信息失败',
      storageUpdateFailed: '更新存储设置失败',
      passwordResetFailed: '重置密码失败',
      emailSentFailed: '发送邮件失败',
      operationFailed: '操作失败',
    },
  },

  search: {
    placeholder: '搜索用户名或邮箱...',
    button: '搜索',
  },

  sendEmail: {
    title: '发送邮件',
    recipients: '收件人',
    recipientsHint: '已选择 {count} 个用户',
    subject: '主题',
    subjectPlaceholder: '请输入邮件主题',
    content: '内容',
    contentPlaceholder: '请输入邮件内容',
    send: '发送',
    cancel: '取消',
    sending: '发送中...',
    sent: '邮件已发送',
    failed: '发送失败',
    templates: '快速模板',
    templateButtons: {
      welcome: '欢迎模板',
      warning: '警告模板',
      notice: '通知模板',
      custom: '自定义模板',
    },
    errors: {
      subjectRequired: '邮件主题不能为空',
      subjectTooLong: '邮件主题长度不能超过100个字符',
      contentRequired: '邮件内容不能为空',
      contentTooLong: '邮件内容长度不能超过5000个字符',
    },
    messages: {
      sent: '邮件发送请求已提交',
      failed: '提交失败，请重试',
    },
    templatesContent: {
      welcome: {
        subject: '欢迎加入我们的平台',
        content: `尊敬的用户，

欢迎加入我们的图床平台！

我们很高兴您选择了我们的服务。如果您在使用过程中遇到任何问题，请随时联系我们的客服团队。

祝您使用愉快！

此致
管理团队`,
      },
      warning: {
        subject: '账户安全提醒',
        content: `尊敬的用户，

我们注意到您的账户存在一些需要注意的安全问题，请及时处理：

1. 请确保您的密码足够安全
2. 不要与他人分享您的账户信息
3. 如发现异常活动，请立即联系我们

如有疑问，请联系客服。

此致
安全团队`,
      },
      notice: {
        subject: '系统通知',
        content: `尊敬的用户，

我们想要通知您以下重要信息：

[请在此处填写具体通知内容]

感谢您的关注。

此致
管理团队`,
      },
      custom: {
        content: '尊敬的用户，\n\n[请在此处填写您的邮件内容]\n\n此致\n管理团队',
      },
    },
  },

  empty: {
    noUsers: '暂无用户',
    noSearchResults: '没有找到符合条件的用户',
    createFirst: '创建第一个用户',
    title: '没有找到符合条件的用户',
    description: '尝试调整筛选条件或创建新用户',
    resetFilter: '重置筛选条件',
  },

  pagination: {
    total: '共 {total} 个用户',
    pageSize: '每页 {size} 条',
  },

  statusText: {
    normal: '正常',
    disabled: '已禁用',
    deleted: '已删除',
    unknown: '未知',
  },

  loading: {
    users: '正在加载用户数据...',
  },

  disableWarning: {
    title: '禁用用户确认',
    warning: '您确定要禁用用户「{username}」吗？',
    consequence1: '禁用后，该用户将立即失去所有操作权限',
    consequence2: '如果用户当前在线，将被强制退出系统',
    consequence3: '用户需要等待重新启用后才能再次登录',
    confirmButton: '确认禁用',
  },

  editDialog: {
    title: '编辑用户',
    saveChanges: '保存更改',
  },

  export: {
    id: '用户ID',
    username: '用户名',
    email: '邮箱',
    status: '状态',
    role: '角色',
    createdAt: '注册时间',
    updatedAt: '最后更新',
    filename: '用户数据',
  },
}
