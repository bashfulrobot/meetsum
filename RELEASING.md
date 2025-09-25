# Release Process

This document describes the release process for meetsum.

## Prerequisites

1. All changes are merged to `main` branch
2. All tests are passing in CI
3. Version number follows [Semantic Versioning](https://semver.org/)

## Release Types

### Stable Release (v1.2.3)
For production-ready releases with new features, bug fixes, or improvements.

### Pre-Release (v1.2.3-alpha.1, v1.2.3-beta.1, v1.2.3-rc.1)
For testing releases before stable release.

## Release Steps

### 1. Prepare the Release

```bash
# Ensure you're on the main branch with latest changes
git checkout main
git pull origin main

# Run final checks
just check
just test
just build-all

# Verify everything is working
./meetsum --help
./meetsum version
```

### 2. Create and Push Release Tag

```bash
# Create a new semver tag (replace with actual version)
export VERSION="v1.0.0"

# Create and push the tag
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION
```

**Important**: The tag MUST follow semver format:
- ✅ `v1.0.0` (stable release)
- ✅ `v1.2.3-alpha.1` (pre-release)
- ✅ `v2.0.0-beta.2` (pre-release)
- ❌ `1.0.0` (missing 'v' prefix)
- ❌ `v1.0` (incomplete semver)

### 3. Automatic Release Process

Once the tag is pushed, GitHub Actions will automatically:

1. **Validate** the semver tag format
2. **Run tests** to ensure quality
3. **Build binaries** for all platforms:
   - `meetsum-darwin-amd64` (macOS Intel)
   - `meetsum-darwin-arm64` (macOS Apple Silicon)
   - `meetsum-linux-amd64` (Linux Intel)
   - `meetsum-linux-arm64` (Linux ARM)
4. **Generate checksums** for all binaries
5. **Create changelog** from git commit messages
6. **Create GitHub release** with binaries attached
7. **Upload artifacts** for download

### 4. Verify Release

After the GitHub Action completes:

1. Check the [Releases page](https://github.com/bashfulrobot/meetsum/releases)
2. Verify all binaries are attached
3. Test download and execution of at least one binary
4. Verify the changelog is accurate

## Versioning Guidelines

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v2.0.0): Incompatible API changes
- **MINOR** version (v1.1.0): New functionality (backward compatible)
- **PATCH** version (v1.0.1): Bug fixes (backward compatible)

### Examples:

- **v1.0.0**: Initial release
- **v1.0.1**: Bug fix release
- **v1.1.0**: New feature added (like new command)
- **v1.1.1**: Bug fix in the new feature
- **v2.0.0**: Breaking change (like config format change)

## Pre-Release Process

For testing major changes:

```bash
# Create pre-release tag
git tag -a v1.1.0-alpha.1 -m "Pre-release v1.1.0-alpha.1"
git push origin v1.1.0-alpha.1
```

Pre-releases are automatically marked as "Pre-release" in GitHub.

## Hotfix Process

For critical bug fixes:

```bash
# Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b hotfix/critical-fix

# Make the fix
# ... make changes ...
git add .
git commit -m "fix: critical security issue in config loading"

# Push and create PR
git push origin hotfix/critical-fix
# Create PR, get approval, merge

# After merge, create release tag
git checkout main
git pull origin main
git tag -a v1.0.1 -m "Hotfix release v1.0.1"
git push origin v1.0.1
```

## Rollback Process

If a release has critical issues:

### Option 1: Quick Fix
```bash
# Fix the issue and create a patch release
git tag -a v1.0.1 -m "Fix for v1.0.0 issues"
git push origin v1.0.1
```

### Option 2: Mark Release as Draft
1. Go to GitHub Releases
2. Edit the problematic release
3. Mark as "Draft" to hide it
4. Fix issues and create new release

## Changelog Generation

The changelog is automatically generated from commit messages. To ensure good changelogs:

### Good Commit Messages:
- `feat: add new --output flag for custom summary location`
- `fix: resolve config file loading on Windows`
- `docs: update installation instructions`
- `refactor: improve error handling in summary generation`

### Commit Message Format:
```
<type>: <description>

[optional body]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests
- `chore`: Changes to build process or auxiliary tools

## Troubleshooting

### Tag Already Exists
```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin --delete v1.0.0

# Create new tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### GitHub Action Failed
1. Check the [Actions tab](https://github.com/bashfulrobot/meetsum/actions)
2. Review the failure logs
3. Fix the issue and push a new tag

### Binary Not Working
1. Test locally with `just build-all`
2. Check the ldflags in the GitHub Action match the justfile
3. Verify Go version compatibility

## Release Checklist

Before creating a release tag:

- [ ] All tests passing in CI
- [ ] Version number follows semver
- [ ] CHANGELOG.md updated (if maintaining manually)
- [ ] No known critical bugs
- [ ] Documentation is up to date
- [ ] Built and tested locally with `just build-all`

After release:
- [ ] Verify GitHub release created successfully
- [ ] Test at least one binary download
- [ ] Update any external documentation mentioning version numbers
- [ ] Announce release (if applicable)