package pagination

type Page struct {
	PageNum  int `form:"page_num" json:"page_num"`   // 页码
	PageSize int `form:"page_size" json:"page_size"` // 页大小
}

func GetPageOffset(pageNum, pageSize int64) int64 {
	return (pageNum - 1) * pageSize
}
