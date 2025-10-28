package common

import "backend/internal/types"

type CommonMountSourceDto struct {
	Alias      string `json:"alias"`
	Name       string `json:"name"`
	SourceType string `json:"sourceType"`
}

type CommonFolderMetaDto struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	MountType   int    `json:"mountType"`
	MountSource string `json:"mountSource"`
}

type CommonFileMetaDto struct {
	Alias         string              `json:"alias"`
	Name          string              `json:"name"`
	DisplayName   string              `json:"displayName"`
	Description   string              `json:"description"`
	SaveToDB      *bool               `json:"save2db"`
	OriginalAlias string              `json:"originalAlias"`
	Fields        []types.FieldDefine `json:"fields"`
}
