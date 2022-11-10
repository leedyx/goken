package main

import (
	"encoding/json"
	"goken/pool"
	"net/http"
)

var sugarLogger = pool.SugarLogger
var tokenPool = pool.New("./token")

func SubmitHandler(w http.ResponseWriter, req *http.Request) {
	data := make([]byte, 4096)

	defer req.Body.Close()
	res, err := req.Body.Read(data)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json; charset=UTF-8")

	var token pool.Token
	err = json.Unmarshal(data[0:res], &token)
	if err != nil {
		sugarLogger.Error(err)
		errRes, _ := json.Marshal(ERROR_RESPONSE)
		w.Write(errRes)
	}

	tokenPool.Offer(token)
	successRes, _ := json.Marshal(SUCCESS_RESPONSE)
	w.Write(successRes)

}

func main() {

	http.HandleFunc("/ali-token/submit", SubmitHandler)

	http.ListenAndServe("localhost:38080", nil)

}
