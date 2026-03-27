import { useEffect, useRef, useState } from "react";
import { Button } from "../common/Button";
import type { DeleteMode, Message } from "../../types/message";
import { ChatMessage } from "./ChatMessage";
import { MessageInput } from "./MessageInput";

type DeleteDialogState =
  | {
      kind: "messages";
      description: string;
      messageIDs: string[];
      title: string;
    }
  | {
      kind: "clear";
      description: string;
      title: string;
    };

type ChatWindowProps = {
  conversationUserID?: string;
  currentUserID: string;
  disabled?: boolean;
  messages: Message[];
  onClearConversation: (userID: string, mode: DeleteMode) => Promise<void>;
  onDeleteMessages: (messageIDs: string[], mode: DeleteMode) => Promise<void>;
  title: string;
  onSend: (content: string) => Promise<void>;
};

export function ChatWindow({
  conversationUserID,
  currentUserID,
  disabled,
  messages,
  onClearConversation,
  onDeleteMessages,
  title,
  onSend,
}: ChatWindowProps) {
  const messageListRef = useRef<HTMLDivElement | null>(null);
  const [deleteDialog, setDeleteDialog] = useState<DeleteDialogState | null>(null);
  const [isSelectionMode, setIsSelectionMode] = useState(false);
  const [isWorking, setIsWorking] = useState(false);
  const [localError, setLocalError] = useState("");
  const [selectedMessageIDs, setSelectedMessageIDs] = useState<string[]>([]);

  useEffect(() => {
    if (!messageListRef.current) {
      return;
    }

    messageListRef.current.scrollTop = messageListRef.current.scrollHeight;
  }, [messages, disabled, title]);

  useEffect(() => {
    setSelectedMessageIDs((current) =>
      current.filter((messageID) => messages.some((message) => message.id === messageID)),
    );
  }, [messages]);

  useEffect(() => {
    setDeleteDialog(null);
    setIsSelectionMode(false);
    setSelectedMessageIDs([]);
    setLocalError("");
  }, [conversationUserID, disabled]);

  function toggleSelection(message: Message) {
    setSelectedMessageIDs((current) =>
      current.includes(message.id)
        ? current.filter((messageID) => messageID !== message.id)
        : [...current, message.id],
    );
  }

  function openSelectedDeleteDialog() {
    if (selectedMessageIDs.length === 0) {
      setLocalError("Choose at least one message first.");
      return;
    }

    setDeleteDialog({
      kind: "messages",
      description:
        "Delete for everyone will remove only your selected messages from both sides. Messages sent by the other person will only be removed for you.",
      messageIDs: selectedMessageIDs,
      title:
        selectedMessageIDs.length === 1
          ? "Delete the selected message?"
          : `Delete ${selectedMessageIDs.length} selected messages?`,
    });
  }

  function openClearDialog() {
    setDeleteDialog({
      kind: "clear",
      description:
        "Delete for me clears the whole chat only from your side. Delete for everyone clears the chat for you and removes only your own sent messages from the other person's side.",
      title: "Clear this chat?",
    });
  }

  async function handleDeleteConfirm(mode: DeleteMode) {
    if (!deleteDialog || deleteDialog.kind !== "messages") {
      return;
    }

    setIsWorking(true);
    setLocalError("");

    try {
      await onDeleteMessages(deleteDialog.messageIDs, mode);
      setDeleteDialog(null);
      setSelectedMessageIDs([]);
      setIsSelectionMode(false);
    } catch (error) {
      setLocalError(error instanceof Error ? error.message : "Unable to delete messages.");
    } finally {
      setIsWorking(false);
    }
  }

  async function handleClearConfirm(mode: DeleteMode) {
    if (!conversationUserID) {
      return;
    }

    setIsWorking(true);
    setLocalError("");

    try {
      await onClearConversation(conversationUserID, mode);
      setDeleteDialog(null);
      setSelectedMessageIDs([]);
      setIsSelectionMode(false);
    } catch (error) {
      setLocalError(error instanceof Error ? error.message : "Unable to clear the conversation.");
    } finally {
      setIsWorking(false);
    }
  }

  return (
    <section className="chat-window">
      <header className="chat-window__header chat-window__header--split">
        <div>
          <p className="eyebrow">Direct line</p>
          <h2>{title}</h2>
        </div>

        {!disabled ? (
          <div className="chat-window__toolbar">
            {isSelectionMode ? (
              <>
                <span className="support-copy">
                  {selectedMessageIDs.length > 0
                    ? `${selectedMessageIDs.length} selected`
                    : "Tap messages to select"}
                </span>
                <Button
                  disabled={selectedMessageIDs.length === 0 || isWorking}
                  onClick={openSelectedDeleteDialog}
                  type="button"
                  variant="ghost"
                >
                  Delete selected
                </Button>
                <Button
                  disabled={isWorking}
                  onClick={() => {
                    setDeleteDialog(null);
                    setIsSelectionMode(false);
                    setLocalError("");
                    setSelectedMessageIDs([]);
                  }}
                  type="button"
                  variant="quiet"
                >
                  Cancel
                </Button>
              </>
            ) : (
              <>
                <Button disabled={isWorking} onClick={() => setIsSelectionMode(true)} type="button" variant="quiet">
                  Select
                </Button>
                <Button disabled={isWorking} onClick={openClearDialog} type="button" variant="ghost">
                  Clear chat
                </Button>
              </>
            )}
          </div>
        ) : null}
      </header>

      {localError ? <p className="form-error">{localError}</p> : null}

      <div className="chat-window__messages" ref={messageListRef}>
        {disabled ? (
          <div className="chat-window__empty">
            <p className="eyebrow">No chat selected</p>
            <p>Pick an existing conversation or search for a username to start chatting.</p>
          </div>
        ) : null}

        {!disabled && messages.length === 0 ? (
          <div className="chat-window__empty">
            <p className="eyebrow">Conversation ready</p>
            <p>Send the first message to start this chat.</p>
          </div>
        ) : null}

        {!disabled
          ? messages.map((message) => (
              <ChatMessage
                key={message.id}
                isOwn={message.sender_id === currentUserID}
                isSelected={selectedMessageIDs.includes(message.id)}
                isSelectionMode={isSelectionMode}
                message={message}
                onToggleSelect={toggleSelection}
              />
            ))
          : null}
      </div>

      <MessageInput disabled={disabled || isWorking} onSend={onSend} />

      {deleteDialog ? (
        <div className="chat-dialog-backdrop" role="presentation">
          <div aria-modal="true" className="chat-dialog" role="dialog">
            <p className="eyebrow">
              {deleteDialog.kind === "clear" ? "Clear conversation" : "Delete options"}
            </p>
            <h3>{deleteDialog.title}</h3>
            <p className="support-copy">{deleteDialog.description}</p>

            <div className="chat-dialog__actions">
              <Button disabled={isWorking} onClick={() => setDeleteDialog(null)} type="button" variant="quiet">
                Cancel
              </Button>
              {deleteDialog.kind === "messages" ? (
                <Button disabled={isWorking} onClick={() => void handleDeleteConfirm("everyone")} type="button" variant="ghost">
                  Delete for everyone
                </Button>
              ) : null}
              {deleteDialog.kind === "messages" ? (
                <Button disabled={isWorking} onClick={() => void handleDeleteConfirm("me")} type="button">
                  Delete for me
                </Button>
              ) : (
                <>
                  <Button disabled={isWorking} onClick={() => void handleClearConfirm("everyone")} type="button" variant="ghost">
                    Delete for everyone
                  </Button>
                  <Button disabled={isWorking} onClick={() => void handleClearConfirm("me")} type="button">
                    Delete for me
                  </Button>
                </>
              )}
            </div>
          </div>
        </div>
      ) : null}
    </section>
  );
}
