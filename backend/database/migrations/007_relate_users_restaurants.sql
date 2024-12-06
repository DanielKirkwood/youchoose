alter table restaurants add column created_by uuid references users(id);

alter table restaurants add constraint chk_creator check (
  (created_by is null and valid = true) or
  (created_by is not null)
);

---- create above / drop below ----

alter table restaurants drop column created_by;
