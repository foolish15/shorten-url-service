package repositories

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//LimitOffset type for limit query
type LimitOffset struct {
	Limit, Offset int
}

//DB implement interface
func (l LimitOffset) DB(db DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Limit(l.Limit).Offset(l.Offset)
}

//LimitPage type for limit query
type LimitPage struct {
	Limit, Page int
}

//DB implement interface
func (l LimitPage) DB(db DB) *gorm.DB {
	g := db.(*gorm.DB)

	// Page
	if l.Page == 0 {
		l.Page = 1
	}

	// Limit
	if l.Limit == 0 {
		l.Limit = 10
	}

	offset := (l.Page - 1) * l.Limit
	return g.Limit(l.Limit).Offset(offset)
}

//OrderInterface interface implement for order
type OrderInterface interface {
	TranslateOrderField(input string) (output string)
}

//OrderBy type for order
type OrderBy struct {
	Order *string
	Repo  OrderInterface
}

//DB implement interface
func (l OrderBy) DB(db DB) *gorm.DB {
	g := db.(*gorm.DB)
	if l.Repo == nil {
		return g
	}
	if l.Order != nil && *l.Order != "" {
		orders := strings.Split(*l.Order, ",")
		for _, v := range orders {
			values := strings.Split(v, ":")
			field := l.Repo.TranslateOrderField(values[0])
			if field == "" {
				continue
			}
			sort := "desc"
			vlen := len(values)
			if vlen == 1 || (vlen > 1 && values[1] != "desc") {
				sort = "asc"
			}

			g = g.Order(fmt.Sprintf("%s %s", field, sort))
		}
	}
	return g
}
