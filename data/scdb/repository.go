package scdb

import (
	"database/sql"

	"github.com/snxamdf/scgo/data"
	"github.com/snxamdf/scgo/data/scsql"
)

type RepositoryInterface interface {
	//返回db
	DB() *sql.DB
	//保存对象,参数 : entity data.EntityInterface
	Save(entity data.EntityInterface) (sql.Result, error)

	//修改对象,参数 : entity data.EntityInterface
	Update(entity data.EntityInterface) (sql.Result, error)

	//保存或修改对象,参数 : entity data.EntityInterface
	SaveOrUpdate(entity data.EntityInterface) (sql.Result, error)

	//查询多条,参数 : entity data.EntityBeanInterface
	Select(entityBean data.EntityBeanInterface) error

	//分页查询,参数 : entity data.EntityBeanInterface
	SelectPage(entityBean data.EntityBeanInterface, page *data.Page) error

	//分页数量,参数 : entity data.EntityInterface
	SelectCount(entity data.EntityInterface) (int, error)

	//查询一条,参数 : entity data.EntityInterface
	SelectOne(entity data.EntityInterface) error

	//删除,参数 : entity data.EntityInterface
	Delete(entity data.EntityInterface) (sql.Result, error)

	//执行自定义DML语言. (DDL,DCL待添加)
	Execute(sql string, args ...interface{})

	//语句解析
	Prepare(csql scsql.SCSQL) (*sql.Stmt, error)
}
