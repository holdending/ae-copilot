package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var mockFiles = [12]string{
	"%s/audience_%d_tenant/bitmap",
	"%s/audience_%d_tenant/metadata.json",
	"%s/audience_%d_tenant/tenant.rebuild",
	"%s/audience_%d_tenant/_SUCCESS",
	"%s/audience_%d_tenant/PART-001",
	"%s/audience_%d_tenant/PART-002",
	"%s/audience_%d_tenant/PART-003",
	"%s/audience_%d_tenant/PART-004",
	"%s/audience_%d_tenant/PART-005",
	"%s/audience_%d_tenant/attrs/PART-001",
	"%s/audience_%d_tenant/attrs/PART-002",
	"%s/audience_%d_tenant/attrs/PART-003",
}

const mockGCSBucket = "gs://bitmap_mock"
const mockContent = "test data"

func mockTestFiles(storageType StorageType, bucket string, t *testing.T) (Storage, [12]string) {
	var files [12]string
	timestamp := time.Now().UnixNano()
	client := NewStorage(storageType, nil)
	for k, v := range mockFiles {
		files[k] = fmt.Sprintf(v, bucket, timestamp)

		if err := client.PutObject(files[k], []byte(mockContent)); err != nil {
			t.Fatal(err)
		}

	}
	return client, files
}

func TestStorage(t *testing.T) {
	tempDir, err := ioutil.TempDir("./", "fileStorage")
	if err != nil {
		t.Fail()
	}
	defer os.RemoveAll(tempDir)

	client, files := mockTestFiles(StorageInLocal, tempDir, t)
	doStorageTestCases(client, files, tempDir, t)
}

func doStorageTestCases(client Storage, files [12]string, tempDir string, t *testing.T) {
	obj, err := parseObj(files[0])
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	prefix := client.PathJoin(obj.Bucket, obj.Prefix)
	defer client.RemoveAll(prefix)

	assert.Nil(t, client.CopyObject(files[0], files[3]))

	assert.Nil(t, client.RemoveObject(files[3]))

	assert.Equal(t, false, client.IsExist(files[3]))

	assert.Nil(t, client.MoveObject(files[3], files[0]))

	assert.Equal(t, true, client.IsExist(files[0]))

	objs, _, _ := client.ListObjects(prefix)
	assert.Equal(t, 11, len(objs))

	objs, _, _ = client.ListChildObjects(prefix)
	assert.Equal(t, 8, len(objs))

	dirs, _ := client.ListDirs(prefix)
	assert.Equal(t, 1, len(dirs))

	localClient := NewFileStorage(nil)
	localDir := localClient.PathJoin(tempDir, "bitmap")
	assert.Nil(t, client.Download(files[0], localDir))

	client.Upload(localDir, files[0])

	assert.Nil(t, client.RemoveAll(prefix))

	//assert.Equal(t,false,client.IsExist(files[0]))
}
