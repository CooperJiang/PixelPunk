/**
 * Activity Log Text - Japanese Normal Theme
 */
export const activity = {
  common: {
    rootFolder: 'ルートフォルダ',
  },
  user: {
    login: 'システムにログインしました',
    loginWithIp: '{ip} からログインしました',
    register: 'アカウントを登録しました',
  },
  file: {
    batchUpload: '{count} 個のファイルをアップロードしました',
    batchUploadWithSize: '{count} 個のファイル（{size}）をアップロードしました',
    delete: 'ファイル「{fileName}」を削除しました',
    rename: 'ファイル「{oldName}」を「{newName}」に名前変更しました',
    move: 'ファイル「{fileName}」を {oldFolder} から {newFolder} に移動しました',
    accessLevelChange: 'ファイル「{fileName}」のアクセスレベルを {level} に変更しました',
    expired: 'ファイル「{fileName}」が期限切れとなり、システムにより自動削除されました',
    expiredBatch: '{count} 個のファイルが期限切れとなり、システムにより自動削除されました',
  },
  folder: {
    create: 'フォルダ「{folderName}」を作成しました',
    createInParent: '「{parentName}」にフォルダ「{folderName}」を作成しました',
    rename: 'フォルダ「{oldName}」を「{newName}」に名前変更しました',
    delete: 'フォルダ「{folderName}」を削除しました',
    deleteWithFiles: 'フォルダ「{folderName}」（{count} 個のファイルを含む）を削除しました',
    accessLevelChange: 'フォルダ「{folderName}」のアクセスレベルを {oldLevel} から {newLevel} に変更しました',
  },
  share: {
    create: '共有リンクを作成しました',
    delete: '共有リンクを削除しました',
    milestone: '共有リンクが {count} 回訪問されました',
  },
  apikey: {
    create: 'API キー「{keyName}」を作成しました',
    delete: 'API キー「{keyName}」を削除しました',
    toggleEnable: 'API キー「{keyName}」を有効にしました',
    toggleDisable: 'API キー「{keyName}」を無効にしました',
    regenerate: 'API キー「{keyName}」を再生成しました',
  },
  randomApi: {
    create: 'ランダム画像 API「{apiName}」を作成しました',
    delete: 'ランダム画像 API「{apiName}」を削除しました',
    toggleEnable: 'ランダム画像 API「{apiName}」を有効にしました',
    toggleEnableWithFolder: 'ランダム画像 API「{apiName}」を有効にしました: {folderName}',
    toggleDisable: 'ランダム画像 API「{apiName}」を無効にしました',
    toggleDisableWithFolder: 'ランダム画像 API「{apiName}」を無効にしました: {folderName}',
  },
  profile: {
    update: 'プロフィールを更新しました',
    avatarUpdate: 'アバターを更新しました',
    usernameUpdate: 'ユーザー名を変更しました',
    bioUpdate: '自己紹介を更新しました',
    websiteUpdate: 'ウェブサイトを更新しました',
    passwordChange: 'パスワードを変更しました',
    emailChange: 'メールアドレスを {email} に変更しました',
  },
  security: {
    hotlinkProtectionEnable: 'ホットリンク保護を有効にしました',
    hotlinkProtectionDisable: 'ホットリンク保護を無効にしました',
    hotlinkProtectionChange: 'ホットリンク保護設定を変更しました',
  },
  system: {
    batchDelete: '{count} 個のファイルを一括削除しました',
    batchDeleteWithFolder: '「{folderName}」から {count} 個のファイルを一括削除しました',
    adminDelete: 'ファイルが管理者により削除されました',
    cleanup: 'システムが {count} 件の {type} レコードをクリーンアップしました',
    guestImageExpired: 'システムが {count} 個の期限切れゲストアップロードをクリーンアップしました',
  },
  accessLevel: {
    private: '非公開',
    public: '公開',
    protected: '保護',
    unlisted: '非表示',
  },
}
