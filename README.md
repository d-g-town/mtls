# mTLS Microservices Demo

This project demonstrates two Go microservices communicating via mutual TLS (mTLS). Each service acts as both a server and client, with proper certificate-based authentication.

## Architecture

- **Service1**: Runs on port 8443, provides user data and can call Service2
- **Service2**: Runs on port 8444, provides metrics and can call Service1
- **mTLS**: Both services require client certificates for authentication

## Features

- ✅ Mutual TLS authentication between services
- ✅ Certificate-based security (CA, server certs, client certs)
- ✅ Docker containerization
- ✅ Health checks and monitoring endpoints
- ✅ Cross-service communication examples

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Make (optional, for convenient commands)
- OpenSSL (for certificate generation)

### Option 1: Using Make (Recommended)

```bash
# Generate certificates, build, and run everything
make setup

# Test the mTLS communication
make test

# View logs
make logs

# Stop services
make stop

# Clean up everything
make clean
```

### Option 2: Manual Steps

```bash
# 1. Generate certificates
chmod +x certificates/generate-certs.sh
cd certificates && ./generate-certs.sh && cd ..

# 2. Build and run services
docker-compose up -d

# 3. Test the services (wait a few seconds for startup)
curl -k --cert certificates/service1-client-cert.pem \
     --key certificates/service1-client-key.pem \
     --cacert certificates/ca.pem \
     https://localhost:8443/health
```

## API Endpoints

### Service1 (Port 8443)

- `GET /health` - Health check
- `GET /data` - Returns user data
- `GET /call-service2` - Makes mTLS call to Service2

### Service2 (Port 8444)

- `GET /health` - Health check
- `GET /metrics` - Returns service metrics
- `GET /call-service1` - Makes mTLS call to Service1

## Testing mTLS Communication

Once the services are running, you can test various scenarios:

```bash
# Test Service1 health
curl -k --cert certificates/service1-client-cert.pem \
     --key certificates/service1-client-key.pem \
     --cacert certificates/ca.pem \
     https://localhost:8443/health

# Test Service2 health
curl -k --cert certificates/service2-client-cert.pem \
     --key certificates/service2-client-key.pem \
     --cacert certificates/ca.pem \
     https://localhost:8444/health

# Test Service1 calling Service2
curl -k --cert certificates/service1-client-cert.pem \
     --key certificates/service1-client-key.pem \
     --cacert certificates/ca.pem \
     https://localhost:8443/call-service2

# Test Service2 calling Service1
curl -k --cert certificates/service2-client-cert.pem \
     --key certificates/service2-client-key.pem \
     --cacert certificates/ca.pem \
     https://localhost:8444/call-service1
```

## Certificate Details

The setup includes:

- **CA Certificate**: Root certificate authority
- **Server Certificates**: For each service's HTTPS endpoint
- **Client Certificates**: For each service when making client calls

Files generated:

- `ca.pem` / `ca-key.pem` - Certificate Authority
- `service1-cert.pem` / `service1-key.pem` - Service1 server cert
- `service2-cert.pem` / `service2-key.pem` - Service2 server cert
- `service1-client-cert.pem` / `service1-client-key.pem` - Service1 client cert
- `service2-client-cert.pem` / `service2-client-key.pem` - Service2 client cert

## Security Features

1. **Mutual Authentication**: Both client and server verify each other's certificates
2. **Certificate Validation**: All certificates are validated against the CA
3. **Encrypted Communication**: All traffic between services is TLS encrypted
4. **Access Control**: Only services with valid certificates can communicate

## Troubleshooting

### Services won't start

- Check certificate generation: `ls -la certificates/`
- Verify Docker is running: `docker ps`
- Check logs: `docker-compose logs`

### mTLS connection failures

- Ensure certificates were generated properly
- Check certificate permissions (should be readable)
- Verify service names match certificate CN fields

### Testing issues

- Wait for services to fully start (10-15 seconds)
- Check if ports 8443/8444 are available
- Verify certificate paths in curl commands

## Development

### Building individual services

```bash
# Build Service1
cd service1
docker build -t service1 .

# Build Service2
cd service2
docker build -t service2 .
```

### Running without Docker

```bash
# Terminal 1 - Service1
cd service1
go run main.go

# Terminal 2 - Service2
cd service2
go run main.go
```

Note: When running without Docker, update certificate paths in the Go code to point to `../certificates/` instead of `/certs/`.

## License

This is a demonstration project for educational purposes.
