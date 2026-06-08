import { useState } from 'react';
import { register } from '../api/client';
import AuthLayout from '../components/AuthLayout.jsx';

function RegisterPage() {
  const [credentials, setCredentials] = useState({ email: '', password: '', confirmPassword: '' });
  const [status, setStatus] = useState('');
  const [isSuccess, setIsSuccess] = useState(false);

  async function handleRegister(event) {
    event.preventDefault();
    setStatus('');
    setIsSuccess(false);

    if (credentials.password !== credentials.confirmPassword) {
      setStatus('Пароли не совпадают');
      return;
    }

    try {
      const data = await register({
        email: credentials.email,
        password: credentials.password,
      });

      localStorage.setItem('verificationToken', data.token);
      localStorage.setItem('verificationEmail', credentials.email);

      window.location.pathname = '/verify-code';
    } catch {
      setStatus('Не удалось зарегистрироваться');
    }
  }

  return (
    <AuthLayout
      footer={
        <a className="auth-switch-link" href="/">
          <span>Есть аккаунт?</span>
          <strong>Войди!</strong>
        </a>
      }
    >
      <form className="auth-form register-form" onSubmit={handleRegister}>
        <h1>Создать аккаунт</h1>

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
          autoComplete="new-password"
          required
        />

        <input
          type="password"
          value={credentials.confirmPassword}
          onChange={(event) => setCredentials((value) => ({ ...value, confirmPassword: event.target.value }))}
          placeholder="Повтори пароль"
          autoComplete="new-password"
          required
        />

        <button type="submit">Регистрация</button>

        {status && <p className={`auth-status ${isSuccess ? 'success' : ''}`}>{status}</p>}
      </form>
    </AuthLayout>
  );
}

export default RegisterPage;
