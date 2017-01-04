package reqman

import (
	"net/http"
)

type hittable interface {
	wasHit() bool
}

type testMuxHandler struct {
	hit, returnVal bool
	remainingRoute string
}

func (t *testMuxHandler) ServeCtxHTTP(
	route string,
	_ http.ResponseWriter,
	_ *http.Request,
	_ *RequestContext,
) bool {
	if t.hit {
		panic("Called ctx handler twice!")
	}

	t.hit = true
	t.remainingRoute = route
	return t.returnVal
}

func (t *testMuxHandler) wasHit() bool {
	return t.hit
}

type testCtxHandler struct {
	hit bool
}

func (t *testCtxHandler) ServeCtxHTTP(
	http.ResponseWriter, *http.Request, *RequestContext) {
	if t.hit {
		panic("Called ctx handler twice")
	}

	t.hit = true
}

func (t *testCtxHandler) wasHit() bool {
	return t.hit
}

func ValidateNotHit(vals ...hittable) int {
	for i, v := range vals {
		if v.wasHit() {
			return i
		}
	}
	return -1
}

func l
