import { useRef, useState } from 'react';
import { verifyCode, resendCode } from '../api/client';
import AuthLayout from '../components/AuthLayout.jsx';

function VerifyCodePage() {
  const [code, setCode] = useState(['', '', '', '', '', '']);
  const [status, setStatus] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const inputRefs = useRef([]);

  const handleInputChange = (index, value) => {
    if (!/^\d*$/.test(value)) return;

    const newCode = [...code];
    newCode[index] = value.slice(0, 1);
    setCode(newCode);

    if (value && index < 5) {
      inputRefs.current[index + 1]?.focus();
    }
  };

  const handleKeyDown = (index, event) => {
    if (event.key === 'Backspace' && !code[index] && index > 0) {
      inputRefs.current[index - 1]?.focus();
    }
  };

  const handlePaste = (event) => {
    event.preventDefault();
    const pastedData = event.clipboardData.getData('text');
    const digits = pastedData.replace(/\D/g, '').slice(0, 6).split('');

    const newCode = [...code];
    digits.forEach((digit, index) => {
      if (index < 6) newCode[index] = digit;
    });
    setCode(newCode);

    if (digits.length === 6) {
      inputRefs.current[5]?.focus();
    }
  };

  const handleVerify = async (event) => {
    event.preventDefault();
    const verificationCode = code.join('');

    if (verificationCode.length !== 6) {
      setStatus('Введите код из 6 цифр');
      return;
    }

    setIsLoading(true);
    setStatus('');

    try {
      const token = localStorage.getItem('verificationToken');
      const data = await verifyCode({ code: verificationCode, token });

      if (data.success) {
        localStorage.removeItem('verificationToken');
        localStorage.removeItem('verificationEmail');

        window.location.pathname = '/';
      }
    } catch (error) {
      setStatus('Неверный код. Попробуй ещё раз');
    } finally {
      setIsLoading(false);
    }
  };

  const handleResend = async () => {
    setIsLoading(true);
    setStatus('');

    try {
      const email = localStorage.getItem('verificationEmail');

      await resendCode({ email });

      setStatus('Код отправлен ещё раз');
      setCode(['', '', '', '', '', '']);

      inputRefs.current[0]?.focus();
    } catch {
      setStatus('Не удалось отправить код');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthLayout
      contentClassName="verify-content"
      footer={
        <button className="resend-link" onClick={handleResend} disabled={isLoading}>
          <span>Не пришел код?</span>
          <strong>Отправить повторно!</strong>
        </button>
      }
    >
      <form className="auth-form verify-form" onSubmit={handleVerify}>
        <h1>Введите код</h1>

        <div className="code-inputs">
          {code.map((digit, index) => (
            <input
              key={index}
              ref={(el) => (inputRefs.current[index] = el)}
              type="text"
              value={digit}
              onChange={(event) => handleInputChange(index, event.target.value)}
              onKeyDown={(event) => handleKeyDown(index, event)}
              onPaste={handlePaste}
              placeholder="0"
              maxLength="1"
              inputMode="numeric"
              autoComplete="off"
              disabled={isLoading}
            />
          ))}
        </div>

        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Ждём...' : 'Тык!'}
        </button>

        {status && <p className={`auth-status ${status.includes('Код отправлен') ? 'success' : ''}`}>{status}</p>}
      </form>
    </AuthLayout>
  );
}

export default VerifyCodePage;
