package logger

import (
	"os"
	"testing"

	"github.com/LiveRamp/ae-copilot/pkg/libs/storage"
)

func TestNewLogger(t *testing.T) {
	loger := NewLogger("file", `{"filename":"./logs.log","daily":true,"maxdays":30}`, 7)
	defer func() {
		os.Remove("./logs.log")
	}()
	loger.Debug("debug: %s", "123")
	localStorage := storage.NewStorage(storage.StorageInLocal, nil)
	if localStorage.IsExist("./logs.log") == false {
		t.Fatalf("Failed to create a new loger.")
	}
}
