#!/usr/bin/env bash

# Variables
API_URL="http://localhost:8000"
UPLOAD_ENDPOINT="$API_URL/files"

# Example file path and content
FILE_PATH="./test.txt"
FILE_CONTENT="Hello World!"
BASE64_CONTENT=$(echo -n "$FILE_CONTENT" | base64)

echo "Uploading file..."
curl -X POST "$UPLOAD_ENDPOINT" \
     -H "Content-Type: application/json" \
     -d '{"filePath":"'"$FILE_PATH"'","file":"'"$BASE64_CONTENT"'"}'
