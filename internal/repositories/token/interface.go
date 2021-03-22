package token

import (
	"database/sql"
	"errors"

	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/schemas"
)

//define error case
var (
	ErrNotFound = errors.New("Token Not Found")
)

//Repository repo interface
type Repository interface {
	repositories.SetDBInterface
	New(idb ...interface{}) (Repository, error)
	StartTransaction() (tx *sql.Tx, rp Repository, err error)
	Create(tkn *schemas.Token) error
	Update(tkn *schemas.Token) error
	Delete(tknID uint) error
	Deletes(queries ...repositories.Query) error
	First(queries ...repositories.Query) (tkn schemas.Token, err error)
	Find(queries ...repositories.Query) (tkns []schemas.Token, err error)
	Count(queries ...repositories.Query) (cnt int64, err error)
}
