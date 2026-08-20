package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	domainuns "backend/internal/domain/uns"
	"backend/internal/infra/outbox"
	"backend/internal/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type unsImportJobRequest struct {
	DryRun       bool  `json:"dryRun"`
	SourceFileID int64 `json:"sourceFileId"`
}

type unsExportJobRequest struct {
	RootNodeID int64 `json:"rootNodeId"`
}

type unsPhysicalEnsureRequest struct {
	Nodes []repo.UnsNode `json:"nodes"`
}

const unsPhysicalEnsureEventType = "uns.physical.ensure"

func newOutboxWorker(app *App) *outbox.Worker {
	worker := outbox.NewWorker(context.Background())
	if app == nil {
		return worker
	}
	worker.Register("uns.import.job", app.handleUnsImportJob)
	worker.Register("uns.export.job", app.handleUnsExportJob)
	worker.Register(unsPhysicalEnsureEventType, app.handleUnsPhysicalEnsure)
	return worker
}

func (a *App) handleUnsPhysicalEnsure(ctx context.Context, event repo.OutboxEvent) error {
	if a == nil || a.DataIngest == nil {
		return errors.New("physical ensure dependencies are not ready")
	}
	var req unsPhysicalEnsureRequest
	if err := json.Unmarshal(event.Payload, &req); err != nil {
		return fmt.Errorf("invalid physical ensure payload: %w", err)
	}
	if len(req.Nodes) == 0 {
		return nil
	}
	started := time.Now()
	metricCount := 0
	aliases := map[string]struct{}{}
	for _, node := range req.Nodes {
		if node.TopicType == 3 {
			metricCount++
		}
		alias := strings.Trim(strings.TrimSpace(node.Alias), "/")
		if alias == "" {
			alias = strings.Trim(strings.TrimSpace(node.Namespace), "/")
		}
		if alias != "" {
			aliases[alias] = struct{}{}
		}
	}
	if err := a.DataIngest.EnsurePhysicalForNodes(ctx, req.Nodes); err != nil {
		return fmt.Errorf("uns physical ensure batch failed nodes=%d metrics=%d aliases=%d duration=%s: %w", len(req.Nodes), metricCount, len(aliases), time.Since(started), err)
	}
	logx.WithContext(ctx).Infof("uns physical ensure batch completed nodes=%d metrics=%d aliases=%d duration=%s", len(req.Nodes), metricCount, len(aliases), time.Since(started))
	return nil
}

func (a *App) handleUnsImportJob(ctx context.Context, event repo.OutboxEvent) error {
	jobs := repo.NewAsyncJobRepo(ctx)
	job, ok, err := loadAsyncJob(ctx, jobs, event, "uns.import")
	if err != nil || !ok {
		return err
	}
	var req unsImportJobRequest
	if err := json.Unmarshal(job.RequestJSON, &req); err != nil {
		return markAsyncJobFailed(ctx, jobs, job, "invalid import job request: "+err.Error(), nil)
	}
	if req.SourceFileID <= 0 {
		return markAsyncJobFailed(ctx, jobs, job, "sourceFileId is required", nil)
	}
	if a == nil || a.Asset == nil || a.UNS == nil {
		return errors.New("async import dependencies are not ready")
	}
	if err := jobs.MarkAsyncJobRunning(ctx, job.JobKey); err != nil {
		return err
	}
	file, reader, _, err := a.Asset.Open(ctx, req.SourceFileID)
	if err != nil {
		return markAsyncJobFailed(ctx, jobs, job, err.Error(), map[string]any{"sourceFileId": req.SourceFileID})
	}
	defer reader.Close()
	raw, err := io.ReadAll(reader)
	if err != nil {
		return markAsyncJobFailed(ctx, jobs, job, err.Error(), map[string]any{"sourceFileId": req.SourceFileID})
	}
	progress := func(status domainuns.ImportStatus) {
		if status.Progress >= 0 {
			if err := jobs.UpdateAsyncJobProgress(ctx, job.JobKey, int(status.Progress)); err != nil {
				logx.WithContext(ctx).Errorf("async import progress update failed jobKey=%s err=%v", job.JobKey, err)
			}
		}
	}
	var result map[string]any
	if req.DryRun {
		result, err = a.UNS.ValidateImportData(ctx, file.OriginalName, raw, progress)
	} else {
		result, err = a.UNS.ImportDataStream(ctx, file.OriginalName, raw, job.CreatedBy, progress)
	}
	if err != nil {
		return markAsyncJobFailed(ctx, jobs, job, err.Error(), result)
	}
	return jobs.MarkAsyncJobSucceeded(ctx, job.JobKey, result)
}

func (a *App) handleUnsExportJob(ctx context.Context, event repo.OutboxEvent) error {
	jobs := repo.NewAsyncJobRepo(ctx)
	job, ok, err := loadAsyncJob(ctx, jobs, event, "uns.export")
	if err != nil || !ok {
		return err
	}
	var req unsExportJobRequest
	if err := json.Unmarshal(job.RequestJSON, &req); err != nil {
		return markAsyncJobFailed(ctx, jobs, job, "invalid export job request: "+err.Error(), nil)
	}
	if a == nil || a.UNS == nil {
		return errors.New("async export dependencies are not ready")
	}
	if err := jobs.MarkAsyncJobRunning(ctx, job.JobKey); err != nil {
		return err
	}
	cmd := domainuns.ExportCommand{ExportType: "ALL"}
	if req.RootNodeID > 0 {
		cmd.ExportType = "PART"
		cmd.Folders = []int64{req.RootNodeID}
	}
	result, err := a.UNS.ExportData(ctx, cmd)
	if err != nil {
		return markAsyncJobFailed(ctx, jobs, job, err.Error(), nil)
	}
	return jobs.MarkAsyncJobSucceeded(ctx, job.JobKey, result)
}

func loadAsyncJob(ctx context.Context, jobs *repo.AsyncJobRepo, event repo.OutboxEvent, expectedType string) (repo.AsyncJob, bool, error) {
	if strings.TrimSpace(event.AggregateType) != "asyncJob" {
		return repo.AsyncJob{}, false, fmt.Errorf("invalid aggregate type %q for async job event", event.AggregateType)
	}
	jobKey := strings.TrimSpace(event.AggregateID)
	if jobKey == "" && len(event.Payload) > 0 {
		var payload repo.AsyncJob
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			jobKey = strings.TrimSpace(payload.JobKey)
		}
	}
	if jobKey == "" {
		return repo.AsyncJob{}, false, errors.New("async job event missing aggregate id")
	}
	job, err := jobs.GetAsyncJob(ctx, jobKey)
	if err != nil {
		return repo.AsyncJob{}, false, err
	}
	if expectedType != "" && job.JobType != expectedType {
		return repo.AsyncJob{}, false, fmt.Errorf("unexpected job type %q, want %q", job.JobType, expectedType)
	}
	if isTerminalAsyncJobStatus(job.Status) {
		return job, false, nil
	}
	return job, true, nil
}

func isTerminalAsyncJobStatus(status string) bool {
	switch status {
	case repo.AsyncJobStatusSucceeded, repo.AsyncJobStatusFailed, repo.AsyncJobStatusCanceled:
		return true
	default:
		return false
	}
}

func markAsyncJobFailed(ctx context.Context, jobs *repo.AsyncJobRepo, job repo.AsyncJob, message string, result map[string]any) error {
	if result == nil {
		result = map[string]any{}
	}
	if strings.TrimSpace(message) != "" {
		result["error"] = message
	}
	return jobs.MarkAsyncJobFailed(ctx, job.JobKey, message, result)
}
