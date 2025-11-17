import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  /* 主应用路由 - 使用标准布局 (已移除多模式功能) */
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        name: 'home',
        component: () => import('@/pages/home/index.vue'),
        meta: {
          title: 'routes.home',
        },
      },
      {
        path: 'explore',
        name: 'explore',
        component: () => import('@/pages/explore/index.vue'),
        meta: {
          title: 'routes.explore',
        },
      },
      {
        path: 'upload',
        name: 'upload',
        component: () => import('@/pages/upload/index.vue'),
        meta: {
          title: 'routes.upload',
        },
      },
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/pages/dashboard/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.dashboard',
        },
      },
      {
        path: 'messages',
        name: 'messages',
        component: () => import('@/pages/messages/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.messages',
        },
      },
      {
        path: 'resource',
        name: 'resource',
        component: () => import('@/pages/resource/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.resource',
        },
      },
      {
        path: 'folders',
        name: 'folders',
        component: () => import('@/pages/folders/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.folders',
        },
      },
      {
        path: 'folders/:folderPath*',
        name: 'foldersWithPath',
        component: () => import('@/pages/folders/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.folders',
        },
      },
      {
        path: 'shares',
        name: 'shares',
        component: () => import('@/pages/shares/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.shares',
        },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/pages/settings/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.settings',
        },
      },
      {
        path: 'tag-manage',
        name: 'tagManage',
        component: () => import('@/pages/TagManage/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.tagManage',
        },
      },
      {
        path: 'category-manage',
        name: 'categoryManage',
        component: () => import('@/pages/CategoryManage/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.categoryManage',
        },
      },
      {
        path: 'automation',
        name: 'automation',
        component: () => import('@/pages/automation/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.automation',
        },
      },
      {
        path: 'open-api',
        name: 'openApi',
        component: () => import('@/pages/open-api/index.vue'),
        meta: {
          requiresAuth: true,
          title: 'routes.openApi',
        },
      },
      {
        path: 'random',
        name: 'random',
        component: () => import('@/pages/random/index.vue'),
        meta: {
          title: 'routes.random',
        },
      },
    ],
  },

  {
    path: '/hive',
    name: 'hive',
    component: () => import('@/pages/hive/index.vue'),
    meta: {
      title: 'routes.hive',
      skipAuth: true,
    },
  },
  {
    path: '/docs',
    name: 'docs',
    component: () => import('@/pages/docs/index.vue'),
    meta: {
      title: 'routes.docs',
    },
  },
  {
    path: '/author/:id',
    name: 'author',
    component: () => import('@/pages/author/index.vue'),
    meta: {
      title: 'routes.author',
      skipAuth: true,
    },
  },
  {
    path: '/author/:authorId/folder/:folderId',
    name: 'authorFolder',
    component: () => import('@/pages/author/folder.vue'),
    meta: {
      title: 'routes.authorFolder',
      skipAuth: true,
    },
  },
  {
    path: '/share/:key',
    name: 'share',
    component: () => import('@/pages/share/index.vue'),
    meta: {
      title: 'routes.share',
      skipAuth: true,
    },
  },
  {
    path: '/auth',
    name: 'auth',
    component: () => import('@/pages/auth/index.vue'),
    meta: {
      guest: true,
      title: 'routes.auth',
    },
  },
  {
    path: '/reset-password',
    name: 'resetPassword',
    component: () => import('@/pages/reset-password/index.vue'),
    meta: {
      guest: true,
      title: 'routes.resetPassword',
      skipAuth: true,
    },
  },
  {
    path: '/setup',
    name: 'setup',
    component: () => import('@/pages/setup/index.vue'),
    meta: {
      title: 'routes.setup',
      skipAuth: true,
    },
  },

  /* 管理中心路由 */
  {
    path: '/admin',
    name: 'admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'routes.admin',
    },
    redirect: '/admin/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'adminDashboard',
        component: () => import('@/pages/admin/dashboard/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminDashboard',
        },
      },
      {
        path: 'users',
        name: 'adminUsers',
        component: () => import('@/pages/admin/users/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminUsers',
        },
      },
      {
        path: 'channels',
        name: 'adminChannels',
        component: () => import('@/pages/admin/channels/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminChannels',
        },
      },
      {
        path: 'settings',
        name: 'adminSettings',
        component: () => import('@/pages/admin/settings/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminSettings',
        },
      },
      {
        path: 'files',
        name: 'adminFiles',
        component: () => import('@/pages/admin/files/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminFiles',
        },
      },
      {
        path: 'tags',
        name: 'adminTags',
        component: () => import('@/pages/admin/tags/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminTags',
        },
      },
      {
        path: 'categories',
        name: 'adminCategories',
        component: () => import('@/pages/admin/categories/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminCategories',
        },
      },
      {
        path: 'tagging',
        name: 'adminTagging',
        component: () => import('@/pages/admin/tagging/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminTagging',
        },
      },
      {
        path: 'vectors',
        name: 'adminVectors',
        component: () => import('@/pages/admin/vectors/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminVectors',
        },
      },
      {
        path: 'ai',
        name: 'adminAI',
        component: () => import('@/pages/admin/ai/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminAI',
        },
      },
      {
        path: 'shares',
        name: 'adminShares',
        component: () => import('@/pages/admin/shares/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminShares',
        },
      },
      {
        path: 'construction',
        name: 'adminConstruction',
        component: () => import('@/pages/admin/construction/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminConstruction',
        },
      },
      {
        path: 'announcements',
        name: 'adminAnnouncements',
        component: () => import('@/pages/admin/announcements/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminAnnouncements',
        },
      },
      {
        path: 'content-review',
        name: 'adminContentReview',
        component: () => import('@/pages/admin/content-review/index.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: 'routes.adminContentReview',
        },
      },
    ],
  },

  /* 错误页面 */
  {
    path: '/404',
    name: 'notFound',
    component: () => import('@/pages/404/index.vue'),
    meta: {
      title: 'routes.notFound',
    },
  },
  {
    path: '/refuse',
    name: 'refuse',
    component: () => import('@/pages/refuse/index.vue'),
    meta: {
      title: 'routes.refuse',
      skipAuth: true,
    },
  },

  /* 通配符路由 */
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404',
  },
]

export default routes
