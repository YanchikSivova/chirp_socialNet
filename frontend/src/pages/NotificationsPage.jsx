import { useEffect, useState } from 'react';

import { getNotifications } from '../api/client';

import NavigationMenu from '../components/NavigationMenu.jsx';

import '../styles/notifications.css';
import MainLayout from "../components/MainLayout.jsx";

function NotificationsPage() {
    const [notifications, setNotifications] = useState([]);

    useEffect(() => {
        async function loadNotifications() {
            try {
                const data = await getNotifications();

                setNotifications(
                    data.notifications || []
                );
            } catch (error) {
                console.error(error);
            }
        }

        loadNotifications();
    }, []);

    function getNotificationText(type) {
        switch (type) {
            case 'like':
                return 'лайкнул ваш пост';

            case 'subscription':
                return 'подписался на вас';

            case 'repost':
                return 'сделал репост вашего поста';

            case 'comment':
                return 'прокомментировал ваш пост';

            case 'comment_like':
                return 'лайкнул ваш комментарий';

            case 'answer':
                return 'ответил на ваш комментарий';

            default:
                return 'взаимодействовал с вами';
        }
    }

    return (
        <>
            <MainLayout>
            <section className="notifications-page">
                <h1 className="notifications-title">
                    Уведомления
                </h1>

                {notifications.length === 0 ? (
                    <p className="notifications-empty">
                        Пока нет уведомлений
                    </p>
                ) : (
                    <div className="notifications-list">
                        {notifications.map((item, index) => (
                            <div
                                key={
                                    item.notification.notification_id ||
                                    index
                                }
                                className="notification-card"
                            >
                                <div className="notification-avatar">
                                    {item.actor.avatar ? (
                                        <img
                                            src={item.actor.avatar}
                                            alt={
                                                item.actor.name
                                            }
                                        />
                                    ) : (
                                        <span>
                                            {item.actor.name?.charAt(
                                                0
                                            ) || '?'}
                                        </span>
                                    )}
                                </div>

                                <div className="notification-content">
                                    <strong>
                                        {item.actor.name}
                                    </strong>

                                    <span className="notification-username">
                                        @{item.actor.username}
                                    </span>

                                    <span className="notification-text">
    {' '}
                                        {getNotificationText(
                                            item.notification.entity_type
                                        )}
                                    </span>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </section>


                </MainLayout>
        </>
    );
}

export default NotificationsPage;