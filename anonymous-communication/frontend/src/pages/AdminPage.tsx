import { AuditLogViewer } from "../components/admin/AuditLogViewer";
import { ImpersonationPanel } from "../components/admin/ImpersonationPanel";
import { UserManagement } from "../components/admin/UserManagement";

export function AdminPage() {
  return (
    <main className="page page--admin">
      <section className="page-header">
        <p className="eyebrow">Restricted operator surface</p>
        <h1>God Mode stays deliberate, logged, and time-boxed.</h1>
        <p className="support-copy">
          This panel is designed around step-up authentication first, not broad destructive controls.
        </p>
      </section>

      <div className="admin-grid">
        <ImpersonationPanel />
        <AuditLogViewer />
        <UserManagement />
      </div>
    </main>
  );
}
