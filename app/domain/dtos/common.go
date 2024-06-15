package dtos

type PageQueryDto struct {
	Page     int `json:"page" form:"page" binding:"required"`           // 页数
	PageSize int `json:"page_size" form:"page_size" binding:"required"` // 每页条数
}

type PageListResDto struct {
	Page      int         `json:"page"`       // 页数
	PageSize  int         `json:"page_size"`  // 每页条数
	Count     int64       `json:"count"`      // 总条数
	PageCount int         `json:"page_count"` // 总页数
	List      interface{} `json:"list"`       // 列表
}
