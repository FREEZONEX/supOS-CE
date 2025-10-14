package relationDB

import (
	"gitee.com/unitedrhino/share/conf"
)

var NeedInitColumn bool

func Migrate(c conf.Database) error {
	return nil
	//if c.IsInitTable == false {
	//	return nil
	//}
	//db := stores.GetCommonConn(ctxs.WithRoot(nil))
	//if !db.Migrator().HasTable(&UnsNamespaceNodeInfo{}) {
	//	//需要初始化表
	//	NeedInitColumn = true
	//}
	//err := db.AutoMigrate(
	//	&UnsNamespaceNodeInfo{},
	//	&UnsNamespaceTemplateInfo{},
	//	&UnsNamespaceLabelInfo{},
	//	&UnsNamespaceLabelNodeID{},
	//	&UnsNamespaceNodeVersion{},
	//	&UnsConnectNodeInfo{},
	//	&UnsNoderedFlow{},
	//	&UnsNoderedFlowNode{},
	//	&UnsNamespaceNodeVersion{},
	//)
	//if err != nil {
	//	return err
	//}
	//
	//return err
}

func migrateTableColumn() error {
	return nil
}
