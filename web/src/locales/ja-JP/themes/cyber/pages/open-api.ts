/**
 * Open API Page Text - Cyber Style
 */
export const openApi = {
  page: {
    title: 'Open API管理',
    subtitle: 'ランダム画像APIインターフェースを管理し、サードパーティアプリケーションにデータサービスを提供',
    create: 'APIを作成',
    search: 'API名を検索...',
  },
  loading: '読み込み中...',
  footer: {
    total: '合計{count}個のAPI',
  },
  empty: {
    title: 'まだAPIがありません',
    description: 'まだAPIインターフェースが作成されていません。上のボタンをクリックして作成を開始',
  },
  toast: {
    loadFailed: '読み込み失敗',
    createSuccess: '作成成功',
    createFailed: '作成失敗',
    updateSuccess: '更新成功',
    updateFailed: '更新失敗',
    deleteSuccess: '削除成功',
    deleteFailed: '削除失敗',
    copySuccess: 'APIアドレスをコピーしました',
    toggleSuccess: 'ステータスを更新しました',
    toggleFailed: 'ステータス更新失敗',
  },
  dialog: {
    create: {
      title: 'ランダム画像APIノードを作成',
      cancel: 'キャンセル',
      submit: '作成',
      form: {
        name: {
          label: 'APIノード名',
          hint: 'このAPIノードを識別するために使用し、管理を容易にします',
        },
        scope: {
          label: '画像スコープノード',
        },
        folder: {
          label: 'フォルダノードを選択',
          hint: 'フォルダノードを選択すると、APIはこのフォルダから公開画像をランダムに返します',
        },
        returnType: {
          label: '返却方法ノード',
          redirect: {
            label: '302リダイレクトノード',
            desc: '実際の画像アドレスノードを返します。パフォーマンスが良く、CDN加速をサポートします',
          },
          direct: {
            label: '画像ノードを直接返却',
            desc: '各更新は新しい画像ノードで、URLは固定のままです。ウェブサイトの背景/壁紙カルーセルに適しています',
          },
        },
      },
      hint: '作成後、一意のAPIリンクノードが生成され、公開アクセス可能な画像のみを返します',
    },
    edit: {
      title: 'API設定ノードを編集',
      cancel: 'キャンセル',
      submit: '保存',
      form: {
        nameLabel: 'APIノード名:',
        scope: {
          label: '画像スコープノード',
        },
        folder: {
          label: 'フォルダノードを選択',
          hint: 'フォルダノードを選択すると、APIはこのフォルダから公開画像をランダムに返します',
        },
        returnType: {
          label: '返却方法ノード',
          redirect: {
            label: '302リダイレクトノード',
            desc: '実際の画像アドレスノードを返します。パフォーマンスが良く、CDN加速をサポートします',
          },
          direct: {
            label: '画像ノードを直接返却',
            desc: '各更新は新しい画像ノードで、URLは固定のままです。ウェブサイトの背景/壁紙カルーセルに適しています',
          },
        },
      },
    },
    delete: {
      title: 'APIノードを削除',
      message: 'APIノード「{name}」を削除してもよろしいですか？',
      warning: '削除後は回復できませんノード。慎重に進めてください',
      cancel: 'キャンセル',
      confirm: '削除を確認',
    },
  },
  form: {
    name: {
      label: 'API名',
      placeholder: 'API名を入力してください',
      required: 'API名を入力してください',
    },
    description: {
      label: 'API説明',
      placeholder: 'API説明を入力してください（オプション）',
    },
    categories: {
      label: 'カテゴリフィルター',
      placeholder: 'カテゴリを選択（複数選択可）',
      all: 'すべてのカテゴリ',
    },
    tags: {
      label: 'タグフィルター',
      placeholder: 'タグを選択（複数選択可）',
      all: 'すべてのタグ',
    },
    status: {
      label: 'ステータス',
      enabled: '有効',
      disabled: '無効',
    },
  },
  card: {
    status: {
      enabled: '有効',
      disabled: '無効',
    },
    actions: {
      open: 'APIを開く',
      copy: 'アドレスをコピー',
      edit: '編集',
      delete: '削除',
      toggle: 'ステータスを切り替え',
    },
    stats: {
      categories: 'カテゴリ',
      tags: 'タグ',
      images: 'データユニット',
      calls: '呼び出し',
    },
    filters: {
      allCategories: 'すべてのカテゴリ',
      allTags: 'すべてのタグ',
    },
    lastCall: {
      never: '呼び出しなし',
      justNow: 'たった今',
      minutesAgo: '{minutes}分前',
      hoursAgo: '{hours}時間前',
      daysAgo: '{days}日前',
    },
    returnType: {
      direct: '直接返却',
      redirect: '302リダイレクト',
    },
    scope: 'スコープ',
    scopeAll: 'すべてのマイ公開画像',
    scopeAllPublic: 'すべての公開画像',
    scopeFolder: '指定フォルダ',
    returnMethod: '返却方法',
    times: '回',
  },
}
