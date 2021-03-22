package authservice

import (
	"fmt"

	"github.com/foolish15/shorten-url-service/internal/schemas"
)

//ErrAuthInterface interface
type ErrAuthInterface interface {
	error
	RegError() string
}

//ErrAuth custom error
type ErrAuth struct {
	err string
}

//NewError new errro
func NewError(err string) ErrAuthInterface {
	return &ErrAuth{
		err: err,
	}
}

//Error implement interface
func (er *ErrAuth) Error() string {
	return er.err
}

//RegError implement interface
func (er *ErrAuth) RegError() string {
	return er.Error()
}

//ServiceAuthChannelInterface auth channel
type ServiceAuthChannelInterface interface {
	Auth(credential map[string]string) (usr schemas.User, err error)
}

//I service interface
type I interface {
	ServiceAuthChannelInterface
	RegisterChannel(authChannel ServiceAuthChannelInterface)
}

//S auth service
type S struct {
	authChan ServiceAuthChannelInterface
}

//RegisterChannel implement service interface
func (sv *S) RegisterChannel(authChannel ServiceAuthChannelInterface) {
	sv.authChan = authChannel
}

//Auth implement interface auth
func (sv *S) Auth(credential map[string]string) (usr schemas.User, err error) {
	if sv.authChan == nil {
		return usr, fmt.Errorf("does not register auth channel yet")
	}
	return sv.authChan.Auth(credential)
}
