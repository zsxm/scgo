package mysql_test

import (
	//	"fmt"
	"log"
	"study/app3/source/chatol/entity"
	"testing"

	//	"sync"

	"github.com/zsxm/scgo/data/scdb"
)

//func TestMysqlSelect(t *testing.T) {
//	var repository = scdb.Connection
//	var e = entity.NewMessageBean()
//	msg := entity.NewMessage()
//	msg.SetId("A070DAE1-E4D5-11E5-B8D4-3010B3A0F15C")
//	msg.Id().FieldExp().Eq().And()
//	e.SetEntity(msg)

//	repository.Select(e)
//	es := e.Entitys()
//	for i, v := range es.Values() {
//		log.Println(i, v.Id().Value(), v.Name().Value())
//	}
//}

//func TestMysqlUpdate(t *testing.T) {
//	var repository = scdb.Connection
//	m := entity.NewMessage()
//	m.SetId("A070DAE1-E4D5-11E5-B8D4-3010B3A0F15C")
//	m.Id().FieldExp().Eq().And()
//	m.SetName("张三1")
//	//m.SetPhone("15164383441")
//	m.SetAge(25)
//	result, err := repository.Update(m)
//	row, err := result.RowsAffected()
//	log.Println("Update", row, err)
//}

//func TestMysqlSave(t *testing.T) {
//	var repository = scdb.Connection
//	for i := 0; i < 20; i++ {
//		m := entity.NewMessage()
//		m.SetName(fmt.Sprint("张三1", i))
//		m.SetPhone(fmt.Sprint("15164383441", i))
//		m.SetAge(25 + i)
//		result, err := repository.Save(m)
//		if err == nil {
//			row, err := result.RowsAffected()
//			log.Println("Save", row, err)
//		}

//	}
//}

//func TestMysqlSaveOrUpdate(t *testing.T) {
//	var repository = scdb.Connection
//	m := entity.NewMessage()
//	m.SetId("A3FC79C5E5C611E59BB63010B3A0F15C")
//	m.Id().FieldExp().Eq().And()
//	m.SetName("张三b")
//	//m.Name().FieldExp().Eq().And()
//	m.SetPhone("15164383441")
//	m.SetAge(25)
//	result, err := repository.SaveOrUpdate(m)
//	row, err := result.RowsAffected()
//	log.Println("SaveOrUpdate", row, err)
//}

//func TestMysqlSelectOne(t *testing.T) {
//	var repository = scdb.Connection
//	m := entity.NewMessage()
//	m.SetId("6CDDE56AE5CD11E5A0233010B3A0F15C")
//	m.Id().FieldExp().Eq().And()
//	//	m.SetName("张三")
//	//	m.Name().FieldExp().Lk().And()

//	err := repository.SelectOne(m)

//	log.Println("SelectOne", err, m.JSON())
//}

//func TestMysqlDelete(t *testing.T) {
//	var repository = scdb.Connection
//	m := entity.NewMessage()
//	m.SetId("4D5CCADFE67011E5A8DE3010B3A0F15C")
//	m.Id().FieldExp().Eq().And()

//	r, err := repository.Delete(m)
//	rows, _ := r.RowsAffected()
//	last, _ := r.LastInsertId()
//	log.Println(rows, last, err)
//}

//func TestMysqlSelectPage(t *testing.T) {
//	var repository = scdb.Connection
//	bean := entity.NewMessageBean()
//	m := entity.NewMessage()
//	m.SetAge(31)
//	m.Age().FieldExp().Eq().And()
//	m.Age().FieldSort().Asc(1)
//	//m.Id().FieldSort().Desc(2)

//	bean.SetEntity(m)
//	page := &data.Page{}
//	page.New(1, 5)
//	repository.SelectPage(bean, page)
//	for _, v := range bean.Entitys().Values() {
//		log.Println(v.Id().Value(), v.Name().Value(), v.Age().Value())
//	}
//}

//func TestMysqlSelectCount(t *testing.T) {
//	config.Conf = &config.Config{
//		FilePath: "../config/db.xml",
//	}
//	mysql.New(config.Conf)
//	var repository = scdb.Connection
//	e := entity.NewMessage()
//	e.SetAge(30)
//	e.Age().FieldExp().Gt().And()
//	c, err := repository.SelectCount(e)
//	log.Println("count=", c, err)
//}
func TestMysqlExecute(t *testing.T) {
	var repository = scdb.Connection
	var sql = "select count(*) from users where u_age>39"
	c, err := repository.Query(sql)
	log.Println("Query=", c.Data, err)
	bean := entity.NewMessageBean()
	repository.Select(bean)
	log.Println("Select=", bean.Entitys().JSON(), err)
	//log.Println("count=", bean.Entitys().JSON())
	//	wg := sync.WaitGroup{}
	//	wg.Add(10)
	//	for i := 0; i < 10; i++ {
	//		go func(i int) {
	//			for j := 0; j < 100; j++ {
	//				m := entity.NewMessage()
	//				m.SetName(fmt.Sprint("张三1", i, j))
	//				m.SetPhone(fmt.Sprint("15164383441", i, j))
	//				m.SetAge(25 + i + j)
	//				result, err := repository.Save(m)
	//				if err == nil {
	//					row, err := result.RowsAffected()
	//					log.Println("Save", row, err)
	//				}
	//			}
	//			wg.Done()
	//		}(i)
	//	}
	//	wg.Wait()

}
