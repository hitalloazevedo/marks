#!/usr/bin/env bash

set -e

APP_NAME="marks"
INSTALL_DIR="$HOME/.local/bin"

echo "Building $APP_NAME..."
go build -o "$APP_NAME"

echo "Creating install dir if needed..."
mkdir -p "$INSTALL_DIR"

echo "Installing to $INSTALL_DIR/$APP_NAME..."
mv "$APP_NAME" "$INSTALL_DIR/$APP_NAME"

chmod +x "$INSTALL_DIR/$APP_NAME"

echo ""
echo "Installed successfully!"
echo "Binary: $INSTALL_DIR/$APP_NAME"
echo ""

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo "WARNING: $INSTALL_DIR is not in your PATH"
  echo ""
  echo "Add this to your shell config:"
  echo 'export PATH="$HOME/.local/bin:$PATH"'
fi
