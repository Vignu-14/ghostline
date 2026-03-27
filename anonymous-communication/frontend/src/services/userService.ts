import type { Post } from "../types/post";
import type { PublicProfile, UserSearchResult } from "../types/user";
import { apiRequest } from "./api";

export function searchUsers(query: string, limit = 8) {
  const params = new URLSearchParams({
    q: query,
    limit: String(limit),
  });

  return apiRequest<{ users: UserSearchResult[] | null; query: string; limit: number }>(
    `/api/users/search?${params.toString()}`,
  ).then((response) => ({
    ...response,
    users: Array.isArray(response.users) ? response.users : [],
  }));
}

export function getProfileByUsername(username: string, page = 1, limit = 20) {
  return apiRequest<{ profile: PublicProfile; posts: Post[] | null; page: number; limit: number }>(
    `/api/users/profile/${encodeURIComponent(username)}?page=${page}&limit=${limit}`,
  ).then((response) => ({
    ...response,
    posts: Array.isArray(response.posts) ? response.posts : [],
  }));
}
