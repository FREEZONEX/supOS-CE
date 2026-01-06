package system

import (
	"backend/internal/logic/supos/uns/system"
	"encoding/json"
	"net/http"
)

func DevtestHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	var params map[string][]string = r.Form
	rs := system.Devtest(r.Context(), params)
	w.Header().Set("content-type", "application/json")
	bs, _ := json.Marshal(rs)
	_, _ = w.Write(bs)
}
