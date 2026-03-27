import { APIError, type APIResponse } from "../types/api";
import { API_BASE_URL } from "../utils/constants";

type RequestOptions = RequestInit & {
  isFormData?: boolean;
};

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

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
    credentials: "include",
  });

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
