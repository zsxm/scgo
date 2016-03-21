package scsql

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/snxamdf/scgo/data"
	"github.com/snxamdf/scgo/tools"
	"github.com/snxamdf/scgo/tools/uuid"
)

const (
	SC_I       = 0
	SC_D       = 1
	SC_U       = 2
	SC_S       = 3
	SC_S_ONE   = 4
	SC_S_PAGE  = 5
	SC_S_COUNT = 6

	SELECTD        = "select"
	FROM           = "from"
	INSERT         = "insert"
	INTO           = "into"
	VALUES         = "values"
	UPDATE         = "update"
	SET            = "set"
	DELETED        = "delete"
	WHERE          = "where"
	ORDER          = "order"
	BY             = "by"
	DESC           = "desc"
	ASC            = "asc"
	SPACE          = " "
	PARENTHESESL_L = "("
	PARENTHESESL_R = ")"
	ASTERISK       = "*"
	LIMIT          = "limit"
)

var (
	NOT_SIFT_CORT = "字段 %s 未设置连接符号,使用方式:field.FieldExp().Eq().And() 或 FieldExpVal('值').Eq().And(), 有 And() 或 Or() 等!"
)

type SCSQL struct {
	S_TYPE       int                   //当前想要执行的操作
	sql          string                //生成的sql
	Args         []interface{}         //参数值
	Operate      int                   //准备执行的操作
	Entity       data.EntityInterface  //实体
	Table        data.TableInformation //表信息
	DataBaseType string                //数据库类型 mysql、oracle ...
	Page         *data.Page
}

//Parse SQL
func (this *SCSQL) ParseSQL() error {
	if this.S_TYPE == SC_I { //insert
		this.genInsert()
	} else if this.S_TYPE == SC_D { //delete
		return this.genDelete()
	} else if this.S_TYPE == SC_U { //update
		return this.genUpdate()
	} else if this.S_TYPE == SC_S { //select
		return this.genSelect()
	} else if this.S_TYPE == SC_S_ONE { //select one
		return this.genSelectOne()
	} else if this.S_TYPE == SC_S_PAGE { //select page
		return this.genSelectPage()
	} else if this.S_TYPE == SC_S_COUNT { //select count
		return this.genSelectCount()
	}
	return nil
}

//Primary Key Is Blank
func (this *SCSQL) PrimaryKeyIsBlank() bool {
	table := this.Table
	entity := this.Entity
	columns := table.Columns()
	for _, v := range columns {
		field := entity.Field(v)
		if field.PrimaryKey() && tools.IsNotBlank(field.Value()) {
			return true
		}
	}
	return false
}

//delete
func (this *SCSQL) genDelete() error {
	var wr bytes.Buffer
	entity := this.Entity
	table := this.Table
	columns := table.Columns()
	wr.WriteString(DELETED + SPACE + FROM + SPACE + table.TableName())
	sft := &sift{stype: this.S_TYPE}
	for _, v := range columns {
		field := entity.Field(v)
		_, ctor := sft.genExp(v, field)
		if ctor {
			return errors.New(fmt.Sprintf(NOT_SIFT_CORT, v))
		}
	}
	sftSql, vals := sft.genSiftSql()
	this.sql = wr.String() + SPACE + sftSql
	this.Args = vals
	return nil
}

//update
func (this *SCSQL) genUpdate() error {
	var wr bytes.Buffer
	entity := this.Entity
	table := this.Table
	columns := table.Columns()
	args := make([]interface{}, 0, len(columns))
	wr.WriteString(UPDATE + SPACE + table.TableName() + SPACE + SET + SPACE)
	sft := &sift{stype: this.S_TYPE}
	var i int
	for _, v := range columns {
		field := entity.Field(v)
		sift, ctor := sft.genExp(v, field)
		if ctor {
			return errors.New(fmt.Sprintf(NOT_SIFT_CORT, v))
		} else if (!field.PrimaryKey() && !field.FieldExp().IsSet() && tools.IsNotBlank(field.Value())) || sift {
			if i > 0 {
				wr.WriteString(", ")
			}
			wr.WriteString(v + " = ?")
			args = append(args, field.Value())
			i++
		}
	}
	sftSql, vals := sft.genSiftSql()
	args = append(args, vals...)
	this.sql = wr.String() + SPACE + sftSql
	this.Args = args
	return nil
}

//insert
func (this *SCSQL) genInsert() {
	var wr bytes.Buffer
	var wrcom bytes.Buffer
	var wrval bytes.Buffer
	entity := this.Entity
	table := this.Table
	columns := table.Columns()
	args := make([]interface{}, 0, len(columns))
	wr.WriteString(INSERT + SPACE + INTO + SPACE + table.TableName() + SPACE)

	for i, v := range columns {
		field := entity.Field(v)
		if i > 0 {
			wrcom.WriteString(",")
			wrval.WriteString(",")
		}
		if field.PrimaryKey() {
			uuid := uuid.NewV1()
			field.SetValue(uuid.String())
		}
		wrcom.WriteString(v)
		wrval.WriteString("?")
		args = append(args, field.Value())
	}
	wr.WriteString(PARENTHESESL_L + wrcom.String() + PARENTHESESL_R + SPACE)
	wr.WriteString(VALUES + PARENTHESESL_L + wrval.String() + PARENTHESESL_R + SPACE)
	this.sql = wr.String()
	this.Args = args
}

//select one
func (this *SCSQL) genSelectCount() error {
	var wr bytes.Buffer
	entity := this.Entity
	table := this.Table
	columns := table.Columns()
	sft := &sift{stype: this.S_TYPE}
	wr.WriteString(SELECTD)
	wr.WriteString(SPACE)
	wr.WriteString("count(")
	for _, v := range columns {
		field := entity.Field(v)
		_, ctor := sft.genExp(v, field)
		if ctor {
			return errors.New(fmt.Sprintf(NOT_SIFT_CORT, v))
		}
		if field.PrimaryKey() {
			wr.WriteString(v)
		}
	}
	wr.WriteString(")")
	wr.WriteString(SPACE)
	wr.WriteString(FROM)
	wr.WriteString(SPACE)
	wr.WriteString(table.TableName())
	wr.WriteString(SPACE)
	sftSql, vals := sft.genSiftSql()
	this.sql = wr.String() + sftSql
	this.Args = vals
	return nil
}

//select one
func (this *SCSQL) genSelectOne() error {
	var wr bytes.Buffer
	entity := this.Entity
	table := this.Table
	columns := table.Columns()
	sft := &sift{stype: this.S_TYPE}
	wr.WriteString(SELECTD)
	wr.WriteString(SPACE)
	for i, v := range columns {
		if i > 0 {
			wr.WriteString(",")
		}
		field := entity.Field(v)
		_, ctor := sft.genExp(v, field)
		if ctor {
			return errors.New(fmt.Sprintf(NOT_SIFT_CORT, v))
		}
		wr.WriteString(v)
	}
	wr.WriteString(SPACE)
	wr.WriteString(FROM)
	wr.WriteString(SPACE)
	wr.WriteString(table.TableName())
	wr.WriteString(SPACE)
	sftSql, vals := sft.genSiftSql()
	this.sql = wr.String() + sftSql
	this.Args = vals
	return nil
}

//select
func (this *SCSQL) genSelectPage() error {
	var wr bytes.Buffer
	entity := this.Entity
	table := this.Table
	columns := table.Columns()
	page := this.Page
	sft := &sift{stype: this.S_TYPE}
	wr.WriteString(SELECTD)
	wr.WriteString(SPACE)
	for i, v := range columns {
		if i > 0 {
			wr.WriteString(",")
		}
		field := entity.Field(v)
		_, ctor := sft.genExp(v, field)
		if ctor {
			return errors.New(fmt.Sprintf(NOT_SIFT_CORT, v))
		}
		wr.WriteString(v)
	}
	wr.WriteString(SPACE)
	wr.WriteString(FROM)
	wr.WriteString(SPACE)
	wr.WriteString(table.TableName())
	wr.WriteString(SPACE)
	sftSql, vals := sft.genSiftSql()
	var pl bytes.Buffer
	if this.DataBaseType == data.DATA_BASE_MYSQL { //mysql page
		pl.WriteString(LIMIT + SPACE + fmt.Sprint(page.FirstResult, ",", page.MaxResults))
	} else if this.DataBaseType == data.DATA_BASE_ORACLE { //oracle page

	}
	this.sql = wr.String() + sftSql + pl.String()
	this.Args = vals
	return nil
}

//select
func (this *SCSQL) genSelect() error {
	var wr bytes.Buffer
	entity := this.Entity

	table := this.Table
	columns := table.Columns()
	sft := &sift{stype: this.S_TYPE}
	wr.WriteString(SELECTD)
	wr.WriteString(SPACE)
	for i, v := range columns {
		if i > 0 {
			wr.WriteString(",")
		}
		field := entity.Field(v)
		_, ctor := sft.genExp(v, field)
		if ctor {
			return errors.New(fmt.Sprintf(NOT_SIFT_CORT, v))
		}
		wr.WriteString(v)
	}
	wr.WriteString(SPACE)
	wr.WriteString(FROM)
	wr.WriteString(SPACE)
	wr.WriteString(table.TableName())
	wr.WriteString(SPACE)
	sftSql, vals := sft.genSiftSql()
	this.sql = wr.String() + sftSql
	this.Args = vals
	return nil
}

type sift struct {
	sifts [][]string
	stype int
	sort  [][]string
}

//exp return exp ctor
//return bool:sift,bool:ctor, sift不使用原有的字段，ctor未设置连接符
func (this *sift) genExp(column string, field data.EntityField) (bool, bool) {
	fieldExp := field.FieldExp()
	fieldSort := field.FieldSort()
	if fieldSort.IsSet() {
		sort := make([]string, 3)
		sort[0] = column
		sort[1] = fieldSort.Value()
		sort[2] = strconv.Itoa(fieldSort.Index())
		this.sort = append(this.sort, sort)
	}
	if fieldExp.IsSet() {
		value := fieldExp.Value()       //比较值
		exp := fieldExp.Exp()           //比较符
		ctor := fieldExp.Ctor().Value() //连接符
		ctorLen := len(ctor)
		for i, v := range value {
			if i <= ctorLen {
				sft := make([]string, 4)
				sft[0] = column
				sft[1] = exp[i]
				sft[2] = v
				sft[3] = ctor[i]
				this.sifts = append(this.sifts, sft)
			}
		}
		//if ctorLen == 1 {
		//调用FieldExp()为0,FieldExpVal("值")为1
		if len(value) == 0 && ctorLen == 1 {
			sft := make([]string, 4)
			sft[0] = column
			sft[1] = exp[0]
			sft[2] = field.Value()
			sft[3] = ctor[0]
			this.sifts = append(this.sifts, sft)
			return false, false
		} else if ctorLen == 0 {
			this.sifts = [][]string{}
			return false, true
		}
	}
	return true, false
}

//排序
func (this *sift) genSort() {
	sort.Sort(this)
}

//gen exp sift sql
func (this *sift) genSiftSql() (string, []interface{}) {
	this.genSort()
	var wr bytes.Buffer
	args := make([]interface{}, 0, 1)
	var ctor string
	var val, whe string
	for i, v := range this.sifts {
		switch v[1] {
		case data.EXP_LK, data.EXP_LK_R, data.EXP_LK_L:
			whe = "like"
			like := v[1]
			if strings.Index(like, "%") == -1 {
				val = "%" + v[2] + "%"
			} else if strings.Index(like, "%") == 0 {
				val = "%" + v[2]
			} else if strings.Index(like, "%") == 4 {
				val = v[2] + "%"
			}
			args = append(args, val)
			break
		default:
			whe = v[1]
			args = append(args, v[2])
			break
		}
		if i > 0 {
			wr.WriteString(SPACE)
		}
		wr.WriteString(ctor + SPACE + v[0] + SPACE + whe + SPACE + "?")
		ctor = v[3]
	}
	var order bytes.Buffer
	for i, v := range this.sort {
		if i > 0 {
			order.WriteString(", ")
		}
		order.WriteString(v[0] + SPACE + v[1])
	}
	var res string
	if tools.IsNotBlank(wr.String()) {
		res = WHERE + wr.String() + SPACE
	}
	if tools.IsNotBlank(order.String()) {
		res = res + ORDER + SPACE + BY + SPACE + order.String() + SPACE
	}
	return res, args
}

//sql
func (this *SCSQL) SQL() string {
	log.Println("SQL :", this.sql, "ARGS :", this.Args)
	return this.sql
}

func (this sift) Len() int {
	return len(this.sort)
}

func (this sift) Less(i, j int) bool {
	is := this.sort[i]
	js := this.sort[j]
	ii := tools.ParseInteger(is[2])
	ji := tools.ParseInteger(js[2])
	if ii < ji {
		return true
	} else if ii > ji {
		return false
	} else {
		return is[0] < js[0]
	}
}

func (this sift) Swap(i, j int) {
	temp := this.sort[i]
	this.sort[i] = this.sort[j]
	this.sort[j] = temp
}
