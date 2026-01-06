package relationDB

import (
	"context"
	"io"
)

func (p UnsLabelRepo) ExportCsv(w io.Writer) error {
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	query := `COPY (SELECT label_name FROM uns_label order by id asc) TO STDOUT WITH CSV HEADER`
	err = conn.CopyTo(context.Background(), w, query)
	return err
}
func (p UnsLabelRepo) Csv2Model(headers, vs []string) *UnsLabel {
	po := &UnsLabel{LabelName: vs[0]}
	return po
}
