/**
 * 设置卡片组件 Props
 */
export interface SettingCardProps {
  /**
   * 卡片图标（FontAwesome类名）
   */
  icon: string

  /**
   * 卡片标题
   */
  title: string

  /**
   * 提示文本（可选）
   */
  hint?: string

  /**
   * 是否显示（可选，默认为true）
   */
  show?: boolean
}
