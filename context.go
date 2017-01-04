package reqman

import (
	"github.com/uber-go/zap"
	"sync"
)

type ContextKey int

var currentKey ContextKey
var currentKeyLock sync.Mutex

func GetCtxKey() ContextKey {
	currentKeyLock.Lock()
	defer currentKeyLock.Unlock()
	currentKey = currentKey + 1
	return currentKey
}

type RequestContext struct {
	LogFields     []zap.Field
	RoutingParams map[string]string
}

func (r *RequestContext) AddField(fields ...zap.Field) {
	r.LogFields = append(r.LogFields, fields...)
}

func CreateNewHTTPContext() *RequestContext {
	return &RequestContext{
		make([]zap.Field, 0, 10),
		make(map[string]string),
	}
}
