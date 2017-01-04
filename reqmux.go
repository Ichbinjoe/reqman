package reqman

import (
	"net/http"
	"strings"
)

type RequestMuxNode interface {
	ServeCtxHTTP(string, http.ResponseWriter,
		*http.Request, *RequestContext) bool
}

type requestMuxEndpoint struct {
	endpoint CtxHttpHandler
	exact    bool
}

func (r *requestMuxEndpoint) ServeCtxHTTP(
	remainingRoute string,
	resp http.ResponseWriter,
	req *http.Request,
	ctx *RequestContext,
) bool {

	if !r.exact || len(remainingRoute) == 0 {
		r.endpoint.ServeCtxHTTP(resp, req, ctx)
		return true
	}
	return false
}

type requestTreeNode struct {
	routingMap         map[string]RequestMuxNode
	defaultNode        *RequestMuxNode
	defaultPlaceholder *string
	capNode            *CtxHttpHandler
}

func (r *requestTreeNode) serveDefault(
	currentRoute, remainingRoute string,
	resp http.ResponseWriter,
	req *http.Request,
	ctx *RequestContext,
) bool {
	if r.defaultNode != nil {
		existed := false
		var oldRoutingParam string
		if r.defaultPlaceholder != nil {
			oldRoutingParam, existed = ctx.RoutingParams[*r.defaultPlaceholder]
			ctx.RoutingParams[*r.defaultPlaceholder] = currentRoute
		}
		handled := (*r.defaultNode).ServeCtxHTTP(
			remainingRoute, resp, req, ctx)
		if !handled {
			if existed {
				ctx.RoutingParams[*r.defaultPlaceholder] = oldRoutingParam
			} else {
				delete(ctx.RoutingParams, *r.defaultPlaceholder)
			}
		}
		return handled
	}
	return false
}

func (r *requestTreeNode) ServeCtxHTTP(
	remainingRoute string,
	resp http.ResponseWriter,
	req *http.Request,
	ctx *RequestContext,
) bool {
	if len(remainingRoute) == 0 {
		if r.capNode != nil {
			(*r.capNode).ServeCtxHTTP(resp, req, ctx)
			return true
		}
	}

	seperatorIdx := strings.IndexRune(remainingRoute, '/')
	if seperatorIdx == -1 {
		muxNode, exists := r.routingMap[remainingRoute]
		if exists {
			return muxNode.ServeCtxHTTP("", resp, req, ctx)
		}
		return r.serveDefault(remainingRoute, "", resp, req, ctx)
	} else {
		currentRoute := remainingRoute[:seperatorIdx+1]
		remainingRoute := remainingRoute[seperatorIdx+1:]
		muxNode, exists := r.routingMap[currentRoute]
		if !exists {
			// trim off the trailing /, and try to match again.
			currentRoute = currentRoute[:len(currentRoute)-1]
			muxNode, exists = r.routingMap[currentRoute]
		}
		if exists {
			return muxNode.ServeCtxHTTP(remainingRoute, resp, req, ctx)
		}
		return r.serveDefault(currentRoute, remainingRoute, resp, req, ctx)
	}
}

type RequestMux struct {
	treeRoot requestTreeNode
}

func (r *RequestMux) ServeMiddleware(
	w http.ResponseWriter,
	req *http.Request,
	ctx *RequestContext,
	fallback NextFunc,
) {
	found := r.treeRoot.ServeCtxHTTP(req.URL.Path[1:], w, req, ctx)
	if !found {
		fallback()
	}
}

func (r *RequestMux) ServeCtxHTTP(
	w http.ResponseWriter,
	req *http.Request,
	ctx *RequestContext,
) {
	r.treeRoot.ServeCtxHTTP(req.URL.Path[1:], w, req, ctx)
}

func (r *RequestMux) ServeHTTP(
	w http.ResponseWriter,
	req *http.Request,
) {
	ctx := CreateNewHTTPContext()
	r.treeRoot.ServeCtxHTTP(req.URL.Path[1:], w, req, ctx)
}
