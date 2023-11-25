package domain

// CustomResponse struct defines the structure for custom API responses.
// It includes the HTTP status code and a flexible data field for response content.
type CustomResponse struct {
	HTTPStatusCode int `json:"http_status_code"`
	ResponseData   any `json:"response_data"`
}
