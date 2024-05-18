package model

type RequestBody struct {
	Data      string `json:"data"`
	Signature []byte `json:"signature"`
}
