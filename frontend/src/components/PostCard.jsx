import { useState } from 'react';
import "../styles/postCard.css"

import likeIcon from '../assets/like.svg';
import likeFilledIcon from '../assets/like-filled.svg';

import repostIcon from '../assets/repost.svg';
import repostFilledIcon from '../assets/repost-filled.svg';

import commentIcon from '../assets/comment.svg';
import shareIcon from '../assets/send.svg';
import { useNavigate } from 'react-router-dom';

import {
    createLike,
    deleteLike,
    createRepost,
    deleteRepost,
} from '../api/client';

function PostCard({ post, onLike, onRepost, disableNavigation = false, }) {
    const [isLiked, setIsLiked] = useState(post.is_liked);
    const [isReposted, setIsReposted] = useState(post.is_reposted);

    const [likesAmount, setLikesAmount] = useState(post.likes_amount);
    const [repostsAmount, setRepostsAmount] = useState(post.reposts_amount);
    const navigate = useNavigate();
    const formattedDate = new Date(
        post.published_at,
    ).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
    });
    const sortedHashtags = [...(post.hashtags || [])]
        .sort((a, b) => a.order_index - b.order_index);

    const sortedImages = [...(post.images || [])]
        .sort((a, b) => a.order_index - b.order_index);

    async function handleLike() {
        try {
            let response;

            if (isLiked) {
                response = await deleteLike(post.post_id);
            } else {
                response = await createLike(post.post_id);
            }

            if (response.success) {
                setIsLiked((prev) => !prev);

                setLikesAmount((prev) =>
                    isLiked ? prev - 1 : prev + 1
                );
            }
        } catch (error) {
            console.error('Like error:', error);
        }
    }

    async function handleRepost() {
        try {
            let response;

            if (isReposted) {
                response = await deleteRepost(post.post_id);
            } else {
                response = await createRepost(post.post_id);
            }

            if (response.success) {
                setIsReposted((prev) => !prev);

                setRepostsAmount((prev) =>
                    isReposted ? prev - 1 : prev + 1
                );
            }
        } catch (error) {
            if (
                error.response?.status === 400 &&
                error.response?.data?.error === 'cannot create repost'
            ) {
                alert('Нельзя репостить собственный пост');
                return;
            }

            console.error('Repost error:', error);
        }
    }

    return (
        <article
            className={`post-card ${
                !disableNavigation ? 'clickable-post' : ''
            }`}
            onClick={() => {
                if (disableNavigation) return;

                const postId =
                    post.post_id || post.id;

                sessionStorage.setItem(
                    `post-${postId}`,
                    JSON.stringify(post)
                );

                navigate(`/posts/${postId}`);
            }}
        >
            <header
                className="post-header clickable-author"
                onClick={(e) => {
                    e.stopPropagation();
                    if (post.author?.is_me) {
                        navigate('/@me');
                    } else {
                        navigate(
                            `/@${post.author.profile_id}`
                        );
                    }
                }}
            >
                <div className="post-author-avatar">
                    {post.author.avatar ? (
                        <img
                            src={post.author.avatar}
                            alt={post.author.name}
                        />
                    ) : (
                        <span>
              {post.author.name?.charAt(0) || '?'}
            </span>
                    )}
                </div>

                <div className="post-author-info">
                    <strong className="post-author-name">
                        {post.author.name}
                    </strong>

                    <div className="post-meta">
        <span className="post-author-username">
            @{post.author.username}
        </span>

                        <span className="post-date">
            · {formattedDate}
        </span>
                    </div>
                </div>
            </header>

            <div className="post-body">
                <p className="post-content">
                    {post.content}
                </p>

                {sortedHashtags.length > 0 && (
                    <div className="post-hashtags">
                        {sortedHashtags.map((tag) => (
                            <span
                                key={tag.hashtag_name}
                                className="post-hashtag"
                            >
                #{tag.hashtag_name}
              </span>
                        ))}
                    </div>
                )}

                {sortedImages.length > 0 && (
                    <div
                        className={`post-images post-images-${sortedImages.length}`}
                    >
                        {sortedImages.map((image) => (
                            <img
                                key={image.url}
                                src={image.url}
                                alt=""
                                className="post-image"
                            />
                        ))}
                    </div>
                )}
            </div>

            <footer className="post-stats">
                <button
                    className="post-stat-button"

                    onClick={(e) => {
                        e.stopPropagation();
                        handleLike();
                    }}
                >
                    <img
                        src={isLiked ? likeFilledIcon : likeIcon}
                        alt=""
                    />

                    <span>{likesAmount}</span>
                </button>

                <button className="post-stat-button">
                    <img src={commentIcon} alt="" />

                    <span>{post.comments_amount}</span>
                </button>

                <button
                    className="post-stat-button"
                    onClick={(e) => {
                        e.stopPropagation();
                        handleRepost();
                    }}
                >
                    <img
                        src={
                            isReposted
                                ? repostFilledIcon
                                : repostIcon
                        }
                        alt=""
                    />

                    <span>{repostsAmount}</span>
                </button>

                <button
                    className="post-stat-button"
                    onClick={(e) => {
                        e.stopPropagation();
                    }}
                >
                    <img src={shareIcon} alt="" />
                </button>
            </footer>
        </article>
    );
}

export default PostCard;