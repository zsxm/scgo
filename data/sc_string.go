package data

//------------------String begin-------------------------------
//字符串类型
type String struct {
	value      string
	primaryKey bool
	exp        *FieldExpression
	sort       *FieldSort
}

func (this *String) SetValue(value string) {
	this.value = value
}

func (this *String) Type() string {
	return "string"
}

func (this *String) StructType() *String {
	return this
}

func (this *String) Pointer() *string {
	return &this.value
}

func (this *String) Value() string {
	return *this.Pointer()
}

func (this *String) SetPrimaryKey(b bool) {
	this.primaryKey = b
}

func (this *String) PrimaryKey() bool {
	return this.primaryKey
}

func (this *String) FieldExp() *FieldExpression {
	if this.exp == nil {
		e := &FieldExpression{
			ctor: &FieldConnector{},
		}
		this.exp = e
	}
	return this.exp
}

func (this *String) FieldExpVal(value string) *FieldExpression {
	this.FieldExp()
	this.exp.value = append(this.exp.value, value)
	return this.exp
}

func (this *String) FieldSort() *FieldSort {
	if this.sort == nil {
		this.sort = &FieldSort{}
	}
	return this.sort
}

//------------------String end-------------------------------
