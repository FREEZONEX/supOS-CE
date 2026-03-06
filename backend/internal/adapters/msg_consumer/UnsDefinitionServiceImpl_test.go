package msg_consumer

import (
	"backend/internal/repo/relationDB"
	"backend/internal/types"
	"context"
	"testing"

	"github.com/karlseguin/ccache/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestGetByAlis(t *testing.T) {
	def := &UnsDefinitionService{
		log:   logx.WithContext(t.Context()),
		cache: ccache.New(ccache.Configure().MaxSize(1100000).Buckets(64).GetsPerPromote(3)),
	}
	logUns := func(rs *types.UnsDefinition) {
		if rs == nil {
			t.Log(rs)
		} else {
			t.Log(rs.Alias)
		}
	}
	r1 := def.getByAliasOrPath(keyAliasPrev, "alias", func(ctx context.Context, arg string) (*relationDB.UnsNamespace, error) {
		return &relationDB.UnsNamespace{Alias: arg}, nil
	})
	logUns(r1)
	r2 := def.getByAliasOrPath(keyAliasPrev, "alias", nil)
	logUns(r2)
}
