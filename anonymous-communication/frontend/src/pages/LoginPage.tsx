import { Link } from "react-router-dom";
import { LoginForm } from "../components/auth/LoginForm";

export function LoginPage() {
  return (
    <main className="auth-page">
      <section className="auth-page__hero">
        <p className="eyebrow">Welcome back</p>
        <h1>Step back into the private side of the feed.</h1>
        <p className="support-copy">
          Your session lives in secure cookies, your messages land in realtime, and admin actions stay audited.
        </p>
      </section>

      <section className="auth-page__panel">
        <LoginForm />
        <p className="auth-page__switch">
          Need an account? <Link to="/register">Register</Link>
        </p>
      </section>
    </main>
  );
}
