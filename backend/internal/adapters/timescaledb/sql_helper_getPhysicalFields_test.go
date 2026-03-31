package timescaledb

import (
	"encoding/json"
	"testing"
)

func TestGetTableInfo(t *testing.T) {
	db := newTsdbPersistentService("postgres://postgres:postgres@100.100.100.20:31014/postgres")
	if db == nil {
		t.Fatal("NoDB")
	}
	fs, ct, qos, tag, err := getPhysicalTableFields(t.Context(), db.dbPool)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ct=%s, qos=%s, tag=%s", ct, qos, tag)
	jsobs, _ := json.MarshalIndent(fs, "", " ")
	t.Logf("fs: %+v\n", string(jsobs))
}
