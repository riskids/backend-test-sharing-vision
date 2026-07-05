# Testing Documentation

## Phase 10: QA & Testing Execution Checklist

This document provides comprehensive testing instructions for the Article Microservice API.

---

## Pre-Execution Verification

### 1. Database Setup

Before running the tests, ensure the MySQL database is set up:

```bash
# Option 1: Use the setup script
chmod +x scripts/setup_db.sh
./scripts/setup_db.sh

# Option 2: Manual setup
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS article_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p article_db < migrations/000001_create_posts_table.up.sql
```

### 2. Configure Environment

Update the `.env` file with your MySQL credentials:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=your_password
DB_NAME=article_db
APP_PORT=8080
```

### 3. Run the Application

```bash
# Using Makefile
make run

# Or directly with Go
go run cmd/api/main.go
```

You should see:
```
Successfully connected to MySQL database
Server is running on port :8080
```

---

## Test Scenarios

### 10.2. Scenario 1: Validation Rejection

**Purpose:** Ensure request is blocked BEFORE hitting the database.

**Test:** Send POST with invalid data

**Request:**
```bash
curl -X POST http://localhost:8080/article/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "too short",
    "content": "too short",
    "category": "ab",
    "status": "invalid_status"
  }'
```

**Expected Output:**
- HTTP Status: `400 Bad Request`
- Body contains validation error messages for:
  - Title (min 20 characters)
  - Content (min 200 characters)
  - Category (min 3 characters)
  - Status (must be one of: publish, draft, thrash)

**Database Verification:**
```sql
SELECT * FROM posts;
-- Should return 0 rows (no invalid data inserted)
```

---

### 10.3. Scenario 2: Successful Article Creation

**Test:** Create article with valid data

**Request:**
```bash
curl -X POST http://localhost:8080/article/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "This is a valid title for the article test",
    "content": "This is a valid content that absolutely needs to exceed the two hundred characters minimum limit requirement as specified in the backend test sharing vision document to ensure the validator passes this specific field correctly.",
    "category": "Technology",
    "status": "draft"
  }'
```

**Expected Output:**
- HTTP Status: `200 OK`
- Body: JSON object with created article details including `id`

**Database Verification:**
```sql
SELECT * FROM posts;
-- Should return 1 row with id=1
```

**Create Second Article (for pagination test):**
```bash
curl -X POST http://localhost:8080/article/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Second article title for pagination testing here",
    "content": "Second article content that also needs to exceed the two hundred characters minimum limit to pass validation. This is additional content to meet the requirement. We need more characters here to ensure it passes.",
    "category": "Science",
    "status": "publish"
  }'
```

---

### 10.4. Scenario 3: Fetch All Articles with Pagination

**Test:** Get all articles

**Request:**
```bash
curl -X GET http://localhost:8080/article/10/0
```

**Expected Output:**
- HTTP Status: `200 OK`
- Body: JSON Array containing article objects

**Test with Pagination:**
```bash
curl -X GET http://localhost:8080/article/1/1
```

**Expected Output:**
- HTTP Status: `200 OK`
- Body: JSON Array with 1 element (second article due to offset)

---

### 10.5. Scenario 4: Fetch Single Article

**Test:** Get article by ID

**Request:**
```bash
curl -X GET http://localhost:8080/article/1
```

**Expected Output:**
- HTTP Status: `200 OK`
- Body: JSON object with article details

**Test Non-existent ID:**
```bash
curl -X GET http://localhost:8080/article/999
```

**Expected Output:**
- HTTP Status: `404 Not Found`

---

### 10.6. Scenario 5: Update Article

**Test:** Update article with valid data

**Request:**
```bash
curl -X PUT http://localhost:8080/article/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated title for the article test that is long enough",
    "content": "Updated content for the article test that is long enough to pass the two hundred characters validation rule enforced by the backend golang microservice architecture.",
    "category": "Science",
    "status": "publish"
  }'
```

**Expected Output:**
- HTTP Status: `200 OK`
- Body: `{}`

**Database Verification:**
```sql
SELECT * FROM posts WHERE id=1;
-- Verify category is "Science", status is "publish"
-- Verify updated_date has changed
```

**Test with Invalid Status:**
```bash
curl -X PUT http://localhost:8080/article/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated title for the article test that is long enough",
    "content": "Updated content for the article test that is long enough to pass the two hundred characters validation rule enforced by the backend golang microservice architecture.",
    "category": "Science",
    "status": "deleted"
  }'
```

**Expected Output:**
- HTTP Status: `400 Bad Request`
- Body contains validation error for status

---

### 10.7. Scenario 6: Delete Article

**Test:** Delete article

**Request:**
```bash
curl -X DELETE http://localhost:8080/article/1
```

**Expected Output:**
- HTTP Status: `200 OK`
- Body: `{}`

**Database Verification:**
```sql
SELECT * FROM posts WHERE id=1;
-- Should return 0 rows
```

**Test Idempotency (Delete again):**
```bash
curl -X DELETE http://localhost:8080/article/1
```

**Expected Output:**
- HTTP Status: `200 OK` (or `404 Not Found` depending on implementation)

---

## Running All Tests

### Using the Test Script

```bash
# Make the script executable
chmod +x scripts/test_api.sh

# Run all tests
./scripts/test_api.sh
```

### Using Makefile

```bash
# Run tests
make test
```

---

## Postman Collection

Import `postman/article_collection.json` into Postman for interactive testing.

The collection includes:
1. **Create Article - Valid**: Creates a valid article
2. **Create Article - Invalid**: Tests validation errors
3. **Get All Articles**: Lists all articles with pagination
4. **Get Article by ID**: Retrieves a single article
5. **Update Article**: Updates an existing article
6. **Delete Article**: Removes an article

---

## Final Sign-off Checklist

- [ ] All 5 endpoints return correct HTTP methods and JSON structures
- [ ] Validation successfully protects database from invalid data
- [ ] POST `/article/` returns 400 for invalid data
- [ ] POST `/article/` returns 200 for valid data
- [ ] GET `/article/:limit/:offset` returns array of articles
- [ ] GET `/article/:id` returns single article or 404
- [ ] PUT `/article/:id` updates article and returns 200
- [ ] DELETE `/article/:id` removes article and returns 200
- [ ] Database timestamps (created_date, updated_date) work correctly