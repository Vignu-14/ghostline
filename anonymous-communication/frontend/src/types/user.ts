export interface User {
  id: string;
  username: string;
  email: string;
  role: "user" | "admin";
  profile_picture_url?: string | null;
  created_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface UserSearchResult {
  id: string;
  username: string;
  profile_picture_url?: string | null;
}

export interface PublicProfile {
  id: string;
  username: string;
  profile_picture_url?: string | null;
  created_at: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}
