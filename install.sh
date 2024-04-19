#!/bin/bash

# Name of the binary
BINARY_NAME="kitty-sesh"

# Destination directory
INSTALL_DIR="/usr/local/bin"

# Temporary directory to clone the repo
TMP_DIR=$(mktemp -d)

# GitHub repository URL
REPO_URL="https://github.com/Raghav-rv28/kitty-sesh.git"

# Clone the repository
echo "Cloning the repository..."
git clone "$REPO_URL" "$TMP_DIR" &>/dev/null

# Check if cloning was successful
if [ $? -ne 0 ]; then
	echo "Error: Failed to clone the repository."
	exit 1
fi

# Navigate to the repository directory
cd "$TMP_DIR"

# Build the Go project
echo "Building $BINARY_NAME..."
go build -o "$BINARY_NAME" ./...

# Check if build was successful
if [ $? -ne 0 ]; then
	echo "Error: Failed to build $BINARY_NAME."
	exit 1
fi

# Move the binary to the install directory
sudo mv "$BINARY_NAME" "$INSTALL_DIR/"

# Make the binary executable
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Clean up
rm -rf "$TMP_DIR"

echo "Installation completed. You can now use '$BINARY_NAME' command."
echo "All sessions are stored at home/<usr>/.config/kitty/sessions"
echo "Read about use here: https://github.com/Raghav-rv28/kitty-sesh/blob/main/README.md"
