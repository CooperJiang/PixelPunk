/**
 * Folder management page texts
 */
export const folders = {
  rootFolder: 'ルート',
  emptyState: {
    title: '現在のフォルダは空です',
    description: 'リソースをアップロードするか、新しいフォルダを作成して開始できます',
    createButton: 'フォルダを作成',
  },
  breadcrumb: {
    label: 'フォルダナビゲーション',
  },
  share: {
    createSuccess: 'シェアが正常に作成されました',
    createSuccessWithCopy: 'シェアが正常に作成されました。リンクをクリップボードにコピーしました',
    manualCopy: 'シェアリンクを手動でコピーしてください',
    description: 'シェア説明',
  },
  shareSelectTip: {
    message: '共有するフォルダまたは画像を選択してください',
  },
  folderGrid: {
    sectionTitle: 'フォルダ',
    dragTip: 'ドラッグしてソート',
    fileCount: '{count}ファイル',
    toast: {
      moveFailed: 'フォルダの移動に失敗しました',
      moveToRootSuccess: 'フォルダをルートに移動しました',
      moveSuccess: 'フォルダを「{name}」に移動しました',
    },
    contextMenu: {
      root: 'ルート',
      moveTo: '移動先...',
      rename: '名前変更/編集',
      delete: '削除',
    },
    actions: {
      selectAll: 'すべて選択',
      deselectAll: 'すべて選択解除',
      selectAllTitle: 'すべてのフォルダを選択',
      deselectAllTitle: 'すべてのフォルダの選択を解除',
      invert: '反転',
      invertTitle: 'フォルダを反転',
      deleteTitle: 'フォルダを削除',
    },
    visibility: {
      toggleTitle: 'クリックして可視性を切り替え',
      setPublic: '公開に設定',
      setPrivate: 'プライベートに設定',
    },
  },
  createFolderDialog: {
    title: {
      create: '新しいフォルダ',
      edit: 'フォルダを編集',
    },
    form: {
      name: {
        label: 'フォルダ名',
        placeholder: 'フォルダ名を入力してください',
        required: 'フォルダ名を入力してください',
        maxLength: 'フォルダ名は30文字を超えることはできません',
      },
      description: {
        label: '説明（オプション）',
        placeholder: 'フォルダの説明を追加',
      },
      permission: {
        label: '権限設定',
        public: {
          name: '公開',
          desc: '他の人が私のホームページを訪問したときにコンテンツを表示できます',
        },
        private: {
          name: 'プライベート',
          desc: 'このフォルダを他の人に表示しない',
        },
      },
    },
    actions: {
      create: '作成',
      save: '保存',
    },
    toast: {
      createSuccess: 'フォルダが正常に作成されました',
      updateSuccess: 'フォルダが正常に更新されました',
      operationFailed: 'フォルダ操作に失敗しました',
    },
  },
  fileGrid: {
    sectionTitle: 'ファイル',
    dragTip: 'ドラッグしてソート',
    selectedCount: '{count}ファイルを選択',
    badge: {
      duplicate: '重複',
    },
    actions: {
      aiInfo: 'AI情報',
      preview: 'プレビュー',
      copyLink: 'リンクをコピー',
      download: 'ダウンロード',
      delete: '削除',
      batchMode: '一括操作',
      batchDelete: '一括削除',
      cancelBatch: '一括操作をキャンセル',
      selectAll: 'すべて選択',
      deselectAll: 'すべて選択解除',
      selectAllTitle: 'すべてのファイルを選択',
      deselectAllTitle: 'すべてのファイルの選択を解除',
      invert: '反転',
      invertTitle: 'ファイルを反転',
      gridView: 'グリッド表示',
      listView: 'リスト表示',
    },
    meta: {
      size: 'ファイルサイズ',
      dimensions: 'ファイルサイズ',
      duplicateFile: '重複ファイル',
    },
    listView: {
      columns: {
        preview: 'プレビュー',
        filename: 'ファイル名',
        size: 'サイズ',
        dimensions: 'サイズ',
        date: 'アップロード日時',
        actions: 'アクション',
      },
      empty: 'ファイルなし',
    },
  },
  pageHeader: {
    title: 'マイフォルダ',
    breadcrumb: {
      root: 'ルート',
      navigation: 'フォルダナビゲーション',
    },
    upload: {
      default: 'ファイルをアップロード',
      toFolder: '「{name}」にアップロード',
    },
    actions: {
      newFolder: '新しいフォルダ',
    },
  },
  createShareDialog: {
    title: 'シェアを作成',
    sections: {
      basic: '基本情報',
      access: 'アクセス制御',
      advanced: '詳細設定',
      content: 'シェアコンテンツ',
    },
    form: {
      name: {
        label: 'シェア名',
        placeholder: 'シェア名を入力してください（オプション）',
      },
      description: {
        placeholder: 'シェア説明を入力してください（オプション）',
      },
      password: {
        label: 'アクセスパスワード',
        placeholder: 'アクセスパスワードを入力してください（オプション）',
      },
      expiredDays: {
        label: '有効期限（日）',
        placeholder: '0は期限なしを意味します',
      },
      maxViews: {
        label: '最大閲覧数',
        placeholder: '0は無制限を意味します',
      },
      collectVisitor: {
        label: '訪問者情報を収集',
        tip: '訪問者は閲覧時に情報フォームのポップアップが表示されます（必須ではありません）',
      },
      notification: {
        label: 'アクセス通知',
        tip: '有効にすると、閲覧数が指定数に達したときにメール通知を受信します',
      },
      notificationThreshold: {
        label: '通知しきい値',
        placeholder: '閲覧数に基づいて通知タイミングを設定',
        tip: '閲覧数がこの値に達すると、システムがメール通知を送信します',
      },
    },
    emptyItems: 'シェアするフォルダまたはファイルを選択してください',
    cancel: 'キャンセル',
    submit: 'シェアを作成',
    toast: {
      createSuccess: 'シェアが正常に作成されました',
      linkCopied: 'シェアが正常に作成されました。リンクをクリップボードにコピーしました',
      copyFailed: 'クリップボードへのコピーに失敗しました',
      copyManually: 'シェアリンクを手動でコピーしてください',
    },
  },
  deleteConfirmDialogs: {
    title: {
      folder: 'フォルダを削除',
      file: 'ファイルを削除',
      batchFile: 'ファイルを一括削除',
    },
    message: {
      folder: 'フォルダ「{name}」を削除してもよろしいですか？',
      file: 'ファイル「{name}」を削除してもよろしいですか？',
      batchFile: '{count}件のリソースを一括削除してもよろしいですか？',
    },
    warning: 'この操作は元に戻せません。慎重に操作してください。',
    confirmText: '削除を確認',
  },
  fileContextMenu: {
    toast: {
      noFileSelected: 'ファイルが選択されていません',
      moveFailed: 'ファイルの移動に失敗しました',
      moveToRootSuccess: 'ファイルをルートに移動しました',
      moveToFolderSuccess: 'ファイルを「{name}」に移動しました',
      loadFoldersFailed: 'フォルダリストの読み込みに失敗しました',
      loadSubfoldersFailed: 'サブフォルダの読み込みに失敗しました',
      noValidImageLink: '有効なオリジナル画像リンクがありません',
      noThumbLink: 'コピーするサムネイルリンクがありません',
      thumbLinkCopied: 'サムネイルリンクをコピーしました',
      copyFailed: 'コピーに失敗しました',
      originalLinkCopied: 'オリジナル画像リンクをコピーしました',
      noOriginalLink: 'コピーするオリジナル画像リンクがありません',
      shortLinkCopied: 'ショートリンクをコピーしました',
      noShortLink: 'コピーするショートリンクがありません',
      fileDeleted: 'ファイルを削除しました',
      deleteFailed: 'ファイルの削除に失敗しました',
    },
    menu: {
      root: 'ルート',
      preview: 'プレビュー',
      openNew: '新しいウィンドウで開く',
      copyLink: 'リンクをコピー',
      copyThumb: 'サムネイルをコピー',
      copyOriginal: 'オリジナルをコピー',
      copyShort: 'ショートリンクをコピー',
      download: 'ダウンロード',
      moveTo: '移動先...',
      delete: '削除',
    },
    confirm: {
      deleteFile: 'ファイル「{name}」を削除してもよろしいですか？この操作は元に戻せません。',
    },
  },
  folderManagement: {
    deleteSuccess: 'フォルダ「{name}」を削除しました',
    visibility: {
      public: '公開',
      private: 'プライベート',
      switched: 'フォルダ「{name}」を{level}に設定しました',
    },
  },
  dragSort: {
    folderSortUpdated: 'フォルダソートを更新しました',
    fileSortUpdated: 'ファイルソートを更新しました',
  },
  fileManagement: {
    copyLink: {
      protectedNotSupported: '保護されたファイルは外部シェアをサポートしていません。アプリケーション内で表示してください',
      publicLinkCopied: '公開リンクをクリップボードにコピーしました',
      privateLinkCopied: 'プライベートリンクをクリップボードにコピーしました',
    },
    download: {
      success: 'ファイル「{name}」のダウンロードが完了しました',
      cancelled: 'ダウンロードがキャンセルされたか失敗しました',
    },
    delete: {
      success: 'ファイル「{name}」を削除しました',
      selectFirst: 'まず削除するファイルを選択してください',
      batchSuccess: '{successCount}件のファイルを正常に削除しました',
      batchPartial: '{successCount}件のファイルを正常に削除しましたが、{failCount}件が失敗しました',
      batchFailed: '削除に失敗しました。{failCount}件のファイルの削除に失敗しました',
    },
    accessLevel: {
      public: '公開',
      private: 'プライベート',
      protected: '保護',
      switched: 'ファイル「{name}」を{level}に設定しました',
    },
  },
}
