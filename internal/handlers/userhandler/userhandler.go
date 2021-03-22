package userhandler

import (
	"net/http"
	"strings"

	"github.com/foolish15/shorten-url-service/internal/handlers"
	"github.com/foolish15/shorten-url-service/internal/repositories/user"
	"github.com/foolish15/shorten-url-service/internal/repositories/userauth"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/foolish15/shorten-url-service/internal/services/registerservice"
	"github.com/foolish15/shorten-url-service/internal/services/registerservice/registerwithpassword"
	"github.com/sirupsen/logrus"
)

type requestRegister struct {
	Name           string                  `json:"name" form:"name"`
	Email          string                  `json:"email" form:"email" validate:"omitempty,email"`
	Phone          string                  `json:"phone" form:"phone" validate:"omitempty,len=10,numeric"`
	AuthChannel    schemas.UserAuthChannel `json:"authChannel" form:"authChannel"`
	AuthCredential map[string]string       `json:"authCredential" from:"authCredential"`
}

//Register create user
func Register(c handlers.ContextHandler, usrRepo user.Repository, usrAuthRepo userauth.Repository, regisSv registerservice.I) error {
	req := requestRegister{}
	err := c.Bind(&req)
	if err != nil {
		handlers.L(c).Debugf("[UserHandler.Register] c.Bind error: %+v", err)
		return handlers.ResponseInvalidData(c, "Invalid input")
	}
	err = c.Validate(req)
	if err != nil {
		return handlers.ErrorValidation("UserHandler.Register", err, c)
	}
	switch req.AuthChannel {
	case schemas.UserAuthChannelPassword:
		regisSv.RegisterChannel(registerwithpassword.New(usrRepo, usrAuthRepo))
	default:
		handlers.L(c).Debugf("[UserHandler.Register] switch req.Channel: %v", req.AuthChannel)
		return handlers.ResponseInvalidData(c, "Invalid register channel")
	}
	usr := schemas.User{}
	if name := strings.TrimSpace(req.Name); name != "" {
		usr.Name = &name
	}
	if email := strings.TrimSpace(req.Email); email != "" {
		usr.Email = &email
	}
	if phone := strings.TrimSpace(req.Phone); phone != "" {
		usr.Phone = &phone
	}
	usr, err = regisSv.Register(usr.Name, usr.Email, usr.Phone, req.AuthCredential)
	if err != nil {
		switch err.(type) {
		case registerservice.ErrRegisterInterface:
			handlers.L(c).
				WithFields(
					logrus.Fields{
						"name":       usr.Name,
						"email":      usr.Email,
						"phone":      usr.Phone,
						"credential": req.AuthCredential,
					},
				).
				Debugf("[UserHandler.Register] regisSv.Register(usr.Name, usr.Email, usr.Phone, req.Credential) error: %+v", err)
			return handlers.ResponseInvalidData(c, err.Error())
		default:
			handlers.L(c).
				WithFields(
					logrus.Fields{
						"name":       usr.Name,
						"email":      usr.Email,
						"phone":      usr.Phone,
						"credential": req.AuthCredential,
					},
				).
				Errorf("[UserHandler.Register] regisSv.Register(usr.Name, usr.Email, usr.Phone, req.Credential) error: %+v", err)
			return handlers.ResponseInternalServerErrer(c, nil)
		}
	}

	return c.JSON(http.StatusCreated, usr)
}
