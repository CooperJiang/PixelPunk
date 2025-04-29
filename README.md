# PixelPunk - 像素朋克图床

<div align="center">
  
![PixelPunk Logo](/image/logo.png)

_🔥 一个赛博朋克风格的AI智能图床系统 🔥_

</div>

<p align="center">
  <a href="https://github.com/CooperJiang/PixelPunk"><img src="https://img.shields.io/badge/官方仓库-PixelPunk-blue" alt="官方仓库"></a>
  <a href="http://cyber.mmmss.com/"><img src="https://img.shields.io/badge/官方网站-cyber.mmmss.com-ff69b4" alt="官方网站"></a>
  <a href="#许可证"><img src="https://img.shields.io/badge/license-免费使用-green" alt="license"></a>
</p>

<p align="center">
  <a href="#项目介绍">项目介绍</a> •
  <a href="#特性">特性</a> •
  <a href="#演示">演示</a> •
  <a href="#部署">部署</a> •
  <a href="#使用">使用</a> •
  <a href="#技术栈">技术栈</a> •
  <a href="#联系我们">联系我们</a> •
  <a href="#许可证">许可证</a>
</p>

## 项目介绍

PixelPunk（像素朋克）是一款拥有赛博朋克视觉风格的AI智能图床系统。本项目结合了现代前端技术和AI图像处理能力，为用户提供一个功能强大、界面酷炫、使用简单的图片存储管理平台。

无论是个人用户还是企业用户，PixelPunk都能为您提供丰富的图片管理功能，包括多种存储方式、AI图像识别、智能标签生成、文件夹管理以及权限控制等。

## 特性

### 🚀 核心功能

- **多种存储方式**：支持本地存储、腾讯云COS、阿里云OSS、七牛云、AWS S3等多种存储方式
- **用户文件隔离**：用户权限控制，文件隔离，不同用户文件分寸不同目录
- **文件夹管理**：完善的文件夹管理系统，支持无限嵌套文件夹创建
- **多种上传方式**：支持拖拽上传、剪贴板上传、URL上传等多种方式
- **AI智能辅助**：可选接入AI，智能打标，智能总结，智能审核，快速定位查找图片

### 🤖 AI智能功能

- **智能图像识别**：自动分析图片内容并生成相关标签
- **内容搜索**：基于图片内容的智能搜索功能
- **图片审核**：智能检测不良图片内容，保障平台安全
- **智能分类**：自动对上传图片进行分类，提升管理效率

### 🎨 用户体验

- **赛博朋克UI**：独特的视觉设计，带来沉浸式赛博朋克体验
- **响应式设计**：完美支持各种设备，从手机到桌面设备
- **自定义分享**：支持自定义图片分享样式和水印
- **快捷操作**：便捷的快捷键和批量操作功能

## 演示

### 在线体验

[官方网站](http://cyber.mmmss.com/)

### 系统截图展示

<div align="center">
  <h4>赛博朋克风格界面</h4>
  <img src="" alt="系统界面总览" width="800px"/>
  <p><i>赛博朋克风格的用户界面设计</i></p>
</div>

#### 核心功能展示

<table>
  <tr>
    <td width="50%">
      <img src="/image/upload.png" alt="图片上传界面"/>
      <p align="center"><b>多样化上传</b> - 拖拽/剪贴板/URL多种上传方式</p>
    </td>
    <td width="50%">
      <img src="/image/manage.png" alt="图片管理界面"/>
      <p align="center"><b>智能管理</b> - 文件夹组织与批量操作</p>
    </td>
  </tr>
  <tr>
    <td width="50%">
      <img src="/image/aiTag.png" alt="AI智能分类与标签"/>
      <p align="center"><b>AI标签</b> - 智能识别图片内容生成标签</p>
    </td>
    <td width="50%">
      <img src="/image/channels.png" alt="多渠道知道"/>
      <p align="center"><b>多渠道知道</b> - 支持同时混用不同渠道</p>
    </td>
  </tr>
</table>

#### 管理后台

<div align="center">
  <img src="/image/admin.png" alt="管理后台界面" width="800px"/>
  <p><i>强大的管理后台，轻松配置系统参数</i></p>
</div>

<!-- ### 功能演示动图

<div align="center">
  <h4>上传与AI识别流程</h4>
  <img src="" alt="上传与AI识别过程" width="800px"/>
  <p><i>图片上传到自动识别标签的完整流程</i></p>
</div> -->

## 部署

PixelPunk 图床系统使用 Go 编写的后端，前端文件已打包并编译为单一二进制文件，可以快速部署到您的服务器上。

### 系统要求

- 操作系统：Linux / Windows / macOS
- 内存：建议 2GB+
- 硬盘：根据存储需求，建议 10GB+

### 部署方式

#### 一键安装脚本（即将支持）

```bash
# 即将提供一键安装脚本
```

#### Docker 安装（即将支持）

```bash
# 即将提供 Docker 安装方式
```

#### 手动安装

1. 从 [GitHub 仓库](https://github.com/CooperJiang/PixelPunk) 下载最新的发布版本
2. 解压下载的文件
3. 赋予可执行权限
   ```bash
   chmod +x pixelpunk
   ```
4. 启动服务
   ```bash
   ./pixelpunk
   ```
5. 访问 `http://您的服务器IP:端口` 开始使用

## 使用

### 系统界面预览

<div align="center">
  <table>
    <tr>
      <td><img src="" alt="登录界面"/></td>
      <td><img src="" alt="仪表盘"/></td>
    </tr>
  </table>
</div>

### 基本配置

在使用前，需要进行基本配置：

1. 配置存储方式（本地/云存储）
2. 设置用户权限和注册选项
3. 配置AI功能（可选）

<div align="center">
  <img src="" alt="存储配置界面" width="600px"/>
  <p><i>存储服务配置界面</i></p>
</div>

详细配置请参考[配置文档]()。

### 快速上传图片

1. 登录系统
2. 进入上传页面
3. 拖拽图片或点击上传按钮
4. 获取图片链接并使用
5. ai将在后台对ai智能分类打标，并进行nsfw安全检测

<div align="center">
  <img src="" alt="上传流程演示" width="600px"/>
  <p><i>简单直观的上传流程</i></p>
</div>

## 技术栈

- **后端框架**: Go + Gin + Mysql
- **前端框架**: Vue 3 + TypeScript
- **UI组件**: 自定义赛博朋克风格组件

## 联系我们

- 官方网站: [http://cyber.mmmss.com/](http://cyber.mmmss.com/)
- 官方仓库: [https://github.com/CooperJiang/PixelPunk](https://github.com/CooperJiang/PixelPunk)
- QQ群: 826708512
- 微信: J_longyan

## 许可证

PixelPunk 图床系统为免费软件，您可以在遵循以下条款的前提下自由使用：

1. 禁止对程序进行逆向工程、反编译或以任何方式尝试提取源代码
2. 禁止删除或修改软件中的版权信息和归属声明
3. 禁止将本软件用于任何违法用途

## 致谢

- 感谢所有为此项目提供反馈和建议的用户