package log

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mercor/payment-service/constants"
	"github.com/mercor/payment-service/pkg/env"
)

type MiddlewareOptions struct {
	Format      string
	Level       string
	LogRequest  bool
	LogResponse bool
	LogHeader   bool
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogMiddleware(opts MiddlewareOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Ignore Health Requests
		if path == "/health" {
			c.Next()

			return
		}

		query := c.Request.URL.RawQuery
		reqID := env.GetRequestID(c)
		start := time.Now()
		requestBodyString := "<Disabled>"
		bodyWriter := &responseWriter{
			body: bytes.NewBufferString("<Disabled>"),
		}

		l := DefaultLogger()
		l = l.With(
			String(constants.HeaderXMercorRequestID, reqID),
		)

		ctx := ContextWithLogger(c, l).(*gin.Context)

		// Create a custom ResponseWriter to capture the response body
		if opts.LogResponse {
			bodyWriter.body = bytes.NewBufferString("")
			bodyWriter.ResponseWriter = c.Writer
			ctx.Writer = bodyWriter
		}

		if opts.LogRequest {
			requestBody, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			requestBodyString = string(requestBody)
		}

		defer func() {
			// Capture the request complete timestamp
			end := time.Now()

			logFields := []Field{
				String("Amzn-Trace-ID", ctx.GetHeader("x-amzn-trace-id")),
				Int("Status", ctx.Writer.Status()),
				String("Method", ctx.Request.Method),
				String("Domain", ctx.Request.Host),
				String("Path", path),
				String("RequestBody", requestBodyString),
				String("Query", query),
				String("ResponseBody", bodyWriter.body.String()),
				String("IP", ctx.ClientIP()),
				String("User-Agent", ctx.Request.UserAgent()),
				Duration("Latency", time.Since(start)),
				String("RequestReceivedAt", start.Format(time.RFC3339)),
				String("RequestCompletedAt", end.Format(time.RFC3339)),
			}

			if opts.LogHeader {
				for k, v := range ctx.Request.Header {
					if len(v) > 0 {
						logFields = append(logFields, String(k, v[0]))
					}
				}
			}

			l.With(logFields...).Info(path)
		}()

		ctx.Next()
	}
}
