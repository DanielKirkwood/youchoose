create table if not exists users (
  id uuid default uuid_generate_v4() not null primary key,
  display_name text NOT NULL,
  email text unique,
  created timestamptz default now(),
  updated timestamptz default now()
);

---- create above / drop below ----

drop table if exists users;
