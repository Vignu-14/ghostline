import { useEffect, useState } from "react";
import { useAuth } from "../../hooks/useAuth";
import * as postService from "../../services/postService";
import type { Post } from "../../types/post";
import { getErrorMessage } from "../../utils/errorHandler";
import { Spinner } from "../common/Spinner";
import { CreatePost } from "./CreatePost";
import { PostCard } from "./PostCard";

export function PostFeed() {
  const { user } = useAuth();
  const [posts, setPosts] = useState<Post[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  async function loadPosts() {
    try {
      setError("");
      const response = await postService.listPosts();
      setPosts(response.posts || []);
    } catch (loadError) {
      setError(getErrorMessage(loadError, "Unable to load feed."));
    } finally {
      setIsLoading(false);
    }
  }

  function handleCreated(post: Post) {
    setPosts((current) => [post, ...current.filter((item) => item.id !== post.id)]);
  }

  async function handleDelete(post: Post) {
    const confirmed = window.confirm("Delete this post from Ghostline?");
    if (!confirmed) {
      return;
    }

    const previousPosts = posts;
    setPosts((current) => current.filter((item) => item.id !== post.id));

    try {
      await postService.deletePost(post.id);
    } catch (deleteError) {
      setPosts(previousPosts);
      setError(getErrorMessage(deleteError, "Unable to delete post."));
    }
  }

  async function handleToggleLike(post: Post) {
    const liked = Boolean(post.liked_by_viewer);

    setPosts((current) =>
      current.map((item) =>
        item.id === post.id
          ? {
              ...item,
              liked_by_viewer: !liked,
              like_count: item.like_count + (liked ? -1 : 1),
            }
          : item,
      ),
    );

    try {
      if (liked) {
        await postService.unlikePost(post.id);
      } else {
        await postService.likePost(post.id);
      }
    } catch {
      setPosts((current) => current.map((item) => (item.id === post.id ? post : item)));
    }
  }

  useEffect(() => {
    void loadPosts();
  }, []);

  return (
    <section className="feed">
      {user ? <CreatePost onCreated={handleCreated} /> : null}

      {isLoading ? (
        <div className="center-stage">
          <Spinner />
        </div>
      ) : null}

      {error ? <p className="form-error">{error}</p> : null}

      {!isLoading && posts.length === 0 ? (
        <section className="panel feed__empty">
          <p className="eyebrow">Feed is quiet</p>
          <h2>No drops yet.</h2>
          <p>
            {user
              ? "Share the first photo or thought and it will appear here immediately."
              : "Once people start posting photos or thoughts, the public feed will fill in here."}
          </p>
        </section>
      ) : (
        <div className="feed__stack">
          {posts.map((post) => (
            <PostCard
              canDelete={user?.id === post.user.id}
              key={post.id}
              onDelete={handleDelete}
              onLikeToggle={handleToggleLike}
              post={post}
            />
          ))}
        </div>
      )}
    </section>
  );
}
