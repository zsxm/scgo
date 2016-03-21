package scdb

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/snxamdf/scgo/data"
	"github.com/snxamdf/scgo/data/dbconfig"
	"github.com/snxamdf/scgo/data/scsql"
)

var Connection RepositoryInterface

type Repository struct {
	dBSource DBSourceInterface
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
		return nil, err
	}

	stmt, err := this.Prepare(csql)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(csql.Args...)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}
	return result, nil
}

func (this *Repository) Update(entity data.EntityInterface) (sql.Result, error) {
	table := entity.Table()
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), Table: table, S_TYPE: scsql.SC_U, Entity: entity}
	err := csql.ParseSQL()
	if err != nil {
		return nil, err
	}

	stmt, err := this.Prepare(csql)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(csql.Args...)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}
	return result, nil
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

func (this *Repository) Delete(entity data.EntityInterface) (sql.Result, error) {
	table := entity.Table()
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), Table: table, S_TYPE: scsql.SC_D, Entity: entity}
	err := csql.ParseSQL()
	if err != nil {
		return nil, err
	}
	stmt, err := this.Prepare(csql)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(csql.Args...)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}
	return result, nil
}

func (this *Repository) SelectOne(entity data.EntityInterface) error {
	table := entity.Table()
	csql := scsql.SCSQL{DataBaseType: this.dBSource.DataBaseType(), S_TYPE: scsql.SC_S_ONE, Table: table, Entity: entity}
	err := csql.ParseSQL()
	if err != nil {
		log.Println(err)
		return err
	}

	stmt, err := this.Prepare(csql)
	if err != nil {
		log.Println("error stmt", err)
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(csql.Args...)
	if err != nil {
		log.Println("error rows", err)
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println("error cols", err)
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
			log.Println("error", err)
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
		log.Println(err)
		return 0, err
	}
	stmt, err := this.Prepare(csql)
	if err != nil {
		log.Println("error stmt", err)
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(csql.Args...)
	if err != nil {
		log.Println("error rows", err)
		return 0, err
	}
	defer rows.Close()
	var resCount int
	for rows.Next() {
		err = rows.Scan(&resCount)
		if err != nil {
			log.Println("error", err)
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
		log.Println(err)
		return err
	}
	count, err := this.SelectCount(entity)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return err
	}

	return this.selected(csql, entityBean)
}

func (this *Repository) selected(csql scsql.SCSQL, entityBean data.EntityBeanInterface) error {
	stmt, err := this.Prepare(csql)
	if err != nil {
		log.Println("error stmt", err)
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(csql.Args...)
	if err != nil {
		log.Println("error rows", err)
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println("error cols", err)
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
			log.Println("error scan", err)
			return err
		}
		beans.Add(bean)
	}
	entityBean.SetEntitys(beans)
	return nil
}

//执行自定义DML语言. (DDL,DCL待添加)
func (this *Repository) Execute(sql string, args ...interface{}) {

}

func (this *Repository) Prepare(csql scsql.SCSQL) (*sql.Stmt, error) {
	var db = this.DB()
	stmt, err := db.Prepare(csql.SQL())
	return stmt, err
}

func init() {
	dbconfig.Conf = &dbconfig.Config{
		FilePath: `conf/db.xml`,
	}
	dbconfig.Conf.Init()
	Connection = NewRepository(dbconfig.Conf.Dbs.Default)
}
