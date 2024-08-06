-- +goose Up
-- +goose StatementBegin
create table audit_log(
  id          serial primary key,
  action      text not null,
  call_params text not null,
  created_at  timestamp default now() not null
);

create index in_audit_log_action on audit_log(action);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table audit_log;
-- +goose StatementEnd
