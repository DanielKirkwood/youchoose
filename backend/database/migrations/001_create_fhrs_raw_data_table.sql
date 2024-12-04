create table if not exists fhrs_raw_data (
  fhrs_id bigint primary key,
  local_authority_business_id text,
  business_name text,
  business_type text,
  business_type_id int,
  address_line1 text,
  address_line2 text,
  address_line3 text,
  postcode text,
  rating_value text,
  rating_key text,
  rating_date text,
  local_authority_code int,
  local_authority_name text,
  local_authority_website text,
  local_authority_email text,
  scheme_type text,
  new_rating_pending boolean,
  longitude float,
  latitude float,
  created timestamptz default now(),
  updated timestamptz default now()
);
---- create above / drop below ----
drop table if exists fhrs_raw_data;
