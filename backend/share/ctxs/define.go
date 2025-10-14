package ctxs

import "strings"

const (
	UserInfoKey      string = "user"
	UserTokenKey     string = "token"
	UserWorkspaceKey string = "workspace-code" //用户租户号
	UserSetTokenKey  string = "set-token"
)

var HttpAllowHeader string

func init() {
	HttpAllowHeader = "Content-Type, Content-Length,Accept-Language, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, sec-ch-ua, sec-ch-ua-mobile, sec-ch-ua-platform," + strings.Join(ContextKeys, ",")
}
