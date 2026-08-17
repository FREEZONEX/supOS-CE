package uns

import (
	"context"
	"errors"
	"testing"

	"backend/internal/repo"
)

func newTestCreateTreeState(existing []repo.UnsNode) *createTreeState {
	state := &createTreeState{service: &Service{}, ctx: context.Background()}
	state.initialize(existing)
	return state
}

func folderInput(name string, children ...CreateTreeNode) CreateTreeNode {
	return CreateTreeNode{
		Command:  SaveCommand{Name: name, NodeType: "folder"},
		Children: children,
	}
}

func metricTopicInput(name string) CreateTreeNode {
	return CreateTreeNode{
		Command: SaveCommand{Name: name, NodeType: "file", TopicType: "metric"},
	}
}

func resultErrFor(t *testing.T, results []CreateTreeResult, name string) error {
	t.Helper()
	for _, result := range results {
		if result.Name == name {
			return result.Err
		}
	}
	t.Fatalf("no result row for node %q", name)
	return nil
}

func TestWalkCreatesFullTreeWhenNothingExists(t *testing.T) {
	state := newTestCreateTreeState(nil)

	state.walk(folderInput("Common",
		folderInput("TAG",
			folderInput("Metric",
				metricTopicInput("TAG_001"),
			),
		),
	), 0)

	if len(state.results) != 4 {
		t.Fatalf("results = %d, want 4", len(state.results))
	}
	for _, result := range state.results {
		if result.Err != nil {
			t.Fatalf("result for %q has unexpected error: %v", result.Name, result.Err)
		}
	}
	if len(state.created) != 4 {
		t.Fatalf("created = %d, want 4", len(state.created))
	}
	leaf := state.created[3]
	if leaf.Name != "TAG_001" || leaf.Namespace != "Common/TAG/Metric/TAG_001" {
		t.Fatalf("leaf = %+v, want TAG_001 under Common/TAG/Metric", leaf)
	}
	if leaf.ParentID != state.created[2].ID {
		t.Fatalf("leaf.ParentID = %d, want Metric folder ID %d", leaf.ParentID, state.created[2].ID)
	}
}

func TestWalkReusesExistingFoldersForNewLeaf(t *testing.T) {
	existing := []repo.UnsNode{
		{ID: 11, ParentID: 0, Name: "Common", Namespace: "Common", Type: 1},
		{ID: 12, ParentID: 11, Name: "TAG", Namespace: "Common/TAG", Type: 1},
		{ID: 13, ParentID: 12, Name: "Metric", Namespace: "Common/TAG/Metric", Type: 1},
	}
	state := newTestCreateTreeState(existing)

	state.walk(folderInput("Common",
		folderInput("TAG",
			folderInput("Metric",
				metricTopicInput("TAG_100"),
			),
		),
	), 0)

	for _, result := range state.results {
		if result.Err != nil {
			t.Fatalf("result for %q has unexpected error: %v", result.Name, result.Err)
		}
	}
	if len(state.created) != 1 {
		t.Fatalf("created = %d, want only the new leaf", len(state.created))
	}
	leaf := state.created[0]
	if leaf.Name != "TAG_100" {
		t.Fatalf("created node = %q, want TAG_100", leaf.Name)
	}
	if leaf.ParentID != 13 {
		t.Fatalf("leaf.ParentID = %d, want existing Metric folder ID 13", leaf.ParentID)
	}
	if leaf.Namespace != "Common/TAG/Metric/TAG_100" {
		t.Fatalf("leaf.Namespace = %q, want Common/TAG/Metric/TAG_100", leaf.Namespace)
	}
}

func TestWalkReusesExistingFoldersCaseInsensitively(t *testing.T) {
	existing := []repo.UnsNode{
		{ID: 11, ParentID: 0, Name: "Common", Namespace: "Common", Type: 1},
	}
	state := newTestCreateTreeState(existing)

	state.walk(folderInput("common",
		folderInput("Metric",
			metricTopicInput("TAG_101"),
		),
	), 0)

	for _, result := range state.results {
		if result.Err != nil {
			t.Fatalf("result for %q has unexpected error: %v", result.Name, result.Err)
		}
	}
	if len(state.created) != 2 {
		t.Fatalf("created = %d, want Metric folder + leaf", len(state.created))
	}
	if state.created[0].ParentID != 11 {
		t.Fatalf("Metric.ParentID = %d, want existing Common ID 11", state.created[0].ParentID)
	}
	if leaf := state.created[1]; leaf.Namespace != "Common/Metric/TAG_101" {
		t.Fatalf("leaf.Namespace = %q, want Common/Metric/TAG_101", leaf.Namespace)
	}
}

func TestWalkStillRejectsDuplicateTopicLeaf(t *testing.T) {
	existing := []repo.UnsNode{
		{ID: 11, ParentID: 0, Name: "Common", Namespace: "Common", Type: 1},
		{ID: 12, ParentID: 11, Name: "Metric", Namespace: "Common/Metric", Type: 1},
		{ID: 13, ParentID: 12, Name: "TAG_001", Namespace: "Common/Metric/TAG_001", Type: 2, TopicType: 3},
	}
	state := newTestCreateTreeState(existing)

	state.walk(folderInput("Common",
		folderInput("Metric",
			metricTopicInput("TAG_001"),
		),
	), 0)

	if err := resultErrFor(t, state.results, "TAG_001"); !errors.Is(err, repo.ErrDuplicate) {
		t.Fatalf("TAG_001 error = %v, want repo.ErrDuplicate", err)
	}
	if len(state.created) != 0 {
		t.Fatalf("created = %d, want 0", len(state.created))
	}
}

func TestWalkRejectsFolderCollidingWithExistingTopic(t *testing.T) {
	existing := []repo.UnsNode{
		{ID: 11, ParentID: 0, Name: "Common", Namespace: "Common", Type: 2, TopicType: 3},
	}
	state := newTestCreateTreeState(existing)

	state.walk(folderInput("Common",
		metricTopicInput("TAG_102"),
	), 0)

	if err := resultErrFor(t, state.results, "Common"); !errors.Is(err, repo.ErrDuplicate) {
		t.Fatalf("Common error = %v, want repo.ErrDuplicate", err)
	}
	if len(state.created) != 0 {
		t.Fatalf("created = %d, want 0", len(state.created))
	}
}

func TestValidateImportExistingSystemFolderAllowsZeroTopicType(t *testing.T) {
	tests := []struct {
		name      string
		topicType int16
	}{
		{name: "State", topicType: 1},
		{name: "Action", topicType: 2},
		{name: "Metric", topicType: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := repo.UnsNode{Name: tt.name, Type: 1, TopicType: 0}
			if err := validateImportExistingNode(node, "folder", true, tt.topicType); err != nil {
				t.Fatalf("validateImportExistingNode() error = %v, want nil", err)
			}
			topic, parentTopicType := importExistingNodeTopicContext(node, "", tt.name, true, tt.topicType)
			if topic != tt.name || parentTopicType != tt.topicType {
				t.Fatalf("importExistingNodeTopicContext() = (%q, %d), want (%q, %d)", topic, parentTopicType, tt.name, tt.topicType)
			}
			if nodeType := importNodeTypeNameForParent("", parentTopicType); nodeType != "file" {
				t.Fatalf("importNodeTypeNameForParent() = %q, want file for child of %s", nodeType, tt.name)
			}
		})
	}
}

func TestValidateImportExistingSystemFolderRejectsConflictingTopicType(t *testing.T) {
	node := repo.UnsNode{Name: "Metric", Type: 1, TopicType: 1}
	if err := validateImportExistingNode(node, "folder", true, 3); !errors.Is(err, ErrInvalid) {
		t.Fatalf("validateImportExistingNode() error = %v, want ErrInvalid", err)
	}
}
