function AuthLayout({ children, contentClassName = '', footer }) {
  return (
    <main className="auth-page">
      <header className="auth-header">
        <div className="logo-text" aria-label="Чирик">
          ЧИРИК
        </div>
      </header>

      <section className={`auth-content ${contentClassName}`}>{children}</section>

      {footer && <footer className="auth-footer">{footer}</footer>}
    </main>
  );
}

export default AuthLayout;
