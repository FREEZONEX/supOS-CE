package msg_consumer

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/karlseguin/ccache/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type UnsDefinitionService struct {
	log                  logx.Logger
	unsMapper            dao.UnsNamespaceRepo
	cache                *ccache.Cache
	locks                [1000]sync.RWMutex
	persistentServiceMap map[types.SrcJdbcType]serviceApi.IPersistentService
}

func init() {
	spring.RegisterBean[*UnsDefinitionService](&UnsDefinitionService{
		log:   logx.WithContext(context.Background()),
		cache: ccache.New(ccache.Configure().MaxSize(1100000).Buckets(64).GetsPerPromote(3)),
	})
}

const keyAliasPrev = "a:"
const keyPathPrev = "p:"

func (u *UnsDefinitionService) GetDefinitionByAlias(alias string) *types.UnsDefinition {
	return u.getByAliasOrPath(keyAliasPrev, alias, u.unsMapper.SelectByAlias)
}

func (u *UnsDefinitionService) GetDefinitionByPath(path string) *types.UnsDefinition {
	return u.getByAliasOrPath(keyPathPrev, path, u.unsMapper.SelectByPath)
}

func (u *UnsDefinitionService) GetDefinitionById(id int64) (rs *types.UnsDefinition) {
	c := u.cache
	idStr := id2key(id)
	item := c.Get(idStr)
	if item == nil || item.Expired() {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
		defer cancel()
		db := dao.GetDb(ctx)
		po, _ := u.unsMapper.SelectById(db, id)
		if po != nil {
			rs = po2dto(po)
			c.Set(idStr, rs, 10*time.Minute)
			key := ""
			if constants.UseAliasAsTopic {
				key = alias2key(po.Alias)
			} else {
				key = path2key(po.Path)
			}
			c.Set(key, rs.Id, 10*time.Minute)
		} else {
			c.Set(idStr, nil, 1*time.Minute) //占位
		}
	} else {
		if item.TTL() < 10*time.Second { //过期时间不到10秒，续期5分钟,避免对象使用中突然被释放遇到NPE
			item.Extend(5 * time.Minute)
		}
		if df, ok := item.Value().(*types.UnsDefinition); ok {
			rs = df
		}
	}
	return rs
}
func (u *UnsDefinitionService) DeleteByIds(ids []int64) error {
	for _, id := range ids {
		u.cache.Delete(id2key(id))
	}
	return nil
}

func (u *UnsDefinitionService) SaveBatch(list []*types.UnsDefinition) error {
	for _, v := range list {
		u.invalidCache(v.Id, v.Alias, v.Path)
	}
	return nil
}

func (u *UnsDefinitionService) DeleteBatch(list []*types.UnsDefinition) error {
	for _, v := range list {
		u.invalidCache(v.Id, v.Alias, v.Path)
	}
	return nil
}

func (u *UnsDefinitionService) OnEventBatchCreateTableEvent0(ev *event.BatchCreateTableEvent) {
	if list := ev.Creates; len(list) > 0 {
		for _, vs := range list {
			for _, v := range vs {
				u.invalidCache(v.GetId(), v.GetAlias(), v.GetPath())
			}
		}
	}
	if list := ev.Updates; len(list) > 0 {
		for _, vs := range list {
			for _, v := range vs {
				u.invalidCache(v.GetId(), v.GetAlias(), v.GetPath())
			}
		}
	}
}
func (u *UnsDefinitionService) OnEventRemoveTopicsEvent0(ev *event.RemoveTopicsEvent) {
	if len(ev.Topics) >= 0 {
		for _, v := range ev.Topics {
			u.invalidCache(v.GetId(), v.GetAlias(), v.GetPath())
		}
	}
}
func (u *UnsDefinitionService) OnEventUpdateInstanceEvent0(ev *event.UpdateInstanceEvent) {
	if len(ev.Topics) >= 0 {
		for _, v := range ev.Topics {
			u.invalidCache(v.Id, v.Alias, v.Path)
		}
	}
}
func (u *UnsDefinitionService) invalidCache(id int64, alias, path string) {
	u.cache.Delete(id2key(id))
	u.cache.Delete(alias2key(alias))
	u.cache.Delete(path2key(path))
}

func id2key(id int64) string {
	return strconv.FormatInt(id, 36)
}
func alias2key(alias string) string {
	return keyAliasPrev + alias
}
func path2key(path string) string {
	return keyPathPrev + path
}

func (u *UnsDefinitionService) getByAliasOrPath(kPrev string, arg string, query func(ctx context.Context, arg string) (*dao.UnsNamespace, error)) (rs *types.UnsDefinition) {
	key := kPrev + arg
	c := u.cache

	idObj := c.Get(key)
	if idObj == nil || idObj.Expired() {
		index := base.Abs(base.HashCode(arg)) % len(u.locks)
		u.locks[index].Lock()

		idObj = c.Get(key)
		if idObj == nil || idObj.Expired() {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
			defer cancel()
			po, _ := query(ctx, arg)
			if po != nil {
				idKey := id2key(po.Id)
				c.Set(key, po.Id, 13*time.Minute)
				rs = po2dto(po)
				c.Set(idKey, rs, 10*time.Minute)
			} else {
				c.Set(key, int64(-1), 2*time.Minute) //占位
			}
			u.locks[index].Unlock()
			return
		} else {
			u.locks[index].Unlock()
			return u.GetDefinitionById(idObj.Value().(int64))
		}
	} else {
		return u.GetDefinitionById(idObj.Value().(int64))
	}
	return
}

func po2dto(po *dao.UnsNamespace) *types.UnsDefinition {
	rs := UnsConverter.Po2Dto(po)
	fields := rs.Fields
	def := &types.UnsDefinition{CreateTopicDto: *rs}
	if len(fields) > 0 {
		for _, field := range fields {
			field.Uns = def
		}
	}
	return def
}
func (u *UnsDefinitionService) OnEventContextRefreshedEvent1(_ *event.ContextRefreshedEvent) {
	u.persistentServiceMap = base.MapArrayToMap(spring.GetBeansOfType[serviceApi.IPersistentService](),
		func(e serviceApi.IPersistentService) (ok bool, k types.SrcJdbcType, v serviceApi.IPersistentService) {
			return true, e.GetDataSrcId(), e
		})
	types.UnsLastValueFill = u.fillUnsLastValue
}
func (u *UnsDefinitionService) fillUnsLastValue(uns *types.UnsDefinition) {
	// 查询数据库表最新的一条数据，填充字段的 lastValue
	psv, has := u.persistentServiceMap[types.SrcJdbcType(uns.DataSrcID)]
	if !has {
		return
	}
	psv.FillLastRecord(uns)
}
