import { computed } from 'vue'
import { TIMING } from '@/constants'

export function useCodeExamples(currentDomain: any) {
  const codeExamples = computed(() => ({
    curl: {
      single: `# 基础上传
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "file=@image.jpg"

# 带参数上传
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "file=@image.jpg" \\
  -F "access_level=public" \\
  -F "optimize=true" \\
  -F "filePath=projects/website"

# 指定文件夹ID上传
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "file=@image.jpg" \\
  -F "folderId=folder_12345" \\
  -F "optimize=1"`,

      batch: `# 批量上传多个文件
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "files[]=@image1.jpg" \\
  -F "files[]=@image2.png" \\
  -F "files[]=@image3.gif" \\
  -F "filePath=batch/$(date +%Y%m%d)"

# 支持的多种字段名格式
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "images[]=@photo1.jpg" \\
  -F "images[]=@photo2.jpg" \\
  -F "access_level=public"`,
    },

    javascript: {
      single: `// 基础上传函数
async function uploadFile(file, options = {}) {
  const formData = new FormData();
  formData.append('file', file);

  if (options.access_level) formData.append('access_level', options.access_level);
  if (options.optimize) formData.append('optimize', options.optimize ? 'true' : 'false');
  if (options.filePath) formData.append('filePath', options.filePath);
  if (options.folderId) formData.append('folderId', options.folderId);

  try {
    const response = await fetch('${currentDomain.value}/api/v1/external/upload', {
      method: 'POST',
      headers: {
        'x-pixelpunk-key': 'YOUR_API_KEY'
      },
      body: formData
    });

    const result = await response.json();

    if (result.code === 0) {
      if (result.data.oversized_files) {
      }
      if (result.data.upload_errors) {
      }

      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('Upload failed:', error);
    throw error;
  }
}

async function uploadMultipleImages(files, options = {}) {
  const formData = new FormData();

  const fieldName = options.fieldName || 'files[]';
  files.forEach(file => {
    formData.append(fieldName, file);
  });

  if (options.filePath) formData.append('filePath', options.filePath);
  if (options.access_level) formData.append('access_level', options.access_level);
  if (options.optimize) formData.append('optimize', options.optimize ? 'true' : 'false');

}`,
    },

    nodejs: {
      single: `// Node.js + Express + multer 示例
const express = require('express');
const multer = require('multer');
const FormData = require('form-data');
const fs = require('fs');
const axios = require('axios');

const app = express();

const upload = multer({ dest: 'uploads/' });

async function uploadToAPI(filePath, options = {}) {
  try {
    const formData = new FormData();
    formData.append('file', fs.createReadStream(filePath));

    if (options.access_level) formData.append('access_level', options.access_level);
    if (options.optimize) formData.append('optimize', options.optimize.toString());
    if (options.filePath) formData.append('filePath', options.filePath);
    if (options.folderId) formData.append('folderId', options.folderId);

    const response = await axios.post('${currentDomain.value}/api/v1/external/upload', formData, {
      headers: {
        'x-pixelpunk-key': 'YOUR_API_KEY',
        ...formData.getHeaders()
      },
      timeout: TIMING.REQUEST.LONG_TIMEOUT
    });

    return response.data;
  } catch (error) {
    console.error('Upload failed:', error.response?.data || error.message);
    throw error;
  }
}

async function uploadMultipleToAPI(filePaths, options = {}) {
  try {
    const formData = new FormData();

    filePaths.forEach(filePath => {
      formData.append('files[]', fs.createReadStream(filePath));
    });

    if (options.access_level) formData.append('access_level', options.access_level);
    if (options.optimize) formData.append('optimize', options.optimize.toString());
    if (options.filePath) formData.append('filePath', options.filePath);
    if (options.folderId) formData.append('folderId', options.folderId);

    const response = await axios.post('${currentDomain.value}/api/v1/external/upload', formData, {
      headers: {
        'x-pixelpunk-key': 'YOUR_API_KEY',
        ...formData.getHeaders()
      },
      timeout: 60000
    });

    return response.data;
  } catch (error) {
    console.error('Batch upload failed:', error.response?.data || error.message);
    throw error;
  }
}

app.post('/upload', upload.single('file'), async (req, res) => {
  try {
    const result = await uploadToAPI(req.file.path, {
      access_level: req.body.access_level || 'private',
      optimize: req.body.optimize === 'true',
      filePath: req.body.filePath,
      folderId: req.body.folderId
    });

    fs.unlinkSync(req.file.path);

    res.json(result);
  } catch (error) {
    res.status(500).json({ _error: error.message });
  }
});

app.post('/upload/batch', upload.array('files', 10), async (req, res) => {
  try {
    const filePaths = req.files.map(file => file.path);

    const result = await uploadMultipleToAPI(filePaths, {
      access_level: req.body.access_level || 'private',
      optimize: req.body.optimize === 'true',
      filePath: req.body.filePath,
      folderId: req.body.folderId
    });

    req.files.forEach(file => fs.unlinkSync(file.path));

    res.json(result);
  } catch (error) {
    res.status(500).json({ _error: error.message });
  }
});

app.listen(3000, () => {
});

`,
    },

    python: {
      basic: `import requests

def upload_image(file_path, api_key):
    url = "${currentDomain.value}/api/v1/external/upload"

    headers = {
        "x-pixelpunk-key": api_key
    }

    files = {
        "file": open(file_path, "rb")
    }

    data = {
        "access_level": "private"
    }

    try:
        response = requests.post(url, headers=headers, files=files, data=data)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"上传失败: {e}")
        return None
    finally:
        files["file"].close()`,
    },

    java: {
      basic: `import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

public class ImageUploader {
    private static final String API_KEY = "YOUR_API_KEY";
    private static final String BASE_URL = "${currentDomain.value}";
    private static final String UPLOAD_URL = BASE_URL + "/api/v1/external/upload";

    public static void main(String[] args) {
        try {
            List<File> images = new ArrayList<>();
            images.add(new File("image1.jpg"));
            images.add(new File("image2.png"));
            images.add(new File("image3.gif"));

            List<String> imagePaths = images.stream()
                .map(File::getAbsolutePath)
                .collect(Collectors.toList());

            String folderId = "folder_12345";
            String filePath = "projects/website";
            String accessLevel = "private";
            boolean optimize = true;

            uploadSingleFile(images.get(0), accessLevel, optimize, filePath, folderId);

            uploadMultipleFiles(imagePaths, accessLevel, optimize, filePath, folderId);
        } catch (IOException | InterruptedException e) {
            e.printStackTrace();
        }
    }

    private static void uploadSingleFile(File file, String accessLevel, boolean optimize, String filePath, String folderId) throws IOException, InterruptedException {
        String command = String.format("curl -X POST \"%s/api/v1/external/upload\" -H \"x-pixelpunk-key: %s\" -F \"file=@%s\" -F \"access_level=%s\" -F \"optimize=%s\" -F \"filePath=%s\" -F \"folderId=%s\"",
            BASE_URL, API_KEY, file.getAbsolutePath(), accessLevel, optimize, filePath, folderId);

        Process process = Runtime.getRuntime().exec(command);
        process.waitFor(30, TimeUnit.SECONDS);

    }

    private static void uploadMultipleFiles(List<String> imagePaths, String accessLevel, boolean optimize, String filePath, String folderId) throws IOException, InterruptedException {
        String command = String.format("curl -X POST \"%s/api/v1/external/upload\" -H \"x-pixelpunk-key: %s\" -F \"files[]=@%s\" -F \"access_level=%s\" -F \"optimize=%s\" -F \"filePath=%s\" -F \"folderId=%s\"",
            BASE_URL, API_KEY, String.join("\" -F \"files[]=@\" ", imagePaths), accessLevel, optimize, filePath, folderId);

        Process process = Runtime.getRuntime().exec(command);
        process.waitFor(30, TimeUnit.SECONDS);

    }
}`,
    },

    go: {
      single: `package main

import (
\t"bytes"
\t"encoding/json"
\t"fmt"
\t"io"
\t"mime/multipart"
\t"net/http"
\t"os"
\t"path/filepath"
\t"time"
)

type UploadResponse struct {
\tCode    int    \`json:"code"\`
\tMessage string \`json:"message"\`
\tData    struct {
\t\tUploaded []UploadedFile \`json:"uploaded"\`
\t\tOversizedFiles []string \`json:"oversized_files,omitempty"\`
\t\tSizeLimit     string   \`json:"size_limit,omitempty"\`
\t\tUploadErrors  []string \`json:"upload_errors,omitempty"\`
\t} \`json:"data"\`
}

type UploadedFile struct {
\tID           string \`json:"id"\`
\tURL          string \`json:"url"\`
\tThumbURL     string \`json:"thumb_url"\`
\tOriginalName string \`json:"original_name"\`
\tSize         int64  \`json:"size"\`
\tWidth        int    \`json:"width"\`
\tHeight       int    \`json:"height"\`
\tFormat       string \`json:"format"\`
\tAccessLevel  string \`json:"access_level"\`
\tCreatedAt    string \`json:"created_at"\`
}

type UploadOptions struct {
\tAPIKey     string
\tAccessLevel string // public/private/protected
\tOptimize   bool
\tFilePath   string
\tFolderID   string
}

type ImageUploader struct {
\tBaseURL    string
\tHTTPClient *http.Client
}

func NewImageUploader(baseURL string) *ImageUploader {
\treturn &ImageUploader{
\t\tBaseURL: baseURL,
\t\tHTTPClient: &http.Client{
\t\t\tTimeout: 30 * time.Second,
\t\t},
\t}
}

func (u *ImageUploader) UploadSingleFile(filePath string, options UploadOptions) (*UploadResponse, error) {
\t// 打开文件
\tfile, err := os.Open(filePath)
\tif err != nil {
\t\treturn nil, fmt.Errorf("无法打开文件: %v", err)
\t}
\tdefer file.Close()

\t// 创建multipart form
\tvar requestBody bytes.Buffer
\twriter := multipart.NewWriter(&requestBody)

\t// 添加文件字段
\tpart, err := writer.CreateFormFile("file", filepath.Base(filePath))
\tif err != nil {
\t\treturn nil, fmt.Errorf("创建文件字段失败: %v", err)
\t}

\t// 拷贝文件内容
\t_, err = io.Copy(part, file)
\tif err != nil {
\t\treturn nil, fmt.Errorf("拷贝文件内容失败: %v", err)
\t}

\t// 添加其他字段
\tif options.AccessLevel != "" {
\t\twriter.WriteField("access_level", options.AccessLevel)
\t}
\tif options.Optimize {
\t\twriter.WriteField("optimize", "true")
\t} else {
\t\twriter.WriteField("optimize", "false")
\t}
\tif options.FilePath != "" {
\t\twriter.WriteField("filePath", options.FilePath)
\t}
\tif options.FolderID != "" {
\t\twriter.WriteField("folderId", options.FolderID)
\t}

\t// 关闭writer
\terr = writer.Close()
\tif err != nil {
\t\treturn nil, fmt.Errorf("关闭writer失败: %v", err)
\t}

\t// 创建请求
\treq, err := http.NewRequest("POST", u.BaseURL+"/api/v1/external/upload", &requestBody)
\tif err != nil {
\t\treturn nil, fmt.Errorf("创建请求失败: %v", err)
\t}

\t// 设置请求头
\treq.Header.Set("Content-Type", writer.FormDataContentType())
\treq.Header.Set("x-pixelpunk-key", options.APIKey)

\t// 发送请求
\tresp, err := u.HTTPClient.Do(req)
\tif err != nil {
\t\treturn nil, fmt.Errorf("发送请求失败: %v", err)
\t}
\tdefer resp.Body.Close()

\t// 读取响应
\tbody, err := io.ReadAll(resp.Body)
\tif err != nil {
\t\treturn nil, fmt.Errorf("读取响应失败: %v", err)
\t}

\t// 解析响应
\tvar uploadResp UploadResponse
\terr = json.Unmarshal(body, &uploadResp)
\tif err != nil {
\t\treturn nil, fmt.Errorf("解析响应失败: %v", err)
\t}

\tif uploadResp.Code != 0 {
\t\treturn &uploadResp, fmt.Errorf("上传失败: %s", uploadResp.Message)
\t}

\treturn &uploadResp, nil
}

func (u *ImageUploader) UploadMultipleFiles(filePaths []string, options UploadOptions) (*UploadResponse, error) {
\tvar requestBody bytes.Buffer
\twriter := multipart.NewWriter(&requestBody)

\t// 添加多个文件
\tfor _, filePath := range filePaths {
\t\tfile, err := os.Open(filePath)
\t\tif err != nil {
\t\t\treturn nil, fmt.Errorf("无法打开文件 %s: %v", filePath, err)
\t\t}

\t\tpart, err := writer.CreateFormFile("files[]", filepath.Base(filePath))
\t\tif err != nil {
\t\t\tfile.Close()
\t\t\treturn nil, fmt.Errorf("创建文件字段失败: %v", err)
\t\t}

\t\t_, err = io.Copy(part, file)
\t\tfile.Close()
\t\tif err != nil {
\t\t\treturn nil, fmt.Errorf("拷贝文件内容失败: %v", err)
\t\t}
\t}

\t// 添加其他字段
\tif options.AccessLevel != "" {
\t\twriter.WriteField("access_level", options.AccessLevel)
\t}
\tif options.Optimize {
\t\twriter.WriteField("optimize", "true")
\t} else {
\t\twriter.WriteField("optimize", "false")
\t}
\tif options.FilePath != "" {
\t\twriter.WriteField("filePath", options.FilePath)
\t}
\tif options.FolderID != "" {
\t\twriter.WriteField("folderId", options.FolderID)
\t}

\terr := writer.Close()
\tif err != nil {
\t\treturn nil, fmt.Errorf("关闭writer失败: %v", err)
\t}

\t// 创建请求
\treq, err := http.NewRequest("POST", u.BaseURL+"/api/v1/external/upload", &requestBody)
\tif err != nil {
\t\treturn nil, fmt.Errorf("创建请求失败: %v", err)
\t}

\t// 设置请求头
\treq.Header.Set("Content-Type", writer.FormDataContentType())
\treq.Header.Set("x-pixelpunk-key", options.APIKey)

\t// 发送请求
\tresp, err := u.HTTPClient.Do(req)
\tif err != nil {
\t\treturn nil, fmt.Errorf("发送请求失败: %v", err)
\t}
\tdefer resp.Body.Close()

\t// 读取响应
\tbody, err := io.ReadAll(resp.Body)
\tif err != nil {
\t\treturn nil, fmt.Errorf("读取响应失败: %v", err)
\t}

\t// 解析响应
\tvar uploadResp UploadResponse
\terr = json.Unmarshal(body, &uploadResp)
\tif err != nil {
\t\treturn nil, fmt.Errorf("解析响应失败: %v", err)
\t}

\treturn &uploadResp, nil
}

func main() {
\t// 创建上传客户端
\tuploader := NewImageUploader("${currentDomain.value}")

\t// 配置上传选项
\toptions := UploadOptions{
\t\tAPIKey:     "YOUR_API_KEY",
\t\tVisibility: "public",
\t\tOptimize:   true,
\t\tFilePath:   "projects/website",
\t\tFolderID:   "", // 可选，优先级高于FilePath
\t}

\t// 单文件上传示例
\tfmt.Println("=== 单文件上传 ===")
\tsingleResult, err := uploader.UploadSingleFile("image.jpg", options)
\tif err != nil {
\t\tfmt.Printf("单文件上传失败: %v\\n", err)
\t} else {
\t\tfmt.Printf("单文件上传成功: %+v\\n", singleResult.Data.Uploaded)
\t}

\t// 批量上传示例
\tfmt.Println("\\n=== 批量上传 ===")
\tfilePaths := []string{"image1.jpg", "image2.png", "image3.gif"}
\tbatchResult, err := uploader.UploadMultipleFiles(filePaths, options)
\tif err != nil {
\t\tfmt.Printf("批量上传失败: %v\\n", err)
\t} else {
\t\tfmt.Printf("批量上传结果: %d个文件成功\\n", len(batchResult.Data.Uploaded))
\t\t
\t\t// 显示成功上传的文件
\t\tfor i, uploaded := range batchResult.Data.Uploaded {
\t\t\tfmt.Printf("[%d] %s -> %s\\n", i+1, uploaded.OriginalName, uploaded.URL)
\t\t}
\t\t
\t\t// 显示警告信息
\t\tif len(batchResult.Data.OversizedFiles) > 0 {
\t\t\tfmt.Printf("超限文件: %v\\n", batchResult.Data.OversizedFiles)
\t\t}
\t\tif len(batchResult.Data.UploadErrors) > 0 {
\t\t\tfmt.Printf("上传错误: %v\\n", batchResult.Data.UploadErrors)
\t\t}
\t}
}

func ConcurrentUpload(uploader *ImageUploader, filePaths []string, options UploadOptions) {
\ttype uploadResult struct {
\t\tFilePath string
\t\tResult   *UploadResponse
\t\tError    error
\t}

\tresultChan := make(chan uploadResult, len(filePaths))

\t// 启动goroutine进行并发上传
\tfor _, filePath := range filePaths {
\t\tgo func(fp string) {
\t\t\tresult, err := uploader.UploadSingleFile(fp, options)
\t\t\tresultChan <- uploadResult{
\t\t\t\tFilePath: fp,
\t\t\t\tResult:   result,
\t\t\t\tError:    err,
\t\t\t}
\t\t}(filePath)
\t}

\t// 收集结果
\tfor i := 0; i < len(filePaths); i++ {
\t\tresult := <-resultChan
\t\tif result.Error != nil {
\t\t\tfmt.Printf("文件 %s 上传失败: %v\\n", result.FilePath, result.Error)
\t\t} else {
\t\t\tfmt.Printf("文件 %s 上传成功: %s\\n",
\t\t\t\tresult.FilePath,
\t\t\t\tresult.Result.Data.Uploaded[0].URL)
\t\t}
\t}
}`,
    },
  }))

  const jsonResponseExamples = computed(() => ({
    single: `{
  "code": 0,
  "message": "上传成功",
  "data": {
    "uploaded": {
      "id": "img_1a2b3c4d5e",
      "url": "${currentDomain.value}/file/user_1/image.jpg",
      "thumb_url": "${currentDomain.value}/thumb/user_1/image.jpg",
      "original_name": "image.jpg",
      "size": 1024000,
      "width": 1920,
      "height": 1080,
      "format": "jpg",
      "access_level": "private",
      "created_at": "2024-01-15T08:30:45Z"
    }
  }
}`,
    batch: `{
  "code": 0,
  "message": "部分文件上传成功",
  "data": {
    "uploaded": [
      {
        "id": "img_1a2b3c4d5e",
        "url": "${currentDomain.value}/file/user_1/image1.jpg",
        "thumb_url": "${currentDomain.value}/thumb/user_1/image1.jpg",
        "original_name": "image1.jpg",
        "size": 1024000,
        "width": 1920,
        "height": 1080,
        "format": "jpg",
        "access_level": "private",
        "created_at": "2024-01-15T08:30:45Z"
      }
    ],
    "oversized_files": ["large_image.jpg"],
    "size_limit": "20.0MB",
    "upload_errors": ["corrupted.jpg: 文件损坏无法处理"]
  }
}`,
  }))

  return {
    codeExamples,
    jsonResponseExamples,
  }
}
