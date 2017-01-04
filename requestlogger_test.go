package reqman

import (
	"net/http"
	"testing"
)

type reqTrigger struct {
	hit, panic bool
}

func (r *reqTrigger) Next() {
	if r.hit {
		panic("reqTrigger hit twice")
	}
	r.hit = true
	if r.panic {
		panic("controlled panic")
	}
}

func TestTimerMiddleware(t *testing.T) {
	respWriter := new(http.ResponseWriter)
	ctx := CreateNewHTTPContext()

	trigger := reqTrigger{false, false}

	requestMiddleware := CreateRequestTimer()

	requestMiddleware.ServeMiddleware(*respWriter, nil, ctx, trigger.Next)

	if !trigger.hit {
		t.Fail()
	}

	if len(ctx.LogFields) != 1 {
		t.Fail()
	}
}

func TestPanicingTimerMiddleware(t *testing.T) {
	respWriter := new(http.ResponseWriter)
	ctx := CreateNewHTTPContext()

	trigger := reqTrigger{false, true}

	requestMiddleware := CreateRequestTimer()
	defer func() {
		if r := recover(); r != nil {
			if r.(string) != "controlled panic" {
				t.Fail()
			}
		}
		if !trigger.hit {
			t.Fail()
		}

		if len(ctx.LogFields) != 1 {
			t.Fail()
		}
	}()
	requestMiddleware.ServeMiddleware(*respWriter, nil, ctx, trigger.Next)
}
