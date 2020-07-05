#!/bin/sh
set -e

if [ -z "$1" ]; then
    echo
    echo "** A version is required"
    echo
    echo "./update_rust.sh <version>"
    echo
    exit -1
fi

mkdir -p .rust-bins
cd .rust-bins/

# Download
if [ ! -f "rust-$1-x86_64-unknown-linux-gnu.tar.gz" ]; then
    curl -L -O "https://static.rust-lang.org/dist/rust-$1-x86_64-unknown-linux-gnu.tar.gz"
fi

if [ ! -f "rustc-$1-src.tar.gz" ]; then
    curl -L -O "https://static.rust-lang.org/dist/rustc-$1-src.tar.gz"
fi

# Calculate Hashes
BIN_HASH=$(shasum -a 256 "rust-$1-x86_64-unknown-linux-gnu.tar.gz" | awk '{print $1}')
SRC_HASH=$(shasum -a 256 "rustc-$1-src.tar.gz" | awk '{print $1}')

cat <<EOF
  [[metadata.dependencies]]
    id = "rust"
    name = "Rust Standalone Installer"
    version = "$1"
    uri = "https://static.rust-lang.org/dist/rust-$1-x86_64-unknown-linux-gnu.tar.gz"
    sha256 = "$BIN_HASH"
    source = "https://static.rust-lang.org/dist/rustc-$1-src.tar.gz"
    source_256 = "$SRC_HASH"
    stacks = ["io.paketo.stacks.tiny", "io.buildpacks.stacks.bionic", "org.cloudfoundry.stacks.cflinuxfs3"]
EOF