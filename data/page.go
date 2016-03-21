package data

import (
	"bytes"
)

//------------------Page begin-------------------------------
//分页
type Page struct {
	FirstResult int    //开始位置
	MaxResults  int    //查询几条
	PageNo      int    //当前页号
	PageSize    int    //每页数
	TotalRow    int    //查询结果：总数量
	TotalPage   int    //查询结果：总页数
	HtmlTag     string //分页标签

	NextPage     bool // 是否有下一页
	PreviousPage bool // 是否有前一页
}

func (this *Page) New(pageNo, pageSize int) {
	if pageNo <= 0 {
		pageNo = 1
	}
	this.FirstResult = (pageNo - 1) * pageSize
	this.MaxResults = pageSize
	this.PageNo = pageNo
	this.PageSize = pageSize
}

func (this *Page) Html() {
	if (this.TotalRow % this.PageSize) == 0 {
		this.TotalPage = this.TotalRow / this.PageSize
	} else {
		this.TotalPage = (this.TotalRow / this.PageSize) + 1
	}
	if this.PageNo >= this.TotalPage {
		this.NextPage = false
		this.PageNo = this.TotalPage
	} else {
		this.NextPage = true
	}
	var htmlTag bytes.Buffer

	this.HtmlTag = htmlTag.String()
}

//------------------Page end-------------------------------
