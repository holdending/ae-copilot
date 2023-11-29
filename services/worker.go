package services

import (
	"context"

	"github.com/LiveRamp/ae-copilot/models"
	"github.com/LiveRamp/ae-copilot/services/job"
)

func Processing(ctx context.Context, task *models.RejectedFileRemediationTask) error {
	hygiene := job.NewHygiene()
	if err := hygiene.Running(task); err != nil {
		return err
	}
	return nil
}
