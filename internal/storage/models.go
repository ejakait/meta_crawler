package storage

import parquet "github.com/parquet-go/parquet-go"

type ContainerMetadata struct {
	Name        string
	ObjectCount int64
	Description string
	Owner       string
	Tags        string
	Location    string
	CreatedAt   string
	UpdatedAt   string
}

type ItemMetadata struct {
	Name    string
	Size    int64
	NumRows int64
	Schema  []parquet.Field
}
