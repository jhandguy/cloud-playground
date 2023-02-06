create table if not exists "users"
(
    id   uuid primary key default gen_random_uuid(),
    name text not null
);
