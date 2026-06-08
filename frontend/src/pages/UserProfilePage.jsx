import { useState, useEffect } from 'react';
import {
  getMe,
  getUserProfile,
  getUserPosts,
  getMyPublishedPosts,
  getMyDraftPosts,
    logout,
  followUser,
  unfollowUser,
} from '../api/client';
import logoutIcon from '../assets/logout.svg';
import DraftCard from '../components/DraftCard.jsx';
import NavigationMenu from '../components/NavigationMenu.jsx';
import PostCard from '../components/PostCard.jsx';
import { useLocation, useNavigate } from 'react-router-dom';

function UserProfilePage({ me }) {
  const [profile, setProfile] = useState(null);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] =
      useState('published');
  const [relationship, setRelationship] =
      useState('');
  const location = useLocation();
  const navigate = useNavigate();
  const pathname = location.pathname;

  const username = pathname.replace('/@', '');

  const isMeRoute = username === 'me';
  const isMyProfile =
      isMeRoute ||
      (me &&
          (
              me.username === username ||
              me.profile_id === username
          ));
  useEffect(() => {
    async function loadProfile() {
      try {
        setLoading(true);

        let profileData;
        let postsData;

        if (isMyProfile) {
          profileData = await getMe();

          postsData =
              activeTab === 'published'
                  ? await getMyPublishedPosts()
                  : await getMyDraftPosts();

          setProfile(profileData);
        } else {
          profileData = await getUserProfile(username);

          postsData = await getUserPosts(username);

          setProfile(
              profileData.profile || profileData
          );

          setRelationship(
              profileData.relationship || ''
          );
        }

        setPosts(
            Array.isArray(postsData)
                ? postsData
                : postsData?.posts || []
        );
      } catch {
        setError('Не удалось загрузить профиль');
      } finally {
        setLoading(false);
      }
    }

    loadProfile();
  }, [username, me, activeTab]);

  async function handleLogout() {
    try {
      await logout();
    } catch {
    } finally {
      localStorage.removeItem('accessToken');

      window.location.href = '/';
    }
  }
  async function handleFollow() {
    try {
      if (relationship === 'subscribed') {
        await unfollowUser(profile.profile_id);

        setRelationship('');
      } else {
        await followUser(profile.profile_id);

        setRelationship('subscribed');
      }
    } catch (error) {
      console.error(error);
    }
  }
  if (loading) {
    return (
        <section className="auth-content">
          <p>Загрузка...</p>
        </section>
    );
  }

  if (error) {
    return (
        <section className="auth-content">
          <p className="auth-status">{error}</p>
        </section>
    );
  }

  if (!profile) {
    return (
        <section className="auth-content">
          <p>Профиль не найден</p>
        </section>
    );
  }

  return (
      <>
        <section className="auth-content profile-content">
          {isMyProfile && (
              <button
                  className="profile-logout-btn"
                  onClick={handleLogout}
              >
                <img
                    src={logoutIcon}
                    alt="Выйти"
                />
              </button>
          )}
          <div className="profile-info">
            <div className="profile-avatar">
              {profile.avatar ? (
                  <img src={profile.avatar} alt={profile.name} />
              ) : (
                  <span>{profile.name?.charAt(0) || '?'}</span>
              )}
            </div>

            <h1 className="profile-name">{profile.name}</h1>

            <p className="profile-username">
              @{profile.username}
            </p>

            {profile.description && (
                <p className="profile-description">
                  {profile.description}
                </p>
            )}
            <div className="profile-stats">
              <div className="profile-stat">
                <strong>{profile.posts_amount}</strong>
                <span>постов</span>
              </div>

              <div className="profile-stat">
                <strong>{profile.subscribers_amount}</strong>
                <span>подписчиков</span>
              </div>

              <div className="profile-stat">
                <strong>{profile.subscribed_amount}</strong>
                <span>подписок</span>
              </div>
            </div>
          </div>

          <div className="profile-posts">
            {isMyProfile ? (
                <div className="profile-tabs">
                  <button
                      className={`profile-tab ${
                          activeTab === 'published'
                              ? 'active'
                              : ''
                      }`}
                      onClick={() =>
                          setActiveTab('published')
                      }
                  >
                    Посты
                  </button>

                  <button
                      className={`profile-tab ${
                          activeTab === 'draft'
                              ? 'active'
                              : ''
                      }`}
                      onClick={() =>
                          setActiveTab('draft')
                      }
                  >
                    Черновики
                  </button>
                </div>
            ) : (
                <div className="profile-actions">
                  <button
                      className={`follow-btn ${
                          relationship === 'subscribed'
                              ? 'following'
                              : ''
                      }`}
                      onClick={handleFollow}
                  >
                    {relationship === 'subscribed'
                        ? 'Вы подписаны'
                        : 'Подписаться'}
                  </button>

                  <button
                      className="message-btn"
                      onClick={() =>
                          navigate('/chats')
                      }
                  >
                    Сообщение
                  </button>
                </div>
            )}

            {posts.length === 0 ? (
                <p className="no-posts">
                  Пока нет постов
                </p>
            ) : (
                posts.map((post) =>
                    activeTab === 'draft' ? (
                        <DraftCard
                            key={post.post_id}
                            post={post}
                        />
                    ) : (
                        <PostCard
                            key={post.post_id}
                            post={post}
                        />
                    ),
                )
            )}
          </div>
        </section>

        <NavigationMenu />
      </>
  );
}

export default UserProfilePage;