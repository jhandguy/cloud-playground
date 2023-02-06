create table if not exists "messages"
(
    id      uuid primary key default gen_random_uuid(),
    content text not null,
    user_id uuid references "users"(id)
);
