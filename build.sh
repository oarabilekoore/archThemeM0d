#!/bin/bash

# ArchThemeM0d Build Script
# Builds the web assets and Go binary

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="archThemeM0d"
BUILD_DIR="build"
WEB_DIR="web"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 is not installed or not in PATH"
        exit 1
    fi
}

# Header
echo -e "${BLUE}"
echo "╔══════════════════════════════════════╗"
echo "║        ArchThemeM0d Build Script     ║"
echo "╚══════════════════════════════════════╝"
echo -e "${NC}"

# Check prerequisites
log_info "Checking prerequisites..."
check_command "bun"
check_command "go"

# Create build directory if it doesn't exist
if [ ! -d "$BUILD_DIR" ]; then
    log_info "Creating build directory..."
    mkdir -p "$BUILD_DIR"
fi

# Build web assets
log_info "Building web assets with Bun..."
if [ -d "$WEB_DIR" ]; then
    cd "$WEB_DIR"

    # Check if package.json exists
    if [ ! -f "package.json" ]; then
        log_warning "No package.json found in $WEB_DIR directory"
        log_warning "Skipping web asset build..."
    else
        # Install dependencies if node_modules doesn't exist
        if [ ! -d "node_modules" ]; then
            log_info "Installing dependencies..."
            bun install
        fi

        # Run the build
        log_info "Running bun build..."
        bun run build
        log_success "Web assets built successfully"
    fi

    cd ..
else
    log_warning "Web directory '$WEB_DIR' not found, skipping web asset build"
fi

# Build Go binary
log_info "Building Go binary..."

# Get version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Detect OS and architecture
OS=$(go env GOOS)
ARCH=$(go env GOARCH)

log_info "Building for ${OS}/${ARCH}..."
log_info "Version: ${VERSION}"
log_info "Build time: ${BUILD_TIME}"
log_info "Git commit: ${GIT_COMMIT}"

# Build the binary
go build -ldflags "${LDFLAGS}" -o "${BUILD_DIR}/${BINARY_NAME}" ./main.go

if [ $? -eq 0 ]; then
    log_success "Go binary built successfully: ${BUILD_DIR}/${BINARY_NAME}"

    # Make binary executable
    chmod +x "${BUILD_DIR}/${BINARY_NAME}"

    # Show binary info
    BINARY_SIZE=$(du -h "${BUILD_DIR}/${BINARY_NAME}" | cut -f1)
    log_info "Binary size: ${BINARY_SIZE}"

    # Test binary
    if ./"${BUILD_DIR}/${BINARY_NAME}" --help >/dev/null 2>&1; then
        log_success "Binary test passed"
    else
        log_warning "Binary test failed - binary might not work correctly"
    fi
else
    log_error "Go build failed"
    exit 1
fi

# Optional: Create release archive
if [ "$1" == "--release" ]; then
    log_info "Creating release archive..."
    ARCHIVE_NAME="${BINARY_NAME}-${VERSION}-${OS}-${ARCH}.tar.gz"

    cd "${BUILD_DIR}"
    tar -czf "${ARCHIVE_NAME}" "${BINARY_NAME}"
    cd ..

    log_success "Release archive created: ${BUILD_DIR}/${ARCHIVE_NAME}"
fi

# Summary
echo
log_success "Build completed successfully!"
echo -e "${GREEN}Output files:${NC}"
echo "  - Binary: ${BUILD_DIR}/${BINARY_NAME}"
if [ "$1" == "--release" ]; then
    echo "  - Archive: ${BUILD_DIR}/${ARCHIVE_NAME}"
fi

echo
log_info "To install the binary system-wide, run:"
echo "  sudo cp ${BUILD_DIR}/${BINARY_NAME} /usr/local/bin/"

echo
log_info "To test the binary, run:"
echo "  ./${BUILD_DIR}/${BINARY_NAME} --help"
