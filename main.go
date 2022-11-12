package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goken/pool"
	"io"
	"net/http"
	"os"
	"strconv"
)

var tokenPool = pool.New("./token")

func main() {

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	// Logging to a file.
	f, _ := os.OpenFile("./gin.log", os.O_CREATE|os.O_APPEND, 0666)
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	router.POST("/ali-token/submit", func(context *gin.Context) {
		data, err := io.ReadAll(context.Request.Body)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"code": "1", "message": "ERROR"})
		}

		var token pool.Token

		err = json.Unmarshal(data, &token)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"code": "1", "message": "ERROR"})
		}

		tokenPool.Offer(token)
		context.JSON(http.StatusOK, gin.H{"code": "0", "message": "SUCCESS"})
	})

	router.GET("/ali-token/take", func(ctx *gin.Context) {
		t := ctx.DefaultQuery("t", "0")
		timestamp, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			timestamp = 0
		}

		token := tokenPool.Get(timestamp)
		if token != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": "0", "message": "SUCCESS", "result": *token})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": "1", "message": "NO ITEM"})
		}

	})

	router.Run("localhost:38080")

}
