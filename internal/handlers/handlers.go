package handlers

import (
	"mime/multipart"
	"net/http"
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
	ResponseMessageSuccess     ResponseMessage = "Success"
	ResponseMessageOK          ResponseMessage = "OK"
	ResponseMessageInvalidData ResponseMessage = "Invalid data"
)

//Response standard response
type Response struct {
	Status  ResponseStatus  `json:"status"`
	Message ResponseMessage `json:"message"`
	Data    interface{}     `json:"data"`
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
