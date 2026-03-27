import { useState, type FormEvent } from "react";
import { apiRequest } from "../../services/api";
import { getErrorMessage } from "../../utils/errorHandler";
import { Button } from "../common/Button";
import { Input } from "../common/Input";

export function ImpersonationPanel() {
  const [targetUserID, setTargetUserID] = useState("");
  const [impersonationPassword, setImpersonationPassword] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  async function handleStart(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setMessage("");

    try {
      await apiRequest("/api/admin/impersonate", {
        method: "POST",
        body: JSON.stringify({
          target_user_id: targetUserID,
          impersonation_password: impersonationPassword,
        }),
      });
      setMessage("Impersonation started. Refreshing the page will make the cookie swap obvious.");
    } catch (requestError) {
      setError(getErrorMessage(requestError, "Unable to start impersonation."));
    }
  }

  async function handleStop() {
    setError("");
    setMessage("");

    try {
      await apiRequest("/api/admin/impersonate/stop", {
        method: "POST",
      });
      setMessage("Impersonation ended.");
    } catch (requestError) {
      setError(getErrorMessage(requestError, "Unable to stop impersonation."));
    }
  }

  return (
    <section className="admin-panel">
      <div>
        <p className="eyebrow">God mode</p>
        <h2>Step up before you cross into a user session.</h2>
        <p className="support-copy">
          Every impersonation request is audited on the backend. Use the exact target user UUID and your admin second password.
        </p>
      </div>

      <form className="auth-form" onSubmit={handleStart}>
        <Input label="Target user UUID" onChange={(event) => setTargetUserID(event.target.value)} value={targetUserID} />
        <Input
          label="Impersonation password"
          onChange={(event) => setImpersonationPassword(event.target.value)}
          type="password"
          value={impersonationPassword}
        />
        {error ? <p className="form-error">{error}</p> : null}
        {message ? <p className="form-success">{message}</p> : null}
        <div className="admin-panel__actions">
          <Button type="submit">Start impersonation</Button>
          <Button onClick={() => void handleStop()} type="button" variant="ghost">
            Stop impersonation
          </Button>
        </div>
      </form>
    </section>
  );
}
