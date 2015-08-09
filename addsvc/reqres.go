package main

type AddRequest struct {
	A int64 `json:"a"`
	B int64 `json:"b"`
}

type AddResponse struct {
	V int64 `json:"v"`
}
