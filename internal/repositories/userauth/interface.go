package userauth

import (
	"database/sql"
	"errors"

	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/schemas"
)

//define error case
var (
	ErrNotFound = errors.New("User Auth Not Found")
)

//Repository repo interface
type Repository interface {
	repositories.SetDBInterface
	New(idb ...interface{}) (Repository, error)
	StartTransaction() (tx *sql.Tx, txRepo Repository, err error)
	Create(usrA *schemas.UserAuth) error
	Update(usrA *schemas.UserAuth) error
	Delete(bannerID uint) error
	Deletes(queries ...repositories.Query) error
	First(queries ...repositories.Query) (usrA schemas.UserAuth, err error)
	Find(queries ...repositories.Query) (usrAs []schemas.UserAuth, err error)
	Count(queries ...repositories.Query) (cnt int64, err error)
}
