import { PostFeed } from "../components/post/PostFeed";
import { UserSearchPanel } from "../components/user/UserSearchPanel";
import { useAuth } from "../hooks/useAuth";

export function HomePage() {
  const { user } = useAuth();

  return (
    <main className="page page--home">
      <section className="hero">
        <div className="hero__copy reveal-up">
          <p className="eyebrow">Quiet social, louder signal</p>
          <h1>Private moments, public images, instant conversations.</h1>
          <p className="support-copy">
            Ghostline blends the intimacy of direct messaging with the rhythm of a visual feed, wrapped in an operator-grade admin layer.
          </p>
        </div>
        <div className="hero__meta reveal-up-delay">
          <span>Feed</span>
          <span>Realtime chat</span>
          <span>Admin step-up auth</span>
        </div>
      </section>

      {user ? <UserSearchPanel /> : null}

      <PostFeed />
    </main>
  );
}
