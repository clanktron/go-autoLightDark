#!/bin/bash
# Set the target operating systems and architectures
os_archs=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64")

# Loop through the target OS and architectures
for os_arch in "${os_archs[@]}"; do
  os="${os_arch%%/*}"
  arch="${os_arch##*/}"

  # Set the GOOS and GOARCH environment variables
  export GOOS="$os"
  export GOARCH="$arch"

  # Build the application
  go build -o autoLightDark-$os-$arch

  echo "Built for $os/$arch"
done
