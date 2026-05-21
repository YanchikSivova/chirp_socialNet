create table if not exists profile(
    profile_id UUID primary key,
    name text,
    username text unique,
    avatar text,
    description text,
    subscribers_amount int not null default 0 check(subscribers_amount >= 0),
    subscribed_amount int not null default 0 check(subscribed_amount >=0),
    posts_amount int not null default 0 check(posts_amount >=0),
    is_completed boolean not null default false
);

create table if not exists credentials(
    credentials_id UUID primary key,
    profile_id UUID,
    email text unique not null,
    hashed_password text not null,
    created_at timestamp not null default now(),
    status text not null default 'pending' check (status in ('pending', 'active')),
    foreign key (profile_id) references profile(profile_id)
);

create table if not exists verification(
    verification_id uuid primary key,
    credentials_id uuid not null unique,
    code varchar(6) not null check (length(code) = 6),
    expires_at timestamp not null,
    foreign key (credentials_id) references credentials(credentials_id) on delete cascade
);

create table if not exists refresh_token(
    refresh_token_id uuid primary key,
    profile_id uuid not null,
    refresh_token text not null,
    expires_at timestamp not null,
    foreign key (profile_id) references profile(profile_id) on delete cascade
);

create table if not exists subscription(
    subscription_id uuid primary key,
    subscriber_id uuid not null,
    subscribed_id uuid not null,
    unique(subscriber_id, subscribed_id),
    foreign key(subscriber_id) references profile(profile_id) on delete cascade,
    foreign key(subscribed_id) references profile(profile_id) on delete cascade
);

create table if not exists blacklist(
    blacklist_id uuid primary key,
    profile_id uuid not null,
    banned_profile_id uuid not null,
    unique(profile_id, banned_profile_id),
    foreign key (profile_id) references profile(profile_id) on delete cascade,
    foreign key(banned_profile_id) references profile(profile_id) on delete cascade
);

create table if not exists email_change(
    email_change_id uuid primary key,
    credentials_id uuid unique not null,
    new_email text not null,
    foreign key (credentials_id) references credentials(credentials_id) on delete cascade
);

create table if not exists password_change(
    password_change_id uuid primary key,
    credentials_id uuid unique not null,
    new_password_hash text not null,
    foreign key (credentials_id) references credentials(credentials_id) on delete cascade
);

create table if not exists processed_events(
    event_id UUID primary key,
    processed_at timestamp not null default now()
);

--Обновление количества подписчиков и подписок после подписки/отписки
create or replace function update_subscription_counts()
returns trigger as $$
begin
if TG_OP = 'INSERT' THEN
    UPDATE profile
    SET subscribed_amount = subscribed_amount+1
    WHERE profile_id = NEW.subscriber_id;

    UPDATE profile
    SET subscribers_amount = subscribers_amount+1
    WHERE profile_id = NEW.subscribed_id;

    RETURN NEW;
ELSIF TG_OP = 'DELETE' THEN
    UPDATE profile
    SET subscribed_amount = subscribed_amount-1
    WHERE profile_id = OLD.subscriber_id;

    UPDATE profile
    SET subscribers_amount = subscribers_amount-1
    WHERE profile_id = OLD.subscribed_id;

    RETURN OLD;
END IF;

RETURN NULL;
END;

$$ LANGUAGE plpgsql;

-- Отписки после добавления в черный список
CREATE TRIGGER trigger_subscription_counts
AFTER INSERT OR DELETE ON subscription
FOR EACH ROW
EXECUTE FUNCTION update_subscription_counts();

CREATE OR REPLACE FUNCTION handle_blacklist_insert()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM subscription
    WHERE subscriber_id = NEW.banned_profile_id
    AND subscribed_id = NEW.profile_id;

    DELETE FROM subscription
    WHERE subscriber_id = NEW.profile_id
    AND subscribed_id = NEW.banned_profile_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_blacklist_insert
AFTER INSERT ON blacklist
FOR EACH ROW
EXECUTE FUNCTION handle_blacklist_insert();