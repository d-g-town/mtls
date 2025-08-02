#!/bin/bash

# Test script for mTLS microservices
set -e

echo "🔧 mTLS Microservices Test Script"
echo "=================================="

# Check if certificates exist
if [ ! -f "certificates/ca.pem" ]; then
    echo "❌ Certificates not found. Generating them now..."
    cd certificates && ./generate-certs.sh && cd ..
else
    echo "✅ Certificates found"
fi

# Check if services are running
echo "🔍 Checking if services are running..."
if ! docker ps | grep -q "service1\|service2"; then
    echo "❌ Services not running. Starting them now..."
    docker-compose up -d
    echo "⏳ Waiting 15 seconds for services to start..."
    sleep 15
else
    echo "✅ Services are running"
fi

echo ""
echo "🧪 Testing mTLS Communication"
echo "=============================="

# Test Service1 health
echo "📡 Testing Service1 health endpoint..."
if curl -s -k --cert certificates/service1-client-cert.pem \
        --key certificates/service1-client-key.pem \
        --cacert certificates/ca.pem \
        https://localhost:8443/health > /dev/null; then
    echo "✅ Service1 health check passed"
else
    echo "❌ Service1 health check failed"
fi

# Test Service2 health
echo "📡 Testing Service2 health endpoint..."
if curl -s -k --cert certificates/service2-client-cert.pem \
        --key certificates/service2-client-key.pem \
        --cacert certificates/ca.pem \
        https://localhost:8444/health > /dev/null; then
    echo "✅ Service2 health check passed"
else
    echo "❌ Service2 health check failed"
fi

# Test Service1 -> Service2 communication
echo "🔄 Testing Service1 -> Service2 communication..."
if curl -s -k --cert certificates/service1-client-cert.pem \
        --key certificates/service1-client-key.pem \
        --cacert certificates/ca.pem \
        https://localhost:8443/call-service2 > /dev/null; then
    echo "✅ Service1 -> Service2 communication successful"
else
    echo "❌ Service1 -> Service2 communication failed"
fi

# Test Service2 -> Service1 communication
echo "🔄 Testing Service2 -> Service1 communication..."
if curl -s -k --cert certificates/service2-client-cert.pem \
        --key certificates/service2-client-key.pem \
        --cacert certificates/ca.pem \
        https://localhost:8444/call-service1 > /dev/null; then
    echo "✅ Service2 -> Service1 communication successful"
else
    echo "❌ Service2 -> Service1 communication failed"
fi

echo ""
echo "📋 Sample API Calls"
echo "==================="

echo "🔗 Service1 data endpoint:"
curl -s -k --cert certificates/service1-client-cert.pem \
     --key certificates/service1-client-key.pem \
     --cacert certificates/ca.pem \
     https://localhost:8443/data | jq '.'

echo ""
echo "📊 Service2 metrics endpoint:"
curl -s -k --cert certificates/service2-client-cert.pem \
     --key certificates/service2-client-key.pem \
     --cacert certificates/ca.pem \
     https://localhost:8444/metrics | jq '.'

echo ""
echo "✨ All tests completed!"
echo "📖 Check README.md for more detailed usage instructions."