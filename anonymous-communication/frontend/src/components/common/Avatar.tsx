type AvatarProps = {
  src?: string | null;
  alt: string;
  size?: "sm" | "md" | "lg";
};

export function Avatar({ src, alt, size = "md" }: AvatarProps) {
  return (
    <span className={`avatar avatar--${size}`}>
      {src ? <img src={src} alt={alt} /> : <span>{alt.slice(0, 1).toUpperCase()}</span>}
    </span>
  );
}
