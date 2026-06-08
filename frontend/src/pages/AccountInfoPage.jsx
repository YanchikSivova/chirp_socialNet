import { useState } from 'react';
import { fillProfile, mockUploadAvatar } from '../api/client';
import AuthLayout from '../components/AuthLayout.jsx';

const initialProfile = {
  name: '',
  username: '',
  description: '',
  avatar: '',
};

function normalizeUsername(value) {
  return value.replace(/^@+/, '').replace(/\s/g, '').toLowerCase();
}

function AccountInfoPage() {
  const [profile, setProfile] = useState(initialProfile);
  const [avatarPreview, setAvatarPreview] = useState('');
  const [status, setStatus] = useState('');
  const [isSuccess, setIsSuccess] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [normalizedUsername, setNormalizedUsername] = useState('');

  async function handleAvatarChange(event) {
    const file = event.target.files?.[0];

    if (!file) {
      return;
    }

    setStatus('');
    setIsUploading(true);
    setAvatarPreview(URL.createObjectURL(file));

    try {
      const avatarUrl = await mockUploadAvatar(file);
      setProfile((value) => ({ ...value, avatar: avatarUrl }));
    } catch {
      setStatus('Не удалось добавить аватар');
    } finally {
      setIsUploading(false);
    }
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setStatus('');
    setIsSuccess(false);

    try {
      const normalized = normalizeUsername(profile.username);
      setNormalizedUsername(normalized);
      
      await fillProfile({
        name: profile.name.trim(),
        username: normalized,
        description: profile.description.trim(),
        avatar: profile.avatar,
      });

      setIsSuccess(true);
      setStatus('Профиль сохранен');
      
      setTimeout(() => {
        window.location.href = '/@' + normalized;
      }, 1500);
    } catch (err) {
      console.error('Failed to save profile:', err);
      setStatus('Не удалось сохранить профиль');
    }
  }

  return (
    <AuthLayout contentClassName="account-content">
      <form className="account-form" onSubmit={handleSubmit}>
        <h1>Заполни профиль</h1>

        <label className="avatar-picker">
          <input type="file" accept="image/*" onChange={handleAvatarChange} />
          <span className="avatar-preview">
            {avatarPreview ? <img src={avatarPreview} alt="" /> : <span>{profile.name.charAt(0) || '+'}</span>}
          </span>
          <strong>{isUploading ? 'Загрузка...' : 'Добавить аватар'}</strong>
        </label>

        <input
          type="text"
          value={profile.name}
          onChange={(event) => setProfile((value) => ({ ...value, name: event.target.value }))}
          placeholder="Name"
          autoComplete="name"
          required
        />

        <label className="username-field">
          <span>@</span>
          <input
            type="text"
            value={profile.username}
            onChange={(event) =>
              setProfile((value) => ({ ...value, username: normalizeUsername(event.target.value) }))
            }
            placeholder="username"
            autoComplete="username"
            required
          />
        </label>

        <textarea
          value={profile.description}
          onChange={(event) => setProfile((value) => ({ ...value, description: event.target.value }))}
          placeholder="Short description"
          maxLength="160"
        />

        <button type="submit">Продолжить</button>

        {status && <p className={`auth-status ${isSuccess ? 'success' : ''}`}>{status}</p>}
      </form>
    </AuthLayout>
  );
}

export default AccountInfoPage;
