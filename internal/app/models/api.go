package m

type Response struct {
	Result interface{} `json:"result"`
	Error  error       `json:"error,omitempty"`
}
