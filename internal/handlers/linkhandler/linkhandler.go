package linkhandler

import (
	"net/http"
	"time"

	"github.com/foolish15/shorten-url-service/internal/handlers"
	"github.com/foolish15/shorten-url-service/internal/repositories/accesstransaction"
	"github.com/foolish15/shorten-url-service/internal/repositories/block"
	"github.com/foolish15/shorten-url-service/internal/repositories/link"
	"github.com/foolish15/shorten-url-service/internal/services/acctxservice"
	"github.com/foolish15/shorten-url-service/internal/services/blockservice"
	"github.com/foolish15/shorten-url-service/internal/services/linkservice"
)

type reqCreate struct {
	URL    string            `json:"url" form:"url" validate:"required,url"`
	Expire handlers.TypeTime `json:"expire" form:"expire" `
}

type reqRedirect struct {
	Code string `param:"code"`
}

func Create(c handlers.ContextHandler, linkRepo link.Repository, blockRepo block.Repository, blockSv blockservice.I, linkSv linkservice.I) error {
	req := reqCreate{}
	err := c.Bind(&req)
	if err != nil {
		handlers.L(c).Debugf("[linkhandler.Create] Bind input error: %+v", err)
		return handlers.ResponseInvalidData(c, "Invalid input")
	}

	err = c.Validate(req)
	if err != nil {
		return handlers.ErrorValidation("linkhandler.Create", err, c)
	}

	lblock, err := blockRepo.Find()
	if err != nil {
		handlers.L(c).Errorf("[linkhandler.Create] blockRepo.Find error: %+v", err)
		return handlers.ResponseInternalServerErrer(c, nil)
	}

	if blockSv.IsBlock(req.URL, lblock) {
		handlers.L(c).Debugf("[linkhandler.Create] url is block")
		return handlers.ResponseInvalidData(c, "URL block")
	}

	expUnix := int64(0)
	tExp := time.Time(req.Expire)
	now := time.Now()
	if tExp.Unix() > 0 {
		expUnix = tExp.Unix()
	}

	if expUnix > 0 && now.After(tExp) {
		handlers.L(c).Debugf("[linkhandler.Create] url is block")
		return handlers.ResponseInvalidData(c, "Expire was passed")
	}

	lnk, err := linkSv.CreateLink(req.URL, expUnix, linkRepo)
	if err != nil {
		handlers.L(c).Errorf("[linkhandler.Create] linkSv.CreateLink error: %+v", err)
		return handlers.ResponseInternalServerErrer(c, "Create link failed")
	}

	return c.JSON(http.StatusCreated, lnk)
}

func Redirect(c handlers.ContextHandler, linkRepo link.Repository, accTxRepo accesstransaction.Repository, accTxSv acctxservice.I) error {
	req := reqRedirect{}
	err := c.Bind(&req)
	if err != nil {
		handlers.L(c).Debugf("[linkhandler.Redirect] Bind input error: %+v", err)
		return handlers.ResponseInvalidData(c, "Invalid input")
	}

	lnk, err := linkRepo.First(link.SelectUnscope{}, link.WhereCode{Code: req.Code})
	if err != nil {
		if err == link.ErrNotFound {
			return handlers.ResponseNotfound(c)
		}
		handlers.L(c).Errorf("[linkhandler.Redirect] linkRepo.First error: %+v", err)
		return handlers.ResponseInternalServerErrer(c, "Find link failed")
	}

	if lnk.DeletedAt.Valid {
		handlers.L(c).Debugf("[linkhandler.Redirect] link gone cause delete")
		return handlers.ResponseGone(c)
	}

	now := time.Now().Unix()
	if lnk.Expire > 0 && now > lnk.Expire {
		handlers.L(c).Debugf("[linkhandler.Redirect] link gone cause expired")
		return handlers.ResponseGone(c)
	}

	err = accTxSv.SaveTx(lnk, c.Request().Header.Get("User-Agent"), linkRepo, accTxRepo)
	if err != nil {
		handlers.L(c).Errorf("[linkhandler.Redirect] epo.Create error: %+v", err)
		return handlers.ResponseInternalServerErrer(c, "Create transaction failed")
	}

	return c.Redirect(http.StatusFound, lnk.Link)
}
