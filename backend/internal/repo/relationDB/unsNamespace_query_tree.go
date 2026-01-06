package relationDB

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"backend/share/base"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsTreeNextLevelQuery struct {
	types.UnsTreeCondition
	LayRecPrev string
}
type LayRecCc struct {
	LayRec        string `gorm:"column:lay_rec" json:"layRec"`
	CountChildren string `gorm:"column:count_children" json:"countChildren"`
}

func (p UnsNamespaceRepo) NextLevelPagedQueryList(db *gorm.DB, q *UnsTreeNextLevelQuery, searchCount *int64) (results []*TreeNodeUns, err error) {
	sql := &base.StringBuilder{}
	sql.Grow(1024)
	sql.Append(` select a.id,a.alias,a.parent_id,a.parent_alias,a.path_type,a.path,a.name,a.data_type,a.mount_type,a.mount_source
       ,bc.ptc as count_children FROM`)
	nextLevelPagedQuery(q, sql)
	db = p.model(db)
	if searchCount != nil {
		curSql := sql.String()
		countSql := "select count(*) " + curSql[strings.Index(curSql, "FROM"):]
		err = db.Raw(countSql).Scan(searchCount).Error
		if err != nil || *searchCount == 0 {
			return
		}
	}
	sql.Append(` order by a.name asc limit `).Int(q.PageSize).Append(` offset `).Int((q.PageNo - 1) * q.PageSize)
	err = db.Raw(sql.String()).Find(&results).Error
	return
}
func (p UnsNamespaceRepo) ParentIdPagedQueryList(db *gorm.DB, parentId *int64, pageNo, pageSize int, searchCount *int64) (results []*UnsNamespace, err error) {
	db = p.model(db)
	if parentId != nil {
		if *parentId == 0 {
			db = db.Where(`parent_id is null`)
		} else {
			db = db.Where(`parent_id =?`, *parentId)
		}
	}
	db = db.Where(`status=1 AND path_type in(0,2) AND id>1000`)
	if searchCount != nil {
		err = db.Count(searchCount).Error
		if err != nil || *searchCount == 0 {
			return
		}
	}
	db = db.Omit("fields", "ref_uns", "refers", "extend", "label_ids")
	db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: false}).
		Limit(pageSize).Offset((pageNo - 1) * pageSize)
	err = db.Find(&results).Error
	return
}
func (p UnsNamespaceRepo) ListCountChildren(db *gorm.DB, layRecPrev string) (results []*LayRecCc, err error) {
	sql := &base.StringBuilder{}
	sql.Grow(1024)
	sql.Append(` select a.lay_rec, b.count_children from
        ( select parent_id, STRING_AGG( path_type::text || ':' || cc::text, ',') as count_children FROM (
        select parent_id ,path_type,count(*) as cc from uns_namespace where `)
	if len(layRecPrev) > 0 {
		sql.Append(` lay_rec like '`).Append(escapeSQL(layRecPrev)).Append(`/%' and status = 1 `)
	} else {
		sql.Append(" status = 1 ")
	}
	sql.Append(` and path_type in(0,2)
        group by parent_id,path_type
        ) t group by parent_id
        )
        b inner join uns_namespace a on b.parent_id = a.id`)
	err = db.Raw(sql.String()).Find(&results).Error
	return
}

func nextLevelPagedQuery(q *UnsTreeNextLevelQuery, sql *base.StringBuilder) {
	sql.Append(` ( SELECT nid as id, STRING_AGG( path_type::text || ':' || count_children::text, ',') as ptc FROM
        (
        WITH ns_data AS (
        SELECT `)
	if layRecPrev := q.LayRecPrev; len(layRecPrev) > 0 {
		sql.Append(`"nextIdLong"('`).Append(escapeSQL(layRecPrev)).Append(`'::TEXT, lay_rec::TEXT)`)
	} else {
		sql.Append(`"nextIdLong"(''::TEXT, lay_rec::TEXT)`)
	}
	sql.Append(` AS nid,
        path_type, id
        FROM uns_namespace`)
	unsTreeNextLevelFilter(q, sql)
	sql.Append(`)
        SELECT nid, path_type, COUNT(CASE WHEN nid != id THEN 1 END) AS count_children
        FROM ns_data
        GROUP BY nid, path_type
        ) t group by nid ) bc
        inner join uns_namespace a on bc.id=a.id`)
}
func unsTreeNextLevelFilter(q *UnsTreeNextLevelQuery, sql *base.StringBuilder) {
	sql.Append(` WHERE `)
	if layRecPrev := q.LayRecPrev; len(layRecPrev) > 0 {
		sql.Append(`lay_rec like '`).Append(escapeSQL(layRecPrev)).Append(`/%' and status = 1 `)
	} else {
		sql.Append(" status = 1 ")
	}
	switch q.SearchType {
	case 1:
		if dataType := q.DataType; dataType != nil {
			sql.Append(`AND data_type =`).Int(*dataType).Append(`and path_type=2`)
		} else {
			sql.Append(`AND path_type in(0,2)`)
		}
		if subscribeEnable := q.SubscribeEnable; subscribeEnable != nil {
			sql.Append(`AND  with_flags&`).Int(constants.UnsFlagWithSubscribeEnable)
			if *subscribeEnable {
				sql.Append(`>0`)
			} else {
				sql.Append(`=0`)
			}
		}
	case 2:
		sql.Append(`AND label_ids IS NOT null and label_ids !='{}'::jsonb `)
	case 3:
		sql.Append(`AND model_id is not null `)
	}
	if keyword := q.Keyword; len(keyword) > 0 {
		keyword = "'%" + escapeLikePattern(escapeSQL(keyword)) + "%'"
		sql.Append(` AND (path ILIKE `).Append(keyword).Append(` OR alias LIKE `).Append(keyword).Append(`)`)
	}
	sql.Append(` AND id>1000`) //排除系统数据
}
