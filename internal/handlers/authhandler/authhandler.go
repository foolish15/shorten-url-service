package authhandler

import (
	"net/http"
	"time"

	"github.com/foolish15/shorten-url-service/internal/handlers"
	"github.com/foolish15/shorten-url-service/internal/repositories/token"
	"github.com/foolish15/shorten-url-service/internal/repositories/user"
	"github.com/foolish15/shorten-url-service/internal/repositories/userauth"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/foolish15/shorten-url-service/internal/services/authservice"
	"github.com/foolish15/shorten-url-service/internal/services/authservice/authpassword"
	"github.com/foolish15/shorten-url-service/internal/services/tokenservice"
	"github.com/sirupsen/logrus"
)

var expireToken = (3 * time.Hour)

type requestAuth struct {
	Channel    schemas.UserAuthChannel `form:"channel" json:"channel" validate:"required,oneof=password"`
	Credential map[string]string       `form:"credential" json:"credential"`
}

//Authenticate handler auth
func Authenticate(c handlers.ContextHandler, userRepo user.Repository, tknRepo token.Repository, userAuthRepo userauth.Repository, authSv authservice.I, tknSV tokenservice.I) (err error) {
	jwtSecrt := c.Get("JWT_SECRET").(string)
	req := requestAuth{}
	err = c.Bind(&req)
	if err != nil {
		handlers.L(c).Debugf("[AuthHandler.Authenticate] c.Bind error: %+v", err)
		return handlers.ResponseInvalidData(c, "Invalid input")
	}
	err = c.Validate(req)
	if err != nil {
		return handlers.ErrorValidation("AuthHandler.Authenticate", err, c)
	}

	switch req.Channel {
	case schemas.UserAuthChannelPassword:
		authSv.RegisterChannel(authpassword.New(userRepo))
	default:
		handlers.L(c).Debugf("[AuthHandler.Authenticate] switch req.Channel: %v", req.Channel)
		return handlers.ResponseInvalidData(c, "Invalid auth channel")
	}

	usr, err := authSv.Auth(req.Credential)
	if err != nil {
		switch err.(type) {
		case authservice.ErrAuthInterface:
			handlers.L(c).Debugf("[AuthHandler.Authenticate] authSv.Auth(req.ChannelCredential) error: %+v", err)
			return handlers.ResponseWithContext(c, http.StatusUnauthorized, handlers.ResponseStatusFail, handlers.ResponseMessageInvalidCredentail, nil)
		case error:
			handlers.L(c).Errorf("[AuthHandler.Authenticate] authSv.Auth(req.ChannelCredential) error: %+v", err)
			return handlers.ResponseInternalServerErrer(c, nil)
		}
	}

	exp := time.Now().Add(expireToken).Unix()
	tkn, err := tknSV.Issue(schemas.TokenSubUser, usr, exp, jwtSecrt, tknRepo)
	if err != nil {
		handlers.L(c).
			WithFields(
				logrus.Fields{
					"usr":     usr,
					"exp":     exp,
					"tknRepo": tknRepo,
				},
			).
			Errorf("[AuthHandler.Authenticate] tknSV.Issue error: %v", err)
		return handlers.ResponseInternalServerErrer(c, nil)
	}

	response := map[string]interface{}{
		"access_token": tkn,
		"expires_in":   exp,
		"token_type":   "bearer",
	}
	return c.JSON(http.StatusOK, response)
}
