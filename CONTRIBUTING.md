# Contributing to Quant-Trader

Thank you for your interest in contributing to Quant-Trader! We welcome contributions from the community to help make this project better.

---

## 🚀 Quick Start

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/your-username/quant-trader.git`
3. **Create** a feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes
5. **Test** your changes
6. **Commit** with descriptive messages
7. **Push** to your fork
8. **Submit** a Pull Request

---

## 📋 Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](https://github.com/your-repo/quant-trader/blob/main/CODE_OF_CONDUCT.md). Please be respectful and constructive.

---

## 🎯 Ways to Contribute

### 🐛 Bug Reports

- Use GitHub Issues to report bugs
- Include detailed reproduction steps
- Attach relevant logs and screenshots
- Use the bug report template

### 💡 Feature Requests

- Open an issue with the `enhancement` tag
- Describe the feature and its use cases
- Explain why this feature would be beneficial

### 🧑‍💻 Code Contributions

- Bug fixes
- New features
- Performance improvements
- Documentation improvements

---

## 🛠 Development Setup

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

## 📝 Coding Standards

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

## 🔄 Pull Request Process

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

## 🏷️ Commit Message Guidelines

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

## 🧪 Testing Requirements

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

## 📚 Documentation

- Update README.md for user-facing changes
- Add code comments for complex logic
- Update API documentation for endpoint changes
- Add examples for new features

---

## 💬 Getting Help

- **GitHub Issues**: For bug reports and feature requests
- **Discussions**: For questions and community support
- **Discord**: Join our community server

---

## 🏆 Contributors

Thank you to all our contributors!

<a href="https://github.com/your-repo/quant-trader/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=your-repo/quant-trader" />
</a>

---

## 📄 License

By contributing to Quant-Trader, you agree that your contributions will be licensed under the [MIT License](LICENSE).
