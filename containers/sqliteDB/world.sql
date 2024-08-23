DROP TABLE IF EXISTS world_chunks;

CREATE TABLE world_chunk (
  world_id INTEGER,
  row_id INTEGER,
  col_id INTEGER,
  last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  locked BOOLEAN,
  data chunk NOT NULL,
  PRIMARY KEY (world_id, row_id, col_id)
);


CREATE TABLE world (
  world_id INTEGER PRIMARY KEY AUTOINCREMENT,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  num_rows INTEGER,
  num_cols INTEGER,
  chunkLen INTEGER
)