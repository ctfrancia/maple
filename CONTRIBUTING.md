# Contributing to [Project Name]

Thank you for your interest in contributing to this project! We welcome contributions from everyone and appreciate your help in making this project better.

*NOTE/IMPORTANT*: this is a work in progress, so please be patient with me as I am working on this project and will be finalized when ready to accept contributing

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Coding Guidelines](#coding-guidelines)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to [email@example.com].

## Getting Started

Before you begin:
- Have you read the [README](README.md)?
- Check out the [existing issues](../../issues)
- Look at the [project roadmap](../../projects) (if available)

## How to Contribute

There are many ways to contribute to this project:

### 🐛 Reporting Bugs
- Use the GitHub issue tracker
- Check if the bug has already been reported
- Include detailed steps to reproduce
- Provide system information and error messages

### 💡 Suggesting Features
- Open an issue with the "enhancement" label
- Clearly describe the feature and its benefits
- Consider if it fits the project's scope and goals

### 📖 Improving Documentation
- Fix typos and clarify existing documentation
- Add examples and use cases
- Translate documentation to other languages

### 🔧 Code Contributions
- Fix bugs
- Implement new features
- Improve performance
- Add tests

## Development Setup

### Prerequisites
- Go 1.23.4 or higher
- Git
- Make (optional, for using Makefile commands)

### Installation
1. Fork the repository
2. Clone your fork:
   ```bash
   git clone 
   cd project-name
   ```
3. Initialize Go modules (if not already done):
   ```bash
   go mod download
   ```
4. Set up development environment:
   ```bash
   # Copy environment file if needed
   cp .env.example .env
   
   # Install development tools
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run specific package tests
go test ./pkg/packagename
```

### Running the Project Locally
```bash
# Build the project
go build -o bin/app ./cmd/main.go

# Run directly with go run
go run ./cmd/main.go

# Or if using Make
make build
make run
```

## Coding Guidelines

### Style Guide
- Follow [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- Use [Effective Go](https://go.dev/doc/effective_go) guidelines
- Follow standard Go naming conventions
- Use meaningful variable and function names
- Write clear, concise comments for exported functions and types
- Keep functions small and focused
- Use Go idioms and standard patterns

### Best Practices
- Write tests for new features and bug fixes (aim for >80% coverage)
- Use table-driven tests when appropriate
- Handle errors explicitly - don't ignore them
- Use context.Context for cancellation and timeouts
- Ensure backward compatibility when possible
- Update documentation for any changes
- Use descriptive commit messages

### Code Formatting
Go has built-in formatting tools:
```bash
# Format code (automatically fixes formatting)
go fmt ./...

# Organize imports
goimports -w .

# Or use gofmt with specific options
gofmt -s -w .
```

### Linting
We use golangci-lint for comprehensive code analysis:
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run with specific configuration
golangci-lint run --config .golangci.yml
```

### Testing Guidelines
- Place tests in the same package as the code being tested
- Use `_test.go` suffix for test files
- Use table-driven tests for multiple test cases
- Use testify/assert for assertions (if using testify)
- Mock external dependencies
- Use `go test -race` to detect race conditions

Example test structure:
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "result",
            wantErr:  false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("FunctionName() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Commit Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification:

### Commit Message Format
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

### Examples
```
feat(api): add user authentication endpoint

fix(parser): handle edge case with empty input

docs: update installation instructions

test(utils): add unit tests for helper functions
```

## Pull Request Process

1. **Create a branch** from `main` for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-number
   ```

2. **Make your changes** following the coding guidelines

3. **Test your changes** thoroughly:
   ```bash
   # Run tests
   go test ./...
   
   # Run tests with race detection
   go test -race ./...
   
   # Check formatting
   go fmt ./...
   
   # Run linter
   golangci-lint run
   
   # Check for common Go issues
   go vet ./...
   ```

4. **Commit your changes** using conventional commit format

5. **Push to your fork**:
   ```bash
   git push origin your-branch-name
   ```

6. **Create a Pull Request**:
   - Use a clear and descriptive title
   - Fill out the PR template completely
   - Link related issues using "Fixes #123" or "Closes #123"
   - Include screenshots for UI changes
   - Ensure all checks pass

### PR Requirements
- [ ] Tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] No race conditions (`go test -race ./...`)
- [ ] Code follows Go best practices
- [ ] Documentation is updated (including godoc comments)
- [ ] Commit messages follow conventions
- [ ] No merge conflicts
- [ ] Reviewed by at least one maintainer

## Issue Guidelines

### Before Creating an Issue
- Search existing issues to avoid duplicates
- Check if it's already fixed in the latest version
- Gather relevant information and examples

### Bug Reports Should Include
- Clear and descriptive title
- Steps to reproduce the issue
- Expected vs actual behavior
- Go version (go version)
- Operating system and architecture
- Relevant environment variables
- Error messages and logs
- Screenshots if applicable

### Feature Requests Should Include
- Clear and descriptive title
- Use case and motivation
- Detailed description of the proposed feature
- Consider alternative solutions
- Mockups or examples if applicable

## Community

### Getting Help
- Check the [documentation](docs/)
- Search [existing issues](../../issues)
- Ask questions in [discussions](../../discussions) or [community forum]
- Join our [Discord/Slack/chat platform]

### Recognition
Contributors are recognized in our:
- [Contributors list](CONTRIBUTORS.md)
- Release notes
- Project documentation

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project. See [LICENSE](LICENSE) for details.

---

**Thank you for contributing!** 🎉

Your contributions help make this project better for everyone. We appreciate your time and effort!
