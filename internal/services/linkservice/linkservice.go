package linkservice

import (
	"github.com/foolish15/shorten-url-service/internal/repositories/link"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/foolish15/shorten-url-service/pkg/hashid"
)

type I interface {
	CreateLink(u string, e int64, linkRepo link.Repository) (*schemas.Link, error)
}

type S struct{}

func (S) CreateLink(u string, e int64, linkRepo link.Repository) (*schemas.Link, error) {
	isComplete := false
	tx, txRepo, err := linkRepo.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer func() {
		if isComplete {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	lnk := schemas.Link{
		Link:   u,
		Expire: e,
	}
	err = txRepo.Create(&lnk)
	if err != nil {
		return nil, err
	}

	lnk.Code, err = hashid.Encrypt(int(lnk.ID))
	if err != nil {
		return nil, err
	}
	err = txRepo.Update(&lnk)
	if err != nil {
		return nil, err
	}
	isComplete = true
	return &lnk, nil
}
