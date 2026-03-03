package crawler

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	stow "github.com/graymeta/stow"
	stowgs "github.com/graymeta/stow/google"
	parquet "github.com/parquet-go/parquet-go"
)

type Metadata struct {
	Name        string
	Size        int64
	ContentType string
}

type Crawler interface {
	InitCrawler(ctx context.Context, config stow.ConfigMap) (stow.Location, error)
	CrawlContainers(ctx context.Context, continer stow.Container) ([]stow.Item, error)
	GetMetaData(ctx context.Context, item stow.Item) (*Metadata, error)
}

type GoogleCrawler struct {
	ProjectID  string
	JsonConfig string
	BucketName string
	Prefix     string
}

type LocationMetadata struct {
	Name        string
	Size        int64
	ContentType string
}

type ContainerMetadata struct {
	Name        string
	Size        int64
	ContentType string
}

type ItemMetadata struct {
	Name    string
	Size    int64
	NumRows int64
	Schema  []parquet.Field
}

func (c *GoogleCrawler) InitCrawler(ctx context.Context) (stow.Location, error) {
	location, err := stow.Dial(stowgs.Kind, stow.ConfigMap{
		stowgs.ConfigJSON:      c.JsonConfig,
		stowgs.ConfigProjectId: c.ProjectID,
	})
	if err != nil {
		slog.Error("failed to dial google storage", "error", err)
		return nil, err
	}
	defer func() {
		slog.Info("Connected to GCS Successfully", "project_id", c.ProjectID)
	}()
	return location, nil
}

func (c *GoogleCrawler) CrawlContainers(ctx context.Context, container stow.Container) ([]stow.Item, error) {
	items := []stow.Item{}
	err := stow.Walk(container, stow.NoPrefix, 100, func(item stow.Item, err error) error {
		if err != nil {
			return err
		}
		items = append(items, item)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *GoogleCrawler) GetMetaData(ctx context.Context, item stow.Item) (*Metadata, error) {

	itemSize, err := item.Size()
	if err != nil {
		return nil, err
	}

	metadata := &Metadata{
		Name:        item.Name(),
		Size:        itemSize,
		ContentType: item.ID(),
	}

	return metadata, nil
}

type stowReaderAt struct {
	item      stow.Item
	ReadCount int64
	BytesRead int64
}

func (r *stowReaderAt) ReadAt(p []byte, off int64) (n int, err error) {

	atomic.AddInt64(&r.ReadCount, 1)
	atomic.AddInt64(&r.BytesRead, int64(len(p)))

	slog.Debug("ReadAt", slog.Int64("read_count", off), "Length", slog.Int64("bytes_read", int64(len(p))), "TotalReads", atomic.LoadInt64(&r.ReadCount))

	rc, err := r.item.Open()

	if err != nil {
		return 0, err
	}
	defer rc.Close()

	if seeker, ok := rc.(io.ReadSeeker); ok {
		_, err := seeker.Seek(off, io.SeekStart)
		if err != nil {
			return 0, err
		}
		return io.ReadFull(seeker, p)
	}

	// Fall back to ReadAt
	_, err = io.CopyN(io.Discard, rc, off)
	if err != nil {
		return 0, err
	}
	return io.ReadFull(rc, p)
}

func (c *GoogleCrawler) GetParquetMetadata(ctx context.Context, item stow.Item) (*ItemMetadata, error) {

	size, err := item.Size()
	if err != nil {
		return nil, err
	}
	readerAt := &stowReaderAt{item: item}

	pf, err := parquet.OpenFile(readerAt, size)
	if err != nil {
		return nil, err
	}

	return &ItemMetadata{
		Name:    item.Name(),
		Size:    size,
		NumRows: pf.NumRows(),
		Schema:  pf.Schema().Fields(),
	}, nil

}

func StartCrawl(ctx context.Context, crawlerParams *GoogleCrawler) error {

	//Time the function
	start := time.Now()
	defer func() {
		slog.Info("Crawl completed", "duration", time.Since(start))
	}()
	wg := sync.WaitGroup{}
	c := &GoogleCrawler{
		ProjectID:  crawlerParams.ProjectID,
		JsonConfig: crawlerParams.JsonConfig,
	}
	location, err := c.InitCrawler(ctx)
	if err != nil {
		return err
	}
	defer location.Close()
	containerList := []stow.Container{}
	slog.Info("Crawling location")
	err = stow.WalkContainers(location, stow.NoPrefix, 100, func(c stow.Container, err error) error {
		if err != nil {
			return err
		}
		slog.Info("container name", "name", c.Name())
		containerList = append(containerList, c)
		return nil
	})
	if err != nil {
		return err
	}
	for _, container := range containerList {

		items, err := c.CrawlContainers(ctx, container)
		if err != nil {
			return err
		}
		// Get only parquet files
		var parquetItems []stow.Item
		for _, item := range items {
			if strings.HasSuffix(item.Name(), ".parquet") {
				parquetItems = append(parquetItems, item)
			}
		}
		items = parquetItems
		for _, item := range items {
			wg.Add(1)

			go func() {
				defer wg.Done()
				metadata, err := c.GetParquetMetadata(ctx, item)
				if err != nil {
					slog.Error("GetParquetMetadata error", "error", err)
					return
				}
				slog.Info("item name", "name", item.Name(), "container", container.Name(), "size", metadata.Size, "rows", metadata.NumRows, "schema", metadata.Schema)
			}()
		}
		wg.Wait()

	}

	return nil
}
