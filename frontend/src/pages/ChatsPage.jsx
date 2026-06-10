import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { getConversations } from '../api/client';

import MainLayout from '../components/MainLayout.jsx';

import '../styles/chats.css';

function ChatsPage() {
    const [conversations, setConversations] =
        useState([]);

    const navigate = useNavigate();

    useEffect(() => {
        async function loadChats() {
            try {
                const data =
                    await getConversations();

                setConversations(
                    data.conversations || []
                );
            } catch (error) {
                console.error(error);
            }
        }

        loadChats();
    }, []);

    return (
        <MainLayout>
            <section className="chats-page">
                <h1 className="chats-title">
                    Сообщения
                </h1>

                {conversations.length === 0 ? (
                    <p className="chats-empty">
                        У вас пока нет диалогов
                    </p>
                ) : (
                    <div className="chats-list">
                        {conversations.map(
                            (conversation) => (
                                <div
                                    key={
                                        conversation.conversation_id
                                    }
                                    className="chat-card"
                                    onClick={() =>
                                        navigate(
                                            `/chats/${conversation.conversation_id}`
                                        )
                                    }
                                >
                                    <div className="chat-left">
                                        <div className="chat-avatar">
                                            {conversation
                                                .profile
                                                .avatar ? (
                                                <img
                                                    src={
                                                        conversation
                                                            .profile
                                                            .avatar
                                                    }
                                                    alt={
                                                        conversation
                                                            .profile
                                                            .name
                                                    }
                                                />
                                            ) : (
                                                <span>
                                                    {conversation.profile.name?.charAt(
                                                            0
                                                        ) ||
                                                        '?'}
                                                </span>
                                            )}
                                        </div>

                                        <div className="chat-info">
                                            <strong>
                                                {
                                                    conversation
                                                        .profile
                                                        .name
                                                }
                                            </strong>

                                            <span className="chat-username">
                                                @
                                                {
                                                    conversation
                                                        .profile
                                                        .username
                                                }
                                            </span>
                                        </div>
                                    </div>

                                    {conversation.unread_count >
                                        0 && (
                                            <div className="chat-unread">
                                                {
                                                    conversation.unread_count
                                                }
                                            </div>
                                        )}
                                </div>
                            )
                        )}
                    </div>
                )}
            </section>
        </MainLayout>
    );
}

export default ChatsPage;