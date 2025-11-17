/**
 * SEO工具类 - 动态管理页面标题和meta标签
 */
export class SEOManager {
  static setTitle(title: string, siteName?: string) {
    let fullTitle = title

    if (siteName && siteName.trim()) {
      fullTitle = title ? `${title} - ${siteName}` : siteName
    }

    document.title = fullTitle
  }

  static setDescription(description: string) {
    if (!description || !description.trim()) return

    let metaDesc = document.querySelector('meta[name="description"]') as HTMLMetaElement
    if (!metaDesc) {
      metaDesc = document.createElement('meta')
      metaDesc.name = 'description'
      document.head.appendChild(metaDesc)
    }
    metaDesc.content = description
  }

  static setKeywords(keywords: string) {
    if (!keywords || !keywords.trim()) return

    let metaKeywords = document.querySelector('meta[name="keywords"]') as HTMLMetaElement
    if (!metaKeywords) {
      metaKeywords = document.createElement('meta')
      metaKeywords.name = 'keywords'
      document.head.appendChild(metaKeywords)
    }
    metaKeywords.content = keywords
  }

  static setSEO(options: { title?: string; siteName?: string; description?: string; keywords?: string }) {
    if (options.title !== undefined || options.siteName !== undefined) {
      this.setTitle(options.title || '', options.siteName)
    }

    if (options.description) {
      this.setDescription(options.description)
    }

    if (options.keywords) {
      this.setKeywords(options.keywords)
    }
  }

  static setDefaultSEO(settingsStore: any, pageTitle?: string) {
    this.setSEO({
      title: pageTitle,
      siteName: settingsStore.siteName,
      description: settingsStore.siteDescription,
      keywords: settingsStore.siteKeywords,
    })
  }
}
