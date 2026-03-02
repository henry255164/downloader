#!/bin/bash

# Configuration
HOST_NAME="com.media.downloader.go"
BINARY_NAME="downloader-host"
TARGET_DIR="$HOME/Library/Application Support/Google/Chrome/NativeMessagingHosts"

# Build the Go binary
echo "Building Go host..."
go build -o $BINARY_NAME main.go

# Get absolute path to the binary
BINARY_PATH="$(pwd)/$BINARY_NAME"

# Create the Native Messaging directory if it doesn't exist
mkdir -p "$TARGET_DIR"

# Generate the manifest from template
# NOTE: The user MUST update the EXTENSION_ID in the manifest or via this script
EXTENSION_ID=${1:-"YOUR_EXTENSION_ID_HERE"}

echo "Generating manifest for Extension ID: $EXTENSION_ID"
sed -e "s|HOST_PATH|$BINARY_PATH|" \
    -e "s|EXTENSION_ID|$EXTENSION_ID|" \
    com.media.downloader.go.json.template > "$TARGET_DIR/$HOST_NAME.json"

echo "Native Messaging Host registered at: $TARGET_DIR/$HOST_NAME.json"
echo "Binary path: $BINARY_PATH"
