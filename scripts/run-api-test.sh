#!/usr/bin/env bash
set -euo pipefail

BASE="http://localhost:3000"
USERNAME="testuser"
PASS="password123"

echo "1) Register"
resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASS\"}")
body=$(echo "$resp" | sed '$d')
code=$(echo "$resp" | tail -n1)
echo "HTTP $code"
echo "$body"
if [ "$code" = "201" ]; then
  echo "Registered OK"
fi

echo
echo "2) Login"
resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASS\"}")
body=$(echo "$resp" | sed '$d')
code=$(echo "$resp" | tail -n1)
echo "HTTP $code"
echo "$body"
if [ "$code" != "200" ]; then
  echo "Login failed"; exit 1
fi

ACCESS_TOKEN=$(echo "$body" | jq -r '.access_token // .token')
REFRESH_TOKEN=$(echo "$body" | jq -r '.refresh_token // empty')
if [ -z "$ACCESS_TOKEN" ]; then
  echo "No access token found in login response"; exit 1
fi
echo "Got access token"

echo
echo "3) Create Post"
resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/posts" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello from automated test"}')
body=$(echo "$resp" | sed '$d')
code=$(echo "$resp" | tail -n1)
echo "HTTP $code"
echo "$body"

echo
echo "4) Get Feed (page=1)"
resp=$(curl -s -w "\n%{http_code}" -X GET "$BASE/api/feed?page=1&limit=5" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
body=$(echo "$resp" | sed '$d')
code=$(echo "$resp" | tail -n1)
echo "HTTP $code"
echo "$body"

echo
echo "Done. Note: This script is minimal — expand it to include negative tests and edge cases."
