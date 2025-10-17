package config

import (
	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database       conf.Database
	DatabaseSchema string `json:",default=supos,env=dbSchema"` //
	//Event      conf.EventConf
	//CacheRedis cache.ClusterConf
	//DevLink    conf.DevLinkConf //和设备交互的设置
	//OssConf    conf.OssConf     `json:",optional"`
}
