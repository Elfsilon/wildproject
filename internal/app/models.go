package app

type Response struct {
	Result interface{} `json:"result"`
	Error  error       `json:"error,omitempty"`
}
