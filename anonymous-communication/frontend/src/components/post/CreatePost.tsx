import type { Post } from "../../types/post";
import { useState, type FormEvent } from "react";
import * as postService from "../../services/postService";
import { getErrorMessage } from "../../utils/errorHandler";
import { optimizeImage } from "../../utils/optimizeImage";
import { validateImage } from "../../utils/validateImage";
import { Button } from "../common/Button";

type CreatePostProps = {
  onCreated: (post: Post) => Promise<void> | void;
};

export function CreatePost({ onCreated }: CreatePostProps) {
  const [caption, setCaption] = useState("");
  const [image, setImage] = useState<File | null>(null);
  const [error, setError] = useState("");
  const [status, setStatus] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!image && !caption.trim()) {
      setError("Add a photo or write a thought before posting.");
      return;
    }

    let uploadFile = image;
    if (uploadFile) {
      const validationError = validateImage(uploadFile);
      if (validationError) {
        setError(validationError);
        return;
      }
    }

    setIsSubmitting(true);
    setError("");
    setStatus("");

    try {
      if (uploadFile) {
        setStatus("Optimizing image...");
        const optimized = await optimizeImage(uploadFile);
        uploadFile = optimized.file;
        if (optimized.didOptimize) {
          const savedKB = Math.max(0, Math.round((image!.size - uploadFile.size) / 1024));
          setStatus(savedKB > 0 ? `Image optimized. Saved about ${savedKB} KB.` : "Image optimized.");
        } else {
          setStatus("Uploading to Supabase...");
        }
      } else {
        setStatus("Publishing thought...");
      }

      const response = await postService.createPost(caption, uploadFile);
      setCaption("");
      setImage(null);
      setStatus("");
      await onCreated(response.post);
    } catch (submitError) {
      setError(getErrorMessage(submitError, "Unable to create post."));
    } finally {
      setStatus("");
      setIsSubmitting(false);
    }
  }

  return (
    <form className="composer" onSubmit={handleSubmit}>
      <div className="composer__lead">
        <p className="eyebrow">Share a private drop</p>
        <h2>Post a photo or a thought.</h2>
      </div>

      <textarea
        className="composer__caption"
        onChange={(event) => setCaption(event.target.value)}
        placeholder="Write a short caption..."
        rows={3}
        value={caption}
      />

      <label className="composer__upload">
        <span>{image ? image.name : "Choose image (optional)"}</span>
        <input
          accept="image/jpeg,image/png,image/gif,image/webp"
          onChange={(event) => setImage(event.target.files?.[0] || null)}
          type="file"
        />
      </label>

      {error ? <p className="form-error">{error}</p> : null}
      {status ? <p className="support-copy">{status}</p> : null}

      <Button disabled={isSubmitting} type="submit">
        {isSubmitting ? "Publishing..." : "Publish to Ghostline"}
      </Button>
    </form>
  );
}
