package linkhandler

import (
	"net/http"
	"time"

	"github.com/foolish15/shorten-url-service/internal/handlers"
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/repositories/accesstransaction"
	"github.com/foolish15/shorten-url-service/internal/repositories/block"
	"github.com/foolish15/shorten-url-service/internal/repositories/link"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/foolish15/shorten-url-service/internal/services/acctxservice"
	"github.com/foolish15/shorten-url-service/internal/services/blockservice"
	"github.com/foolish15/shorten-url-service/internal/services/linkservice"
	"github.com/foolish15/shorten-url-service/pkg/paging"
)

type reqCreate struct {
	URL    string            `json:"url" form:"url" validate:"required,url"`
	Expire handlers.TypeTime `json:"expire" form:"expire" `
}

type reqRedirect struct {
	Code string `param:"code"`
}

type reqList struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type reqDelete struct {
	ID uint `param:"id"`
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

func List(c handlers.ContextHandler, linkRepo link.Repository) error {
	req := reqList{}
	err := c.Bind(&req)
	if err != nil {
		handlers.L(c).Debugf("[linkhandler.List] c.Bind error: %+v", err)
		return handlers.ResponseInvalidData(c, "Invalid input")
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 200
	}

	count, err := linkRepo.Count()
	if err != nil {
		handlers.L(c).Errorf("[linkhandler.List] count product error: %+v", err)
		return handlers.ResponseInternalServerErrer(c, nil)
	}

	var data []schemas.Link
	if count == 0 {
		data = []schemas.Link{}
	} else {
		data, err = linkRepo.Find(repositories.LimitPage{Page: req.Page, Limit: req.Limit})
		if err != nil {
			handlers.L(c).Errorf("[ProductHandler.List] find products error: %+v", err)
			return handlers.ResponseInternalServerErrer(c, nil)
		}
	}

	resp := paging.Pack(c.Request().URL, req.Page, req.Limit, (req.Page-1)*req.Limit, int(count), data)
	return c.JSON(http.StatusOK, resp)
}

func Delete(c handlers.ContextHandler, linkRepo link.Repository) error {
	req := reqDelete{}
	err := c.Bind(&req)
	if err != nil {
		handlers.L(c).Debugf("[linkhandler.Delete] c.Bind error: %+v", err)
		return handlers.ResponseInvalidData(c, "Invalid input")
	}

	err = linkRepo.Delete(req.ID)
	if err != nil {
		if err == link.ErrNotFound {
			return handlers.ResponseNotfound(c)
		}
		handlers.L(c).Errorf("[linkhandler.Delete] linkRepo.Delete error: %+v", err)
		return handlers.ResponseInternalServerErrer(c, "Delete failed")
	}

	return c.JSON(http.StatusNoContent, nil)
}
