/**
 * 活动日志多语言文本 - 中文普通主题
 */
export const activity = {
  common: {
    rootFolder: '根目录',
  },
  user: {
    login: '登录了系统',
    loginWithIp: '从 {ip} 登录了系统',
    register: '注册了账号',
  },
  file: {
    batchUpload: '上传了 {count} 个文件',
    batchUploadWithSize: '上传了 {count} 个文件（{size}）',
    delete: '删除了文件「{fileName}」',
    rename: '将文件「{oldName}」重命名为「{newName}」',
    move: '将文件「{fileName}」从 {oldFolder} 移动到 {newFolder}',
    accessLevelChange: '将文件「{fileName}」的访问权限修改为{level}',
    expired: '文件「{fileName}」已过期被系统自动删除',
    expiredBatch: '有 {count} 个文件已过期被系统自动删除',
  },
  folder: {
    create: '创建了文件夹「{folderName}」',
    createInParent: '在「{parentName}」中创建了文件夹「{folderName}」',
    rename: '将文件夹「{oldName}」重命名为「{newName}」',
    delete: '删除了文件夹「{folderName}」',
    deleteWithFiles: '删除了文件夹「{folderName}」（包含 {count} 个文件）',
    accessLevelChange: '将文件夹「{folderName}」的权限从{oldLevel}改为{newLevel}',
  },
  share: {
    create: '创建了分享链接',
    delete: '删除了分享链接',
    milestone: '分享链接达到了 {count} 次访问',
  },
  apikey: {
    create: '创建了 API 密钥「{keyName}」',
    delete: '删除了 API 密钥「{keyName}」',
    toggleEnable: '启用了 API 密钥「{keyName}」',
    toggleDisable: '停用了 API 密钥「{keyName}」',
    regenerate: '重新生成了 API 密钥「{keyName}」',
  },
  randomApi: {
    create: '创建了随机图片 API「{apiName}」',
    delete: '删除了随机图片 API「{apiName}」',
    toggleEnable: '启用了随机图片 API「{apiName}」',
    toggleEnableWithFolder: '启用了随机图片 API「{apiName}」: {folderName}',
    toggleDisable: '停用了随机图片 API「{apiName}」',
    toggleDisableWithFolder: '停用了随机图片 API「{apiName}」: {folderName}',
  },
  profile: {
    update: '更新了个人资料',
    avatarUpdate: '更新了头像',
    usernameUpdate: '修改了用户名',
    bioUpdate: '更新了个人简介',
    websiteUpdate: '更新了个人网站',
    passwordChange: '修改了登录密码',
    emailChange: '邮箱已更改为 {email}',
  },
  security: {
    hotlinkProtectionEnable: '启用了防盗链保护',
    hotlinkProtectionDisable: '关闭了防盗链保护',
    hotlinkProtectionChange: '修改了防盗链设置',
  },
  system: {
    batchDelete: '批量删除了 {count} 个文件',
    batchDeleteWithFolder: '从「{folderName}」批量删除了 {count} 个文件',
    adminDelete: '文件被管理员删除',
    cleanup: '系统清理了 {count} 条{type}数据',
    guestImageExpired: '系统清理了 {count} 个过期的访客上传文件',
  },
  accessLevel: {
    private: '私有',
    public: '公开',
    protected: '受保护',
    unlisted: '不公开列表',
  },
}
