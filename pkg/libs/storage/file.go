package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// FileStorage is local storage
type FileStorage struct {
	protocol string
}

// NewFileStorage return a new file storage client
func NewFileStorage(opts map[string]interface{}) *FileStorage {
	return new(FileStorage)
}

func (f *FileStorage) PathJoin(items ...string) string {
	return filepath.Join(items...)
}

// GetObject return a file by node.
func (f *FileStorage) GetObject(node string) ([]byte, error) {
	if _, err := os.Stat(node); err != nil {
		return nil, ErrCodeNoSuchKey
	}
	file, err := os.Open(node)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// PutObject save a file via node.
func (f *FileStorage) PutObject(node string, data []byte) error {
	if err := mkDirs(node); err != nil {
		return err
	}
	return ioutil.WriteFile(node, data, 0750)
}

// RemoveObject remove a file via node.
func (f *FileStorage) RemoveObject(node string) error {
	return os.Remove(node)
}

// RemoveDir remove a folder.
func (f *FileStorage) RemoveDir(path string) error {
	return f.RemoveObject(path)
}

// RemoveAll remove a folder via path
func (f *FileStorage) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// CopyObject backup this node file
func (f *FileStorage) CopyObject(from, to string) error {
	var cmd string
	if isDir(from) {
		to = to + "/"
		cmd = "cp -rf %s/* %s"
	} else {
		cmd = "cp -f %s %s"
	}
	if err := mkDirs(to); err != nil {
		return err
	}
	_, err := runCMD(fmt.Sprintf(cmd, from, to))
	return err
}

// MoveObject rename this object
func (f *FileStorage) MoveObject(from, to string) error {
	if isDir(from) {
		from += "/*"
		to += "/"
	}
	if err := mkDirs(to); err != nil {
		return err
	}
	_, err := runCMD(fmt.Sprintf("mv -f %s %s", from, to))
	return err
}

// IsExist return false if node doesn't exist
func (f *FileStorage) IsExist(node string) bool {
	_, err := os.Stat(node)
	return err == nil || os.IsExist(err)
}

// ListObjects return all files via prefix dir
func (f *FileStorage) ListObjects(dir string) ([]*Object, int64, error) {
	return f.listByPrefix(dir, "")
}

func (f *FileStorage) ListChildObjects(dir string) ([]*Object, int64, error) {
	return f.listByPrefix(dir, "/")
}

// ListDirs return all dirs via prefix dir
func (f *FileStorage) ListDirs(dir string) ([]string, error) {
	var files []string
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}
	for _, fi := range dirs {
		if fi.IsDir() == true {
			files = append(files, filepath.Join(dir, fi.Name()))
		}
	}
	return files, nil
}

func (f *FileStorage) listByPrefix(prefix, delim string) ([]*Object, int64, error) {
	objs := make([]*Object, 0, 10000)
	var size int64
	if delim == "/" {
		dirs, err := ioutil.ReadDir(prefix)
		if err != nil {
			return nil, 0, err
		}
		for _, f := range dirs {
			if f.IsDir() == false {
				objs = append(objs, &Object{
					FileName: filepath.Join(prefix, f.Name()),
					Size:     f.Size(),
					ModTime:  f.ModTime().Unix(),
				})
				size += f.Size()
			}
		}
	} else {
		err := filepath.Walk(prefix, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() == false {
				objs = append(objs, &Object{
					FileName: path,
					Size:     f.Size(),
					ModTime:  f.ModTime().Unix(),
				})
				size += f.Size()
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
	}
	return objs, size, nil
}

// Download download file to local
func (f *FileStorage) Download(from, to string) error {
	return f.CopyObject(from, to)
}

// Upload put file to remote
func (f *FileStorage) Upload(from, to string) error {
	return f.CopyObject(from, to)
}

func runCMD(shell string) (string, error) {
	cmd := exec.Command("sh", "-c", shell)
	result, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("run cmd %s: %v", shell, err)
	}
	return string(result), nil
}

func mkDirs(node string) error {
	dir := filepath.Dir(node)
	if !isExist(dir) {
		os.MkdirAll(dir, 0750)
	}
	return nil
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func isExist(node string) bool {
	_, err := os.Stat(node)
	return err == nil || os.IsExist(err)
}
