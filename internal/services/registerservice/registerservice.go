package registerservice

import (
	"fmt"

	"github.com/foolish15/shorten-url-service/internal/schemas"
)

//define error case
var (
	ErrDuplicateEmail = NewError("This email already registered")
	ErrDuplicatePhone = NewError("This phone number already registered")
)

//ErrRegisterInterface interface
type ErrRegisterInterface interface {
	error
	RegError() string
}

//ErrRegister custom error
type ErrRegister struct {
	err string
}

//NewError new errro
func NewError(err string) ErrRegisterInterface {
	return &ErrRegister{
		err: err,
	}
}

//Error implement interface
func (er *ErrRegister) Error() string {
	return er.err
}

//RegError implement interface
func (er *ErrRegister) RegError() string {
	return er.Error()
}

//ServiceRegisterChannelInterface register channel
type ServiceRegisterChannelInterface interface {
	Register(name, email, phone *string, credential map[string]string) (usr schemas.User, err error)
}

//I service interface
type I interface {
	ServiceRegisterChannelInterface
	RegisterChannel(regChannel ServiceRegisterChannelInterface)
}

//S implement service interface
type S struct {
	regChan ServiceRegisterChannelInterface
}

//RegisterChannel implement service interface
func (sv *S) RegisterChannel(regChan ServiceRegisterChannelInterface) {
	sv.regChan = regChan
}

//Register implement interface auth
func (sv *S) Register(name, email, phone *string, credential map[string]string) (usr schemas.User, err error) {
	if sv.regChan == nil {
		return usr, fmt.Errorf("does not register register channel yet")
	}
	return sv.regChan.Register(name, email, phone, credential)
}
