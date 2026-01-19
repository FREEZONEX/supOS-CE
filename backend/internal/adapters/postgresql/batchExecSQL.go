package postgresql

import (
	"backend/share/base"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BatchExecSQL(log errLogger, dbPool *pgxpool.Pool, sqlList []string, retry func(er error, index int, sql string) error) error {
	if len(sqlList) == 0 {
		return nil
	}
	log.Debugf("BatchExecSQL: %v", sqlList)

	var wg sync.WaitGroup
	var lock sync.Mutex
	const BATCH_SIZE = 500
	listParts := base.Partition(sqlList, BATCH_SIZE)
	var allErrors []error
	wg.Add(len(listParts))
	addError := func(err error) {
		lock.Lock()
		allErrors = append(allErrors, err)
		lock.Unlock()
	}
	for i, segment := range listParts {
		index := i
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			conn, err := dbPool.Acquire(ctx)
			defer func() {
				if conn != nil {
					conn.Release()
				}
				wg.Done()
			}()
			cancel()

			if err != nil {
				addError(err)
				return
			} else if conn == nil {
				err = fmt.Errorf("conn is nil")
				addError(err)
				return
			}
			var batch = &pgx.Batch{}
			for _, sql := range segment {
				batch.Queue(sql)
			}
			rs := conn.SendBatch(context.Background(), batch)
			for j, sql := range segment {
				_, err = rs.Exec()
				if err != nil {
					if retry != nil {
						err = retry(err, index*BATCH_SIZE+j, sql)
					}
					if err != nil {
						log.Error("Failed to execute sql: ", err, ", sql: ", sql)
						if err != nil {
							addError(err)
						}
					}
				}
			}
		}()
	}
	wg.Wait()
	if len(allErrors) > 0 {
		return allErrors[0]
	}
	return nil
}
func ExecSQLs(log errLogger, dbPool *pgxpool.Pool, sqlList []string) error {
	log.Debugf("ExecSQLs: %+v", sqlList)

	var wg sync.WaitGroup
	var lock sync.Mutex
	listParts := base.Partition(sqlList, 500)
	var allErrors []error
	wg.Add(len(listParts))
	addError := func(err error) {
		lock.Lock()
		allErrors = append(allErrors, err)
		lock.Unlock()
	}
	for _, segment := range listParts {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			conn, err := dbPool.Acquire(ctx)
			defer func() {
				if conn != nil {
					conn.Release()
				}
				wg.Done()
			}()
			cancel()

			if err != nil {
				addError(err)
				return
			} else if conn == nil {
				err = fmt.Errorf("conn is nil")
				addError(err)
				return
			}
			for _, sql := range segment {
				_, err = conn.Exec(context.Background(), sql)
				if err != nil {
					addError(err)

					log.Error("Failed to execute sql: ", err, ", sql: ", sql)
					break
				}
			}
		}()
	}
	wg.Wait()
	if len(allErrors) > 0 {
		return allErrors[0]
	}
	return nil
}
