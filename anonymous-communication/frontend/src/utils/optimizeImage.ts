const MAX_DIMENSION = 1600;
const TARGET_SIZE_BYTES = 900 * 1024;
const JPEG_QUALITY = 0.84;
const WEBP_QUALITY = 0.82;

export async function optimizeImage(file: File) {
  if (!file.type.startsWith("image/") || file.type === "image/gif") {
    return {
      file,
      didOptimize: false,
    };
  }

  const imageBitmap = await createImageBitmap(file);
  try {
    const scale = Math.min(1, MAX_DIMENSION / Math.max(imageBitmap.width, imageBitmap.height));
    const width = Math.max(1, Math.round(imageBitmap.width * scale));
    const height = Math.max(1, Math.round(imageBitmap.height * scale));

    const canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;

    const context = canvas.getContext("2d");
    if (!context) {
      return {
        file,
        didOptimize: false,
      };
    }

    context.drawImage(imageBitmap, 0, 0, width, height);

    const preferredType =
      file.type === "image/png" && file.size < TARGET_SIZE_BYTES ? "image/png" : "image/webp";
    const quality = preferredType === "image/png" ? undefined : preferredType === "image/webp" ? WEBP_QUALITY : JPEG_QUALITY;

    const optimizedBlob = await canvasToBlob(canvas, preferredType, quality);
    if (!optimizedBlob || optimizedBlob.size >= file.size) {
      return {
        file,
        didOptimize: false,
      };
    }

    const extension = preferredType === "image/png" ? "png" : preferredType === "image/webp" ? "webp" : "jpg";
    const optimizedFile = new File([optimizedBlob], replaceExtension(file.name, extension), {
      type: preferredType,
      lastModified: file.lastModified,
    });

    return {
      file: optimizedFile,
      didOptimize: true,
    };
  } finally {
    imageBitmap.close();
  }
}

function canvasToBlob(canvas: HTMLCanvasElement, type: string, quality?: number) {
  return new Promise<Blob | null>((resolve) => {
    canvas.toBlob((blob) => resolve(blob), type, quality);
  });
}

function replaceExtension(fileName: string, nextExtension: string) {
  const trimmedName = fileName.trim();
  if (!trimmedName.includes(".")) {
    return `${trimmedName || "upload"}.${nextExtension}`;
  }

  return `${trimmedName.replace(/\.[^.]+$/, "")}.${nextExtension}`;
}
