package shared

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type SqliteDAO struct {
	db   *sql.DB
	name string
}

const PATH = "containers/sqliteDB/"

func NewSqliteDAO(name string) *SqliteDAO {
	return &SqliteDAO{
		name: name,
	}
}

func (d *SqliteDAO) OpenDb(dataSourceName string) error {
	var err error
	path := PATH + dataSourceName
	d.db, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	return nil
}

func (d *SqliteDAO) InsertNewWorld(rows int, cols int, chunkLen int) (int64, error) {
	db := d.db

	query := `INSERT INTO world ("date_created", "rows", "columns", "chunk_length") VALUES (CURRENT_TIMESTAMP, ?, ?, ?)`

	result, err := db.Exec(query, rows, cols, chunkLen)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *SqliteDAO) SaveWorldChunk(worldId int, chunkRow int, chunkCol int, chunk []byte, lockChunk bool) error {
	db := d.db

	// SQL statement to insert or update a world chunk
	query := `
	INSERT INTO world_chunk (world_id, row_id, col_id, data, last_updated, locked)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, 0)
	ON CONFLICT(world_id, row_id, col_id)
	DO UPDATE SET
		data = excluded.data,
		last_updated = CURRENT_TIMESTAMP,
		locked = excluded.locked`

	_, err := db.Exec(query, worldId, chunkRow, chunkCol, chunk, lockChunk)
	if err != nil {
		return fmt.Errorf("failed to save world chunk: %w", err)
	}

	return nil
}

func (d *SqliteDAO) FetchWorldChunk(worldId int, chunkRow int, chunkCol int) ([]byte, error) {
	db := d.db
	var locked bool
	query := "SELECT data, locked FROM world_chunk WHERE world_id = ? AND row_id = ? AND col_id = ?"
	var chunk []byte
	err := db.QueryRow(query, worldId, chunkRow, chunkCol).Scan(&chunk, &locked)
	if err != nil {
		return nil, err
	}

	if locked {
		err = fmt.Errorf("fetched row is already claimed")
	}

	return chunk, err
}

func (d *SqliteDAO) CloseDb() error {
	return d.db.Close()
}
