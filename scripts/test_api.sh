#!/bin/bash

# Phase 10: QA & Testing Execution Script
# Article Microservice API Test Suite

BASE_URL="http://localhost:8080"
CONTENT_TYPE="Content-Type: application/json"

echo "========================================"
echo "Article Microservice API Test Suite"
echo "========================================"
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $2"
    else
        echo -e "${RED}✗ FAIL${NC}: $2"
    fi
}

# Function to check HTTP status
check_status() {
    expected=$1
    actual=$2
    if [ "$expected" -eq "$actual" ]; then
        return 0
    else
        return 1
    fi
}

echo "========================================"
echo "10.2. Scenario 1: Validation Rejection"
echo "========================================"
echo ""

# Test 1: Validation Rejection - Invalid data
echo "Test: POST /article/ with invalid data (should fail validation)"
echo "Request Body: { title: 'too short', content: 'too short', category: 'ab', status: 'invalid_status' }"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X POST "$BASE_URL/article/" \
  -H "$CONTENT_TYPE" \
  -d '{
    "title": "too short",
    "content": "too short",
    "category": "ab",
    "status": "invalid_status"
  }')

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"
echo ""

# Check if validation errors are present
if grep -q "errors" /tmp/response.txt || grep -q "error" /tmp/response.txt; then
    print_result 0 "Validation error messages returned"
else
    print_result 1 "No validation error messages found"
fi

check_status 400 $HTTP_CODE
print_result $? "HTTP Status is 400 Bad Request"
echo ""

echo "========================================"
echo "10.3. Scenario 2: Successful Article Creation"
echo "========================================"
echo ""

# Test 2: Create Article - Valid data
echo "Test: POST /article/ with valid data"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X POST "$BASE_URL/article/" \
  -H "$CONTENT_TYPE" \
  -d '{
    "title": "This is a valid title for the article test",
    "content": "This is a valid content that absolutely needs to exceed the two hundred characters minimum limit requirement as specified in the backend test sharing vision document to ensure the validator passes this specific field correctly.",
    "category": "Technology",
    "status": "draft"
  }')

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 200 $HTTP_CODE
print_result $? "HTTP Status is 200 OK"
echo ""

# Test 3: Create another article for pagination test
echo "Test: POST /article/ with second article"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X POST "$BASE_URL/article/" \
  -H "$CONTENT_TYPE" \
  -d '{
    "title": "Second article title for pagination testing here",
    "content": "Second article content that also needs to exceed the two hundred characters minimum limit to pass validation. This is additional content to meet the requirement. We need more characters here to ensure it passes.",
    "category": "Science",
    "status": "publish"
  }')

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 200 $HTTP_CODE
print_result $? "Second article created successfully"
echo ""

echo "========================================"
echo "10.4. Scenario 3: Fetch All Articles with Pagination"
echo "========================================"
echo ""

# Test 4: Get all articles
echo "Test: GET /article/10/0 (limit=10, offset=0)"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X GET "$BASE_URL/article/10/0")

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 200 $HTTP_CODE
print_result $? "HTTP Status is 200 OK"

# Check if response is an array
if grep -q "\[" /tmp/response.txt; then
    print_result 0 "Response is a JSON array"
else
    print_result 1 "Response is not a JSON array"
fi
echo ""

# Test 5: Get articles with pagination
echo "Test: GET /article/1/1 (limit=1, offset=1)"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X GET "$BASE_URL/article/1/1")

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 200 $HTTP_CODE
print_result $? "HTTP Status is 200 OK"
echo ""

echo "========================================"
echo "10.5. Scenario 4: Fetch Single Article"
echo "========================================"
echo ""

# Test 6: Get article by ID
echo "Test: GET /article/1"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X GET "$BASE_URL/article/1")

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 200 $HTTP_CODE
print_result $? "HTTP Status is 200 OK"
echo ""

# Test 7: Get non-existent article
echo "Test: GET /article/999 (non-existent ID)"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X GET "$BASE_URL/article/999")

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 404 $HTTP_CODE
print_result $? "HTTP Status is 404 Not Found"
echo ""

echo "========================================"
echo "10.6. Scenario 5: Update Article"
echo "========================================"
echo ""

# Test 8: Update article
echo "Test: PUT /article/1 with valid data"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X PUT "$BASE_URL/article/1" \
  -H "$CONTENT_TYPE" \
  -d '{
    "title": "Updated title for the article test that is long enough",
    "content": "Updated content for the article test that is long enough to pass the two hundred characters validation rule enforced by the backend golang microservice architecture.",
    "category": "Science",
    "status": "publish"
  }')

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 200 $HTTP_CODE
print_result $? "HTTP Status is 200 OK"
echo ""

# Test 9: Update with invalid status
echo "Test: PUT /article/1 with invalid status (should fail validation)"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X PUT "$BASE_URL/article/1" \
  -H "$CONTENT_TYPE" \
  -d '{
    "title": "Updated title for the article test that is long enough",
    "content": "Updated content for the article test that is long enough to pass the two hundred characters validation rule enforced by the backend golang microservice architecture.",
    "category": "Science",
    "status": "deleted"
  }')

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 400 $HTTP_CODE
print_result $? "HTTP Status is 400 Bad Request (validation error)"
echo ""

echo "========================================"
echo "10.7. Scenario 6: Delete Article"
echo "========================================"
echo ""

# Test 10: Delete article
echo "Test: DELETE /article/1"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X DELETE "$BASE_URL/article/1")

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

check_status 200 $HTTP_CODE
print_result $? "HTTP Status is 200 OK"
echo ""

# Test 11: Delete same article again (idempotency check)
echo "Test: DELETE /article/1 again (should return error or 404)"
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X DELETE "$BASE_URL/article/1")

echo "Response Body:"
cat /tmp/response.txt
echo ""
echo "HTTP Status Code: $HTTP_CODE"

# Note: Depending on implementation, this might return 200 or 404
echo "Note: DELETE on non-existent ID returned HTTP $HTTP_CODE"
echo ""

echo "========================================"
echo "TEST SUMMARY"
echo "========================================"
echo ""
echo "All tests completed. Please review the results above."
echo "If all tests pass, the implementation is correct."
echo ""
echo "To manually verify database state, run:"
echo "  mysql -u root -p article_db -e 'SELECT * FROM posts;'"
echo ""