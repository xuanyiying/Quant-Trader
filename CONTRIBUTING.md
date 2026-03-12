# Contributing to Quant-Trader

[English](#english) | [简体中文](#中文)

---

# English

Thank you for your interest in contributing to Quant-Trader! We welcome contributions from the community to help make this project better.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Code of Conduct](#code-of-conduct)
- [Ways to Contribute](#ways-to-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Pull Request Process](#pull-request-process)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)
- [Getting Help](#getting-help)

---

## Quick Start

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/your-username/quant-trader.git`
3. **Create** a feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes
5. **Test** your changes
6. **Commit** with descriptive messages
7. **Push** to your fork
8. **Submit** a Pull Request

---

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](https://github.com/your-repo/quant-trader/blob/main/CODE_OF_CONDUCT.md). Please be respectful and constructive.

---

## Ways to Contribute

### Bug Reports

- Use GitHub Issues to report bugs
- Include detailed reproduction steps
- Attach relevant logs and screenshots
- Use the bug report template

### Feature Requests

- Open an issue with the `enhancement` tag
- Describe the feature and its use cases
- Explain why this feature would be beneficial

### Code Contributions

- Bug fixes
- New features
- Performance improvements
- Documentation improvements

---

## Development Setup

### Prerequisites

- Go 1.24+
- Node.js 18+ (for frontend)
- Docker & Docker Compose
- PostgreSQL 15+
- NATS Server

### Backend Setup

```bash
# Clone the repository
git clone https://github.com/your-repo/quant-trader.git
cd quant-trader/backend

# Install dependencies
go mod download

# Run tests
go test ./...

# Start development server
go run cmd/main.go
```

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

---

## Coding Standards

### Go (Backend)

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for code formatting
- Run `go vet` before committing
- Add unit tests for new features
- Use meaningful variable and function names
- Keep functions focused and small

```bash
# Format code
gofmt -w .

# Run linter
go vet ./...

# Run all tests with coverage
go test -v -cover ./...
```

### React/TypeScript (Frontend)

- Follow ESLint rules
- Use TypeScript for type safety
- Use functional components with hooks
- Keep components small and focused

```bash
# Check linting
npm run lint

# Run type checking
npm run typecheck
```

---

## Pull Request Process

### 1. Before Submitting

- [ ] Code follows our coding standards
- [ ] Tests pass locally
- [ ] Documentation is updated (if needed)
- [ ] Commits are clean and descriptive

### 2. PR Description

Include in your PR description:

```
## Changes Made
- Brief description of changes

## Testing
- How did you test these changes?

## Screenshots (if applicable)
- Add screenshots for UI changes
```

### 3. PR Review

- At least one maintainer approval required
- Address feedback promptly
- Keep PRs focused and atomic

---

## Commit Message Guidelines

Use conventional commit format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Code restructuring
- `test`: Tests
- `chore`: Maintenance

### Examples

```
feat(matching): add order book depth snapshot
fix(connector): resolve reconnection timeout issue
docs(readme): update installation instructions
refactor(engine): simplify order processing logic
```

---

## Testing Requirements

### Backend Tests

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/matching/...

# Run with coverage
go test -cover ./...
```

### Frontend Tests

```bash
# Run all tests
npm test

# Run with coverage
npm run test:coverage
```

---

## Documentation

- Update README.md for user-facing changes
- Add code comments for complex logic
- Update API documentation for endpoint changes
- Add examples for new features

---

## Getting Help

- **GitHub Issues**: For bug reports and feature requests
- **Discussions**: For questions and community support
- **Discord**: Join our community server

---

## Contributors

Thank you to all our contributors!

<a href="https://github.com/your-repo/quant-trader/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=your-repo/quant-trader" />
</a>

---

## License

By contributing to Quant-Trader, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

# 中文

感谢您对 Quant-Trader 的兴趣！我们欢迎社区贡献，帮助改进这个项目。

---

## 目录

- [快速开始](#快速开始)
- [行为准则](#行为准则)
- [贡献方式](#贡献方式)
- [开发环境设置](#开发环境设置)
- [编码规范](#编码规范)
- [Pull Request 流程](#pull-request-流程)
- [提交信息规范](#提交信息规范)
- [测试要求](#测试要求)
- [文档](#文档-1)
- [获取帮助](#获取帮助)

---

## 快速开始

1. **Fork** 本仓库
2. **克隆** 您的 Fork: `git clone https://github.com/your-username/quant-trader.git`
3. **创建** 功能分支: `git checkout -b feature/amazing-feature`
4. **进行** 您的修改
5. **测试** 您的修改
6. **提交** 带有描述性信息的提交
7. **推送** 到您的 Fork
8. **提交** Pull Request

---

## 行为准则

参与本项目即表示您同意遵守我们的[行为准则](https://github.com/your-repo/quant-trader/blob/main/CODE_OF_CONDUCT.md)。请保持尊重和建设性。

---

## 贡献方式

### 报告 Bug

- 使用 GitHub Issues 报告 Bug
- 包含详细的复现步骤
- 附上相关日志和截图
- 使用 Bug 报告模板

### 功能请求

- 创建带有 `enhancement` 标签的 Issue
- 描述功能及其用例
- 解释为什么这个功能会有益

### 代码贡献

- Bug 修复
- 新功能
- 性能改进
- 文档改进

---

## 开发环境设置

### 前置要求

- Go 1.24+
- Node.js 18+（前端）
- Docker & Docker Compose
- PostgreSQL 15+
- NATS Server

### 后端设置

```bash
# 克隆仓库
git clone https://github.com/your-repo/quant-trader.git
cd quant-trader/backend

# 安装依赖
go mod download

# 运行测试
go test ./...

# 启动开发服务器
go run cmd/main.go
```

### 前端设置

```bash
cd frontend
npm install
npm run dev
```

---

## 编码规范

### Go（后端）

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `gofmt` 格式化代码
- 提交前运行 `go vet`
- 为新功能添加单元测试
- 使用有意义的变量和函数名
- 保持函数专注且精简

```bash
# 格式化代码
gofmt -w .

# 运行静态检查
go vet ./...

# 运行所有测试并生成覆盖率报告
go test -v -cover ./...
```

### React/TypeScript（前端）

- 遵循 ESLint 规则
- 使用 TypeScript 确保类型安全
- 使用函数组件和 Hooks
- 保持组件小巧专注

```bash
# 检查代码规范
npm run lint

# 运行类型检查
npm run typecheck
```

---

## Pull Request 流程

### 1. 提交前

- [ ] 代码符合编码规范
- [ ] 本地测试通过
- [ ] 文档已更新（如需要）
- [ ] 提交信息清晰描述

### 2. PR 描述

在 PR 描述中包含：

```
## 变更内容
- 变更的简要描述

## 测试
- 如何测试这些变更？

## 截图（如适用）
- 添加 UI 变更的截图
```

### 3. PR 审核

- 需要至少一名维护者批准
- 及时回应反馈
- 保持 PR 专注和原子性

---

## 提交信息规范

使用约定式提交格式：

```
<类型>(<范围>): <描述>

[可选正文]

[可选脚注]
```

### 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档
- `style`: 格式
- `refactor`: 代码重构
- `test`: 测试
- `chore`: 维护

### 示例

```
feat(matching): 添加订单簿深度快照功能
fix(connector): 解决重连超时问题
docs(readme): 更新安装说明
refactor(engine): 简化订单处理逻辑
```

---

## 测试要求

### 后端测试

```bash
# 运行所有测试
go test ./...

# 运行特定包
go test ./internal/matching/...

# 生成覆盖率报告
go test -cover ./...
```

### 前端测试

```bash
# 运行所有测试
npm test

# 生成覆盖率报告
npm run test:coverage
```

---

## 文档

- 面向用户的变更更新 README.md
- 为复杂逻辑添加代码注释
- 端点变更更新 API 文档
- 为新功能添加示例

---

## 获取帮助

- **GitHub Issues**: 用于 Bug 报告和功能请求
- **Discussions**: 用于问题和社区支持
- **Discord**: 加入我们的社区服务器

---

## 贡献者

感谢所有贡献者！

<a href="https://github.com/your-repo/quant-trader/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=your-repo/quant-trader" />
</a>

---

## 许可证

向 Quant-Trader 贡献代码即表示您同意将您的贡献在 [MIT 许可证](LICENSE) 下授权。
