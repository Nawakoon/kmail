package model

type RequestBody struct {
	Data      string `json:"data"`
	Signature []byte `json:"signature"`
}

type QueryJson struct {
	Page  *int `json:"page"`
	Limit *int `json:"limit"`
}
