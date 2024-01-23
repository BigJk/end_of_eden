#!/bin/bash

GOEXE=$(which go)
PKGXEXE=$(which pkgx)
PKGX_GO_VERSION="1.21.6"

# Function to setup pkgx
function setup_pkgx() {
    eval "$(curl -Ssf https://pkgx.sh)"
    export GOROOT="$(dirname $(dirname $(pkgx +go@$PKGX_GO_VERSION which go)))"
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

# Create bin folder
mkdir -p ./bin

# Build binaries
echo "1. Building terminal version..."
$GOEXE build -o ./bin/end_of_eden$EXT ./cmd/game/

echo "2. Building windowed gl version..."
$GOEXE build -o ./bin/end_of_eden_win$EXT ./cmd/game_win/

echo "3. Building testing util..."
$GOEXE build -o ./bin/tester$EXT ./cmd/internal/tester/

echo "4. Building fuzzy testing util..."
$GOEXE build -o ./bin/fuzzy_tester$EXT ./cmd/internal/fuzzy_tester/

echo "5. Building wasm version..."
GOOS=js GOARCH=wasm $GOEXE build -o ./bin/eoe.wasm ./cmd/game_wasm/
cp "$($GOEXE env GOROOT)/misc/wasm/wasm_exec.js" "./bin/wasm_exec.js"
cp ./cmd/game_wasm/index.html ./bin/index.html

# Disable SSH version for now:
# go build -o ./bin/end_of_eden_ssh$EXT ./cmd/game_ssh/

# Build asset index
echo "6. Building asset index..."
./internal/misc/build_index.sh ./assets

# Copy /assets to /bin
echo "7. Copying assets..."
cp -r ./assets ./bin/assets

# Finished!
echo "Done! Binaries are in ./bin/"