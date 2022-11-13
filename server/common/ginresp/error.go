package ginresp

type errBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
