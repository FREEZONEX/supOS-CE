package timescaledb

import (
	"backend/internal/common/utils/dbpool"
	"context"
	"encoding/json"
	"testing"
)

func TestPgQueryParserView(t *testing.T) {
	pool, err := dbpool.NewPool(context.Background(), "postgres://postgres:postgres@100.100.100.20:31014/postgres", "test")
	if err != nil {
		t.Fatal(err)
	}
	rs, err := parseViews(pool, context.Background(), "public", "opcua_demo_file_5483")
	if err != nil {
		t.Fatal(err)
	}
	bs, _ := json.Marshal(rs)
	t.Log(string(bs))
}
