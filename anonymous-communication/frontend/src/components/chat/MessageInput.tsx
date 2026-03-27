import { useState, type FormEvent } from "react";
import { Button } from "../common/Button";

type MessageInputProps = {
  disabled?: boolean;
  onSend: (content: string) => Promise<void> | void;
};

export function MessageInput({ disabled, onSend }: MessageInputProps) {
  const [content, setContent] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!content.trim()) {
      return;
    }

    await onSend(content);
    setContent("");
  }

  return (
    <form className="message-input" onSubmit={handleSubmit}>
      <input
        disabled={disabled}
        onChange={(event) => setContent(event.target.value)}
        placeholder="Write a message..."
        value={content}
      />
      <Button disabled={disabled || !content.trim()} type="submit">
        Send
      </Button>
    </form>
  );
}
