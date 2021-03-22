package user

import (
	"database/sql"
	"errors"

	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/schemas"
)

//define error case
var (
	ErrNotFound = errors.New("user not found")
)

//Repository repo interface
type Repository interface {
	repositories.SetDBInterface
	New(idb ...interface{}) (Repository, error)
	StartTransaction() (tx *sql.Tx, txRepo Repository, err error)
	Create(usr *schemas.User) error
	Update(usr *schemas.User) error
	Delete(id uint) error
	Deletes(queries ...repositories.Query) error
	First(queries ...repositories.Query) (usr schemas.User, err error)
	Find(queries ...repositories.Query) (usrs []schemas.User, err error)
	Count(queries ...repositories.Query) (cnt int64, err error)
}
