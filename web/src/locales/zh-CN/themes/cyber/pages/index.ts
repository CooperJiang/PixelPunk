/**
 * 用户页面文案统一导出 - 赛博风格
 */
import { notFound } from './404'
import { auth } from './auth'
import { author } from './author'
import { automation } from './automation'
import { category } from './category'
import { dashboard } from './dashboard'
import { docs } from './docs'
import { explore } from './explore'
import { find } from './find'
import { folders } from './folders'
import { gallery } from './gallery'
import { hive } from './hive'
import { home } from './home'
import { openApi } from './open-api'
import { random } from './random'
import { refuse } from './refuse'
import { resetPassword } from './reset-password'
import { resource } from './resource'
import { settings } from './settings'
import { setup } from './setup'
import { share } from './share'
import { shares } from './shares'
import { tagManage } from './tag-manage'
import { upload } from './upload'

/**
 * 页面翻译统一收集（不展开）
 * 最终在主题 index.ts 中通过 ...pages 展开
 */
export const pages = {
  notFound,
  auth,
  author,
  automation,
  category,
  dashboard,
  docs,
  explore,
  find,
  folders,
  gallery,
  hive,
  home,
  openApi,
  random,
  refuse,
  resetPassword,
  resource,
  settings,
  setup,
  share,
  shares,
  tagManage,
  upload,
}
