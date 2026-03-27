import { NavLink } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { APP_NAME } from "../../utils/constants";
import { LogoutButton } from "../auth/LogoutButton";

export function Navbar() {
  const { user } = useAuth();

  return (
    <header className="navbar">
      <NavLink className="brand" to="/">
        <span className="brand__kicker">Private social</span>
        <span className="brand__name">{APP_NAME}</span>
      </NavLink>

      <nav className="navbar__links">
        <NavLink to="/">Feed</NavLink>
        <NavLink to="/chat">Chat</NavLink>
        {user?.role === "admin" ? <NavLink to="/admin">God Mode</NavLink> : null}
      </nav>

      <div className="navbar__meta">
        {user ? (
          <>
            <NavLink className="navbar__user" to={`/u/${user.username}`}>
              @{user.username}
            </NavLink>
            <LogoutButton />
          </>
        ) : (
          <>
            <NavLink to="/login">Login</NavLink>
            <NavLink to="/register">Register</NavLink>
          </>
        )}
      </div>
    </header>
  );
}
