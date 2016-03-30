package scdb

import (
	"database/sql"
	"errors"
	"strconv"

	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zsxm/scgo/data"
	"github.com/zsxm/scgo/data/dbconfig"
	"github.com/zsxm/scgo/data/scsql"
	"github.com/zsxm/scgo/log"
)

var Connection RepositoryInterface

type Repository struct {
	dBSource DBSourceInterface
	sync.Mutex
}

func NewRepository(dbName string) *Repository {
	var db = dbconfig.Conf.Db(dbName)
	var dbType = db.DataBaseType
	c := &Config{
		DriverName:   db.DriverName,
		UserName:     db.UserName,
		PassWord:     db.PassWord,
		Ip:           db.IP,
		Prot:         db.Prot,
		DBName:       db.Database,
		Charset:      db.Charset,
		MaxIdleConns: db.MaxIdleConns,
		MaxOpenConns: db.MaxOpenConns,
	}
	if dbType == data.DATA_BASE_MYSQL {
		c.MySqlInit()
	} else if dbType == data.DATA_BASE_ORACLE {
		c.OracleInit()
	}
	e := &Repository{}
	e.SetDBSource(c)
	return e
}

func (this *Repository) SetDBSource(dBSource DBSourceInterface) {
	this.dBSource = dBSource
}

func (this *Repository) DB() *sql.DB {
	db := this.dBSource.DB()
	return db
}

func (this *Repository) Save(entity data.EntityInterface) (sql.Result, error) {
	table := entity.Table()
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), Table: table, Entity: entity}
	if csql.PrimaryKeyIsBlank() {
		csql.S_TYPE = scsql.SC_U
	} else {
		csql.S_TYPE = scsql.SC_I
	}
	err := csql.ParseSQL()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return this.Execute(csql.SQL(), csql.Args...)
}

func (this *Repository) Update(entity data.EntityInterface) (sql.Result, error) {
	table := entity.Table()
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), Table: table, S_TYPE: scsql.SC_U, Entity: entity}
	err := csql.ParseSQL()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return this.Execute(csql.SQL(), csql.Args...)
}

func (this *Repository) SaveOrUpdate(entity data.EntityInterface) (sql.Result, error) {
	table := entity.Table()
	csql := scsql.SCSQL{Table: table, Entity: entity}
	if csql.PrimaryKeyIsBlank() {
		return this.Update(entity)
	} else {
		return this.Save(entity)
	}
}

func (this *Repository) Delete(entity data.EntityInterface, deleted ...bool) (sql.Result, error) {
	if len(deleted) > 0 && deleted[0] { //逻辑删除
		if entity.Field("deleted") == nil {
			return nil, errors.New("逻辑删除需要有 deleted 字段")
		}
		entity.SetValue("deleted", "1")
		return this.Update(entity)
	}
	table := entity.Table()
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), Table: table, S_TYPE: scsql.SC_D, Entity: entity}
	err := csql.ParseSQL()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return this.Execute(csql.SQL(), csql.Args...)
}

func (this *Repository) SelectOne(entity data.EntityInterface) error {
	table := entity.Table()
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), S_TYPE: scsql.SC_S_ONE, Table: table, Entity: entity}
	err := csql.ParseSQL()
	if err != nil {
		log.Error(err)
		return err
	}

	stmt, err := this.Prepare(csql.SQL())
	if err != nil {
		log.Error(" stmt", err)
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(csql.Args...)
	if err != nil {
		log.Error(" rows", err)
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Error(" cols", err)
		return err
	}

	colsLen := len(cols)

	for rows.Next() {
		vals := make([]interface{}, colsLen)
		for i := 0; i < colsLen; i++ {
			colm := cols[i]
			if field := entity.Field(colm); field != nil {
				vals[i] = field.Pointer()
			}
		}
		err = rows.Scan(vals...)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	}
	return nil
}

func (this *Repository) SelectCount(entity data.EntityInterface) (int, error) {
	table := entity.Table()
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), S_TYPE: scsql.SC_S_COUNT, Table: table, Entity: entity}
	err := csql.ParseSQL()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	stmt, err := this.Prepare(csql.SQL())
	if err != nil {
		log.Error(" stmt", err)
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(csql.Args...)
	if err != nil {
		log.Error(" rows", err)
		return 0, err
	}
	defer rows.Close()
	var resCount int
	for rows.Next() {
		err = rows.Scan(&resCount)
		if err != nil {
			log.Error(err)
			return 0, err
		}
		return resCount, nil
	}
	return 0, nil
}

func (this *Repository) SelectPage(entityBean data.EntityBeanInterface, page *data.Page) error {
	table := entityBean.Table()
	entity := entityBean.Entity()
	if entity == nil {
		entity = entityBean.NewEntity()
	}
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), S_TYPE: scsql.SC_S_PAGE, Table: table, Entity: entity, Page: page}
	err := csql.ParseSQL()
	if err != nil {
		log.Error(err)
		return err
	}
	count, err := this.SelectCount(entity)
	if err != nil {
		log.Error(err)
		return err
	}
	page.TotalRow = count
	err = this.selected(csql, entityBean)
	if entityBean.Entitys() != nil {
		entityBean.Entitys().SetPage(page)
	}
	return err
}

func (this *Repository) Select(entityBean data.EntityBeanInterface) error {
	table := entityBean.Table()
	entity := entityBean.Entity()
	if entity == nil {
		entity = entityBean.NewEntity()
	}
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), S_TYPE: scsql.SC_S, Table: table, Entity: entity}
	err := csql.ParseSQL()
	if err != nil {
		log.Error(err)
		return err
	}

	return this.selected(csql, entityBean)
}

func (this *Repository) selected(csql scsql.SCSQL, entityBean data.EntityBeanInterface) error {
	stmt, err := this.Prepare(csql.SQL())
	if err != nil {
		log.Error(" stmt", err)
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(csql.Args...)
	if err != nil {
		log.Error(" rows", err)
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Error(" cols", err)
		return err
	}

	colsLen := len(cols)
	beans := entityBean.NewEntitys(5)

	for rows.Next() {
		vals := make([]interface{}, colsLen)
		bean := entityBean.NewEntity()
		for i := 0; i < colsLen; i++ {
			colm := cols[i]
			if field := bean.Field(colm); field != nil {
				vals[i] = field.Pointer()
			}
		}
		err = rows.Scan(vals...)
		if err != nil {
			log.Error(" scan", err)
			return err
		}
		beans.Add(bean)
	}
	entityBean.SetEntitys(beans)
	return nil
}

func (this *Repository) Execute(sql string, args ...interface{}) (sql.Result, error) {
	stmt, err := this.Prepare(sql)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

//执行自定义DML语言. (DDL,DCL待添加)
//return []slice,error
func (this *Repository) Query(sql string, args ...interface{}) (data.QueryResult, error) {
	var result data.QueryResult
	stmt, err := this.Prepare(sql)
	if err != nil {
		log.Error(" stmt", err)
		return result, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Error(" rows", err)
		return result, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		log.Error(" cols", err)
		return result, err
	}
	colsLen := len(cols)
	scanArgs := make([]interface{}, colsLen)
	values := make([]interface{}, colsLen)
	for j := range values {
		scanArgs[j] = &values[j]
	}
	result = data.QueryResult{Data: make([]data.Map, 0, 1)}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Error(" scan", err)
			return result, err
		}
		tmpMap := make(data.Map, colsLen)
		for i, val := range values {
			if val != nil {
				if v, ok := val.(int64); ok {
					tmpMap[cols[i]] = strconv.Itoa(int(v))
				} else if v, ok := val.(float64); ok {
					tmpMap[cols[i]] = strconv.FormatFloat(v, 'f', -1, 64)
				} else if v, ok := val.(uint64); ok {
					tmpMap[cols[i]] = strconv.FormatUint(v, 10)
				} else if v, ok := val.(bool); ok {
					tmpMap[cols[i]] = strconv.FormatBool(v)
				} else {
					tmpMap[cols[i]] = string(val.([]byte))
				}

			}
		}
		result.Data = append(result.Data, tmpMap)
	}
	return result, nil
}

func (this *Repository) Prepare(sql string) (*sql.Stmt, error) {
	var db = this.DB()
	stmt, err := db.Prepare(sql)
	return stmt, err
}

func init() {
	dbconfig.Conf = &dbconfig.Config{
		FilePath: `conf/db.xml`,
	}
	dbconfig.Conf.Init()
	Connection = NewRepository(dbconfig.Conf.Dbs.Default)
}
