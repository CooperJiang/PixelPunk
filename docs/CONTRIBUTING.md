# 贡献指南

感谢你有兴趣为 PixelPunk 贡献代码！

## 🎯 贡献方式

### 1. 报告 Bug 🐛

在提交 Bug 前，请先：
- 搜索 [现有 Issues](https://github.com/CooperJiang/PixelPunk-v1/issues)，避免重复
- 确认问题可复现
- 准备好复现步骤

**Bug 报告应包含**：
- **环境信息**：操作系统、浏览器版本、PixelPunk 版本
- **复现步骤**：详细的步骤说明
- **预期行为**：应该发生什么
- **实际行为**：实际发生了什么
- **截图/日志**：相关的截图或错误日志

[→ 报告 Bug](https://github.com/CooperJiang/PixelPunk-v1/issues/new?labels=bug)

### 2. 提出新功能 💡

我们欢迎功能建议！在提交前：
- 检查功能是否已在 [Roadmap](https://github.com/CooperJiang/PixelPunk-v1/issues?q=label%3Aenhancement) 中
- 说明功能的使用场景
- 描述期望的实现方式

**功能请求应包含**：
- **问题描述**：当前遇到的问题
- **解决方案**：建议的功能设计
- **使用场景**：谁会使用这个功能，如何使用
- **替代方案**：是否考虑过其他实现方式

[→ 提出新功能](https://github.com/CooperJiang/PixelPunk-v1/issues/new?labels=enhancement)

### 3. 改进文档 📖

文档的贡献同样重要！你可以：
- 修正错别字和语法错误
- 补充缺失的文档
- 改进文档结构和可读性
- 翻译文档到其他语言

[→ 改进文档](https://github.com/CooperJiang/PixelPunk-v1/issues/new?labels=documentation)

### 4. 提交代码 💻

请按照以下流程提交代码贡献。

---

## 🚀 代码贡献流程

### 步骤 1: Fork 仓库

点击项目页面右上角的 **Fork** 按钮，将项目复制到你的 GitHub 账户。

### 步骤 2: 克隆到本地

```bash
git clone https://github.com/YOUR_USERNAME/PixelPunk-v1.git
cd PixelPunk-v1

# 添加上游仓库
git remote add upstream https://github.com/CooperJiang/PixelPunk-v1.git
```

### 步骤 3: 创建分支

```bash
# 更新主分支
git checkout main
git pull upstream main

# 创建功能分支
git checkout -b feature/your-feature-name
# 或修复分支
git checkout -b fix/your-bug-fix
```

**分支命名规范**：
- `feature/xxx` - 新功能
- `fix/xxx` - Bug 修复
- `docs/xxx` - 文档改进
- `refactor/xxx` - 代码重构
- `style/xxx` - 样式调整
- `test/xxx` - 测试相关

### 步骤 4: 搭建开发环境

详见 [开发文档](./DEVELOPMENT.md)。

```bash
# 1. 配置文件
cp configs/config.example.yaml configs/config.yaml

# 2. 启动向量数据库
make qdrant-start

# 3. 启动后端（新终端）
make dev

# 4. 启动前端（新终端）
make dev-frontend
```

### 步骤 5: 编写代码

**编码规范**：

#### 后端（Go）
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 规范
- 使用 `gofmt` 格式化代码
- 函数和方法添加注释
- 单元测试覆盖关键逻辑

```go
// ✅ 好的示例
// UploadFile 上传文件到服务器
// 参数：file - 文件数据，userID - 用户ID
// 返回：文件ID 和错误信息
func UploadFile(file *multipart.FileHeader, userID uint) (string, error) {
    // 实现逻辑
}

// ❌ 不好的示例
func upload(f *multipart.FileHeader, u uint) (string, error) {
    // 缺少注释，命名不清晰
}
```

#### 前端（Vue 3 + TypeScript）
- 使用 `ESLint` 和 `Prettier` 检查代码
- 组件使用 `<script setup>` 语法
- Props 和 Emits 定义类型
- 避免使用 `any` 类型

```vue
<!-- ✅ 好的示例 -->
<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  title: string
  count?: number
}

interface Emits {
  (e: 'update', value: number): void
}

const props = withDefaults(defineProps<Props>(), {
  count: 0
})

const emit = defineEmits<Emits>()
</script>

<!-- ❌ 不好的示例 -->
<script setup>
const props = defineProps(['title', 'count']) // 缺少类型
</script>
```

### 步骤 6: 提交代码

```bash
# 添加文件
git add .

# 提交（遵循提交规范）
git commit -m "feat: 添加图片批量下载功能"
```

**提交信息规范**（遵循 [Conventional Commits](https://www.conventionalcommits.org/)）：

```
<类型>: <简短描述>

[可选的详细描述]

[可选的 Issue 引用]
```

**类型**：
- `feat` - 新功能
- `fix` - Bug 修复
- `docs` - 文档更新
- `style` - 代码格式调整（不影响功能）
- `refactor` - 重构代码
- `perf` - 性能优化
- `test` - 测试相关
- `chore` - 构建工具、依赖更新

**示例**：
```bash
# 新功能
git commit -m "feat: 添加图片水印功能"

# Bug 修复
git commit -m "fix: 修复上传大文件时内存溢出问题"

# 文档
git commit -m "docs: 更新部署文档中的 Docker 配置"

# 带详细说明
git commit -m "feat: 支持 WebP 图片格式

- 添加 WebP 解码器
- 支持 WebP 格式转换
- 更新上传白名单

Closes #123"
```

### 步骤 7: 推送分支

```bash
git push origin feature/your-feature-name
```

### 步骤 8: 创建 Pull Request

1. 访问你的 Fork 仓库页面
2. 点击 **Compare & pull request** 按钮
3. 填写 PR 标题和描述

**PR 描述模板**：

```markdown
## 变更类型
- [ ] 新功能
- [ ] Bug 修复
- [ ] 文档更新
- [ ] 代码重构
- [ ] 性能优化

## 变更说明
简要说明这个 PR 做了什么。

## 相关 Issue
Closes #123

## 测试
- [ ] 已通过本地测试
- [ ] 已添加单元测试
- [ ] 已测试不同浏览器

## 截图（如适用）
贴上相关截图。

## Checklist
- [ ] 代码遵循项目规范
- [ ] 已更新相关文档
- [ ] 已添加必要的注释
- [ ] 无明显的性能问题
```

### 步骤 9: 代码审查

维护者会审查你的代码，可能会：
- 提出修改建议
- 请求补充测试
- 要求优化实现

请及时回应审查意见，并根据反馈修改代码：

```bash
# 继续在功能分支上修改
git add .
git commit -m "fix: 根据审查意见优化代码"
git push origin feature/your-feature-name
```

### 步骤 10: 合并

审查通过后，维护者会合并你的 PR。恭喜你成为贡献者！🎉

---

## 📋 开发指南

### 项目结构

```
PixelPunk-v1/
├── cmd/                    # 程序入口
├── internal/               # 内部代码
│   ├── controllers/        # 控制器
│   ├── services/           # 业务逻辑
│   ├── models/             # 数据模型
│   └── middleware/         # 中间件
├── pkg/                    # 公共包
├── web/                    # 前端代码
│   ├── src/
│   │   ├── pages/          # 页面
│   │   ├── components/     # 组件
│   │   ├── composables/    # 组合式函数
│   │   └── stores/         # 状态管理
├── configs/                # 配置文件
├── docs/                   # 文档
└── scripts/                # 脚本
```

### 代码检查

**后端**：
```bash
# 格式化
gofmt -w .

# 静态检查
go vet ./...

# 运行测试
go test ./...
```

**前端**：
```bash
cd web

# 代码检查
pnpm lint:check

# 自动修复
pnpm lint:fix

# 类型检查
pnpm run type-check
```

### 调试技巧

**后端调试**：
- 使用 `make dev` 启动热重载
- 查看日志：`tail -f logs/app.log`
- 使用 Postman 测试 API

**前端调试**：
- 使用浏览器开发者工具
- 安装 Vue DevTools 扩展
- 查看控制台错误

---

## 🤝 行为准则

### 我们的承诺

为了营造开放和友好的环境，我们承诺：
- 尊重不同的观点和经验
- 接受建设性的批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

### 不可接受的行为

- 使用性别化语言或图像
- 发表侮辱性/贬损性评论
- 骚扰行为（公开或私下）
- 未经明确许可发布他人私人信息
- 其他不道德或不专业的行为

### 举报

如遇到不当行为，请通过 Issue 或邮件联系维护者。

---

## 💬 交流渠道

- **GitHub Issues**: 报告 Bug 和功能请求
- **GitHub Discussions**: 一般性讨论和问答
- **Pull Request**: 代码贡献

---

## 📚 其他资源

- [开发文档](./DEVELOPMENT.md) - 环境搭建
- [部署文档](./DEPLOYMENT.md) - 生产部署
- [主题开发](../web/docs/THEME-DEVELOPMENT-GUIDE.md) - 自定义主题
- [国际化](../web/docs/I18N-GUIDE.md) - 添加语言

---

## 🎖️ 贡献者

感谢所有贡献者！

<!--
可以使用 https://contrib.rocks/ 生成贡献者墙
![Contributors](https://contrib.rocks/image?repo=CooperJiang/PixelPunk-v1)
-->

---

再次感谢你的贡献！每一个 PR 都让 PixelPunk 变得更好。❤️
