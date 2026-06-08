import { useEffect, useState } from 'react';

import { getFeed } from '../api/client';

import PostCard from '../components/PostCard.jsx';
import "../styles/feed.css";
import MainLayout from '../components/MainLayout.jsx';
import searchIcon from "../assets/search.svg";
import {useNavigate} from "react-router-dom";
function FeedPage() {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    const navigate = useNavigate();

    useEffect(() => {
        async function loadFeed() {
            try {
                setLoading(true);

                const response = await getFeed();

                setPosts(
                    Array.isArray(response)
                        ? response
                        : response?.posts || []
                );
            } catch {
                setError('Не удалось загрузить ленту');
            } finally {
                setLoading(false);
            }
        }

        loadFeed();
    }, []);

    return (
        <MainLayout>
            <section className="feed-page">
                <div className="feed-topbar">
                    <button
                        className="feed-search-btn"
                        onClick={() => navigate('/search')}
                    >
                        <img
                            src={searchIcon}
                            alt="Поиск"
                        />
                    </button>
                </div>

                {loading ? (
                    <p>Загрузка...</p>
                ) : error ? (
                    <p className="auth-status">
                        {error}
                    </p>
                ) : posts.length === 0 ? (
                    <p className="no-posts">
                        Постов пока нет
                    </p>
                ) : (
                    <div className="feed-posts">
                        {posts.map((post) => (
                            <PostCard
                                key={post.post_id}
                                post={post}
                            />
                        ))}
                    </div>
                )}
            </section>
        </MainLayout>
    );
}

export default FeedPage;