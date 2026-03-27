import type {
  ClearConversationRequest,
  Conversation,
  DeleteMessagesRequest,
  Message,
  SendMessageRequest,
} from "../types/message";
import type { UserSearchResult } from "../types/user";
import { apiRequest } from "./api";

type ConversationResponse = {
  user?: {
    id?: string;
    username?: string;
    profile_picture_url?: string | null;
  };
  unread_count?: number;
  last_message?: {
    content?: string;
    created_at?: string;
  };
};

export function listConversations(page = 1, limit = 20) {
  return apiRequest<{ conversations: ConversationResponse[] | null; page: number; limit: number }>(
    `/api/messages/conversations?page=${page}&limit=${limit}`,
  ).then((response) => ({
    ...response,
    conversations: Array.isArray(response.conversations)
      ? response.conversations
          .map((conversation) => normalizeConversation(conversation))
          .filter((conversation): conversation is Conversation => Boolean(conversation))
      : [],
  }));
}

export function getConversation(userID: string, page = 1, limit = 50) {
  return apiRequest<{ messages: Message[] | null; page: number; limit: number }>(
    `/api/messages/${userID}?page=${page}&limit=${limit}`,
  ).then((response) => ({
    ...response,
    messages: Array.isArray(response.messages) ? response.messages : [],
  }));
}

export function sendMessage(input: SendMessageRequest) {
  return apiRequest<{ message: Message }>("/api/messages", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function searchUsers(query: string, limit = 8) {
  const params = new URLSearchParams({
    q: query,
    limit: String(limit),
  });

  return apiRequest<{ users: UserSearchResult[] | null; query: string; limit: number }>(
    `/api/users/search?${params.toString()}`,
  ).then((response) => ({
    ...response,
    users: Array.isArray(response.users) ? response.users : [],
  }));
}

export function deleteMessages(input: DeleteMessagesRequest) {
  return apiRequest<{ deleted_count: number; mode: string }>("/api/messages/delete", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function clearConversation(userID: string, input: ClearConversationRequest) {
  return apiRequest<{ deleted_count: number }>(`/api/messages/${userID}/clear`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

function normalizeConversation(conversation: ConversationResponse): Conversation | null {
  if (!conversation.user?.id || !conversation.user.username) {
    return null;
  }

  return {
    user_id: conversation.user.id,
    username: conversation.user.username,
    profile_picture_url: conversation.user.profile_picture_url,
    unread_count: conversation.unread_count || 0,
    last_message: conversation.last_message?.content || "",
    last_message_at: conversation.last_message?.created_at || "",
  };
}
