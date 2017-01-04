package reqman

import (
	"github.com/uber-go/zap"
	"net/http"
	"time"
)

type requestLogger struct {
	logger zap.Logger
}

func (r *requestLogger) ServeMiddleware(
	w http.ResponseWriter,
	req *http.Request,
	context *RequestContext,
	next NextFunc,
) {
	defer func(reqLog *requestLogger, ctx *RequestContext) {
		// using defer to log in case of panics
		reqLog.logger.Info("request made", ctx.LogFields...)
	}(r, context)
	next()
}

func CreateRequestLogger(logger zap.Logger) PipelineHandler {
	return &requestLogger{logger}
}

type requestTimer struct {
}

func (*requestTimer) ServeMiddleware(
	w http.ResponseWriter,
	req *http.Request,
	context *RequestContext,
	next NextFunc,
) {
	currentTime := time.Now().UnixNano()
	defer func(ctx *RequestContext, cTime int64) {
		ctx.AddField(zap.Int64("request-time", time.Now().UnixNano()-cTime))
	}(context, currentTime)
	next()
}

func CreateRequestTimer() PipelineHandler {
	return &requestTimer{}
}
