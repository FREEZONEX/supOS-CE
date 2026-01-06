package connection

type ConnectionResp struct {
	Connections []ConnectionDto `json:"connections"`
}

type ConnectionReq struct {
	Name string `json:"name"`
}

type ConnectionDto struct {
	Name   string         `json:"name"`
	Config map[string]any `json:"config"`
}
