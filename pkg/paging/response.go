package paging

//ResponseInterface for hook response
type ResponseInterface interface {
	Set(total, perPage, curPage, lastPage int, fpURL, lpURL, npURL, ppURL *string, from, to int, datas interface{})
}

//Response type for pagination
type Response struct {
	Total        int         `json:"total"`
	PerPage      int         `json:"perPage"`
	CurrentPage  int         `json:"currentPage"`
	LastPage     int         `json:"lastPage"`
	FirstPageURL *string     `json:"firstPageUrl"`
	LastPageURL  *string     `json:"lastPageUrl"`
	NextPageURL  *string     `json:"nextPageUrl"`
	PrevPageURL  *string     `json:"prevPageUrl"`
	From         int         `json:"from"`
	To           int         `json:"to"`
	Data         interface{} `json:"data"`
}

//Set implement interface
func (r *Response) Set(total, perPage, curPage, lastPage int, fpURL, lpURL, npURL, ppURL *string, from, to int, datas interface{}) {
	r.Total = total
	r.PerPage = perPage
	r.CurrentPage = curPage
	r.LastPage = lastPage
	r.FirstPageURL = fpURL
	r.LastPageURL = lpURL
	r.NextPageURL = npURL
	r.PrevPageURL = ppURL
	r.From = from
	r.To = to
	r.Data = datas
}
