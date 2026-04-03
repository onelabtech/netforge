#!/bin/bash

# NetForge Installer
# This script installs NetForge to /usr/local/bin/netforge

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🔨 Building NetForge...${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.21+ first.${NC}"
    exit 1
fi

# Build
go build -o netforge main.go

# Install
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo -e "${BLUE}Permission required to install to $INSTALL_DIR...${NC}"
    sudo mv netforge "$INSTALL_DIR/netforge"
else
    mv netforge "$INSTALL_DIR/netforge"
fi

echo -e "${GREEN}✅ NetForge successfully installed to $INSTALL_DIR/netforge${NC}"
echo -e "You can now run 'netforge doctor google.com' from anywhere!"
