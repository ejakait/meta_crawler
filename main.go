package main

import (
	extractor "github.com/ejakait/meta_crawler/internal/extractor"
)

func parquetReader() {
}

func main() {
	// ctx := context.Background()

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

	extractor.ExtractParquetMetadata("25-csv-20250827124850_dat.parquet")
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
