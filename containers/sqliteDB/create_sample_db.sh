#!/bin/bash

# Check if a parameter is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <database_name.db>"
    exit 1
fi

# Path to the SQLite database file from the first argument
DATABASE="$1"

# Path to the SQL script file
SQL_SCRIPT="world.sql"

# Check if the SQLite database file does not exist
if [ ! -f "$DATABASE" ]; then
    echo "Database file not found. Creating new database."
    touch "$DATABASE"
fi

# Execute the SQL script using sqlite3
sqlite3 $DATABASE < $SQL_SCRIPT

# Verify and inform the user
if [ $? -eq 0 ]; then
    echo "SQL script executed successfully."
else
    echo "Failed to execute SQL script."
fi
