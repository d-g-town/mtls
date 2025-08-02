# PaaS Deployment Guide for mTLS Microservices

This guide shows you how to deploy the mTLS microservices to various Platform-as-a-Service (PaaS) providers like Railway, Render, Heroku, DigitalOcean App Platform, etc.

## 📋 Prerequisites

1. Generate certificates locally first:
   ```bash
   cd certificates && ./generate-certs.sh
   ```

2. Build the deployment-ready Docker images:
   ```bash
   docker build -t service1-deploy ./deploy/service1
   docker build -t service2-deploy ./deploy/service2
   ```

## 🔐 Required Secret Files

Upload these certificate files as **secret files** in your PaaS:

### For Service1:
- `ca.pem` → Mount at `/certs/ca.pem`
- `service1-cert.pem` → Mount at `/certs/service1-cert.pem`
- `service1-key.pem` → Mount at `/certs/service1-key.pem`
- `service1-client-cert.pem` → Mount at `/certs/service1-client-cert.pem`
- `service1-client-key.pem` → Mount at `/certs/service1-client-key.pem`

### For Service2:
- `ca.pem` → Mount at `/certs/ca.pem`
- `service2-cert.pem` → Mount at `/certs/service2-cert.pem`
- `service2-key.pem` → Mount at `/certs/service2-key.pem`
- `service2-client-cert.pem` → Mount at `/certs/service2-client-cert.pem`
- `service2-client-key.pem` → Mount at `/certs/service2-client-key.pem`

## ⚙️ Environment Variables

### Service1 Configuration:
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8443` | Port for the service |
| `TLS_CERT_PATH` | `/certs/service1-cert.pem` | Server certificate path |
| `TLS_KEY_PATH` | `/certs/service1-key.pem` | Server private key path |
| `CA_CERT_PATH` | `/certs/ca.pem` | CA certificate path |
| `CLIENT_CERT_PATH` | `/certs/service1-client-cert.pem` | Client certificate path |
| `CLIENT_KEY_PATH` | `/certs/service1-client-key.pem` | Client private key path |
| `SERVICE2_URL` | `https://service2:8444` | URL of Service2 (⚠️ **REQUIRED**) |

### Service2 Configuration:
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8444` | Port for the service |
| `TLS_CERT_PATH` | `/certs/service2-cert.pem` | Server certificate path |
| `TLS_KEY_PATH` | `/certs/service2-key.pem` | Server private key path |
| `CA_CERT_PATH` | `/certs/ca.pem` | CA certificate path |
| `CLIENT_CERT_PATH` | `/certs/service2-client-cert.pem` | Client certificate path |
| `CLIENT_KEY_PATH` | `/certs/service2-client-key.pem` | Client private key path |
| `SERVICE1_URL` | `https://service1:8443` | URL of Service1 (⚠️ **REQUIRED**) |

## 🚀 Platform-Specific Deployment

### Railway
```yaml
# railway.toml
[build]
builder = "DOCKER"
dockerfilePath = "deploy/service1/Dockerfile"

[deploy]
healthcheckPath = "/health"
healthcheckTimeout = 100
restartPolicyType = "ON_FAILURE"

# Environment Variables:
SERVICE2_URL=https://service2-production.up.railway.app
```

### Render
```yaml
# render.yaml
services:
  - type: web
    name: service1
    env: docker
    dockerfilePath: ./deploy/service1/Dockerfile
    envVars:
      - key: SERVICE2_URL
        value: https://service2.onrender.com
    disk:
      name: certs
      mountPath: /certs
      sizeGB: 1
```

### Heroku
```bash
# Using Heroku CLI
heroku create your-service1
heroku config:set SERVICE2_URL=https://your-service2.herokuapp.com
heroku container:push web
heroku container:release web
```

### DigitalOcean App Platform
```yaml
# .do/app.yaml
name: mtls-microservices
services:
- name: service1
  source_dir: /deploy/service1
  dockerfile_path: Dockerfile
  environment_slug: docker
  envs:
  - key: SERVICE2_URL
    value: https://service2-abc123.ondigitalocean.app
```

## 📁 How to Get Certificate Content

### Option 1: Base64 Encode (Most Compatible)
```bash
# For each certificate file:
base64 -i certificates/ca.pem
base64 -i certificates/service1-cert.pem
# ... etc
```

Then in your PaaS, create secret files with the base64 content.

### Option 2: Direct File Content
```bash
# Copy file content directly:
cat certificates/ca.pem
cat certificates/service1-cert.pem
# ... etc
```

Paste the content (including `-----BEGIN CERTIFICATE-----` headers) into your PaaS secret files.

### Option 3: Environment Variables (If File Mounting Not Supported)
Some PaaS providers don't support file mounting. In that case, modify the code to read from environment variables:

```bash
# Set these as environment variables with full certificate content:
CA_CERT_CONTENT="-----BEGIN CERTIFICATE-----
MIIFbzCCA1egAwIBAgIUJ...
-----END CERTIFICATE-----"

SERVICE1_CERT_CONTENT="-----BEGIN CERTIFICATE-----
MIIFdjCCA16gAwIBAgIUZ...
-----END CERTIFICATE-----"
```

## 🔧 Service URLs Configuration

**CRITICAL**: You must set the correct service URLs for inter-service communication.

### Pattern 1: Both services in same PaaS
```bash
# If your PaaS assigns URLs like:
# Service1: https://service1-abc123.yourpaas.com
# Service2: https://service2-def456.yourpaas.com

# Service1 environment:
SERVICE2_URL=https://service2-def456.yourpaas.com

# Service2 environment:
SERVICE1_URL=https://service1-abc123.yourpaas.com
```

### Pattern 2: Custom domains
```bash
# Service1 environment:
SERVICE2_URL=https://api2.yourdomain.com

# Service2 environment:
SERVICE1_URL=https://api1.yourdomain.com
```

## 🧪 Testing Your Deployment

Once deployed, test the mTLS communication:

```bash
# Test Service1 health
curl https://your-service1-url.com/health

# Test Service2 health
curl https://your-service2-url.com/health

# Test inter-service communication
curl https://your-service1-url.com/call-service2
curl https://your-service2-url.com/call-service1
```

## 🚨 Common Issues

### 1. Certificate path issues
**Error**: `failed to load server certificate`
**Solution**: Verify secret files are mounted at correct paths

### 2. Service communication fails
**Error**: `Failed to call service2: dial tcp: lookup service2`
**Solution**: Check `SERVICE1_URL`/`SERVICE2_URL` environment variables

### 3. Certificate validation fails
**Error**: `x509: certificate signed by unknown authority`
**Solution**: Ensure all services use the same `ca.pem` file

### 4. Port binding issues
**Error**: `bind: permission denied`
**Solution**: Most PaaS providers require port 80/443 or use `PORT` env var

## 📊 Monitoring & Logging

Add these environment variables for better observability:

```bash
# Optional monitoring variables
LOG_LEVEL=info
ENABLE_METRICS=true
HEALTH_CHECK_INTERVAL=30s
```

## 🔒 Security Best Practices

1. **Never commit certificates** to your repository
2. **Use secret management** features of your PaaS
3. **Rotate certificates** regularly (recommend 90 days)
4. **Monitor certificate expiry** dates
5. **Use least privilege** access controls
6. **Enable audit logging** if available

## 📞 Support

If you encounter issues:

1. Check service logs for certificate loading errors
2. Verify environment variables are set correctly
3. Test certificate validity locally first
4. Ensure inter-service URLs are accessible

Example debug commands:
```bash
# Check if certificates are loaded
ls -la /certs/

# Verify environment variables
env | grep -E "(SERVICE|TLS|CA|CLIENT)"

# Test certificate validity
openssl x509 -in /certs/service1-cert.pem -text -noout
```