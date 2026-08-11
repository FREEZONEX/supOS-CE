package dataingest

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"backend/internal/repo"
)

func TestStatePayloadPreservesCompleteJSONObject(t *testing.T) {
	def := &Definition{
		ID:             1,
		Namespace:      "State/action1",
		Alias:          "action1",
		TopicType:      1,
		Fields:         []Field{{Name: "aaa", Type: "integer"}, {Name: "timestamp", Type: "integer"}},
		TimestampField: "timestamp",
	}
	before := time.Now().UTC()
	records, err := NormalizePayload(context.Background(), def, []byte(`{"aaa":123222333452,"timestamp":1732223334512,"ccc":222}`))
	after := time.Now().UTC()
	if err != nil {
		t.Fatalf("NormalizePayload failed: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("records = %d, want 1", len(records))
	}
	if records[0].Time.Before(before) || records[0].Time.After(after) {
		t.Fatalf("record time = %s, want current processing time in [%s, %s]", records[0].Time, before, after)
	}

	stored, err := eventRecordJSON(records[0])
	if err != nil {
		t.Fatalf("eventRecordJSON failed: %v", err)
	}
	want := `{"aaa":123222333452,"ccc":222,"timestamp":1732223334512}`
	if string(stored) != want {
		t.Fatalf("stored payload = %s, want %s", stored, want)
	}

	live := livePayloadFromRecord(records[0])
	assertCompleteEventPayload(t, live["payload"], "live")
	assertCompleteEventPayload(t, eventPayloadFromSinkJSON(stored), "sink")
}

func TestMetricPayloadRemainsSchemaDriven(t *testing.T) {
	def := &Definition{
		ID:             2,
		Namespace:      "Metric/action1",
		Alias:          "metric-action1",
		TopicType:      metricTopicType,
		Fields:         []Field{{Name: "aaa", Type: "long"}, {Name: metricTimeColumn, Type: "datetime"}},
		TimestampField: metricTimeColumn,
	}
	before := time.Now().UTC()
	records, err := NormalizePayload(context.Background(), def, []byte(`{"aaa":123222333452,"_timestamp":1732223334512,"ccc":222}`))
	after := time.Now().UTC()
	if err != nil {
		t.Fatalf("NormalizePayload failed: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("records = %d, want 1", len(records))
	}
	if records[0].Time.Before(before) || records[0].Time.After(after) {
		t.Fatalf("metric record time = %s, want current processing time in [%s, %s]", records[0].Time, before, after)
	}
	payload, ok := livePayloadFromRecord(records[0])["payload"].(map[string]any)
	if !ok {
		t.Fatalf("metric payload type = %T", livePayloadFromRecord(records[0])["payload"])
	}
	if len(payload) != 1 || payload["aaa"] != int64(123222333452) {
		t.Fatalf("metric payload = %#v, want only schema field aaa", payload)
	}
}

func assertCompleteEventPayload(t *testing.T, value any, source string) {
	t.Helper()
	payload, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s payload type = %T", source, value)
	}
	if len(payload) != 3 {
		t.Fatalf("%s payload = %#v, want 3 original fields", source, payload)
	}
	for key, want := range map[string]string{
		"aaa":       "123222333452",
		"timestamp": "1732223334512",
		"ccc":       "222",
	} {
		got, ok := payload[key].(json.Number)
		if !ok || got.String() != want {
			t.Fatalf("%s payload[%s] = %#v, want json number %s", source, key, payload[key], want)
		}
	}
}

func TestWriteOpenAPIValuePublishesMQTT(t *testing.T) {
	s := &Service{}

	var gotTopic string
	var gotPayload []byte
	var gotQoS byte
	var gotRetain bool
	s.publishOverride = func(topic string, qos byte, retain bool, payload []byte) error {
		gotTopic = topic
		gotQoS = qos
		gotRetain = retain
		gotPayload = payload
		return nil
	}

	node := repo.UnsNode{
		ID:        1,
		Namespace: "Factory/line1/Metric/temperature",
		Type:      2,
		TopicType: 3,
	}

	if err := s.WriteOpenAPIValue(context.Background(), node, "hello", 0, 2, true); err != nil {
		t.Fatalf("WriteOpenAPIValue failed: %v", err)
	}

	if gotTopic != "Factory/line1/Metric/temperature" {
		t.Errorf("topic = %q, want %q", gotTopic, "Factory/line1/Metric/temperature")
	}
	if gotQoS != 2 {
		t.Errorf("qos = %d, want 2", gotQoS)
	}
	if !gotRetain {
		t.Error("retain = false, want true")
	}
	if string(gotPayload) != "hello" {
		t.Errorf("payload = %q, want raw string %q", string(gotPayload), "hello")
	}
}

func TestWriteOpenAPIValuePublishesJSONObject(t *testing.T) {
	s := &Service{}

	var gotPayload []byte
	s.publishOverride = func(topic string, qos byte, retain bool, payload []byte) error {
		gotPayload = payload
		return nil
	}

	node := repo.UnsNode{
		ID:        2,
		Namespace: "Factory/line1/Metric/pressure",
		Type:      2,
		TopicType: 3,
	}

	value := map[string]any{"pressure": 101.3}
	if err := s.WriteOpenAPIValue(context.Background(), node, value, 0, 0, false); err != nil {
		t.Fatalf("WriteOpenAPIValue failed: %v", err)
	}

	want := `{"pressure":101.3}`
	if string(gotPayload) != want {
		t.Errorf("payload = %q, want %q", string(gotPayload), want)
	}
}

func TestWriteOpenAPIValueNoMQTTClientFails(t *testing.T) {
	s := &Service{}
	node := repo.UnsNode{
		ID:        3,
		Namespace: "Factory/line1/Metric/flow",
		Type:      2,
		TopicType: 3,
	}

	err := s.WriteOpenAPIValue(context.Background(), node, "x", 0, 1, false)
	if err == nil || err.Error() != "mqtt client is not initialized" {
		t.Fatalf("WriteOpenAPIValue error = %v, want mqtt client is not initialized", err)
	}
}

func TestWriteOpenAPIValuePublishErrorReturned(t *testing.T) {
	s := &Service{}
	s.publishOverride = func(topic string, qos byte, retain bool, payload []byte) error {
		return errors.New("broker down")
	}

	node := repo.UnsNode{
		ID:        4,
		Namespace: "Factory/line1/Metric/level",
		Type:      2,
		TopicType: 3,
	}

	err := s.WriteOpenAPIValue(context.Background(), node, "y", 0, 1, false)
	if err == nil || err.Error() != "broker down" {
		t.Fatalf("WriteOpenAPIValue error = %v, want broker down", err)
	}
}
