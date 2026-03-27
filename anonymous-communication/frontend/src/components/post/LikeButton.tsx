import { Button } from "../common/Button";

type LikeButtonProps = {
  liked: boolean;
  count: number;
  onToggle: () => void;
  disabled?: boolean;
};

export function LikeButton({ liked, count, onToggle, disabled }: LikeButtonProps) {
  return (
    <Button className="like-button" disabled={disabled} variant="quiet" onClick={onToggle} type="button">
      <span aria-hidden="true">{liked ? "♥" : "♡"}</span>
      <span>{count}</span>
    </Button>
  );
}
