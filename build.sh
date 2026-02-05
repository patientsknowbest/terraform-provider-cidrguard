#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Clean up previous Terraform environment and state
echo "Cleaning up Terraform environment..."
rm -rf .terraform
rm -f .terraform.lock.hcl
rm -f terraform.tfstate
rm -f terraform.tfstate.backup

# Build the provider binary with the conventional name
echo "Building the provider..."
go build -o terraform-provider-cidrguard

# Set up the local provider directory that mirrors the final registry path
PROVIDER_DIR="$HOME/.terraform.d/plugins/patientsknowbest/cidrguard/1.0.0/$(go env GOOS)_$(go env GOARCH)"
echo "Installing provider to $PROVIDER_DIR..."
mkdir -p "$PROVIDER_DIR"

# Move the binary to the local provider directory
mv terraform-provider-cidrguard "$PROVIDER_DIR/"

echo "Build and installation complete."

echo "Setup complete. You can now run 'terraform apply'."
