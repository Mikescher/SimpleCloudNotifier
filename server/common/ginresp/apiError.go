package ginresp

type sendAPIError struct {
	Success        bool   `json:"success"`
	Error          int    `json:"error"`
	ErrorHighlight int    `json:"errhighlight"`
	Message        string `json:"message"`
}

type internAPIError struct {
	Success bool   `json:"success"`
	ErrorID int    `json:"errid,omitempty"`
	Message string `json:"message"`
}
