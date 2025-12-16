# meetsum Go CLI Project
# Meeting Summary Generator
#
# Note: Production builds use CGO_ENABLED=0 for static binaries.
# Development builds may use CGO for faster compilation.

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
    @echo "  build-all [version]        - Build for all platforms (darwin/linux arm64/amd64)"
    @echo "  generate-changelog [tag]   - Generate changelog since tag or last release"
    @echo "  release <version>          - Create GitHub release with binaries and changelog"
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
build-all version="{{version}}": clean
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🏗️  Building {{app_name}} for all platforms..."
    mkdir -p dist

    # Use provided version or default to git describe
    build_version="{{version}}"
    if [[ "{{version}}" == "dev" || -z "{{version}}" ]]; then
        build_version="{{version}}"
    fi

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

        # Use CGO_ENABLED=0 for all builds to create statically linked binaries
        # This ensures binaries work on all Linux distros including NixOS
        CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
            -ldflags="-X main.Version=${build_version} -X main.BuildTime={{build_time}} -X main.GitCommit={{git_commit}}" \
            -o "$output"
    done

    echo "✅ All builds complete:"
    ls -la dist/

# Generate changelog from conventional commits since last tag
[group('prod')]
generate-changelog since_tag="":
    #!/usr/bin/env bash
    set -euo pipefail

    # Determine the range for changelog generation
    if [[ -n "{{since_tag}}" ]]; then
        range="{{since_tag}}..HEAD"
        echo "📝 Generating changelog since {{since_tag}}..."
    else
        # Get the latest tag, or use initial commit if no tags exist
        latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
        if [[ -n "$latest_tag" ]]; then
            range="${latest_tag}..HEAD"
            echo "📝 Generating changelog since ${latest_tag}..."
        else
            range="$(git rev-list --max-parents=0 HEAD)..HEAD"
            echo "📝 Generating changelog from initial commit..."
        fi
    fi

    # Generate changelog
    echo "## Changelog"
    echo ""

    # Get commits in reverse chronological order
    commits=$(git log --pretty=format:"%h|%s" --no-merges $range)

    # Initialize sections
    declare -a features=()
    declare -a fixes=()
    declare -a docs=()
    declare -a build=()
    declare -a refactor=()
    declare -a other=()

    # Parse commits and categorize
    while IFS='|' read -r hash message; do
        [[ -z "$hash" ]] && continue

        # Extract conventional commit type using simple string matching
        if [[ "$message" == feat:* ]] || [[ "$message" == feat\(*\):* ]]; then
            desc=$(echo "$message" | sed 's/^feat[^:]*:[[:space:]]*//')
            desc=$(echo "$desc" | sed 's/^[^[:alnum:][:space:]]*[[:space:]]*//')
            features+=("- ${desc} (${hash})")
        elif [[ "$message" == fix:* ]] || [[ "$message" == fix\(*\):* ]]; then
            desc=$(echo "$message" | sed 's/^fix[^:]*:[[:space:]]*//')
            desc=$(echo "$desc" | sed 's/^[^[:alnum:][:space:]]*[[:space:]]*//')
            fixes+=("- ${desc} (${hash})")
        elif [[ "$message" == docs:* ]] || [[ "$message" == docs\(*\):* ]]; then
            desc=$(echo "$message" | sed 's/^docs[^:]*:[[:space:]]*//')
            desc=$(echo "$desc" | sed 's/^[^[:alnum:][:space:]]*[[:space:]]*//')
            docs+=("- ${desc} (${hash})")
        elif [[ "$message" == build:* ]] || [[ "$message" == build\(*\):* ]] || [[ "$message" == ci:* ]] || [[ "$message" == ci\(*\):* ]] || [[ "$message" == chore:* ]] || [[ "$message" == chore\(*\):* ]]; then
            desc=$(echo "$message" | sed 's/^[^:]*:[[:space:]]*//')
            desc=$(echo "$desc" | sed 's/^[^[:alnum:][:space:]]*[[:space:]]*//')
            build+=("- ${desc} (${hash})")
        elif [[ "$message" == refactor:* ]] || [[ "$message" == refactor\(*\):* ]] || [[ "$message" == perf:* ]] || [[ "$message" == perf\(*\):* ]] || [[ "$message" == style:* ]] || [[ "$message" == style\(*\):* ]]; then
            desc=$(echo "$message" | sed 's/^[^:]*:[[:space:]]*//')
            desc=$(echo "$desc" | sed 's/^[^[:alnum:][:space:]]*[[:space:]]*//')
            refactor+=("- ${desc} (${hash})")
        else
            # Non-conventional commit or other types
            other+=("- ${message} (${hash})")
        fi
    done <<< "$commits"

    # Output sections
    if [[ ${#features[@]} -gt 0 ]]; then
        echo "### ✨ Features"
        printf '%s\n' "${features[@]}"
        echo ""
    fi

    if [[ ${#fixes[@]} -gt 0 ]]; then
        echo "### 🐛 Bug Fixes"
        printf '%s\n' "${fixes[@]}"
        echo ""
    fi

    if [[ ${#refactor[@]} -gt 0 ]]; then
        echo "### ♻️ Code Changes"
        printf '%s\n' "${refactor[@]}"
        echo ""
    fi

    if [[ ${#build[@]} -gt 0 ]]; then
        echo "### 🔧 Build & CI"
        printf '%s\n' "${build[@]}"
        echo ""
    fi

    if [[ ${#docs[@]} -gt 0 ]]; then
        echo "### 📝 Documentation"
        printf '%s\n' "${docs[@]}"
        echo ""
    fi

    if [[ ${#other[@]} -gt 0 ]]; then
        echo "### 🔄 Other Changes"
        printf '%s\n' "${other[@]}"
        echo ""
    fi

# Create a GitHub release with binaries
[group('prod')]
release version: (build-all version)
    #!/usr/bin/env bash
    set -euo pipefail

    if [[ -z "{{version}}" ]]; then
        echo "❌ Version argument required. Usage: just release v1.0.0"
        exit 1
    fi

    echo "🚀 Creating release {{version}}..."

    # Validate version format
    if [[ ! "{{version}}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "❌ Version must be in format v1.0.0"
        exit 1
    fi

    # Check if tag already exists
    if git tag -l | grep -q "^{{version}}$"; then
        echo "❌ Tag {{version}} already exists"
        exit 1
    fi

    # Verify gh CLI is available
    if ! command -v gh >/dev/null 2>&1; then
        echo "❌ gh CLI not found. Install with: brew install gh"
        exit 1
    fi

    # Verify we're authenticated with gh
    if ! gh auth status >/dev/null 2>&1; then
        echo "❌ Not authenticated with GitHub. Run: gh auth login"
        exit 1
    fi

    # Create and push git tag
    echo "📝 Creating git tag {{version}}..."
    git tag -a "{{version}}" -m "Release {{version}}"
    git push origin "{{version}}"

    # Build all binaries with the release version
    echo "🔨 Building release binaries..."

    # Generate changelog since last release
    echo "📝 Generating changelog..."
    changelog=$(just generate-changelog)

    # Create GitHub release
    echo "📦 Creating GitHub release..."

    # Create GitHub release with generated changelog
    gh release create "{{version}}" \
        --title "{{app_name}} {{version}}" \
        --notes "$changelog" \
        dist/{{app_name}}-darwin-arm64 \
        dist/{{app_name}}-darwin-amd64 \
        dist/{{app_name}}-linux-amd64 \
        dist/{{app_name}}-linux-arm64

    echo "✅ Release {{version}} created successfully!"
    echo "🔗 View at: $(gh release view {{version}} --json url --jq .url)"

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