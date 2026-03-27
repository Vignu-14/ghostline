import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as userService from "../../services/userService";
import type { UserSearchResult } from "../../types/user";
import { getErrorMessage } from "../../utils/errorHandler";
import { Avatar } from "../common/Avatar";

export function UserSearchPanel() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<UserSearchResult[]>([]);
  const [error, setError] = useState("");
  const [isSearching, setIsSearching] = useState(false);

  useEffect(() => {
    const trimmedQuery = query.trim();
    if (trimmedQuery.length < 1) {
      setResults([]);
      setError("");
      setIsSearching(false);
      return;
    }

    setIsSearching(true);

    const timeoutID = window.setTimeout(() => {
      void (async () => {
        try {
          const response = await userService.searchUsers(trimmedQuery, 10);
          setResults(response.users);
          setError("");
        } catch (searchError) {
          setError(getErrorMessage(searchError, "Unable to search users."));
        } finally {
          setIsSearching(false);
        }
      })();
    }, 180);

    return () => {
      window.clearTimeout(timeoutID);
    };
  }, [query]);

  return (
    <section className="panel user-search-panel">
      <p className="eyebrow">Discover people</p>
      <h2>Search Ghostline users</h2>
      <label className="field">
        <span className="field__label">Username</span>
        <input
          className="field__input"
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Type even 1 letter..."
          value={query}
        />
      </label>

      {error ? <p className="form-error">{error}</p> : null}
      {isSearching ? <p className="support-copy">Searching users...</p> : null}

      {query.trim().length >= 1 && !isSearching && results.length === 0 ? (
        <p className="support-copy">No matching users found.</p>
      ) : null}

      <div className="user-search-panel__results">
        {results.map((user) => (
          <Link className="user-search-panel__result" key={user.id} to={`/u/${user.username}`}>
            <Avatar alt={user.username} src={user.profile_picture_url} />
            <span className="chat-list__copy">
              <strong>@{user.username}</strong>
              <span>Open profile</span>
            </span>
          </Link>
        ))}
      </div>
    </section>
  );
}
