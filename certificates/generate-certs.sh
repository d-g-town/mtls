#!/bin/bash

# Generate certificates for mTLS communication between microservices
set -e

CERT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$CERT_DIR"

echo "Generating certificates for mTLS..."

# Clean up any existing certificates
rm -f *.pem *.key *.crt *.csr

# Generate CA private key
openssl genrsa -out ca-key.pem 4096

# Generate CA certificate
openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem -subj "/C=US/ST=CA/L=San Francisco/O=MTLS Demo/CN=CA"

# Generate server key for service1
openssl genrsa -out service1-key.pem 4096

# Generate server certificate request for service1
openssl req -subj "/C=US/ST=CA/L=San Francisco/O=MTLS Demo/CN=service1" -new -key service1-key.pem -out service1.csr

# Generate server certificate for service1
echo "subjectAltName = DNS:service1,DNS:localhost,IP:127.0.0.1" > service1-extfile.cnf
openssl x509 -req -days 365 -in service1.csr -CA ca.pem -CAkey ca-key.pem -out service1-cert.pem -extfile service1-extfile.cnf -CAcreateserial

# Generate server key for service2
openssl genrsa -out service2-key.pem 4096

# Generate server certificate request for service2
openssl req -subj "/C=US/ST=CA/L=San Francisco/O=MTLS Demo/CN=service2" -new -key service2-key.pem -out service2.csr

# Generate server certificate for service2
echo "subjectAltName = DNS:service2,DNS:localhost,IP:127.0.0.1" > service2-extfile.cnf
openssl x509 -req -days 365 -in service2.csr -CA ca.pem -CAkey ca-key.pem -out service2-cert.pem -extfile service2-extfile.cnf -CAcreateserial

# Generate client key for service1 (when it acts as client)
openssl genrsa -out service1-client-key.pem 4096

# Generate client certificate request for service1
openssl req -subj "/C=US/ST=CA/L=San Francisco/O=MTLS Demo/CN=service1-client" -new -key service1-client-key.pem -out service1-client.csr

# Generate client certificate for service1
openssl x509 -req -days 365 -in service1-client.csr -CA ca.pem -CAkey ca-key.pem -out service1-client-cert.pem -CAcreateserial

# Generate client key for service2 (when it acts as client)
openssl genrsa -out service2-client-key.pem 4096

# Generate client certificate request for service2
openssl req -subj "/C=US/ST=CA/L=San Francisco/O=MTLS Demo/CN=service2-client" -new -key service2-client-key.pem -out service2-client.csr

# Generate client certificate for service2
openssl x509 -req -days 365 -in service2-client.csr -CA ca.pem -CAkey ca-key.pem -out service2-client-cert.pem -CAcreateserial

# Clean up temporary files
rm -f *.csr *-extfile.cnf

echo "Certificates generated successfully!"
echo "Files created:"
ls -la *.pem

# Set appropriate permissions
chmod 600 *-key.pem
chmod 644 *.pem

echo "Certificate generation complete!"