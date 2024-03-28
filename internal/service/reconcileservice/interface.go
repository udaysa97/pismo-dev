package reconcileservice

import (
	"context"
	"pismo-dev/internal/types"
)

type ReconcileServiceInterface interface {
	ProcessJobStatus(ctx context.Context, jobId string) bool
	GetJobStatusFromFM(ctx context.Context, jobId string) (types.FMGetStatusResult, error) //TODO: might remove it later
	SetRequiredRepos(repos RequiredRepos)
	SetRequiredServices(services RequiredServices)
	InitiateSqsConsumer()
}
