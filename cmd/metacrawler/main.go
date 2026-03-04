package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ejakait/meta_crawler/crawler"
	logging "github.com/ejakait/meta_crawler/internal/logging"
	storage "github.com/ejakait/meta_crawler/internal/storage"
	config "github.com/ejakait/meta_crawler/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		slog.Error("Failed to load .env file", "error", err)
		os.Exit(1)
	}
	env := os.Getenv("env")
	GOOGLE_CREDENTIALS_JSON := os.Getenv("GOOGLE_CREDENTIALS_JSON")
	if env == "prod" {
		logging.InitLogger(true)
	} else {
		logging.InitLogger(false)
	}
	ctx := context.Background()

	slog.Info("Application Started")
	storage.InitDB(ctx, config.DbPath)

	config := make(map[string]string)

	config["project_id"] = "vbcdata"
	config["json_config"] = string(GOOGLE_CREDENTIALS_JSON)

	crawlerParams := &crawler.GoogleCrawler{
		ProjectID:  config["project_id"],
		JsonConfig: config["json_config"],
	}
	err = crawler.StartCrawl(ctx, crawlerParams)
	if err != nil {
		slog.Error("failed to start crawl", "error", err)
		return
	}
	// containerList, err := storage.ListGCSContainers("google")

	// slog.Info("containers found", "count", len(containerList))
	// if err != nil {
	// 	slog.Error("failed to list files", "error", err)

	// }

	// containerItems, err := storage.ListContainerItems("pharmaccess-cs")
	// if err != nil {
	// 	slog.Error("failed to list container items", "error", err)
	// }
	// for _, item := range containerItems {
	// 	size, _ := item.Size()
	// 	rawReader := &storage.StowReaderAt{Item: item}

	// 	bufferedReader := bufra.NewBufReaderAt(rawReader, 1024*1024)

	// 	pf, err := parquet.OpenFile(bufferedReader, size)

	// 	if err != nil {
	// 		slog.Error("failed to open parquet file", "error", err)
	// 		continue
	// 	}

	// 	slog.Info("file found", "name", item.Name(), "size", size, "rows", pf.NumRows(), "columns", len(pf.Schema().Fields()))
	// }
}
