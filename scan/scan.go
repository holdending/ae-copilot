package scan

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/LiveRamp/ae-copilot/config"
	"github.com/LiveRamp/ae-copilot/models"
	"github.com/LiveRamp/ae-copilot/pkg/libs/storage"
	"github.com/LiveRamp/ae-copilot/services"
	constant "github.com/LiveRamp/ae-copilot/utils"
	"github.com/astaxie/beego/logs"
)

var rejectedFileScannerClient *rejectedFileScanner

type Scanner interface {
	AsyncRunning(ctx context.Context)
	Close()
}

func AsyncRunning() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.TODO())
	ingestionJob := NewRejectedFileScanner(time.Second * time.Duration(config.Agent.ScanIntervalTime))
	ingestionJob.AsyncRunning(ctx)
	return cancel
}

func NewRejectedFileScanner(duration time.Duration) *rejectedFileScanner {
	if rejectedFileScannerClient == nil {
		rejectedFileScannerClient = &rejectedFileScanner{
			duration: duration,
			skip:     map[string]bool{},
		}
	}
	return rejectedFileScannerClient
}

type rejectedFileScanner struct {
	duration time.Duration
	stopped  bool
	skip     map[string]bool
}

func (s *rejectedFileScanner) AsyncRunning(ctx context.Context) {
	logs.Info("rejected file scanner starts running.")
	logs.Info(s.duration)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Info("rejectedFileScanner.AsyncRunning", r)
				go s.AsyncRunning(ctx)
			}
		}()
		t1 := time.NewTimer(s.duration)
		for {
			if s.stopped {
				return
			}
			select {
			case <-t1.C:
				s.scanning()
				t1.Reset(s.duration)
			case <-ctx.Done():
				logs.Info("rejected file scanner stopped.")
				return
			}
		}
	}()
}

func (s *rejectedFileScanner) scanning() {
	for _, tenant := range config.Agent.Tenants {
		files := s.listObjects(fmt.Sprintf(config.Agent.RejectPath, tenant, constant.REJECT_PATH_PREFIX))
		for _, file := range files {
			if strings.HasSuffix(file, constant.CSV_SUFFIX) && !s.isExist(file+constant.SCANED_SUFFIX, files) {
				task := &models.RejectedFileRemediationTask{
					TaskName:       file,
					RejectedPrefix: file,
					InPrefix:       strings.Replace(file, constant.REJECT_PATH_PREFIX, constant.IN_PATH_PREFIX, 1),
				}
				s.tryToDoTheTask(task)
			}
		}
	}
}

func (s *rejectedFileScanner) tryToDoTheTask(task *models.RejectedFileRemediationTask) {
	logs.Info("try to do the task %s.", task.TaskName)
	s.skip[task.TaskName] = true
	if err := s.putObject(task.RejectedPrefix+constant.SCANED_SUFFIX, []byte{}); err != nil {
		logs.Error("put scanned file error:" + err.Error())
	}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	services.Processing(ctx, task)
}

func (s *rejectedFileScanner) listObjects(path string) []string {
	files := []string{}
	fs := storage.NewStorageClient(path, config.Agent.GCSCredentials)
	if folders, err := fs.ListDirs(path); err == nil {
		for _, folder := range folders {
			if ls, _, err := fs.ListChildObjects(folder); err == nil {
				for _, v := range ls {
					files = append(files, v.FileName)
				}
			}
		}
	}
	return files
}

func (s *rejectedFileScanner) isExist(file string, files []string) bool {
	exist := false
	for _, f := range files {
		if f == file {
			exist = true
			break
		}
	}
	return exist
}

func (s *rejectedFileScanner) putObject(path string, data []byte) error {
	fs := storage.NewStorageClient(path, config.Agent.GCSCredentials)
	return fs.PutObject(path, data)
}
