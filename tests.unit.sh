#!/bin/sh
set -e

# Run tests
echo "Running Go tests..."
cd api
go test ./... -v
echo "Tests passed."

# Start the main application
exec "$@"
