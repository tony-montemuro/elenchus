#!/bin/bash

echo "📦 Installing project dependencies..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ Failed to install dependencies. Ensure Go is properly installed on your system."
    exit 1
else
    echo "✅ Dependencies installed!"
fi

echo -e "\n"
echo "🔐 Creating local TLS certificate..."

# Create tls directory if it does not exist
echo "Creating the ./tls directory..."
mkdir -p tls
cd tls

# Generate TLS certificate
echo "Generating tls certificate in ./tls directory..."
go run "$(go env GOROOT)/src/crypto/tls/generate_cert.go" --rsa-bits=2048 --host=localhost
if [ $? -eq 0 ]; then
    echo "✅ TLS certificate generated!"
else
    echo "❌ Failed to generate TLS certificate. Ensure Go is properly installed on your system."
    exit 1
fi

echo -e "\n"
echo "========================================"
echo -e "\n"
echo "✅ Project initialization complete! Use:"
echo -e "go run ./cmd/web/\n"
echo -e "ℹ️ To view all configurable options, run the above command with the -help flag.\n"
echo "⚠️ Your browser will show a security warning for using a self-signed certificate."
echo "This is normal for development, just accept the warning to proceed."
