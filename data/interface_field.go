package data

const (
	EXP_EQ       = "="
	EXP_LK       = "like"
	EXP_LK_R     = "%like"
	EXP_LK_L     = "like%"
	EXP_GT       = ">"
	EXP_LT       = "<"
	EXP_GT_EQ    = ">="
	EXP_LT_EQ    = "<="
	EXP_AND      = "and"
	EXP_OR       = "or"
	EXP_IN       = "in"
	EXP_Buckle_R = "("
	EXP_Buckle_L = ")"
	SORT_ASC     = "asc"
	SORT_DESC    = "desc"
)

//------------------interface begin-------------------------------
//实体字段接口
type EntityField interface {
	Type() string
	SetValue(value string)
	Pointer() *string
	Value() string
	PrimaryKey() bool
	FieldExp() *FieldExpression
	FieldExpVal(value string) *FieldExpression
	FieldSort() *FieldSort
}

//------------------ begin-------------------------------
//字段表达式
type FieldExpression struct {
	value []string
	exp   []string
	b     bool
	ctor  *FieldConnector
}

//字段连接符
type FieldConnector struct {
	b     bool
	value []string
}

//字段排序
type FieldSort struct {
	value string
	b     bool
	index int
}

func (this *FieldSort) Asc(index int) {
	this.b = true
	this.value = SORT_ASC
	this.index = index
}

func (this *FieldSort) Desc(index int) {
	this.b = true
	this.value = SORT_DESC
	this.index = index
}

func (this *FieldSort) Value() string {
	return this.value
}

func (this *FieldSort) Index() int {
	return this.index
}

func (this *FieldSort) IsSet() bool {
	return this.b
}

//and
func (this *FieldConnector) And() *FieldConnector {
	this.b = true
	this.value = append(this.value, EXP_AND)
	return this
}

//or
func (this *FieldConnector) Or() *FieldConnector {
	this.b = true
	this.value = append(this.value, EXP_OR)
	return this
}

//in
func (this *FieldConnector) In() *FieldConnector {
	this.b = true
	this.value = append(this.value, EXP_IN)
	return this
}

//Buckle R
func (this *FieldConnector) BuckleR() *FieldConnector {
	this.b = true
	this.value = append(this.value, EXP_Buckle_R)
	return this
}

//Buckle L
func (this *FieldConnector) BuckleL() *FieldConnector {
	this.b = true
	this.value = append(this.value, EXP_Buckle_L)
	return this
}

//result and or or
func (this *FieldConnector) Value() []string {
	return this.value
}

func (this *FieldConnector) IsSet() bool {
	return this.b
}

//=
func (this *FieldExpression) Eq() *FieldConnector {
	this.b = true
	this.exp = append(this.exp, EXP_EQ)
	return this.ctor
}

//like
func (this *FieldExpression) Lk() *FieldConnector {
	this.b = true
	this.exp = append(this.exp, EXP_LK)
	return this.ctor
}

//like%
func (this *FieldExpression) LkR() *FieldConnector {
	this.b = true
	this.exp = append(this.exp, EXP_LK_R)
	return this.ctor
}

//%like
func (this *FieldExpression) LkL() *FieldConnector {
	this.b = true
	this.exp = append(this.exp, EXP_LK_L)
	return this.ctor
}

//>
func (this *FieldExpression) Gt() *FieldConnector {
	this.b = true
	this.exp = append(this.exp, EXP_GT)
	return this.ctor
}

//<
func (this *FieldExpression) Lt() *FieldConnector {
	this.b = true
	this.exp = append(this.exp, EXP_LT)
	return this.ctor
}

//>=
func (this *FieldExpression) GtEq() *FieldConnector {
	this.b = true
	this.exp = append(this.exp, EXP_GT_EQ)
	return this.ctor
}

//<=
func (this *FieldExpression) LtEq() *FieldConnector {
	this.b = true
	this.exp = append(this.exp, EXP_LT_EQ)
	return this.ctor
}

func (this *FieldExpression) Ctor() *FieldConnector {
	return this.ctor
}

func (this *FieldExpression) IsSet() bool {
	return this.b
}

//result exp []value
func (this *FieldExpression) Value() []string {
	return this.value
}

//result exp []exp
func (this *FieldExpression) Exp() []string {
	return this.exp
}

//------------------ end-------------------------------
