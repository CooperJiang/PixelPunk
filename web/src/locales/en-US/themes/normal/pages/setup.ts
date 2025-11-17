/**
 * System installation page texts
 */
export const setup = {
  loading: 'SYSTEM INITIALIZING...',
  title: 'PIXELPUNK SYSTEM INSTALLATION',
  subtitle: 'Configure your image management platform',
  success: {
    title: 'System Installation Successful',
    subtitle: 'Welcome to PixelPunk Intelligent Resource Management System',
    entries: {
      home: {
        title: 'Go to Homepage',
        desc: 'Basic entry and data display',
      },
      upload: {
        title: 'Upload Images',
        desc: 'Start uploading your first image',
      },
      dashboard: {
        title: 'Control Panel',
        desc: 'View system statistics and data',
      },
      settings: {
        title: 'System Configuration',
        desc: 'Customize your image hosting settings',
      },
    },
    tip: 'Select a quick entry to start using the system',
  },
  database: {
    title: 'DATABASE',
    mysql: 'MySQL',
    sqlite: 'SQLite',
    fields: {
      path: 'PATH',
      host: 'HOST',
      port: 'PORT',
      user: 'USER',
      password: 'PASSWORD',
      name: 'DATABASE NAME',
    },
    placeholders: {
      path: './pixelpunk.db',
      host: '127.0.0.1',
      port: '3306',
      user: 'root',
      password: 'mysql password',
      name: 'pixelpunk',
    },
    validation: {
      mysqlRequired: 'Please complete MySQL connection information',
      sqliteRequired: 'Please specify SQLite file path',
      connectionFailed: 'Connection failed',
      unknownError: 'Unknown error',
    },
    test: {
      success: {
        sqlite: 'SQLite ready',
        mysql: 'MySQL connection successful',
      },
      testing: 'Testing...',
      button: 'Test Database Configuration',
    },
  },
  admin: {
    title: 'ADMIN ACCOUNT',
    fields: {
      username: 'USERNAME',
      password: 'PASSWORD',
    },
    placeholders: {
      username: 'root',
      password: 'At least 6 characters',
    },
    validation: {
      mysqlInfoRequired: 'Please fill in complete MySQL information',
      sqlitePathRequired: 'Please fill in SQLite path',
      usernameRequired: 'Please enter administrator username (at least 3 characters)',
      passwordRequired: 'Please enter administrator password',
      passwordTooShort: 'Administrator password must be at least 6 characters',
    },
  },
  redis: {
    title: 'REDIS',
    optional: 'OPTIONAL',
    fields: {
      host: 'HOST',
      port: 'PORT',
      password: 'PASSWORD',
      db: 'DB',
    },
    placeholders: {
      host: 'localhost',
      port: '6379',
      password: '(empty)',
      db: '0',
    },
  },
  vector: {
    title: 'QDRANT VECTOR DB',
    optional: 'OPTIONAL',
    hint: 'Vector database is used to support image search and semantic search features. Not configuring it will not affect other features',
    mode: {
      builtin: {
        title: 'Use Built-in',
        desc: 'System will automatically start and manage Qdrant',
        badge: 'Recommended',
      },
      external: {
        title: 'Use External',
        desc: 'Connect to existing Qdrant service',
      },
    },
    advanced: {
      title: '⚙️ Advanced Configuration',
      hint: 'Configure Qdrant ports (modify if ports are occupied)',
      httpPort: {
        label: 'HTTP Port',
        desc: 'For API calls (main usage)',
      },
      grpcPort: {
        label: 'gRPC Port',
        desc: "Used internally by Qdrant (we don't use it directly, but Qdrant needs this port to run)",
      },
    },
    fields: {
      url: 'QDRANT URL',
      timeout: 'TIMEOUT (seconds)',
    },
    placeholders: {
      url: 'http://your-server:6333',
      timeout: '30',
      httpPort: '6333',
      grpcPort: '6334',
    },
    validation: {
      urlRequired: 'Please fill in Qdrant URL',
      invalidUrl: 'Invalid URL format',
      urlInvalid: 'URL invalid',
      formatError: 'Format error',
    },
    test: {
      success: {
        builtin: 'URL format correct, built-in Qdrant will be automatically started during installation',
        external: 'URL format correct (actual connection will be verified during installation)',
      },
      testing: 'Testing...',
      button: {
        builtin: 'Verify Configuration',
        external: 'Test Connection',
      },
    },
  },
  actions: {
    install: 'System Installation',
    installing: 'Installing...',
  },
}
