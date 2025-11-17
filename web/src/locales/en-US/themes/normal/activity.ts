/**
 * Activity Log Text - English Normal Theme
 */
export const activity = {
  common: {
    rootFolder: 'Root folder',
  },
  user: {
    login: 'Logged into the system',
    loginWithIp: 'Logged in from {ip}',
    register: 'Registered an account',
  },
  file: {
    batchUpload: 'Uploaded {count} files',
    batchUploadWithSize: 'Uploaded {count} files ({size})',
    delete: 'Deleted file「{fileName}」',
    rename: 'Renamed file「{oldName}」to「{newName}」',
    move: 'Moved file「{fileName}」from {oldFolder} to {newFolder}',
    accessLevelChange: 'Changed access level of file「{fileName}」to {level}',
    expired: 'File「{fileName}」expired and was automatically deleted by the system',
    expiredBatch: '{count} files expired and were automatically deleted by the system',
  },
  folder: {
    create: 'Created folder「{folderName}」',
    createInParent: 'Created folder「{folderName}」in「{parentName}」',
    rename: 'Renamed folder「{oldName}」to「{newName}」',
    delete: 'Deleted folder「{folderName}」',
    deleteWithFiles: 'Deleted folder「{folderName}」(containing {count} files)',
    accessLevelChange: 'Changed folder「{folderName}」access level from {oldLevel} to {newLevel}',
  },
  share: {
    create: 'Created a share link',
    delete: 'Deleted a share link',
    milestone: 'Share link reached {count} visits',
  },
  apikey: {
    create: 'Created API key「{keyName}」',
    delete: 'Deleted API key「{keyName}」',
    toggleEnable: 'Enabled API key「{keyName}」',
    toggleDisable: 'Disabled API key「{keyName}」',
    regenerate: 'Regenerated API key「{keyName}」',
  },
  randomApi: {
    create: 'Created random image API「{apiName}」',
    delete: 'Deleted random image API「{apiName}」',
    toggleEnable: 'Enabled random image API「{apiName}」',
    toggleEnableWithFolder: 'Enabled random image API「{apiName}」: {folderName}',
    toggleDisable: 'Disabled random image API「{apiName}」',
    toggleDisableWithFolder: 'Disabled random image API「{apiName}」: {folderName}',
  },
  profile: {
    update: 'Updated profile',
    avatarUpdate: 'Updated avatar',
    usernameUpdate: 'Changed username',
    bioUpdate: 'Updated bio',
    websiteUpdate: 'Updated website',
    passwordChange: 'Changed password',
    emailChange: 'Email changed to {email}',
  },
  security: {
    hotlinkProtectionEnable: 'Enabled hotlink protection',
    hotlinkProtectionDisable: 'Disabled hotlink protection',
    hotlinkProtectionChange: 'Modified hotlink protection settings',
  },
  system: {
    batchDelete: 'Batch deleted {count} files',
    batchDeleteWithFolder: 'Batch deleted {count} files from「{folderName}」',
    adminDelete: 'File was deleted by administrator',
    cleanup: 'System cleaned up {count} {type} records',
    guestImageExpired: 'System cleaned up {count} expired guest uploads',
  },
  accessLevel: {
    private: 'Private',
    public: 'Public',
    protected: 'Protected',
    unlisted: 'Unlisted',
  },
}
