package timescaledb

import (
	"backend/internal/common/utils/dbpool"
	"context"
	"encoding/json"
	"testing"
)

func Test_parseUnsViews(t *testing.T) {
	pool, err := dbpool.NewPool(context.Background(), "postgres://postgres:postgres@192.168.235.152:2345/postgres", "test")
	if err != nil {
		t.Fatal(err)
	}
	info, err := parseUnsViews(pool, context.Background(), "public", []string{"table2", "table3"})
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range info {
		bs, _ := json.Marshal(v)
		t.Log(k, " = ", string(bs))
	}
}
