package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/utils/PathUtil"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"context"
	"strconv"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

const maxCreatePathSegmentLen = 63

var reservedCreateFolderNames = map[string]struct{}{
	"label":    {},
	"template": {},
}

func isAbsoluteCreatePathName(name string) bool {
	return strings.Contains(strings.TrimSpace(name), "/")
}

func isCategoryFolderDataType(dataType *int16) bool {
	return dataType != nil && *dataType > 0
}

func tryReuseExistingCategoryDirLeaf(topicDto *types.CreateTopicDto, existing *dao.UnsNamespace) bool {
	if topicDto.PathType != constants.PathTypeDir || existing == nil || existing.PathType != constants.PathTypeDir {
		return false
	}
	leafDataType, ok := deriveCategoryDataType(topicDto.Name)
	if !ok || existing.DataType == nil || *existing.DataType != *leafDataType {
		return false
	}
	topicDto.Id = existing.Id
	topicDto.Alias = existing.Alias
	topicDto.ParentId = existing.ParentId
	topicDto.ParentAlias = existing.ParentAlias
	topicDto.DataType = existing.DataType
	return true
}

func (u *UnsAddService) resolveCreateParent(
	ctx context.Context,
	db *gorm.DB,
	topicDto *types.CreateTopicDto,
) (parent *dao.UnsNamespace, msg string, err error) {
	if topicDto.ParentId != nil && *topicDto.ParentId == 0 {
		topicDto.ParentId = nil
	}
	switch {
	case topicDto.ParentId != nil && topicDto.ParentAlias == nil:
		folder, queryErr := u.unsMapper.SelectById(db, *topicDto.ParentId)
		if queryErr != nil {
			return nil, "", queryErr
		}
		if folder == nil {
			return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.folder.not.found") + ":parent=" + strconv.FormatInt(*topicDto.ParentId, 10), nil
		}
		topicDto.ParentAlias = &folder.Alias
		return folder, "", nil
	case topicDto.ParentAlias != nil && *topicDto.ParentAlias != "" && topicDto.ParentId == nil:
		folder, queryErr := u.unsMapper.FindOneByAliasOrNil(db, *topicDto.ParentAlias)
		if queryErr != nil {
			return nil, "", queryErr
		}
		if folder == nil {
			return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.folder.not.found") + ":parentAlias=" + *topicDto.ParentAlias, nil
		}
		topicDto.ParentId = &folder.Id
		return folder, "", nil
	}
	return nil, "", nil
}

func validateCreateParent(ctx context.Context, topicDto *types.CreateTopicDto, parent *dao.UnsNamespace) string {
	if topicDto.PathType == constants.PathTypeDir && parent != nil && isCategoryFolderDataType(parent.DataType) {
		return I18nUtils.GetMessageWithCtx(ctx, "uns.category.folder.child.forbidden")
	}
	return ""
}

func (u *UnsAddService) expandCreatePathTopics(
	ctx context.Context,
	db *gorm.DB,
	topicDto *types.CreateTopicDto,
	parentFolder *dao.UnsNamespace,
) (topics []*types.CreateTopicDto, msg string, err error) {
	topicDto.Name = strings.TrimSpace(topicDto.Name)
	if !isAbsoluteCreatePathName(topicDto.Name) {
		return []*types.CreateTopicDto{topicDto}, "", nil
	}

	segments, msg := validateCreatePathSegments(ctx, topicDto.Name, topicDto.PathType)
	if msg != "" {
		return nil, msg, nil
	}

	var categoryDataType *int16
	if topicDto.PathType == constants.PathTypeFile {
		categoryDataType, msg = validateMultiSegmentFilePath(ctx, segments, topicDto.ParentDataType, topicDto.DataType)
		if msg != "" {
			return nil, msg, nil
		}
	}

	currentParentAlias := topicDto.ParentAlias
	currentParentID := topicDto.ParentId
	parentExistsInDB := true
	currentParentFolderType := (*int16)(nil)
	if parentFolder != nil {
		currentParentFolderType = parentFolder.DataType
	}
	folderDtos := make([]*types.CreateTopicDto, 0, len(segments)-1)
	lastDirIndex := len(segments) - 2

	for i := 0; i < len(segments)-1; i++ {
		if topicDto.PathType == constants.PathTypeDir && isCategoryFolderDataType(currentParentFolderType) {
			return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.category.folder.child.forbidden"), nil
		}

		segment := segments[i]
		requiredCategory := topicDto.PathType == constants.PathTypeFile && i == lastDirIndex
		var currentFolderType *int16
		if requiredCategory {
			currentFolderType = categoryDataType
		}

		if parentExistsInDB {
			existing, queryErr := u.unsMapper.FindOneByParentIDAndName(db, currentParentID, segment)
			if queryErr != nil {
				return nil, "", queryErr
			}
			if existing != nil {
				if existing.PathType != constants.PathTypeDir {
					return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.path.segment.not.folder", segment), nil
				}
				if requiredCategory && (existing.DataType == nil || *existing.DataType != *categoryDataType) {
					return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.path.category.folder.invalid", segment), nil
				}
				currentParentID = &existing.Id
				currentParentAlias = &existing.Alias
				currentParentFolderType = existing.DataType
				continue
			}
		}

		folderDto := buildCreatePathFolderDto(segment, currentParentAlias, currentParentID, currentFolderType)
		folderDtos = append(folderDtos, folderDto)
		currentParentID = &folderDto.Id
		currentParentAlias = &folderDto.Alias
		currentParentFolderType = folderDto.DataType
		parentExistsInDB = false
	}

	if topicDto.PathType == constants.PathTypeDir && isCategoryFolderDataType(currentParentFolderType) {
		return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.category.folder.child.forbidden"), nil
	}

	topicDto.Name = segments[len(segments)-1]
	if topicDto.PathType == constants.PathTypeDir {
		existingLeaf, queryErr := u.unsMapper.FindOneByParentIDAndName(db, currentParentID, topicDto.Name)
		if queryErr != nil {
			return nil, "", queryErr
		}
		if tryReuseExistingCategoryDirLeaf(topicDto, existingLeaf) {
			return folderDtos, "", nil
		}
	}
	topicDto.ParentAlias = currentParentAlias
	topicDto.ParentId = currentParentID
	if categoryDataType != nil {
		topicDto.ParentDataType = categoryDataType
	}

	return append(folderDtos, topicDto), "", nil
}

func buildCreatePathFolderDto(name string, parentAlias *string, parentID *int64, dataType *int16) *types.CreateTopicDto {
	if dataType != nil && *dataType > 0 {
		dto := buildCategoryFolderDto(parentAlias, nil, nil, dataType)
		dto.ParentId = parentID
		return dto
	}
	return &types.CreateTopicDto{
		Id:          common.NextId(),
		Name:        name,
		Alias:       PathUtil.GenerateAliasWithRandom(name),
		ParentAlias: parentAlias,
		ParentId:    parentID,
		PathType:    constants.PathTypeDir,
	}
}

func validateCreatePathSegments(ctx context.Context, name string, pathType int16) ([]string, string) {
	msgKey := "uns.topic.format.invalid"
	if pathType == constants.PathTypeDir {
		msgKey = "uns.folder.format.invalid"
	}
	if !PathUtil.ValidTopicFormat(name, nil) {
		return nil, I18nUtils.GetMessageWithCtx(ctx, msgKey)
	}

	segments := strings.Split(name, "/")
	lastDirIndex := len(segments) - 2
	for i, segment := range segments {
		if utf8.RuneCountInString(segment) > maxCreatePathSegmentLen {
			return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.path.segment.too.long", segment, maxCreatePathSegmentLen)
		}
		if pathType == constants.PathTypeDir || i <= lastDirIndex {
			if _, exists := reservedCreateFolderNames[strings.ToLower(segment)]; exists {
				return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.prohibitKeywords")
			}
		}
		if pathType == constants.PathTypeDir && i < len(segments)-1 {
			if _, ok := deriveCategoryDataType(segment); ok {
				return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.category.folder.child.forbidden")
			}
		}
	}
	return segments, ""
}

func validateMultiSegmentFilePath(
	ctx context.Context,
	segments []string,
	parentDataType *int16,
	dataType *int16,
) (*int16, string) {
	categoryDataType, ok := deriveCategoryDataType(segments[len(segments)-2])
	if !ok {
		return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.path.category.required")
	}
	for _, segment := range segments[:len(segments)-2] {
		if _, ok := deriveCategoryDataType(segment); ok {
			return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.category.folder.child.forbidden")
		}
	}
	if parentDataType != nil && *parentDataType != *categoryDataType {
		return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.category.type.not.eq")
	}
	if !enums.IsTypeMatched(categoryDataType, dataType) {
		return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.file.type.invalid")
	}
	return categoryDataType, ""
}

func deriveCategoryDataType(segment string) (*int16, bool) {
	folderType, ok := enums.GetFolderDataTypeByName(segment)
	if !ok || folderType == enums.NORMAL {
		return nil, false
	}
	dataType := folderType.TypeIndex()
	return &dataType, true
}
