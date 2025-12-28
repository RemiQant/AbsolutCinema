# Deploy Frontend to DigitalOcean App Platform

This guide walks you through deploying the React frontend to DigitalOcean App Platform alongside your existing backend.

## Architecture Overview

After deployment, you'll have:
- **Frontend:** `https://absolut-cinema-umwih.ondigitalocean.app` (or custom domain)
- **Backend:** `https://absolut-cinema-umwih.ondigitalocean.app/api` (via same app) OR separate URL
- **Benefits:** Same registrable domain = simpler cookie handling, no CORS issues with SameSite=Lax

## Prerequisites

1. DigitalOcean account with existing backend app
2. GitHub repository connected to DigitalOcean
3. Changes committed and pushed to your repository

## Deployment Steps

### Step 1: Add Frontend Component to Your App

1. **Go to DigitalOcean Dashboard**
   - Navigate to your existing app: `absolut-cinema-umwih`
   - Click **"Components"** tab
   - Click **"Create Component"**

2. **Select Component Type**
   - Choose **"Static Site"**
   - Source: Your GitHub repository
   - Branch: `main` (or your production branch)

3. **Configure Build Settings**
   ```
   Build Command: cd frontend && npm install && npm run build
   Output Directory: frontend/dist
   ```

4. **Environment Variables** (if needed)
   - Add any frontend-specific env vars
   - Vite env vars must be prefixed with `VITE_`
   - Example:
     ```
     VITE_API_URL=https://absolut-cinema-umwih.ondigitalocean.app/api
     ```

5. **Resource Settings**
   - Instance Type: Basic (usually sufficient for static sites)
   - Instance Size: 512MB or 1GB

6. **Click "Create Component"**

### Step 2: Configure Routing (Important!)

Since your frontend is a SPA (Single Page Application), you need to handle client-side routing:

**Option A: Via DigitalOcean Console**
1. After component is created, go to **Settings**
2. Under **Routes**, ensure all paths (`/*`) serve `index.html`
3. DigitalOcean usually auto-configures this for SPAs

**Option B: Add `_redirects` file (Recommended)**
Create this file in your frontend's `public/` folder:

```
/*  /index.html  200
```

This tells the server to serve `index.html` for all routes (required for React Router).

### Step 3: Update Backend CORS (Already Done ✅)

The backend CORS has been updated to allow:
- `http://localhost:5173` (local dev)
- `https://absolut-cinema.vercel.app` (legacy)
- `https://absolut-cinema-umwih.ondigitalocean.app` (new frontend)

### Step 4: Domain Configuration Options

**Option A: Use Default DigitalOcean Domain**
- Frontend: `https://absolut-cinema-umwih.ondigitalocean.app`
- Backend: `https://absolut-cinema-umwih.ondigitalocean.app/api`
- Cookies work perfectly (same domain!)

**Option B: Use Custom Domain**
If you have a custom domain (e.g., `absolutcinema.com`):
1. Go to App Settings → Domains
2. Add custom domain
3. Update DNS records as instructed
4. Both frontend and backend will use: `https://absolutcinema.com`

### Step 5: Deploy and Test

1. **Commit and push changes:**
   ```bash
   git add .
   git commit -m "feat: add DigitalOcean frontend deployment config"
   git push origin main
   ```

2. **Trigger Deployment:**
   - DigitalOcean will auto-deploy on push
   - Monitor deployment in the console
   - Check build logs if any errors occur

3. **Test the Application:**
   - Visit your frontend URL
   - Register/Login
   - **Refresh the page** → You should stay logged in! ✅
   - Test on different browsers (Chrome, Brave, Safari)

### Step 6: Environment-Specific Behavior

Your app now handles environments automatically:

| Environment | Frontend | Backend | Behavior |
|-------------|----------|---------|----------|
| **Local Dev** | `localhost:5173` | `localhost:8080` | SameSite=Lax, Secure=false |
| **Production** | DigitalOcean | DigitalOcean | SameSite=Lax or None, Secure=true |

## Troubleshooting

### Build Fails
- Check `npm install` works locally in `frontend/` folder
- Ensure `package.json` has correct build script
- Check DigitalOcean build logs for specific errors

### Blank Page After Deploy
- Verify Output Directory is set to `frontend/dist`
- Check browser console for errors
- Ensure `_redirects` file exists in `public/` folder

### Cookies Still Not Working
- Verify backend has `APP_ENV=production`
- Check browser Network tab → Response headers should have `Set-Cookie`
- Ensure both frontend and backend use HTTPS

### CORS Errors
- Backend CORS should allow your frontend domain
- Check [backend/internal/server/routes.go](backend/internal/server/routes.go) has correct origins
- Redeploy backend if you changed CORS config

## Next Steps

1. **Remove Vercel Deployment** (optional)
   - Once DigitalOcean frontend is stable
   - Cancel Vercel subscription if not needed

2. **Set Up CI/CD**
   - DigitalOcean auto-deploys on git push
   - Configure deployment branches in app settings

3. **Monitor Performance**
   - Use DigitalOcean monitoring dashboard
   - Set up alerts for downtime

4. **Custom Domain** (optional but recommended)
   - Purchase domain via DigitalOcean or external registrar
   - Point to your DigitalOcean app
   - Enables cleaner URLs and better branding

## Cost Estimate

**DigitalOcean App Platform Pricing:**
- Static Site (Frontend): ~$5/month (512MB instance)
- Backend (Already running): ~$12-25/month depending on plan
- Total: ~$17-30/month for full stack

Compare to:
- Vercel Free Tier: Limited builds/bandwidth
- Current setup: Vercel + DO = potential costs on both platforms

## Support

If you encounter issues:
1. Check DigitalOcean build logs
2. Review browser console errors
3. Test locally first: `npm run build && npm run preview`
4. DigitalOcean Support: https://www.digitalocean.com/support

---

**Ready to deploy?** Follow Step 1 above! 🚀
