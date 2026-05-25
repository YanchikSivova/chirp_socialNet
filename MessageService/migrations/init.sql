create table if not exists conversation(
    conversation_id  uuid primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    last_message_id uuid
);

create table if not exists message(
    message_id uuid primary key,
    conversation_id uuid not null,
    sender_id uuid not null,
    type text not null check(type in ('text', 'post', 'image')),
    content text not null check (length(content)>0),
    created_at timestamp not null default now(),
    foreign key(conversation_id) references conversation(conversation_id) on delete cascade
);

do $$
begin
  if not exists(
      select 1 from pg_constraint where conname = 'fk_last_message'
  ) then
      alter table conversation add constraint fk_last_message
      foreign key (last_message_id) references message(message_id)
      on delete set null;
  end if;
end $$;

create table if not exists conversation_member(
    conversation_member_id uuid primary key,
    conversation_id uuid not null,
    profile_id uuid not null,
    last_read_message_id uuid,
    unread_count int not null default 0,
    foreign key (conversation_id) references conversation(conversation_id) on delete cascade,
    foreign key (last_read_message_id) references message(message_id),
    unique(conversation_id, profile_id)
);

create index if not exists idx_message_conversation_created on message(conversation_id, created_at desc);

create index if not exists idx_conversation_member_profile on conversation_member(profile_id);

create index if not exists idx_conversation_update_at on conversation(updated_at desc);

create or replace function update_conversation_timestamp()
returns trigger as $$
begin
    update conversation
    set updated_at=NEW.created_at, last_message_id = NEW.message_id
    where conversation_id = NEW.conversation_id;

    return NEW;
end;
$$ language plpgsql;

create trigger trigger_update_conversation_timestamp
after insert on message
for each row
execute function update_conversation_timestamp();

create or replace function increment_unread_count()
returns trigger as $$
begin
    update conversation_member
    set unread_count = unread_count+1
    where conversation_id =  NEW.conversation_id and profile_id != NEW.sender_id;

    return NEW;
end;
$$ language plpgsql;

create trigger trigger_increment_unread_count
after insert on message
for each row
execute function increment_unread_count();

create or replace function check_conversation_member_limit()
returns trigger as $$
begin
    if( select count(*) from conversation_member
        where conversation_id = NEW.conversation_id) >= 2 then
        raise exception 'conversation cannot have more than 2 members';
    end if;
    return NEW;
end;
$$ language plpgsql;

create trigger trigger_check_conversation_member_limit
before insert on conversation_member
for each row
execute function check_conversation_member_limit();

create or replace function check_sender_in_conversation()
returns trigger as $$
begin if not exists(
    select 1 from conversation_member
    where conversation_id = NEW.conversation_id and profile_id = NEW.sender_id)
    then raise exception 'sender is not member of conversation';
    end if;
    return NEW;
end;
$$ language plpgsql;

create trigger trigger_check_sender_in_conversation
before insert on message
for each row
execute function check_sender_in_conversation();