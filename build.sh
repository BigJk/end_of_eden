#!/bin/bash

GOEXE=$(which go)
PKGXEXE=$(which pkgx)
PKGX_GO_VERSION="1.20"

# Function to setup pkgx
function setup_pkgx() {
    eval "$(curl -Ssf https://pkgx.sh)"
    export GOROOT="$(dirname $(dirname $(pkgx +go@1.20 which go)))"
    GOEXE="pkgx go@${PKGX_GO_VERSION}"
}

# Check if go is installed
# - If go is not installed / is not in PATH ask the user if they want to run it via pkgx.
# - If go is not installed but pkgx is, directly use pkgx.
# - If pkgx is not installed, ask the user if they want to install it.
if [[ -z $GOEXE ]]; then
    # Check if pkgx is installed
    if [[ -z $PKGXEXE ]]; then
        # pkgx is not installed
        echo "Go is not installed or is not in PATH."
        read -p "Do you want to install and run it automatically via pkgx (https://pkgx.sh/)? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Installing pkgx..."
            setup_pkgx
        else
            echo "Can't build, exiting..."
            exit 1
        fi
    else
        # pkgx is installed
        setup_pkgx
    fi
fi

# If env EOE_FORCE_PKGX is set, use pkgx
if [[ ! -z $EOE_FORCE_PKGX ]]; then
    setup_pkgx
fi

# Check if windows extension is needed
EXT=""
if [[ $GOOS == "windows" ]]; then
    EXT=".exe"
fi

# Print build info
echo "======================================="
echo "> End of Eden Build Script"
echo "======================================="
echo "Build Time  : $(date -Iseconds)"
echo "Go Version  : $($GOEXE version)"
echo "Build OS    : ${GOOS:=$($GOEXE env GOOS)}"
echo "Build Arch  : ${GOARCH:=$($GOEXE env GOARCH)}"
echo "Go Root     : ${GOROOT:=$($GOEXE env GOROOT)}"
echo "Go Binary   : ${GOEXE:=$($GOEXE env GOEXE)}"
echo "======================================="

# Delete old binaries
rm -rf ./bin
mkdir -p ./bin

# Build binaries
$GOEXE build -o ./bin/end_of_eden$EXT ./cmd/game/
$GOEXE build -o ./bin/end_of_eden_win$EXT ./cmd/game/
$GOEXE build -o ./bin/end_of_eden_browser$EXT ./cmd/game_browser/
$GOEXE build -o ./bin/fuzzy_tester$EXT ./cmd/fuzzy-tester/

# Disable SSH version for now:
# go build -o ./bin/end_of_eden_ssh$EXT ./cmd/game_ssh/

# Copy /assets to /bin
cp -r ./assets ./bin/assets

# Finished!
echo "Done! Binaries are in ./bin/"