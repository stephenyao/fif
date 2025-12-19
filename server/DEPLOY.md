# FIF Server - Digital Ocean Deployment Guide

## Quick Deploy to Digital Ocean App Platform

### 1. Prerequisites
- GitHub/GitLab account with your code pushed
- Digital Ocean account

### 2. Push Your Code
```bash
git add .
git commit -m "Add deployment configuration"
git push origin main
```

### 3. Deploy on Digital Ocean

1. **Go to [Digital Ocean App Platform](https://cloud.digitalocean.com/apps)**

2. **Click "Create App"**

3. **Connect Your Repository:**
   - Select your GitHub/GitLab account
   - Choose the `fif/server` repository
   - Select the `main` branch
   - Click "Next"

4. **Configure App:**
   - App name: `fif-server` (or your choice)
   - Region: Choose closest to your users
   - Autodeploy: ✅ Enable (deploys on git push)

5. **Set Environment Variables:**
   Click "Edit" next to environment variables and add:
   ```
   FIREBASE_KEY_B64=<your-base64-firebase-key>
   ALLOWED_ORIGINS=https://your-frontend-domain.com
   ```

6. **Configure Resources:**
   - Instance type: Basic ($12/month recommended)
   - Instance size: 1GB RAM / 1 vCPU

7. **Review & Deploy:**
   - Review settings
   - Click "Create Resources"
   - Wait 5-10 minutes for deployment

### 4. After Deployment

Your app will be available at:
```
https://fif-server-xxxxx.ondigitalocean.app
```

Update your frontend's API URL to point to this domain.

### 5. Custom Domain (Optional)

1. Go to your app settings
2. Click "Domains"
3. Add your custom domain
4. Update DNS records as instructed
5. SSL certificate auto-provisioned

## Local Testing with Docker

Test your Docker build locally before deploying:

```bash
# Build the image
docker build -t fif-server .

# Run locally
docker run -p 8080:8080 \
  -e FIREBASE_KEY_B64="your-key" \
  -e ALLOWED_ORIGINS="http://localhost:5173" \
  fif-server
```

Visit: http://localhost:8080/health

## Environment Variables Reference

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `FIREBASE_KEY_B64` | ✅ Yes | Base64-encoded Firebase credentials | `ewogICJ0eXBlIjogI...` |
| `ALLOWED_ORIGINS` | ✅ Yes | Comma-separated CORS origins | `https://app.com,https://www.app.com` |
| `PORT` | No | Server port (auto-set by DO) | `8080` |

## Monitoring

Digital Ocean App Platform provides:
- ✅ Build logs
- ✅ Runtime logs
- ✅ Metrics (CPU, memory, requests)
- ✅ Alerts

Access via: App → Logs/Insights tabs

## Costs

- **Basic Plan:** $12/month
  - 1 GB RAM
  - 1 vCPU
  - Perfect for getting started

## Troubleshooting

**Build fails:**
- Check build logs in Digital Ocean console
- Verify `go.mod` and `go.sum` are committed

**App crashes:**
- Check runtime logs
- Verify environment variables are set correctly
- Test Docker build locally first

**CORS errors:**
- Ensure `ALLOWED_ORIGINS` includes your frontend domain
- Check browser console for exact origin being sent

## Next Steps

1. ✅ Deploy app
2. Test `/health` endpoint
3. Update frontend API URL
4. Add custom domain (optional)
5. Set up monitoring alerts
6. Add database when needed (Supabase recommended)
