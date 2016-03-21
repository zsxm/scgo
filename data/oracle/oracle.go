package oracle

import (
	"scgo/sc/data/scdb"
)

//alias 别名,driverName 驱动名称
func New(alias, driverName string) *scdb.Repository {
	c := &scdb.Config{
		Alias:        alias,
		DriverName:   driverName,
		UserName:     "root",
		PassWord:     "root",
		Ip:           "localhost",
		Prot:         "3306",
		DBName:       "testdb",
		Charset:      "UTF8",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}
	c.MySqlInit()
	return &scdb.Repository{
		DBSource: c,
	}
}
