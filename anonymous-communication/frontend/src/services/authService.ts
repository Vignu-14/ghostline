import type { LoginRequest, RegisterRequest, User } from "../types/user";
import { apiRequest } from "./api";

export function login(input: LoginRequest) {
  return apiRequest<{ user: User }>("/api/auth/login", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function register(input: RegisterRequest) {
  return apiRequest<{ user: User }>("/api/auth/register", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function logout() {
  return apiRequest<null>("/api/auth/logout", {
    method: "POST",
  });
}

export function getCurrentUser() {
  return apiRequest<{ user: User }>("/api/auth/me");
}
