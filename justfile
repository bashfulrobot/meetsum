# meetsum Go CLI Project
# Meeting Summary Generator
#
# Note: This project requires CGO due to some dependencies.
# Ensure you have a C compiler (gcc, clang) installed.
# For clang users, set CC=clang before running build commands.

# === Settings ===
set dotenv-load := true
set ignore-comments := true
set fallback := true
set shell := ["bash", "-euo", "pipefail", "-c"]

# === Variables ===
app_name := "meetsum"
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
build_time := `date +%Y-%m-%d_%H:%M:%S`
git_commit := `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`

# === Help ===
# Show available commands
default:
    @echo "🤖 meetsum - Meeting Summary Generator"
    @echo "====================================="
    @just --list --unsorted
    @echo ""
    @echo "🔧 Commands with Parameters:"
    @echo "  run [args...]              - Run with arguments (e.g., just run --trace /path/to/meeting)"
    @echo "  build-all                  - Build for all platforms (darwin/linux arm64/amd64)"
    @echo "  lint [fix=false]           - Add fix=true to auto-fix issues"
    @echo ""
    @echo "💡 Pro Tips:"
    @echo "  • Run 'just check' before committing"
    @echo "  • Use 'just run --help' to see CLI options"
    @echo "  • Common workflow: just check → just test → just build"

# === Development Commands ===
# Quick validation without building
[group('dev')]
check:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🔍 Running development checks..."

    echo "📋 Formatting..."
    go fmt ./...

    echo "📦 Tidying modules..."
    go mod tidy

    echo "🔍 Vetting code..."
    go vet ./...

    echo "✅ All checks passed"

# Build for current platform
[group('dev')]
build:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🔨 Building {{app_name}} for current platform..."
    CC=clang go build -ldflags="-X main.Version={{version}} -X main.BuildTime={{build_time}} -X main.GitCommit={{git_commit}}" -o {{app_name}}
    echo "✅ Build complete: ./{{app_name}}"
    echo "💡 Test version: ./{{app_name}} version"

# Run the application with arguments
[group('dev')]
run *args:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🚀 Running {{app_name}}..."
    CC=clang go run . {{args}}

# Run tests
[group('dev')]
test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🧪 Running tests..."
    go test ./...

# Run tests with coverage
[group('dev')]
test-coverage:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🧪 Running tests with coverage..."
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    echo "📊 Coverage report: ./coverage.html"

# === Production Commands ===
# Build for all platforms
[group('prod')]
build-all: clean
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🏗️  Building {{app_name}} for all platforms..."
    mkdir -p dist

    platforms=(
        "darwin/arm64"
        "darwin/amd64"
        "linux/amd64"
        "linux/arm64"
    )

    for platform in "${platforms[@]}"; do
        GOOS=${platform%/*}
        GOARCH=${platform#*/}
        output="dist/{{app_name}}-${GOOS}-${GOARCH}"

        echo "Building for ${GOOS}/${GOARCH}..."
        CC=clang GOOS=$GOOS GOARCH=$GOARCH go build \
            -ldflags="-X main.Version={{version}} -X main.BuildTime={{build_time}} -X main.GitCommit={{git_commit}}" \
            -o "$output"
    done

    echo "✅ All builds complete:"
    ls -la dist/

# Install locally
[group('prod')]
install:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "📦 Installing {{app_name}}..."
    CC=clang go install -ldflags="-X main.Version={{version}} -X main.BuildTime={{build_time}} -X main.GitCommit={{git_commit}}"
    echo "✅ Installed to $(go env GOPATH)/bin/{{app_name}}"

# === Maintenance Commands ===
# Clean build artifacts
[group('maintenance')]
clean:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🧹 Cleaning build artifacts..."
    rm -f {{app_name}}
    rm -rf dist/
    rm -f coverage.out coverage.html
    go clean

# Format code
[group('maintenance')]
fmt:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "📋 Formatting code..."
    go fmt ./...
    echo "✅ Code formatted"

# Lint code with optional auto-fix
[group('maintenance')]
lint fix="false":
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v golangci-lint >/dev/null 2>&1; then
        echo "❌ golangci-lint not found. Install with:"
        echo "   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        exit 1
    fi

    if [[ "{{fix}}" == "true" ]]; then
        echo "🔧 Linting code with auto-fix..."
        golangci-lint run --fix
    else
        echo "🔍 Linting code..."
        golangci-lint run
    fi

# Tidy dependencies
[group('maintenance')]
tidy:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "📦 Tidying dependencies..."
    go mod tidy
    echo "✅ Dependencies tidied"

# Update dependencies
[group('maintenance')]
update:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "⬆️  Updating dependencies..."
    go get -u ./...
    go mod tidy
    echo "✅ Dependencies updated"

# Check for security vulnerabilities
[group('maintenance')]
security:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v govulncheck >/dev/null 2>&1; then
        echo "📦 Installing govulncheck..."
        go install golang.org/x/vuln/cmd/govulncheck@latest
    fi
    echo "🔒 Checking for security vulnerabilities..."
    govulncheck ./...

# === Info Commands ===
# Show build and module info
[group('info')]
info:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "📊 {{app_name}} Build Information"
    echo "================================"
    echo "Go version: $(go version)"
    echo "Module: $(go list -m)"
    echo "Git version: {{version}}"
    echo "Git commit: {{git_commit}}"
    echo "Build time: {{build_time}}"
    echo ""
    echo "📦 Dependencies:"
    go list -m all | head -20
    dep_count=$(go list -m all | wc -l)
    if [ $dep_count -gt 20 ]; then
        echo "... and $((dep_count - 20)) more"
    fi

# Development setup
[group('info')]
setup:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🛠️  Setting up development environment..."
    go mod tidy
    go mod download

    # Install helpful tools
    echo "📦 Installing development tools..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install golang.org/x/vuln/cmd/govulncheck@latest

    echo "✅ Development environment ready"

# === Workflow Aliases ===
alias c := check
alias b := build
alias r := run
alias t := test
alias tc := test-coverage