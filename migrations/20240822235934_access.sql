-- +goose Up
-- +goose StatementBegin
create table accessible_roles(
      endpoint_address  text not null,
      role              text,
      created_at        timestamp default now() not null,
      primary key(endpoint_address, role)
);

insert into accessible_roles(endpoint_address, role) values ('/chat_v1.ChatV1/Delete', 'ADMIN');
insert into accessible_roles(endpoint_address, role) values ('/chat_v1.ChatV1/Create', 'ADMIN');
insert into accessible_roles(endpoint_address, role) values ('/chat_v1.ChatV1/SendMessage', 'ADMIN');
insert into accessible_roles(endpoint_address, role) values ('/chat_v1.ChatV1/SendMessage', 'USER');
commit;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table accessible_roles;
-- +goose StatementEnd
