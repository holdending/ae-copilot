package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CopyObject(t *testing.T) {
	tempDir, err := ioutil.TempDir("./", "fileStorage")
	if err != nil {
		t.Fail()
	}
	defer os.RemoveAll(tempDir)

	local := NewFileStorage(nil)
	local.PutObject(local.PathJoin(tempDir, "f0"), []byte("a1"))
	local.CopyObject(local.PathJoin(tempDir, "f0"), local.PathJoin(tempDir, "f01"))
	bs, err := local.GetObject(local.PathJoin(tempDir, "f01"))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, string(bs) == "a1")

	local.PutObject(local.PathJoin(tempDir, "f1", "a1"), []byte("a1"))
	local.PutObject(local.PathJoin(tempDir, "f1", "f2", "a2"), []byte("a2"))

	local.CopyObject(local.PathJoin(tempDir, "f1"), local.PathJoin(tempDir, "t1"))

	bs, err = local.GetObject(local.PathJoin(tempDir, "t1", "a1"))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, string(bs) == "a1")

	bs, err = local.GetObject(local.PathJoin(tempDir, "t1", "f2", "a2"))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, string(bs) == "a2")

	local.PutObject(local.PathJoin(tempDir, "f1", "f3", "a3"), []byte("a3"))
	local.CopyObject(local.PathJoin(tempDir, "f1"), local.PathJoin(tempDir, "t1"))

	bs, err = local.GetObject(local.PathJoin(tempDir, "t1", "a1"))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, string(bs) == "a1")

	bs, err = local.GetObject(local.PathJoin(tempDir, "t1", "f2", "a2"))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, string(bs) == "a2")

	bs, err = local.GetObject(local.PathJoin(tempDir, "t1", "f3", "a3"))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, string(bs) == "a3")
}

func Test_CopyObject2(t *testing.T) {
	tempDir, err := ioutil.TempDir("./", "fileStorage")
	if err != nil {
		t.Fail()
	}
	//defer os.RemoveAll(tempDir)

	local := NewFileStorage(nil)
	//local.PutObject(local.PathJoin(tempDir, "d01", "f0"), []byte("a1"))
	local.PutObject(local.PathJoin(tempDir, "d0", "f0"), []byte("a1"))
	local.CopyObject(local.PathJoin(tempDir, "d0"), local.PathJoin(tempDir, "d01"))

	bs, err := local.GetObject(local.PathJoin(tempDir, "d01", "f0"))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, string(bs) == "a1")

}
