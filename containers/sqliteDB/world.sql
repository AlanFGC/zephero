DROP TABLE IF EXISTS world_chunks;

CREATE TABLE world_chunk (
  world_id INTEGER,
  row_id INTEGER,
  col_id INTEGER,
  PRIMARY KEY (world_id, row_id, col_id),
  last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  locked BOOLEAN,
  data chunk NOT NULL
);