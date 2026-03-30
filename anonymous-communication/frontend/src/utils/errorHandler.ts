import { APIError } from "../types/api";

/**
 * Extracts a user-friendly error message from various error types.
 * Provides specific guidance for common authentication failures.
 */
export function getErrorMessage(error: unknown, fallback = "Something went wrong."): string {
  if (error instanceof APIError) {
    // Network errors (status 0) from failed fetch requests
    if (error.status === 0) {
      return `${error.message} If this persists, please contact support.`;
    }

    // Rate limiting
    if (error.status === 429) {
      return `${error.message} Please wait a few minutes before trying again.`;
    }

    // Authentication errors
    if (error.status === 401) {
      return "Invalid credentials. Please check your username and password.";
    }

    // Conflict errors (duplicate username/email)
    if (error.status === 409) {
      return error.message;
    }

    return error.message;
  }

  if (error instanceof Error) {
    // Handle generic JavaScript errors
    if (error.message.includes("fetch") || error.message.includes("network")) {
      return "Network error. Please check your internet connection and try again.";
    }
    return error.message;
  }

  return fallback;
}

/**
 * Determines if an error is a network/connection error
 */
export function isNetworkError(error: unknown): boolean {
  if (error instanceof APIError) {
    return error.status === 0;
  }
  return error instanceof Error && (
    error.message.includes("fetch") ||
    error.message.includes("network") ||
    error.message.includes("Failed to fetch")
  );
}
