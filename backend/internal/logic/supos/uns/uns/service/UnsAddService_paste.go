package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"fmt"
	"strconv"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

/**
 * 文件或者文件夹复制黏贴
 * @param req.sourceId 被复制的顶级文件（夹）ID
 * @param req.targetParentId 粘贴的目的地顶层文件夹ID
 * @param req.newFile 源顶层文件夹重命名
 */
func (u *UnsAddService) PasteFolderOrFile(ctx context.Context, req *types.PasteRequestVO) (resp *types.CreateUnsResp, err error) {
	db := dao.GetDb(ctx)
	resp = &types.CreateUnsResp{}
	resp.Code, resp.Msg = 200, "ok"
	if req.SourceId == 0 {
		resp.Code = 400
		resp.Msg = "NoneSrcId"
		return
	}
	var ids = make([]int64, 0, 2)
	ids = append(ids, req.SourceId)
	if targetId := req.TargetId; targetId > 0 && targetId != req.SourceId {
		ids = append(ids, targetId)
	}
	var idMap map[int64]*dao.UnsNamespace
	{
		list, _ := u.unsMapper.SelectByIds(db, ids)
		idMap = base.MapArrayToMap(list, func(e *dao.UnsNamespace) (ok bool, k int64, v *dao.UnsNamespace) {
			return true, e.Id, e
		})
	}
	if len(idMap) < len(ids) {
		resp.Code = 400
		if len(idMap) == 0 {
			resp.Msg = I18nUtils.GetMessage("uns.folder.or.file.not.found") + ":" + fmt.Sprintf("%+v", ids)
		} else {
			for _, id := range ids {
				if !base.MapContainsKey(idMap, id) {
					resp.Msg = I18nUtils.GetMessage("uns.folder.or.file.not.found") + ":" + strconv.FormatInt(id, 10)
					break
				}
			}
		}
		return
	}
	src, tar := idMap[req.SourceId], idMap[req.TargetId]
	srcUns, parentAliasMap := u.getSrcUns(src, req, tar)
	var positioningUns = srcUns //返回给前端的定位UNS
	countChildren := int64(0)
	if src.PathType == constants.PathTypeDir {
		var tipMap map[string]string
		var countList int
		countChildren, countList, tipMap = u.copyChildren(ctx, db, src, srcUns, &positioningUns, parentAliasMap, resp)
		if len(tipMap) > 0 {
			if len(tipMap) < countList {
				setPasteData(positioningUns, resp)
				for _, tip := range tipMap {
					resp.Code, resp.Msg = 206, tip
					break
				}
			} else {
				for _, tip := range tipMap {
					resp.Code, resp.Msg = 400, tip
					break
				}
			}
			return resp, nil
		}
	}
	if countChildren == 0 {
		rs := u.CreateModelInstance(ctx, srcUns)
		resp.Code, resp.Msg = rs.Code, rs.Msg
	}
	if resp.Code == 200 {
		setPasteData(positioningUns, resp)
	}
	return
}

func setPasteData(positioningUns *types.CreateTopicDto, resp *types.CreateUnsResp) {
	parentId := ""
	if pid := positioningUns.ParentId; pid != nil {
		parentId = strconv.FormatInt(*pid, 10)
	}
	resp.Data = types.CreateUnsResult{
		Id:       strconv.FormatInt(positioningUns.Id, 10),
		ParentId: parentId,
	}
}

func (u *UnsAddService) copyChildren(ctx context.Context,
	db *gorm.DB,
	src *dao.UnsNamespace,
	srcUns *types.CreateTopicDto,
	positioningUns **types.CreateTopicDto,
	parentAliasMap map[string]string,
	resp *types.CreateUnsResp) (int64, int, map[string]string) {

	countChildren := int64(0)
	page := &stores.PageInfo{Page: 1, Size: 1000, Orders: []stores.OrderBy{{Field: "id", Sort: stores.OrderAsc}}}
	MaxId := int64(0)
	queryMaxId := &MaxId
	for loop := true; loop; {
		poList, _ := u.unsMapper.ListByLayRec(db, src.LayRec+"/", page, queryMaxId)
		if countChildren == 0 {
			queryMaxId = nil
		}
		if len(poList) == 0 {
			break
		}
		page.Page++
		list := make([]*types.CreateTopicDto, 0, 1+len(poList))
		if countChildren == 0 {
			list = append(list, srcUns)
		}

		for _, po := range poList {
			if po.Id > MaxId {
				u.log.Infof("复制自己文件结束: id=%d > %d, len(list)=%d", po.Id, MaxId, len(list))
				loop = false
				break
			}
			newAlias := PathUtil.GenerateAliasWithRandom(po.Name)
			if po.PathType == constants.PathTypeDir {
				parentAliasMap[po.Alias] = newAlias
			}
			unsDto := UnsConverter.Po2Dto(po)
			list = append(list, unsDto)
			unsDto.Id = 0
			unsDto.ParentId = nil
			unsDto.Alias = newAlias
			if pAlias := unsDto.ParentAlias; pAlias != nil {
				if newA := parentAliasMap[*pAlias]; newA != "" {
					unsDto.ParentAlias = &newA
				}
			}
			flags := unsDto.WithFlags
			if flags != nil {
				var fl int32 = (*flags) & (^constants.UnsFlagWithFlow)
				unsDto.WithFlags = &fl
			}
		}
		if len(list) == 0 {
			loop = false
			break
		}
		if countChildren == 0 && u.sysConfig.EnableAutoCategorization && base.P2v(src.DataType) > 0 && len(list) > 1 {
			*positioningUns = list[1]
		}
		srcAlias := srcUns.Alias
		tipMap := u.CreateModelAndInstance(ctx, list, false)
		if len(tipMap) > 0 {
			resp.Code, resp.Msg = 400, fmt.Sprintf("%+v", tipMap)
			return countChildren, len(list), tipMap
		}
		if countChildren == 0 && srcAlias != srcUns.Alias {
			u.log.Infof("修正别名: %s -> %s", srcAlias, srcUns.Alias)
			parentAliasMap[src.Alias] = srcUns.Alias
		}
		listSize := int64(len(poList))
		countChildren += listSize
		if listSize < page.Size {
			break
		}
	}
	return countChildren, 0, nil
}

func (u *UnsAddService) getSrcUns(src *dao.UnsNamespace, req *types.PasteRequestVO, tar *dao.UnsNamespace) (*types.CreateTopicDto, map[string]string) {
	srcUns := UnsConverter.Po2Dto(src)
	parentAliasMap := make(map[string]string, 64)
	if nf := req.NewFile; nf != nil {
		nf.Id = 0
		nf.Alias = ""
		nf.ParentAlias = nil
		nf.ParentId = nil
		srcUns = nf
	}
	{
		srcUns.Id = 0
		srcUns.ParentId = nil
		srcUns.ParentAlias = nil
		srcUns.PathType = src.PathType

		if req.TargetId != req.SourceId || src.PathType == constants.PathTypeFile || base.P2v(src.DataType) == 0 || !u.sysConfig.EnableAutoCategorization {
			newAlias := PathUtil.GenerateAliasWithRandom(srcUns.Name)
			if src.PathType == constants.PathTypeDir {
				parentAliasMap[src.Alias] = newAlias
			}
			srcUns.Alias = newAlias
			if tar != nil {
				srcUns.ParentAlias = &tar.Alias
			} else {
				srcUns.ParentAlias = nil
			}
		} else {
			srcUns.Alias = src.Alias
			srcUns.ParentAlias = src.ParentAlias
			parentAliasMap[src.Alias] = src.Alias
		}
	}
	return srcUns, parentAliasMap
}
