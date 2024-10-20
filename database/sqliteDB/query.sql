-- name: GetWorld :one
SELECT *
FROM world
WHERE world_id = ?;

-- name: UpdateWorld :exec
UPDATE world
SET row_length = ?, column_length = ?, chunk_length = ?
WHERE world_id = ?;

-- name: InsertWorld :execlastid
INSERT INTO world (row_length, column_length, chunk_length)
VALUES (?, ?, ?);

-- name: GetWorldChunk :one
SELECT *
FROM world_chunk
WHERE world_id = ? AND row_id = ? AND col_id = ?;

-- name: GetWorldChunkByWorldId :many
SELECT *
FROM world_chunk
WHERE world_id = ?
ORDER BY row_id, col_id;

-- name: InsertWorldChunk :execlastid
INSERT INTO world_chunk (world_id, row_id, col_id, locked, chunk)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateWorldChunk :exec
UPDATE world_chunk
SET last_updated = CURRENT_TIMESTAMP, locked = ?, chunk = ?
WHERE world_id = ? AND row_id = ? AND col_id = ?;

-- name: DeleteWorldChunk :exec
DELETE FROM world_chunk
WHERE world_id = ? AND row_id = ? AND col_id = ?;