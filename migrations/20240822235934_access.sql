-- +goose Up
-- +goose StatementBegin
create table accessible_roles(
      endpoint_address  text not null,
      role              text,
      created_at        timestamp default now() not null,
      primary key(endpoint_address, role)
);

insert into accessible_roles(endpoint_address, role) values ('/chat/v1/delete', 'ADMIN');
insert into accessible_roles(endpoint_address, role) values ('/chat/v1/create', 'ADMIN');
insert into accessible_roles(endpoint_address, role) values ('/chat/v1/send', 'ADMIN');
insert into accessible_roles(endpoint_address, role) values ('/chat/v1/send', 'USER');
commit;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table accessible_roles;
-- +goose StatementEnd
