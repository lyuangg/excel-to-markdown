#!/bin/bash

# macOS installation helper script
# This script helps remove the quarantine attribute from the downloaded binary

set -e

BINARY_NAME="excel-to-markdown"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🔧 macOS Installation Helper"
echo ""

# Check if binary exists
if [ ! -f "${BINARY_NAME}" ]; then
    echo "❌ Error: ${BINARY_NAME} not found in current directory"
    echo ""
    echo "Please:"
    echo "1. Download the binary from: https://github.com/lyuangg/excel-to-markdown/releases"
    echo "2. Extract it to this directory"
    echo "3. Run this script again"
    exit 1
fi

# Check if quarantine attribute exists
if xattr -l "${BINARY_NAME}" 2>/dev/null | grep -q "com.apple.quarantine"; then
    echo "📦 Removing quarantine attribute from ${BINARY_NAME}..."
    xattr -d com.apple.quarantine "${BINARY_NAME}"
    echo "✅ Quarantine attribute removed successfully!"
else
    echo "ℹ️  No quarantine attribute found. The binary is ready to use."
fi

# Make it executable
chmod +x "${BINARY_NAME}"

echo ""
echo "✅ Installation complete!"
echo ""
echo "You can now run: ./${BINARY_NAME} -h"

