package handlers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/foolish15/shorten-url-service/internal/services/tokenservice"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

//ResponseStatus handler response statusssss
type ResponseStatus string

//define response status
const (
	ResponseStatusSuccess ResponseStatus = "success"
	ResponseStatusFail    ResponseStatus = "fail"
	ResponseStatusError   ResponseStatus = "error"
)

//ResponseMessage handler response statusssss
type ResponseMessage string

//define standard message
const (
	ResponseMessageSuccess             ResponseMessage = "Success"
	ResponseMessageOK                  ResponseMessage = "OK"
	ResponseMessageInvalidData         ResponseMessage = "Invalid data"
	ResponseMessageInternalServerError ResponseMessage = "Internal server error"
	ResponseMessageNotfound            ResponseMessage = "Page not found"
	ResponseMessageInvalidJWT          ResponseMessage = "Invalid JWT"
	ResponseMessageInvalidCredentail   ResponseMessage = "Invalid credential"
)

//Response standard response
type Response struct {
	Status  ResponseStatus  `json:"status"`
	Message ResponseMessage `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
}

type TypeTime time.Time

func (t *TypeTime) UnmarshalParam(st string) error {
	st = strings.TrimSpace(st)
	if st == "" {
		return nil
	}
	ti, err := time.ParseInLocation("2006-01-02 15:04:05", st, time.Local)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	*t = TypeTime(ti)
	return nil
}

//UnmarshalJSON  custom json unmarshal
func (t *TypeTime) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}
	st := strings.TrimSpace(strings.Trim(string(data), "\""))
	if st == "" {
		return nil
	}

	ti, err := time.ParseInLocation("2006-01-02 15:04:05", st, time.Local)
	if err != nil {
		return err
	}
	*t = TypeTime(ti)
	return nil
}

//MarshalJSON  custom json unmarshal
func (t *TypeTime) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}

	st := time.Time(*t).Format("2006-01-02 15:04:05")
	return json.Marshal(st)
}

func ResponseInternalServerErrer(c ContextHandler, data interface{}) error {
	return ResponseWithContext(c, http.StatusInternalServerError, ResponseStatusError, ResponseMessageInternalServerError, data)
}

func ResponseNotfound(c ContextHandler) error {
	return ResponseWithContext(c, http.StatusNotFound, ResponseStatusFail, ResponseMessageNotfound, nil)
}

func ResponseInvalidData(c ContextHandler, data interface{}) error {
	return ResponseWithContext(c, http.StatusUnprocessableEntity, ResponseStatusFail, ResponseMessageInvalidData, data)
}

func ResponseGone(c ContextHandler) error {
	return ResponseWithContext(c, http.StatusGone, ResponseStatusFail, ResponseMessageInvalidData, nil)
}

func ResponseInvalidJWT(c ContextHandler) error {
	return ResponseWithContext(c, http.StatusUnprocessableEntity, ResponseStatusFail, ResponseMessageInvalidJWT, nil)
}

func ResponseWithContext(c ContextHandler, httpStatus int, status ResponseStatus, msg ResponseMessage, data interface{}) error {
	resp := Response{
		Status:  status,
		Message: msg,
		Data:    data,
	}
	return c.JSON(httpStatus, resp)
}

func ErrorValidation(caller string, err error, c ContextHandler) error {
	switch err := err.(type) {
	case validator.ValidationErrors:
		L(c).Debugf("[%s.ErrorValidation] c.Validate error: %+v", caller, err)
		errs := []validator.FieldError(err)
		dataE := []string{}
		for _, ers := range errs {
			dataE = append(dataE, ers.Error())
		}
		return ResponseWithContext(c, http.StatusUnprocessableEntity, ResponseStatusFail, ResponseMessageInvalidData, dataE)
	default:
		L(c).Errorf("[%s.ErrorValidation] c.Validate error: %+v", caller, err)
		return ResponseWithContext(c, http.StatusInternalServerError, ResponseStatusError, ResponseMessageInternalServerError, nil)
	}
}

func L(c ContextHandler) *logrus.Entry {
	return logrus.WithContext(c.Request().Context())
}

func UnpackToken(c ContextHandler) (*tokenservice.SystemClaims, error) {
	tk, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("user parse token error")
	}

	claims, ok := tk.Claims.(*tokenservice.SystemClaims)
	if !ok {
		return nil, fmt.Errorf("jwt parse system claims error")
	}

	return claims, nil
}

//ContextHandler type context for handler
type ContextHandler interface {
	Get(string) interface{}
	Request() *http.Request
	QueryParam(string) string
	Bind(interface{}) error
	FormFile(name string) (*multipart.FileHeader, error)
	Set(key string, val interface{})
	Validate(i interface{}) error
	JSON(int, interface{}) error
	Param(name string) string
	FormValue(name string) string
	MultipartForm() (*multipart.Form, error)
	Redirect(int, string) error
}
