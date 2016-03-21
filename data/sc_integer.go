package data

//------------------Integer begin-------------------------------
//整型类型
type Integer struct {
	value      string
	primaryKey bool
	exp        *FieldExpression
	sort       *FieldSort
}

func (this *Integer) SetValue(value string) {
	this.value = value
}

func (this *Integer) Type() string {
	return "int"
}

func (this *Integer) StructType() *Integer {
	return this
}

func (this *Integer) Pointer() *string {
	return &this.value
}

func (this *Integer) Value() string {
	return *this.Pointer()
}

func (this *Integer) SetPrimaryKey(b bool) {
	this.primaryKey = b
}

func (this *Integer) PrimaryKey() bool {
	return this.primaryKey
}

func (this *Integer) FieldExp() *FieldExpression {
	if this.exp == nil {
		e := &FieldExpression{
			ctor: &FieldConnector{},
		}
		this.exp = e
	}
	return this.exp
}

func (this *Integer) FieldExpVal(value string) *FieldExpression {
	this.FieldExp()
	this.exp.value = append(this.exp.value, value)
	return this.exp
}

func (this *Integer) FieldSort() *FieldSort {
	if this.sort == nil {
		this.sort = &FieldSort{}
	}
	return this.sort
}

//------------------Integer end-------------------------------
