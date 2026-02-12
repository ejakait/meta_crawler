package storage

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/graymeta/stow"
	// support Azure storage
	_ "github.com/graymeta/stow/azure"
	// support Google storage
	stowgs "github.com/graymeta/stow/google"
	// support local storage
	_ "github.com/graymeta/stow/local"
	// support swift storage
	_ "github.com/graymeta/stow/swift"
	// support s3 storage
	_ "github.com/graymeta/stow/s3"
	// support oracle storage
	_ "github.com/graymeta/stow/oracle"
)

type StowReaderAt struct {
	Item stow.Item
}

func (r *StowReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	ranger, ok := r.Item.(stow.ItemRanger)

	if !ok {
		return 0, fmt.Errorf("item does not support range")
	}

	end := uint64(off) + uint64(len(p)) - 1
	rc, err := ranger.OpenRange(uint64(off), end)
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	return io.ReadFull(rc, p)
}

func connectGCS() (stow.Location, error) {

	var gcsConfigJson string = ""

	location, err := stow.Dial(stowgs.Kind, stow.ConfigMap{
		stowgs.ConfigProjectId: "vbcdata",
		stowgs.ConfigJSON:      gcsConfigJson,
	})
	if err != nil {
		return nil, err
	}
	return location, nil
}

// func ListGCSContainers(kind string, logger *slog.Logger) ([]string, error) {

// 	location, err := connectGCS()
// 	if err != nil {
// 		logger.Info("failed to connect to GCS", "error", err)
// 		return nil, err
// 	}

// 	fileList := []string{}
// 	err = stow.WalkContainers(location, stow.NoPrefix, 100, func(c stow.Container, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		logger.Info("container name", "name", c.Name())
// 		fileList = append(fileList, c.Name())
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return fileList, nil
// }

func ListGCSContainers(kind string, logger *slog.Logger) ([]stow.Container, error) {

	location, err := connectGCS()
	if err != nil {
		logger.Info("failed to connect to GCS", "error", err)
		return nil, err
	}

	containerList := []stow.Container{}
	err = stow.WalkContainers(location, stow.NoPrefix, 100, func(c stow.Container, err error) error {
		if err != nil {
			return err
		}
		logger.Info("container name", "name", c.Name())
		containerList = append(containerList, c)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return containerList, nil
}

func ListContainerItems(logger *slog.Logger, locContainer string) ([]stow.Item, error) {
	loc, err := connectGCS()
	if err != nil {
		logger.Info("failed to connect to GCS", "error", err)
		return nil, err
	}
	container, _ := loc.Container(locContainer)
	fileList := []stow.Item{}
	stow.Walk(container, stow.NoPrefix, 100, func(item stow.Item, err error) error {
		if err != nil {
			return err
		}
		logger.Info("item name", "name", item.Name())
		fileList = append(fileList, item)
		return nil
	})
	return fileList, nil
}
