import '../styles/draftCard.css'
function DraftCard({ post }) {
    const formattedDate = new Date(
        post.last_edited_at,
    ).toLocaleString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });

    const hashtags = [...(post.hashtags || [])]
        .sort((a, b) => a.order_index - b.order_index);

    const images = [...(post.images || [])]
        .sort((a, b) => a.order_index - b.order_index);

    return (
        <article className="draft-card">
            <div className="draft-card-header">
        <span className="draft-badge">
          Черновик
        </span>

                <time className="draft-date">
                    {formattedDate}
                </time>
            </div>

            <p className="draft-content">
                {post.content || 'Пустой черновик'}
            </p>

            {hashtags.length > 0 && (
                <div className="draft-hashtags">
                    {hashtags.map((hashtag) => (
                        <span
                            key={hashtag.order_index}
                            className="draft-hashtag"
                        >
              #{hashtag.hashtag_name}
            </span>
                    ))}
                </div>
            )}

            {images.length > 0 && (
                <div
                    className={`draft-images images-count-${images.length}`}
                >
                    {images.map((image) => (
                        <img
                            key={image.order_index}
                            src={image.url}
                            alt=""
                            className="draft-image"
                        />
                    ))}
                </div>
            )}
        </article>
    );
}

export default DraftCard;