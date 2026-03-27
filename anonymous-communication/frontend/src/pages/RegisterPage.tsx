import { Link } from "react-router-dom";
import { RegisterForm } from "../components/auth/RegisterForm";

export function RegisterPage() {
  return (
    <main className="auth-page">
      <section className="auth-page__hero">
        <p className="eyebrow">Create your lane</p>
        <h1>Open a clean account and start sharing in minutes.</h1>
        <p className="support-copy">
          Ghostline gives you a public visual feed, private one-to-one messaging, and a security model built for serious demos.
        </p>
      </section>

      <section className="auth-page__panel">
        <RegisterForm />
        <p className="auth-page__switch">
          Already registered? <Link to="/login">Login</Link>
        </p>
      </section>
    </main>
  );
}
