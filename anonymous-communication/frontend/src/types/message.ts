export interface Conversation {
  user_id: string;
  username: string;
  profile_picture_url?: string | null;
  unread_count: number;
  last_message?: string;
  last_message_at?: string;
}

export interface Message {
  id: string;
  sender_id: string;
  receiver_id: string;
  content: string;
  is_read: boolean;
  created_at: string;
  deleted_for_everyone: boolean;
}

export interface SendMessageRequest {
  receiver_id: string;
  content: string;
}

export type DeleteMode = "me" | "everyone";

export interface DeleteMessagesRequest {
  message_ids: string[];
  mode: DeleteMode;
}

export interface ClearConversationRequest {
  mode: DeleteMode;
}
