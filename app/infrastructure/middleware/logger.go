package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"

	"jade-mes/app/infrastructure/constant"
	"jade-mes/app/infrastructure/log"
	"jade-mes/app/infrastructure/util"
)

// Logger ...
func Logger(pathBlacklist ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(pathBlacklist); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range pathBlacklist {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.Query()
		body, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // write back, otherwise no data will be received afterward

		spanID, _ := util.NewUUID()
		c.Set("SpanId", spanID)

		c.Next()

		if _, ok := skip[path]; !ok {
			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					log.Error(e)
				}
			} else {
				end := time.Now()
				clientIP := c.ClientIP()
				method := c.Request.Method
				statusCode := c.Writer.Status()
				comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
				latency := end.Sub(start)

				log.AccessLog(
					fmt.Sprintf("%3d | %13v | %-15s | %-7s | %s | %s",
						statusCode,
						latency,
						clientIP,
						method,
						path,
						comment,
					),
					log.Reflect("accessPayload", map[string]interface{}{
						"statusCode": statusCode,
						"latency":    latency,
						"clientIP":   clientIP,
						"path":       path,
					}),
					log.String("spanId", spanID),
					log.String("category", constant.LogCategoryAccess),
					log.Reflect("query", query),
					log.String("rawBody", string(body)),
					log.String("user-agent", c.Request.UserAgent()),
				)
			}
		}
	}
}
