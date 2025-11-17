export interface ParticleBackgroundProps {
  maxParticles?: number
  colors?: string[]
  theme?: 'default' | 'docs'
  connectionDistance?: number // 粒子间连接的最大距离
  connectionOpacity?: number // 连接线的透明度
  showConnections?: boolean // 是否显示连接线
  particleSpeed?: number // 粒子移动速度倍数
  showGrid?: boolean // 是否显示网格
}
