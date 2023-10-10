#!/usr/bin/env sh

######
# Improved script to download the latest release from GitHub and install it at /usr/local/bin.
#
# Author: Thales Pinheiro
######

######
# Variables & Setup
######

ORG_NAME="thalesfsp"
APP_NAME="configurer"
BIN_DIR="${BIN_DIR:-/usr/local/bin}"

if [ $# -gt 0 ]; then
  BIN_DIR=$1
fi

### Logging & Helper Functions ###

log() {
  printf "[%s] %s\n" "$(date +"%Y-%m-%d %H:%M:%S")" "$*"
}

info() {
  log "INFO: $*"
}

warn() {
  log "WARNING: $*"
}

error_exit() {
  log "ERROR: $*"
  exit 1
}

clean_up() {
  info "Cleaning up temporary directory: $tmp_dir"
  rm -rf "$tmp_dir"
}

trap clean_up EXIT

check_dependency() {
  command -v "$1" >/dev/null 2>&1 || error_exit "Command not found: $1"
}

has() {
  command -v "$1" 1>/dev/null 2>&1
}

######
# Main Execution
######

# Check dependencies
check_dependency curl
check_dependency tar
check_dependency mktemp
check_dependency uname

# Check if sudo is available and provide a warning if not
SUDO=""
[ ! has "sudo" ] && warn "sudo not found. Please run the script with appropriate permissions if required."

# Get the latest release version from GitHub.
version=$(curl -s https://api.github.com/repos/${ORG_NAME}/${APP_NAME}/releases/latest | grep tag_name | cut -d '"' -f 4)

# Detect the architecture.
arch=$(uname -m)
case $arch in
  x86_64)
    arch="amd64"
    ;;
  arm64)
    arch="arm64"
    ;;
  *)
    error_exit "Unsupported architecture: $arch"
    ;;
esac

# Detect the OS.
os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux*)
    os="linux"
    ;;
  darwin*)
    os="darwin"
    ;;
  *)
    error_exit "Unsupported operating system: $os"
    ;;
esac

# Fetcher function
fetcher() {
    if has "curl"; then
        printf "%s" "curl -L --fail --silent --show-error -o"
    elif has "wget"; then
        printf "%s" "wget --quiet --output-document"
    else
        error_exit "curl or wget is required"
    fi
}

# Remove "v" from the version string.
versionWithoutV=${version#v}

# Parse URL.
final_url=$(printf "https://github.com/%s/%s/releases/download/%s/%s_%s_%s_%s.tar.gz" "$ORG_NAME" "$APP_NAME" "$version" "$APP_NAME" "$versionWithoutV" "$os" "$arch")

# Create a temp directory.
tmp_dir=$(mktemp -d)

info "Architecture: $arch"
info "OS: $os"
info "Temporary Filepath: $tmp_dir/$APP_NAME.tar.gz"
info "Tarball URL: $final_url"

# Download the latest release using fetcher
info "Downloading $final_url"
eval "$(fetcher)" "$tmp_dir/$APP_NAME.tar.gz" "$final_url"

# Unpack the archive in a temp directory.
info "Unpacking archive"
tar -xzf "$tmp_dir/$APP_NAME.tar.gz" -C "$tmp_dir"

# Move the binary to BIN_DIR, use sudo only if necessary.
if [ -w "$BIN_DIR" ]; then
  info "Moving binary to $BIN_DIR"
  mv "$tmp_dir/$APP_NAME" "$BIN_DIR"
else
  info "Moving binary to $BIN_DIR using sudo"
  $SUDO mv "$tmp_dir/$APP_NAME" "$BIN_DIR"
fi

# Notify the user of successful installation.
info "$APP_NAME installed successfully"
