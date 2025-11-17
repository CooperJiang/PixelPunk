/**
 * 基础类型定义
 * 包含项目中最基础的通用类型

/* ID类型 */
export type ID = string | number

/* 时间戳类型 */
export type Timestamp = number

/* 空值类型 */
export type Nullable<T> = T | null
export type Optional<T> = T | undefined
export type Maybe<T> = T | null | undefined

/* 函数类型 */
export type VoidFunction = () => void
export type AsyncVoidFunction = () => Promise<void>
export type EventHandler<T = Event> = (event: T) => void
export type AsyncEventHandler<T = Event> = (event: T) => Promise<void>
