import { NavLink } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";

export function MobileNav() {
  const { user } = useAuth();

  return (
    <nav aria-label="Mobile navigation" className="mobile-nav">
      <div className="mobile-nav__shell">
        <NavLink to="/">Feed</NavLink>
        <NavLink to="/chat">Chat</NavLink>
        {user ? (
          <>
            <NavLink to={`/u/${user.username}`}>Profile</NavLink>
            {user.role === "admin" ? <NavLink to="/admin">Admin</NavLink> : null}
          </>
        ) : (
          <>
            <NavLink to="/login">Login</NavLink>
            <NavLink to="/register">Register</NavLink>
          </>
        )}
      </div>
    </nav>
  );
}
