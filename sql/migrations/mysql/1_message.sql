create table if not exists message
(
    id      binary(16) primary key not null,
    content text not null
);
