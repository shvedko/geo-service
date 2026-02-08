CREATE EXTENSION IF NOT EXISTS postgis;
CREATE TABLE IF NOT EXISTS points
(
    id    SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    point GEOMETRY(Point, 4326) NOT NULL -- 4326 это WGS84 (стандарт для GPS)
);
CREATE INDEX IF NOT EXISTS points_point_idx ON points USING GIST (point);
