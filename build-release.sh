#!/bin/bash

# Build script for creating release packages
# This script builds binaries for multiple platforms

set -e

VERSION=${1:-"v1.0.0"}
APP_NAME="excel-to-markdown"
BUILD_DIR="dist"
RELEASE_DIR="release"

# Clean previous builds
rm -rf ${BUILD_DIR} ${RELEASE_DIR}
mkdir -p ${BUILD_DIR} ${RELEASE_DIR}

# Build function
build() {
    local GOOS=$1
    local GOARCH=$2
    local EXT=$3
    local OUTPUT="${BUILD_DIR}/${APP_NAME}-${GOOS}-${GOARCH}${EXT}"
    
    echo "Building ${GOOS}/${GOARCH}..."
    GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o ${OUTPUT} .
    
    # Create release package
    local RELEASE_NAME="${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    local RELEASE_PATH="${RELEASE_DIR}/${RELEASE_NAME}"
    mkdir -p ${RELEASE_PATH}
    
    # Copy binary
    cp ${OUTPUT} ${RELEASE_PATH}/${APP_NAME}${EXT}
    
    # Copy documentation
    cp README.md ${RELEASE_PATH}/
    cp README.zh.md ${RELEASE_PATH}/
    cp LICENSE ${RELEASE_PATH}/
    
    # Create archive
    if [ "${GOOS}" = "windows" ]; then
        cd ${RELEASE_DIR}
        zip -r ${RELEASE_NAME}.zip ${RELEASE_NAME}
        cd ..
    else
        cd ${RELEASE_DIR}
        tar -czf ${RELEASE_NAME}.tar.gz ${RELEASE_NAME}
        cd ..
    fi
    
    echo "✓ Built ${RELEASE_NAME}.tar.gz / ${RELEASE_NAME}.zip"
}

# Build for all platforms
echo "Building release packages for version ${VERSION}..."
echo ""

# macOS
build darwin amd64 ""
build darwin arm64 ""

# Linux
build linux amd64 ""
build linux arm64 ""

# Windows
build windows amd64 ".exe"

echo ""
echo "✅ All builds completed!"
echo "Release packages are in the '${RELEASE_DIR}' directory:"
ls -lh ${RELEASE_DIR}/*.{tar.gz,zip} 2>/dev/null || true

