package util

import (
	"context"
	"strings"

	"github.com/lanwenhong/lgobase/dbpool"
	"github.com/lanwenhong/lgobase/logger"
)

type PageDataDb struct {
	Dbs      *dbpool.Dbpool
	DbName   string
	Sql      string
	CountSql string
	Records  int64
	Data     []map[string]interface{}
}

type Page struct {
	Count    int64
	Pages    int64
	PageSize int64
	Pdd      *PageDataDb
	Page     int64
}

func NewPageDataDb(ctx context.Context, db *dbpool.Dbpool, sql string, csql string, dbname string) *PageDataDb {
	pdd := PageDataDb{
		Dbs:     db,
		Records: -1,
		Data:    []map[string]interface{}{},
		DbName:  dbname,
	}
	q_sql := " limit ?,?"
	pdd.Sql = sql + q_sql
	logger.Debugf(ctx, "q_sql: %s", pdd.Sql)

	if csql != "" {
		pdd.CountSql = csql
	} else {
		f_index := strings.Index(sql, " from ")
		f_sql := sql[f_index:]
		o_index := strings.Index(f_sql, " order by ")
		if o_index > 0 {
			f_sql = f_sql[:o_index]
		}
		pdd.CountSql = "select count(*) as count " + f_sql
		logger.Debugf(ctx, "count_sql: %s", pdd.CountSql)
	}
	return &pdd
}

func (pdd *PageDataDb) Load(ctx context.Context, cur int64, pagesize int64) ([]map[string]interface{}, error) {
	if len(pdd.Data) > 0 {
		return pdd.Data, nil
	}
	page := (cur - 1) * pagesize
	db := pdd.Dbs.OrmPools[pdd.DbName]
	ret := db.WithContext(ctx).Raw(pdd.Sql, page, pagesize).Scan(&pdd.Data)
	if ret.Error != nil {
		logger.Warnf(ctx, "db query: %s", ret.Error.Error())
		return pdd.Data, ret.Error
	}
	return pdd.Data, nil
}

func (pdd *PageDataDb) Count(ctx context.Context, pagesize int64) (int64, int64, error) {
	row := map[string]interface{}{}
	db := pdd.Dbs.OrmPools[pdd.DbName]
	ret := db.WithContext(ctx).Raw(pdd.CountSql).Find(&row)
	if ret.Error != nil {
		logger.Warnf(ctx, "db query: %s", ret.Error.Error())
		return -1, -1, ret.Error
	}
	//pdd.Records = row["count"].(int64)
	rds, _ := row["count"]
	pdd.Records = rds.(int64)
	var page_count int64 = 0
	if pdd.Records%pagesize == 0 {
		page_count = pdd.Records % pagesize
	} else {
		page_count = pdd.Records/pagesize + 1
	}
	logger.Debugf(ctx, "records: %d count: %d", pdd.Records, page_count)
	return pdd.Records, page_count, nil
}

func NewPage(ctx context.Context, page int64, page_size int64, pdd *PageDataDb) *Page {
	p := Page{
		Page:     page,
		PageSize: page_size,
		Pdd:      pdd,
		Count:    -1,
	}
	return &p
}
func (pg *Page) Split(ctx context.Context) error {
	var err error = nil
	if pg.Count < 0 {
		logger.Debugf(ctx, "=====")
		pg.Count, pg.Pages, err = pg.Pdd.Count(ctx, pg.PageSize)
		if err != nil {
			return err
		}
	}
	_, err = pg.Pdd.Load(ctx, pg.Page, pg.PageSize)
	return err
}
