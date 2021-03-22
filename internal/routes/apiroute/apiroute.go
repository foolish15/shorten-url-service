package apiroute

import (
	"github.com/foolish15/shorten-url-service/internal/handlers/linkhandler"
	"github.com/foolish15/shorten-url-service/internal/repositories/accesstransaction"
	"github.com/foolish15/shorten-url-service/internal/repositories/block"
	"github.com/foolish15/shorten-url-service/internal/repositories/link"
	"github.com/foolish15/shorten-url-service/internal/services/acctxservice"
	"github.com/foolish15/shorten-url-service/internal/services/blockservice"
	"github.com/foolish15/shorten-url-service/internal/services/linkservice"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type R struct {
	DB *gorm.DB
}

//R implement route
func (r R) Route(e *echo.Echo) {
	{ // group api
		a := e.Group("api")
		a.POST("/links", func(c echo.Context) error {
			db := c.Get("DB").(*gorm.DB)
			linkRepo := link.New(db)
			blockRepo := block.New(db)
			linkSv := &linkservice.S{}
			blockSv := &blockservice.S{}
			return linkhandler.Create(c, linkRepo, blockRepo, blockSv, linkSv)
		})
	}

	{ // group link
		l := e.Group("link")
		l.GET("/:code", func(c echo.Context) error {
			db := c.Get("DB").(*gorm.DB)
			linkRepo := link.New(db)
			accTxRepo := accesstransaction.New(db)
			accTxSv := &acctxservice.S{}
			return linkhandler.Redirect(c, linkRepo, accTxRepo, accTxSv)
		})
	}
}
