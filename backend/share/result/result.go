package result

import (
	"backend/internal/common/utils/apiutil"
	"backend/share/ctxs"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func init() {
	result.Http = func(w http.ResponseWriter, r *http.Request, resp any, err error) {
		var code int
		var msg string
		if err == nil {
			//成功返回
			re := result.Success(resp)
			httpx.WriteJson(w, http.StatusOK, re)
			code = 200
			msg = "success"

		} else {
			//错误返回
			er := errors.Fmt(err)
			u := apiutil.GetUserFromContext(r.Context())
			var lang string
			if u != nil {
				lang = u.MainLanguage
			}
			msg = er.GetI18nMsg(lang)

			logx.WithContext(r.Context()).Errorf("【http handle err】router:%v err: %#v ",
				r.URL.Path, err)
			httpx.WriteJson(w, http.StatusOK, result.Error(er.Code, msg))
			code = int(er.Code)
		}
		ret := ctxs.GetResp(r)
		if ret != nil {
			//将接口的应答结果写入r.Response，为操作日志记录接口提供应答信息
			var temp http.Response
			temp.StatusCode = code
			temp.Status = msg
			if resp != nil {
				bs, _ := json.Marshal(resp)
				temp.Body = ioutil.NopCloser(bytes.NewReader(bs))
			}
			*ret = temp
		}
	}
}
