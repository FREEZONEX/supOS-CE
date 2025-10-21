package fuxautil

import (
	"bytes"
	"io"
	"net/http"

	"backend/internal/common/utils/runtimeutil"

	"github.com/zeromicro/go-zero/core/logx"
)

// FuxaUtils provides utility functions for FUXA operations.
//
// GetFuxaURL returns the FUXA URL based on the runtime environment.
func GetFuxaURL() string {
	if runtimeutil.IsLocalEnv() {
		return "http://100.100.100.22:33893/fuxa/home"
	}
	return "http://fuxa:1881"
}

// Create creates a FUXA dashboard from JSON.
func Create(dashboardJSON string) (bool, error) {
	logx.Debugf(">>>>>>>>>>>>>>>fuxa 创建 dashboards 请求: %s", dashboardJSON)

	resp, err := http.Post(GetFuxaURL()+"/api/project", "application/json", bytes.NewBufferString(dashboardJSON))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf(">>>>>>>>>>>>>>>fuxa 创建 dashboards 返回结果: %s", string(body))

	return resp.StatusCode == 200, nil
}

// Get retrieves a FUXA dashboard by layout ID.
func Get(layoutID string) (string, error) {
	url := GetFuxaURL() + "/api/project?layoutId=" + layoutID
	logx.Debugf(">>>>>>>>>>>>>>>fuxa 查询 dashboards 请求: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf(">>>>>>>>>>>>>>>fuxa 查询 dashboards 返回结果: %s", string(body))

	if resp.StatusCode != 200 {
		return "", nil
	}
	return string(body), nil
}
