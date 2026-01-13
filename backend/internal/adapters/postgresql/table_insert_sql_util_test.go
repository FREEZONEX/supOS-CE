package postgresql

import (
	"backend/internal/types"
	"backend/share/base"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestInsertSQL(t *testing.T) {
	def := &types.CreateTopicDto{Alias: "_dingdanceshi_edd1ebaeb5da4155bdbf", Fields: []*types.FieldDefine{
		{
			Name: "_id", Type: types.FieldTypeLong, Unique: base.OptionalTrue,
		},
		{
			Name: "timeStamp", Type: types.FieldTypeDatetime,
		},
		{
			Name: "wst", Type: types.FieldTypeLong,
		},
	}}
	data := []map[string]string{
		{"timeStamp": "1763707635481", "wst": "25"}, {"_id": "10", "wst": "27"}, {"wst": "28"},
	}
	sql, args := getInsertStatement(def, data)
	t.Log(sql, args)

	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@100.100.100.20:31014/postgres")
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	tag, err := pool.Exec(context.Background(), sql, args...)
	if err != nil {
		panic(err)
	}
	t.Log(tag.RowsAffected())
}
func TestInsertWithArgs(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@pgsql-ha.supos.app:50000/tier0")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	for i := 0; i < 3000; i++ {
		time.Sleep(time.Second)
		tag, err := pool.Exec(context.Background(), `INSERT INTO "supos_timeserial_long"("timeStamp", "tag", "value") 
        VALUES (NOW(), 1, (random() * 100)::int8)`)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(tag.RowsAffected())
		}
	}
}
