package timescaledb

import (
	"backend/internal/adapters/postgresql"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/dbpool"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		url := svCtx.Config.TimescaledbUrl
		if len(url) == 0 {
			logx.Info("timescaledb url not found in config")
			return nil
		}
		return newTsdbPersistentService(url)
	})
}
func newTsdbPersistentService(url string) *TsdbPersistentService {
	ctx := context.Background()
	log := logx.WithContext(ctx)
	pool, er := dbpool.NewPool(ctx, url, "uns_tsDB")
	if er != nil {
		log.Error("timescaledb init fail", er)
		return nil
	}
	log.Info("timescaledb init success, url: ", url)
	return &TsdbPersistentService{
		log:           log,
		dbPool:        pool,
		currentSchema: "public",
		batchSize:     1000,
		maxRetries:    1,
		urlProperties: postgresql.ParseDbUrlProperties(url),
	}
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
func (p *TsdbPersistentService) FillLastRecord(uns *types.CreateTopicDto) {
	query := func(ctx context.Context) (pgx.Rows, error) {
		sql := base.StringBuilder{}
		sql.Grow(256)
		sql.Append(`SELECT * FROM "`).Append(uns.Alias).Append(`" `)
		sql.Append(` ORDER BY "`).Append(uns.GetTimestampField()).Append(`" DESC LIMIT 1`)
		return p.dbPool.Query(ctx, sql.String())
	}
	postgresql.FillLastRecord(p.log, uns, query)
}

//	func (p *TsdbPersistentService) OnEventBatchCreateTableEvent9(evt *event.BatchCreateTableEvent) error {
//		creates := evt.GetCreateFiles(dsId)
//		undates := evt.GetUpdateFiles(dsId)
//		if len(undates) > 0 {
//			if len(creates) > 0 {
//				creates = append(creates, undates...)
//			} else {
//				creates = undates
//			}
//		}
//		return p.Save(creates)
//	}
//
//	func (p *TsdbPersistentService) OnEventUpdateInstanceEvent9(evt *event.UpdateInstanceEvent) error {
//		topicList := base.Filter(evt.Topics, func(e *types.CreateTopicDto) bool {
//			return e.GetSrcJdbcType() == dsId
//		})
//		return p.Save(topicList)
//	}
func (p *TsdbPersistentService) OnEventRemoveTopicsEvent9(evt *event.RemoveTopicsEvent) {
	p.Remove(base.Map[*types.CreateTopicDto, types.UnsInfo](evt.Topics, func(e *types.CreateTopicDto) types.UnsInfo {
		return e
	}))
}

// Save 保存 UNS
func (p *TsdbPersistentService) Save(topics []types.UnsInfo) error {
	if len(topics) == 0 {
		return nil
	}
	topics = base.Filter(topics, func(e types.UnsInfo) bool {
		return base.P2v(e.GetDataType()) != constants.AlarmRuleType
	})
	ctx := context.Background()
	viewNames := base.Map[types.UnsInfo, string](topics, func(e types.UnsInfo) string {
		return e.GetAlias()
	})
	views, err := parseUnsViews(p.dbPool, ctx, p.currentSchema, viewNames)
	if err != nil {
		return err
	}
	physicsTableFields, err := getPhysicalTableFields(ctx, p.dbPool)
	if err != nil {
		return err
	}
	// 创建 SQL 生成器
	sqlGen := NewSQLGenerator()
	unsList := make([]UnsViewInfo, len(topics))
	for i, create := range topics {
		unsList[i] = UnsViewInfo{Uns: create, View: views[create.GetAlias()]}
	}
	s := sqlGen.GenerateSyncSQLs(physicsTableFields, unsList)
	if createSQLs := s.CreateTableSQL; len(createSQLs) > 0 {
		// 1. 创建物理表
		err = postgresql.ExecSQLs(p.log, p.dbPool, createSQLs)
		if err != nil {
			return err
		}
	}
	var orderedSQLs []string
	// 2. 修改表
	orderedSQLs = append(orderedSQLs, s.AlterTableSQL...)
	// 3. 更新数据
	orderedSQLs = append(orderedSQLs, s.UpdateDataSQL...)
	err = postgresql.BatchExecSQL(p.log, p.dbPool, orderedSQLs, nil)
	if err != nil {
		return err
	}
	retryCreateView := func(er error, index int, sql string) error {
		logx.Errorf("CreateViewErr: %v, index=%d, sql=%s\n", er, index, sql)
		var pgEr *pgconn.PgError
		if errors.As(er, &pgEr) && pgEr.Code == "42809" {
			sqlUp := strings.ToUpper(sql)
			st := strings.Index(sqlUp, "VIEW")
			end := -1
			if st > 0 {
				for i := st + 5; i <= len(sqlUp); i++ {
					if unicode.IsSpace(rune(sqlUp[i])) {
						end = i
						break
					}
				}
			}
			if end > st {
				viewName := strings.TrimSpace(sql[st+5 : end])
				dot := strings.Index(viewName, ".")
				alias := viewName
				if dot > 0 {
					alias = viewName[dot+1:]
				}
				matches := validNamePattern.FindAllStringSubmatch(alias, 1)
				if len(matches) > 0 {
					alias = matches[0][0]
				}
				uns := unsList[index].Uns
				logx.Infof("%v ~准备迁移物理表: %s, id=%d, path=%s", alias == uns.GetAlias(), alias, uns.GetId(), uns.GetPath())
				if alias == uns.GetAlias() {
					var execErr error
					// ALTER TABLE public.test_ds RENAME TO test_ds1;
					_, execErr = p.dbPool.Exec(context.Background(),
						fmt.Sprintf(`ALTER TABLE %s."%s" RENAME TO "bk_%s"`, p.currentSchema, alias, alias))
					if execErr != nil {
						logx.Errorf("%s rename fail %v", alias, execErr)
					} else {
						_, execErr = p.dbPool.Exec(context.Background(), sql)
						if execErr != nil {
							logx.Errorf("%s 重建视图失败 %v", alias, execErr)
						}
					}
					if execErr == nil {
						mSql := base.StringBuilder{}
						mSql.Grow(200)
						mSql.Append(`INSERT INTO "`).Append(uns.GetTable()).Append(`"(`)
						for _, f := range uns.GetFields() {
							name := f.Name
							if idx := f.Index; idx != nil && len(*idx) > 0 {
								name = *idx
							}
							mSql.Append(`"` + name + `"`).Append(",")
						}
						mSql.SetLast(')')
						mSql.Append(` SELECT `)
						for i, f := range uns.GetFields() {
							if i > 0 {
								mSql.Append(`,`)
							}
							if f.TbValueName == nil {
								mSql.Append(`"` + f.Name + `"`)
								if idx := f.Index; idx != nil && len(*idx) > 0 {
									mSql.Append(` AS `).Append(*idx)
								}
							} else {
								mSql.Long(uns.GetId()).Append(" AS ").Append(*f.TbValueName)
							}
						}
						mSql.Append(` FROM "bk_`).Append(alias).Append(`"`)
						mSql.Append(` ON CONFLICT DO NOTHING`)
						migSql := mSql.String()
						_, migrateErr := p.dbPool.Exec(context.Background(), migSql)
						if migrateErr != nil {
							logx.Errorf("数据迁移失败： %s %v, sql= %v", alias, migrateErr, migSql)
						} else {
							_, dropErr := p.dbPool.Exec(context.Background(), `drop table if exists "bk_`+alias+`"`)
							if dropErr != nil {
								logx.Errorf("数据迁移 drop 空表失败： %s %v", alias, dropErr)
							}
						}
					}
					return execErr
				}
			}
		}
		return er
	}
	// 4. 创建视图, 当同名的表已存在时 需要处理冲突
	return postgresql.BatchExecSQL(p.log, p.dbPool, s.CreateViewSQL, retryCreateView)
}

func (p *TsdbPersistentService) Remove(topics []types.UnsInfo) error {
	sqls := make([]string, 0, len(topics))
	countIds := 0
	idsStr := base.StringBuilder{}
	idsStr.Grow(128)
	for _, e := range topics {
		if e.GetSrcJdbcType() == dsId {
			tbf := e.GetTbFieldName()
			table := postgresql.GetFullTableName(e.GetTable())
			if tbf != "" {
				if countIds == 0 {
					idsStr.Append("delete from ").Append(table).Append(" where ").Append(tbf).Append(" IN (")
				}
				countIds++
				idsStr.Long(e.GetId()).Append(",")
				if countIds == 500 {
					idsStr.SetLast(')')
					sqls = append(sqls, idsStr.String())
					countIds = 0
					idsStr.Reset()
				}
				sqls = append(sqls, "drop view if exists \""+e.GetAlias()+"\"")
			} else {
				sqls = append(sqls, "drop table if exists "+table)
			}
		}
	}
	if countIds > 0 {
		idsStr.SetLast(')')
		sqls = append(sqls, idsStr.String())
	}
	if len(sqls) > 0 {
		return postgresql.BatchExecSQL(p.log, p.dbPool, sqls, nil)
	}
	return nil
}
