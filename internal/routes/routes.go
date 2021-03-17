package routes

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"github.com/sirupsen/logrus"
)

//RouteInterface interface route
type RouteInterface interface {
	Route(e *echo.Echo)
}

//Register register router
func Register(e *echo.Echo, rs ...RouteInterface) {
	for _, r := range rs {
		r.Route(e)
	}
}

//MiddlewareBodyDump dump request response
func MiddlewareBodyDump() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		var reqURI = c.Request().RequestURI
		var skipRegx = regexp.MustCompile(`.(?:jpg|jpeg|gif|png|ico|cur|gz|svg|svgz|mp4|ogg|ogv|webm|htc|svg|woff|woff2|ttf)$`)
		if len(skipRegx.FindStringIndex(reqURI)) > 0 {
			return
		}

		route := fmt.Sprintf("%s[%s]", reqURI, c.Request().Method)
		logrus.WithContext(c.Request().Context()).WithFields(logrus.Fields{
			"timeEnd":    time.Now(),
			"uri":        route,
			"reqIP":      c.RealIP(),
			"reqHeader":  c.Request().Header,
			"reqBody":    fmt.Sprintf("%.512s", reqBody),
			"respStatus": c.Response().Status,
			"respHeader": c.Response().Header(),
			"respBody":   fmt.Sprintf("%.1024s", resBody),
		}).Infof("response")
	})
}

//MiddlewareRequestID generate request id header
func MiddlewareRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			req := ec.Request()
			res := ec.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = fmt.Sprintf("%d-%s", time.Now().UnixNano(), random.String(32))
				req.Header.Set(echo.HeaderXRequestID, rid)
			}
			res.Header().Set(echo.HeaderXRequestID, rid)
			ec.SetRequest(
				ec.Request().WithContext(context.WithValue(ec.Request().Context(), echo.HeaderXRequestID, rid)),
			)

			return next(ec)
		}
	}
}
