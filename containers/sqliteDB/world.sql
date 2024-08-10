DROP TABLE IF EXISTS world_chunks;

CREATE TABLE world_chunks (
  id INTEGER PRIMARY KEY,
  row_id INTEGER,
  col_id INTEGER,
  data chunk NOT NULL
);