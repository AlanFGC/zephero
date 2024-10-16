-- name: WorldChunk :one
SELECT * FROM world_chunk WHERE
world_id = $1 AND row_id = $2 AND cold_id = $3 LIMIT 1;

-- name: WorldChunkbyWorldId :many
SELECT * FROM world_chunk WHERE
world_id = ? ORDER BY row_id, col_id;


-- name: CreateWorldChunk :one
INSERT INTO world_chunk (
  world_id, row_id, cold_id, locked, data
)