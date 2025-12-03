#!/bin/sh
set -e

REPO="johnhorton/annotr"
BINARY="annotr"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

get_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) echo "unknown" ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) echo "unknown" ;;
    esac
}

main() {
    OS=$(get_os)
    ARCH=$(get_arch)

    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        echo "Error: Unsupported OS or architecture"
        exit 1
    fi

    echo "Detected: ${OS}-${ARCH}"

    LATEST=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST" ]; then
        LATEST="v0.1.0"
        echo "Note: No releases found, using ${LATEST}"
    fi

    if [ "$OS" = "windows" ]; then
        FILENAME="${BINARY}-${OS}-${ARCH}.exe"
    else
        FILENAME="${BINARY}-${OS}-${ARCH}"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST}/${FILENAME}"

    echo "Downloading ${BINARY} ${LATEST}..."
    
    TMP_DIR=$(mktemp -d)
    trap "rm -rf ${TMP_DIR}" EXIT

    if command -v curl >/dev/null 2>&1; then
        curl -sL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${BINARY}"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "${DOWNLOAD_URL}" -O "${TMP_DIR}/${BINARY}"
    else
        echo "Error: curl or wget required"
        exit 1
    fi

    chmod +x "${TMP_DIR}/${BINARY}"

    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    else
        echo "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    fi

    echo "✓ Installed ${BINARY} to ${INSTALL_DIR}/${BINARY}"
    echo ""
    echo "Run 'annotr init' to get started!"
}

main
