# Contributing to GoScraper

Thank you for your interest in contributing to GoScraper! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional, for testing)

### Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/ramusaaa/goscraper.git
   cd goscraper
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run tests:
   ```bash
   go test -v ./...
   ```

5. Run benchmarks:
   ```bash
   go test -bench=. -benchmem
   ```

## 🛠️ Development Guidelines

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write clear, self-documenting code
- Add comments for complex logic

### Testing

- Write tests for new features
- Maintain test coverage above 80%
- Include benchmarks for performance-critical code
- Test both success and error cases

### Commit Messages

Use conventional commit format:
```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/tooling changes

Examples:
```
feat(scraper): add retry mechanism with exponential backoff
fix(parser): handle malformed HTML gracefully
docs(readme): update installation instructions
```

## 📝 Pull Request Process

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -m "feat: add your feature"
   ```

3. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

4. Create a Pull Request with:
   - Clear title and description
   - Reference any related issues
   - Include tests for new functionality
   - Update documentation if needed

### PR Requirements

- [ ] Tests pass (`go test ./...`)
- [ ] Code is formatted (`gofmt`)
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
- [ ] Benchmarks included for performance changes

## 🐛 Bug Reports

When reporting bugs, please include:

- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant code snippets
- Error messages/logs

Use the bug report template in GitHub Issues.

## 💡 Feature Requests

For new features:

- Check existing issues first
- Describe the use case
- Explain why it's valuable
- Consider implementation complexity
- Be open to discussion

## 🏗️ Architecture

### Project Structure

```
goscraper/
├── pkg/                    # Core packages
│   ├── ai/                # AI extraction
│   ├── browser/           # Browser automation
│   ├── cache/             # Caching layer
│   ├── cluster/           # Distributed coordination
│   ├── monitoring/        # Metrics & observability
│   └── queue/             # Message queuing
├── cmd/                   # Applications
├── examples/              # Usage examples
├── docs/                  # Documentation
└── tests/                 # Integration tests
```

### Key Principles

- **Modularity**: Each package has a single responsibility
- **Interfaces**: Use interfaces for testability and flexibility
- **Error Handling**: Explicit error handling with context
- **Performance**: Optimize for high throughput and low latency
- **Observability**: Comprehensive metrics and logging

## 🔒 Security

- Report security issues privately to security@goscraper.com
- Follow responsible disclosure practices
- Security fixes get priority review

## 📚 Documentation

- Update README.md for user-facing changes
- Add godoc comments for public APIs
- Include examples in documentation
- Update architecture docs for significant changes

## 🎯 Areas for Contribution

### High Priority

- [ ] Additional browser engines
- [ ] More AI model integrations
- [ ] Performance optimizations
- [ ] Better error handling
- [ ] Documentation improvements

### Medium Priority

- [ ] Additional cache backends
- [ ] More queue systems
- [ ] Enhanced monitoring
- [ ] CLI improvements
- [ ] Example applications

### Low Priority

- [ ] UI/Dashboard
- [ ] Additional parsers
- [ ] Plugin system
- [ ] Mobile support

## 🏆 Recognition

Contributors are recognized in:

- README.md contributors section
- Release notes
- GitHub contributors page
- Special thanks in documentation

## 📞 Getting Help

- GitHub Discussions for questions
- GitHub Issues for bugs/features
- Code review feedback in PRs
- Community Discord (coming soon)

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to GoScraper! 🚀
