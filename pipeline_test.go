package reqman

import (
	"net/http"
	"testing"
)

type testPipeline struct {
	passthrough bool
	hit         bool
}

func (p *testPipeline) ServeMiddleware(
	w http.ResponseWriter,
	req *http.Request,
	ctx *RequestContext,
	f NextFunc,
) {
	if p.hit {
		panic("PipelineHandler called twice")
	}

	p.hit = true

	if p.passthrough {
		f()
	}
}

func TestPipelineReqThrough(t *testing.T) {

	respWriter := new(http.ResponseWriter)
	req := http.Request{}

	middleware := testPipeline{passthrough: true}
	endCap := testPipeline{passthrough: false}
	overstepTrigger := testPipeline{passthrough: false}

	pipeline := CreatePipeline(&middleware, &endCap, &overstepTrigger)

	pipeline.ServeHTTP(*respWriter, &req)

	if !middleware.hit {
		t.Error("middleware never hit")
	}

	if !endCap.hit {
		t.Error("terminal node never hit")
	}

	if overstepTrigger.hit {
		t.Error("pipeline overstepped")
	}
}
