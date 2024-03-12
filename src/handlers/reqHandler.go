package handlers

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"shellProxy/data_defs"
	"strings"

	"github.com/gin-gonic/gin"
)

func replaceDoubleDots(oriString string) string {
	dirKeys := strings.Split(oriString, "/")

	for i, key := range dirKeys {
		if key == ".." {
			dirKeys[i] = "X"
		}
	}

	return strings.Join(dirKeys, "/")
}

func sanitizeReq(reqParams *data_defs.ReqParams, ctx *gin.Context) bool {
	pattSlash := regexp.MustCompile(`\\`)
	pattEndSlash := regexp.MustCompile(`^\/*|\/*$`)
	pattMidSlash := regexp.MustCompile(`\/+`)

	sanitizedStr, unescapeErr := url.QueryUnescape(reqParams.ShellName)
	if unescapeErr != nil {
		SendResponse(
			ctx,
			500,
			"Internal server error",
			[]string{},
			fmt.Sprintf("%s - (sanitizeReq) ERR: Fail to decode the string", ctx.ClientIP()))
		return false
	}

	sanitizedStr = pattSlash.ReplaceAllString(sanitizedStr, "/")
	sanitizedStr = pattEndSlash.ReplaceAllString(sanitizedStr, "")
	sanitizedStr = pattMidSlash.ReplaceAllString(sanitizedStr, "/")
	sanitizedStr = replaceDoubleDots(sanitizedStr)

	var scriptExt string
	if runtime.GOOS == "windows" {
		scriptExt = ".ps1"
	} else {
		scriptExt = ".sh"
	}

	reqParams.ShellName = fmt.Sprintf("./%s%s", sanitizedStr, scriptExt)
	return true
}

func ExamReq(reqParams *data_defs.ReqParams, ctx *gin.Context) bool {
	if err := ctx.BindQuery(reqParams); err != nil || reqParams.ShellName == "" {
		SendResponse(
			ctx,
			400,
			"Invalid request for shell execution",
			[]string{},
			fmt.Sprintf("%s - (examConfig) ERR: Invalid params", ctx.ClientIP()))
		return false
	}

	isSuccess := sanitizeReq(reqParams, ctx)
	if isSuccess {
		if _, err := os.Stat(reqParams.ShellName); err != nil {
			SendResponse(
				ctx,
				404,
				"Fail to access the shell script",
				[]string{},
				fmt.Sprintf("%s - (examConfig) ERR: %s", ctx.ClientIP(), err))
			return false
		}
	}
	return isSuccess
}
