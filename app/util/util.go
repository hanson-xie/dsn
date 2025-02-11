package util

import (
	"github.com/Bedrock-Technology/Dsn/app/proto"
	"github.com/Bedrock-Technology/Dsn/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorMsg(c *gin.Context, err string) {
	log.Errorf("errorMsg", "error", err)
	c.AbortWithStatusJSON(http.StatusBadRequest,
		proto.ResponseMsg{
			Error: err,
		})
}

func SuccessMsg(c *gin.Context, code int, msg string, data interface{}) {
	log.Debugf("successMsg", "code", code, "msg", msg, "data", data)
	c.AbortWithStatusJSON(http.StatusOK, proto.ResponseSuccessMsg{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
