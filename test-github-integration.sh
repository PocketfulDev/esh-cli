#!/bin/bash
# Test GitHub Actions integration locally

set -e

echo "🧪 Testing GitHub Actions Integration Locally"
echo "=============================================="

# Check if bc (basic calculator) is available for coverage thresholds
if ! command -v bc &> /dev/null; then
    echo "❌ bc (basic calculator) is required for coverage threshold checking"
    echo "   Install with: brew install bc (macOS) or apt install bc (Ubuntu)"
    exit 1
fi

echo ""
echo "1️⃣  Running tests with JSON output (CI-style)..."
echo "------------------------------------------------"
make test-coverage-json

echo ""
echo "2️⃣  Checking coverage thresholds..."
echo "-----------------------------------"
make test-coverage-check

echo ""
echo "3️⃣  Validating test result files..."
echo "-----------------------------------"
if [ -f "test-results.json" ]; then
    echo "✅ test-results.json created"
    echo "   📊 Test count: $(grep -c '"Test":' test-results.json || echo "0")"
else
    echo "❌ test-results.json not found"
fi

if [ -f "coverage.out" ]; then
    echo "✅ coverage.out created"
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo "   📈 Total coverage: $COVERAGE"
else
    echo "❌ coverage.out not found"
fi

if [ -f "coverage.html" ]; then
    echo "✅ coverage.html created"
    echo "   🌐 Open with: open coverage.html"
else
    echo "❌ coverage.html not found"
fi

if [ -f "coverage-func.txt" ]; then
    echo "✅ coverage-func.txt created"
    echo "   📋 Function breakdown available"
else
    echo "❌ coverage-func.txt not found"
fi

echo ""
echo "4️⃣  Testing cross-platform builds..."
echo "------------------------------------"
platforms=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")

for platform in "${platforms[@]}"; do
    IFS='/' read -r os arch <<< "$platform"
    binary_name="esh-cli-${os}-${arch}-test"
    
    if GOOS=$os GOARCH=$arch go build -o $binary_name .; then
        echo "✅ Built $binary_name"
        rm -f $binary_name
    else
        echo "❌ Failed to build $binary_name"
        exit 1
    fi
done

echo ""
echo "5️⃣  Simulating GitHub Actions summary..."
echo "----------------------------------------"
if [ -f "coverage-func.txt" ]; then
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    UTILS_COVERAGE=$(go test -coverprofile=utils_coverage.out -covermode=atomic ./pkg/utils >/dev/null 2>&1 && go tool cover -func=utils_coverage.out | grep total | awk '{print $3}' | sed 's/%//' || echo "0")
    
    echo "## 📊 Test Coverage Report"
    echo ""
    echo "| Metric | Value | Status |"
    echo "|--------|--------|--------|"
    echo "| Overall Coverage | ${COVERAGE}% | $([ $(echo "$COVERAGE >= 30" | bc -l) -eq 1 ] && echo "✅ Pass" || echo "⚠️ Below Target") |"
    echo "| Utils Package | ${UTILS_COVERAGE}% | $([ $(echo "$UTILS_COVERAGE >= 60" | bc -l) -eq 1 ] && echo "✅ Pass" || echo "⚠️ Below Target") |"
    echo ""
    echo "### 📋 Coverage Details"
    echo "\`\`\`"
    cat coverage-func.txt
    echo "\`\`\`"
fi

echo ""
echo "🎉 GitHub Actions integration test completed!"
echo ""
echo "📁 Generated files:"
echo "   - test-results.json (for dorny/test-reporter)"
echo "   - coverage.out (Go coverage profile)"
echo "   - coverage.html (Visual coverage report)"
echo "   - coverage-func.txt (Function breakdown)"
echo ""
echo "💡 Next steps:"
echo "   1. View coverage: open coverage.html"
echo "   2. Check test results: cat test-results.json | jq"
echo "   3. Run CI tests anytime: make test-ci"
