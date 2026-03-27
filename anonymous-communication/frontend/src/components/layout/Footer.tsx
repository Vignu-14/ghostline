const socialLinks = [
  {
    href: "https://www.instagram.com/vignu_14/",
    label: "Instagram",
    icon: (
      <svg aria-hidden="true" viewBox="0 0 24 24">
        <path d="M7 2h10a5 5 0 0 1 5 5v10a5 5 0 0 1-5 5H7a5 5 0 0 1-5-5V7a5 5 0 0 1 5-5Zm0 2a3 3 0 0 0-3 3v10a3 3 0 0 0 3 3h10a3 3 0 0 0 3-3V7a3 3 0 0 0-3-3Zm5 3.5A4.5 4.5 0 1 1 7.5 12 4.5 4.5 0 0 1 12 7.5Zm0 2A2.5 2.5 0 1 0 14.5 12 2.5 2.5 0 0 0 12 9.5Zm5-3.25a1.25 1.25 0 1 1-1.25 1.25A1.25 1.25 0 0 1 17 6.25Z" />
      </svg>
    ),
  },
  {
    href: "https://github.com/Vignu-14  ",
    label: "GitHub",
    icon: (
      <svg aria-hidden="true" viewBox="0 0 24 24">
        <path d="M12 2a10 10 0 0 0-3.16 19.49c.5.09.68-.21.68-.48v-1.68c-2.78.6-3.37-1.18-3.37-1.18a2.66 2.66 0 0 0-1.11-1.46c-.91-.62.07-.61.07-.61a2.11 2.11 0 0 1 1.53 1 2.14 2.14 0 0 0 2.92.84 2.14 2.14 0 0 1 .64-1.35c-2.22-.25-4.55-1.11-4.55-4.94a3.86 3.86 0 0 1 1-2.68 3.58 3.58 0 0 1 .1-2.64s.84-.27 2.75 1a9.55 9.55 0 0 1 5 0c1.91-1.3 2.75-1 2.75-1a3.58 3.58 0 0 1 .1 2.64 3.86 3.86 0 0 1 1 2.68c0 3.84-2.34 4.69-4.57 4.93a2.39 2.39 0 0 1 .68 1.86V21c0 .27.18.58.69.48A10 10 0 0 0 12 2Z" />
      </svg>
    ),
  },
  {
    href: "https://www.linkedin.com/in/vijnesh-m-ab1468313/",
    label: "LinkedIn",
    icon: (
      <svg aria-hidden="true" viewBox="0 0 24 24">
        <path d="M6.94 8.5H3.56V20h3.38ZM5.25 3A1.97 1.97 0 1 0 7.22 5 1.98 1.98 0 0 0 5.25 3ZM20.44 12.33c0-3.08-1.64-4.51-3.83-4.51a3.32 3.32 0 0 0-3 1.65V8.5h-3.24c0 .64.04 11.5 0 11.5h3.24v-6.42c0-.34 0-.68.13-.92a1.76 1.76 0 0 1 1.65-1.18c1.17 0 1.64.89 1.64 2.19V20h3.24Z" />
      </svg>
    ),
  },
];

export function Footer() {
  return (
    <footer className="site-footer">
      <div className="site-footer__copy">
        <p className="eyebrow">Ghostline</p>
        <p>Private social, realtime chat, and audited admin access.</p>
      </div>

      <nav aria-label="Social links" className="site-footer__socials">
        {socialLinks.map((link) => (
          <a
            className="site-footer__social"
            href={link.href}
            key={link.label}
            rel="noreferrer"
            target="_blank"
          >
            {link.icon}
            <span>{link.label}</span>
          </a>
        ))}
      </nav>
    </footer>
  );
}
