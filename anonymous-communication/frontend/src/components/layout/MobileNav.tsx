import { NavLink } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";

export function MobileNav() {
  const { user } = useAuth();

  return (
    <nav className="mobile-nav">
      <NavLink to="/">Feed</NavLink>
      <NavLink to="/chat">Chat</NavLink>
      {user ? <NavLink to={`/u/${user.username}`}>Profile</NavLink> : null}
      {user?.role === "admin" ? <NavLink to="/admin">Admin</NavLink> : null}
    </nav>
  );
}
