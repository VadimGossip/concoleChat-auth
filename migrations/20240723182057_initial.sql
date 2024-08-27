-- +goose Up
-- +goose StatementBegin
create table users(
  id         serial primary key,
  username   text not null unique,
  password   text not null,
  email      text not null unique,
  role       text not null,
  created_at timestamp default now() not null,
  updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
