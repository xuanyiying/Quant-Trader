#!/bin/bash

# Configuration
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DB_SERVICE_NAME="timescaledb"
DB_NAME="postgres"
DB_USER="postgres"
export PGPASSWORD="password"

echo "🚀 Starting SQL migration process..."
echo "📂 Script directory: $SCRIPT_DIR"

# Find container ID for the timescaledb service

if [ -z "$DB_CONTAINER" ]; then
    echo "❌ Error: Could not find a running container for service '${DB_SERVICE_NAME}'."
    echo "Make sure you have run 'docker compose up -d' first."
    exit 1
fi

echo "📦 Found database container: $DB_CONTAINER"

# 1. Apply init.sql
echo "--------------------------------------"
echo "📄 Applying base schema: init.sql"
docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < "$SCRIPT_DIR/init.sql"
if [ $? -eq 0 ]; then
    echo "✅ init.sql applied successfully."
else
    echo "❌ Error applying init.sql"
    exit 1
fi

# 2. Apply migrations in order
echo "--------------------------------------"
echo "📁 Applying migrations from migrations/ directory..."

# Get all .sql files in migrations/, sort them, and apply
for file in $(ls "$SCRIPT_DIR/migrations"/*.sql | sort); do
    echo "📄 Applying migration: $(basename "$file")"
    docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < "$file"
    
    if [ $? -eq 0 ]; then
        echo "✅ $(basename "$file") applied successfully."
    else
        echo "❌ Error applying $file. Stopping."
        exit 1
    fi
done

echo "--------------------------------------"
echo "🎉 All SQL scripts have been executed successfully!"
