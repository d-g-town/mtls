# Makefile for mTLS microservices

.PHONY: help generate-certs build run clean test logs

# Default target
help:
	@echo "Available targets:"
	@echo "  generate-certs  - Generate mTLS certificates"
	@echo "  build          - Build both microservices"
	@echo "  run            - Run both services with docker-compose"
	@echo "  clean          - Clean up containers and certificates"
	@echo "  test           - Test mTLS communication between services"
	@echo "  logs           - Show logs from both services"
	@echo "  stop           - Stop all services"

# Generate certificates
generate-certs:
	@echo "Generating mTLS certificates..."
	@chmod +x certificates/generate-certs.sh
	@cd certificates && ./generate-certs.sh

# Build services
build:
	@echo "Building microservices..."
	@docker-compose build

# Run services
run:
	@echo "Starting microservices with mTLS..."
	@docker-compose up -d
	@echo "Services started:"
	@echo "  Service1: https://localhost:8443"
	@echo "  Service2: https://localhost:8444"

# Clean up
clean:
	@echo "Cleaning up..."
	@docker-compose down -v
	@docker system prune -f
	@rm -f certificates/*.pem certificates/*.crt certificates/*.key certificates/*.srl

# Test mTLS communication
test:
	@echo "Testing mTLS communication..."
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "\n=== Testing Service1 Health ==="
	@curl -k --cert certificates/service1-client-cert.pem \
		--key certificates/service1-client-key.pem \
		--cacert certificates/ca.pem \
		https://localhost:8443/health || echo "Service1 health check failed"
	@echo "\n=== Testing Service2 Health ==="
	@curl -k --cert certificates/service2-client-cert.pem \
		--key certificates/service2-client-key.pem \
		--cacert certificates/ca.pem \
		https://localhost:8444/health || echo "Service2 health check failed"
	@echo "\n=== Testing Service1 -> Service2 Communication ==="
	@curl -k --cert certificates/service1-client-cert.pem \
		--key certificates/service1-client-key.pem \
		--cacert certificates/ca.pem \
		https://localhost:8443/call-service2 || echo "Service1->Service2 call failed"
	@echo "\n=== Testing Service2 -> Service1 Communication ==="
	@curl -k --cert certificates/service2-client-cert.pem \
		--key certificates/service2-client-key.pem \
		--cacert certificates/ca.pem \
		https://localhost:8444/call-service1 || echo "Service2->Service1 call failed"

# Show logs
logs:
	@docker-compose logs -f

# Stop services
stop:
	@docker-compose down

# Full setup (certificates + build + run)
setup: generate-certs build run
	@echo "Full setup complete!"