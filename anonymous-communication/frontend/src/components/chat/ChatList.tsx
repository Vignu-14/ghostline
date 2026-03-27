import type { Conversation } from "../../types/message";
import type { UserSearchResult } from "../../types/user";
import { Avatar } from "../common/Avatar";
import { UnreadBadge } from "./UnreadBadge";

type ChatListProps = {
  conversations: Conversation[];
  activeUserID?: string;
  onSelect: (conversation: Conversation) => void;
  searchQuery: string;
  onSearchChange: (value: string) => void;
  searchResults: UserSearchResult[];
  isSearching: boolean;
  onStartConversation: (user: UserSearchResult) => void;
};

export function ChatList({
  conversations,
  activeUserID,
  onSelect,
  searchQuery,
  onSearchChange,
  searchResults,
  isSearching,
  onStartConversation,
}: ChatListProps) {
  return (
    <aside className="chat-list">
      <div className="chat-search">
        <label className="field">
          <span className="field__label">Search usernames</span>
          <input
            className="field__input"
            onChange={(event) => onSearchChange(event.target.value)}
            placeholder="Type even 1 letter..."
            value={searchQuery}
          />
        </label>

        {searchQuery.trim().length >= 1 ? (
          <div className="chat-search__results">
            {isSearching ? <p className="support-copy">Searching users...</p> : null}

            {!isSearching && searchResults.length === 0 ? (
              <p className="support-copy">No matching users found.</p>
            ) : null}

            {searchResults.map((user) => (
              <button
                key={user.id}
                className={`chat-search__item ${activeUserID === user.id ? "is-active" : ""}`}
                onClick={() => onStartConversation(user)}
                type="button"
              >
                <Avatar alt={user.username} src={user.profile_picture_url} />
                <span className="chat-list__copy">
                  <strong>@{user.username}</strong>
                  <span>Start a new conversation</span>
                </span>
              </button>
            ))}
          </div>
        ) : null}
      </div>

      {conversations.length === 0 ? (
        <div className="chat-list__empty">
          <p className="eyebrow">No chats yet</p>
          <p>Search for a username above to start your first conversation.</p>
        </div>
      ) : null}

      {conversations.map((conversation) => (
        <button
          key={conversation.user_id}
          className={`chat-list__item ${activeUserID === conversation.user_id ? "is-active" : ""}`}
          onClick={() => onSelect(conversation)}
          type="button"
        >
          <Avatar alt={conversation.username} src={conversation.profile_picture_url} />
          <span className="chat-list__copy">
            <strong>@{conversation.username}</strong>
            <span>{conversation.last_message || "No messages yet"}</span>
          </span>
          <UnreadBadge count={conversation.unread_count} />
        </button>
      ))}
    </aside>
  );
}
