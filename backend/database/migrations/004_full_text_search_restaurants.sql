-- add a full-text search column
alter table restaurants add column search_vector tsvector generated always as (to_tsvector('english', name || ' ' || address_line2 || ' ' || postcode)) stored;

-- create an index on the full-text search vector for fast searches
create index restaurants_search_vector_idx on restaurants using gin (search_vector);

---- create above / drop below ----

drop index if exists restaurants_search_vector_idx;
alter table restaurants drop column if exists search_vector;
