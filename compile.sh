#!/bin/bash

# Seanime Compilation Script
# This script compiles both the backend and frontend components of the Seanime application

set -e  # Exit immediately if a command exits with a non-zero status

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=======================================${NC}"
echo -e "${BLUE}  Seanime Compilation Script          ${NC}"
echo -e "${BLUE}=======================================${NC}"

# Set the working directories
ROOT_DIR="$(pwd)"
BACKEND_DIR="$ROOT_DIR"
FRONTEND_DIR="$ROOT_DIR/seanime-web"
BIN_DIR="$ROOT_DIR/bin"

# Create bin directory if it doesn't exist
mkdir -p "$BIN_DIR"

# Compile the backend
echo -e "\n${YELLOW}=== Compiling Backend ===${NC}"
cd "$BACKEND_DIR"
echo -e "${BLUE}> Running go build...${NC}"
go build -o "$BIN_DIR/seanime"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backend compilation successful${NC}"
else
    echo -e "${RED}✗ Backend compilation failed${NC}"
    exit 1
fi

# Check if the frontend directory exists and has package.json
if [ -d "$FRONTEND_DIR" ] && [ -f "$FRONTEND_DIR/package.json" ]; then
    echo -e "\n${YELLOW}=== Building Frontend ===${NC}"
    cd "$FRONTEND_DIR"
    
    # Check if yarn is present
    if command -v yarn &> /dev/null; then
        echo -e "${BLUE}> Running yarn build...${NC}"
        yarn build
        BUILD_STATUS=$?
    # Check if npm is present
    elif command -v npm &> /dev/null; then
        echo -e "${BLUE}> Running npm run build...${NC}"
        npm run build
        BUILD_STATUS=$?
    else
        echo -e "${YELLOW}! Skipping frontend build: Neither yarn nor npm is installed${NC}"
        BUILD_STATUS=0
    fi
    
    if [ $BUILD_STATUS -eq 0 ]; then
        echo -e "${GREEN}✓ Frontend build successful${NC}"
    else
        echo -e "${RED}✗ Frontend build failed${NC}"
        exit 1
    fi
else
    echo -e "\n${YELLOW}! Skipping frontend build: Frontend directory or package.json not found${NC}"
fi

# Go back to root directory
cd "$ROOT_DIR"

echo -e "\n${GREEN}=======================================${NC}"
echo -e "${GREEN}  Compilation Complete!                ${NC}"
echo -e "${GREEN}=======================================${NC}"
echo -e "${BLUE}Executable: ${NC}$BIN_DIR/seanime"
echo -e "${BLUE}Run with:   ${NC}$BIN_DIR/seanime server"
echo -e "${GREEN}=======================================${NC}"
