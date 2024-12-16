#!/usr/bin/env bash

# Variables
API_URL="http://localhost:8000"
GET_ENDPOINT="$API_URL/files"

# Example file path and content
FILE_PATH="./test.txt"

echo "Getting file CID..."
curl "$GET_ENDPOINT?filePath=$FILE_PATH"
