package storage

import (
	"encoding/json"
	"strings"
	"time"
)

// Do not change it
const maximum = 255

type ObjectType int

const (
	ObjectTypeIsDir ObjectType = iota
	ObjectTypeIsObject
)

type StorageType int

func (s StorageType) ToString() string {
	switch s {
	case StorageInLocal:
		return "local"
	case StorageOnAWS:
		return "aws"
	case StorageOnGCP:
		return "gcp"
	}
	return ""
}

func (s StorageType) Protocol() string {
	switch s {
	case StorageInLocal:
		return ""
	case StorageOnAWS:
		return "s3://"
	case StorageOnGCP:
		return "gs://"
	}
	return ""
}

const (
	StorageInLocal StorageType = iota
	StorageOnAWS
	StorageOnGCP
)

type Object struct {
	FileName string
	Size     int64
	ModTime  int64
	Sum      string
	Created  time.Time
	Updated  time.Time
}

// 100 ... 10000 => 100M ... 10000M
// attribute node
// gs://bucket/tenant/universe/bitmap/node/100M ... 10000M
// calculate node
// gs://bucket/tenant/universe/calculate/node/100M ... 10000M

type Storage interface {
	GetObject(node string) ([]byte, error)
	PutObject(node string, data []byte) error
	RemoveObject(node string) error
	RemoveDir(node string) error
	RemoveAll(node string) error
	CopyObject(from, to string) error
	MoveObject(from, to string) error
	IsExist(node string) bool
	ListObjects(dir string) ([]*Object, int64, error)
	ListChildObjects(dir string) ([]*Object, int64, error)
	ListDirs(dir string) ([]string, error)
	Download(from, to string) error
	Upload(from, to string) error

	PathJoin(items ...string) string
}

// NewStorage return a new Storage
func NewStorage(t StorageType, opts map[string]interface{}) Storage {
	switch t {
	case StorageOnGCP:
		return NewGCSStorage(opts)
	default:
		return NewFileStorage(nil)
	}
}

func NewStorageClient(ossPath, credentials string) Storage {
	if strings.HasPrefix(ossPath, StorageOnGCP.Protocol()) {
		return NewStorage(StorageOnGCP, generateOpts(credentials))
	}
	return NewStorage(StorageInLocal, nil)
}

func generateOpts(credentials string) map[string]interface{} {
	var opts map[string]interface{}
	if err := json.Unmarshal([]byte(credentials), &opts); err != nil {
		panic(err)
	}
	return opts
}
