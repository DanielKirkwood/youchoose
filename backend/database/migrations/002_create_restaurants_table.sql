create extension if not exists "uuid-ossp";

create table restaurants (
  id uuid default uuid_generate_v4() not null primary key,
  fhrs_id int unique not null,
  name text not null,
  address_line1 text,
  address_line2 text,
  address_line3 text,
  address_line4 text,
  postcode text,
  latitude double precision,
  longitude double precision,
  business_type text,
  valid boolean default false,
  created timestamptz default now(),
  updated timestamptz default now()
);

---- create above / drop below ----

drop table if exists restaurants;
