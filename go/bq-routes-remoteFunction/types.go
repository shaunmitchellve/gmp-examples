package function

type bqRequest struct {
	RequestId          string            `json:"requestId"`
	Caller             string            `json:"caller"`
	SessionUser        string            `json:"sessionUser"`
	UserDefinedContext map[string]string `json:"userDefinedContext"`
	Calls              [][]interface{}   `json:"calls"`
}

type bqResponse struct {
	Replies      []int64	`json:"replies,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

type bqResponseJSON struct {
	Replies      []route	`json:"replies,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

type route struct {
	Distance int64 `json:"distance"`
	Duration int64 `json:"duration"`
}