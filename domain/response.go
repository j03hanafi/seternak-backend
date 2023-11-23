package domain

type CustomResponse struct {
	HTTPStatusCode int `json:"http_status_code"`
	ResponseData   any `json:"response_data"`
}
