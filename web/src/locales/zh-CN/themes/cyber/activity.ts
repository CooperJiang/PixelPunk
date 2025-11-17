/**
 * 活动日志多语言文本 - 中文赛博主题
 */
export const activity = {
  common: {
    rootFolder: '根节点',
  },
  user: {
    login: '接入了系统',
    loginWithIp: '从 {ip} 建立了神经连接',
    register: '注册了新节点',
  },
  file: {
    batchUpload: '注入了 {count} 个数据单元',
    batchUploadWithSize: '注入了 {count} 个数据单元（{size}）',
    delete: '清除了数据单元「{fileName}」',
    rename: '将数据单元「{oldName}」重命名为「{newName}」',
    move: '将数据单元「{fileName}」从 {oldFolder} 转移到 {newFolder}',
    accessLevelChange: '将数据单元「{fileName}」的访问协议修改为{level}',
    expired: '数据单元「{fileName}」已过期被系统自动清除',
    expiredBatch: '有 {count} 个数据单元已过期被系统自动清除',
  },
  folder: {
    create: '创建了文件夹「{folderName}」',
    createInParent: '在「{parentName}」中创建了文件夹「{folderName}」',
    rename: '将文件夹「{oldName}」重命名为「{newName}」',
    delete: '清除了文件夹「{folderName}」',
    deleteWithFiles: '清除了文件夹「{folderName}」（包含 {count} 个数据单元）',
    accessLevelChange: '将文件夹「{folderName}」的访问协议从{oldLevel}改为{newLevel}',
  },
  share: {
    create: '建立了共享链接',
    delete: '断开了共享链接',
    milestone: '共享链接达到了 {count} 次访问',
  },
  apikey: {
    create: '生成了 API 密钥「{keyName}」',
    delete: '销毁了 API 密钥「{keyName}」',
    toggleEnable: '激活了 API 密钥「{keyName}」',
    toggleDisable: '停用了 API 密钥「{keyName}」',
    regenerate: '重新生成了 API 密钥「{keyName}」',
  },
  randomApi: {
    create: '创建了随机数据 API「{apiName}」',
    delete: '删除了随机数据 API「{apiName}」',
    toggleEnable: '激活了随机数据 API「{apiName}」',
    toggleEnableWithFolder: '激活了随机数据 API「{apiName}」: {folderName}',
    toggleDisable: '停用了随机数据 API「{apiName}」',
    toggleDisableWithFolder: '停用了随机数据 API「{apiName}」: {folderName}',
  },
  profile: {
    update: '更新了节点资料',
    avatarUpdate: '更新了节点头像',
    usernameUpdate: '修改了节点标识',
    bioUpdate: '更新了节点简介',
    websiteUpdate: '更新了节点链接',
    passwordChange: '修改了访问密钥',
    emailChange: '联系节点已更改为 {email}',
  },
  security: {
    hotlinkProtectionEnable: '启用了防盗链防护',
    hotlinkProtectionDisable: '关闭了防盗链防护',
    hotlinkProtectionChange: '修改了防盗链策略',
  },
  system: {
    batchDelete: '批量清除了 {count} 个数据单元',
    batchDeleteWithFolder: '从「{folderName}」批量清除了 {count} 个数据单元',
    adminDelete: '数据单元被系统管理员清除',
    cleanup: '系统清理了 {count} 条{type}数据',
    guestImageExpired: '系统清理了 {count} 个过期的访客数据单元',
  },
  accessLevel: {
    private: '私有',
    public: '公开',
    protected: '受保护',
    unlisted: '隐藏',
  },
}
