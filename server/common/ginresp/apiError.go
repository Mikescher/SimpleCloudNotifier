package ginresp

type apiError struct {
	Success        bool    `json:"success"`
	Error          int     `json:"error"`
	ErrorHighlight int     `json:"errhighlight"`
	Message        string  `json:"message"`
	RawError       *string `json:"errorObj,omitempty"`
	Trace          string  `json:"traceObj,omitempty"`
}

type compatAPIError struct {
	Success bool   `json:"success"`
	ErrorID int    `json:"errid,omitempty"`
	Message string `json:"message"`
}
