import { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

import {
    getPostComments,
    createComment,
    createCommentLike,
    deleteCommentLike,
} from '../api/client';

import PostCard from '../components/PostCard.jsx';
import NavigationMenu from '../components/NavigationMenu.jsx';

import likeIcon from '../assets/like.svg';
import likeFilledIcon from '../assets/like-filled.svg';

import '../styles/postPage.css';

function PostPage() {
    const location = useLocation();

    const id = location.pathname.replace('/posts/', '');

    const [post, setPost] = useState(null);
    const [comments, setComments] = useState([]);
    const [commentText, setCommentText] = useState('');

    async function loadComments() {
        const data = await getPostComments(id);

        setComments(data?.comments || []);
    }

    useEffect(() => {
        const storedPost =
            sessionStorage.getItem(`post-${id}`);

        if (storedPost) {
            setPost(JSON.parse(storedPost));
        }

        loadComments();
    }, [id]);

    async function handleSendComment() {
        if (!commentText.trim()) return;

        await createComment(id, {
            content: commentText,
        });

        setCommentText('');

        loadComments();
    }

    async function handleCommentLike(commentId, isLiked) {
        if (isLiked) {
            await deleteCommentLike(commentId);
        } else {
            await createCommentLike(commentId);
        }

        loadComments();
    }

    if (!post) {
        return null;
    }

    return (
        <>
            <section className="post-page">
                <PostCard
                    post={post}
                    disableNavigation={true}
                />

                <div className="comments-list">
                    {comments.map((comment) => (
                        <div
                            key={comment.comment_id}
                            className="comment-card"
                        >
                            <div className="comment-header">
                                <strong>
                                    {comment.author.name}
                                </strong>

                                <span>
                                    @{comment.author.username}
                                </span>
                            </div>

                            <p className="comment-content">
                                {comment.content}
                            </p>

                            <button
                                className="comment-like-btn"
                                onClick={() =>
                                    handleCommentLike(
                                        comment.comment_id,
                                        comment.is_liked
                                    )
                                }
                            >
                                <img
                                    src={
                                        comment.is_liked
                                            ? likeFilledIcon
                                            : likeIcon
                                    }
                                    alt=""
                                />

                                <span>
                                    {comment.likes_amount}
                                </span>
                            </button>
                        </div>
                    ))}
                </div>
            </section>

            <div className="comment-input-wrapper">
                <input
                    value={commentText}
                    onChange={(e) =>
                        setCommentText(e.target.value)
                    }
                    placeholder="Написать комментарий..."
                />

                <button onClick={handleSendComment}>
                    Отправить
                </button>
            </div>
        </>
    );
}

export default PostPage;