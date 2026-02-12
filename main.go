package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	storage "github.com/ejakait/meta_crawler/internal/storage"

	bufra "github.com/avvmoto/buf-readerat"
	parquet "github.com/parquet-go/parquet-go"
)

func main() {
	ctx := context.Background()
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("Application Started")
	requestID := "req-" + time.Now().Format("20060102150405")
	requestLogger := logger.With(slog.String("requestID", requestID))
	storage.InitDB(ctx, requestLogger, "meta.db")

	containerList, err := storage.ListGCSContainers("google", requestLogger)

	logger.Info("containers found", "count", len(containerList))
	if err != nil {
		logger.Error("failed to list files", "error", err)

	}

	containerItems, err := storage.ListContainerItems(requestLogger, "pharmaccess-cs")
	if err != nil {
		logger.Error("failed to list container items", "error", err)
	}
	for _, item := range containerItems {
		size, _ := item.Size()
		rawReader := &storage.StowReaderAt{Item: item}

		bufferedReader := bufra.NewBufReaderAt(rawReader, 1024*1024)

		pf, err := parquet.OpenFile(bufferedReader, size)

		if err != nil {
			logger.Error("failed to open parquet file", "error", err)
			continue
		}

		logger.Info("file found", "name", item.Name(), "size", size, "rows", pf.NumRows(), "columns", len(pf.Schema().Fields()))
	}
}
