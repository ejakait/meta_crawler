package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/apache/arrow/go/v18/arrow/memory"
	parquet "github.com/apache/arrow/go/v18/parquet"
	"github.com/apache/arrow/go/v18/parquet/file"
)

func parquetReader() {

}
func main() {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}
	defer client.Close()

	bkt := client.Bucket("pharmaccess-cs")

	attrs, err := bkt.Attrs(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Printf("bucket %s, created at %s, is located in %s with storage class %s\n",
		attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)
	var object string = "28_dat.parquet"
	it := bkt.Object(object)

	reader, err := it.NewReader(ctx)

	if err != nil {
		panic(err)
	}

	defer reader.Close()

	rdr, err := file.OpenParquetFile(reader, false, file.WithReadProps(parquet.NewReaderProperties(memory.DefaultAllocator)))

	if err != nil {
		panic(err)
	}

	defer rdr.Close()

	fmt.Println(rdr.MetaData())
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
