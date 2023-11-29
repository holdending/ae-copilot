package storage

import (
	"errors"
	"fmt"
	usr "os/user"
	"path"
	"path/filepath"
	"strings"
)

// IllegalPath ...
var IllegalPath = errors.New("illegal file path")
var ErrCodeNoSuchKey = errors.New("no such key")

const protocolFlag = "://"
const slash = "/"

type objOpt struct {
	Bucket string
	Key    string
	Prefix string
}

func parseObj(node string) (*objOpt, error) {
	prefixIdx := strings.Index(node, protocolFlag)
	if prefixIdx > -1 {
		node = strings.Split(node, protocolFlag)[1]
	}
	nodes := strings.Split(node, slash)
	if len(nodes) < 2 {
		return nil, fmt.Errorf("%v %s", IllegalPath, node)
	}
	return &objOpt{
		Bucket: nodes[0],
		Key:    path.Join(nodes[1:]...),
		Prefix: path.Join(nodes[1:len(nodes)-1]...) + slash,
	}, nil
}

func ParseObj(node string) (*objOpt, error) {
	return parseObj(node)
}

func appendPathSuffix(dir string) string {
	if dir == "" {
		return ""
	}
	if strings.HasSuffix(dir, slash) {
		return dir
	}
	return dir + slash
}

func ObjectsToStrings(objs []*Object) []string {
	files := make([]string, len(objs))
	for k, v := range objs {
		files[k] = v.FileName
	}
	return files
}

func FilterFilesByPrefix(objs []*Object, prefix string) ([]*Object, int64) {
	var files []*Object
	var size int64
	for _, v := range objs {
		if _, filename := filepath.Split(v.FileName); strings.HasPrefix(strings.ToUpper(filename), strings.ToUpper(prefix)) {
			files = append(files, v)
			size += v.Size
		}
	}
	return files, size
}

// ExpandUserDir returns the argument with an initial component of ~
func ExpandUserDir(path string) string {
	usr, err := usr.Current()
	if err != nil {
		return path
	}
	dir := usr.HomeDir
	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}
	return path
}
