package models

import (
	"time"
)

type JobType int16

const (
	_ JobType = iota
	JobHygiene
)

func (j JobType) String() string {
	switch j {
	case JobHygiene:
		return "hygiene"
	default:
		return "UNKNOWN JOB TYPE"
	}
}

type TaskStatus int

const (
	Ready TaskStatus = iota
	Running
	Success
	Failed
)

func (s TaskStatus) String() string {
	switch s {
	case Ready:
		return "Ready"
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

type JobTask struct {
	ID        int64
	Name      string
	Type      string
	Status    string
	Heartbeat int64
	ErrorInfo string
	CreatedAt time.Time
	UpdatedAt time.Time
}
