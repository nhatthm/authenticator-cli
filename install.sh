#!/usr/bin/env bash

VERSION=${VERSION:-latest}

if [[ "$VERSION" == "latest" ]]; then
    API_PATH="latest"
else
    API_PATH="tags/$VERSION"
fi

INSTALL_SCRIPT=$(
    curl -s "https://api.github.com/repos/nhatthm/authenticator-cli/releases/${API_PATH}" \
        | grep "browser_download_url.*install\.sh" \
        | cut -d : -f 2,3 \
        | tr -d \" \
        | tr -d ' '
)

if [[ -z "$INSTALL_SCRIPT" ]]; then
    echo "could not find install script for version $VERSION"
    exit 1
fi

bash -c "$(curl -fsSL "$INSTALL_SCRIPT")"
