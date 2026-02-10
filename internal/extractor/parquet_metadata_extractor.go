package internal

import (
	"encoding/json"
	"fmt"

	"github.com/apache/arrow/go/v18/arrow/memory"
	"github.com/apache/arrow/go/v18/parquet"
	"github.com/apache/arrow/go/v18/parquet/file"
	"github.com/apache/arrow/go/v18/parquet/schema"
)

type RowGroupMetadata struct {
	RowCount              int
	ColumnCount           int
	TotalCompressedSize   *int64
	TotalUncompressedSize int
}

type BaseParquetMetadata struct {
	NumColumns   int32
	NumRowGroups int32
	RowGroups    []RowGroupMetadata
	Columns      []ColumnMetadata
}

type ColumnMetadata struct {
	Name         string
	PhysicalType parquet.Type
	LogicalType  schema.LogicalType
}

func ExtractParquetMetadata(filePath string) ([]byte, error) {
	// Implementation of ExtractParquetMetadata function

	// Read the Parquet file
	rdr, err := file.OpenParquetFile(filePath, false, file.WithReadProps(parquet.NewReaderProperties(memory.DefaultAllocator)))
	if err != nil {
		panic(err)
	}
	defer rdr.Close()

	metadata := rdr.MetaData()
	fmt.Print(len(metadata.RowGroups))
	basemeta := &BaseParquetMetadata{
		NumColumns:   int32(metadata.Schema.NumColumns()),
		NumRowGroups: int32(metadata.NumRows),
		RowGroups:    make([]RowGroupMetadata, len(metadata.RowGroups)),
		Columns:      make([]ColumnMetadata, metadata.Schema.NumColumns()),
	}

	fmt.Println(basemeta.NumRowGroups)
	// Extract row group metadata

	if basemeta.NumRowGroups > 1 {
		fmt.Println("More than one row group found")

		for i := 0; i < int(len(basemeta.RowGroups)); i++ {
			fmt.Println("RowGroup:", i)
			rowGroup := metadata.RowGroups[i]
			basemeta.RowGroups[i] = RowGroupMetadata{
				RowCount:              int(rowGroup.NumRows),
				ColumnCount:           len(rowGroup.GetColumns()),
				TotalCompressedSize:   rowGroup.TotalCompressedSize,
				TotalUncompressedSize: int(rowGroup.TotalByteSize),
			}
		}
	}

	// Extract column metadata
	for i := 0; i < metadata.Schema.NumColumns(); i++ {
		column := metadata.Schema.Column(i)
		basemeta.Columns[i] = ColumnMetadata{
			Name:         column.Name(),
			PhysicalType: column.PhysicalType(),
			LogicalType:  column.LogicalType(),
		}
	}

	meta_json, err := json.Marshal(basemeta)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(meta_json))
	return meta_json, nil
}
