#!/bin/bash

set -e

BINARY_NAME="telecode"
INSTALL_DIR="$HOME/.local/bin"
CONFIG_DIR="$HOME/.telecode"
REPO="anomalyco/telecode"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    case "$OS" in
        linux|darwin)
            PLATFORM="${OS}-${ARCH}"
            ;;
        *)
            echo -e "${RED}❌ Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac
    
    echo "$PLATFORM"
}

# Download binary from GitHub releases
download_binary() {
    local platform=$1
    local download_url="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${platform}"
    
    echo "📥 Downloading ${BINARY_NAME} for ${platform}..."
    
    if command -v curl &> /dev/null; then
        curl -fsSL "$download_url" -o "/tmp/${BINARY_NAME}"
    elif command -v wget &> /dev/null; then
        wget -q "$download_url" -O "/tmp/${BINARY_NAME}"
    else
        echo -e "${RED}❌ curl or wget is required${NC}"
        exit 1
    fi
    
    chmod +x "/tmp/${BINARY_NAME}"
}

# Install from local binary (for development)
install_local() {
    if [ ! -f "./${BINARY_NAME}" ]; then
        echo -e "${RED}❌ Error: ${BINARY_NAME} binary not found in current directory${NC}"
        echo "Usage:"
        echo "  Remote install: curl -sSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash"
        echo "  Local install:  ./install.sh --local (requires built binary in current directory)"
        exit 1
    fi
    
    echo "📦 Installing local binary..."
    cp "./${BINARY_NAME}" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/${BINARY_NAME}"
}

echo "🚀 Installing ${BINARY_NAME}..."

# Parse arguments
LOCAL_INSTALL=false
if [ "$1" == "--local" ]; then
    LOCAL_INSTALL=true
fi

# Create directories
echo "📁 Creating directories..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"

# Install binary
if [ "$LOCAL_INSTALL" = true ]; then
    install_local
else
    PLATFORM=$(detect_platform)
    download_binary "$PLATFORM"
    mv "/tmp/${BINARY_NAME}" "$INSTALL_DIR/"
fi

echo "📦 Binary installed to ${INSTALL_DIR}/${BINARY_NAME}"

# Create config file if it doesn't exist
if [ ! -f "$CONFIG_DIR/config.yml" ]; then
    echo "⚙️  Creating config file..."
    
    if [ "$LOCAL_INSTALL" = true ] && [ -f "./telecode.yml" ]; then
        cp "./telecode.yml" "$CONFIG_DIR/config.yml"
    else
        cat > "$CONFIG_DIR/config.yml" << 'EOF'
# Telecode Multi-Bot Configuration
# Each workspace represents a separate project with its own bot

workspaces:
  - name: project-a
    working_dir: /home/user/project-a
    bot_token: "YOUR_BOT_TOKEN_1"
    allowed_chats:
      - 123456789
    default_cli: opencode

  - name: project-b
    working_dir: /home/user/project-b
    bot_token: "YOUR_BOT_TOKEN_2"
    allowed_chats:
      - 987654321
    default_cli: claude
EOF
    fi
    
    echo -e "${YELLOW}⚠️  Please edit ${CONFIG_DIR}/config.yml to add your bot tokens${NC}"
else
    echo -e "${GREEN}✅ Config file already exists at ${CONFIG_DIR}/config.yml${NC}"
fi

# Check PATH
echo ""
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo -e "${YELLOW}⚠️  Warning: ${INSTALL_DIR} is not in your PATH${NC}"
    echo "Add this to your ~/.bashrc or ~/.zshrc:"
    echo -e "${BLUE}   export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
    echo ""
fi

echo -e "${GREEN}✅ Installation complete!${NC}"
echo ""
echo "📋 Installed:"
echo "   Binary: ${INSTALL_DIR}/${BINARY_NAME}"
echo "   Config: ${CONFIG_DIR}/config.yml"
echo ""
echo "📝 Next steps:"
echo "   1. Edit ${CONFIG_DIR}/config.yml"
echo "   2. Run: ${BINARY_NAME}"
echo ""
echo "📖 Documentation: https://github.com/${REPO}"
