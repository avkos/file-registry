#!/bin/sh
set -e

# Run tests
echo "Running Go tests..."
cd api
go test tests/e2e/e2e_test.go
echo "Tests passed."

# Start the main application
exec "$@"
