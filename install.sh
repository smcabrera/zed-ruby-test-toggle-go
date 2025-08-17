#!/usr/bin/env bash

# go-zed-test-toggle installation script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default installation directory
DEFAULT_INSTALL_DIR="/usr/local/bin"
INSTALL_DIR="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"

# Binary name
BINARY_NAME="go-zed-test-toggle"

# Print colored output
print_error() {
    echo -e "${RED}ERROR: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}SUCCESS: $1${NC}"
}

print_info() {
    echo -e "${YELLOW}INFO: $1${NC}"
}

# Check if running as root for system-wide installation
check_permissions() {
    if [[ "$INSTALL_DIR" == "/usr/local/bin" ]] && [[ $EUID -ne 0 ]]; then
        print_error "Installing to $INSTALL_DIR requires root privileges."
        print_info "Run with sudo: sudo ./install.sh"
        print_info "Or install to user directory: INSTALL_DIR=~/bin ./install.sh"
        exit 1
    fi
}

# Check if Go is installed
check_go_installation() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        print_info "Please install Go from https://golang.org/dl/"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}')
    print_info "Found Go version: $GO_VERSION"
}

# Create install directory if it doesn't exist
create_install_dir() {
    if [[ ! -d "$INSTALL_DIR" ]]; then
        print_info "Creating installation directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi
}

# Build the binary
build_binary() {
    print_info "Building $BINARY_NAME..."

    # Get version information
    VERSION="${VERSION:-1.0.0}"
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

    # Build with version information
    LDFLAGS="-s -w"
    LDFLAGS="$LDFLAGS -X main.Version=$VERSION"
    LDFLAGS="$LDFLAGS -X main.Commit=$COMMIT"
    LDFLAGS="$LDFLAGS -X main.BuildTime=$BUILD_TIME"

    if go build -ldflags "$LDFLAGS" -o "$BINARY_NAME"; then
        print_success "Binary built successfully"
    else
        print_error "Failed to build binary"
        exit 1
    fi
}

# Install the binary
install_binary() {
    print_info "Installing $BINARY_NAME to $INSTALL_DIR..."

    # Expand tilde in path
    EXPANDED_INSTALL_DIR="${INSTALL_DIR/#\~/$HOME}"

    # Copy binary to installation directory
    if cp "$BINARY_NAME" "$EXPANDED_INSTALL_DIR/"; then
        # Make it executable
        chmod +x "$EXPANDED_INSTALL_DIR/$BINARY_NAME"
        print_success "Installation completed!"
    else
        print_error "Failed to copy binary to $INSTALL_DIR"
        exit 1
    fi
}

# Check if installation directory is in PATH
check_path() {
    EXPANDED_INSTALL_DIR="${INSTALL_DIR/#\~/$HOME}"

    if [[ ":$PATH:" != *":$EXPANDED_INSTALL_DIR:"* ]]; then
        print_info "WARNING: $INSTALL_DIR is not in your PATH"
        print_info "Add it to your PATH by adding this line to your shell configuration:"

        # Detect shell
        SHELL_NAME=$(basename "$SHELL")
        case "$SHELL_NAME" in
            bash)
                CONFIG_FILE="~/.bashrc"
                ;;
            zsh)
                CONFIG_FILE="~/.zshrc"
                ;;
            fish)
                CONFIG_FILE="~/.config/fish/config.fish"
                print_info "set -gx PATH $INSTALL_DIR \$PATH"
                return
                ;;
            *)
                CONFIG_FILE="your shell configuration file"
                ;;
        esac

        print_info "export PATH=\"$INSTALL_DIR:\$PATH\""
        print_info "Then reload your shell configuration or start a new terminal session."
    fi
}

# Verify installation
verify_installation() {
    EXPANDED_INSTALL_DIR="${INSTALL_DIR/#\~/$HOME}"

    if [[ -x "$EXPANDED_INSTALL_DIR/$BINARY_NAME" ]]; then
        print_success "$BINARY_NAME has been installed successfully!"

        # Show version if binary is in PATH
        if command -v "$BINARY_NAME" &> /dev/null; then
            print_info "Installed version:"
            "$BINARY_NAME" version
        else
            print_info "Run the following to see the version:"
            print_info "$EXPANDED_INSTALL_DIR/$BINARY_NAME version"
        fi
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Main installation process
main() {
    print_info "Installing go-zed-test-toggle..."
    echo

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --install-dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --help|-h)
                echo "Usage: ./install.sh [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --install-dir DIR    Installation directory (default: $DEFAULT_INSTALL_DIR)"
                echo "  --help, -h          Show this help message"
                echo ""
                echo "Environment variables:"
                echo "  INSTALL_DIR         Alternative way to set installation directory"
                echo "  VERSION             Set version number (default: 1.0.0)"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done

    # Run installation steps
    check_go_installation
    check_permissions
    create_install_dir
    build_binary
    install_binary
    check_path
    echo
    verify_installation

    echo
    print_success "Installation complete!"
    print_info "Usage: $BINARY_NAME lookup -p <file> -r <project-root>"
}

# Run main function
main "$@"
