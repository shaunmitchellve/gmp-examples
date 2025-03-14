package function

type BqRequest struct {
	RequestId          string            `json:"requestId"`
	Caller             string            `json:"caller"`
	SessionUser        string            `json:"sessionUser"`
	UserDefinedContext map[string]string `json:"userDefinedContext"`
	Calls              [][]interface{}   `json:"calls"`
}

type BqResponse struct {
	Replies      []int64`json:"replies,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}