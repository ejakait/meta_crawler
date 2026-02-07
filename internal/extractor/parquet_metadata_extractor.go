package internal

import (
	"fmt"

	"github.com/apache/arrow/go/v18/arrow/memory"
	"github.com/apache/arrow/go/v18/parquet"
	"github.com/apache/arrow/go/v18/parquet/file"
)

type RowGroupMetadata struct {
	RowCount              int64
	ColumnCount           int32
	TotalCompressedSize   int64
	TotalUncompressedSize int64
}

type BaseParquetMetadata struct {
	NumColumns   int32
	NumRowGroups int32
	RowGroups    []RowGroupMetadata
	Columns      []ColumnMetadata
}

type ColumnMetadata struct {
	Name        string
	LogicalType parquet.Type
}

func ExtractParquetMetadata(filePath string) {
	// Implementation of ExtractParquetMetadata function

	// Read the Parquet file
	rdr, err := file.OpenParquetFile(filePath, false, file.WithReadProps(parquet.NewReaderProperties(memory.DefaultAllocator)))
	if err != nil {
		panic(err)
	}
	defer rdr.Close()

	metadata := rdr.MetaData()
	basemeta := &BaseParquetMetadata{
		NumColumns:   int32(metadata.Schema.NumColumns()),
		NumRowGroups: int32(metadata.NumRows),
		RowGroups:    make([]RowGroupMetadata, len(metadata.RowGroups)),
		Columns:      make([]ColumnMetadata, metadata.Schema.NumColumns()),
	}

	// Extract row group metadata
	// for i := 0; i < int(basemeta.NumRowGroups); i++ {
	// 	rowGroup := metadata.RowGroups[i]
	// 	basemeta.RowGroups[i] = RowGroupMetadata{
	// 		RowCount:              rowGroup.NumRows,
	// 		ColumnCount:           rowGroup.Columns,
	// 		TotalCompressedSize:   rowGroup.TotalByteSize,
	// 		TotalUncompressedSize: rowGroup.TotalUncompressedSize,
	// 	}
	// }

	// // Extract column metadata
	// for i := 0; i < metadata.NumColumns; i++ {
	// 	column := metadata.Columns[i]
	// 	metadata.Columns[i] = ColumnMetadata{
	// 		Name:        column.Name,
	// 		LogicalType: column.LogicalType,
	// 	}
	// }

	// return metadata, nil
	fmt.Println(basemeta)
	return
}
