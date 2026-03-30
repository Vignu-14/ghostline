export const APP_NAME = "Ghostline";

// API Base URL configuration with validation
const rawApiUrl = import.meta.env.VITE_API_BASE_URL || "";
export const API_BASE_URL = rawApiUrl.replace(/\/$/, "");

// Validate API URL in production
if (import.meta.env.PROD && !API_BASE_URL) {
  console.error(
    "CRITICAL: VITE_API_BASE_URL is not set! " +
    "Authentication requests will fail. " +
    "Please set this environment variable in your deployment platform (Vercel)."
  );
}

// Log API configuration in development
if (import.meta.env.DEV) {
  console.log("API Configuration:", {
    baseUrl: API_BASE_URL || "(using Vite proxy)",
    mode: import.meta.env.MODE,
  });
}

export const MAX_IMAGE_SIZE_MB = 5;
