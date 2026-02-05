-- name: AddPoint :exec
INSERT INTO points (title, point) -- lon/lat (X/Y)
VALUES (@name, st_setsrid(st_makepoint(@lon::float8, @lat::float8), 4326));

-- name: GetPointsFromBox :many
SELECT id, title AS name, st_y(point)::float8 AS lat, st_x(point)::float8 AS lon
FROM points
WHERE point && st_makeenvelope(@lon1::float8, @lat1::float8, @lon2::float8, @lat2::float8, 4326);