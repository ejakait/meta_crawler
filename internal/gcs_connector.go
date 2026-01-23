package internal

import "time"

type Bucket struct {
	ObjAttributes
}
type ObjAttributes struct {
	name         string
	location     Region
	createdAt    time.Time
	storageClass string
}
type BlobObj struct {
	ObjAttributes
	GcsUri          string
	SizeBytes       int
	RowCount        int
	PartitionValues string
}

type Dataset struct {
	Id          int
	Name        string
	Description string
	Owner       string
}
type Tags struct {
}
type Region struct {
	name string
}
type GCSConnector struct {
	credentialsPath string
	project         string
	Region
	buckets []Bucket
}

type MetaData struct {
}
