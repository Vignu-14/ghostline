import { APIError, type APIResponse } from "../types/api";
import { API_BASE_URL } from "../utils/constants";

type RequestOptions = RequestInit & {
  isFormData?: boolean;
};

// Helper to check if we're in production
const isProduction = import.meta.env.PROD;

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers = new Headers(options.headers);

  if (!options.isFormData) {
    headers.set("Content-Type", "application/json");
  }

  // Add token to Authorization header if available
  const token = localStorage.getItem('auth_token');
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  // Generate request ID for tracing
  const requestId = crypto.randomUUID();
  headers.set("X-Request-ID", requestId);

  const url = `${API_BASE_URL}${path}`;

  let response: Response;

  try {
    response = await fetch(url, {
      ...options,
      headers,
      credentials: "include",
    });
  } catch (networkError) {
    // Handle network errors (CORS, connection refused, DNS failures, etc.)
    const errorMessage = isProduction
      ? "Unable to connect to server. Please check your internet connection and try again."
      : `Network error: ${networkError instanceof Error ? networkError.message : 'Unknown error'}`;

    console.error("Network request failed:", {
      url,
      method: options.method || 'GET',
      requestId,
      error: networkError,
      apiBaseUrl: API_BASE_URL,
    });

    throw new APIError(errorMessage, 0);
  }

  let payload: APIResponse<T>;

  try {
    payload = (await response.json()) as APIResponse<T>;
  } catch {
    throw new APIError("Unable to read server response.", response.status);
  }

  if (!response.ok || payload.status === "error") {
    throw new APIError(payload.error || "Request failed.", response.status, payload.details);
  }

  return payload.data as T;
}
