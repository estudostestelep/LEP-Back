#!/bin/bash

# Test script for Menu Name Uniqueness Validation
# Tests that menu names must be unique within a project
# and validates the 409 Conflict response when duplicates are attempted

set -e

# Base URL
BASE_URL="http://localhost:8080"

# Test user email/password
TEST_EMAIL="test@example.com"
TEST_PASS="testpass123"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "═══════════════════════════════════════════════════════════════"
echo "🍽️  Menu Name Uniqueness Validation Tests"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Try to login with test user
echo "🔑 Logging in with test user..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASS\"}")

# Extract token from login response
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4 || true)

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token - user may not exist"
  echo "Creating test user..."

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

# TEST 1: Create first menu with unique name
echo "📋 TEST 1: POST /menu - Create first menu with name 'Almoço'"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/menu" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Almoço",
    "description": "Menu do almoço",
    "order": 1,
    "active": true
  }')

echo "Response: $CREATE_RESPONSE"
echo ""

MENU1_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4 || true)

if [ -z "$MENU1_ID" ]; then
  echo -e "${RED}❌ TEST 1 FAILED: Could not create first menu${NC}"
  echo "Response: $CREATE_RESPONSE"
else
  echo -e "${GREEN}✅ TEST 1 PASSED: Menu created with ID: ${MENU1_ID:0:20}...${NC}"
fi
echo ""

# TEST 2: Try to create duplicate menu with same name (should fail with 409)
echo "🚫 TEST 2: POST /menu - Try to create duplicate menu with same name (should fail)"
echo "Expected: 409 Conflict error"
DUPLICATE_RESPONSE=$(curl -s -X POST "$BASE_URL/menu" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Almoço",
    "description": "Attempting duplicate menu name",
    "order": 2,
    "active": true
  }')

echo "Response: $DUPLICATE_RESPONSE"
echo ""

if echo "$DUPLICATE_RESPONSE" | grep -q "409\|Conflict\|already exists"; then
  echo -e "${GREEN}✅ TEST 2 PASSED: Duplicate creation correctly rejected${NC}"
  if echo "$DUPLICATE_RESPONSE" | grep -q "already exists in this project"; then
    echo -e "${GREEN}   ✓ Error message is correct${NC}"
  fi
else
  echo -e "${RED}❌ TEST 2 FAILED: Duplicate should have been rejected${NC}"
fi
echo ""

# TEST 3: Create second menu with different name (should succeed)
echo "📋 TEST 3: POST /menu - Create second menu with different name 'Jantar'"
MENU2_RESPONSE=$(curl -s -X POST "$BASE_URL/menu" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jantar",
    "description": "Menu da noite",
    "order": 2,
    "active": true
  }')

echo "Response: $MENU2_RESPONSE"
echo ""

MENU2_ID=$(echo "$MENU2_RESPONSE" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4 || true)

if [ -z "$MENU2_ID" ]; then
  echo -e "${RED}❌ TEST 3 FAILED: Could not create second menu with different name${NC}"
else
  echo -e "${GREEN}✅ TEST 3 PASSED: Second menu created successfully${NC}"
fi
echo ""

# TEST 4: Try to rename existing menu to duplicate name (should fail with 409)
if [ -n "$MENU2_ID" ]; then
  echo "🔄 TEST 4: PUT /menu/{id} - Try to rename menu to existing name (should fail)"
  echo "Attempting to rename Menu 2 from 'Jantar' to 'Almoço' (taken name)"
  UPDATE_DUPLICATE=$(curl -s -X PUT "$BASE_URL/menu/$MENU2_ID" \
    -H "X-Lpe-Organization-Id: $ORG_ID" \
    -H "X-Lpe-Project-Id: $PROJ_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Almoço",
      "description": "Attempting to rename to existing name",
      "order": 2,
      "active": true
    }')

  echo "Response: $UPDATE_DUPLICATE"
  echo ""

  if echo "$UPDATE_DUPLICATE" | grep -q "409\|Conflict\|already exists"; then
    echo -e "${GREEN}✅ TEST 4 PASSED: Update to duplicate name correctly rejected${NC}"
  else
    echo -e "${RED}❌ TEST 4 FAILED: Update should have been rejected${NC}"
  fi
  echo ""
fi

# TEST 5: Rename menu to a new name (should succeed)
if [ -n "$MENU2_ID" ]; then
  echo "✏️  TEST 5: PUT /menu/{id} - Rename menu to new unique name"
  echo "Renaming Menu 2 from 'Jantar' to 'Café da Manhã'"
  UPDATE_VALID=$(curl -s -X PUT "$BASE_URL/menu/$MENU2_ID" \
    -H "X-Lpe-Organization-Id: $ORG_ID" \
    -H "X-Lpe-Project-Id: $PROJ_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Café da Manhã",
      "description": "Menu da manhã",
      "order": 0,
      "active": true
    }')

  echo "Response: $UPDATE_VALID"
  echo ""

  if echo "$UPDATE_VALID" | grep -q "Café da Manhã"; then
    echo -e "${GREEN}✅ TEST 5 PASSED: Menu successfully renamed to new name${NC}"
  else
    echo -e "${RED}❌ TEST 5 FAILED: Menu rename to unique name should succeed${NC}"
  fi
  echo ""
fi

# TEST 6: Case-insensitive uniqueness check
echo "🔤 TEST 6: POST /menu - Verify case-insensitive uniqueness"
echo "Attempting to create menu 'almoço' (lowercase) when 'Almoço' exists"
CASE_DUPLICATE=$(curl -s -X POST "$BASE_URL/menu" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "almoço",
    "description": "Testing case-insensitive uniqueness",
    "order": 3,
    "active": true
  }')

echo "Response: $CASE_DUPLICATE"
echo ""

if echo "$CASE_DUPLICATE" | grep -q "409\|Conflict\|already exists"; then
  echo -e "${GREEN}✅ TEST 6 PASSED: Case-insensitive uniqueness working${NC}"
else
  echo -e "${YELLOW}⚠️  TEST 6: Case-insensitive check may not be working${NC}"
fi
echo ""

# TEST 7: List menus and verify created menus
echo "📊 TEST 7: GET /menu - List all menus for project"
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/menu" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $LIST_RESPONSE"
echo ""

# Count menus in response (should have at least 2: Almoço and Café da Manhã)
MENU_COUNT=$(echo "$LIST_RESPONSE" | grep -o '"name"' | wc -l)
echo "Total menus in project: $MENU_COUNT"

if [ "$MENU_COUNT" -ge 2 ]; then
  echo -e "${GREEN}✅ TEST 7 PASSED: Multiple menus exist in project${NC}"
else
  echo -e "${YELLOW}⚠️  TEST 7: Expected at least 2 menus, found $MENU_COUNT${NC}"
fi
echo ""

# TEST 8: Verify menu can keep its own name when updated
if [ -n "$MENU1_ID" ]; then
  echo "🔁 TEST 8: PUT /menu/{id} - Menu can keep its own name when updated"
  echo "Updating Menu 1 with same name 'Almoço' (should succeed)"
  UPDATE_SAME_NAME=$(curl -s -X PUT "$BASE_URL/menu/$MENU1_ID" \
    -H "X-Lpe-Organization-Id: $ORG_ID" \
    -H "X-Lpe-Project-Id: $PROJ_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Almoço",
      "description": "Updated description for lunch menu",
      "order": 1,
      "active": true
    }')

  echo "Response: $UPDATE_SAME_NAME"
  echo ""

  if echo "$UPDATE_SAME_NAME" | grep -q "Almoço" && ! echo "$UPDATE_SAME_NAME" | grep -q "409\|Conflict"; then
    echo -e "${GREEN}✅ TEST 8 PASSED: Menu can keep its own name when updated${NC}"
  else
    echo -e "${RED}❌ TEST 8 FAILED: Menu should be able to keep its own name${NC}"
  fi
  echo ""
fi

echo "═══════════════════════════════════════════════════════════════"
echo "🎉 Menu uniqueness validation tests completed!"
echo ""
echo "Summary:"
echo "- TEST 1: Create first menu with unique name"
echo "- TEST 2: Reject duplicate menu creation (409 Conflict)"
echo "- TEST 3: Create second menu with different name"
echo "- TEST 4: Reject update to existing name (409 Conflict)"
echo "- TEST 5: Allow update to new unique name"
echo "- TEST 6: Case-insensitive uniqueness check"
echo "- TEST 7: List menus and verify creations"
echo "- TEST 8: Menu can keep its own name when updated"
echo "═══════════════════════════════════════════════════════════════"
