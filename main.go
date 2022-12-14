package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goken/pool"
	"io"
	"net/http"
	"os"
	"strconv"
)

var tokenPool = pool.New("./token")

func main() {

	// Logging to a file.
	pwd, _ := os.Getwd()

	logName := fmt.Sprintf("%s/%s", pwd, "log/gin.log")
	f, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		fmt.Errorf("error ! %v", err)
	}

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

	router.Run(":38080")

}
