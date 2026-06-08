import { useState } from 'react';

import { createPost } from '../api/client';
import NavigationMenu from '../components/NavigationMenu.jsx';
import "../styles/createPost.css";

const MAX_CONTENT_LENGTH = 280;
const MAX_HASHTAGS = 3;

function CreatePostPage() {
    const [content, setContent] = useState('');
    const [hashtags, setHashtags] = useState([]);
    const [hashtagInput, setHashtagInput] = useState('');
    const [images, setImages] = useState([]);

    const [status, setStatus] = useState('');
    const [loading, setLoading] = useState(false);

    function handleAddHashtag() {
        const normalized = hashtagInput
            .trim()
            .replace('#', '')
            .toLowerCase();

        if (!normalized) {
            return;
        }

        if (hashtags.length >= MAX_HASHTAGS) {
            return;
        }

        if (hashtags.includes(normalized)) {
            return;
        }

        setHashtags((prev) => [...prev, normalized]);
        setHashtagInput('');
    }

    function handleRemoveHashtag(tag) {
        setHashtags((prev) =>
            prev.filter((item) => item !== tag)
        );
    }

    function handleImageUpload(event) {
        const files = Array.from(event.target.files || []);

        const mapped = files.slice(0, 3).map((file, index) => ({
            url: URL.createObjectURL(file),
            order_index: index + 1,
        }));

        setImages(mapped);
    }

    async function submitPost(postStatus) {
        if (!content.trim()) {
            setStatus('Введите текст поста');
            return;
        }

        try {
            setLoading(true);
            setStatus('');

            await createPost({
                content: content.trim(),

                images: images.map((image, index) => ({
                    url: image.url,
                    order_index: index + 1,
                })),

                hashtags: hashtags.map((tag, index) => ({
                    hashtag_name: tag,
                    order_index: index + 1,
                })),

                status: postStatus,
            });

            window.location.href = '/@me';
        } catch {
            setStatus('Не удалось создать пост');
        } finally {
            setLoading(false);
        }
    }

    return (
        <>
            <section className="create-post-page">
                <div className="create-post-card">
                    <h1>Создать пост</h1>

                    <textarea
                        value={content}
                        onChange={(event) =>
                            setContent(
                                event.target.value.slice(
                                    0,
                                    MAX_CONTENT_LENGTH,
                                ),
                            )
                        }
                        placeholder="Что происходит?"
                        className="create-post-textarea"
                    />

                    <div className="create-post-counter">
                        {content.length}/{MAX_CONTENT_LENGTH}
                    </div>

                    <div className="hashtags-section">
                        <div className="hashtags-input-row">
                            <input
                                type="text"
                                value={hashtagInput}
                                onChange={(event) =>
                                    setHashtagInput(event.target.value)
                                }
                                placeholder="Добавить хэштег"
                            />

                            <button
                                type="button"
                                onClick={handleAddHashtag}
                                disabled={hashtags.length >= 3}
                            >
                                Добавить
                            </button>
                        </div>

                        <div className="hashtags-list">
                            {hashtags.map((tag) => (
                                <button
                                    key={tag}
                                    type="button"
                                    className="hashtag-item"
                                    onClick={() =>
                                        handleRemoveHashtag(tag)
                                    }
                                >
                                    #{tag}
                                </button>
                            ))}
                        </div>
                    </div>

                    <div className="images-section">
                        <label className="upload-images-button">
                            Добавить изображения

                            <input
                                type="file"
                                multiple
                                accept="image/*"
                                hidden
                                onChange={handleImageUpload}
                            />
                        </label>

                        <div className="images-preview">
                            {images.map((image) => (
                                <img
                                    key={image.url}
                                    src={image.url}
                                    alt=""
                                />
                            ))}
                        </div>
                    </div>

                    <div className="create-post-actions">
                        <button
                            type="button"
                            className="draft-button"
                            disabled={loading}
                            onClick={() =>
                                submitPost('draft')
                            }
                        >
                            Сохранить в черновики
                        </button>

                        <button
                            type="button"
                            className="publish-button"
                            disabled={loading}
                            onClick={() =>
                                submitPost('published')
                            }
                        >
                            Опубликовать
                        </button>
                    </div>

                    {status && (
                        <p className="auth-status">
                            {status}
                        </p>
                    )}
                </div>
            </section>

            <NavigationMenu />
        </>
    );
}

export default CreatePostPage;