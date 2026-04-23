package types

type OpenSourceFlowCreateReq = SourceFlowCreateReq

type OpenSourceFlowListQuery = SourceFlowListQuery

type OpenSourceFlowUpdateReq = SourceFlowUpdateReq

type OpenSourceFlowDataQuery struct {
	ID string `form:"id"`
}

type OpenSourceFlowDataResp struct {
	Flows []map[string]interface{} `json:"flows,optional"`
	Rev   string                   `json:"rev,optional"`
}

type OpenSourceFlowDeployReq = SourceFlowDeployReq

type OpenSourceFlowDeleteReq = SourceFlowDeleteReq

type OpenSourceFlowInfo = SourceFlowInfo

type OpenSourceFlowPageResult = SourceFlowPageResult

type OpenSourceFlowDeployResult = SourceFlowDeployResult
