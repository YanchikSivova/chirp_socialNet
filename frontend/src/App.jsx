import { useEffect, useState } from 'react';
import AccountInfoPage from './pages/AccountInfoPage.jsx';
import LoginPage from './pages/LoginPage.jsx';
import RegisterPage from './pages/RegisterPage.jsx';
import VerifyCodePage from './pages/VerifyCodePage.jsx';
import UserProfilePage from './pages/UserProfilePage.jsx';
import { getMe } from './api/client';
import CreatePostPage from "./pages/CreatePostPage.jsx";
import { useLocation } from 'react-router-dom';
import FeedPage from './pages/FeedPage.jsx';
import SearchPage from "./pages/SearchPage.jsx";
import PostPage from "./pages/PostPage.jsx";
function App() {
  const [loading, setLoading] = useState(true);
  const [me, setMe] = useState(null);

  const location = useLocation();

  const path = location.pathname;

  useEffect(() => {
    async function loadMe() {
      const token = localStorage.getItem('accessToken');

      if (!token) {
        setLoading(false);
        return;
      }

      try {
        const profile = await getMe();
        setMe(profile);
      } catch (error) {
        const message =
            error?.response?.data?.error ||
            error?.response?.data?.message ||
            '';

        if (message.includes('profile is not completed')) {
          setMe({ needsCompletion: true });
        } else {
          localStorage.removeItem('accessToken');
        }
      } finally {
        setLoading(false);
      }
    }

    loadMe();
  }, []);

  if (loading) {
    return null;
  }

  if (path === '/register') {
    return <RegisterPage />;
  }

  if (path === '/verify-code') {
    return <VerifyCodePage />;
  }

  if (path === '/account-info') {
    return <AccountInfoPage />;
  }

  if (!me) {
    return <LoginPage />;
  }

  if (me.needsCompletion) {
    window.location.href = '/account-info';
    return null;
  }

  if (path === '/') {
    return <FeedPage />;
  }

  if (path === '/add') {
    return <CreatePostPage />;
  }
  if (path === '/search') {
    return <SearchPage me={me} />;
  }
  if (path.startsWith('/posts/')) {
    return <PostPage />;
  }

  if (path.startsWith('/@')) {
    return <UserProfilePage me={me} />;
  }

  window.location.href = '/';
  return null;
}

export default App;