package gintool

import "math"

//分页对象
type Pager struct {
	Page      int `form:"page"  json:"page"`           //当前页
	PageSize  int `form:"pageSize"  json:"pageSize"`   //每页条数
	Total     int `form:"total"  json:"total"`         //总条数
	PageCount int `form:"pageCount"  json:"pageCount"` //总页数
	NumStart  int `form:"numStart"  json:"numStart"`   //开始序数
}

func CreatePager(page, pagesize int) *Pager {
	if page < 1 {
		page = 1
	}
	if pagesize < 1 {
		pagesize = 10
	}
	pager := new(Pager)
	pager.Page = page
	pager.PageSize = pagesize
	pager.setNumStart()
	return pager
}

func (p *Pager) setNumStart() {
	p.NumStart = (p.Page - 1) * p.PageSize
}

func (p *Pager) SetTotal(total int) {
	p.Total = total
	p.PageCount = int(math.Ceil(float64(total) / float64(p.PageSize)))
}
