# FIF - Digital Ocean Deployment Guide (Unified App)

This project is now configured as a **unified application**. The Go backend serves the React frontend static files, allowing you to deploy both as a single service on Digital Ocean App Platform.

## Quick Deploy to Digital Ocean App Platform

### 1. Prerequisites
- GitHub/GitLab account with your code pushed
- Digital Ocean account

### 2. Push Your Code
Ensure your changes (including the new root `Dockerfile`) are pushed to your repository.
```bash
git add .
git commit -m "Configure unified deployment"
git push origin main
```

### 3. Deploy on Digital Ocean

1. **Go to [Digital Ocean App Platform](https://cloud.digitalocean.com/apps)**

2. **Click "Create App"**

3. **Connect Your Repository:**
   - Select your GitHub/GitLab account
   - Choose the `fif` repository
   - Select the `main` branch
   - **Source Directory:** Leave as root `/` (it will use the root Dockerfile)
   - Click "Next"

4. **Configure App:**
   - App name: `fif-app`
   - Region: Choose closest to your users
   - Autodeploy: ✅ Enable

5. **Set Environment Variables:**
   Click "Edit" next to environment variables and add:
   ```
   FIREBASE_KEY_B64=<your-base64-firebase-key>
   ALLOWED_ORIGINS=https://your-app-domain.ondigitalocean.app
   ```
   > [!NOTE]
   > Since the frontend is served from the same origin as the backend, you can set `ALLOWED_ORIGINS` to your app's domain.

6. **Configure Resources:**
   - Instance type: Basic ($12/month recommended)
   - Instance size: 1GB RAM / 1 vCPU

7. **Review & Deploy:**
   - Click "Create Resources"
   - Wait 5-10 minutes for the build to complete.

### 4. Local Testing with Docker

To test the unified build locally:

```bash
# From the root directory
docker build -t fif-app .

# Run locally
docker run -p 8080:8080 \
  -e FIREBASE_KEY_B64="your-key" \
  -e ALLOWED_ORIGINS="http://localhost:8080" \
  fif-app
```

Visit: `http://localhost:8080` (Frontend) and `http://localhost:8080/api/health` (Backend).

## Environment Variables Reference

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `FIREBASE_KEY_B64` | ✅ Yes | Base64-encoded Firebase credentials | `ewogICJ0eXBlIjogI...` |
| `ALLOWED_ORIGINS` | ✅ Yes | Comma-separated CORS origins | `https://your-app.com` |
| `PORT` | No | Server port (auto-set by DO) | `8080` |
