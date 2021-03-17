package paging

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strings"
)

var defaultResponse ResponseInterface

func init() {
	ResetResponse()
}

//SetResponse set struct of response
func SetResponse(resp ResponseInterface) {
	defaultResponse = resp
}

//GetResponse get struct of response
func GetResponse() ResponseInterface {
	return defaultResponse
}

//ResetResponse get struct of response
func ResetResponse() {
	defaultResponse = &Response{}
}

//Pack pack everything return as pagination
func Pack(url *url.URL, page, limit, offset, count int, datas interface{}) (resp ResponseInterface) {
	query := url.Query()
	query.Set("page", "")
	url.RawQuery = query.Encode()
	reqURL := strings.Replace(url.String(), "page=", "page=%d", 1)
	URL := fmt.Sprintf("%s%s", os.Getenv("APP_URL"), reqURL)

	var nextPageURL, lastPageURL, firstPageURL, prevPageURL *string
	lastPage := int(math.Ceil(float64(count) / float64(limit)))
	from := offset + 1
	to := ((page - 1) * limit) + count
	nURL := fmt.Sprintf(URL, page+1)
	lURL := fmt.Sprintf(URL, lastPage)
	fURL := fmt.Sprintf(URL, 1)
	pURL := fmt.Sprintf(URL, page-1)
	nextPageURL = &nURL
	lastPageURL = &lURL
	firstPageURL = &fURL
	prevPageURL = &pURL

	if page == lastPage {
		nextPageURL = nil
	}
	if page > lastPage {
		nextPageURL = nil
		prevPageURL = lastPageURL
		from = 0
		to = 0
	}
	if page == 1 {
		prevPageURL = nil
	}
	if count == 0 {
		firstPageURL = nil
		nextPageURL = nil
		lastPageURL = nil
		prevPageURL = nil
		from = 0
		to = 0
	}

	defaultResponse.Set(count, limit, page, lastPage, firstPageURL, lastPageURL, nextPageURL, prevPageURL, from, to, datas)
	return defaultResponse
}
