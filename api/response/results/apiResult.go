package results

type APIResult struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   *string     `json:"error"`
}
