# 🤝 Contributing to MCP RAPI

Thank you for your interest in contributing to **MCP RAPI**! We welcome contributions from developers of all skill levels. This guide will help you get started.

## 📋 Table of Contents

- [🎯 Ways to Contribute](#-ways-to-contribute)
- [🚀 Getting Started](#-getting-started)
- [🔧 Development Setup](#-development-setup)
- [📝 Code Standards](#-code-standards)
- [🧪 Testing Guidelines](#-testing-guidelines)
- [📤 Submitting Changes](#-submitting-changes)
- [🐛 Bug Reports](#-bug-reports)
- [💡 Feature Requests](#-feature-requests)
- [💬 Community Guidelines](#-community-guidelines)

## 🎯 Ways to Contribute

### 🔧 Code Contributions
- **🐛 Bug fixes** - Help us squash bugs
- **✨ New features** - Add exciting new functionality
- **🚀 Performance improvements** - Make MCP RAPI faster
- **📚 Documentation** - Improve our docs
- **🧪 Tests** - Increase test coverage

### 📝 Non-Code Contributions
- **🐛 Bug reports** - Help us identify issues
- **💡 Feature requests** - Suggest new features
- **📖 Documentation improvements** - Fix typos, add examples
- **🎨 Design feedback** - UI/UX suggestions
- **💬 Community support** - Help other users

## 🚀 Getting Started

### 1. 🍴 Fork the Repository
```bash
# Fork the repo on GitHub, then clone your fork
git clone https://github.com/chuangyeshuo/mcprapi.git
cd mcp-rapi

# Add the original repo as upstream
git remote add upstream https://github.com/chuangyeshuo/mcprapi.git
```

### 2. 🌿 Create a Branch
```bash
# Create a new branch for your feature/fix
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/bug-description
```

### 3. 🔄 Stay Updated
```bash
# Regularly sync with upstream
git fetch upstream
git checkout main
git merge upstream/main
```

## 🔧 Development Setup

### 📋 Prerequisites
- **Go** 1.21+
- **Node.js** 18+
- **MySQL** 8.0+
- **Redis** 6.0+
- **Docker** & **Docker Compose**
- **Git**

### 🐳 Quick Setup with Docker
```bash
# Copy environment file
cp .env.example .env

# Start all services
docker-compose -f docker-compose.dev.yml up -d

# Your development environment is ready! 🎉
```

### 💻 Local Development Setup

#### Backend Setup
```bash
cd backend

# Install dependencies
go mod download

# Copy config
cp configs/dev.yaml.example configs/dev.yaml

# Initialize database
go run scripts/init_admin.go

# Start with hot reload
go run cmd/main.go --config configs/dev.yaml
```

#### Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run serve
```

## 📝 Code Standards

### 🔧 Go Code Standards
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Run `go vet` before committing
- Add comments for exported functions
- Use meaningful variable names

**Example:**
```go
// UserService handles user-related business logic
type UserService struct {
    repo repository.UserRepository
}

// CreateUser creates a new user with the given information
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Validate input
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // Business logic here...
}
```

### 🎨 Vue.js Code Standards
- Follow [Vue.js Style Guide](https://vuejs.org/style-guide/)
- Use ESLint configuration provided
- Use meaningful component names
- Add JSDoc comments for complex functions

**Example:**
```vue
<template>
  <div class="user-card">
    <h3>{{ user.name }}</h3>
    <p>{{ user.email }}</p>
  </div>
</template>

<script>
/**
 * UserCard component displays user information
 */
export default {
  name: 'UserCard',
  props: {
    user: {
      type: Object,
      required: true,
      validator: (user) => user.name && user.email
    }
  }
}
</script>
```

### 📝 Commit Message Standards
We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(auth): add JWT token refresh functionality

fix(api): resolve user permission check bug

docs(readme): update installation instructions

test(user): add unit tests for user service
```

## 🧪 Testing Guidelines

### 🔧 Backend Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestUserService ./internal/domain/service
```

**Test Structure:**
```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &mocks.UserRepository{}
    service := NewUserService(mockRepo)
    
    // Act
    user, err := service.CreateUser(context.Background(), &CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
}
```

### 🎨 Frontend Testing
```bash
# Run unit tests
npm run test:unit

# Run e2e tests
npm run test:e2e

# Run tests with coverage
npm run test:coverage
```

## 📤 Submitting Changes

### 1. ✅ Pre-submission Checklist
- [ ] Code follows style guidelines
- [ ] Tests pass locally
- [ ] New tests added for new features
- [ ] Documentation updated
- [ ] Commit messages follow convention
- [ ] No merge conflicts

### 2. 🔍 Code Review Process
1. **Create Pull Request** with clear description
2. **Automated checks** must pass (CI/CD)
3. **Code review** by maintainers
4. **Address feedback** if any
5. **Merge** after approval

### 3. 📝 Pull Request Template
```markdown
## 📝 Description
Brief description of changes

## 🔧 Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## 🧪 Testing
- [ ] Tests pass
- [ ] New tests added

## 📋 Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
```

## 🐛 Bug Reports

### 📝 Bug Report Template
```markdown
**🐛 Bug Description**
Clear description of the bug

**🔄 Steps to Reproduce**
1. Go to '...'
2. Click on '...'
3. See error

**✅ Expected Behavior**
What should happen

**❌ Actual Behavior**
What actually happens

**🖥️ Environment**
- OS: [e.g., macOS 12.0]
- Browser: [e.g., Chrome 95]
- Version: [e.g., 1.0.0]

**📸 Screenshots**
If applicable, add screenshots
```

## 💡 Feature Requests

### 📝 Feature Request Template
```markdown
**💡 Feature Description**
Clear description of the feature

**🎯 Problem Statement**
What problem does this solve?

**💭 Proposed Solution**
How should this work?

**🔄 Alternatives Considered**
Other solutions you've considered

**📈 Additional Context**
Any other context or screenshots
```

## 💬 Community Guidelines

### 🤝 Code of Conduct
- **Be respectful** and inclusive
- **Be constructive** in feedback
- **Be patient** with newcomers
- **Be collaborative** and helpful

### 🗣️ Communication Channels
- **GitHub Issues** - Bug reports and feature requests
- **GitHub Issues** - General questions and ideas
- **Discord** - Real-time chat and support
- **Email** - Private matters

### 🏆 Recognition
Contributors will be recognized in:
- **README.md** contributors section
- **Release notes** for significant contributions
- **Special badges** for long-term contributors

## 🎉 Thank You!

Every contribution, no matter how small, makes MCP RAPI better for everyone. We appreciate your time and effort in helping us build an amazing API permission management system!

---

**Questions?** Feel free to reach out:
- 💬 [GitHub Issues](https://github.com/chuangyeshuo/mcprapi/issues)
- 📧 Email: contributors@mcp-rapi.com
- 📱 Discord: [Join our Discord](https://discord.gg/DmyRA3Nj)

**Happy Contributing! 🚀**