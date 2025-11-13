package msg_consumer

import (
	"backend/internal/common/event"
	"backend/internal/common/service"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UnsDefinitionService struct {
	log       logx.Logger
	unsMapper dao.UnsNamespaceRepo
	cache     *cache.Cache // map[int64]*types.CreateTopicDto
}

func init() {
	spring.RegisterBean[service.IUnsDefinitionService](&UnsDefinitionService{
		log:   logx.WithContext(context.Background()),
		cache: cache.New(10*time.Minute, 5*time.Minute),
	})
}

const keyAliasPrev = "a:"
const keyPathPrev = "p:"

func (u *UnsDefinitionService) GetDefinitionByAlias(alias string) *types.CreateTopicDto {
	return getByAliasOrPath(keyAliasPrev, u.cache, alias, u.unsMapper.GetByAlias)
}

func (u *UnsDefinitionService) GetDefinitionByPath(path string) *types.CreateTopicDto {
	return getByAliasOrPath(keyPathPrev, u.cache, path, u.unsMapper.GetByPath)
}

func (u *UnsDefinitionService) GetDefinitionById(id int64) (rs *types.CreateTopicDto) {
	c := u.cache
	idStr := strconv.FormatInt(id, 10)
	vu, exist := c.Get(idStr)
	if exist && vu != nil {
		rs = vu.(*types.CreateTopicDto)
	} else if !exist {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
		defer cancel()
		db := dao.GetDb(ctx)
		po, _ := u.unsMapper.SelectById(db, id)
		if po != nil {
			rs = UnsConverter.Po2Dto(po)
			c.Set(idStr, rs, 10*time.Minute)
			c.Set(keyAliasPrev+po.Alias, idStr, 10*time.Minute)
		} else {
			c.Set(idStr, nil, 1*time.Minute) //占位
		}
	}
	return
}

func (u *UnsDefinitionService) DeleteByIds(ids []int64) error {
	for _, id := range ids {
		u.cache.Delete(strconv.FormatInt(id, 10))
	}
	return nil
}

func (u *UnsDefinitionService) SaveBatch(list []*types.CreateTopicDto) error {
	for _, v := range list {
		u.cache.Delete(strconv.FormatInt(v.Id, 10))
		u.cache.Delete(v.Alias)
		u.cache.Delete(v.Path)
	}
	return nil
}

func (u *UnsDefinitionService) DeleteBatch(list []*types.CreateTopicDto) error {
	for _, v := range list {
		u.cache.Delete(strconv.FormatInt(v.Id, 10))
		u.cache.Delete(v.Alias)
		u.cache.Delete(v.Path)
	}
	return nil
}

func (u *UnsDefinitionService) OnBatchCreateTableEvent0(ev *event.BatchCreateTableEvent) {
	if list := ev.Creates; len(list) > 0 {
		for _, vs := range list {
			for _, v := range vs {
				u.cache.Delete(strconv.FormatInt(v.GetId(), 10))
				u.cache.Delete(v.GetAlias())
				u.cache.Delete(v.GetPath())
			}
		}
	}
	if list := ev.Updates; len(list) > 0 {
		for _, vs := range list {
			for _, v := range vs {
				u.cache.Delete(strconv.FormatInt(v.GetId(), 10))
				u.cache.Delete(v.GetAlias())
				u.cache.Delete(v.GetPath())
			}
		}
	}
}
func (u *UnsDefinitionService) OnRemoveTopicsEvent0(ev *event.RemoveTopicsEvent) {
	if len(ev.Topics) >= 0 {
		for _, v := range ev.Topics {
			u.cache.Delete(strconv.FormatInt(v.GetId(), 10))
			u.cache.Delete(v.GetAlias())
			u.cache.Delete(v.GetPath())
		}
	}
}
func (u *UnsDefinitionService) OnUpdateInstanceEvent0(ev *event.UpdateInstanceEvent) {
	if len(ev.Topics) >= 0 {
		for _, v := range ev.Topics {
			u.cache.Delete(strconv.FormatInt(v.Id, 10))
			u.cache.Delete(v.Alias)
			u.cache.Delete(v.Path)
		}
	}
}

func getByAliasOrPath(kPrev string, c *cache.Cache, arg string, query func(db *gorm.DB, arg string) (*dao.UnsNamespace, error)) (rs *types.CreateTopicDto) {
	key := kPrev + arg
	idObj, has := c.Get(key)
	if !has {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
		defer cancel()
		db := dao.GetDb(ctx)
		po, _ := query(db, arg)
		if po != nil {
			rs = UnsConverter.Po2Dto(po)
			idStr := strconv.FormatInt(po.Id, 10)
			c.Set(key, idStr, 10*time.Minute)
			c.Set(idStr, rs, 10*time.Minute)
		} else {
			c.Set(key, "-1", 2*time.Minute) //占位
		}
	} else {
		vu, exist := c.Get(idObj.(string))
		if exist && vu != nil {
			rs = vu.(*types.CreateTopicDto)
		}
	}
	return
}
