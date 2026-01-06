package dbpool

import (
	"context"
	"io"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Exporter interface {
	CopyTo(ctx context.Context, w io.Writer, sql string) error
	Release()
}

type ExporterPool struct {
	pool *pgxpool.Pool
}

// NewExporterPool 让应用层导出业务和pgx解耦
func NewExporterPool(ctx context.Context, connString, appName string) (*ExporterPool, error) {
	pool, err := NewPool(ctx, connString, appName)
	if err != nil {
		return nil, err
	}
	return &ExporterPool{pool: pool}, nil
}
func (p *ExporterPool) Acquire(ctx context.Context) (Exporter, error) {
	c, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return pgExporter{conn: c}, nil
}

type pgExporter struct {
	conn *pgxpool.Conn
}

func (exporter pgExporter) CopyTo(ctx context.Context, w io.Writer, sql string) error {
	_, err := exporter.conn.Conn().PgConn().CopyTo(ctx, w, sql)
	return err
}
func (exporter pgExporter) Release() {
	exporter.conn.Release()
}
