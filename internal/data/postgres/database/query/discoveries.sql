-- name: CreatePlanetDiscovery :one
INSERT INTO planet_discoveries (system_x, system_y, orbit, discovered_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPlanetDiscovery :one
SELECT * FROM planet_discoveries
WHERE system_x = $1 AND system_y = $2 AND orbit = $3;

-- name: GetPlanetDiscoveriesBySystem :many
SELECT * FROM planet_discoveries
WHERE system_x = $1 AND system_y = $2;

-- name: GetPlanetDiscoveriesByAgent :many
SELECT * FROM planet_discoveries
WHERE discovered_by = $1;

-- name: IsPlanetDiscovered :one
SELECT EXISTS (
    SELECT 1
    FROM planet_discoveries
    WHERE system_x = $1 AND system_y = $2 AND orbit = $3
);

-- name: CreateSystemDiscovery :one
INSERT INTO system_discoveries (system_x, system_y, discovered_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSystemDiscovery :one
SELECT * FROM system_discoveries
WHERE system_x = $1 AND system_y = $2;

-- name: GetSystemDiscoveriesByAgent :many
SELECT * FROM system_discoveries
WHERE discovered_by = $1;

-- name: IsSystemDiscovered :one
SELECT EXISTS (
    SELECT 1
    FROM system_discoveries
    WHERE system_x = $1 AND system_y = $2
);