package turingpi

type SDCard struct {
	Total int64 `json:"total"`
	Free  int64 `json:"free"`
	Use   int64 `json:"use"`
}

type SDCardResponse struct {
	Response []SDCard `json:"response"`
}
