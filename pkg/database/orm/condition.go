package orm

import (
	"errors"

	"gorm.io/gorm"
)

type findPageOptions struct {
	countQuery *gorm.DB
	columns    string
	order      string
}

type FindPageOption func(*findPageOptions)

func WithCountQuery(countQuery *gorm.DB) FindPageOption {
	return func(o *findPageOptions) {
		o.countQuery = countQuery
	}
}

func WithFindPageColumns(columns string) FindPageOption {
	return func(o *findPageOptions) {
		o.columns = columns
	}
}

func WithFindPageOrder(order string) FindPageOption {
	return func(o *findPageOptions) {
		o.order = order
	}
}

type findOneOptions struct {
	ignoreErrRecordNotFind bool
	columns                string
	order                  string
}

type FindOneOption func(*findOneOptions)

func WithIgnoreErrRecordNotFind() FindOneOption {
	return func(o *findOneOptions) {
		o.ignoreErrRecordNotFind = true
	}
}

func WithFindOneColumns(columns string) FindOneOption {
	return func(o *findOneOptions) {
		o.columns = columns
	}
}

func WithFindOneOrder(order string) FindOneOption {
	return func(o *findOneOptions) {
		o.order = order
	}
}

// Where .
func Where(query *gorm.DB, cond string, value interface{}) *gorm.DB {
	query = query.Where(cond, value)
	return query
}

// WhereIgnoreBlank .
func WhereIgnoreBlank(query *gorm.DB, cond string, value string) *gorm.DB {
	if value != "" {
		query = query.Where(cond, value)
	}
	return query
}

// WhereIgnoreZero .
func WhereIgnoreZero(query *gorm.DB, cond string, value interface{}) *gorm.DB {
	switch t := value.(type) {
	case int:
		if t <= 0 {
			return query
		}
	case int32:
		if t <= 0 {
			return query
		}
	case int64:
		if t <= 0 {
			return query
		}
	case float64:
		if t <= 0 {
			return query
		}
	case *int:
		if t == nil || *t == -1 {
			return query
		}
	case *int64:
		if t == nil || *t == -1 {
			return query
		}
	case *float64:
		if t == nil || *t == -1 {
			return query
		}
	default:
		return query
	}
	query = query.Where(cond, value)
	return query
}

// WhereLikeIgnoreBlank .
func WhereLikeIgnoreBlank(query *gorm.DB, cond string, value string) *gorm.DB {
	if value != "" {
		query = query.Where(cond, "%"+value+"%")
	}
	return query
}

// FindPage 查询分页数据
func FindPage(query *gorm.DB, pageNum, pageSize int, out interface{}, opts ...FindPageOption) (count int64, err error) {
	var o findPageOptions
	for _, opt := range opts {
		opt(&o)
	}
	if o.countQuery != nil {
		if err = o.countQuery.Count(&count).Error; err != nil {
			return 0, err
		}
	} else {
		if err = query.Count(&count).Error; err != nil {
			return 0, err
		}
	}
	if count == 0 {
		return 0, nil
	}
	// 如果分页大小小于0，则不查询数据
	if pageSize < 0 || pageNum < 0 {
		return count, nil
	}
	if o.columns != "" {
		query = query.Select(o.columns)
	}
	if o.order != "" {
		query = query.Order(o.order)
	}
	if pageNum > 0 && pageSize > 0 {
		query = query.Offset((pageNum - 1) * pageSize)
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize > 0 {
		query = query.Limit(pageSize)
	}
	if err = query.Find(out).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindOne 查询单条数据
func FindOne(query *gorm.DB, out interface{}, opts ...FindOneOption) (exist bool, err error) {
	var o findOneOptions
	for _, opt := range opts {
		opt(&o)
	}
	if o.columns != "" {
		query = query.Select(o.columns)
	}
	if o.order != "" {
		query = query.Order(o.order)
	}
	if err := query.First(out).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && o.ignoreErrRecordNotFind {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckExist 检查数据是否存在
func CheckExist(query *gorm.DB) (exist bool, err error) {
	var count int64
	if err = query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
