import { Link } from "react-router-dom";

export function NotFoundPage() {
  return (
    <main className="center-stage">
      <div className="panel">
        <p className="eyebrow">404</p>
        <h1>This route drifted out of range.</h1>
        <p>Head back to the feed and pick up the thread again.</p>
        <Link to="/">Return home</Link>
      </div>
    </main>
  );
}
