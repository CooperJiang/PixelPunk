/* 文档页面相关类型定义 */

/* 导航项类型 */
export interface NavItem {
  id: string
  title: string
  href: string
  icon?: string
}

/* 代码示例标签类型 */
export interface CodeTab {
  id: string
  name: string
  language: string
}

/* 代码示例类型 */
export interface CodeExample {
  id: string
  title: string
  description?: string
  tabs: {
    [key: string]: string
  }
}

/* 复制状态类型 */
export interface CopyStatus {
  [key: string]: boolean
}

/* 文档页面状态类型 */
export interface DocsPageState {
  activeSection: string
  activeCodeTab: string
  copyStatus: CopyStatus
  mobileMenuOpen: boolean
  showBackToTop: boolean
  readingProgress: number
}

/* JSON响应示例类型 */
export interface JsonResponseExample {
  success: string
  error: string
}

/* 功能特性类型 */
export interface FeatureItem {
  id: string
  icon: string
  title: string
  description: string
}
