package controller

import (
	"net/http"
	"jade-mes/app/infrastructure/constant"
	"jade-mes/app/infrastructure/log"
	"jade-mes/ecode"

	"github.com/gin-gonic/gin"
)

type responseBody struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
	Data    interface{} `json:"data"`
}

func handleError(err error, target *responseBody) {
	if err != nil {
		var e ecode.ECode
		var ok bool

		if e, ok = err.(ecode.ECode); !ok {
			e = ecode.ErrUnhandledException(err)
		}

		target.Code = e.Code()
		target.Message = e.Message()
		target.Error = e.Errors()
	}
}

func response(c *gin.Context, data interface{}, err error) {
	var resp responseBody
	defer func() {
		c.JSON(http.StatusOK, &resp)
	}()

	if err != nil {
		spanID, _ := c.Get("SpanId")

		log.Error(
			"app exception",
			log.Err(err),
			log.Reflect("spanId", spanID),
			log.String("category", constant.LogCategoryError),
		)

		handleError(err, &resp)
		return
	}

	resp.Data = data
}
