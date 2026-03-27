import { Link } from "react-router-dom";
import type { Post } from "../../types/post";
import { formatRelativeDate } from "../../utils/formatDate";
import { Avatar } from "../common/Avatar";
import { Button } from "../common/Button";
import { LikeButton } from "./LikeButton";

type PostCardProps = {
  canDelete?: boolean;
  onDelete?: (post: Post) => void;
  post: Post;
  onLikeToggle?: (post: Post) => void;
};

export function PostCard({ canDelete, onDelete, post, onLikeToggle }: PostCardProps) {
  return (
    <article className="post-card">
      <header className="post-card__header">
        <Link className="post-card__identity post-card__identity--link" to={`/u/${post.user.username}`}>
          <Avatar alt={post.user.username} src={post.user.profile_picture_url} />
          <div>
            <strong>@{post.user.username}</strong>
            <p>{formatRelativeDate(post.created_at)}</p>
          </div>
        </Link>
        <div className="post-card__actions">
          <LikeButton
            count={post.like_count}
            liked={Boolean(post.liked_by_viewer)}
            onToggle={() => onLikeToggle?.(post)}
          />
          {canDelete ? (
            <Button onClick={() => onDelete?.(post)} type="button" variant="quiet">
              Delete
            </Button>
          ) : null}
        </div>
      </header>

      {post.image_url ? (
        <img
          alt={post.caption || `${post.user.username}'s post`}
          className="post-card__image"
          src={post.image_url}
        />
      ) : null}

      {post.caption ? <p className="post-card__caption">{post.caption}</p> : null}
    </article>
  );
}
