package models

type Numbers struct {
	Token  string    `json:"token"`
	Values []float64 `json:"values"`
}

type SumResponse struct {
	Token string  `json:"token"`
	Sum   float64 `json:"sum"`
}

type MultiplyResponse struct {
	Token    string  `json:"token"`
	Multiply float64 `json:"multiply"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}
