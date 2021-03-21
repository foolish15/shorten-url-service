package blockservice

import (
	"net/url"
	"regexp"

	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/sirupsen/logrus"
)

//I interface define block service
type I interface {
	IsBlock(u string, lblock []schemas.Block) bool
}

//S struct implement block interface
type S struct {
}

func (*S) IsBlock(u string, lblock []schemas.Block) bool {
	ul, err := url.Parse(u)
	if err != nil {
		return true //block if cannot parse to url
	}
	if ul.Host == "" || ul.Scheme == "" {
		return true //block if cannot parse to url
	}

	for _, block := range lblock {
		switch block.Type {
		case schemas.BlockTypeScheme:
			if IsBlockScheme(ul, block.Value) {
				logrus.Tracef("[blockservice.IsBlock] loop for _, block := range block.Type[%v] block.Value[%v] blocked", block.Type, block.Value)
				return true
			}
		case schemas.BlockTypeDomain:
			if IsBlockDomain(ul, block.Value) {
				logrus.Tracef("[blockservice.IsBlock] loop for _, block := range block.Type[%v] block.Value[%v] blocked", block.Type, block.Value)
				return true
			}
		case schemas.BlockTypeRegex:
			if IsBlockRegex(u, block.Value) {
				logrus.Tracef("[blockservice.IsBlock] loop for _, block := range block.Type[%v] block.Value[%v] blocked", block.Type, block.Value)
				return true
			}
		}
		logrus.Tracef("[blockservice.IsBlock] loop for _, block := range block.Type[%v] block.Value[%v] pass", block.Type, block.Value)
	}

	return false
}

func IsBlockScheme(u *url.URL, block string) bool {
	return u.Scheme == block
}

func IsBlockDomain(u *url.URL, block string) bool {
	return u.Host == block
}

func IsBlockRegex(u string, reg string) bool {
	skipRegx := regexp.MustCompile(reg)
	if skipRegx == nil {
		return false
	}

	dcode, err := url.QueryUnescape(u)
	if err != nil {
		return len(skipRegx.FindStringIndex(u)) > 0
	}
	return len(skipRegx.FindStringIndex(dcode)) > 0 || len(skipRegx.FindStringIndex(u)) > 0
}
