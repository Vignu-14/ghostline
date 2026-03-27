type UnreadBadgeProps = {
  count: number;
};

export function UnreadBadge({ count }: UnreadBadgeProps) {
  if (count < 1) {
    return null;
  }

  return <span className="badge">{count}</span>;
}
