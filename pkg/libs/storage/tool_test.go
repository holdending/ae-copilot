package storage

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseObj(t *testing.T) {
	node := "s3://bucket/data/audience_1001_comcast/object"
	obj, err := parseObj(node)
	assert.Nil(t, err)
	assert.Equal(t, "data/audience_1001_comcast/", obj.Prefix)
	assert.Equal(t, "data/audience_1001_comcast/object", obj.Key)
	assert.Equal(t, "bucket", obj.Bucket)
}

func TestAppendPathSuffix(t *testing.T) {
	node := "s3://bucket/data"
	assert.Equal(t, node+"/", appendPathSuffix(node))
}

func TestObjectsToStrings(t *testing.T) {
	objs := []*Object{
		&Object{
			FileName: "1.txt",
			Size:     10,
		},
		&Object{
			FileName: "2.txt",
			Size:     10,
		},
	}
	files := ObjectsToStrings(objs)
	assert.Equal(t, "1.txt,2.txt", strings.Join(files, ","))
}

func TestFilterFilesByPrefix(t *testing.T) {
	objs := []*Object{
		&Object{
			FileName: "part-001.txt",
			Size:     10,
		},
		&Object{
			FileName: "PART-002.txt",
			Size:     10,
		},
	}
	partObjs, size := FilterFilesByPrefix(objs, "part-")
	assert.Equal(t, int64(20), size)
	assert.Equal(t, 2, len(partObjs))
}
