// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package system

import (
	"backend/internal/common/utils/loggerlevel"
	"context"
	"log"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 日志级别配置
func NewLogConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogConfigLogic {
	return &LogConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const (
	levelAlert  = "alert"
	levelInfo   = "info"
	levelError  = "error"
	levelSevere = "severe"
	levelFatal  = "fatal"
	levelSlow   = "slow"
	levelStat   = "stat"
	levelDebug  = "debug"
)

func (l *LogConfigLogic) LogConfig(req *types.LogConfigRequest) (resp *types.LogConfigResponse, err error) {
	if req.Level != "" {
		SetLogLevel(req.Level)
	}
	return &types.LogConfigResponse{Level: loggerlevel.CurrentLevel}, nil
}

func SetLogLevel(level string) bool {
	level = strings.ToLower(strings.TrimSpace(level))
	var levelInt = uint32(0)
	ok := true
	switch level {
	case levelDebug:
		levelInt = logx.DebugLevel
	case levelInfo:
		levelInt = logx.InfoLevel
	case levelError:
		levelInt = logx.ErrorLevel
	case levelSevere:
		levelInt = logx.SevereLevel
	default:
		ok = false
	}
	if ok {
		loggerlevel.CurrentLevel = level
		logx.SetLevel(levelInt)
	}
	log.Println("Set log level to ", loggerlevel.CurrentLevel, ok)
	return ok
}
