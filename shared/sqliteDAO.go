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

func (d *SqliteDAO) SaveWorldChunk(worldId int, chunkRow int, chunkCol int, chunk []byte) error {
	db := d.db
	// SQL statement to insert or update a world chunk
	query := "INSERT INTO world_chunk (world_id, row_id, col_id, data, last_updated, locked) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, 0)"
	_, err := db.Exec(query, worldId, chunkRow, chunkCol, chunk)
	if err != nil {
		return fmt.Errorf("failed to save world chunk: %w", err)
	}

	return nil
}

func (d *SqliteDAO) FetchWorldChunk(worldId int, chunkRow int, chunkCol int) ([]byte, error) {
	db := d.db
	query := "SELECT data FROM world_chunk WHERE world_id = ? AND row_id = ? AND col_id = ? AND locked != 1"
	var chunk []byte
	err := db.QueryRow(query, worldId, chunkRow, chunkCol)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch world chunk: %w", err)
	}

	// Return the retrieved chunk
	return chunk, nil
}

func (d *SqliteDAO) CloseDb() error {
	return d.db.Close()
}
