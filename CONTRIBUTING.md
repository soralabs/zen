# Contributing to Zen
Thank you for considering contributing to **zen**! This document provides guidelines to ensure your contributions are smooth and effective.

## Table of Contents
1. [How to Contribute](#how-to-contribute)
    - [Reporting Issues](#reporting-issues)
    - [Submitting Pull Requests](#submitting-pull-requests)
    - [Suggesting Enhancements](#suggesting-enhancements)
2. [Development Setup](#development-setup)
3. [Code Guidelines](#code-guidelines)
4. [Testing](#testing)
5. [Style Guide](#style-guide)

---

## How to Contribute

### Reporting Issues
If you encounter a bug, have a question, or want to request a feature:
1. Check the [issue tracker](https://github.com/soralabs/zen/issues) to see if it has already been reported.
2. If it's a new issue, create one. Include:
   - A clear, descriptive title.
   - Steps to reproduce the issue.
   - Expected and actual results.
   - Relevant logs or screenshots, if applicable.

### Submitting Pull Requests
1. Fork the repository and clone it locally.
2. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Make your changes. Ensure your code adheres to the [Style Guide](#style-guide) and includes tests.
4. Commit your changes using conventional commit messages:
   ```bash
   git commit -m "type(scope): description"
   ```
   
   Conventional commit format:
   - Format: `type(scope): description`
   - Types:
     - `feat`: New feature
     - `fix`: Bug fix
     - `docs`: Documentation changes
     - `style`: Code style changes (formatting, missing semicolons, etc.)
     - `refactor`: Code refactoring
     - `perf`: Performance improvements
     - `test`: Adding or modifying tests
     - `chore`: Maintenance tasks, dependencies, etc.
   - Scope: Optional component/module name (e.g., `api`, `cli`, `core`)
   - Description: Present tense, lowercase, no period at end

   Examples:
   ```bash
   git commit -m "feat(api): add user authentication endpoint"
   git commit -m "fix(core): resolve null pointer in config parser"
   git commit -m "docs: update installation instructions"
   git commit -m "test(cli): add integration tests for command parsing"
   ```

5. Push your changes:
   ```bash
   git push origin feature/your-feature-name
   ```
6. Open a pull request on GitHub. Include:
   - A link to the related issue, if applicable.
   - A summary of the changes.
   - Any additional context or details.

### Suggesting Enhancements
Enhancement suggestions can be submitted as issues. Include:
- A clear title.
- The problem the enhancement addresses.
- Your proposed solution or approach.

---

## Development Setup

### Prerequisites
Ensure you have the following installed:
- [Go](https://golang.org/doc/install) (version 1.23.3 or later)
- [Git](https://git-scm.com/)
- Any additional dependencies listed in the `README.md`.

### Setup Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/soralabs/zen.git
   ```
2. Navigate to the project directory:
   ```bash
   cd zen
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Build the project:
   ```bash
   go build
   ```
5. Run the application:
   ```bash
   go run main.go
   ```

---

## Code Guidelines
- Follow Go's idiomatic patterns. Refer to the [Effective Go](https://go.dev/doc/effective_go) guide.
- Write clear, concise, and well-documented code.
- Keep functions small and focused.
- Use meaningful variable and function names.

### Folder Structure
- `/cmd`: Main applications for the project.
- `/pkg`: Library code that can be used by external applications.
- `/internal`: Code not intended for external use.
- `/test`: Additional testing utilities.

---

## Testing
- Write tests for all new features and bug fixes.
- Run tests before submitting a pull request:
  ```bash
  go test ./...
  ```
- Use [Go's testing package](https://pkg.go.dev/testing) and ensure your tests cover edge cases.

---

## Style Guide
- Use `gofmt` to format your code:
  ```bash
  gofmt -s -w .
  ```
- Use `golint` to check for stylistic issues:
  ```bash
  golint ./...
  ```
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

---

We're excited to see your contributions! If you have questions, feel free to reach out by opening an issue or joining our community discussions.