package reqman

import (
	"net/http"
)

type NextFunc func()

type PipelineHandler interface {
	ServeMiddleware(http.ResponseWriter, *http.Request,
		*RequestContext, NextFunc)
}

type CtxHttpHandler interface {
	ServeCtxHTTP(http.ResponseWriter, *http.Request, *RequestContext)
}

type httpPipelineHandler struct {
	handler http.Handler
}

func (p *httpPipelineHandler) ServeMiddleware(
	resp http.ResponseWriter,
	req *http.Request,
	_ *RequestContext,
	_ NextFunc,
) {
	p.handler.ServeHTTP(resp, req)
}

func WrapHttpHandler(handler http.Handler) PipelineHandler {
	return &httpPipelineHandler{handler}
}

type Pipeline struct {
	handlers []PipelineHandler
}

type pipelineContext struct {
	handlerIndex int
	w            http.ResponseWriter
	req          *http.Request
	context      *RequestContext
}

func (p *Pipeline) getNext(context pipelineContext) func() {
	context.handlerIndex++
	if context.handlerIndex >= len(p.handlers) {
		return func() {} // end of the pipeline
	}

	return func() {
		p.handlers[context.handlerIndex].ServeMiddleware(
			context.w, context.req, context.context, p.getNext(context))
	}
}

func (p *Pipeline) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.ServeContextualHTTP(w, req, CreateNewHTTPContext())
}

func (p *Pipeline) ServeContextualHTTP(
	w http.ResponseWriter,
	req *http.Request,
	ctx *RequestContext,
) {
	context := pipelineContext{0, w, req, ctx}
	p.handlers[0].ServeMiddleware(w, req, context.context, p.getNext(context))
}

func CreatePipeline(pipelineHandlers ...PipelineHandler) http.Handler {
	if len(pipelineHandlers) < 1 {
		panic("no pipeline handlers passed")
	}

	return &Pipeline{pipelineHandlers}
}
