package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SendResponse(ctx *gin.Context, statusCode int, resMsg string, scriptLogs []string, logMsg string) {
	ctx.JSON(statusCode, gin.H{
		"result": resMsg,
		"logs":   scriptLogs,
	})
	log.Println(logMsg)
}
