/**
 * 系统安装页面文案 - 赛博风格
 */
export const setup = {
  loading: 'SYSTEM INITIALIZING...',
  title: 'PIXELPUNK SYSTEM INSTALLATION',
  subtitle: '配置您的智能图片管理平台',
  success: {
    title: '系统部署成功',
    subtitle: '欢迎接入 PixelPunk 智能资源管理系统',
    entries: {
      home: {
        title: '进入首页',
        desc: '基础入口与数据展示',
      },
      upload: {
        title: '数据注入',
        desc: '开始注入你的第一个数据单元',
      },
      dashboard: {
        title: '控制面板',
        desc: '查看系统统计与数据',
      },
      settings: {
        title: '系统配置',
        desc: '自定义你的图床设置',
      },
    },
    tip: '选择一个快捷入口开始使用系统',
  },
  database: {
    title: 'DATABASE',
    mysql: 'MySQL',
    sqlite: 'SQLite',
    fields: {
      path: 'PATH',
      host: 'HOST',
      port: 'PORT',
      user: 'USER',
      password: 'PASSWORD',
      name: 'DATABASE NAME',
    },
    placeholders: {
      path: './pixelpunk.db',
      host: '127.0.0.1',
      port: '3306',
      user: 'root',
      password: 'mysql password',
      name: 'pixelpunk',
    },
    validation: {
      mysqlRequired: '请完善MySQL连接信息',
      sqliteRequired: '请指定SQLite文件路径',
      connectionFailed: '连接失败',
      unknownError: '未知错误',
    },
    test: {
      success: {
        sqlite: 'SQLite就绪',
        mysql: 'MySQL连接成功',
      },
      testing: '测试中...',
      button: '数据库配置测试',
    },
  },
  admin: {
    title: 'ADMIN ACCOUNT',
    fields: {
      username: 'USERNAME',
      password: 'PASSWORD',
    },
    placeholders: {
      username: 'root',
      password: '最少6位数',
    },
    validation: {
      mysqlInfoRequired: '请填写完整的MySQL信息',
      sqlitePathRequired: '请填写SQLite路径',
      usernameRequired: '请输入管理员用户名（至少3位）',
      passwordRequired: '请输入管理员密码',
      passwordTooShort: '管理员密码至少6位',
    },
  },
  redis: {
    title: 'REDIS',
    optional: 'OPTIONAL',
    fields: {
      host: 'HOST',
      port: 'PORT',
      password: 'PASSWORD',
      db: 'DB',
    },
    placeholders: {
      host: 'localhost',
      port: '6379',
      password: '(empty)',
      db: '0',
    },
  },
  vector: {
    title: 'QDRANT VECTOR DB',
    optional: 'OPTIONAL',
    hint: '向量数据库用于支持以图搜图和语义搜索功能，不配置不影响其他功能使用',
    mode: {
      builtin: {
        title: '使用内置',
        desc: '系统将自动启动和管理 Qdrant',
        badge: '推荐',
      },
      external: {
        title: '使用外部',
        desc: '连接已存在的 Qdrant 服务',
      },
    },
    advanced: {
      title: '⚙️ 高级配置',
      hint: '配置 Qdrant 端口（如端口被占用请修改）',
      httpPort: {
        label: 'HTTP 端口',
        desc: '用于 API 调用（主要使用）',
      },
      grpcPort: {
        label: 'gRPC 端口',
        desc: 'Qdrant 内部使用（虽然我们不直接使用，但 Qdrant 需要此端口才能运行）',
      },
    },
    fields: {
      url: 'QDRANT URL',
      timeout: 'TIMEOUT (秒)',
    },
    placeholders: {
      url: 'http://your-server:6333',
      timeout: '30',
      httpPort: '6333',
      grpcPort: '6334',
    },
    validation: {
      urlRequired: '请填写Qdrant URL',
      invalidUrl: '无效的URL格式',
      urlInvalid: 'URL 无效',
      formatError: '格式错误',
    },
    test: {
      success: {
        builtin: 'URL 格式正确，安装时将自动启动内置 Qdrant',
        external: 'URL 格式正确（实际连接需要在安装时验证）',
      },
      testing: '测试中...',
      button: {
        builtin: '验证配置',
        external: '测试连接',
      },
    },
  },
  actions: {
    install: '系统部署',
    installing: '部署中...',
  },
}
