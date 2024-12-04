create index restaurants_name_idx on restaurants (name);
create index restaurants_address_idx on restaurants (address_line2);
create index restaurants_postcode_idx on restaurants (postcode);

---- create above / drop below ----

drop index if exists restaurants_name_idx;
drop index if exists restaurants_address_idx;
drop index if exists restaurants_postcode_idx;
