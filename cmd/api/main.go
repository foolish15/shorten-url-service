package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/foolish15/shorten-url-service/internal/routes"
	"github.com/foolish15/shorten-url-service/internal/routes/adminroute"
	"github.com/foolish15/shorten-url-service/internal/routes/apiroute"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/foolish15/shorten-url-service/pkg/logrus/hook/writer"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	logrus.SetOutput(ioutil.Discard) // Send all logs to nowhere by default
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
	logrus.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.TraceLevel,
		},
	})

	err := godotenv.Load()
	if err != nil {
		logrus.Debugf("Get environment variable from OS")
	}
	setLogLevel()
}

func setLogLevel() {
	lv := os.Getenv("LOG_LEVEL")
	switch lv {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fetal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func connectDB() *gorm.DB {
	connectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	logrus.Debugf(connectString)
	for {
		db, err := gorm.Open(
			mysql.Open(connectString),
			&gorm.Config{},
		)

		if err != nil {
			time.Sleep(5 * time.Second)
			logrus.Errorf("ConnectDB error: %+v", err)
		} else {
			return db
		}
	}
}

//CustomValidator define custom validator
type CustomValidator struct {
	validator *validator.Validate
}

//Validate implement validate
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	validate := validator.New()
	e.Validator = &CustomValidator{validator: validate}
	db := connectDB()
	err := db.AutoMigrate(
		schemas.AccessTransaction{},
		schemas.Block{},
		schemas.Link{},
		schemas.Token{},
		schemas.User{},
		schemas.UserAuth{},
	)
	if err != nil {
		logrus.Errorf("Auto migrate error: %+v", err)
		os.Exit(1)
	}

	e.Use(
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ec echo.Context) error {
				ec.Set("DB", db)
				return next(ec)
			}
		},
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		}),
		routes.MiddlewareRequestID(),
		routes.MiddlewareBodyDump(),
	)
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "")
	})
	routes.Register(e, apiroute.R{DB: db}, adminroute.R{DB: db})

	defaultIP := "0.0.0.0"
	ip := os.Getenv("SERVICE_IP")
	if ip == "" {
		ip = defaultIP
	}
	if net.ParseIP(ip) == nil { //validate IP
		ip = defaultIP
	}

	defaultPort := "80"
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}
	_, err = strconv.Atoi(port)
	if err != nil {
		port = defaultPort
	}

	err = e.Start(fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		logrus.Errorf("e.Start error: %+v", err)
	}
}
