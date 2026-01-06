package svc

import (
	"backend/internal/common"
	"backend/internal/middleware"
	"backend/internal/repo/relationDB"
	_ "backend/share/result"
	"os"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/i18ns"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"golang.org/x/text/language"

	cache "backend/internal/common/cache"
	"backend/internal/common/constants"
	"backend/internal/config"
	keycloakrepo "backend/internal/repo/keycloak"
	"backend/share/clients"
	noderedclient "backend/share/clients/nodered"
)

type ServiceContext struct {
	Config         config.Config
	InitCtxsWare   rest.Middleware
	CheckTokenWare rest.Middleware
	SnowFlake      *utils.SnowFlake
	OssClient      *oss.Client
	Keycloak       *clients.KeycloakClient
	SourceNodeRed  *noderedclient.Client
	EventNodeRed   *noderedclient.Client
	I18n           map[language.Tag]*i18n.MessageFile
}

func NewServiceContext(c config.Config) *ServiceContext {
	mf, err := i18ns.InitWithFS("etc/i18n")
	logx.Must(err)
	stores.InitConn(c.Database)
	relationDB.Migrate(c.Database)
	if err := cache.InitCaches(); err != nil {
		logx.Errorf("failed to init cache: %v", err)
		panic(err)
	}

	if dbConn, err := stores.GetConn(conf.Database{
		DBType: "pgsql",
		DSN:    c.KeycloakDSN,
	}); err != nil {
		logx.Errorf("failed to init keycloak database: %v", err)
	} else {
		keycloakrepo.InitWithDB(dbConn)
	}

	keycloakClient := clients.InitKeycloakClient(c.OAuthKeyCloak)

	common.InitSnowflake(1)
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:         c,
		OssClient:      ossClient,
		CheckTokenWare: middleware.NewCheckTokenWareMiddleware(keycloakClient, constants.DefaultHomepage, c.OAuthKeyCloak.Realm).Handle,
		InitCtxsWare:   middleware.NewInitCtxsWareMiddleware().Handle,
		SnowFlake:      utils.NewSnowFlake(1),
		I18n:           mf,
		Keycloak:       keycloakClient,
		SourceNodeRed:  noderedclient.NewClient(c.NodeRed.Source.Host, c.NodeRed.Source.Port, ""),
		EventNodeRed:   noderedclient.NewClient(c.NodeRed.Event.Host, c.NodeRed.Event.Port, ""),
	}
}
