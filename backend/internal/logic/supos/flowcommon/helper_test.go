package flowcommon

import (
	"context"
	"testing"

	"backend/internal/common"
	"backend/internal/repo/relationDB"
)

type fakeFlowRepo struct {
	source   *relationDB.NoderedFlow
	inserted *relationDB.NoderedFlow
}

func (f *fakeFlowRepo) FindOne(_ context.Context, id int64) (*relationDB.NoderedFlow, error) {
	if f.source != nil && f.source.ID == id {
		return f.source, nil
	}
	return nil, nil
}

func (f *fakeFlowRepo) Insert(_ context.Context, data *relationDB.NoderedFlow) error {
	f.inserted = data
	return nil
}

func (f *fakeFlowRepo) Update(_ context.Context, _ *relationDB.NoderedFlow) error {
	return nil
}

func (f *fakeFlowRepo) ReplaceModels(_ context.Context, _ int64, _ []string) error {
	return nil
}

func TestCopyFlowSetsCreatorFromInput(t *testing.T) {
	common.InitSnowflake(1)
	repo := &fakeFlowRepo{
		source: &relationDB.NoderedFlow{
			ID:       1001,
			FlowData: `[{"id":"node-1","z":"tab-1","type":"debug"}]`,
		},
	}

	record, err := CopyFlow(context.Background(), nil, repo, 1001, FlowCopyInput{
		FlowName: "copied-flow",
		Template: "node-red",
		Creator:  "tier0",
	}, nil)
	if err != nil {
		t.Fatalf("CopyFlow returned error: %v", err)
	}
	if record == nil {
		t.Fatalf("CopyFlow returned nil record")
	}
	if record.Creator != "tier0" {
		t.Fatalf("expected creator to be tier0, got %q", record.Creator)
	}
	if repo.inserted == nil {
		t.Fatalf("expected inserted record to be captured")
	}
	if repo.inserted.Creator != "tier0" {
		t.Fatalf("expected inserted creator to be tier0, got %q", repo.inserted.Creator)
	}
}
