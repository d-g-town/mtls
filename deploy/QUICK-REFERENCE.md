# 🚀 Quick PaaS Deployment Reference

## Step 1: Extract Certificates

```bash
cd deploy
chmod +x extract-secrets.sh
./extract-secrets.sh
```

## Step 2: Upload Secret Files

### For Service1:

Mount these files in your PaaS:

- `secrets/service1-ca.pem` → `/certs/ca.pem`
- `secrets/service1-cert.pem` → `/certs/service1-cert.pem`
- `secrets/service1-key.pem` → `/certs/service1-key.pem`
- `secrets/service1-client-cert.pem` → `/certs/service1-client-cert.pem`
- `secrets/service1-client-key.pem` → `/certs/service1-client-key.pem`

### For Service2:

Mount these files in your PaaS:

- `secrets/service2-ca.pem` → `/certs/ca.pem`
- `secrets/service2-cert.pem` → `/certs/service2-cert.pem`
- `secrets/service2-key.pem` → `/certs/service2-key.pem`
- `secrets/service2-client-cert.pem` → `/certs/service2-client-cert.pem`
- `secrets/service2-client-key.pem` → `/certs/service2-client-key.pem`

## Step 3: Set Environment Variables

### Service1 Environment Variables:

```
SERVICE2_URL=https://your-service2-url.paas.com
```

### Service2 Environment Variables:

```
SERVICE1_URL=https://your-service1-url.paas.com
```

## Step 4: Deploy

### Railway:

```bash
railway login
railway link
railway up
```

### Render:

```bash
# Push to GitHub, connect repo in Render dashboard
# Use Dockerfile: deploy/service1/Dockerfile
```

### Heroku:

```bash
heroku container:push web --app your-service1
heroku container:release web --app your-service1
```

## Step 5: Test

```bash
curl https://your-service1-url/health
curl https://your-service1-url/call-service2
```

## 🔧 Common PaaS URL Patterns

| Platform     | URL Pattern                                         |
| ------------ | --------------------------------------------------- |
| Railway      | `https://service1-production-abc123.up.railway.app` |
| Render       | `https://service1-abc123.onrender.com`              |
| Heroku       | `https://your-service1.herokuapp.com`               |
| DigitalOcean | `https://service1-abc123.ondigitalocean.app`        |

**Replace the URLs above with your actual service URLs!**
