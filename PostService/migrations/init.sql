create table if not exists post(
    post_id UUID primary key,
    profile_id UUID not null,
    content varchar(250) not null check (length(content) > 0),
    published_at timestamp,
    last_edited_at timestamp not null default now(),
    status text not null default 'draft' check(status in ('draft', 'published', 'banned','deleted')),
    likes_amount int not null default 0 check(likes_amount >= 0),
    comments_amount int not null default 0 check(comments_amount >= 0),
    reposts_amount int not null default 0 check(reposts_amount >= 0),
    reports_amount int not null default 0 check(reports_amount >= 0)
);

create table if not exists image(
    image_id UUID primary key,
    post_id UUID not null,
    url text,
    order_index int not null default 1 check(order_index in (1, 2, 3)),
    foreign key (post_id) references post(post_id) on delete cascade
);

create table if not exists hashtag(
    hashtag_id UUID primary key,
    hashtag_name text not null check(length(hashtag_name)>0)
);

create table if not exists post_hashtag(
    post_hashtag_id UUID primary key,
    post_id UUID not null,
    hashtag_id UUID not null,
    order_index int not null default 1 check(order_index in (1, 2, 3)),
    foreign key(post_id) references post(post_id) on delete cascade,
    foreign key(hashtag_id) references hashtag(hashtag_id)
);

create table if not exists repost(
    repost_id UUID primary key,
    post_id UUID not null,
    profile_id UUID not null,
    unique(post_id, profile_id),
    created_at timestamp not null default now(),
    foreign key (post_id) references post(post_id) on delete cascade
);

create table if not exists likes(
    likes_id UUID primary key,
    post_id UUID not null,
    profile_id UUID not null,
    unique(post_id, profile_id),
    liked_at timestamp not null default now(),
    foreign key(post_id) references post(post_id) on delete cascade
);

create table if not exists report(
    report_id UUID primary key,
    post_id UUID not null,
    profile_id UUID,
    unique(post_id, profile_id),
    reason text not null check(length(reason)>0),
    created_at timestamp not null default now(),
    foreign key(post_id) references post(post_id) on delete cascade
);

create table if not exists comment(
    comment_id UUID primary key,
    post_id UUID not null,
    profile_id UUID not null,
    parent_comment_id UUID default null,
    content varchar(250) not null check(length(content)>0),
    likes_amount int not null default 0 check(likes_amount >= 0),
    created_at timestamp not null default now(),
    foreign key(post_id) references post(post_id) on delete cascade,
    foreign key(parent_comment_id) references comment(comment_id) on delete cascade
);

create table if not exists comment_like(
    comment_like_id UUID primary key,
    comment_id UUID not null,
    profile_id UUID not null,
    unique(profile_id, comment_id),
    liked_at timestamp not null default now(),
    foreign key(comment_id) references comment(comment_id) on delete cascade
);

--Обновление количества лайков
create or replace function update_likes_amount()
returns trigger as $$
BEGIN
if TG_OP = 'INSERT' THEN
    update post
    set likes_amount = likes_amount + 1
    where post_id = NEW.post_id;

    return NEW;
elsif TG_OP = 'DELETE' THEN
    update post
    set likes_amount = likes_amount - 1
    where post_id = OLD.post_id;

    return OLD;
end if;
return null;
end;
$$ LANGUAGE plpgsql;

create trigger trigger_likes_amount
after insert or delete on likes
for each row
execute function update_likes_amount();

-- Обновление количества репостов
create or replace function update_reposts_amount()
returns trigger as $$
begin
    if tg_op='INSERT' then
        update post
        set reposts_amount = reposts_amount + 1
        where post_id = NEW.post_id;

        return NEW;

    elsif tg_op='DELETE' then
        update post
        set reposts_amount = reposts_amount - 1
        where post_id = OLD.post_id;

        return OLD;
    end if;
return null;
end;
$$ language plpgsql;

create trigger trigger_reposts_amount
after insert or delete on repost
for each row
execute function update_reposts_amount();

-- Обновление количества комментариев
create or replace function update_comments_amount()
returns trigger as $$
begin
    if tg_op='INSERT' then
        update post
        set comments_amount = comments_amount + 1
        where post_id = NEW.post_id;

        return NEW;

    elsif tg_op='DELETE' then
        update post
        set comments_amount = comments_amount - 1
        where post_id = OLD.post_id;

        return OLD;
    end if;
return null;
end;
$$ language plpgsql;

create trigger trigger_comments_amount
after insert or delete on comment
for each row
execute function update_comments_amount();

-- Обновление количества лайков на комментарии
create or replace function update_comment_likes_amount()
returns trigger as $$
begin
    if tg_op='INSERT' then
        update comment
        set likes_amount = likes_amount+1
        where comment_id = NEW.comment_id;

        return NEW;

    elsif tg_op='DELETE' then
        update comment
        set likes_amount = likes_amount-1
        where comment_id = OLD.comment_id;

        return OLD;
    end if;
return null;
end;
$$ language plpgsql;

create trigger trigger_comment_likes_amount
after insert or delete on comment_like
for each row
execute function update_comment_likes_amount();

-- Обновление количества жалоб с изменением статуса
create or replace function update_reports_amount()
returns trigger as $$
declare
    reports_count int;
begin
    if tg_op='INSERT' then
        update post
        set reports_amount = reports_amount + 1
        where post_id = NEW.post_id
        returning reports_amount into reports_count;

        if reports_count >= 10 then
            update post
            set status = 'banned'
            where post_id = NEW.post_id;
        end if;

        return NEW;

    elsif tg_op='DELETE' then
        update post
        set reports_amount = reports_amount -1
        where post_id = OLD.post_id;

        return OLD;
    end if;
return null;
end;
$$ language plpgsql;

create trigger trigger_reports_amount
after insert or delete on report
for each row
execute function update_reports_amount();

create or replace function set_published_at()
returns trigger as $$
begin
    if tg_op = 'UPDATE' then
        if NEW.status = 'published' and OLD.status != 'published' then
            NEW.published_at = NOW();
        end if;

        return NEW;

    elsif tg_op = 'INSERT' then
        if NEW.status = 'published' then
            NEW.published_at = NOW();
        end if;
        return NEW;
    end if;
    return null;
end;
$$ language plpgsql;

create trigger trigger_set_published_at
before update or insert on post
for each row
execute function set_published_at();