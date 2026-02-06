-- name: AddPoint :one
INSERT INTO points (title, point) -- lon/lat (X/Y)
VALUES (@name, st_setsrid(st_makepoint(@lon::float8, @lat::float8), 4326))
RETURNING id;

-- name: GetPointsFromBox :many
SELECT id, title AS name, st_x(point)::float8 AS lon, st_y(point)::float8 AS lat
FROM points
WHERE point && st_makeenvelope(@lon1::float8, @lat1::float8, @lon2::float8, @lat2::float8, 4326);

-- name: GetPoint :one
SELECT id, title AS name, st_x(point)::float8 AS lon, st_y(point)::float8 AS lat
FROM points
WHERE id = @id;
