#!/bin/bash

# Check if a parameter is provided
if [[ "$#" -ne 1 ]]; then
    printf "Usage: %s <database_name.db>\n" "$0" >&2
    exit 1
fi

# Path to the SQLite database file from the first argument
DATABASE="$1"

# Path to the SQL script file
SQL_SCRIPT="world.sql"

# Function to reset or delete the database
reset_database() {
    if [[ -f "$DATABASE" ]]; then
        rm -f "$DATABASE" || {
            printf "Failed to delete the existing database file.\n" >&2
            return 1
        }
    fi
    touch "$DATABASE" || {
        printf "Failed to create a new database file.\n" >&2
        return 1
    }
}

# Execute the reset function
if ! reset_database; then
    exit 1
fi

# Execute the SQL script using sqlite3
if ! sqlite3 "$DATABASE" < "$SQL_SCRIPT"; then
    printf "Failed to execute SQL script.\n" >&2
    exit 1
fi

# Verify and inform the user
printf "SQL script executed successfully.\n"
