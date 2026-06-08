import { useState } from 'react';
import { login } from '../api/client';
import AuthLayout from '../components/AuthLayout.jsx';

function LoginPage() {
  const [credentials, setCredentials] = useState({ email: '', password: '' });
  const [status, setStatus] = useState('');

  async function handleLogin(event) {
    event.preventDefault();
    setStatus('');

    try {
      const data = await login(credentials);

      const token =
          data.accessToken ||
          data.access_token ||
          data.token;

      if (token) {
        localStorage.setItem('accessToken', token);
        window.location.href = '/';
      }
    } catch {
      setStatus('Не удалось войти');
    }
  }

  return (
    <AuthLayout
      footer={
        <a className="auth-switch-link" href="/register">
          <span>Нет аккаунта?</span>
          <strong>Зарегистрируйся!</strong>
        </a>
      }
    >
      <form className="auth-form" onSubmit={handleLogin}>
        <h1>Войти в аккаунт</h1>

        <input
          type="email"
          value={credentials.email}
          onChange={(event) => setCredentials((value) => ({ ...value, email: event.target.value }))}
          placeholder="Email"
          autoComplete="email"
          required
        />

        <input
          type="password"
          value={credentials.password}
          onChange={(event) => setCredentials((value) => ({ ...value, password: event.target.value }))}
          placeholder="Пароль"
          autoComplete="current-password"
          required
        />

        <button type="submit">Войти</button>

        {status && <p className="auth-status">{status}</p>}
      </form>
    </AuthLayout>
  );
}

export default LoginPage;
