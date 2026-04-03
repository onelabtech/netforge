#!/bin/bash

# NetForge Installer
# This script installs NetForge to /usr/local/bin/netforge

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🔨 Preparing to build NetForge...${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.21+ first.${NC}"
    exit 1
fi

# Create a temporary directory for building
TEMP_DIR=$(mktemp -d)
echo -e "${BLUE}📂 Created temporary build directory: $TEMP_DIR${NC}"

# Clone the repository
echo -e "${BLUE}🔌 Cloning NetForge repository...${NC}"
git clone https://github.com/onelabtech/netforge.git "$TEMP_DIR/netforge" --depth 1

# Build
echo -e "${BLUE}⚙️  Building binary...${NC}"
cd "$TEMP_DIR/netforge"
go build -o netforge main.go

# Install
INSTALL_DIR="/usr/local/bin"
echo -e "${BLUE}📦 Installing to $INSTALL_DIR...${NC}"

if [ ! -w "$INSTALL_DIR" ]; then
    echo -e "${BLUE}🔐 Permission required to install to $INSTALL_DIR (using sudo)...${NC}"
    sudo mv netforge "$INSTALL_DIR/netforge"
else
    mv netforge "$INSTALL_DIR/netforge"
fi

# Cleanup
echo -e "${BLUE}🧹 Cleaning up...${NC}"
rm -rf "$TEMP_DIR"

echo -e "${GREEN}✅ NetForge successfully installed to $INSTALL_DIR/netforge${NC}"
echo -e "You can now run 'netforge doctor google.com' from anywhere!"
