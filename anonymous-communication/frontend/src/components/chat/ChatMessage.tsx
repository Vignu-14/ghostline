import type { Message } from "../../types/message";
import { formatChatTime } from "../../utils/formatDate";

type ChatMessageProps = {
  isOwn: boolean;
  isSelected: boolean;
  isSelectionMode: boolean;
  message: Message;
  onToggleSelect: (message: Message) => void;
};

export function ChatMessage({
  isOwn,
  isSelected,
  isSelectionMode,
  message,
  onToggleSelect,
}: ChatMessageProps) {
  const bubbleClassName = [
    "chat-message",
    isOwn ? "chat-message--sent" : "chat-message--received",
    message.deleted_for_everyone ? "chat-message--deleted" : "",
    isSelectionMode ? "chat-message--selectable" : "",
    isSelected ? "is-selected" : "",
  ]
    .filter(Boolean)
    .join(" ");

  const bubbleContent = (
    <p className={message.deleted_for_everyone ? "chat-message__deleted-copy" : ""}>{message.content}</p>
  );

  return (
    <div className={`chat-message-row ${isOwn ? "chat-message-row--sent" : "chat-message-row--received"}`}>
      {isOwn ? <span className="chat-message__time">{formatChatTime(message.created_at)}</span> : null}
      {isSelectionMode ? (
        <button
          aria-pressed={isSelected}
          className={`${bubbleClassName} chat-message__button`}
          onClick={() => onToggleSelect(message)}
          type="button"
        >
          {bubbleContent}
        </button>
      ) : (
        <div className={bubbleClassName}>{bubbleContent}</div>
      )}
      {!isOwn ? <span className="chat-message__time">{formatChatTime(message.created_at)}</span> : null}
    </div>
  );
}
