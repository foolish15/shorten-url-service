package block

import (
	"database/sql"
	"errors"

	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/schemas"
)

//define error case
var (
	ErrNotFound = errors.New("block not found")
)

//Repository repo interface
type Repository interface {
	repositories.SetDBInterface
	repositories.OrderInterface
	New(idb ...interface{}) (Repository, error)
	StartTransaction() (tx *sql.Tx, txRepo Repository, err error)
	Create(bl *schemas.Block) error
	Update(bl *schemas.Block) error
	Delete(t schemas.BlockType, v string) error
	First(queries ...repositories.Query) (bl schemas.Block, err error)
	Find(queries ...repositories.Query) (bls []schemas.Block, err error)
	Count(queries ...repositories.Query) (cnt int64, err error)
}
