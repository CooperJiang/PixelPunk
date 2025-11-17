/**
 * 开放API页面文案 - 赛博风格
 */
export const openApi = {
  page: {
    title: '开放API管理',
    subtitle: '管理随机图片API接口，为第三方应用提供数据服务',
    create: '创建API',
    search: '检索API名称...',
  },
  loading: '加载中...',
  footer: {
    total: '共 {count} 个API',
  },
  empty: {
    title: '暂无API',
    description: '尚未创建任何API接口，点击上方按钮开始创建',
  },
  toast: {
    loadFailed: '加载失败',
    createSuccess: '创建成功',
    createFailed: '创建失败',
    updateSuccess: '更新成功',
    updateFailed: '更新失败',
    deleteSuccess: '删除成功',
    deleteFailed: '删除失败',
    copySuccess: 'API地址已复制',
    toggleSuccess: '状态已更新',
    toggleFailed: '状态更新失败',
  },
  dialog: {
    create: {
      title: '创建随机图片API节点',
      cancel: '取消',
      submit: '创建',
      form: {
        name: {
          label: 'API节点名称',
          hint: '用于标识这个API节点，方便管理',
        },
        scope: {
          label: '图片范围节点',
        },
        folder: {
          label: '选择文件夹节点',
          hint: '选择一个文件夹节点，API将从该文件夹的公开图片中随机返回',
        },
        returnType: {
          label: '返回方式节点',
          redirect: {
            label: '302重定向节点',
            desc: '返回图片真实地址节点，性能更好，支持CDN加速',
          },
          direct: {
            label: '直接返回图片节点',
            desc: '每次刷新都是新图片节点，URL固定不变，适合网站背景/壁纸轮播',
          },
        },
      },
      hint: '创建后将生成唯一的API链接节点，只返回公开访问权限的图片',
    },
    edit: {
      title: '编辑API配置节点',
      cancel: '取消',
      submit: '保存',
      form: {
        nameLabel: 'API节点名称:',
        scope: {
          label: '图片范围节点',
        },
        folder: {
          label: '选择文件夹节点',
          hint: '选择一个文件夹节点，API将从该文件夹的公开图片中随机返回',
        },
        returnType: {
          label: '返回方式节点',
          redirect: {
            label: '302重定向节点',
            desc: '返回图片真实地址节点，性能更好，支持CDN加速',
          },
          direct: {
            label: '直接返回图片节点',
            desc: '每次刷新都是新图片节点，URL固定不变，适合网站背景/壁纸轮播',
          },
        },
      },
    },
    delete: {
      title: '删除API节点',
      message: '确定要删除API节点「{name}」吗？',
      warning: '删除后无法恢复节点，请谨慎操作',
      cancel: '取消',
      confirm: '确认删除',
    },
  },
  form: {
    name: {
      label: 'API名称',
      placeholder: '请输入API名称',
      required: '请输入API名称',
    },
    description: {
      label: 'API描述',
      placeholder: '请输入API描述（可选）',
    },
    categories: {
      label: '分类筛选',
      placeholder: '选择分类（可多选）',
      all: '全部分类',
    },
    tags: {
      label: '标签筛选',
      placeholder: '选择标签（可多选）',
      all: '全部标签',
    },
    status: {
      label: '状态',
      enabled: '启用',
      disabled: '禁用',
    },
  },
  card: {
    status: {
      enabled: '启用',
      disabled: '禁用',
    },
    actions: {
      open: '打开接口',
      copy: '复制地址',
      edit: '编辑',
      delete: '删除',
      toggle: '切换状态',
    },
    stats: {
      categories: '分类',
      tags: '标签',
      images: '数据单元',
      calls: '调用',
    },
    filters: {
      allCategories: '全部分类',
      allTags: '全部标签',
    },
    lastCall: {
      never: '从未调用',
      justNow: '刚刚',
      minutesAgo: '{minutes}分钟前',
      hoursAgo: '{hours}小时前',
      daysAgo: '{days}天前',
    },
    returnType: {
      direct: '直接返回',
      redirect: '302重定向',
    },
    scope: '范围',
    scopeAll: '我的全部公开图片',
    scopeAllPublic: '全部公开图片',
    scopeFolder: '指定文件夹',
    returnMethod: '返回方式',
    times: '次',
  },
}
