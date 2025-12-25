package postgresql

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/dbpool"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

type PgPersistentService struct {
	log           logx.Logger
	dbPool        *pgxpool.Pool
	currentSchema string
	batchSize     int
	maxRetries    int

	urlProperties serviceApi.DataSourceProperties
}

const dsId = types.SrcJdbcTypePostgresql

func init() {
	spring.RegisterLazy[*PgPersistentService](func() *PgPersistentService {
		svCtx := spring.GetBean[*svc.ServiceContext]()
		pgUrl, has := svCtx.Config.PersistentUrl["postgresql"]
		ctx := context.Background()
		log := logx.WithContext(ctx)
		if !has || len(pgUrl) == 0 {
			log.Info("postgresql url not found in config")
			return nil
		}
		pool, er := dbpool.NewPool(ctx, pgUrl, "uns_relation")
		if er != nil {
			log.Error("postgresql init fail", er)
			return nil
		}
		log.Info("postgresql init success, url: ", pgUrl)
		return &PgPersistentService{
			log:           log,
			dbPool:        pool,
			currentSchema: "public",
			batchSize:     200,
			maxRetries:    1,
			urlProperties: ParseDbUrlProperties(pgUrl),
		}
	})
}
func (p *PgPersistentService) Persistent(unsData []serviceApi.UnsData) {
	err := persistence(p.dbPool, p.currentSchema, p.batchSize, unsData)
	if err != nil {
		p.log.Error("persistence fail", err, len(unsData))
	}
}
func (p *PgPersistentService) GetDataSourceProperties() (rs serviceApi.DataSourceProperties) {
	return p.urlProperties
}
func (p *PgPersistentService) GetDataSrcId() types.SrcJdbcType {
	return dsId
}
func (p *PgPersistentService) OnEventBatchCreateTableEvent7(evt *event.BatchCreateTableEvent) error {
	return OnCreate(p.log, p.dbPool, p.currentSchema, dsId, evt)
}
func (p *PgPersistentService) OnEventUpdateInstanceEvent7(evt *event.UpdateInstanceEvent) error {
	return OnUpdate(p.log, p.dbPool, p.currentSchema, dsId, evt)
}
func (p *PgPersistentService) OnEventRemoveTopicsEvent7(evt *event.RemoveTopicsEvent) {
	OnRemove(p.log, p.dbPool, dsId, evt)
}

func (p *PgPersistentService) FillLastRecord(uns *types.CreateTopicDto) {
	FillLastRecord(p.log, uns, func(ctx context.Context) (pgx.Rows, error) {
		return p.dbPool.Query(ctx, fmt.Sprintf(`select * from "%s" ORDER BY "%s" DESC LIMIT 1`, uns.GetTable(), uns.GetTimestampField()))
	})
}
func FillLastRecord(log errLogger, uns *types.CreateTopicDto, query func(ctx context.Context) (pgx.Rows, error)) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, er := query(ctx)
	if rows != nil {
		defer rows.Close()
	}
	if er != nil {
		log.Error("fail to get last record", er, uns.GetAlias())
		return
	}

	if rows.Next() {
		values, er := rows.Values()
		if er != nil {
			log.Error("fail to Scan last record", er, uns.GetAlias())
		} else {
			ct := uns.GetTimestampField()
			var updateTime int64
			for i, f := range rows.FieldDescriptions() {
				if f.Name == ct {
					if tm, is := values[i].(time.Time); is {
						updateTime = tm.UnixMilli()
					}
					break
				}
			}
			fd := uns.GetFieldDefines().FieldsMap
			for i, f := range rows.FieldDescriptions() {
				if field, has := fd[f.Name]; has {
					field.LastValue = values[i]
					field.LastTime = updateTime
				}
			}
		}
	}
}

type errLogger interface {
	Error(...any)
}

func OnCreate(log errLogger, dbPool *pgxpool.Pool, defaultSchema string, dataSrcId types.SrcJdbcType, evt *event.BatchCreateTableEvent) error {
	creates := evt.GetCreateFiles(dataSrcId)
	undates := evt.GetUpdateFiles(dataSrcId)
	if len(undates) > 0 {
		if len(creates) > 0 {
			creates = append(creates, undates...)
		} else {
			creates = undates
		}
	}
	if len(creates) > 0 {
		creates = base.Filter(creates, func(e *types.CreateTopicDto) bool {
			return base.P2v(e.DataType) != constants.AlarmRuleType
		})
		conn, er := dbPool.Acquire(context.Background())
		if er != nil {
			logPoolError("OnCreate:"+dataSrcId.Alias(), time.Time{}, dbPool, "getConn", er)
			log.Error("ListTableInfos fail", er)
			return er
		}
		defer conn.Release()
		tableInfoMap, er := ListTableInfos(conn, creates)
		if er != nil {
			log.Error("ListTableInfos fail", er)
			return er
		}
		Errors := BatchCreateTables(conn, defaultSchema, creates, tableInfoMap)
		if len(Errors) > 0 {
			log.Error("BatchCreateTables fail", Errors)
			return Errors[0]
		}
	}
	return nil
}
func OnUpdate(log errLogger, dbPool *pgxpool.Pool, defaultSchema string, dataSrcId types.SrcJdbcType, evt *event.UpdateInstanceEvent) error {
	topicList := base.Filter(evt.Topics, func(e *types.CreateTopicDto) bool {
		return e.FieldsChanged && e.GetSrcJdbcType() == dataSrcId
	})
	if len(topicList) > 0 {
		conn, er := dbPool.Acquire(context.Background())
		if er != nil {
			logPoolError("OnUpdate:"+dataSrcId.Alias(), time.Time{}, dbPool, "getConn", er)
			log.Error("OnUpdate fail", er)
			return er
		}
		defer conn.Release()
		tableInfoMap, er := ListTableInfos(conn, topicList)
		if er != nil {
			log.Error("ListTableInfos fail", er)
			return er
		}
		Errors := BatchCreateTables(conn, defaultSchema, topicList, tableInfoMap)
		if len(Errors) > 0 {
			log.Error("BatchUpdate fail", Errors)
			return Errors[0]
		}
	}
	return nil
}
func OnRemove(log errLogger, dbPool *pgxpool.Pool, dataSrcId types.SrcJdbcType, evt *event.RemoveTopicsEvent) {
	batch := &pgx.Batch{}
	for _, e := range evt.Topics {
		if e.GetSrcJdbcType() == dataSrcId {
			tbf := e.GetTbFieldName()
			sql := ""
			table := getFullTableName(e.GetTable())
			if tbf == "" {
				sql = "drop table if exists " + table
			} else {
				sql = fmt.Sprintf("delete from %s where %s=%v", table, tbf, e.Id)
			}
			batch.Queue(sql)
		}
	}
	if len(batch.QueuedQueries) > 0 {
		conn, er := dbPool.Acquire(context.Background())
		if er != nil {
			logPoolError("OnRemove:"+dataSrcId.Alias(), time.Time{}, dbPool, "getConn", er)
			log.Error("OnRemove fail", er)
			return
		}
		defer conn.Release()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		br := conn.SendBatch(ctx, batch)
		defer cancel()
		for i := 0; i < batch.Len(); i++ {
			_, err := br.Exec()
			if err != nil {
				log.Error("delete fail", err, batch.QueuedQueries[i].SQL)
			}
		}
	}
}
