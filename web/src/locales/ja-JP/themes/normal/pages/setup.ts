/**
 * System installation page texts
 */
export const setup = {
  loading: 'システム初期化中...',
  title: 'PIXELPUNKシステムインストール',
  subtitle: 'インテリジェント画像管理プラットフォームを設定する',
  success: {
    title: 'システムインストール成功',
    subtitle: 'PixelPunkインテリジェントリソース管理システムへようこそ',
    entries: {
      home: {
        title: 'ホームページに移動',
        desc: '基本的なエントリーとデータ表示',
      },
      upload: {
        title: '画像をアップロード',
        desc: '最初の画像をアップロード開始',
      },
      dashboard: {
        title: 'コントロールパネル',
        desc: 'システム統計とデータを表示',
      },
      settings: {
        title: 'システム設定',
        desc: '画像ホスティング設定をカスタマイズ',
      },
    },
    tip: 'クイックエントリーを選択してシステムの使用を開始',
  },
  database: {
    title: 'データベース',
    mysql: 'MySQL',
    sqlite: 'SQLite',
    fields: {
      path: 'パス',
      host: 'ホスト',
      port: 'ポート',
      user: 'ユーザー',
      password: 'パスワード',
      name: 'データベース名',
    },
    placeholders: {
      path: './pixelpunk.db',
      host: '127.0.0.1',
      port: '3306',
      user: 'root',
      password: 'mysqlパスワード',
      name: 'pixelpunk',
    },
    validation: {
      mysqlRequired: 'MySQL接続情報を入力してください',
      sqliteRequired: 'SQLiteファイルパスを指定してください',
      connectionFailed: '接続失敗',
      unknownError: '不明なエラー',
    },
    test: {
      success: {
        sqlite: 'SQLite準備完了',
        mysql: 'MySQL接続成功',
      },
      testing: 'テスト中...',
      button: 'データベース設定をテスト',
    },
  },
  admin: {
    title: '管理者アカウント',
    fields: {
      username: 'ユーザー名',
      password: 'パスワード',
    },
    placeholders: {
      username: 'root',
      password: '6文字以上',
    },
    validation: {
      mysqlInfoRequired: 'MySQL情報を完全に入力してください',
      sqlitePathRequired: 'SQLiteパスを入力してください',
      usernameRequired: '管理者ユーザー名を入力してください（3文字以上）',
      passwordRequired: '管理者パスワードを入力してください',
      passwordTooShort: '管理者パスワードは6文字以上である必要があります',
    },
  },
  redis: {
    title: 'REDIS',
    optional: 'オプション',
    fields: {
      host: 'ホスト',
      port: 'ポート',
      password: 'パスワード',
      db: 'DB',
    },
    placeholders: {
      host: 'localhost',
      port: '6379',
      password: '(空)',
      db: '0',
    },
  },
  vector: {
    title: 'QDRANTベクトルDB',
    optional: 'オプション',
    hint: 'ベクトルデータベースは画像検索とセマンティック検索機能をサポートするために使用されます。設定しない場合でも他の機能には影響しません',
    mode: {
      builtin: {
        title: 'ビルトインを使用',
        desc: 'システムが自動的にQdrantを起動して管理します',
        badge: '推奨',
      },
      external: {
        title: '外部を使用',
        desc: '既存のQdrantサービスに接続',
      },
    },
    advanced: {
      title: '⚙️ 詳細設定',
      hint: 'Qdrantポートを設定（ポートが使用中の場合は変更）',
      httpPort: {
        label: 'HTTPポート',
        desc: 'API呼び出し用（主な使用）',
      },
      grpcPort: {
        label: 'gRPCポート',
        desc: 'Qdrantが内部的に使用（直接使用しませんが、Qdrantが動作するにはこのポートが必要）',
      },
    },
    fields: {
      url: 'QDRANT URL',
      timeout: 'タイムアウト（秒）',
    },
    placeholders: {
      url: 'http://your-server:6333',
      timeout: '30',
      httpPort: '6333',
      grpcPort: '6334',
    },
    validation: {
      urlRequired: 'Qdrant URLを入力してください',
      invalidUrl: '無効なURL形式',
      urlInvalid: 'URLが無効です',
      formatError: '形式エラー',
    },
    test: {
      success: {
        builtin: 'URL形式が正しいです。インストール中にビルトインQdrantが自動的に起動します',
        external: 'URL形式が正しいです（実際の接続はインストール中に確認されます）',
      },
      testing: 'テスト中...',
      button: {
        builtin: '設定を確認',
        external: '接続をテスト',
      },
    },
  },
  actions: {
    install: 'システムインストール',
    installing: 'インストール中...',
  },
}

