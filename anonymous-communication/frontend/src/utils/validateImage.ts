import { MAX_IMAGE_SIZE_MB } from "./constants";

const ALLOWED_IMAGE_TYPES = new Set(["image/jpeg", "image/png", "image/gif", "image/webp"]);

export function validateImage(file: File) {
  if (!ALLOWED_IMAGE_TYPES.has(file.type)) {
    return "Please choose a JPEG, PNG, GIF, or WebP image.";
  }

  if (file.size > MAX_IMAGE_SIZE_MB * 1024 * 1024) {
    return `Image must be smaller than ${MAX_IMAGE_SIZE_MB}MB.`;
  }

  return "";
}
