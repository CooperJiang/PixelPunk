/**
 * Activity Log Text - English Cyber Theme
 */
export const activity = {
  common: {
    rootFolder: 'Root node',
  },
  user: {
    login: 'System access initiated',
    loginWithIp: 'Neural link established from {ip}',
    register: 'New node registered',
  },
  file: {
    batchUpload: 'Injected {count} data units',
    batchUploadWithSize: 'Injected {count} data units ({size})',
    delete: 'Purged data unit「{fileName}」',
    rename: 'Renamed data unit「{oldName}」to「{newName}」',
    move: 'Relocated data unit「{fileName}」from {oldFolder} to {newFolder}',
    accessLevelChange: 'Modified access protocol for data unit「{fileName}」to {level}',
    expired: 'Data unit「{fileName}」expired and was auto-purged by the system',
    expiredBatch: '{count} data units expired and were auto-purged by the system',
  },
  folder: {
    create: 'Created folder「{folderName}」',
    createInParent: 'Created folder「{folderName}」in「{parentName}」',
    rename: 'Renamed folder「{oldName}」to「{newName}」',
    delete: 'Purged folder「{folderName}」',
    deleteWithFiles: 'Purged folder「{folderName}」(containing {count} data units)',
    accessLevelChange: 'Modified folder「{folderName}」access protocol from {oldLevel} to {newLevel}',
  },
  share: {
    create: 'Established share link',
    delete: 'Terminated share link',
    milestone: 'Share link reached {count} access events',
  },
  apikey: {
    create: 'Generated API key「{keyName}」',
    delete: 'Destroyed API key「{keyName}」',
    toggleEnable: 'Activated API key「{keyName}」',
    toggleDisable: 'Deactivated API key「{keyName}」',
    regenerate: 'Regenerated API key「{keyName}」',
  },
  randomApi: {
    create: 'Created random data API「{apiName}」',
    delete: 'Deleted random data API「{apiName}」',
    toggleEnable: 'Activated random data API「{apiName}」',
    toggleEnableWithFolder: 'Activated random data API「{apiName}」: {folderName}',
    toggleDisable: 'Deactivated random data API「{apiName}」',
    toggleDisableWithFolder: 'Deactivated random data API「{apiName}」: {folderName}',
  },
  profile: {
    update: 'Updated node profile',
    avatarUpdate: 'Updated node avatar',
    usernameUpdate: 'Modified node identifier',
    bioUpdate: 'Updated node bio',
    websiteUpdate: 'Updated node link',
    passwordChange: 'Modified access key',
    emailChange: 'Contact node changed to {email}',
  },
  security: {
    hotlinkProtectionEnable: 'Enabled hotlink protection',
    hotlinkProtectionDisable: 'Disabled hotlink protection',
    hotlinkProtectionChange: 'Modified hotlink protection strategy',
  },
  system: {
    batchDelete: 'Batch purged {count} data units',
    batchDeleteWithFolder: 'Batch purged {count} data units from「{folderName}」',
    adminDelete: 'Data unit was purged by system administrator',
    cleanup: 'System cleaned up {count} {type} records',
    guestImageExpired: 'System purged {count} expired guest data units',
  },
  accessLevel: {
    private: 'Private',
    public: 'Public',
    protected: 'Protected',
    unlisted: 'Hidden',
  },
}
