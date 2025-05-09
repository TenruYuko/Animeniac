#!/bin/bash

# Set the data directory path
DATA_DIR="/aeternae/functional/dockers/animechanical/data"

# Create data directory if it doesn't exist
mkdir -p "$DATA_DIR"

echo "======================================================="
echo "Starting Seanime with data directory: $DATA_DIR"
echo "======================================================="
echo "If you see 'invalid token' errors:"  
echo "1. Click on the logout button (if available)"
echo "2. Then log in again to AniList through the app"
echo "3. If errors persist, try deleting the database:"
echo "   rm \"$DATA_DIR/seanime.db\""
echo "   and then restart the application"
echo "======================================================="

# Run with verbose logging and listen on all interfaces
./seanime --datadir "/aeternae/functional/dockers/animechanical/data"
