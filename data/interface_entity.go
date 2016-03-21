package data

//实体bean
type EntityBeanInterface interface {
	NewEntity() EntityInterface
	NewEntitys(cap int) EntitysInterface
	SetEntity(bean EntityInterface)
	SetEntitys(beans EntitysInterface)
	Table() TableInformation
	FieldNames() FieldNames
	Entity() EntityInterface
	Entitys() EntitysInterface
}

//单个实体
type EntityInterface interface {
	SetValue(filed, value string)
	Field(filedName string) EntityField
	JSON() string
	Table() TableInformation
	FieldNames() FieldNames
}

//多个个实体
type EntitysInterface interface {
	SetPage(page *Page)
	Add(e EntityInterface)
	Table() TableInformation
	FieldNames() FieldNames
	JSON() string
	Len() int
	Values() []EntityInterface
}

//------------------FieldNames begin-------------------------------
type FieldNames struct {
	names []string
}

func (this *FieldNames) SetNames(names []string) {
	this.names = names
}

func (this *FieldNames) Names() []string {
	return this.names
}

//------------------FieldNames begin-------------------------------

//------------------TableInformation begin-------------------------------
type TableInformation struct {
	tableName string
	columns   []string
}

func (this *TableInformation) SetTableName(tableName string) {
	this.tableName = tableName
}

func (this *TableInformation) SetColumns(columns []string) {
	this.columns = columns
}

func (this *TableInformation) TableName() string {
	return this.tableName
}

func (this *TableInformation) Columns() []string {
	return this.columns
}

//------------------TableInformation begin-------------------------------
