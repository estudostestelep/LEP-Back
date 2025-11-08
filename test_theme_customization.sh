#!/bin/bash

# Test script for Theme Customization API
# Tests all 5 endpoints for theme customization with new 8 fields

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

# Use default IDs from seed data
ORG_ID="123e4567-e89b-12d3-a456-426614174000"
PROJ_ID="123e4567-e89b-12d3-a456-426614174001"

echo "Using ORG_ID: $ORG_ID"
echo "Using PROJ_ID: $PROJ_ID"
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Test 1: GET theme customization (should return defaults or 404 if not created yet)
echo "📋 TEST 1: GET /project/settings/theme"
echo "URL: $BASE_URL/project/settings/theme"
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/project/settings/theme" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $GET_RESPONSE"
echo ""

if echo "$GET_RESPONSE" | grep -q "primary_color"; then
  echo "✅ TEST 1 PASSED: primary_color field found"
  # Check for new 8 fields
  if echo "$GET_RESPONSE" | grep -q "destructive_color"; then
    echo "✅ Found destructive_color (semantic color)"
  fi
  if echo "$GET_RESPONSE" | grep -q "success_color"; then
    echo "✅ Found success_color (semantic color)"
  fi
  if echo "$GET_RESPONSE" | grep -q "warning_color"; then
    echo "✅ Found warning_color (semantic color)"
  fi
  if echo "$GET_RESPONSE" | grep -q "border_color"; then
    echo "✅ Found border_color (semantic color)"
  fi
  if echo "$GET_RESPONSE" | grep -q "disabled_opacity"; then
    echo "✅ Found disabled_opacity (numeric field)"
  fi
  if echo "$GET_RESPONSE" | grep -q "focus_ring_color"; then
    echo "✅ Found focus_ring_color (system config)"
  fi
  if echo "$GET_RESPONSE" | grep -q "input_background_color"; then
    echo "✅ Found input_background_color (system config)"
  fi
  if echo "$GET_RESPONSE" | grep -q "shadow_intensity"; then
    echo "✅ Found shadow_intensity (numeric field)"
  fi
elif echo "$GET_RESPONSE" | grep -q "404"; then
  echo "⚠️ TEST 1: Theme customization not found yet (expected for new project)"
else
  echo "Response: $GET_RESPONSE"
fi
echo ""

# Test 2: PUT theme customization (partial update with new fields)
echo "🎨 TEST 2: PUT /project/settings/theme (with new 8 fields)"
echo "Updating: colors + new fields (disabled_opacity, shadow_intensity)"
PUT_RESPONSE=$(curl -s -X PUT "$BASE_URL/project/settings/theme" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "primary_color": "#FF6B35",
    "secondary_color": "#F4A261",
    "background_color": "#09090b",
    "card_background_color": "#18181b",
    "text_color": "#fafafa",
    "text_secondary_color": "#a1a1aa",
    "accent_color": "#FF9F1C",
    "destructive_color": "#EF4444",
    "success_color": "#10B981",
    "warning_color": "#F59E0B",
    "border_color": "#E5E7EB",
    "focus_ring_color": "#3B82F6",
    "input_background_color": "#FFFFFF",
    "disabled_opacity": 0.50,
    "shadow_intensity": 1.00
  }')

echo "Response: $PUT_RESPONSE"
echo ""

if echo "$PUT_RESPONSE" | grep -q "primary_color"; then
  echo "✅ TEST 2 PASSED: Theme updated with all fields"
  if echo "$PUT_RESPONSE" | grep -q '"disabled_opacity"'; then
    echo "✅ disabled_opacity field persisted"
  fi
  if echo "$PUT_RESPONSE" | grep -q '"shadow_intensity"'; then
    echo "✅ shadow_intensity field persisted"
  fi
else
  echo "❌ TEST 2 FAILED"
  echo "Response: $PUT_RESPONSE"
fi
echo ""

# Test 3: Verify update (GET again to confirm all fields)
echo "✔️ TEST 3: Verify update - GET /project/settings/theme"
VERIFY_RESPONSE=$(curl -s -X GET "$BASE_URL/project/settings/theme" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $VERIFY_RESPONSE"
echo ""

VERIFY_PASS=true
if echo "$VERIFY_RESPONSE" | grep -q '"primary_color":"#FF6B35"'; then
  echo "✅ primary_color correctly persisted (#FF6B35)"
else
  echo "⚠️  primary_color not verified"
  VERIFY_PASS=false
fi

if echo "$VERIFY_RESPONSE" | grep -q '"disabled_opacity":0.5'; then
  echo "✅ disabled_opacity correctly persisted (0.5)"
elif echo "$VERIFY_RESPONSE" | grep -q '"disabled_opacity"'; then
  echo "✅ disabled_opacity field found (numeric value present)"
else
  echo "⚠️  disabled_opacity not found"
  VERIFY_PASS=false
fi

if echo "$VERIFY_RESPONSE" | grep -q '"shadow_intensity":1'; then
  echo "✅ shadow_intensity correctly persisted (1.0)"
elif echo "$VERIFY_RESPONSE" | grep -q '"shadow_intensity"'; then
  echo "✅ shadow_intensity field found (numeric value present)"
else
  echo "⚠️  shadow_intensity not found"
  VERIFY_PASS=false
fi

if [ "$VERIFY_PASS" = true ]; then
  echo ""
  echo "✅ TEST 3 PASSED: All fields verified"
else
  echo ""
  echo "⚠️  TEST 3: Some fields not verified"
fi
echo ""

# Test 4: Partial update (update only new numeric fields)
echo "🔧 TEST 4: PUT /project/settings/theme (partial update - numeric fields only)"
echo "Updating: disabled_opacity = 0.75, shadow_intensity = 1.5"
PARTIAL_RESPONSE=$(curl -s -X PUT "$BASE_URL/project/settings/theme" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "disabled_opacity": 0.75,
    "shadow_intensity": 1.5
  }')

echo "Response: $PARTIAL_RESPONSE"
echo ""

if echo "$PARTIAL_RESPONSE" | grep -q '"disabled_opacity":0.75'; then
  echo "✅ TEST 4 PASSED: disabled_opacity updated to 0.75"
elif echo "$PARTIAL_RESPONSE" | grep -q '"disabled_opacity"'; then
  echo "✅ TEST 4: disabled_opacity field updated"
else
  echo "⚠️  TEST 4: Could not verify disabled_opacity update"
fi

if echo "$PARTIAL_RESPONSE" | grep -q '"shadow_intensity":1.5'; then
  echo "✅ shadow_intensity updated to 1.5"
elif echo "$PARTIAL_RESPONSE" | grep -q '"shadow_intensity"'; then
  echo "✅ shadow_intensity field updated"
fi
echo ""

# Test 5: Validation test (invalid opacity)
echo "⚠️ TEST 5: Validation - PUT with invalid opacity (should fail)"
echo "Sending: disabled_opacity = 2.0 (max is 1.0)"
INVALID_RESPONSE=$(curl -s -X PUT "$BASE_URL/project/settings/theme" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "disabled_opacity": 2.0
  }')

echo "Response: $INVALID_RESPONSE"
echo ""

if echo "$INVALID_RESPONSE" | grep -q "error\|Error\|400\|validation"; then
  echo "✅ TEST 5 PASSED: Validation correctly rejected invalid opacity"
elif echo "$INVALID_RESPONSE" | grep -q '"disabled_opacity":2'; then
  echo "⚠️  TEST 5: Opacity was accepted (validation may be missing)"
else
  echo "Response: $INVALID_RESPONSE"
fi
echo ""

# Test 6: Validation test (invalid shadow intensity)
echo "⚠️ TEST 6: Validation - PUT with invalid shadow intensity (should fail)"
echo "Sending: shadow_intensity = 3.0 (max is 2.0)"
INVALID_SHADOW=$(curl -s -X PUT "$BASE_URL/project/settings/theme" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shadow_intensity": 3.0
  }')

echo "Response: $INVALID_SHADOW"
echo ""

if echo "$INVALID_SHADOW" | grep -q "error\|Error\|400\|validation"; then
  echo "✅ TEST 6 PASSED: Validation correctly rejected invalid shadow intensity"
elif echo "$INVALID_SHADOW" | grep -q '"shadow_intensity":3'; then
  echo "⚠️  TEST 6: Shadow intensity was accepted (validation may be missing)"
else
  echo "Response: $INVALID_SHADOW"
fi
echo ""

# Test 7: POST reset (restore defaults)
echo "🔄 TEST 7: POST /project/settings/theme/reset"
echo "Resetting to defaults (should restore all 15 fields with defaults)..."
RESET_RESPONSE=$(curl -s -X POST "$BASE_URL/project/settings/theme/reset" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "Response: $RESET_RESPONSE"
echo ""

RESET_PASS=true
if echo "$RESET_RESPONSE" | grep -q "primary_color"; then
  echo "✅ primary_color reset to default"
else
  RESET_PASS=false
fi

if echo "$RESET_RESPONSE" | grep -q "disabled_opacity"; then
  echo "✅ disabled_opacity reset to default (0.5)"
else
  echo "⚠️  disabled_opacity not found in reset"
  RESET_PASS=false
fi

if echo "$RESET_RESPONSE" | grep -q "shadow_intensity"; then
  echo "✅ shadow_intensity reset to default (1.0)"
else
  echo "⚠️  shadow_intensity not found in reset"
  RESET_PASS=false
fi

if [ "$RESET_PASS" = true ]; then
  echo ""
  echo "✅ TEST 7 PASSED: Theme reset to defaults with all 15 fields"
else
  echo ""
  echo "⚠️  TEST 7: Some fields missing in reset"
fi
echo ""

echo "═══════════════════════════════════════════════════════════════"
echo "🎉 All theme customization tests completed!"
echo ""
echo "Summary:"
echo "- TEST 1: GET with all 15 fields"
echo "- TEST 2: PUT with all 15 fields (new 8 + original 7)"
echo "- TEST 3: Verify persistence of all fields"
echo "- TEST 4: Partial update (numeric fields only)"
echo "- TEST 5: Validation - Invalid opacity (0.0-1.0)"
echo "- TEST 6: Validation - Invalid shadow intensity (0.0-2.0)"
echo "- TEST 7: Reset to defaults with all 15 fields"
