package crawler

import (
	"context"
	"log/slog"

	stow "github.com/graymeta/stow"
	stowgs "github.com/graymeta/stow/google"
)

type Metadata struct {
	Name        string
	Size        int64
	ContentType string
}

type Crawler interface {
	InitCrawler(ctx context.Context, config stow.ConfigMap) (stow.Location, error)
	CrawlLocation(ctx context.Context, containerLocation stow.Location) ([]stow.Container, error)
	CrawlContainers(ctx context.Context, continer stow.Container) ([]stow.Item, error)
	GetMetaData(ctx context.Context, item stow.Item) (*Metadata, error)
}

type googleCrawler struct {
	projectID  string
	jsonConfig string
	bucketName string
	prefix     string
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
	Name        string
	Size        int64
	ContentType string
}

func (c *googleCrawler) InitCrawler(ctx context.Context) (stow.Location, error) {
	location, err := stow.Dial(stowgs.Kind, stow.ConfigMap{
		stowgs.ConfigJSON:      c.jsonConfig,
		stowgs.ConfigProjectId: c.projectID,
	})
	if err != nil {
		slog.Error("failed to dial google storage", "error", err)
		return nil, err
	}
	defer func() {
		slog.Info("Connected to GCS Successfully", "project_id", c.projectID)
	}()
	return location, nil
}

func (c *googleCrawler) CrawlLocation(ctx context.Context, containerLocation stow.Location) ([]stow.Container, error) {
	return nil, nil
}

func (c *googleCrawler) CrawlContainers(ctx context.Context, container stow.Container) ([]stow.Item, error) {
	return nil, nil
}

func (c *googleCrawler) GetMetaData(ctx context.Context, item stow.Item) (*Metadata, error) {
	return nil, nil
}

func StartCrawl(ctx context.Context, config map[string]string) error {
	c := &googleCrawler{
		projectID:  config["project_id"],
		jsonConfig: config["json_config"],
	}
	location, err := c.InitCrawler(ctx)
	if err != nil {
		return err
	}
	defer location.Close()

	return nil
}
