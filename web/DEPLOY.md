# FIF Frontend - Digital Ocean Deployment Guide

## Quick Deploy to Digital Ocean App Platform (Static Site)

### 1. Prerequisites
- GitHub/GitLab account with your code pushed
- Digital Ocean account
- Backend already deployed (need the URL)

### 2. Prepare Environment Variables
Make sure your API calls use `import.meta.env.VITE_API_URL` (this will be handled by the update).

### 3. Deploy on Digital Ocean

1. **Go to [Digital Ocean App Platform](https://cloud.digitalocean.com/apps)**

2. **Click "Create App"**

3. **Connect Your Repository:**
   - Select your GitHub/GitLab account
   - Choose the `fif` repository
   - Select the `main` branch
   - **Source Directory:** Set this to `web`
   - Click "Next"

4. **Configure Resources:**
   - Digital Ocean should detect it as a **Static Site**.
   - If not, ensure the build command is `npm run build` and the output directory is `dist`.
   - Instance type: Free (Starter) or Basic.

5. **Set Environment Variables:**
   Click "Edit" next to environment variables and add:
   ```
   VITE_API_URL=https://your-backend-domain.ondigitalocean.app
   ```
   > [!IMPORTANT]
   > Ensure there is no trailing slash in the URL.

6. **Review & Deploy:**
   - Review settings
   - Click "Create Resources"

### 4. Update Backend CORS
Once your frontend is deployed, you must update the backend's `ALLOWED_ORIGINS` environment variable to include the new frontend URL.

1. Go to your **Backend App** in Digital Ocean.
2. Settings -> Environment Variables.
3. Update `ALLOWED_ORIGINS` (e.g., `https://fif-web-xxxxx.ondigitalocean.app`).
4. Redeploy the backend.

### 5. Troubleshooting
- **404 on Refresh**: If you use React Router and get a 404 when refreshing a subpage (like `/dashboard`), ensure you have a "Catch-all" or "Error Page" configuration in DO pointing to `index.html`. 
- **CORS Errors**: Check the backend logs and ensure `ALLOWED_ORIGINS` exactly matches the frontend origin.
