/**
 * 文件夹读取工具
 * 支持递归读取文件夹中的所有文件
 */

export interface FolderStats {
  totalFiles: number
  validFiles: number
  invalidFiles: number
}

/**
 * 递归读取文件夹中的所有文件
 * @param entry FileSystemEntry 对象
 * @returns Promise<File[]> 文件数组
 */
export async function readFolderRecursively(entry: FileSystemEntry): Promise<File[]> {
  const files: File[] = []

  if (entry.isFile) {
    const fileEntry = entry as FileSystemFileEntry
    try {
      const file = await new Promise<File>((resolve, reject) => {
        fileEntry.file(resolve, reject)
      })
      files.push(file)
    } catch (error) {
      console.warn('[FolderReader] Failed to read file:', entry.name, error)
    }
  } else if (entry.isDirectory) {
    const dirEntry = entry as FileSystemDirectoryEntry
    const reader = dirEntry.createReader()

    try {
      const entries = await readAllEntries(reader)

      for (const subEntry of entries) {
        const subFiles = await readFolderRecursively(subEntry)
        files.push(...subFiles)
      }
    } catch (error) {
      console.warn('[FolderReader] Failed to read directory:', entry.name, error)
    }
  }

  return files
}

/**
 * 读取目录的所有条目（处理分批返回的情况）
 * @param reader FileSystemDirectoryReader 对象
 * @returns Promise<FileSystemEntry[]> 条目数组
 */
function readAllEntries(reader: FileSystemDirectoryReader): Promise<FileSystemEntry[]> {
  return new Promise((resolve, reject) => {
    const allEntries: FileSystemEntry[] = []

    function readBatch() {
      reader.readEntries(
        (entries) => {
          if (entries.length === 0) {
            resolve(allEntries)
          } else {
            allEntries.push(...entries)
            readBatch()
          }
        },
        (error) => reject(error)
      )
    }

    readBatch()
  })
}

/**
 * 从拖拽 Items 中提取所有文件（支持文件夹）
 * @param items DataTransferItemList
 * @returns Promise<File[]> 文件数组
 */
export async function extractFilesFromItems(items: DataTransferItemList): Promise<File[]> {
  const allFiles: File[] = []

  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    const entry = item.webkitGetAsEntry()

    if (entry) {
      const files = await readFolderRecursively(entry)
      allFiles.push(...files)
    }
  }

  return allFiles
}

/**
 * 从 FileList 中读取所有文件（支持 webkitdirectory）
 * @param fileList FileList 对象
 * @returns File[] 文件数组
 */
export function extractFilesFromFileList(fileList: FileList): File[] {
  return Array.from(fileList)
}

/**
 * 过滤有效的图片文件
 * @param files 文件数组
 * @param isAllowedType 类型验证函数
 * @param isAllowedSize 大小验证函数
 * @returns 有效文件数组
 */
export function filterValidFiles(
  files: File[],
  isAllowedType: (type: string) => boolean,
  isAllowedSize: (size: number) => boolean
): File[] {
  return files.filter((file) => isAllowedType(file.type) && isAllowedSize(file.size))
}

/**
 * 获取文件夹统计信息
 * @param allFiles 所有文件
 * @param validFiles 有效文件
 * @returns FolderStats 统计信息
 */
export function getFolderStats(allFiles: File[], validFiles: File[]): FolderStats {
  return {
    totalFiles: allFiles.length,
    validFiles: validFiles.length,
    invalidFiles: allFiles.length - validFiles.length,
  }
}

/**
 * 检查是否包含文件夹
 * @param items DataTransferItemList
 * @returns boolean
 */
export function hasFolder(items: DataTransferItemList): boolean {
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    const entry = item.webkitGetAsEntry()
    if (entry?.isDirectory) {
      return true
    }
  }
  return false
}

/**
 * 检查浏览器是否支持文件夹上传
 * @returns boolean
 */
export function supportsFolderUpload(): boolean {
  const input = document.createElement('input')
  return 'webkitdirectory' in input
}
