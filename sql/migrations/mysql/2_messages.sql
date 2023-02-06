create table if not exists messages
(
    id      binary(16) primary key not null,
    content text not null,
    user_id binary(16) references users(id)
);
