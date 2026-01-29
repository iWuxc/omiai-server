package biz

// WhereClause .
type WhereClause struct {
	OrderBy string
	Where   string
	Args    []interface{}
}

// Pagination .
type Pagination struct {
	Total       int64 `json:"total"`        // 总条数
	CurrentPage int   `json:"current_page"` // 当前页
	PageSize    int   `json:"page_size"`    // 每页条数
}

const (
	operation = " and "
)

// JoinCondition 拼接查询条件
func JoinCondition(clause *WhereClause, where string, args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case string:
			if len(arg.(string)) == 0 {
				return
			}
		}
	}

	if len(clause.Where) > 0 {
		clause.Where += operation
	}
	clause.Where += where
	clause.Args = append(clause.Args, args...)
}
