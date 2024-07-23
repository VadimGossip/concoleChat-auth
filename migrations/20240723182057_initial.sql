-- +goose Up
-- +goose StatementBegin
create table users(
  id         serial primary key,
  username   varchar(255) not null,
  password   varchar(255) not null,
  email      varchar(255) not null unique,
  role       varchar(20) not null,
  created_at timestamp default now() not null,
  updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
