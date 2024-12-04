-- install postgis extension if it's not already installed
create extension if not exists postgis;

-- add the geospatial columns for latitude and longitude in the restaurants table
alter table restaurants add column geolocation geography(point, 4326) generated always as (st_setsrid(st_makepoint(longitude, latitude), 4326)) stored;

-- create a gist index on the geolocation column for fast location-based queries
create index restaurants_geolocation_idx on restaurants using gist (geolocation);

---- create above / drop below ----

drop index if exists restaurants_geolocation_idx;
alter table restaurants drop column if exists geolocation;
drop extension if exists postgis cascade;
