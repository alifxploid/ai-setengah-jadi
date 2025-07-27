#!/bin/bash

# Simple endpoint test script
# Removed set -e to continue testing even if endpoints fail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="https://lipdev.id"
ACCESS_KEY="NBzbgYqdG6oErJPOKzi4JkaFK3eka8C5TPcz4uLikuY="
JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM2OTIzMDYsImlhdCI6MTc1MzYwNTkwNiwidXNlcl9pZCI6ImFkbWluLTAwMSJ9.fjet7ZZzqxOgJMdpyJN9RU1EPFWYKPvmV-Ay55q9JkA"

# Initialize counters
SUCCESS_COUNT=0
ERROR_COUNT=0
FAILED_COUNT=0
TOTAL_COUNT=0

echo -e "${BLUE}🔥 GRYT Backend Simple Test${NC}"
echo -e "${YELLOW}Base URL: $BASE_URL${NC}"
echo ""

# Test function
test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local headers="$4"
    local name="$5"
    
    echo -e "${BLUE}Testing: $name${NC}"
    
    local cmd="curl -s -w '%{http_code}' -o /dev/null"
    
    if [ "$method" = "POST" ] && [ -n "$data" ]; then
        cmd="$cmd -X POST -H 'Content-Type: application/json' -d '$data'"
    elif [ "$method" = "POST" ]; then
        cmd="$cmd -X POST"
    elif [ "$method" = "DELETE" ]; then
        cmd="$cmd -X DELETE"
    fi
    
    if [ -n "$headers" ]; then
        cmd="$cmd $headers"
    fi
    
    cmd="$cmd '$BASE_URL$endpoint'"
    
    local response_code=$(eval $cmd 2>/dev/null || echo "000")
    
    # Count results
    if [[ "$response_code" =~ ^2[0-9][0-9]$ ]]; then
        echo -e "${GREEN}✅ $name: HTTP $response_code - SUCCESS${NC}"
        ((SUCCESS_COUNT++))
    elif [[ "$response_code" =~ ^[45][0-9][0-9]$ ]]; then
        echo -e "${YELLOW}⚠️  $name: HTTP $response_code - ERROR${NC}"
        ((ERROR_COUNT++))
    else
        echo -e "${RED}❌ $name: HTTP $response_code - FAILED${NC}"
        ((FAILED_COUNT++))
    fi
    
    ((TOTAL_COUNT++))
    
    sleep 1.0
}

echo -e "${GREEN}🚀 Starting endpoint tests...${NC}"
echo ""

# 1. Health Check
test_endpoint "GET" "/health" "" "" "health-check"

# 2. Validate Access Key
test_endpoint "POST" "/api/auth/validate-key" "{\"access_key\":\"$ACCESS_KEY\"}" "" "validate-access-key"

# 3. Login
test_endpoint "POST" "/api/auth/login" "{\"email\":\"admin@gryt.id\",\"password\":\"admin123\"}" "" "login"

# 4. Refresh Token
test_endpoint "POST" "/api/auth/refresh" "{\"refresh_token\":\"refresh_token_here\"}" "" "refresh-token"

# Protected endpoints with JWT
AUTH_HEADER="-H 'Authorization: Bearer $JWT_TOKEN'"

# 5. Create AI Chat Session
test_endpoint "POST" "/api/ai/chat/sessions" "{\"title\":\"Test Session\"}" "$AUTH_HEADER" "create-ai-session"

# 6. Get AI Chat Sessions
test_endpoint "GET" "/api/ai/chat/sessions" "" "$AUTH_HEADER" "get-ai-sessions"

# 7. Get AI Chat Session by ID
test_endpoint "GET" "/api/ai/chat/sessions/session-123" "" "$AUTH_HEADER" "get-ai-session-by-id"

# 8. Delete AI Chat Session
test_endpoint "DELETE" "/api/ai/chat/sessions/session-123" "" "$AUTH_HEADER" "delete-ai-session"

# 9. Send AI Message (Text)
test_endpoint "POST" "/api/ai/chat/sessions/session-123/messages" "{\"content\":\"Hello AI\",\"type\":\"text\"}" "$AUTH_HEADER" "send-ai-message-text"

# 10. Get AI Messages
test_endpoint "GET" "/api/ai/chat/sessions/session-123/messages?limit=20&offset=0" "" "$AUTH_HEADER" "get-ai-messages"

# 11. AI Stream Chat
test_endpoint "POST" "/api/ai/chat/sessions/session-123/stream" "{\"content\":\"Stream test\",\"type\":\"text\"}" "$AUTH_HEADER" "ai-stream-chat"

# 12. AI Search
test_endpoint "POST" "/api/ai/search/" "{\"query\":\"test search\",\"limit\":10}" "$AUTH_HEADER" "ai-search"

# 13. Get AI Search History
test_endpoint "GET" "/api/ai/search/history?limit=20&offset=0" "" "$AUTH_HEADER" "get-ai-search-history"

# 14. Get AI Search by ID
test_endpoint "GET" "/api/ai/search/search-123" "" "$AUTH_HEADER" "get-ai-search-by-id"

# 15. Get AI Models
test_endpoint "GET" "/api/ai/models/" "" "$AUTH_HEADER" "get-ai-models"

# 16. Create Legacy Chat Session
test_endpoint "POST" "/api/chat/sessions" "{\"title\":\"Legacy Test\"}" "$AUTH_HEADER" "create-legacy-session"

# 17. Get Legacy Chat Sessions
test_endpoint "GET" "/api/chat/sessions" "" "$AUTH_HEADER" "get-legacy-sessions"

# 18. Send Legacy Message
test_endpoint "POST" "/api/chat/sessions/session-123/messages" "{\"content\":\"Legacy message\"}" "$AUTH_HEADER" "send-legacy-message"

# 19. Get Legacy Messages
test_endpoint "GET" "/api/chat/sessions/session-123/messages" "" "$AUTH_HEADER" "get-legacy-messages"

# 20. Legacy Search
test_endpoint "POST" "/api/search/" "{\"query\":\"legacy search\"}" "$AUTH_HEADER" "legacy-search"

# 21. Get Legacy Search History
test_endpoint "GET" "/api/search/history" "" "$AUTH_HEADER" "get-legacy-search-history"

# 22. Get User Profile
test_endpoint "GET" "/api/user/profile" "" "$AUTH_HEADER" "get-user-profile"

# 23. Get User Tokens
test_endpoint "GET" "/api/user/tokens" "" "$AUTH_HEADER" "get-user-tokens"

echo ""
echo -e "${BLUE}📊 TEST SUMMARY${NC}"
echo -e "${GREEN}✅ SUCCESS: $SUCCESS_COUNT endpoints${NC}"
echo -e "${YELLOW}⚠️  ERRORS: $ERROR_COUNT endpoints${NC}"
echo -e "${RED}❌ FAILED: $FAILED_COUNT endpoints${NC}"
echo -e "${BLUE}📈 TOTAL: $TOTAL_COUNT endpoints tested${NC}"
echo ""

if [ $SUCCESS_COUNT -gt 0 ]; then
    echo -e "${GREEN}🎉 Test completed with $SUCCESS_COUNT successful endpoints!${NC}"
else
    echo -e "${RED}💥 No endpoints working! Check server status!${NC}"
fi