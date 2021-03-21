package publicroute

import (
	"github.com/foolish15/shorten-url-service/internal/handlers/linkhandler"
	"github.com/foolish15/shorten-url-service/internal/repositories/accesstransaction"
	"github.com/foolish15/shorten-url-service/internal/repositories/link"
	"github.com/foolish15/shorten-url-service/internal/services/acctxservice"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type R struct {
	DB *gorm.DB
}

//R implement route
func (r R) Route(e *echo.Echo) {
	api := e.Group("p")

	api.GET("/:code", func(c echo.Context) error {
		db := c.Get("DB").(*gorm.DB)
		linkRepo := link.New(db)
		accTxRepo := accesstransaction.New(db)
		accTxSv := &acctxservice.S{}
		return linkhandler.Redirect(c, linkRepo, accTxRepo, accTxSv)
	})
}
