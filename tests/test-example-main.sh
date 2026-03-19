#!/bin/bash
set -e

# Configuration
ORG_ID="test-tf"
API_TOKEN="f1l1v68jvs2j8ix.34fvctzav5t46kdnchztxz6u5ajfxt5wobs4iulv"
API_URL="http://localhost:3003/v1"
PROVIDER_PATH=$(pwd)

# Reset DB to ensure clean state
make test-reset

# Create a temporary directory for the test
TEST_DIR=$(mktemp -d)
echo "Testing in $TEST_DIR"
export ASDF_TERRAFORM_VERSION=1.10.0

# Configure local plugin mirror
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
PLUGIN_DIR="$TEST_DIR/plugins/registry.terraform.io/gmpinder/pangolin/0.1.0/${OS}_${ARCH}"
mkdir -p "$PLUGIN_DIR"
cp "$PROVIDER_PATH/terraform-provider-pangolin" "$PLUGIN_DIR/terraform-provider-pangolin_v0.1.0"

# Copy main example
cp "$PROVIDER_PATH/examples/main.tf" "$TEST_DIR/main.tf"
cp "$PROVIDER_PATH/.tool-versions" "$TEST_DIR/"

# Replace placeholders with test values
sed -i '' "s/YOUR_API_TOKEN/$API_TOKEN/g" "$TEST_DIR/main.tf"
sed -i '' "s/your-org-id/$ORG_ID/g" "$TEST_DIR/main.tf"

# Inject the base_url
sed -i '' 's|token = "'"$API_TOKEN"'"|token = "'"$API_TOKEN"'"\n  base_url = "'"$API_URL"'"|' "$TEST_DIR/main.tf"

cd "$TEST_DIR"
ls -la
cat main.tf
echo "Initializing Terraform..."
terraform init -plugin-dir="$TEST_DIR/plugins"

echo "Applying Terraform..."
terraform apply -auto-approve

echo "Destroying Terraform..."
terraform destroy -auto-approve

echo "Test successful!"
rm -rf "$TEST_DIR"
