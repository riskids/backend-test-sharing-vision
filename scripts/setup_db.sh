#!/bin/bash

# Database Setup Script for Article Microservice
# This script creates the database and runs migrations

DB_NAME="article_db"
DB_USER="${DB_USER:-root}"
DB_PASS="${DB_PASS:-}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"

echo "========================================"
echo "Article Microservice Database Setup"
echo "========================================"
echo ""
echo "Database Name: $DB_NAME"
echo "Database User: $DB_USER"
echo "Database Host: $DB_HOST"
echo "Database Port: $DB_PORT"
echo ""

# Check if mysql command is available
if ! command -v mysql &> /dev/null; then
    echo "ERROR: MySQL client not found. Please install MySQL client."
    echo ""
    echo "On macOS with Homebrew:"
    echo "  brew install mysql-client"
    echo ""
    echo "On Ubuntu/Debian:"
    echo "  sudo apt-get install mysql-client"
    echo ""
    exit 1
fi

# Create database
echo "Step 1: Creating database '$DB_NAME'..."
if [ -z "$DB_PASS" ]; then
    mysql -u "$DB_USER" -h "$DB_HOST" -P "$DB_PORT" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
else
    mysql -u "$DB_USER" -p"$DB_PASS" -h "$DB_HOST" -P "$DB_PORT" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
fi

if [ $? -eq 0 ]; then
    echo "✓ Database '$DB_NAME' created successfully"
else
    echo "✗ Failed to create database. Please check your MySQL credentials."
    exit 1
fi

echo ""

# Run migration
echo "Step 2: Running database migration..."
MIGRATION_FILE="migrations/000001_create_posts_table.up.sql"

if [ ! -f "$MIGRATION_FILE" ]; then
    echo "ERROR: Migration file not found: $MIGRATION_FILE"
    exit 1
fi

if [ -z "$DB_PASS" ]; then
    mysql -u "$DB_USER" -h "$DB_HOST" -P "$DB_PORT" "$DB_NAME" < "$MIGRATION_FILE"
else
    mysql -u "$DB_USER" -p"$DB_PASS" -h "$DB_HOST" -P "$DB_PORT" "$DB_NAME" < "$MIGRATION_FILE"
fi

if [ $? -eq 0 ]; then
    echo "✓ Migration completed successfully"
else
    echo "✗ Failed to run migration"
    exit 1
fi

echo ""
echo "========================================"
echo "Database setup complete!"
echo "========================================"
echo ""
echo "You can now run the application:"
echo "  make run"
echo "  OR"
echo "  go run cmd/api/main.go"
echo ""
echo "To test the API, run:"
echo "  make test"
echo "  OR"
echo "  ./scripts/test_api.sh"
echo ""