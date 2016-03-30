package data

type Map map[string]string

func (this Map) Get(key string) string {
	return this[key]
}

type QueryResult struct {
	Data []Map
}

func (this QueryResult) Get(i int) Map {
	return this.Data[i]
}

func (this QueryResult) Size() int {
	return len(this.Data)
}
