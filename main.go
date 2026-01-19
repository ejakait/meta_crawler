package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	bkt := client.Bucket("pharmaccess-cs")

	attrs, err := bkt.Attrs(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Printf("bucket %s, created at %s, is located in %s with storage class %s\n",
		attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)

	it := bkt.Objects(ctx, nil)
	for {
		attrs, err := it.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", attrs.Created)
	}
}
