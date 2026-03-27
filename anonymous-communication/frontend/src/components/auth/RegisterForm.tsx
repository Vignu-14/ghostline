import { useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { getErrorMessage } from "../../utils/errorHandler";
import { Button } from "../common/Button";
import { Input } from "../common/Input";

export function RegisterForm() {
  const navigate = useNavigate();
  const { register } = useAuth();
  const [form, setForm] = useState({
    username: "",
    email: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);
    setError("");

    try {
      await register(form);
      navigate("/");
    } catch (submitError) {
      setError(getErrorMessage(submitError, "Unable to create account."));
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form className="auth-form" onSubmit={handleSubmit}>
      <Input label="Username" value={form.username} onChange={(event) => setForm((current) => ({ ...current, username: event.target.value }))} />
      <Input label="Email" type="email" value={form.email} onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))} />
      <Input label="Password" type="password" value={form.password} onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))} />
      {error ? <p className="form-error">{error}</p> : null}
      <Button disabled={isSubmitting} type="submit">
        {isSubmitting ? "Creating..." : "Create account"}
      </Button>
    </form>
  );
}
