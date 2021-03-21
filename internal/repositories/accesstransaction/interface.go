package accesstransaction

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
	Create(acc *schemas.AccessTransaction) error
	Update(acc *schemas.AccessTransaction) error
	Delete(ID uint) error
	First(queries ...repositories.Query) (acc schemas.AccessTransaction, err error)
	Find(queries ...repositories.Query) (accs []schemas.AccessTransaction, err error)
	Count(queries ...repositories.Query) (cnt int64, err error)
}
