#!/bin/bash

# Script to extract certificate content for PaaS deployment
set -e

CERT_DIR="../certificates"
OUTPUT_DIR="./secrets"

echo "🔐 Extracting certificate content for PaaS deployment..."

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Function to create secret file with content
create_secret_file() {
    local cert_file="$1"
    local output_file="$2"
    local description="$3"
    
    if [ -f "$CERT_DIR/$cert_file" ]; then
        echo "📄 Creating $output_file ($description)"
        cat "$CERT_DIR/$cert_file" > "$OUTPUT_DIR/$output_file"
        echo "   ✅ Content saved to secrets/$output_file"
    else
        echo "   ❌ Certificate $cert_file not found!"
        echo "   Run: cd ../certificates && ./generate-certs.sh"
        exit 1
    fi
}

# Function to create base64 encoded version
create_base64_file() {
    local cert_file="$1"
    local output_file="$2"
    local description="$3"
    
    if [ -f "$CERT_DIR/$cert_file" ]; then
        echo "📄 Creating $output_file.b64 ($description - Base64)"
        base64 -i "$CERT_DIR/$cert_file" > "$OUTPUT_DIR/$output_file.b64"
        echo "   ✅ Base64 content saved to secrets/$output_file.b64"
    fi
}

echo ""
echo "🚀 SERVICE1 CERTIFICATES"
echo "========================"
create_secret_file "ca.pem" "service1-ca.pem" "CA Certificate"
create_secret_file "service1-cert.pem" "service1-cert.pem" "Server Certificate"
create_secret_file "service1-key.pem" "service1-key.pem" "Server Private Key"
create_secret_file "service1-client-cert.pem" "service1-client-cert.pem" "Client Certificate"
create_secret_file "service1-client-key.pem" "service1-client-key.pem" "Client Private Key"

echo ""
echo "🚀 SERVICE2 CERTIFICATES"
echo "========================"
create_secret_file "ca.pem" "service2-ca.pem" "CA Certificate"
create_secret_file "service2-cert.pem" "service2-cert.pem" "Server Certificate"
create_secret_file "service2-key.pem" "service2-key.pem" "Server Private Key"
create_secret_file "service2-client-cert.pem" "service2-client-cert.pem" "Client Certificate"
create_secret_file "service2-client-key.pem" "service2-client-key.pem" "Client Private Key"

echo ""
echo "🔄 Creating Base64 versions (for some PaaS providers)..."
create_base64_file "ca.pem" "service1-ca.pem" "CA Certificate"
create_base64_file "service1-cert.pem" "service1-cert.pem" "Server Certificate"
create_base64_file "service1-key.pem" "service1-key.pem" "Server Private Key"
create_base64_file "service1-client-cert.pem" "service1-client-cert.pem" "Client Certificate"
create_base64_file "service1-client-key.pem" "service1-client-key.pem" "Client Private Key"

create_base64_file "ca.pem" "service2-ca.pem" "CA Certificate"
create_base64_file "service2-cert.pem" "service2-cert.pem" "Server Certificate"
create_base64_file "service2-key.pem" "service2-key.pem" "Server Private Key"
create_base64_file "service2-client-cert.pem" "service2-client-cert.pem" "Client Certificate"
create_base64_file "service2-client-key.pem" "service2-client-key.pem" "Client Private Key"

echo ""
echo "📋 SUMMARY"
echo "=========="
echo "Certificate files extracted to: $(pwd)/secrets/"
echo ""
echo "For your PaaS, upload these as secret files:"
echo ""
echo "SERVICE1 Secret Files:"
echo "  - /certs/ca.pem                       ← secrets/service1-ca.pem"
echo "  - /certs/service1-cert.pem           ← secrets/service1-cert.pem"
echo "  - /certs/service1-key.pem            ← secrets/service1-key.pem"
echo "  - /certs/service1-client-cert.pem    ← secrets/service1-client-cert.pem"
echo "  - /certs/service1-client-key.pem     ← secrets/service1-client-key.pem"
echo ""
echo "SERVICE2 Secret Files:"
echo "  - /certs/ca.pem                       ← secrets/service2-ca.pem"
echo "  - /certs/service2-cert.pem           ← secrets/service2-cert.pem"
echo "  - /certs/service2-key.pem            ← secrets/service2-key.pem"
echo "  - /certs/service2-client-cert.pem    ← secrets/service2-client-cert.pem"
echo "  - /certs/service2-client-key.pem     ← secrets/service2-client-key.pem"
echo ""
echo "⚠️  IMPORTANT: Set these environment variables:"
echo ""
echo "SERVICE1:"
echo "  SERVICE2_URL=https://your-service2-url.com"
echo ""
echo "SERVICE2:"
echo "  SERVICE1_URL=https://your-service1-url.com"
echo ""
echo "✅ Ready for PaaS deployment!"

# List the files created
echo ""
echo "📁 Files created:"
ls -la "$OUTPUT_DIR/"