create table if not exists notification(
    notification_id uuid primary key,
    profile_id uuid not null,
    actor_id uuid not null,
    type text not null check(type in ('subscription', 'like', 'repost', 'comment', 'answer', 'comment_like')),
    entity_id uuid not null,
    created_at timestamp not null default now()
);

create table if not exists processed_events(
    event_id UUID primary key,
    processed_at timestamp not null default now()
);