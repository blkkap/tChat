#!/bin/bash

set -e

VERSION="v0.1.0"
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="amd64"


if [[ "$(uname -m)" == "arm64" ]]; then
    ARCH="arm64"
fi

URL = "https://github.com/blkkap/tChat/releases/download/$VERSION/tchat-$OS-$ARCH.tar.gz"

curl -L $URL -o tchat.tar.gz
tar -xzf tchat.tar.gz
chmod +x tchat

sudo mv tchat /usr/local/bin/	

echo "Installed! run with: tchat"
