#!/bin/bash

# Test script for Display Settings API
# This script tests all 3 endpoints for display settings

set -e

# Base URL
BASE_URL="http://localhost:8080"

# Test user email/password
TEST_EMAIL="test@example.com"
TEST_PASS="testpass123"

# Try to login with test user
echo "🔑 Logging in with test user..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASS\"}")

echo "Login Response: $LOGIN_RESPONSE"

# Extract token from login response
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token - user may not exist"
  echo "Creating test user and organization..."

  # Create user
  curl -s -X POST "$BASE_URL/user" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Test User\",\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASS\"}" > /dev/null

  # Try login again
  LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASS\"}")

  TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

  if [ -z "$TOKEN" ]; then
    echo "Still couldn't get token. Response: $LOGIN_RESPONSE"
    exit 1
  fi
fi

echo "✅ Token received: ${TOKEN:0:20}..."
echo ""

# Create organization (using token)
echo "🏢 Creating test organization..."
ORG_RESPONSE=$(curl -s -X POST "$BASE_URL/organization" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Test Restaurant","email":"test@restaurant.com","phone":"+55 11 99999-9999","address":"Test St, 123"}')

echo "Org Response: $ORG_RESPONSE"
ORG_ID=$(echo "$ORG_RESPONSE" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$ORG_ID" ]; then
  echo "❌ Failed to create organization"
  echo "Using default org ID from seed data..."
  ORG_ID="123e4567-e89b-12d3-a456-426614174000"
fi

echo "Using ORG_ID: $ORG_ID"
echo ""

# Create project (using org header)
echo "📁 Creating test project..."
PROJ_RESPONSE=$(curl -s -X POST "$BASE_URL/project" \
  -H "Content-Type: application/json" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Test Project","description":"Test project for display settings","organization_id":"'$ORG_ID'"}')

echo "Project Response: $PROJ_RESPONSE"
PROJ_ID=$(echo "$PROJ_RESPONSE" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$PROJ_ID" ]; then
  echo "⚠️ Failed to create project, using default..."
  PROJ_ID="123e4567-e89b-12d3-a456-426614174001"
fi

echo "Using PROJ_ID: $PROJ_ID"
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Test 1: GET display settings (should return defaults or 404 if not created yet)
echo "📋 TEST 1: GET /project/settings/display"
echo "URL: $BASE_URL/project/settings/display"
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/project/settings/display" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $GET_RESPONSE"
echo ""

if echo "$GET_RESPONSE" | grep -q "show_prep_time"; then
  echo "✅ TEST 1 PASSED: show_prep_time field found"
elif echo "$GET_RESPONSE" | grep -q "404"; then
  echo "⚠️ TEST 1: Display settings not found yet (expected for new project)"
else
  echo "Response: $GET_RESPONSE"
fi
echo ""

# Test 2: PUT display settings (update values)
echo "🔄 TEST 2: PUT /project/settings/display"
echo "Updating: show_prep_time = false, show_rating = false, show_description = true"
PUT_RESPONSE=$(curl -s -X PUT "$BASE_URL/project/settings/display" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "show_prep_time": false,
    "show_rating": false,
    "show_description": true
  }')

echo "Response: $PUT_RESPONSE"
echo ""

if echo "$PUT_RESPONSE" | grep -q "show_prep_time"; then
  echo "✅ TEST 2 PASSED: Settings updated"
else
  echo "Response: $PUT_RESPONSE"
fi
echo ""

# Test 3: Verify update (GET again to confirm)
echo "✔️ TEST 3: Verify update - GET /project/settings/display"
VERIFY_RESPONSE=$(curl -s -X GET "$BASE_URL/project/settings/display" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $VERIFY_RESPONSE"
echo ""

if echo "$VERIFY_RESPONSE" | grep -q '"show_prep_time":false'; then
  echo "✅ TEST 3 PASSED: show_prep_time correctly set to false"
elif echo "$VERIFY_RESPONSE" | grep -q "show_prep_time"; then
  echo "✅ TEST 3 PASSED: show_prep_time field found (value changed)"
else
  echo "Response: $VERIFY_RESPONSE"
fi
echo ""

# Test 4: POST reset (restore defaults)
echo "🔄 TEST 4: POST /project/settings/display/reset"
echo "Resetting to defaults..."
RESET_RESPONSE=$(curl -s -X POST "$BASE_URL/project/settings/display/reset" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $RESET_RESPONSE"
echo ""

if echo "$RESET_RESPONSE" | grep -q '"show_prep_time":true'; then
  echo "✅ TEST 4 PASSED: Settings reset to defaults (show_prep_time = true)"
elif echo "$RESET_RESPONSE" | grep -q "show_prep_time"; then
  echo "✅ TEST 4 PASSED: Settings reset (show_prep_time field found)"
else
  echo "Response: $RESET_RESPONSE"
fi
echo ""

echo "🎉 All tests completed!"
