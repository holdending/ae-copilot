package storage

import "testing"

func TestGCSStorage_ListChildObjects(t *testing.T) {
	dir := "gs://lranalytics-au-endpoint-select-vm/721211/REJECT/"
	client := NewGCSStorage(nil)
	ls, size, err := client.ListChildObjects(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(size)
	for k, v := range ls {
		t.Log(k, v)
	}
}

func TestGCSStorage_ListDirs(t *testing.T) {
	dir := "gs://lr-select-vm-us-qa-etl/data/ccpa/output/"
	client := NewGCSStorage(nil)
	ls, err := client.ListDirs(dir)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range ls {
		t.Log(k, v)
	}
}

func TestGCSStorage_ListObjects(t *testing.T) {
	dir := "gs://lr-select-vm-us-qa-etl/data/ccpa/output/"
	client := NewGCSStorage(nil)
	ls, size, err := client.ListObjects(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(size)
	for k, v := range ls {
		t.Log(k, v)
	}
}

func TestGCSStorage_GetObject(t *testing.T) {
	dir := "gs://lr-select-vm-us-qa-etl/data/ccpa/output//AUDIENCE_1577347828545319742_pel///metadata.json"
	client := NewGCSStorage(nil)
	data, err := client.GetObject(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

	if _, err := client.GetObject(client.PathJoin("gs://lr-select-vm-us-qa-etl/data/dev-test/ingestion/output/AUDIENCE_1603419314179682000_clusterdemo", "_SUCCESS")); err != nil {
		t.Fatal(err)
	}

	client.MoveObject("gs://lr-select-vm-us-qa-etl/data/dev-test/ingestion/output/AUDIENCE_1603419314179682000_clusterdemo/_SUCCESS", "gs://lr-select-vm-us-qa-etl/data/dev-test/ingestion/output/AUDIENCE_1603419314179682000_clusterdemo/_ETLSUCCESS")
}

func TestGCSStorage_RemoveAllObject(t *testing.T) {
	dir := "gs://lr-select-tv-us-qa-etl/data/smartaudience/outbound/tv/AUDIENCE_1639445694069150203_autovivian/"
	client := NewGCSStorage(nil)
	if err := client.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
}

func Test_listByPrefix(t *testing.T) {
	client := NewGCSStorage(nil)
	objs, _, err := client.listByPrefix("lr-select-vm-us-qa-etl", "", "/", ObjectTypeIsDir)
	if err != nil {
		t.Fatal(err)
	}
	strSlice := make([]string, len(objs))
	for k, v := range objs {
		strSlice[k] = v.FileName
	}
	t.Log(strSlice)
}
