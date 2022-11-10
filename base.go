package main

import "goken/pool"

type SubmitResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TakeResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Token   pool.Token `json:"token"`
}

var ERROR_RESPONSE = SubmitResponse{
	Code:    1,
	Message: "ERROR",
}

var SUCCESS_RESPONSE = SubmitResponse{
	Code:    0,
	Message: "SUCCESS",
}
