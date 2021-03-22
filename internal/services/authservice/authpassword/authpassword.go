package authpassword

import (
	"strings"

	"github.com/foolish15/shorten-url-service/internal/repositories/user"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/foolish15/shorten-url-service/internal/services/authservice"
	"golang.org/x/crypto/bcrypt"
)

//Service implement service interface
type Service struct {
	userRepo user.Repository
}

//New new struct
func New(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

//Auth implement auth
func (sv *Service) Auth(credential map[string]string) (usr schemas.User, err error) {
	eml, ok := credential["email"]
	if !ok {
		return usr, authservice.NewError("Please input email")
	}
	eml = strings.TrimSpace(eml)
	if eml == "" {
		return usr, authservice.NewError("Please input email")
	}

	psswd, ok := credential["password"]
	if !ok {
		return usr, authservice.NewError("Please input password")
	}
	if psswd == "" {
		return usr, authservice.NewError("Please input password")
	}

	usr, err = sv.userRepo.First(user.PreloadUserAuth{}, user.WhereEmail{Email: eml})
	if err != nil {
		if err == user.ErrNotFound {
			return usr, authservice.NewError(user.ErrNotFound.Error())
		}
		return usr, err
	}

	chkChannel := false
	hashPsswd := ""
	for _, uAuth := range usr.UserAuths {
		if uAuth.Channel == schemas.UserAuthChannelPassword {
			chkChannel = true
			hashPsswd = uAuth.ChannelSecret
		}
	}
	if !chkChannel {
		return usr, authservice.NewError("This user not allow to auth with channel password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPsswd), []byte(psswd))
	if err != nil {
		return usr, err
	}

	return usr, nil
}
