# Article Microservice

A Golang-based article microservice using Gin framework, MySQL database, and implementing Controller-Service-Repository pattern.

## Project Structure

```
article-microservice/
├── cmd/api/main.go                 # Application entry point
├── internal/
│   ├── config/config.go            # Configuration management
│   ├── delivery/
│   │   └── http/
│   │       ├── dto/                # Data Transfer Objects
│   │       │   ├── article_request.go
│   │       │   └── article_response.go
│   │       ├── article_handler.go  # HTTP Handlers
│   │       └── router.go           # Routing configuration
│   ├── model/article.go            # Domain models
│   ├── repository/
│   │   └── article_repository.go   # Data access layer
│   └── service/
│       └── article_service.go      # Business logic layer
├── migrations/
│   ├── 000001_create_posts_table.up.sql    # Migration up
│   └── 000001_create_posts_table.down.sql  # Migration down
├── pkg/mysql/db.go                 # MySQL connection
├── postman/
│   └── article_collection.json      # Postman collection
├── .env                             # Environment variables
└── go.mod                           # Go module file
```

## Requirements

- Go 1.20+ (tested with Go 1.25)
- MySQL 5.7+ or 8.0+

## Installation

### 1. Clone and Install Dependencies

```bash
go mod download
```

### 2. Configure Environment

Edit the `.env` file with your MySQL credentials:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=your_password
DB_NAME=article_db
APP_PORT=8080
```

### 3. Create MySQL Database

```sql
CREATE DATABASE article_db;
```

### 4. Run Database Migration

Execute the SQL file located at `migrations/000001_create_posts_table.up.sql`:

```bash
mysql -u root -p article_db < migrations/000001_create_posts_table.up.sql
```

Or run directly in MySQL:

```sql
CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status VARCHAR(100) NOT NULL
);
```

### 5. Run the Application

```bash
go run cmd/api/main.go
```

Or build and run:

```bash
go build -o bin/article-microservice ./cmd/api
./bin/article-microservice
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/article/` | Create a new article |
| GET    | `/article/:limit/:offset` | Get all articles with pagination |
| GET    | `/article/:id` | Get article by ID |
| PUT    | `/article/:id` | Update article by ID |
| DELETE | `/article/:id` | Delete article by ID |

## Validation Rules

The following validation rules are enforced on Create and Update requests:

- **Title**: Required, minimum 20 characters
- **Content**: Required, minimum 200 characters
- **Category**: Required, minimum 3 characters
- **Status**: Required, must be one of: `publish`, `draft`, `thrash`

### Example Request Body

```json
{
    "title": "This is a valid article title with more than 20 characters",
    "content": "This is the content of the article. It must be at least 200 characters long to pass the validation requirement. This content is long enough to meet the minimum character requirement for the content field validation.",
    "category": "Technology",
    "status": "draft"
}
```

## Postman Collection

Import the `postman/article_collection.json` file into Postman to test all API endpoints.

The collection includes:
- **Create Article - Valid**: A valid request that passes validation
- **Create Article - Invalid**: A request that fails validation (demonstrates validation errors)
- **Get All Articles**: Retrieve articles with pagination
- **Get Article by ID**: Retrieve a single article
- **Update Article**: Update an existing article
- **Delete Article**: Remove an article

## Architecture

This project follows the **Controller-Service-Repository (CSR)** pattern:

1. **Handler (Controller)**: Handles HTTP requests/responses and input validation
2. **Service**: Contains business logic
3. **Repository**: Manages database operations

## Technologies Used

- **Gin**: HTTP web framework
- **go-playground/validator/v10**: Request validation
- **go-sql-driver/mysql**: MySQL driver
- **joho/godotenv**: Environment variable management

## Testing

### Automated Testing

Run the test script to execute all API test scenarios:

```bash
# Make the script executable
chmod +x scripts/test_api.sh

# Run all tests
./scripts/test_api.sh
```

### Manual Testing with curl

See [TESTING.md](TESTING.md) for comprehensive test scenarios and expected results.

### Using Makefile

```bash
# Install dependencies
make deps

# Build the application
make build

# Run the application
make run

# Run API tests
make test
```

### Database Setup

```bash
# Make the script executable
chmod +x scripts/setup_db.sh

# Run database setup
./scripts/setup_db.sh
```

## API Documentation

See [TESTING.md](TESTING.md) for detailed API documentation including:
- Request/Response formats
- Validation rules
- Error handling
- All test scenarios
# backend-test-sharing-vision
