package jsonstream

import (
	"context"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCsv2JsonStream(t *testing.T) {
	// 初始化数据库连接池
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@100.100.100.22:34099/postgres")
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Release()
	csvExporter := func(w io.Writer) error {
		query := `COPY (SELECT * FROM supos.uns_namespace WHERE path_type in(0,2) and status=1 order by lay_rec asc) TO STDOUT WITH CSV HEADER`
		_, err = conn.Conn().PgConn().CopyTo(context.Background(), w, query)
		return err
	}
	idIndex := -1
	parentIdIndex := -1
	nameIndex := 4
	csv2node := func(headers, values []string) *IdNode {
		if idIndex < 0 {
			for i, h := range headers {
				switch h {
				case "id":
					idIndex = i
				case "parent_id":
					parentIdIndex = i
				case "name":
					nameIndex = i
				}
			}
		}
		id, _ := strconv.ParseInt(values[idIndex], 10, 64)
		parentId, _ := strconv.ParseInt(values[parentIdIndex], 10, 64)
		return &IdNode{ID: id, ParentId: parentId, Name: values[nameIndex]}
	}

	jsonWriter := &strings.Builder{}
	countNodes, err := Csv2JsonStream(csvExporter, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, csv2node, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("countNodes:", countNodes)
	t.Log(jsonWriter.String())
}
