package link

import (
	"database/sql"
	"errors"

	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/schemas"
)

//define error case
var (
	ErrNotFound = errors.New("link not found")
)

//Repository repo interface
type Repository interface {
	repositories.SetDBInterface
	repositories.OrderInterface
	New(idb ...interface{}) (Repository, error)
	StartTransaction() (tx *sql.Tx, txRepo Repository, err error)
	Create(lnk *schemas.Link) error
	Update(lnk *schemas.Link) error
	Delete(ID uint) error
	First(queries ...repositories.Query) (lnk schemas.Link, err error)
	Find(queries ...repositories.Query) (lnks []schemas.Link, err error)
	Count(queries ...repositories.Query) (cnt int64, err error)
}
