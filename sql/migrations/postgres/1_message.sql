create table if not exists "message"
(
    id      uuid primary key default gen_random_uuid(),
    content text not null
);
