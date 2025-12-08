package timescaledb

import (
	"backend/internal/adapters/postgresql"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

type TsdbPersistentService struct {
	log           logx.Logger
	dbPool        *pgxpool.Pool
	currentSchema string
	batchSize     int
	maxRetries    int
	urlProperties serviceApi.DataSourceProperties
}

func init() {
	spring.RegisterLazy[*TsdbPersistentService](func() *TsdbPersistentService {
		svCtx := spring.GetBean[*svc.ServiceContext]()
		url, has := svCtx.Config.PersistentUrl["timescaledb"]
		ctx := context.Background()
		log := logx.WithContext(ctx)
		if !has || len(url) == 0 {
			log.Info("timescaledb url not found in config")
			return nil
		}
		pool, er := pgxpool.New(ctx, url)
		if er != nil {
			log.Error("timescaledb init fail", er)
			return nil
		}
		log.Info("timescaledb init success, url: ", url)
		return &TsdbPersistentService{
			log:           log,
			dbPool:        pool,
			currentSchema: "public",
			batchSize:     8000,
			maxRetries:    1,
			urlProperties: postgresql.ParseDbUrlProperties(url),
		}
	})
}

const dsId = types.SrcJdbcTypeTimeScaleDB

func (p *TsdbPersistentService) Persistent(unsData []serviceApi.UnsData) {
	err := persistence(p.dbPool, p.currentSchema, p.batchSize, unsData)
	if err != nil {
		p.log.Error("persistence fail", err, len(unsData))
	}
}
func (p *TsdbPersistentService) GetDataSourceProperties() (rs serviceApi.DataSourceProperties) {
	return p.urlProperties
}
func (p *TsdbPersistentService) GetDataSrcId() types.SrcJdbcType {
	return dsId
}
func (p *TsdbPersistentService) OnEventBatchCreateTableEvent9(evt *event.BatchCreateTableEvent) error {
	return postgresql.OnCreate(p.log, p.dbPool, p.currentSchema, dsId, evt)
}
func (p *TsdbPersistentService) OnEventUpdateInstanceEvent9(evt *event.UpdateInstanceEvent) error {
	return postgresql.OnUpdate(p.log, p.dbPool, p.currentSchema, dsId, evt)
}
func (p *TsdbPersistentService) OnEventRemoveTopicsEvent9(evt *event.RemoveTopicsEvent) {
	postgresql.OnRemove(p.log, p.dbPool, dsId, evt)
}
func (p *TsdbPersistentService) FillLastRecord(uns *types.CreateTopicDto) {
	sql := base.StringBuilder{}
	sql.Grow(256)
	sql.Append(`SELECT * FROM "`).Append(uns.GetTable()).Append(`" `)
	if tbF := uns.GetTbFieldName(); tbF != "" {
		sql.Append(`WHERE "`).Append(tbF).Append(`"= `).Long(uns.Id)
	}
	sql.Append(` ORDER BY "`).Append(uns.GetTimestampField()).Append(`" DESC LIMIT 1`)
	rows, er := p.dbPool.Query(context.Background(), sql.String())
	if er != nil {
		p.log.Error("fail to get last record", er, uns.GetAlias())
		return
	}
	valuePointers := make([]interface{}, len(rows.FieldDescriptions()))
	fd := uns.GetFieldDefines().FieldsMap
	for i, f := range rows.FieldDescriptions() {
		if field, has := fd[f.Name]; has {
			field.LastValue = types.FieldType(field.Type).ZeroValue()
			valuePointers[i] = &field.LastValue
		} else {
			var obj any
			valuePointers[i] = &obj
		}
	}
	if rows.Next() {
		er = rows.Scan(valuePointers...)
		if er != nil {
			p.log.Error("fail to Scan last record", er, uns.GetAlias())
		}
	}
}
