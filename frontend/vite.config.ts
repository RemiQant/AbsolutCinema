import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc' // or just '@vitejs/plugin-react'
import { sentryVitePlugin } from "@sentry/vite-plugin";

export default defineConfig({
  plugins: [
    react(), 
    // Put the Sentry plugin last
    sentryVitePlugin({
      org: "remiqant",
      project: "absolutcinema-frontend",
      // The auth token needs to be set as a secret in your CI/CD (Digital Ocean)
      authToken: process.env.SENTRY_AUTH_TOKEN, 
      sourcemaps: {
          // Delete source maps after upload to Sentry (recommended)
          filesToDeleteAfterUpload: ["./dist/**/*.map"],
      },
    }),
  ],
  build: {
    sourcemap: true, // MUST be set to true
  },
})