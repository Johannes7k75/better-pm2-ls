#!/bin/sh

# Create the bin directory if it doesn't exist
mkdir -p bin

# Define OS targets
OS_TARGETS=("windows" "linux" "darwin")

# Define Architecture targets corresponding to each OS
ARCH_TARGETS_windows=("386" "amd64")
ARCH_TARGETS_linux=("386" "amd64" "arm" "arm64")
ARCH_TARGETS_darwin=("amd64" "arm64")

# Iterate over OS targets
for GOOS in "${OS_TARGETS[@]}"; do
    # Get corresponding architectures based on the OS
    ARCH_LIST="ARCH_TARGETS_$GOOS[@]"
    for GOARCH in "${!ARCH_LIST}"; do
        # Set binary name
        BINARY_NAME="pmls-$GOOS-$GOARCH"
        
        echo "Building $BINARY_NAME"

        # Add .exe extension for Windows builds
        if [ "$GOOS" = "windows" ]; then
            BINARY_NAME="$BINARY_NAME.exe"
        fi

        # Build the Go binary (assuming main.go is your entry file)
        go build -trimpath -o "bin/$BINARY_NAME" main.go

        echo "Built for $GOOS ($GOARCH) as bin/$BINARY_NAME"
    done
done
