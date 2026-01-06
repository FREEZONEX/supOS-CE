package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"strings"
)

// buildCategoryFolderDto 构建分类文件夹DTO
func buildCategoryFolderDto(parentAlias *string, mountType *int16, mountSource *string, folderDataType *int16) *types.CreateTopicDto {
	dto := &types.CreateTopicDto{}
	fdt := enums.GetFolderDataType(base.P2v(folderDataType))
	if parentAlias != nil {
		// 使用父别名的情况
		dto.Alias = strings.ToLower(fdt.Name()) + "_" + *parentAlias
		dto.ParentAlias = parentAlias
		dto.MountSource = mountSource
		dto.MountType = mountType
	} else {
		// 没有父别名的情况
		dto.Alias = PathUtil.GenerateAliasWithRandom("_" + strings.ToLower(fdt.Name()) + "_")
		dto.ParentAlias = nil
	}

	// 设置其他属性
	dto.Id = common.NextId()
	//  edge 固定用分类文件夹的英文名 https://zentao.bluetron.cn/biz/bug-view-100705.html
	name := I18nUtils.GetMessageWithLang("en_US", fdt.String())
	dto.Name = name
	dto.DisplayName = &name
	dto.DataType = folderDataType
	dto.PathType = constants.PathTypeDir
	return dto
}

// appendCategoryFolders 追加分类文件夹
func (u *UnsAddService) appendCategoryFolders(ctx context.Context, dtos []*types.CreateTopicDto, errorTip map[string]string) []*types.CreateTopicDto {
	// 收集所有非空的 parentAlias
	parentAliasList := base.MapKeys(base.MapArrayToMap[*types.CreateTopicDto, string, bool](dtos, func(dto *types.CreateTopicDto) (ok bool, k string, v bool) {
		pa := dto.ParentAlias
		if pa == nil {
			return false, "", false
		}
		return true, *pa, true
	}))
	var paramParentGroups map[string][]*types.CreateTopicDto

	parentAliasMap := make(map[string]*dao.UnsNamespace)
	categoryFolderMap := make(map[string][]*dao.UnsNamespace)
	db := dao.GetDb(ctx)
	// 如果 parentAliasSet 不为空，查询数据库
	if len(parentAliasList) > 0 {
		// 查询父级 UnsPo 列表
		parentUnsPoList, err := u.unsMapper.ListByAlias(db, parentAliasList)
		if err == nil && len(parentUnsPoList) > 0 {
			for _, unsPo := range parentUnsPoList {
				parentAliasMap[unsPo.Alias] = unsPo
			}
		}

		// 查询归类子文件夹
		categoryFolders, err := u.unsMapper.ListCategoryFolders(db, parentAliasList)
		if err == nil {
			for _, folder := range categoryFolders {
				if folder.ParentAlias != nil {
					parentAlias := *folder.ParentAlias
					categoryFolderMap[parentAlias] = append(categoryFolderMap[parentAlias], folder)
				}
			}
		}
	}

	// 查询根归类文件夹
	rootUnsPoList, err := u.unsMapper.ListRootCategoryFolders(db)
	if err == nil {
		categoryFolderMap[""] = rootUnsPoList // Go 中用空字符串代替 null
	}

	// 创建 alias 到 CreateTopicDto 的映射
	aliasMap := make(map[string]*types.CreateTopicDto)
	for _, dto := range dtos {
		aliasMap[dto.Alias] = dto
	}

	newCategoryAliasMap := make(map[string]*types.CreateTopicDto)

	// 使用新的切片来存储有效元素，而不是在迭代中删除
	validTopicDtos := make([]*types.CreateTopicDto, 0, len(dtos))

	overrideFolders := make(map[string]bool)
	for _, topicDto := range dtos {
		if topicDto.PathType == constants.PathTypeFile && base.P2v(topicDto.DataType) != constants.AlarmRuleType {
			// 验证 parentDataType 是否有效
			if topicDto.ParentDataType == nil || *topicDto.ParentDataType < 1 || *topicDto.ParentDataType > 3 {
				if topicDto.Id == 0 {
					errorTip[topicDto.GainBatchIndex()] = I18nUtils.GetMessage("uns.file.type.invalid")
				} else {
					validTopicDtos = append(validTopicDtos, topicDto)
				}
				continue // 跳过这个元素
			}
			// 判断父级类型和文件类型是否匹配
			if !enums.IsTypeMatched(topicDto.ParentDataType, topicDto.DataType) {
				errorTip[topicDto.GainBatchIndex()] = I18nUtils.GetMessage("uns.category.type.not.eq")
				continue
			}
			parentAlias := base.P2v(topicDto.ParentAlias)
			parentUnsPo := parentAliasMap[parentAlias]
			var PDT = *topicDto.ParentDataType
			// 如果父级目录是归类文件夹并且已经在数据库中存在
			if parentUnsPo != nil && parentUnsPo.DataType != nil {
				// 检查文件类型和父级文件夹类型是否一致
				if *parentUnsPo.DataType > 0 && *parentUnsPo.DataType != PDT {
					errorTip[topicDto.GainBatchIndex()] = I18nUtils.GetMessage("uns.category.type.not.eq")
					continue
				}
				if *parentUnsPo.DataType == PDT {
					validTopicDtos = append(validTopicDtos, topicDto)
					continue
				}
			}

			// 如果数据库不存在则从当前列表中查询
			parentDto := aliasMap[parentAlias]
			// 如果当前列表中已存在父级文件夹，并且类型为归类文件夹
			if parentDto != nil && parentDto.DataType != nil {
				if *parentDto.DataType > 0 && *parentDto.DataType != PDT {
					errorTip[topicDto.GainBatchIndex()] = I18nUtils.GetMessage("uns.category.type.not.eq")
					continue
				}
				if *parentDto.DataType == PDT {
					validTopicDtos = append(validTopicDtos, topicDto)
					continue
				}
			}

			// 文件：判断父节点下面有没有归类子文件夹
			existsCategoryFolders := categoryFolderMap[parentAlias]
			if len(existsCategoryFolders) > 0 && topicDto.ParentDataType != nil {
				bingo := false
				for _, existsCategoryFolder := range existsCategoryFolders {
					if base.P2v(existsCategoryFolder.DataType) == PDT {
						topicDto.ParentAlias = &existsCategoryFolder.Alias
						bingo = true
						break
					}
				}
				if bingo {
					validTopicDtos = append(validTopicDtos, topicDto)
					continue
				}
			}

			// 获取 mountSource 和 mountType
			var mountSource *string
			var mountType *int16

			if parentUnsPo != nil {
				mountSource = parentUnsPo.MountSource
				mountType = parentUnsPo.MountType
			} else if parentDto != nil {
				mountSource = parentDto.MountSource
				mountType = parentDto.MountType
			}

			// 在文件和其父节点之间插入一层归类文件夹
			categoryDto := buildCategoryFolderDto(topicDto.ParentAlias, mountType, mountSource, topicDto.ParentDataType)
			topicDto.ParentAlias = &categoryDto.Alias
			if _, exists := newCategoryAliasMap[categoryDto.Alias]; !exists {
				newCategoryAliasMap[categoryDto.Alias] = categoryDto
			}

			validTopicDtos = append(validTopicDtos, topicDto)
		} else {
			if topicDto.PathType == constants.PathTypeDir && base.P2v(topicDto.DataType) > 0 {
				parentAlias := base.P2v(topicDto.ParentAlias)
				// 目录：判断父节点下面有没有归类子文件夹
				existsCategoryFolders := categoryFolderMap[parentAlias]
				if len(existsCategoryFolders) > 0 {
					dt := *topicDto.DataType
					findSameCategoryDir := false
					for _, existsCategoryFolder := range existsCategoryFolders {
						if base.P2v(existsCategoryFolder.DataType) == dt {
							// 发现同类的分类文件夹alias=ad1，则当前文件夹不保存,且修正直接子节点ParentAlias指向 ad1
							topicDto.Id = existsCategoryFolder.Id
							topicDto.ParentId = existsCategoryFolder.ParentId
							delete(aliasMap, topicDto.Alias)
							aliasMap[existsCategoryFolder.Alias] = UnsConverter.Po2Dto(existsCategoryFolder)
							if paramParentGroups == nil {
								paramParentGroups = base.GroupBy[*types.CreateTopicDto, string](dtos, func(e *types.CreateTopicDto) string {
									return base.P2v(e.ParentAlias)
								})
							}
							children := paramParentGroups[topicDto.Alias]
							if len(children) > 0 {
								for _, child := range children {
									child.ParentAlias = &existsCategoryFolder.Alias
								}
							}
							topicDto.Alias = existsCategoryFolder.Alias
							if !overrideFolders[existsCategoryFolder.Alias] {
								overrideFolders[existsCategoryFolder.Alias] = true
								validTopicDtos = append(validTopicDtos, topicDto)
							}
							findSameCategoryDir = true
							break
						}
					}
					if findSameCategoryDir {
						continue //当前文件夹不保存
					}
				}
			}
			// 非文件类型或告警规则类型，直接保留
			validTopicDtos = append(validTopicDtos, topicDto)
		}
	}

	// 添加新的分类文件夹
	if len(newCategoryAliasMap) > 0 {
		for _, categoryDto := range newCategoryAliasMap {
			validTopicDtos = append(validTopicDtos, categoryDto)
		}
	}

	return validTopicDtos
}
