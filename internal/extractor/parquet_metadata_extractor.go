package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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

func ExtractParquetMetadata(ctx context.Context, logger *slog.Logger, filePath string) ([]byte, error) {

	// Implementation of ExtractParquetMetadata function
	extractorLogger := logger.With(slog.String("file", filePath))
	// Read the Parquet file
	extractorLogger.InfoContext(ctx, "Opening file", "status", "initated")
	rdr, err := file.OpenParquetFile(filePath, false, file.WithReadProps(parquet.NewReaderProperties(memory.DefaultAllocator)))
	if err != nil {
		extractorLogger.ErrorContext(ctx, "Failed to open file", "error", err)
		return nil, err
	}
	defer func() {
		rdr.Close()
		extractorLogger.InfoContext(ctx, "Closing file", "status", "completed")
	}()

	metadata := rdr.MetaData()

	extractorLogger.InfoContext(ctx, "Extracting metadata")
	basemeta := &BaseParquetMetadata{
		NumColumns:   int32(metadata.Schema.NumColumns()),
		NumRowGroups: int32(metadata.NumRows),
		RowGroups:    make([]RowGroupMetadata, len(metadata.RowGroups)),
		Columns:      make([]ColumnMetadata, metadata.Schema.NumColumns()),
	}

	extractorLogger.InfoContext(ctx, "Extracting row group metadata")
	// Extract row group metadata
	if basemeta.NumRowGroups > 1 {
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
	extractorLogger.InfoContext(ctx, "Extracted row group metadata", "status", "completed", "row_groups", len(basemeta.RowGroups))

	extractorLogger.InfoContext(ctx, "Extracting column metadata")
	// Extract column metadata
	for i := 0; i < metadata.Schema.NumColumns(); i++ {
		column := metadata.Schema.Column(i)
		basemeta.Columns[i] = ColumnMetadata{
			Name:         column.Name(),
			PhysicalType: column.PhysicalType(),
			LogicalType:  column.LogicalType(),
		}
	}
	extractorLogger.InfoContext(ctx, "Extracted column metadata", "status", "completed", "columns", len(basemeta.Columns))

	meta_json, err := json.Marshal(basemeta)
	if err != nil {
		extractorLogger.ErrorContext(ctx, "Failed to marshal metadata", "error", err)
	}

	return meta_json, nil
}
