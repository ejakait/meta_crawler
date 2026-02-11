package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	storage "github.com/ejakait/meta_crawler/internal/storage"
)

type requestKey string

const requestIDKey requestKey = "requestID"

func main() {
	ctx := context.Background()
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("Application Started")
	// client, err := storage.NewClient(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.Close()

	// bkt := client.Bucket("pharmaccess-cs")

	// attrs, err := bkt.Attrs(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("bucket %s, created at %s, is located in %s with storage class %s\n",
	// 	attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)
	// var object string = "28_dat.parquet"
	// it := bkt.Object(object)

	// reader, err := it.NewReader(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// defer reader.Close()

	// rdr, err := file.OpenParquetFile(reader, false, file.WithReadProps(parquet.NewReaderProperties(memory.DefaultAllocator)))
	// if err != nil {
	// 	panic(err)
	// }

	// defer rdr.Close()

	// fmt.Println(rdr.MetaData())
	requestID := "req-" + time.Now().Format("20060102150405")
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	requestLogger := logger.With(slog.String("requestID", requestID))
	storage.InitDB(ctx, requestLogger, "meta.db")

	files, err := storage.ListGCSFiles("Google Cloud Storage")
	if err != nil {
		fmt.Errorf("failed to list files: %v", err)

	}
	if len(files) == 0 {
		fmt.Errorf("no files found")
	}
	fmt.Print(files)
	// extractor.ExtractParquetMetadata(ctx, requestLogger, "25-csv-20250827124850_dat.parquet")

	// for {
	// 	attrs, err := it.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("%s", attrs.Created)
	// 	reader, err := it.NewReader(ctx)
	// }
}
