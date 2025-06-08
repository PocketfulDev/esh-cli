#!/bin/bash

# Script to get clean coverage values for CI/CD
# Usage: ./scripts/get-coverage.sh [utils|total]

set -e

case "${1:-total}" in
    "utils")
        # Get utils package coverage
        go test -coverprofile=utils_coverage.out -covermode=atomic ./pkg/utils >/dev/null 2>&1
        go tool cover -func=utils_coverage.out | grep total | awk '{print $3}' | sed 's/%//'
        ;;
    "total")
        # Get total coverage (assumes coverage.out already exists)
        if [ -f coverage.out ]; then
            go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//'
        else
            echo "0"
        fi
        ;;
    *)
        echo "Usage: $0 [utils|total]"
        exit 1
        ;;
esac
