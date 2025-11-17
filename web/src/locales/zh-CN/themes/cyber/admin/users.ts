export const users = {
  title: '节点管理',
  subtitle: '管理和监控所有接入节点',

  columns: {
    user: '节点',
    email: '通信地址',
    registeredAt: '接入时间',
    lastActivity: '最后活动',
    status: '状态',
    role: '权限',
    storage: '存储/带宽',
    actions: '操作',
    userId: 'ID',
    neverActive: '从未活动',
    storageLabel: '存储',
    bandwidthLabel: '带宽',
  },

  status: {
    all: '全部',
    active: '在线',
    disabled: '已禁用',
    suspended: '已暂停',
  },

  role: {
    all: '全部',
    superAdmin: '超级管理员',
    admin: '管理节点',
    user: '普通节点',
    guest: '访客节点',
    unknown: '未知',
  },

  roles: {
    all: '全部',
    superAdmin: '超级管理员',
    admin: '管理节点',
    normal: '普通节点',
    guest: '访客节点',
    unknown: '未知',
  },

  actions: {
    view: '查看详情',
    edit: '编辑节点',
    disable: '禁用节点',
    enable: '启用节点',
    delete: '删除节点',
    export: '导出',
    import: '导入',
    refresh: '刷新',
  },

  buttons: {
    createUser: '新增节点',
    create: '创建节点',
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
    userCreated: '节点创建成功',
    userUpdated: '节点信息已更新',
    userDeleted: '节点删除成功',
    statusChanged: '状态已更改',
    batchOperationSuccess: '批量操作完成',
    operationFailed: '操作失败',
    pleaseSelectUsers: '请选择要操作的节点',
    userEnabled: '节点已启用',
    userDisabled: '节点已禁用',
    userDataNotFound: '未找到节点数据',
    success: '{operation}成功',
    error: '{operation}失败',
    exportSuccess: '节点数据导出成功',
    exportError: '节点数据导出失败',
  },

  confirm: {
    enableUser: '确定要启用此节点吗？',
    disableUser: '确定要禁用此节点吗？',
    deleteUser: '确定要删除此节点吗？',
    delete: '确定要删除选中的节点吗？',
    batchDelete: '确定要删除选中的 {count} 个节点吗？',
    batchDisable: '确定要禁用选中的节点吗？',
    batchEnable: '确定要启用选中的 {count} 个节点吗？',
    disableTitle: '即将禁用 {count} 个节点',
    disableMessage1: '禁用后，这些节点将无法接入系统',
    disableMessage2: '当前在线的节点会话将立即终止',
    disableMessage3: '所有相关的协议链接将失效',
  },

  disableDialog: {
    title: '禁用节点确认',
    warning: '您即将禁用节点',
    consequences: '此操作将导致:',
    consequence1: '节点无法登录系统',
    consequence2: '节点的所有会话将被立即终止',
    consequence3: '节点的所有协议链接将失效',
    consequence4: '节点无法访问其数据',
    inputLabel: '请输入节点名以确认',
    inputPlaceholder: '输入节点名',
    inputHint: '请输入 "{username}" 以确认禁用',
    usernameNotMatch: '节点名不匹配',
    confirmButton: '确认禁用',
    cancelButton: '取消',
  },

  filter: {
    title: '筛选节点',
    keyword: '关键词',
    keywordPlaceholder: '搜索节点名或通信地址',
    status: {
      label: '节点状态',
      all: '全部节点',
      normal: '正常节点',
      disabled: '已禁用节点',
      deleted: '已删除节点',
    },
    role: {
      label: '节点权限',
      all: '全部权限',
      superAdmin: '超级管理节点',
      admin: '管理节点',
      user: '普通节点',
    },
    registrationTime: '接入时间',
    registrationDate: '接入日期',
    dateRangePlaceholder: '选择时间范围',
    startDatePlaceholder: '开始日期',
    endDatePlaceholder: '结束日期',
    startDate: '开始日期',
    endDate: '结束日期',
    sortBy: '排序方式',
    sort: {
      label: '排序方式',
      placeholder: '选择排序方式',
      default: '默认',
      newest: '最新接入',
      oldest: '最早接入',
      lastActivityNewest: '最近活跃',
      lastActivityOldest: '最早活跃',
      usernameAZ: '节点名A-Z',
      usernameZA: '节点名Z-A',
    },
  },

  batch: {
    title: '批量操作',
    selected: '已选择 {count} 个节点',
    selectAll: '全选',
    deselectAll: '取消全选',
    clearSelection: '清除选择',
    enable: '启用',
    disable: '禁用',
    delete: '删除',
    export: '导出',
    sendEmail: '发送通信',
    role: '设置权限',
    apply: '应用',
    cancel: '取消',
    noSelection: '请选择要操作的节点',
    setRole: '设置权限',
    selectRole: '请选择权限',
    selectUsersToExport: '请选择要导出的节点',
    confirmDisable: '确认禁用',
  },

  form: {
    createTitle: '创建新节点',
    editTitle: '编辑节点',
    username: '节点名',
    usernamePlaceholder: '请输入节点名',
    usernameHint: '3-20个字符，只能包含字母、数字和下划线',
    email: '通信地址',
    emailPlaceholder: '请输入通信地址',
    emailHint: '用于登录和接收通知',
    password: '密码',
    passwordPlaceholder: '请输入密码',
    passwordHint: '至少8个字符',
    confirmPassword: '确认密码',
    confirmPasswordPlaceholder: '请再次输入密码',
    role: '权限',
    rolePlaceholder: '选择节点权限',
    roleUser: '普通节点',
    roleAdmin: '管理节点',
    roleSuperAdmin: '超级管理节点',
    status: '状态',
    statusActive: '正常节点',
    statusDisabled: '已禁用节点',
    storageConfig: '存储配置',
    storageLimit: '存储限制',
    storageLimitPlaceholder: '请输入存储限制',
    storageLimitDefault: '默认5GB',
    storageLimitHint: '单位: GB, 0表示无限制',
    bandwidthLimit: '带宽限制',
    bandwidthLimitPlaceholder: '请输入带宽限制',
    bandwidthLimitDefault: '默认100GB',
    bandwidthLimitHint: '单位: GB, 0表示无限制',
    avatar: '头像',
    avatarUpload: '上传头像',
    usernameRequired: '节点名不能为空',
    usernameMinLength: '节点名至少2个字符',
    usernameMaxLength: '节点名不能超过20个字符',
    emailRequired: '通信地址不能为空',
    emailInvalid: '通信地址格式不正确',
    passwordRequired: '密码不能为空',
    passwordMinLength: '密码至少6个字符',
    passwordMaxLength: '密码不能超过20个字符',
    required: '此项为必填',
    invalidEmail: '通信地址格式不正确',
    passwordTooShort: '密码至少需要8个字符',
    passwordNotMatch: '两次输入的密码不一致',
  },

  detail: {
    title: '节点详情',
    loading: '数据加载中...',
    noData: '暂无数据',
    basicInfo: '基本信息',
    accountInfo: '账号信息',
    statisticsInfo: '统计信息',
    activityInfo: '活动信息',
    usageStats: '使用统计',
    storageManagement: '存储管理',
    quickActions: '快捷操作',
    userId: '节点ID',
    username: '节点名',
    email: '通信地址',
    role: '权限',
    registeredAt: '接入时间',
    lastLogin: '最后登录',
    lastLoginAt: '最后登录',
    lastActivity: '最后活动',
    lastActivityAt: '最后活动',
    lastActivityIp: 'IP地址',
    uploadedFiles: '上传数据数',
    shareCount: '协议数量',
    totalViews: '总访问量',
    usedStorage: '已用存储',
    storageLimit: '存储限制',
    monthlyBandwidth: '月带宽',
    bandwidthLimit: '带宽限制',
    fileCount: '数据数量',
    unlimited: '无限制',
    never: '从未',
    currentSetting: '当前配置',
    enterValue: '输入数值',
    unit: '单位',
    updateStorageSettings: '更新存储配置',
    status: {
      label: '状态',
      active: '在线',
      disabled: '已禁用',
      deleted: '已删除',
      unknown: '未知',
    },
    roles: {
      superAdmin: '超级管理节点',
      admin: '管理节点',
      user: '普通节点',
      unknown: '未知',
    },
    units: {
      mb: 'MB',
      gb: 'GB',
      tb: 'TB',
    },
    actions: {
      resetPassword: '重置密码',
      resetPasswordDesc: '为该节点生成新密码',
      sendEmail: '发送邮件',
      sendEmailDesc: '向该节点发送邮件',
      enableUser: '启用节点',
      enableUserDesc: '允许该节点登录',
      disableUser: '禁用节点',
      disableUserDesc: '阻止该节点登录',
    },
    messages: {
      storageUpdateSuccess: '存储配置已更新',
      passwordResetSuccess: '密码已重置为: {password}',
      emailSentSuccess: '邮件发送成功',
      userEnabled: '节点已启用',
      userDisabled: '节点已禁用',
    },
    errors: {
      fetchFailed: '获取节点信息失败',
      storageUpdateFailed: '更新存储配置失败',
      passwordResetFailed: '重置密码失败',
      emailSentFailed: '发送邮件失败',
      operationFailed: '操作失败',
    },
  },

  search: {
    placeholder: '搜索节点名或通信地址...',
    button: '搜索',
  },

  sendEmail: {
    title: '发送通信',
    recipients: '接收节点',
    recipientsHint: '已选择 {count} 个节点',
    subject: '主题',
    subjectPlaceholder: '请输入通信主题',
    content: '内容',
    contentPlaceholder: '请输入通信内容',
    send: '发送',
    cancel: '取消',
    sending: '发送中...',
    sent: '通信已发送',
    failed: '发送失败',
    templates: '快速模板',
    templateButtons: {
      welcome: '欢迎模板',
      warning: '警告模板',
      notice: '通知模板',
      custom: '自定义模板',
    },
    errors: {
      subjectRequired: '通信主题不能为空',
      subjectTooLong: '通信主题长度不能超过100个字符',
      contentRequired: '通信内容不能为空',
      contentTooLong: '通信内容长度不能超过5000个字符',
    },
    messages: {
      sent: '通信发送请求已提交',
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
    noUsers: '暂无节点',
    noSearchResults: '没有找到符合条件的节点',
    createFirst: '创建第一个节点',
    title: '没有找到符合条件的节点',
    description: '尝试调整筛选条件或创建新节点',
    resetFilter: '重置筛选条件',
  },

  pagination: {
    total: '共 {total} 个节点',
    pageSize: '每页 {size} 条',
  },

  statusText: {
    normal: '在线',
    disabled: '已禁用',
    deleted: '已删除',
    unknown: '未知',
  },

  loading: {
    users: '正在加载节点数据...',
  },

  disableWarning: {
    title: '禁用节点确认',
    warning: '您确定要禁用节点「{username}」吗？',
    consequence1: '禁用后，该节点将立即失去所有操作权限',
    consequence2: '如果节点当前在线，将被强制退出系统',
    consequence3: '节点需要等待重新启用后才能再次登录',
    confirmButton: '确认禁用',
  },

  editDialog: {
    title: '编辑节点',
    saveChanges: '保存更改',
  },

  export: {
    id: '节点ID',
    username: '节点名',
    email: '通信地址',
    status: '状态',
    role: '权限',
    createdAt: '接入时间',
    updatedAt: '最后更新',
    filename: '节点数据',
  },
}
