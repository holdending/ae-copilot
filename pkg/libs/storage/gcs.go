package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	gs "cloud.google.com/go/storage"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCSStorage is remote storage by gcs
type GCSStorage struct {
	ProjectID string
	Token     string
	protocol  string
}

// NewGCSStorage return a new GCS storage client
func NewGCSStorage(opts map[string]interface{}) *GCSStorage {
	gcpStorage := new(GCSStorage)
	if ProjectID, ok := opts["ProjectID"]; ok {
		gcpStorage.ProjectID = ProjectID.(string)
	}
	if Token, ok := opts["SecretAccessKey"]; ok {
		gcpStorage.Token = Token.(string)
	}
	gcpStorage.protocol = StorageOnGCP.Protocol()
	return gcpStorage
}

func (g *GCSStorage) PathJoin(items ...string) string {
	if len(items) <= 0 {
		return ""
	}
	items[0] = strings.Replace(items[0], g.protocol, "", 1)
	return g.protocol + path.Join(items...)
}

// GetObject return a data object by node.
func (g *GCSStorage) GetObject(node string) ([]byte, error) {
	opts, err := parseObj(node)
	if err != nil {
		return nil, err
	}
	data, err := g.read(opts.Bucket, opts.Key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// PutObject save a data object via node.
func (g *GCSStorage) PutObject(node string, data []byte) error {
	opts, err := parseObj(node)
	if err != nil {
		return err
	}
	return g.write(data, opts.Bucket, opts.Key)
}

// RemoveObject remove a data object via node.
func (g *GCSStorage) RemoveObject(node string) error {
	opts, err := parseObj(node)
	if err != nil {
		return err
	}
	return g.delete(opts.Bucket, opts.Key)
}

// RemoveDir remove a folder.
func (g *GCSStorage) RemoveDir(node string) error {
	opts, err := parseObj(node)
	if err != nil {
		return err
	}
	return g.delete(opts.Bucket, opts.Prefix)
}

// RemoveAll remove a folder via path
func (g *GCSStorage) RemoveAll(path string) error {
	// todo: detail logic
	return nil
}

// CopyObject backup this object
func (g *GCSStorage) CopyObject(src, dst string) error {
	srcOpts, err := parseObj(src)
	if err != nil {
		return err
	}
	dstOpts, err := parseObj(dst)
	if err != nil {
		return err
	}
	return g.copyToBucket(srcOpts.Bucket, srcOpts.Key, dstOpts.Bucket, dstOpts.Key)
}

// MoveObject rename this object
func (g *GCSStorage) MoveObject(src, dst string) error {
	srcOpts, err := parseObj(src)
	if err != nil {
		return err
	}
	dstOpts, err := parseObj(dst)
	if err != nil {
		return err
	}
	return g.move(srcOpts.Bucket, srcOpts.Key, dstOpts.Bucket, dstOpts.Key)
}

// IsExist return false if node doesn't exist
func (g *GCSStorage) IsExist(node string) bool {
	opts, err := parseObj(node)
	if err != nil {
		return false
	}
	_, err = g.attrs(opts.Bucket, opts.Key)
	if err != nil {
		return false
	}
	return true
}

// ListObjects return all files via prefix dir
func (g *GCSStorage) ListObjects(dir string) ([]*Object, int64, error) {
	opts, err := parseObj(dir)
	if err != nil {
		return nil, 0, err
	}
	list, size, err := g.listByPrefix(opts.Bucket, appendPathSuffix(opts.Key), "", ObjectTypeIsObject)
	return list, size, err
}

func (g *GCSStorage) ListChildObjects(dir string) ([]*Object, int64, error) {
	opts, err := parseObj(dir)
	if err != nil {
		return nil, 0, err
	}
	list, size, err := g.listByPrefix(opts.Bucket, appendPathSuffix(opts.Key), "/", ObjectTypeIsObject)
	return list, size, err
}

// ListDirs return all dirs via prefix dir
func (g *GCSStorage) ListDirs(dir string) ([]string, error) {
	opts, err := parseObj(dir)
	if err != nil {
		return nil, err
	}
	objs, _, err := g.listByPrefix(opts.Bucket, appendPathSuffix(opts.Key), "/", ObjectTypeIsDir)
	strSlice := make([]string, len(objs))
	for k, v := range objs {
		strSlice[k] = v.FileName
	}
	return strSlice, err
}

func (g *GCSStorage) conn() (*gs.Client, error) {
	ctx := context.Background()
	var client *gs.Client
	var err error
	if g.Token != "" {
		cred := []byte(g.Token)
		opts := make([]option.ClientOption, 0)
		if cred != nil && len(cred) > 0 {
			opts = append(opts, option.WithCredentialsJSON(cred))
		}
		client, err = gs.NewClient(ctx, opts...)
	} else {
		client, err = gs.NewClient(ctx)
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (g *GCSStorage) write(data []byte, bucket, object string) error {
	client, err := g.conn()
	if err != nil {
		return err
	}
	ctx := context.Background()
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewReader(data)); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func (g *GCSStorage) read(bucket, object string) ([]byte, error) {
	client, err := g.conn()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (g *GCSStorage) move(srcBucket, srcObject, dstBucket, dstObject string) error {
	client, err := g.conn()
	if err != nil {
		return err
	}
	ctx := context.Background()
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	if err := src.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func (g *GCSStorage) copyToBucket(srcBucket, srcObject, dstBucket, dstObject string) error {
	client, err := g.conn()
	if err != nil {
		return err
	}
	ctx := context.Background()
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	return nil
}

func (g *GCSStorage) delete(bucket, object string) error {
	client, err := g.conn()
	if err != nil {
		return err
	}
	ctx := context.Background()
	o := client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func (g *GCSStorage) attrs(bucket, object string) (map[string]string, error) {
	client, err := g.conn()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for key, value := range attrs.Metadata {
		m[key] = value
	}
	return m, nil
}

func (g *GCSStorage) listByPrefix(bucket, prefix, delim string, types ...ObjectType) ([]*Object, int64, error) {
	var dirType, objectType bool
	for _, v := range types {
		if v == ObjectTypeIsDir {
			dirType = true
		}
		if v == ObjectTypeIsObject {
			objectType = true
		}
	}

	client, err := g.conn()
	if err != nil {
		return nil, 0, err
	}
	ctx := context.Background()
	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimiter, the entire  tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="a/", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="a/" and delim="/", you'll get back:
	//   /a/1.txt
	it := client.Bucket(bucket).Objects(ctx, &gs.Query{
		Prefix:    prefix,
		Delimiter: delim,
		Versions:  false,
	})
	m := make([]*Object, 0, 10000)
	var size int64
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		if attrs.Prefix == "" && attrs.Name == "" {
			continue
		}
		size += attrs.Size
		if attrs.Prefix != "" && dirType {
			m = append(m, &Object{
				FileName: fmt.Sprintf("gs://%s/%s", bucket, attrs.Prefix),
				Size:     attrs.Size,
				Sum:      fmt.Sprintf("%x", attrs.MD5),
				Created:  attrs.Created,
				Updated:  attrs.Updated,
			})
		}
		if attrs.Name != "" && objectType {
			m = append(m, &Object{
				FileName: fmt.Sprintf("gs://%s/%s", bucket, attrs.Name),
				Size:     attrs.Size,
				Sum:      fmt.Sprintf("%x", attrs.MD5),
				Created:  attrs.Created,
				Updated:  attrs.Updated,
			})
		}
	}
	return m, size, nil
}

// Download download file to local
func (g *GCSStorage) Download(from, to string) error {
	opts, err := parseObj(from)
	if err != nil {
		return err
	}
	file, err := os.Create(to)
	if err != nil {
		return err
	}
	defer file.Close()

	client, err := g.conn()
	if err != nil {
		return err
	}
	r, err := client.Bucket(opts.Bucket).Object(opts.Key).NewReader(context.Background())
	if err != nil {
		if err == gs.ErrObjectNotExist {
			return ErrCodeNoSuchKey
		}
		return err
	}
	defer r.Close()
	buf := make([]byte, 5*1024*1024) //5MB
	_, err = io.CopyBuffer(file, r, buf)
	if err != nil {
		return err
	}
	return nil
}

// Upload put file to remote
func (g *GCSStorage) Upload(from, to string) error {
	file, err := os.Open(from)
	if err != nil {
		return err
	}
	defer file.Close()

	opts, err := parseObj(to)
	if err != nil {
		return err
	}
	client, err := g.conn()
	if err != nil {
		return err
	}
	w := client.Bucket(opts.Bucket).Object(opts.Key).NewWriter(context.Background())
	defer w.Close()
	buf := make([]byte, 5*1024*1024) //5MB
	_, err = io.CopyBuffer(w, file, buf)
	if err != nil {
		return err
	}
	return nil
}
