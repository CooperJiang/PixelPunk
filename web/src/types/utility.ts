/**
 * 工具类型定义
 * 包含TypeScript工具类型和通用类型辅助

/* 深度只读 */
export type DeepReadonly<T> = {
  readonly [P in keyof T]: T[P] extends object ? DeepReadonly<T[P]> : T[P]
}

/* 深度部分 */
export type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P]
}

/* 排除字段 */
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

/* 选择字段 */
export type Pick<T, K extends keyof T> = {
  [P in K]: T[P]
}

/* 必需字段 */
export type Required<T> = {
  [P in keyof T]-?: T[P]
}

/* 键值对 */
export interface KeyValuePair<K = string, V = any> {
  key: K
  value: V
}

/* 字典类型 */
export type Dictionary<T = any> = Record<string, T>

/* Promise 解包 */
export type Awaited<T> = T extends Promise<infer U> ? U : T

/* 函数参数类型 */
export type Parameters<T extends (...args: any) => any> = T extends (...args: infer P) => any ? P : never

/* 函数返回类型 */
export type ReturnType<T extends (...args: any) => any> = T extends (...args: any) => infer R ? R : any

/* 构造函数类型 */
export type Constructor<T = any> = new (...args: any[]) => T

/* 类实例类型 */
export type Instance<T extends Constructor> = T extends Constructor<infer I> ? I : never
