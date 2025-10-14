package ctxs

import (
	"context"
	"encoding/json"
	"net/http"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type UserCtx struct {
	WorkspaceCode         string `json:"workspaceCode,omitempty"` //目前选择的空间编码
	AcceptLanguage        string `json:"acceptLanguage,omitempty"`
	IsSuperAdmin          bool   `json:"isSuperAdmin,omitempty"`          //全局超管
	IsAdmin               bool   `json:"isAdmin,omitempty"`               //全局管理员
	IsWorkspaceSuperAdmin bool   `json:"isWorkspaceSuperAdmin,omitempty"` //工作空间超管,拥有者才能是超管
	IsWorkspaceAdmin      bool   `json:"isWorkspaceAdmin,omitempty"`      //工作空间管理员
	UserID                int64  `json:"userID,string"`                   //用户id
	UserName              string `json:"userName"`
	Email                 string `json:"email"`
	IP                    string `json:"ip,omitempty"` //用户的ip地址
	InnerCtx
}

func (u *UserCtx) ClearInner() *UserCtx {
	if u == nil {
		return nil
	}
	newCtx := *u
	newCtx.InnerCtx = InnerCtx{}
	return &newCtx
}

type InnerCtx struct {
}

func GetHandle(r *http.Request, keys ...string) string {
	var val string
	for _, v := range keys {
		val = r.Header.Get(v)
		if val != "" {
			return val
		}
		val = r.URL.Query().Get(v)
		if val != "" {
			return val
		}
	}
	return val
}

func BindWorkspaceCode(ctx context.Context, workspaceCode string) context.Context {
	uc := GetUserCtx(ctx)
	if uc == nil {
		uc = &UserCtx{
			WorkspaceCode: workspaceCode,
		}
		ctx = context.WithValue(ctx, UserInfoKey, uc)
	} else {
		uc.WorkspaceCode = workspaceCode
	}
	ctx = stores.SetTenantCode(ctx, workspaceCode)
	return ctx
}

func UpdateUserCtx(ctx context.Context) context.Context {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return ctx
	}
	return SetUserCtx(ctx, uc)
}

type ctxStr struct {
	UserCtx *UserCtx
	Trace   string
}

func ToString(ctx context.Context) string {
	uc := GetUserCtx(ctx)
	span := trace.SpanFromContext(ctx)
	traceinfo, _ := span.SpanContext().MarshalJSON()
	ctxstr := ctxStr{
		UserCtx: uc,
		Trace:   string(traceinfo),
	}
	return utils.MarshalNoErr(ctxstr)
}

type mySpanContextConfig struct {
	TraceID string
	SpanID  string
}

func StringParse(ctx context.Context, str string) (context.Context, bool) {
	var cs ctxStr
	err := json.Unmarshal([]byte(str), &cs)
	if err != nil {
		return ctx, false
	}
	var msg mySpanContextConfig
	err = json.Unmarshal([]byte(cs.Trace), &msg)
	if err != nil {
		logx.Errorf("[GetCtx]|json Unmarshal trace.SpanContextConfig err:%v", err)
		return ctx, false
	}
	//将MsgHead 中的msg链路信息 重新注入ctx中并返回
	t, err := trace.TraceIDFromHex(msg.TraceID)
	if err != nil {
		logx.Errorf("[GetCtx]|TraceIDFromHex err:%v", err)
		return ctx, false
	}
	s, err := trace.SpanIDFromHex(msg.SpanID)
	if err != nil {
		logx.Errorf("[GetCtx]|SpanIDFromHex err:%v", err)
		return ctx, false
	}
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    t,
		SpanID:     s,
		TraceFlags: 0x1,
	})
	ctx2 := trace.ContextWithRemoteSpanContext(ctx, parent)
	return SetUserCtx(ctx2, cs.UserCtx), true

}

func SetUserCtx(ctx context.Context, userCtx *UserCtx) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if userCtx == nil {
		return ctx
	}
	info, _ := json.Marshal(userCtx)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		UserInfoKey, string(info),
	))
	ctx = stores.SetTenantCode(ctx, userCtx.WorkspaceCode)
	return context.WithValue(ctx, UserInfoKey, userCtx)
}

func SetInnerCtx(ctx context.Context, inner InnerCtx) context.Context {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return ctx
	}
	uc.InnerCtx = inner
	return SetUserCtx(ctx, uc)
}

func GetInnerCtx(ctx context.Context) InnerCtx {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return InnerCtx{}
	}
	return uc.InnerCtx
}

// 使用该函数前必须传了UserCtx
func GetUserCtx(ctx context.Context) *UserCtx {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return nil
	}
	return val
}

func GetUserCtxNoNil(ctx context.Context) *UserCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return &UserCtx{}
	}
	return val
}

func WithRoot(ctx context.Context) context.Context {
	uc := *GetUserCtxNoNil(ctx)
	uc.IsSuperAdmin = true
	uc.IsAdmin = true
	return SetUserCtx(ctx, &uc)
}

func WithAdmin(ctx context.Context) context.Context {
	uc := *GetUserCtxNoNil(ctx)
	uc.IsAdmin = true
	return SetUserCtx(ctx, &uc)
}

// 如果是default租户直接给root权限
func WithDefaultRoot(ctx context.Context) context.Context {
	uc := GetUserCtxNoNil(ctx)
	if uc.WorkspaceCode != "" || !uc.IsAdmin { //传了项目ID则不是root权限
		return ctx
	}
	return WithRoot(ctx)
}

func IsAdmin(ctx context.Context) error {
	uc := GetUserCtxNoNil(ctx)
	if uc.IsAdmin || uc.IsSuperAdmin {
		return nil
	}
	return errors.Permissions.AddMsg("只允许管理员操作")
}

func NewUserCtx(ctx context.Context) context.Context {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return ctx
	}

	var newUc UserCtx
	newUc = *val
	newCtx := context.WithValue(context.Background(), UserInfoKey, &newUc)
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		newCtx = metadata.NewOutgoingContext(newCtx, md)
	}
	return newCtx
}

func IsRoot(ctx context.Context) error {
	uc := GetUserCtx(ctx)
	if uc == nil || uc.WorkspaceCode != "" || !uc.IsAdmin {
		return errors.Permissions.AddDetailf("需要超管才能操作")
	}
	return nil
}

// 使用该函数前必须传了UserCtx
func GetUserCtxOrNil(ctx context.Context) *UserCtx {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return nil
	}
	return val
}

//// 指定项目id（企业版功能）
//func SetMetaProjectID(ctx context.Context, projectID int64) {
//	mc := GetMetaCtx(ctx)
//	projectIDStr := utils.ToString(projectID)
//	mc[string(MetaFieldProjectID)] = []string{projectIDStr}
//}
//
//// 获取meta里的项目ID（企业版功能）
//func ClearMetaProjectID(ctx context.Context) {
//	mc := GetMetaCtx(ctx)
//	delete(mc, string(MetaFieldProjectID))
//}
