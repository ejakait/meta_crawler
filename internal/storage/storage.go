package storage

import (
	"fmt"

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

// Dial dials stow storage.
// See stow.Dial for more information.
func Dial(kind string, config stow.Config) (stow.Location, error) {
	return stow.Dial(kind, config)
}

func ListGCSFiles(kind string) ([]string, error) {
	stowLoc, err := stow.Dial(stowgs.Kind, stow.ConfigMap{
		"bucket": "pharmaccess-cs",
	})

	fmt.Print(stowLoc)
	if err != nil {
		return nil, err
	}

	fileList := []string{}
	err = stow.WalkContainers(stowLoc, stow.NoPrefix, 100, func(c stow.Container, err error) error {
		if err != nil {
			return err
		}
		fmt.Print(c.Name())
		fileList = append(fileList, c.Name())

		return nil
	})
	if err != nil {
		return nil, err
	}

	// containers, _, err := location.Containers(stow.NoPrefix,"",int(0))
	// if err != nil {
	// 	return nil, err
	// }

	// var files []string
	// for _, container := range containers {
	// 	items, err := container.Items()
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	for _, item := range items {
	// 		files = append(files, item.ID())
	// 	}
	// }

	// return files, nil
	return fileList, nil
}
