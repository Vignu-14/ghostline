import { APIError } from "../types/api";
import type { Post, PostUploadTarget } from "../types/post";
import { apiRequest } from "./api";

export function listPosts(page = 1, limit = 20) {
  return apiRequest<{ posts: Post[] | null; page: number; limit: number }>(`/api/posts?page=${page}&limit=${limit}`).then(
    (response) => ({
      ...response,
      posts: Array.isArray(response.posts) ? response.posts : [],
    }),
  );
}

export function likePost(postID: string) {
  return apiRequest<null>(`/api/posts/${postID}/like`, {
    method: "POST",
  });
}

export function unlikePost(postID: string) {
  return apiRequest<null>(`/api/posts/${postID}/like`, {
    method: "DELETE",
  });
}

export function deletePost(postID: string) {
  return apiRequest<null>(`/api/posts/${postID}`, {
    method: "DELETE",
  });
}

export function createPost(caption: string, image?: File | null) {
  return createPostWithDirectUpload(caption, image);
}

export function createUploadURL(image: File) {
  return apiRequest<{ upload: PostUploadTarget }>("/api/posts/upload-url", {
    method: "POST",
    body: JSON.stringify({
      file_name: image.name,
      content_type: image.type,
      file_size: image.size,
    }),
  });
}

export function finalizePost(caption: string, objectPath: string) {
  return apiRequest<{ post: Post }>("/api/posts/finalize", {
    method: "POST",
    body: JSON.stringify({
      caption,
      object_path: objectPath,
    }),
  });
}

export async function uploadImageToSupabase(upload: PostUploadTarget, image: File) {
  let response: Response;

  try {
    response = await fetch(upload.upload_url, {
      method: upload.method || "PUT",
      headers: {
        "Content-Type": image.type || "application/octet-stream",
      },
      body: image,
    });
  } catch {
    throw new APIError("Unable to upload image directly to Supabase.", 0);
  }

  if (!response.ok) {
    throw new APIError("Unable to upload image directly to Supabase.", response.status);
  }
}

async function createPostWithDirectUpload(caption: string, image?: File | null) {
  if (!image) {
    return finalizePost(caption, "");
  }

  const { upload } = await createUploadURL(image);
  await uploadImageToSupabase(upload, image);
  return finalizePost(caption, upload.object_path);
}
