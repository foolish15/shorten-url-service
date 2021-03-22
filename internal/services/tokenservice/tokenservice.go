package tokenservice

import (
	"encoding/json"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/foolish15/shorten-url-service/internal/repositories/token"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/pkg/errors"
)

type TokenUser struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// SystemClaims type
type SystemClaims struct {
	U TokenUser `json:"u"`
	jwt.StandardClaims
}

//ErrTokenInterface interface
type ErrTokenInterface interface {
	error
	RegError() string
}

//ErrToken custom error
type ErrToken struct {
	err string
}

//NewError new errro
func NewError(err string) ErrTokenInterface {
	return &ErrToken{
		err: err,
	}
}

//Error implement interface
func (er *ErrToken) Error() string {
	return er.err
}

//RegError implement interface
func (er *ErrToken) RegError() string {
	return er.Error()
}

//I service interface
type I interface {
	Issue(sub schemas.TokenSub, usr schemas.User, exp int64, secret string, tknRepo token.Repository) (t string, err error)
}

//S implement service interface
type S struct{}

//Issue issue token
func (sv *S) Issue(sub schemas.TokenSub, usr schemas.User, exp int64, secret string, tknRepo token.Repository) (t string, err error) {
	u, err := json.Marshal(usr)
	if err != nil {
		return t, errors.Wrap(err, "json marshal usr error")
	}

	var tknUsr TokenUser
	err = json.Unmarshal(u, &tknUsr)
	if err != nil {
		return t, errors.Wrap(err, "json unmarshal tknUsr error")
	}

	tu, err := json.Marshal(tknUsr)
	if err != nil {
		return t, errors.Wrap(err, "json marshal tknUsr error")
	}

	var code *string

	tkn := schemas.Token{
		Sub:  sub,
		U:    tu,
		Exp:  exp,
		Code: code,
	}

	err = tknRepo.Create(&tkn)
	if err != nil {
		return t, err
	}

	claims := SystemClaims{
		U: tknUsr,
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.Itoa(int(tkn.ID)),
			Subject:   string(sub),
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(secret))

	return t, err
}
