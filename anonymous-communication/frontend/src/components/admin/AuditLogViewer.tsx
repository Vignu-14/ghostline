export function AuditLogViewer() {
  return (
    <section className="panel">
      <p className="eyebrow">Audit visibility</p>
      <h3>Server audit logs are active.</h3>
      <p>
        The backend already records impersonation activity, but this UI view needs a dedicated audit listing endpoint before it can render live records.
      </p>
    </section>
  );
}
