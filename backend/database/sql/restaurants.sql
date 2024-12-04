-- name: SearchRestaurants :many
select id, name, address_line1, address_line2, postcode, ts_rank(search_vector, websearch_to_tsquery(@search_term::text)) as rank
from restaurants
where search_vector @@ websearch_to_tsquery('english', @search_term::text)
and valid = true
order by rank desc
limit $1 offset $2;


-- name: GetNearestRestaurants :many
with user_location as (
    select st_setsrid(st_makepoint(@user_longitude::float, @user_latitude::float), 4326) as location
)
select
    r.id,
    r.name,
    r.address_line1,
    r.address_line2,
    r.postcode,
    r.longitude,
    r.latitude,
    st_distance(ul.location, r.geolocation) as distance_meters
from
    restaurants r,
    user_location ul
where
    r.geolocation is not null
    and st_dwithin(ul.location, r.geolocation, @max_radius::float)
    and valid = true
order by
    distance_meters asc
limit $1 offset $2;
