package acctxservice

import (
	"github.com/foolish15/shorten-url-service/internal/repositories/accesstransaction"
	"github.com/foolish15/shorten-url-service/internal/repositories/link"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/mssola/user_agent"
)

type I interface {
	SaveTx(lnk schemas.Link, uas string, linkRepo link.Repository, accTxRepo accesstransaction.Repository) error
}

type S struct{}

func (*S) SaveTx(lnk schemas.Link, uas string, linkRepo link.Repository, accTxRepo accesstransaction.Repository) error {
	isComplete := false
	tx, txLinkRepo, err := linkRepo.StartTransaction() // start db transaction
	if err != nil {
		return err
	}
	defer func() {
		if isComplete {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	txAccTxRepo, err := accTxRepo.New(tx) // new repo with db transaction
	if err != nil {
		return err
	}
	lnk, err = txLinkRepo.First(link.SelectForUpdate{}, link.WhereID{ID: lnk.ID}) // select and lock for update
	if err != nil {
		return err
	}

	ua := user_agent.New(uas)
	browser, browserVer := ua.Browser()
	osInfo := ua.OSInfo()
	os := osInfo.FullName
	osVersion := osInfo.Version
	ua.Platform()

	accTr := schemas.AccessTransaction{
		LinkID:         lnk.ID,
		LinkURL:        lnk.Link,
		Browser:        browser,
		BrowserVersion: browserVer,
		OS:             os,
		OSVersion:      osVersion,
		DeviceType:     ua.Platform(),
		UserAgent:      uas,
	}
	err = txAccTxRepo.Create(&accTr)
	if err != nil {
		return err
	}

	lnk.Hit = lnk.Hit + 1
	err = txLinkRepo.Update(&lnk)
	if err != nil {
		return err
	}
	isComplete = true
	return nil
}
