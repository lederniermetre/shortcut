#!/bin/bash

GITHUB_REPO="lederniermetre/shortcut"
RELEASE_TAG=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | jq -r '.tag_name')

# Define the installation directory
INSTALL_DIR="."

# Determine the current OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map OS and ARCH to GitHub release assets
case $OS in
  "linux")
    OS="Linux"
    ;;
  "darwin")
    OS="Darwin"
    ;;
  *)
    echo "Unsupported operating system: $OS"
    exit 1
    ;;
esac

case $ARCH in
  "x86_64")
    ARCH="amd64"
    ;;
  "aarch64"|"arm64")
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

echo $RELEASE_TAG

# Construct the download URL
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$RELEASE_TAG/shortcut_${OS}_${ARCH}.tar.gz"

# Download and install the release
echo "Downloading $RELEASE_TAG for $OS-$ARCH..."
echo $DOWNLOAD_URL
curl -s -L -o "$INSTALL_DIR/$(basename $DOWNLOAD_URL)" "$DOWNLOAD_URL"

tar xf "$INSTALL_DIR/$(basename $DOWNLOAD_URL)"

chmod u+x "shortcut"

echo "Installed $RELEASE_TAG to $INSTALL_DIR/shortcut"
