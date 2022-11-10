package main

import (
	"goken/pool"
	"net/http"
)

func main() {

	tokenPool := pool.New("./token")

	for i := 0; i < 5; i++ {
		tokenPool.Offer(pool.Token{
			Sig: "lee",
		})

	}

	http.ListenAndServe(":8000", nil)
}
