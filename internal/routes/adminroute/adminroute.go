package adminroute

import (
	"fmt"
	"os"
	"regexp"

	"github.com/foolish15/shorten-url-service/internal/handlers/authhandler"
	"github.com/foolish15/shorten-url-service/internal/handlers/linkhandler"
	"github.com/foolish15/shorten-url-service/internal/handlers/userhandler"
	"github.com/foolish15/shorten-url-service/internal/repositories/accesstransaction"
	"github.com/foolish15/shorten-url-service/internal/repositories/link"
	"github.com/foolish15/shorten-url-service/internal/repositories/token"
	"github.com/foolish15/shorten-url-service/internal/repositories/user"
	"github.com/foolish15/shorten-url-service/internal/repositories/userauth"
	"github.com/foolish15/shorten-url-service/internal/services/acctxservice"
	"github.com/foolish15/shorten-url-service/internal/services/authservice"
	"github.com/foolish15/shorten-url-service/internal/services/registerservice"
	"github.com/foolish15/shorten-url-service/internal/services/tokenservice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type R struct {
	DB *gorm.DB
}

func skipJWT(ec echo.Context) bool {
	url := fmt.Sprintf("%s[%s]", ec.Request().RequestURI, ec.Request().Method)
	re := regexp.MustCompile(`(?m)^\/admin\/auth\[POST\]$`)
	if len(re.FindStringIndex(url)) > 0 {
		return true //skip validate jwt
	}
	re = regexp.MustCompile(`(?m)^\/admin\/register\[POST\]$`)
	if len(re.FindStringIndex(url)) > 0 {
		return true //skip validate jwt
	}
	re = regexp.MustCompile(`(?m).(?:jpg|jpeg|gif|png|ico|cur|gz|svg|svgz|mp4|ogg|ogv|webm|htc|svg|woff|woff2|ttf)\[GET\]$`)
	return len(re.FindStringIndex(url)) > 0 //skip validate jwt
}

//R implement route
func (r R) Route(e *echo.Echo) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtClaim := tokenservice.SystemClaims{}
	a := e.Group("admin",
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ec echo.Context) error {
				ec.Set("JWT_SECRET", jwtSecret)
				return next(ec)
			}
		},
		middleware.JWTWithConfig(middleware.JWTConfig{
			Skipper:    skipJWT,
			SigningKey: []byte(jwtSecret),
			Claims:     &jwtClaim,
		}))
	a.POST("/register", func(c echo.Context) error {
		db := c.Get("DB").(*gorm.DB)
		userRepo := user.New(db)
		userAuthRepo := userauth.New(db)
		registerSv := &registerservice.S{}
		return userhandler.Register(c, userRepo, userAuthRepo, registerSv)
	})
	a.POST("/auth", func(c echo.Context) error {
		db := c.Get("DB").(*gorm.DB)
		userRepo := user.New(db)
		tknRepo := token.New(db)
		userAuthRepo := userauth.New(db)
		authSv := &authservice.S{}
		tokenSv := &tokenservice.S{}
		return authhandler.Authenticate(c, userRepo, tknRepo, userAuthRepo, authSv, tokenSv)
	})

	a.GET("/links", func(c echo.Context) error {
		db := c.Get("DB").(*gorm.DB)
		linkRepo := link.New(db)
		accTxRepo := accesstransaction.New(db)
		accTxSv := &acctxservice.S{}
		return linkhandler.Redirect(c, linkRepo, accTxRepo, accTxSv)
	})
}
