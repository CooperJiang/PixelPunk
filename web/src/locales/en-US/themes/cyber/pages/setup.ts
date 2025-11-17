/**
 * System Installation Page Text - Cyber Style
 */
export const setup = {
  loading: 'SYSTEM INITIALIZING...',
  title: 'PIXELPUNK SYSTEM INSTALLATION',
  subtitle: 'Configure your intelligent image management platform',
  success: {
    title: 'System Deployment Successful',
    subtitle: 'Welcome to PixelPunk Intelligent Resource Management System',
    entries: {
      home: {
        title: 'Go to Homepage',
        desc: 'Basic entry and data display',
      },
      upload: {
        title: 'Data Injection',
        desc: 'Start injecting your first data unit',
      },
      dashboard: {
        title: 'Control Panel',
        desc: 'View system statistics and data',
      },
      settings: {
        title: 'System Config',
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
      mysqlRequired: 'Please complete MySQL connection info',
      sqliteRequired: 'Please specify SQLite file path',
      connectionFailed: 'Connection failed',
      unknownError: 'Unknown error',
    },
    test: {
      success: {
        sqlite: 'SQLite Ready',
        mysql: 'MySQL connection successful',
      },
      testing: 'Testing...',
      button: 'Database Config Test',
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
      password: 'at least 6 digits',
    },
    validation: {
      mysqlInfoRequired: 'Please fill in complete MySQL info',
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
    hint: 'Vector database supports image search and semantic search features, not configuring does not affect other features',
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
      title: '⚙️ Advanced Config',
      hint: 'Configure Qdrant ports (modify if ports are occupied)',
      httpPort: {
        label: 'HTTP Port',
        desc: 'Used for API calls (primary use)',
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
      urlInvalid: 'URL Invalid',
      formatError: 'Format Error',
    },
    test: {
      success: {
        builtin: 'URL format correct, built-in Qdrant will be automatically started during installation',
        external: 'URL format correct (actual connection needs to be verified during installation)',
      },
      testing: 'Testing...',
      button: {
        builtin: 'Verify Config',
        external: 'Test Connection',
      },
    },
  },
  actions: {
    install: 'System Deployment',
    installing: 'Deploying...',
  },
}
