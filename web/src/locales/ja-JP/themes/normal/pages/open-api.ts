/**
 * Open API page texts
 */
export const openApi = {
  page: {
    title: 'オープンAPI管理',
    subtitle: 'ランダム画像APIエンドポイントを管理し、サードパーティアプリケーションに画像サービスを提供',
    create: 'APIを作成',
    search: 'API名を検索...',
  },
  loading: '読み込み中...',
  footer: {
    total: '合計 {count}件のAPI',
  },
  empty: {
    title: 'APIなし',
    description: 'まだAPIエンドポイントを作成していません。上記のボタンをクリックして開始してください',
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
    toggleSuccess: 'ステータス更新成功',
    toggleFailed: 'ステータス更新失敗',
  },
  dialog: {
    create: {
      title: 'ランダム画像APIを作成',
      cancel: 'キャンセル',
      submit: '作成',
      form: {
        name: {
          label: 'API名',
          hint: 'このAPIを識別するために使用（管理しやすくするため）',
        },
        scope: {
          label: '画像スコープ',
        },
        folder: {
          label: 'フォルダを選択',
          hint: 'フォルダを選択すると、APIはそのフォルダ内の公開画像からランダムに返します',
        },
        returnType: {
          label: '返却方法',
          redirect: {
            label: '302リダイレクト',
            desc: '実際の画像URLを返す、パフォーマンスが良く、CDN加速をサポート',
          },
          direct: {
            label: '画像を直接返す',
            desc: '各更新で新しい画像、URLは固定されたまま、ウェブサイトの背景/壁紙ローテーションに適しています',
          },
        },
      },
      hint: '作成後、一意のAPIリンクが生成され、公開アクセスのある画像のみを返します',
    },
    edit: {
      title: 'API設定を編集',
      cancel: 'キャンセル',
      submit: '保存',
      form: {
        nameLabel: 'API名:',
        scope: {
          label: '画像スコープ',
        },
        folder: {
          label: 'フォルダを選択',
          hint: 'フォルダを選択すると、APIはそのフォルダ内の公開画像からランダムに返します',
        },
        returnType: {
          label: '返却方法',
          redirect: {
            label: '302リダイレクト',
            desc: '実際の画像URLを返す、パフォーマンスが良く、CDN加速をサポート',
          },
          direct: {
            label: '画像を直接返す',
            desc: '各更新で新しい画像、URLは固定されたまま、ウェブサイトの背景/壁紙ローテーションに適しています',
          },
        },
      },
    },
    delete: {
      title: 'APIを削除',
      message: 'API「{name}」を削除してもよろしいですか？',
      warning: '削除後は復元できません。慎重に操作してください',
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
      placeholder: 'カテゴリを選択（複数可）',
      all: 'すべてのカテゴリ',
    },
    tags: {
      label: 'タグフィルター',
      placeholder: 'タグを選択（複数可）',
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
      images: '画像',
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

