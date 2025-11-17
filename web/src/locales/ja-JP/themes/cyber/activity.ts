/**
 * Activity Log Text - Japanese Cyber Theme
 */
export const activity = {
  common: {
    rootFolder: 'ルートノード',
  },
  user: {
    login: 'システムアクセス開始',
    loginWithIp: '{ip} からニューラルリンク確立',
    register: '新規ノード登録完了',
  },
  file: {
    batchUpload: '{count} 個のデータユニットをインジェクト',
    batchUploadWithSize: '{count} 個のデータユニット（{size}）をインジェクト',
    delete: 'データユニット「{fileName}」をパージ',
    rename: 'データユニット「{oldName}」を「{newName}」にリネーム',
    move: 'データユニット「{fileName}」を {oldFolder} から {newFolder} に再配置',
    accessLevelChange: 'データユニット「{fileName}」のアクセスプロトコルを {level} に変更',
    expired: 'データユニット「{fileName}」が期限切れとなり、システムにより自動パージされました',
    expiredBatch: '{count} 個のデータユニットが期限切れとなり、システムにより自動パージされました',
  },
  folder: {
    create: 'フォルダ「{folderName}」を作成',
    createInParent: '「{parentName}」にフォルダ「{folderName}」を作成',
    rename: 'フォルダ「{oldName}」を「{newName}」にリネーム',
    delete: 'フォルダ「{folderName}」をパージ',
    deleteWithFiles: 'フォルダ「{folderName}」（{count} 個のデータユニットを含む）をパージ',
    accessLevelChange: 'フォルダ「{folderName}」のアクセスプロトコルを {oldLevel} から {newLevel} に変更',
  },
  share: {
    create: '共有リンクを確立',
    delete: '共有リンクを終了',
    milestone: '共有リンクが {count} 回のアクセスイベントに到達',
  },
  apikey: {
    create: 'API キー「{keyName}」を生成',
    delete: 'API キー「{keyName}」を破棄',
    toggleEnable: 'API キー「{keyName}」をアクティブ化',
    toggleDisable: 'API キー「{keyName}」を非アクティブ化',
    regenerate: 'API キー「{keyName}」を再生成',
  },
  randomApi: {
    create: 'ランダムデータ API「{apiName}」を作成',
    delete: 'ランダムデータ API「{apiName}」を削除',
    toggleEnable: 'ランダムデータ API「{apiName}」をアクティブ化',
    toggleEnableWithFolder: 'ランダムデータ API「{apiName}」をアクティブ化: {folderName}',
    toggleDisable: 'ランダムデータ API「{apiName}」を非アクティブ化',
    toggleDisableWithFolder: 'ランダムデータ API「{apiName}」を非アクティブ化: {folderName}',
  },
  profile: {
    update: 'ノードプロフィールを更新',
    avatarUpdate: 'ノードアバターを更新',
    usernameUpdate: 'ノード識別子を変更',
    bioUpdate: 'ノードバイオを更新',
    websiteUpdate: 'ノードリンクを更新',
    passwordChange: 'アクセスキーを変更',
    emailChange: 'コンタクトノードを {email} に変更',
  },
  security: {
    hotlinkProtectionEnable: 'ホットリンク保護を有効化',
    hotlinkProtectionDisable: 'ホットリンク保護を無効化',
    hotlinkProtectionChange: 'ホットリンク保護戦略を変更',
  },
  system: {
    batchDelete: '{count} 個のデータユニットを一括パージ',
    batchDeleteWithFolder: '「{folderName}」から {count} 個のデータユニットを一括パージ',
    adminDelete: 'データユニットがシステム管理者によりパージされました',
    cleanup: 'システムが {count} 件の {type} レコードをクリーンアップ',
    guestImageExpired: 'システムが {count} 個の期限切れゲストデータユニットをパージ',
  },
  accessLevel: {
    private: 'プライベート',
    public: 'パブリック',
    protected: 'プロテクト',
    unlisted: 'ヒドゥン',
  },
}
