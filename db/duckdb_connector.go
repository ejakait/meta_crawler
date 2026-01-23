package db

import (
	"context"

	_ "github.com/duckdb/duckdb-go/v2"
)

type FileDBConn struct {
	Path string
}

func createTables(dbPath *FileDBConn) {
	ctx := context.Background()
	connector, err := duckdb.NewConnector(dbPath.Path)
	if err != nil {
		panic(err)
	}
}
