package timescaledb

import (
	"backend/internal/common/utils/dbpool"
	"context"
	"encoding/json"
	"testing"
)

func TestPgQueryParserView(t *testing.T) {
	pool, err := dbpool.NewPool(context.Background(), "postgres://postgres:postgres@192.168.235.152:2345/postgres", "test")
	if err != nil {
		t.Fatal(err)
	}
	rs, err := parseViews(pool, context.Background(), "public", "table3")
	if err != nil {
		t.Fatal(err)
	}
	bs, _ := json.Marshal(rs)
	t.Log(string(bs))
}
