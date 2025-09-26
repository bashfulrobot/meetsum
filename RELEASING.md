# Release Process

## Prerequisites

- All changes merged to `main` branch
- All tests passing locally
- Version follows [Semantic Versioning](https://semver.org/)

## Release Commands

### 1. Prepare

```bash
# Ensure clean state
git checkout main
git pull origin main

# Run checks
just check
just test
```

### 2. Create Release

```bash
# Replace with your desired version
just release v1.0.0
```

That's it! The `just release` command will:
- Validate version format
- Check for existing tags
- Build all platform binaries
- Create and push git tag
- Create GitHub release with binaries

### 3. Verify

```bash
# Check the release was created
gh release view v1.0.0
```

## Version Format

Use semantic versioning with `v` prefix:
- ✅ `v1.0.0` (stable)
- ✅ `v1.2.3-alpha.1` (pre-release)
- ❌ `1.0.0` (missing v)
- ❌ `v1.0` (incomplete)