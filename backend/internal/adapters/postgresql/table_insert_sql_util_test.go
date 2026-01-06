package postgresql

import (
	"backend/internal/types"
	"backend/share/base"
	"context"
	"testing"

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
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@100.100.100.20:31014/postgres")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	args := []any{20, 21}
	tag, err := pool.Exec(context.Background(), `insert into public._dingdanceshi_edd1ebaeb5da4155bdbf (_id,"timeStamp",wst)
        values(default,default,$1),(default,default,$2)`, args...)
	if err != nil {
		panic(err)
	}
	t.Log(tag.RowsAffected())
}
