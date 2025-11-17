/**
 * User Pages Text Unified Export
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
 * Page translations unified collection (not expanded)
 * Finally expanded via ...pages in theme index.ts
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
