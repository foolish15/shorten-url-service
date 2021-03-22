package registerwithpassword

import (
	"strconv"

	"github.com/foolish15/shorten-url-service/internal/repositories/user"
	"github.com/foolish15/shorten-url-service/internal/repositories/userauth"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/foolish15/shorten-url-service/internal/services/registerservice"
	"golang.org/x/crypto/bcrypt"
)

//Channel register with password channel
type Channel struct {
	userRepo     user.Repository
	userAuthRepo userauth.Repository
}

//New return register channel
func New(userRepo user.Repository, userAuthRepo userauth.Repository) *Channel {
	return &Channel{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}

//Register implement interface
func (c *Channel) Register(name, email, phone *string, credential map[string]string) (usr schemas.User, err error) {
	passwd, ok := credential["password"]
	if !ok {
		return usr, registerservice.NewError("Please set credential password")
	}

	if email != nil && *email != "" {
		cnt, err := c.userRepo.Count(user.WhereEmail{Email: *email})
		if err != nil {
			return usr, err
		}
		if cnt > 0 {
			return usr, registerservice.ErrDuplicateEmail
		}
	}

	isComplete := false
	tx, txUserRepo, err := c.userRepo.StartTransaction()
	if err != nil {
		return usr, err
	}
	defer func() {
		if isComplete {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	txUserAuthRepo, err := c.userAuthRepo.New(tx)
	if err != nil {
		return usr, err
	}

	usr = schemas.User{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	err = txUserRepo.Create(&usr)
	if err != nil {
		return usr, err
	}

	hpsswd, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return usr, err
	}
	usrA := schemas.UserAuth{
		UserID:        usr.ID,
		Channel:       schemas.UserAuthChannelPassword,
		ChannelID:     strconv.Itoa(int(usr.ID)),
		ChannelSecret: string(hpsswd),
	}
	err = txUserAuthRepo.Create(&usrA)
	if err != nil {
		return usr, err
	}

	isComplete = true
	return usr, err
}
