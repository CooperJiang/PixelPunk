/**
 * Folder Management Page Text - Cyber Style
 */
export const folders = {
  rootFolder: 'ルートノード',
  emptyState: {
    title: '現在のデータディレクトリが空です',
    description: 'データリソースをアップロードするか、新しいデータディレクトリを作成して開始できます',
    createButton: 'データディレクトリを作成',
  },
  breadcrumb: {
    label: 'データディレクトリナビゲーション',
  },
  share: {
    createSuccess: 'データ共有を作成しました',
    createSuccessWithCopy: 'データ共有を作成しました。リンクをクリップボードにコピーしました',
    manualCopy: '共有リンクを手動でコピーしてください',
    description: '共有説明',
  },
  shareSelectTip: {
    message: '共有するデータディレクトリまたはデータユニットを選択してください',
  },
  folderGrid: {
    sectionTitle: 'データディレクトリ',
    dragTip: 'ドラッグして並べ替え',
    fileCount: '{count}件',
    toast: {
      moveFailed: 'データディレクトリの移動に失敗しました',
      moveToRootSuccess: 'データディレクトリをルートノードに移動しました',
      moveSuccess: 'データディレクトリを「{name}」に移動しました',
    },
    contextMenu: {
      root: 'ルートノード',
      moveTo: '移動先...',
      rename: '名前変更/編集',
      delete: '削除',
    },
    actions: {
      selectAll: 'すべて選択',
      deselectAll: 'すべての選択を解除',
      selectAllTitle: 'すべてのデータディレクトリを選択',
      deselectAllTitle: 'すべてのデータディレクトリの選択を解除',
      invert: '選択を反転',
      invertTitle: 'データディレクトリの選択を反転',
      deleteTitle: 'データディレクトリを削除',
    },
    visibility: {
      toggleTitle: 'クリックして可視性を切り替え',
      setPublic: '公開に設定',
      setPrivate: 'プライベートに設定',
    },
  },
  createFolderDialog: {
    title: {
      create: '新しいデータディレクトリ',
      edit: 'データディレクトリを編集',
    },
    form: {
      name: {
        label: 'ディレクトリ名',
        placeholder: 'データディレクトリ名を入力してください',
        required: 'データディレクトリ名を入力してください',
        maxLength: 'データディレクトリ名は30文字を超えることはできません',
      },
      description: {
        label: '説明（オプション）',
        placeholder: 'データディレクトリの説明を追加',
      },
      permission: {
        label: '権限設定',
        public: {
          name: '公開',
          desc: '他の人が私のホームページを訪問する際にコンテンツを閲覧できます',
        },
        private: {
          name: 'プライベート',
          desc: 'このデータディレクトリは他の人に表示されません',
        },
      },
    },
    actions: {
      create: '作成',
      save: '保存',
    },
    toast: {
      createSuccess: 'データディレクトリを作成しました',
      updateSuccess: 'データディレクトリを更新しました',
      operationFailed: 'データディレクトリ操作に失敗しました',
    },
  },
  fileGrid: {
    sectionTitle: 'データファイル',
    dragTip: 'ドラッグして並べ替え',
    selectedCount: '{count}件を選択',
    badge: {
      duplicate: '重複',
    },
    actions: {
      aiInfo: 'AI分析',
      preview: 'プレビュー',
      copyLink: 'リンクをコピー',
      download: 'ダウンロード',
      delete: '削除',
      batchMode: '一括プロトコル操作',
      batchDelete: '一括削除',
      cancelBatch: '一括プロトコル操作をキャンセル',
      selectAll: 'すべて選択',
      deselectAll: 'すべての選択を解除',
      selectAllTitle: 'すべてのデータユニットを選択',
      deselectAllTitle: 'すべてのデータユニットの選択を解除',
      invert: '選択を反転',
      invertTitle: 'データユニットの選択を反転',
      gridView: 'グリッドビュー',
      listView: 'リストビュー',
    },
    meta: {
      size: 'データユニットサイズ',
      dimensions: 'データユニットサイズ',
      duplicateFile: '重複データユニット',
    },
    listView: {
      columns: {
        preview: 'プレビュー',
        filename: 'データユニット名',
        size: 'サイズ',
        dimensions: 'サイズ',
        date: '注入時間',
        actions: 'アクション',
      },
      empty: 'データユニットがまだありません',
    },
  },
  pageHeader: {
    title: 'マイデータディレクトリ',
    breadcrumb: {
      root: 'ルートノード',
      navigation: 'データディレクトリナビゲーション',
    },
    upload: {
      default: 'データユニットを注入',
      toFolder: '「{name}」に注入',
    },
    actions: {
      newFolder: '新しいデータディレクトリ',
    },
  },
  createShareDialog: {
    title: '共有リンクを作成',
    sections: {
      basic: '基本情報ノード',
      access: 'アクセス制御ノード',
      advanced: '高度な設定ノード',
      content: '共有コンテンツノード',
    },
    form: {
      name: {
        label: '共有名',
        placeholder: '共有名を入力してください（オプション）',
      },
      description: {
        placeholder: '共有説明を入力してください（オプション）',
      },
      password: {
        label: 'アクセスパスワード',
        placeholder: 'アクセスパスワードを入力してください（オプション）',
      },
      expiredDays: {
        label: '有効期限日数ノード',
        placeholder: '0は期限なし',
      },
      maxViews: {
        label: '最大アクセス数ノード',
        placeholder: '0は無制限',
      },
      collectVisitor: {
        label: '訪問者情報収集ノード',
        tip: '訪問者が閲覧する際にポップアップウィンドウノードが表示されます（オプション）',
      },
      notification: {
        label: 'アクセス通知ノード',
        tip: '有効にすると、アクセス数が指定された数に達したときにメール通知ノードを受信します',
      },
      notificationThreshold: {
        label: '通知しきい値ノード',
        placeholder: '通知をトリガーするアクセス数を設定',
        tip: 'アクセス数がこの値に達すると、システムがメール通知ノードを送信します',
      },
    },
    emptyItems: '共有するデータディレクトリまたはデータユニットを選択してください',
    cancel: 'キャンセル',
    submit: '共有リンクを作成',
    toast: {
      createSuccess: '共有リンクを作成しました',
      linkCopied: '共有リンクを作成しました。リンクをクリップボードにコピーしました',
      copyFailed: 'クリップボードへのコピーに失敗しました',
      copyManually: '共有リンクを手動でコピーしてください',
    },
  },
  deleteConfirmDialogs: {
    title: {
      folder: 'データディレクトリを削除',
      file: 'データユニットを削除',
      batchFile: 'データユニットを一括削除',
    },
    message: {
      folder: 'データディレクトリ「{name}」を削除してもよろしいですか？',
      file: 'データユニット「{name}」を削除してもよろしいですか？',
      batchFile: '{count}個のリソースノードを一括削除してもよろしいですか？',
    },
    warning: 'この操作は元に戻せないノードです。慎重に進めてください。',
    confirmText: '削除を確認',
  },
  fileContextMenu: {
    toast: {
      noFileSelected: 'データユニットが選択されていません',
      moveFailed: 'データユニットの移動に失敗しました',
      moveToRootSuccess: 'データユニットをルートノードに移動しました',
      moveToFolderSuccess: 'データユニットを「{name}」に移動しました',
      loadFoldersFailed: 'データディレクトリリストの読み込みに失敗しました',
      loadSubfoldersFailed: 'サブデータディレクトリの読み込みに失敗しました',
      noValidImageLink: '有効な元のデータリンクがありません',
      noThumbLink: 'コピーするサムネイルリンクがありません',
      thumbLinkCopied: 'サムネイルリンクをコピーしました',
      copyFailed: 'コピー失敗',
      originalLinkCopied: '元のデータリンクをコピーしました',
      noOriginalLink: 'コピーする元のデータリンクがありません',
      shortLinkCopied: 'ショートリンクをコピーしました',
      noShortLink: 'コピーするショートリンクがありません',
      markdownLinkCopied: 'Markdown リンクをコピーしました',
      htmlLinkCopied: 'HTML タグをコピーしました',
      fileDeleted: 'データユニットを削除しました',
      deleteFailed: 'データユニットの削除に失敗しました',
    },
    menu: {
      root: 'ルートノード',
      preview: 'プレビュー',
      openNew: '新しいウィンドウで開く',
      copyLink: 'リンクをコピー',
      copyThumb: 'サムネイルをコピー',
      copyOriginal: '元のデータをコピー',
      copyShort: 'ショートリンクをコピー',
      copyMarkdown: 'Markdown をコピー',
      copyHtml: 'HTML をコピー',
      download: 'ダウンロード',
      moveTo: '移動先...',
      delete: '削除',
    },
    confirm: {
      deleteFile: 'データユニット「{name}」を削除してもよろしいですか？この操作は元に戻せません。',
    },
  },
  folderManagement: {
    deleteSuccess: 'データディレクトリ「{name}」を削除しました',
    visibility: {
      public: '公開プロトコル',
      private: 'プライベートアーカイブ',
      switched: 'データディレクトリ「{name}」を{level}に設定しました',
    },
  },
  dragSort: {
    folderSortUpdated: 'データディレクトリの並び順を更新しました',
    fileSortUpdated: 'データユニットの並び順を更新しました',
  },
  fileManagement: {
    copyLink: {
      protectedNotSupported: '保護されたデータは外部転送をサポートしていません。ターミナル内で閲覧してください',
      publicLinkCopied: '公開リンクをクリップボードにコピーしました',
      privateLinkCopied: 'プライベートリンクをクリップボードにコピーしました',
    },
    download: {
      success: 'データユニット「{name}」のダウンロードが完了しました',
      cancelled: 'ダウンロードがキャンセルされたか失敗しました',
    },
    delete: {
      success: 'データユニット「{name}」を削除しました',
      selectFirst: 'まず削除するデータユニットを選択してください',
      batchSuccess: '{successCount}個のデータユニットを削除しました',
      batchPartial: '{successCount}個のデータユニットを削除しました。{failCount}個が失敗しました',
      batchFailed: '削除に失敗しました。{failCount}個のデータユニットの削除に失敗しました',
    },
    accessLevel: {
      public: '公開プロトコル',
      private: 'プライベートアーカイブ',
      protected: 'セキュリティシールド',
      switched: 'データユニット「{name}」を{level}に設定しました',
    },
  },
}
