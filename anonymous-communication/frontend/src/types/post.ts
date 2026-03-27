export interface PostUser {
  id: string;
  username: string;
  profile_picture_url?: string | null;
}

export interface Post {
  id: string;
  user_id?: string;
  image_url?: string | null;
  caption?: string | null;
  created_at: string;
  like_count: number;
  liked_by_viewer?: boolean;
  user: PostUser;
}

export interface PostUploadTarget {
  object_path: string;
  upload_url: string;
  public_url: string;
  method: string;
}
