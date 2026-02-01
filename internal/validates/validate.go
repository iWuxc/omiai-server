package validates

// Paginate 分页公共参数 .
type Paginate struct {
	Page     int `form:"page" query:"page" json:"page" binding:"numeric" field:"页码"`                  // 页码  默认1
	PageSize int `form:"page_size" query:"page_size" json:"page_size" binding:"numeric" field:"分页大小"` // 每页展示数量   默认20
}

func (p *Paginate) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Paginate) Limit() int {
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	return p.PageSize
}
