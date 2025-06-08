# ESH CLI Testing Scenarios

## 🧪 Comprehensive Testing Strategy

### Test Categories

1. **Unit Tests** - Individual function testing
2. **Integration Tests** - Command-level testing  
3. **End-to-End Tests** - Full workflow testing
4. **Performance Tests** - Large repository testing
5. **Compatibility Tests** - Backward compatibility
6. **Edge Case Tests** - Error handling and edge cases

---

## 📋 Unit Test Scenarios

### Semantic Version Parsing Tests
```go
// pkg/utils/semver_test.go additions needed
func TestParseSemanticVersionEdgeCases(t *testing.T) {
    tests := []struct {
        name      string
        version   string
        shouldErr bool
        expected  *SemanticVersion
    }{
        {"Version with v prefix", "v1.2.3", false, &SemanticVersion{1, 2, 3, ""}},
        {"Version with prerelease", "1.2.3-alpha.1", false, &SemanticVersion{1, 2, 3, "alpha.1"}},
        {"Version with build metadata", "1.2.3+build.1", false, &SemanticVersion{1, 2, 3, ""}},
        {"Leading zeros", "01.02.03", true, nil},
        {"Negative numbers", "-1.2.3", true, nil},
        {"Empty prerelease", "1.2.3-", false, &SemanticVersion{1, 2, 3, ""}},
    }
    // Test implementation...
}
```

### Version Bump Logic Tests
```go
func TestBumpSemanticVersionEdgeCases(t *testing.T) {
    tests := []struct {
        name        string
        version     string
        bumpType    BumpType
        expected    string
        shouldErr   bool
    }{
        {"Major bump resets minor and patch", "1.2.3", BumpMajor, "2.0.0", false},
        {"Minor bump resets patch", "1.2.3", BumpMinor, "1.3.0", false},
        {"Patch bump increments patch", "1.2.3", BumpPatch, "1.2.4", false},
        {"Bump with prerelease", "1.2.3-alpha", BumpPatch, "1.2.4", false},
        {"Invalid bump type", "1.2.3", "invalid", "", true},
    }
    // Test implementation...
}
```

### Tag Format Validation Tests
```go
func TestBumpTagVersionWithServices(t *testing.T) {
    tests := []struct {
        name        string
        tag         string
        bumpType    BumpType
        environment string
        service     string
        expected    string
        shouldErr   bool
    }{
        {"Service tag major bump", "api_stg6_1.2.3-1", BumpMajor, "stg6", "api", "api_stg6_2.0.0-1", false},
        {"Non-service tag minor bump", "stg6_1.2.3-1", BumpMinor, "stg6", "", "stg6_1.3.0-1", false},
        {"Cross-environment error", "dev_1.2.3-1", BumpPatch, "stg6", "", "", true},
    }
    // Test implementation...
}
```

---

## 🔄 Integration Test Scenarios

### Command Integration Tests

Create test file: `cmd/semantic_versioning_integration_test.go`

```go
func TestBumpVersionCommand(t *testing.T) {
    // Setup test repository
    tempDir := setupTestRepo(t)
    defer cleanup(tempDir)
    
    tests := []struct {
        name           string
        args           []string
        expectError    bool
        expectOutput   string
        preConditions  func(t *testing.T)
        postConditions func(t *testing.T)
    }{
        {
            name: "Major version bump",
            args: []string{"bump-version", "dev", "--major", "--preview"},
            preConditions: func(t *testing.T) {
                // Create initial tag
                createTag(t, "dev_1.2.3-1")
            },
            expectOutput: "dev_2.0.0-1",
        },
        {
            name: "Auto-detect from conventional commits",
            args: []string{"bump-version", "dev", "--auto"},
            preConditions: func(t *testing.T) {
                createTag(t, "dev_1.2.3-1")
                createConventionalCommit(t, "feat: add new feature")
            },
            expectOutput: "dev_1.3.0-1", // Should detect minor bump
        },
    }
    // Test implementation...
}
```

### Workflow Integration Tests

```go
func TestFullVersioningWorkflow(t *testing.T) {
    scenarios := []struct {
        name     string
        workflow func(t *testing.T)
    }{
        {
            name: "Feature branch to production workflow",
            workflow: func(t *testing.T) {
                // 1. Create feature branch
                checkoutBranch(t, "feature/new-api")
                
                // 2. Make conventional commits
                makeCommit(t, "feat: add new API endpoint")
                makeCommit(t, "test: add API tests")
                
                // 3. Get branch suggestions
                output := runCommand(t, "branch-version", "--suggest")
                assert.Contains(t, output, "MINOR")
                
                // 4. Auto-tag for dev
                runCommand(t, "branch-version", "--auto-tag", "--env", "dev")
                
                // 5. Promote through environments
                runCommand(t, "add-tag", "stg6", "1.3.0-1", "--from", "dev_1.3.0-1")
                runCommand(t, "add-tag", "production2", "1.3.0-1", "--from", "stg6_1.3.0-1")
                
                // 6. Generate changelog
                output = runCommand(t, "changelog", "--conventional-commits")
                assert.Contains(t, output, "new API endpoint")
            },
        },
    }
}
```

---

## 🏁 End-to-End Test Scenarios

### Real Repository Testing

Create test script: `test-e2e-semantic-versioning.sh`

```bash
#!/bin/bash
# End-to-end semantic versioning tests

set -e

# Test repository setup
TEST_REPO="/tmp/esh-cli-e2e-test"
rm -rf $TEST_REPO
mkdir -p $TEST_REPO
cd $TEST_REPO

git init
git config user.name "Test User"
git config user.email "test@example.com"

# Create initial commit and tag
echo "Initial version" > README.md
git add README.md
git commit -m "initial commit"

# Test 1: Bump version commands
echo "=== Testing bump-version commands ==="

# Create initial tag
esh-cli add-tag dev 1.0.0-1

# Test major bump
echo "Feature 1" >> README.md
git add README.md
git commit -m "feat!: breaking change"
esh-cli bump-version dev --auto --preview | grep "2.0.0-1"

# Test minor bump  
echo "Feature 2" >> README.md
git add README.md
git commit -m "feat: new feature"
esh-cli bump-version dev --auto --preview | grep "1.1.0-1"

# Test patch bump
echo "Fix 1" >> README.md
git add README.md  
git commit -m "fix: bug fix"
esh-cli bump-version dev --auto --preview | grep "1.0.1-1"

# Test 2: Version listing and filtering
echo "=== Testing version-list commands ==="
esh-cli version-list dev --format json | jq '.[]'
esh-cli version-list dev --major 1 | grep "1\."

# Test 3: Version comparison
echo "=== Testing version-diff commands ==="
esh-cli version-diff dev --history
esh-cli version-diff dev_1.0.0-1 --commits

# Test 4: Changelog generation
echo "=== Testing changelog commands ==="
esh-cli changelog dev --conventional-commits
esh-cli changelog --full --format markdown

# Test 5: Branch-based workflows
echo "=== Testing branch-version commands ==="
git checkout -b feature/test-feature
echo "New feature" >> README.md
git add README.md
git commit -m "feat: implement test feature"

esh-cli branch-version --suggest | grep "MINOR"
esh-cli branch-version --auto-tag --env dev

# Cleanup
cd /
rm -rf $TEST_REPO
echo "✅ All end-to-end tests passed!"
```

---

## ⚡ Performance Test Scenarios

### Large Repository Testing

```bash
#!/bin/bash
# Performance tests for large repositories

# Test with repository containing:
# - 10,000+ commits
# - 1,000+ tags  
# - 100+ branches
# - Multiple environments and services

echo "=== Performance Testing ==="

time esh-cli version-list --all > /dev/null
time esh-cli version-diff dev --history > /dev/null
time esh-cli changelog --full > /dev/null

# Memory usage testing
echo "=== Memory Usage Testing ==="
/usr/bin/time -l esh-cli version-list --all --format json > large-repo-test.json

# Concurrent operations testing
echo "=== Concurrent Operations Testing ==="
for env in dev stg6 production2; do
    esh-cli version-list $env --format json > ${env}-versions.json &
done
wait
```

### Load Testing Script

```go
// cmd/load_test.go
func BenchmarkBumpVersion(b *testing.B) {
    setupBenchmarkRepo(b)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Benchmark version bumping operations
        err := runBumpVersionCommand("dev", "--patch", "--preview")
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkVersionList(b *testing.B) {
    setupLargeRepo(b) // Create repo with 1000+ tags
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        err := runVersionListCommand("--all", "--format", "json")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

## 🔍 Edge Case Testing

### Error Handling Tests

```go
func TestErrorHandling(t *testing.T) {
    tests := []struct {
        name        string
        command     []string
        expectError string
    }{
        {
            name:        "Invalid environment",
            command:     []string{"bump-version", "invalid-env", "--major"},
            expectError: "invalid environment",
        },
        {
            name:        "No existing tags",
            command:     []string{"bump-version", "dev", "--major"},
            expectError: "no existing tags found",
        },
        {
            name:        "Conflicting flags",
            command:     []string{"bump-version", "dev", "--major", "--minor"},
            expectError: "mutually exclusive",
        },
        {
            name:        "Invalid tag format in diff",
            command:     []string{"version-diff", "invalid-tag"},
            expectError: "invalid tag format",
        },
    }
}
```

### Boundary Condition Tests

```go
func TestBoundaryConditions(t *testing.T) {
    tests := []struct {
        name        string
        scenario    func(t *testing.T)
    }{
        {
            name: "Maximum version numbers",
            scenario: func(t *testing.T) {
                // Test with version 999.999.999
                tag := "dev_999.999.999-1"
                newTag, err := utils.BumpTagVersion(tag, utils.BumpPatch, "dev", "")
                assert.Equal(t, "dev_999.999.1000-1", newTag)
            },
        },
        {
            name: "Very long service names",
            scenario: func(t *testing.T) {
                longService := strings.Repeat("a", 100)
                // Test behavior with long service names
            },
        },
        {
            name: "Empty repository",
            scenario: func(t *testing.T) {
                // Test commands on empty git repository
            },
        },
    }
}
```

---

## 🔄 Compatibility Testing

### Backward Compatibility Tests

```bash
#!/bin/bash
# Test backward compatibility with existing workflows

echo "=== Backward Compatibility Testing ==="

# Test that existing add-tag commands still work
esh-cli add-tag dev 1.0-0  # Old format
esh-cli add-tag dev last   # Should still work

# Test promotion workflows
esh-cli add-tag stg6 1.0-0 --from dev_1.0-0

# Test service-specific tags
esh-cli add-tag dev 1.0-0 --service myservice

# Test hot-fix workflows
git checkout -b release/1.0
esh-cli add-tag dev 1.0-0 --hot-fix

echo "✅ Backward compatibility maintained"
```

### Migration Testing

```bash
#!/bin/bash
# Test migration scenarios from old to new versioning

echo "=== Migration Testing ==="

# Scenario: Repository with old-style tags
esh-cli add-tag dev 1.0-0      # Old style
esh-cli add-tag dev 1.1-0      # Old style

# Migrate to semantic versioning
esh-cli bump-version dev --patch  # Should work with existing tags

# Verify mixed tag formats work
esh-cli version-list dev          # Should show both old and new formats
esh-cli version-diff dev --history # Should handle mixed formats

echo "✅ Migration scenarios tested"
```

---

## 📊 Test Reporting

### Test Coverage Analysis

```bash
#!/bin/bash
# Generate comprehensive test coverage report

echo "=== Generating Test Coverage Report ==="

# Run all tests with coverage
go test -v -coverprofile=coverage.out -covermode=atomic ./...

# Generate detailed coverage reports
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out > coverage-func.txt

# Extract coverage metrics
TOTAL_COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo "Total coverage: $TOTAL_COVERAGE"

# Coverage by package
echo "=== Coverage by Package ==="
go tool cover -func=coverage.out | grep -E "^[^/]*\.go:" | sort

# Identify uncovered lines
echo "=== Uncovered Code Analysis ==="
go tool cover -func=coverage.out | grep -E "\s+0\.0%" || echo "All code covered!"
```

### Continuous Testing Setup

```yaml
# .github/workflows/comprehensive-testing.yml
name: Comprehensive Testing

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.21'
    
    - name: Run unit tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  integration-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.21'
    
    - name: Build CLI
      run: make build
    
    - name: Run integration tests
      run: ./test-e2e-semantic-versioning.sh

  performance-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.21'
    
    - name: Run performance tests
      run: go test -bench=. -benchmem ./...
```

---

## 🎯 Test Execution Strategy

### Development Testing
1. **Unit tests** after each function change
2. **Integration tests** after command modifications  
3. **Manual testing** for user experience validation

### Pre-Release Testing
1. **Full test suite** execution
2. **Performance benchmarks** 
3. **Compatibility testing** with real repositories
4. **End-to-end workflows** validation

### Production Monitoring
1. **Error tracking** for command failures
2. **Performance monitoring** for large repositories
3. **User feedback** collection for improvement areas

This comprehensive testing strategy ensures the semantic versioning features are robust, performant, and maintain backward compatibility while providing extensive coverage for all use cases.
